/*
Copyright 2021 Upbound Inc.
*/

package main

import (
	"context"
	"flag"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/crossplane/crossplane-runtime/pkg/certificates"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/statemetrics"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	xpcontroller "github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	tjcontroller "github.com/crossplane/upjet/pkg/controller"
	"gopkg.in/alecthomas/kingpin.v2"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/dana-team/provider-avi/apis"
	"github.com/dana-team/provider-avi/apis/v1alpha1"
	"github.com/dana-team/provider-avi/config"
	"github.com/dana-team/provider-avi/internal/clients"
	"github.com/dana-team/provider-avi/internal/controller"
	"github.com/dana-team/provider-avi/internal/features"
)

func main() {
	var (
		debug                      bool
		syncPeriod                 time.Duration
		pollInterval               time.Duration
		pollStateMetricInterval    time.Duration
		leaderElection             bool
		maxReconcileRate           int
		namespace                  string
		enableExternalSecretStores bool
		essTLSCertsPath            string
		enableManagementPolicies   bool
	)

	flag.BoolVar(&debug, "debug", false, "Run with debug logging.")
	flag.BoolVar(&debug, "d", false, "Run with debug logging (shorthand).")
	flag.DurationVar(&syncPeriod, "sync", time.Hour, "Controller manager sync period such as 300ms, 1.5h, or 2h45m")
	flag.DurationVar(&pollInterval, "poll", 10*time.Minute, "Poll interval controls how often an individual resource should be checked for drift.")
	flag.DurationVar(&pollStateMetricInterval, "poll-state-metric", 5*time.Second, "State metric recording interval")
	flag.BoolVar(&leaderElection, "leader-election", false, "Use leader election for the controller manager.")
	flag.IntVar(&maxReconcileRate, "max-reconcile-rate", 10, "The global maximum rate per second at which resources may be checked for drift from the desired state.")
	flag.StringVar(&namespace, "namespace", "crossplane-system", "Namespace used to set as default scope in default secret store config.")
	flag.BoolVar(&enableExternalSecretStores, "enable-external-secret-stores", false, "Enable support for ExternalSecretStores.")
	flag.StringVar(&essTLSCertsPath, "ess-tls-cert-dir", "", "Path of ESS TLS certificates.")
	flag.BoolVar(&enableManagementPolicies, "enable-management-policies", true, "Enable support for Management Policies.")

	// Parse the command-line flags
	flag.Parse()

	log.Default().SetOutput(io.Discard)
	ctrl.SetLogger(zap.New(zap.WriteTo(io.Discard)))

	zl := zap.New(zap.UseDevMode(debug))
	logr := logging.NewLogrLogger(zl.WithName("provider-avi"))
	if debug {
		// The controller-runtime runs with a no-op logger by default. It is
		// *very* verbose even at info level, so we only provide it a real
		// logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	}

	// currently, we configure the jitter to be the 5% of the poll interval
	pollJitter := time.Duration(float64(pollInterval) * 0.05)
	logr.Debug("Starting", "sync-period", syncPeriod.String(),
		"poll-interval", pollInterval.String(), "poll-jitter", pollJitter, "max-reconcile-rate", maxReconcileRate)

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		LeaderElection:   leaderElection,
		LeaderElectionID: "crossplane-leader-election-provider-avi",
		Cache: cache.Options{
			SyncPeriod: &syncPeriod,
		},
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")
	kingpin.FatalIfError(apis.AddToScheme(mgr.GetScheme()), "Cannot add Avi APIs to scheme")

	metricRecorder := managed.NewMRMetricRecorder()
	stateMetrics := statemetrics.NewMRStateMetrics()

	metrics.Registry.MustRegister(metricRecorder)
	metrics.Registry.MustRegister(stateMetrics)

	provider := config.GetProvider()
	o := tjcontroller.Options{
		Options: xpcontroller.Options{
			Logger:                  logr,
			GlobalRateLimiter:       ratelimiter.NewGlobal(maxReconcileRate),
			PollInterval:            pollInterval,
			MaxConcurrentReconciles: maxReconcileRate,
			Features:                &feature.Flags{},
			MetricOptions: &xpcontroller.MetricOptions{
				PollStateMetricInterval: pollStateMetricInterval,
				MRMetrics:               metricRecorder,
				MRStateMetrics:          stateMetrics,
			},
		},
		Provider:              provider,
		SetupFn:               clients.TerraformSetupBuilder(provider.TerraformProvider),
		PollJitter:            pollJitter,
		OperationTrackerStore: tjcontroller.NewOperationStore(logr),
	}

	if enableManagementPolicies {
		o.Features.Enable(features.EnableBetaManagementPolicies)
		logr.Info("Beta feature enabled", "flag", features.EnableBetaManagementPolicies)
	}

	if enableExternalSecretStores {
		o.SecretStoreConfigGVK = &v1alpha1.StoreConfigGroupVersionKind
		logr.Info("Alpha feature enabled", "flag", features.EnableAlphaExternalSecretStores)

		o.ESSOptions = &tjcontroller.ESSOptions{}
		if essTLSCertsPath != "" {
			logr.Info("ESS TLS certificates path is set. Loading mTLS configuration.")
			tCfg, err := certificates.LoadMTLSConfig(filepath.Join(essTLSCertsPath, "ca.crt"), filepath.Join(essTLSCertsPath, "tls.crt"), filepath.Join(essTLSCertsPath, "tls.key"), false)
			kingpin.FatalIfError(err, "Cannot load ESS TLS config.")

			o.ESSOptions.TLSConfig = tCfg
		}

		// Ensure default store config exists.
		kingpin.FatalIfError(resource.Ignore(kerrors.IsAlreadyExists, mgr.GetClient().Create(context.Background(), &v1alpha1.StoreConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
			Spec: v1alpha1.StoreConfigSpec{
				SecretStoreConfig: xpv1.SecretStoreConfig{
					DefaultScope: namespace,
				},
			},
		})), "cannot create default store config")
	}

	kingpin.FatalIfError(controller.Setup(mgr, o), "Cannot setup Avi controllers")
	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}
