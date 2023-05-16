// Package registry provides a dynamic api service router
package registry

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"git.xswitch.cn/xswitch/xctrl/xctrl/api"
	"git.xswitch.cn/xswitch/xctrl/xctrl/api/router"
	"git.xswitch.cn/xswitch/xctrl/xctrl/api/router/util"
	"git.xswitch.cn/xswitch/xctrl/xctrl/logger"
	"git.xswitch.cn/xswitch/xctrl/xctrl/metadata"
)

// endpoint struct, that holds compiled pcre
type endpoint struct {
	hostregs []*regexp.Regexp
	pathregs []util.Pattern
	pcreregs []*regexp.Regexp
}

// router is the default router
type registryRouter struct {
	exit chan bool
	opts router.Options

	// registry cache
	//rc cache.Cache

	sync.RWMutex
	eps map[string]*api.Service
	// compiled regexp for host and path
	ceps map[string]*endpoint
}

func (r *registryRouter) isClosed() bool {
	select {
	case <-r.exit:
		return true
	default:
		return false
	}
}

func (r *registryRouter) Options() router.Options {
	return r.opts
}

func (r *registryRouter) Close() error {
	select {
	case <-r.exit:
		return nil
	default:
		close(r.exit)
		//r.rc.Stop()
	}
	return nil
}

func (r *registryRouter) Register(ep *api.Endpoint) error {
	return nil
}

func (r *registryRouter) Deregister(ep *api.Endpoint) error {
	return nil
}

func (r *registryRouter) Endpoint(req *http.Request) (*api.Service, error) {
	if r.isClosed() {
		return nil, errors.New("router closed")
	}

	r.RLock()
	defer r.RUnlock()

	var idx int
	if len(req.URL.Path) > 0 && req.URL.Path != "/" {
		idx = 1
	}
	path := strings.Split(req.URL.Path[idx:], "/")

	// use the first match
	// TODO: weighted matching
	for n, e := range r.eps {
		cep, ok := r.ceps[n]
		if !ok {
			continue
		}
		ep := e.Endpoint
		var mMatch, hMatch, pMatch bool
		// 1. try method
		for _, m := range ep.Method {
			if m == req.Method {
				mMatch = true
				break
			}
		}
		if !mMatch {
			continue
		}
		if logger.V(logger.DebugLevel, logger.DefaultLogger) {
			logger.Debugf("api method match %s", req.Method)
		}

		// 2. try host
		if len(ep.Host) == 0 {
			hMatch = true
		} else {
			for idx, h := range ep.Host {
				if h == "" || h == "*" {
					hMatch = true
					break
				} else {
					if cep.hostregs[idx].MatchString(req.URL.Host) {
						hMatch = true
						break
					}
				}
			}
		}
		if !hMatch {
			continue
		}
		if logger.V(logger.DebugLevel, logger.DefaultLogger) {
			logger.Debugf("api host match %s", req.URL.Host)
		}

		// 3. try path via google.api path matching
		for _, pathreg := range cep.pathregs {
			matches, err := pathreg.Match(path, "")
			if err != nil {
				if logger.V(logger.DebugLevel, logger.DefaultLogger) {
					logger.Debugf("api gpath not match %s != %v", path, pathreg)
				}
				continue
			}
			if logger.V(logger.DebugLevel, logger.DefaultLogger) {
				logger.Debugf("api gpath match %s = %v", path, pathreg)
			}
			pMatch = true
			ctx := req.Context()
			md, ok := metadata.FromContext(ctx)
			if !ok {
				md = make(metadata.Metadata)
			}
			for k, v := range matches {
				md[fmt.Sprintf("x-api-field-%s", k)] = v
			}
			md["x-api-body"] = ep.Body
			*req = *req.Clone(metadata.NewContext(ctx, md))
			break
		}

		if !pMatch {
			// 4. try path via pcre path matching
			for _, pathreg := range cep.pcreregs {
				if !pathreg.MatchString(req.URL.Path) {
					if logger.V(logger.DebugLevel, logger.DefaultLogger) {
						logger.Debugf("api pcre path not match %s != %v", path, pathreg)
					}
					continue
				}
				if logger.V(logger.DebugLevel, logger.DefaultLogger) {
					logger.Debugf("api pcre path match %s != %v", path, pathreg)
				}
				pMatch = true
				break
			}
		}

		if !pMatch {
			continue
		}

		// TODO: Percentage traffic
		// we got here, so its a match
		return e, nil
	}

	// no match
	return nil, errors.New("not found")
}

func (r *registryRouter) Route(req *http.Request) (*api.Service, error) {
	if r.isClosed() {
		return nil, errors.New("router closed")
	}

	// try get an endpoint
	ep, err := r.Endpoint(req)
	if err == nil {
		return ep, nil
	}

	// error not nil
	// ignore that shit
	// TODO: don't ignore that shit

	// get the service name
	rp, err := r.opts.Resolver.Resolve(req)
	if err != nil {
		return nil, err
	}

	// service name
	name := rp.Name

	// get service
	//services, err := r.rc.GetService(name)
	if err != nil {
		return nil, err
	}

	// only use endpoint matching when the meta handler is set aka api.Default
	switch r.opts.Handler {
	// rpc handlers
	case "meta", "api", "rpc":
		handler := r.opts.Handler

		// set default handler to api
		if r.opts.Handler == "meta" {
			handler = "rpc"
		}

		// construct api service
		return &api.Service{
			Name: name,
			Endpoint: &api.Endpoint{
				Name:    rp.Method,
				Handler: handler,
			},
			Services: nil,
		}, nil
	// http handler
	case "http", "proxy", "web":
		// construct api service
		return &api.Service{
			Name: name,
			Endpoint: &api.Endpoint{
				Name:    req.URL.String(),
				Handler: r.opts.Handler,
				Host:    []string{req.Host},
				Method:  []string{req.Method},
				Path:    []string{req.URL.Path},
			},
			Services: nil,
		}, nil
	}

	return nil, errors.New("unknown handler")
}
