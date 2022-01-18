package http

import (
	"net/http"

	"git.xswitch.cn/xswitch/xctrl/stack/registry"
	"git.xswitch.cn/xswitch/xctrl/stack/selector"
)

func NewRoundTripper(opts ...Option) http.RoundTripper {
	options := Options{
		Registry: registry.DefaultRegistry,
	}
	for _, o := range opts {
		o(&options)
	}

	return &roundTripper{
		rt:   http.DefaultTransport,
		st:   selector.Random,
		opts: options,
	}
}
