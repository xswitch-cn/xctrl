package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.xswitch.cn/xswitch/xctrl/stack"

	"git.xswitch.cn/xswitch/xctrl/stack/broker"
	"git.xswitch.cn/xswitch/xctrl/stack/broker/nats"
	"git.xswitch.cn/xswitch/xctrl/stack/client"
	"git.xswitch.cn/xswitch/xctrl/stack/registry"
	"git.xswitch.cn/xswitch/xctrl/stack/registry/etcd"
	"git.xswitch.cn/xswitch/xctrl/stack/selector"
	"git.xswitch.cn/xswitch/xctrl/stack/server"
)

type serviceKey struct{}

// Service is an interface that wraps the lower level libraries
// within stack-rpc. Its a convenience method for building
// and initialising services.
type Service interface {
	// The service name
	Name() string
	// Init initialises options
	Init(...stack.Option)
	// Options returns the current options
	Options() stack.Options
	// Client is used to call services
	Client() client.Client
	// Server is for handling requests and events
	Server() server.Server
	// Run the service
	Run() error
	// The service implementation
	String() string
}

// Option set option callback
type Option func(*stack.Options)

func Broker(b broker.Broker) stack.Option {
	return func(o *stack.Options) {
		o.Broker = b
		// Update Client and Server
		o.Client.Init(client.Broker(b))
		o.Server.Init(server.Broker(b))
	}
}

func Client(c client.Client) stack.Option {
	return func(o *stack.Options) {
		o.Client = c
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) stack.Option {
	return func(o *stack.Options) {
		o.Context = ctx
	}
}

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) stack.Option {
	return func(o *stack.Options) {
		o.Signal = b
	}
}

func Server(s server.Server) stack.Option {
	return func(o *stack.Options) {
		o.Server = s
	}
}

// Registry sets the registry for the service
// and the underlying components
func Registry(r registry.Registry) stack.Option {
	return func(o *stack.Options) {
		o.Registry = r
		// Update Client and Server
		o.Client.Init(client.Registry(r))
		o.Server.Init(server.Registry(r))
		// Update Selector
		o.Client.Options().Selector.Init(selector.Registry(r))
		// Update Broker
		o.Broker.Init(broker.Registry(r))
	}
}

// Selector sets the selector for the service client
func Selector(s selector.Selector) stack.Option {
	return func(o *stack.Options) {
		o.Client.Init(client.Selector(s))
	}
}

// Address sets the address of the server
func Address(addr string) stack.Option {
	return func(o *stack.Options) {
		o.Server.Init(server.Address(addr))
	}
}

// Name of the service
func Name(n string) stack.Option {
	return func(o *stack.Options) {
		o.Server.Init(server.Name(n))
	}
}

// Version of the service
func Version(v string) stack.Option {
	return func(o *stack.Options) {
		o.Server.Init(server.Version(v))
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) stack.Option {
	return func(o *stack.Options) {
		o.Server.Init(server.Metadata(md))
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) stack.Option {
	return func(o *stack.Options) {
		o.Server.Init(server.RegisterTTL(t))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) stack.Option {
	return func(o *stack.Options) {
		o.Server.Init(server.RegisterInterval(t))
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) stack.Option {
	return func(o *stack.Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			o.Client = w[i-1](o.Client)
		}
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) stack.Option {
	return func(o *stack.Options) {
		o.Client.Init(client.WrapCall(w...))
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...server.HandlerWrapper) stack.Option {
	return func(o *stack.Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// Init once
		o.Server.Init(wrappers...)
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w ...server.SubscriberWrapper) stack.Option {
	return func(o *stack.Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}

		// Init once
		o.Server.Init(wrappers...)
	}
}

// Before and Afters

func BeforeStart(fn func() error) stack.Option {
	return func(o *stack.Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) stack.Option {
	return func(o *stack.Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func AfterStart(fn func() error) stack.Option {
	return func(o *stack.Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

func AfterStop(fn func() error) stack.Option {
	return func(o *stack.Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}

// FromContext retrieves a Service from the Context.
func FromContext(ctx context.Context) (Service, bool) {
	s, ok := ctx.Value(serviceKey{}).(Service)
	return s, ok
}

// NewContext returns a new Context with the Service embedded within it.
func NewContext(ctx context.Context, s Service) context.Context {
	return context.WithValue(ctx, serviceKey{}, s)
}

// RegisterHandler is syntactic sugar for registering a handler
func RegisterHandler(s server.Server, h interface{}, opts ...server.HandlerOption) error {
	return s.Handle(s.NewHandler(h, opts...))
}

// RegisterSubscriber is syntactic sugar for registering a subscriber
func RegisterSubscriber(topic string, s server.Server, h interface{}, opts ...server.SubscriberOption) error {
	return s.Subscribe(s.NewSubscriber(topic, h, opts...))
}

// LogMiddlewareWrapper 请求日志中间件
func LogMiddlewareWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		start := time.Now()
		err := fn(ctx, req, rsp)
		duration := time.Since(start)
		if req.Method() != "Acd.Getmakecallevents" && req.Method() != "Acd.Getworklist" { // this method might be called too much
			fmt.Printf("| %s | %v | %s", req.Method(), duration, string(ParseRequest(req)))

		}
		return err
	}
}

// NewService returns a new mucp service
func NewService(name string, version string, brokerAddress string, registryAddress string) Service {
	if name == "mock" {
		return stack.NewService(
			stack.Name(name),
			stack.Version(version),
		)
	}

	b := nats.NewBroker(broker.Addrs(strings.Split(brokerAddress, ",")...))
	r := etcd.NewRegistry(registry.Addrs(strings.Split(registryAddress, ",")...))

	srv := stack.NewService(
		stack.Name(name),
		stack.Version(version),
		stack.Broker(b),
		stack.WrapHandler(LogMiddlewareWrapper),
		stack.Registry(r),
		stack.RegisterInterval(15*time.Second),
		stack.RegisterTTL(30*time.Second),
	)
	return srv
}
