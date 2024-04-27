/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	_ "embed"

	ujconfig "github.com/crossplane/upjet/pkg/config"
	"github.com/crossplane/upjet/pkg/registry/reference"
	"github.com/dana-team/provider-avi/config/avi"
	aviconfig "github.com/vmware/terraform-provider-avi/avi"
)

const (
	resourcePrefix = "avi"
	modulePath     = "github.com/dana-team/provider-avi"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("crossplane.io"),
		ujconfig.WithShortName("avi"),
		ujconfig.WithIncludeList(resourceList(cliReconciledExternalNameConfigs)),
		ujconfig.WithTerraformPluginSDKIncludeList(resourceList(terraformPluginSDKExternalNameConfigs)),
		ujconfig.WithDefaultResourceOptions(
			resourceConfigurator(),
		),
		ujconfig.WithReferenceInjectors([]ujconfig.ReferenceInjector{reference.NewInjector(modulePath)}),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithTerraformProvider(aviconfig.Provider()),
	)

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		avi.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// resourceList returns the list of resources that have external
// name configured in the specified table.
func resourceList(t map[string]ujconfig.ExternalName) []string {
	l := make([]string, len(t))
	i := 0
	for n := range t {
		// Expected format is regex and we'd like to have exact matches.
		l[i] = n + "$"
		i++
	}
	return l
}
