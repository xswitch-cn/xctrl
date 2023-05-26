// Package micro is a pluggable framework for microservices
package xctrl

import (
	"context"

	"git.xswitch.cn/xswitch/xctrl/xctrl/client"
	"git.xswitch.cn/xswitch/xctrl/xctrl/server"
)

type serviceKey struct{}

// Service is an interface that wraps the lower level libraries
// within go-micro. Its a convenience method for building
// and initialising services.
type Service interface {
	// The service name
	Name() string
	// Init initialises options
	Init(...Option)
	// Options returns the current options
	Options() Options
	// Client is used to call services
	Client() client.Client
	// Server is for handling requests and events
	Server() server.Server
	// Run the service
	Run() error
	// The service implementation
	String() string
}

// Function is a one time executing Service
type Function interface {
	// Inherits Service interface
	Service
	// Done signals to complete execution
	Done() error
	// Handle registers an RPC handler
	Handle(v interface{}) error
	// Subscribe registers a subscriber
	Subscribe(topic string, v interface{}) error
}

// Publisher is syntactic sugar for publishing
type Publisher interface {
	Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error
}

type Option func(*Options)

var (
	HeaderPrefix = "Micro-"
)

// NewService creates and returns a new Service based on the packages within.
func NewService(opts ...Option) Service {
	return newService(opts...)
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
