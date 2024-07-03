package avi

import "github.com/crossplane/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("avi_pool", func(r *config.Resource) {
		r.ShortGroup = "avi"
		r.Kind = "Pool"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("avi_gslb", func(r *config.Resource) {
		r.ShortGroup = "avi"
		r.Kind = "GSLB"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("avi_serviceengine", func(r *config.Resource) {
		r.ShortGroup = "avi"
		r.Kind = "ServiceEngine"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("avi_serviceenginegroup", func(r *config.Resource) {
		r.ShortGroup = "avi"
		r.Kind = "ServiceEngineGroup"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("avi_virtualservice", func(r *config.Resource) {
		r.ShortGroup = "avi"
		r.Kind = "VirtualService"
		r.Version = "v1alpha1"
	})
}
