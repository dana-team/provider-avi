package avi

import "github.com/crossplane/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("avi_pool", func(r *config.Resource) {
		r.ShortGroup = "avi"
		r.Kind = "Pool"
		r.Version = "v1alpha1"
	})
}
