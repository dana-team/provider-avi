package avi

import "github.com/crossplane/upjet/pkg/config"

const (
	apiVersion = "v1alpha1"
	shortGroup = "avi"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("avi_pool", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "Pool"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("avi_gslb", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "GSLB"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("avi_serviceengine", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ServiceEngine"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("avi_serviceenginegroup", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ServiceEngineGroup"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("avi_virtualservice", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "VirtualService"
		r.Version = apiVersion
	})
	p.AddResourceConfigurator("avi_healthmonitor", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "HealthMonitor"
		r.Version = apiVersion
	})
	p.AddResourceConfigurator("avi_vsvip", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "VsVip"
		r.Version = apiVersion
	})

}
