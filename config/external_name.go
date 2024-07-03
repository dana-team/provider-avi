/*
Copyright 2022 Upbound Inc.
*/

package config

import "github.com/crossplane/upjet/pkg/config"

// terraformPluginSDKExternalNameConfigs contains all external name configurations for this
// provider.
var terraformPluginSDKExternalNameConfigs = map[string]config.ExternalName{
	"avi_pool":               config.IdentifierFromProvider,
	"avi_gslb":               config.IdentifierFromProvider,
	"avi_serviceengine":      config.IdentifierFromProvider,
	"avi_serviceenginegroup": config.IdentifierFromProvider,
	"avi_virtualservice":     config.IdentifierFromProvider,
}

var terraformPluginFrameworkExternalNameConfigs = map[string]config.ExternalName{}

// cliReconciledExternalNameConfigs contains all external name configurations
// belonging to Terraform resources to be reconciled under the CLI-based
// architecture for this provider.
var cliReconciledExternalNameConfigs = map[string]config.ExternalName{}

// resourceConfigurator applies all external name configs listed in
// the table terraformPluginSDKExternalNameConfigs,
// cliReconciledExternalNameConfigs, and
// terraformPluginFrameworkExternalNameConfigs and sets the version of
// those resources to v1beta1.
func resourceConfigurator() config.ResourceOption {
	return func(r *config.Resource) {
		// If an external name is configured for multiple architectures,
		// Terraform Plugin Framework takes precedence over Terraform
		// Plugin SDKv2, which takes precedence over CLI architecture.
		e, configured := terraformPluginFrameworkExternalNameConfigs[r.Name]
		if !configured {
			e, configured = terraformPluginSDKExternalNameConfigs[r.Name]
			if !configured {
				e, configured = cliReconciledExternalNameConfigs[r.Name]
			}
		}
		if !configured {
			return
		}
		r.Version = "v1beta1"
		r.ExternalName = e
	}
}
