/*
Copyright 2022 Upbound Inc.
*/

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	gslb "github.com/dana-team/provider-avi/internal/controller/avi/gslb"
	pool "github.com/dana-team/provider-avi/internal/controller/avi/pool"
	serviceengine "github.com/dana-team/provider-avi/internal/controller/avi/serviceengine"
	serviceenginegroup "github.com/dana-team/provider-avi/internal/controller/avi/serviceenginegroup"
	virtualservice "github.com/dana-team/provider-avi/internal/controller/avi/virtualservice"
	providerconfig "github.com/dana-team/provider-avi/internal/controller/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		gslb.Setup,
		pool.Setup,
		serviceengine.Setup,
		serviceenginegroup.Setup,
		virtualservice.Setup,
		providerconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
