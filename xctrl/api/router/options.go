package router

import (
	"git.xswitch.cn/xswitch/xctrl/xctrl/api/resolver"
	"git.xswitch.cn/xswitch/xctrl/xctrl/api/resolver/vpath"
	//"git.xswitch.cn/xswitch/xctrl/xctrl/registry"
)

type Options struct {
	Handler  string
	Resolver resolver.Resolver
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Handler: "meta",
		//Registry: registry.NewRegistry(),
	}

	for _, o := range opts {
		o(&options)
	}

	if options.Resolver == nil {
		options.Resolver = vpath.NewResolver(
			resolver.WithHandler(options.Handler),
		)
	}

	return options
}

func WithHandler(h string) Option {
	return func(o *Options) {
		o.Handler = h
	}
}

func WithResolver(r resolver.Resolver) Option {
	return func(o *Options) {
		o.Resolver = r
	}
}
