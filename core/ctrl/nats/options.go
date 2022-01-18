package nats

import (
	"context"
	"crypto/tls"

	"git.xswitch.cn/xswitch/xctrl/stack/codec"
	"git.xswitch.cn/xswitch/xctrl/stack/registry"
)

type optionsKey struct{}
type drainConnectionKey struct{}
type drainSubscriptionKey struct{}

// DrainConnection will drain subscription on close
func DrainConnection() Option {
	return setBrokerOption(drainConnectionKey{}, true)
}

// DrainSubscription will drain pending messages when unsubscribe
func DrainSubscription() SubscribeOption {
	return setSubscribeOption(drainSubscriptionKey{}, true)
}

// Options .
type Options struct {
	Addrs  []string
	Secure bool
	Codec  codec.Marshaler

	// Handler executed when error happens in broker mesage
	// processing
	ErrorHandler Handler

	TLSConfig *tls.Config
	// Registry used for clustering
	Registry registry.Registry
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
	Trace   bool
}

// PublishOptions  is Publish meesage options
type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// SubscribeOptions is Subscribe options
type SubscribeOptions struct {
	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option is conn option
type Option func(*Options)

// PublishOption is Publish Option callbcak
type PublishOption func(*PublishOptions)

// PublishContext set context
func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

// SubscribeOption is subscribe option
type SubscribeOption func(*SubscribeOptions)

// NewSubscribeOptions new subscribe options
func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Addrs sets the host addresses to be used by the broker
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Codec sets the codec used for encoding/decoding used where
// a broker does not support headers
func Codec(c codec.Marshaler) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

// DisableAutoAck will disable auto acking of messages
// after they have been handled.
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

// ErrorHandler will catch all broker errors that cant be handled
// in normal way, for example Codec errors
func ErrorHandler(h Handler) Option {
	return func(o *Options) {
		o.ErrorHandler = h
	}
}

func Trace(enable bool) Option {
	return func(o *Options) {
		o.Trace = enable
	}
}

// Queue sets the name of the queue to share messages on
func Queue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}

// Registry is registry
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Secure communication with the broker
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// TLSConfig Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// SubscribeContext set context
func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(Event) error

// Message is data
type Message struct {
	Header map[string]string
	Body   []byte
}

// Event is given to a subscription handler for processing
type Event interface {
	Topic() string
	Message() *Message
	Reply() string
	Ack() error
	Error() error
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
	SetPendingLimits(msgLimit, bytesLimit int) error
}

type EventCallback func(context.Context, Event) error
