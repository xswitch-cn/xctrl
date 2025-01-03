// Package nats provides a NATS Conn
package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"git.xswitch.cn/xswitch/proto/xctrl/util/log"
	nats "github.com/nats-io/nats.go"
)

type nConn struct {
	sync.RWMutex
	addrs []string
	conn  *nats.Conn
	opts  Options
	nopts nats.Options
	drain bool
	trace bool
}

type subscriber struct {
	s     *nats.Subscription
	opts  SubscribeOptions
	drain bool
}

// Publication .
type Publication struct {
	t     string
	err   error
	m     *Message
	reply string
}

func (p *Publication) Error() error {
	return p.err
}

// Topic .
func (p *Publication) Topic() string {
	return p.t
}

// Message .
func (p *Publication) Message() *Message {
	return p.m
}

// Ack .
func (p *Publication) Ack() error {
	return nil
}

// Reply .
func (p *Publication) Reply() string {
	return p.reply
}

func (n *subscriber) Options() SubscribeOptions {
	return n.opts
}

func (n *subscriber) Topic() string {
	return n.s.Subject
}

func (n *subscriber) Unsubscribe() error {
	if n.drain {
		return n.s.Drain()
	}
	return n.s.Unsubscribe()
}

func (n *subscriber) SetPendingLimits(msgLimit, bytesLimit int) error {
	return n.s.SetPendingLimits(msgLimit, bytesLimit)
}

func (n *nConn) Address() string {
	if n.conn != nil && n.conn.IsConnected() {
		return n.conn.ConnectedUrl()
	}
	if len(n.addrs) > 0 {
		return n.addrs[0]
	}

	return ""
}

func setAddrs(addrs []string) []string {
	var cAddrs []string
	for _, addr := range addrs {
		if len(addr) == 0 {
			continue
		}
		if !strings.HasPrefix(addr, "nats://") {
			addr = "nats://" + addr
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{nats.DefaultURL}
	}
	return cAddrs
}

func natsErrHandler(nc *nats.Conn, sub *nats.Subscription, natsErr error) {
	fmt.Printf("%v\n", natsErr)
	if natsErr == nats.ErrSlowConsumer {
		pendingMsgs, _, err := sub.Pending()
		if err != nil {
			log.Error(fmt.Errorf("couldn't get pending messages: %v", err))
			return
		}
		log.Error(fmt.Errorf("falling behind with %d pending messages on subject %q ", pendingMsgs, sub.Subject))
	}
}

func (n *nConn) Connect() error {
	n.Lock()
	defer n.Unlock()

	status := nats.CLOSED
	if n.conn != nil {
		status = n.conn.Status()
	}

	switch status {
	case nats.CONNECTED, nats.RECONNECTING, nats.CONNECTING:
		return nil
	default: // DISCONNECTED or CLOSED or DRAINING
		opts := n.nopts
		opts.Servers = n.addrs
		opts.Secure = n.opts.Secure
		opts.TLSConfig = n.opts.TLSConfig
		opts.AsyncErrorCB = natsErrHandler

		// secure might not be set
		if n.opts.TLSConfig != nil {
			opts.Secure = true
		}

		c, err := opts.Connect()
		if err != nil {
			return err
		}
		n.conn = c
		return nil
	}
}

func (n *nConn) Disconnect() error {
	n.RLock()
	if n.drain {
		n.conn.Drain()
	} else {
		n.conn.Close()
	}
	n.RUnlock()
	return nil
}

func (n *nConn) GetConn() *nats.Conn {
	return n.conn
}

func (n *nConn) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	n.addrs = setAddrs(n.opts.Addrs)
	return nil
}

func (n *nConn) Options() Options {
	return n.opts
}

func (n *nConn) Publish(topic string, msg []byte, opts ...PublishOption) error {
	if n.trace {
		log.Tracef("NATS Publish: %s \n%s\n", topic, string(msg))
	}
	return n.conn.Publish(topic, msg)
}

func (n *nConn) Request(topic string, data []byte, timeout time.Duration) (*Message, error) {
	if n.trace {
		log.Tracef("NATS Request: topic=%s\n%s\n", topic, string(data))
	}
	msg, err := n.conn.Request(topic, data, timeout)
	var m Message
	if err != nil {
		return &m, err
	}
	m.Body = msg.Data
	m.Header = make(map[string]string)
	m.Header["topic"] = msg.Subject
	if n.trace {
		log.Tracef("NATS Response: \n%s\n", string(m.Body))
	}
	return &m, nil
}

func (n *nConn) RequestWithContext(ctx context.Context, topic string, data []byte) (*Message, error) {
	if n.trace {
		log.Tracef("Request: %s \n%s", topic, string(data))
	}
	msg, err := n.conn.RequestWithContext(ctx, topic, data)
	var m Message
	if err != nil {
		return &m, err
	}
	m.Body = msg.Data
	m.Header = make(map[string]string)
	m.Header["topic"] = msg.Subject
	if n.trace {
		log.Tracef("NATS Response: \n%s", string(m.Body))
	}
	return &m, nil
}

func (n *nConn) Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	if n.conn == nil {
		return nil, errors.New("not connected")
	}

	opt := SubscribeOptions{
		AutoAck: true,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	var drain bool
	if _, ok := opt.Context.Value(drainSubscriptionKey{}).(bool); ok {
		drain = true
	}

	fn := func(msg *nats.Msg) {
		if n.trace && msg.Subject != "cn.xswitch.ctrl" {
			var data bytes.Buffer
			err := json.Indent(&data, msg.Data, "", "\t")
			if err != nil {
				fmt.Printf("nats MarshalIndent error:%v\n", err)
			} else {
				if msg.Reply != "" {
					fmt.Printf("Message: %s|%s \n%s\n", msg.Subject, msg.Reply, data.String())
				} else {
					fmt.Printf("Message: %s \n%s\n", msg.Subject, string(msg.Data))
				}
			}

		}

		var m Message
		m.Header = map[string]string{"Content-Type": "application/json"}
		m.Body = msg.Data
		handler(&Publication{m: &m, t: msg.Subject, reply: msg.Reply})
	}

	var sub *nats.Subscription
	var err error

	n.RLock()
	if len(opt.Queue) > 0 {
		sub, err = n.conn.QueueSubscribe(topic, opt.Queue, fn)
	} else {
		sub, err = n.conn.Subscribe(topic, fn)
	}
	n.RUnlock()
	if err != nil {
		return nil, err
	}
	return &subscriber{s: sub, opts: opt, drain: drain}, nil
}

func (n *nConn) String() string {
	return "nats"
}

// Conn is an interface used for asynchronous messaging.
type Conn interface {
	Init(...Option) error
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Publish(topic string, m []byte, opts ...PublishOption) error
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	Request(topic string, data []byte, timeout time.Duration) (*Message, error)
	RequestWithContext(ctx context.Context, topic string, data []byte) (*Message, error)
	String() string
	GetConn() *nats.Conn
}

// NewConn .
func NewConn(opts ...Option) Conn {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	natsOpts := nats.GetDefaultOptions()
	if n, ok := options.Context.Value(optionsKey{}).(nats.Options); ok {
		natsOpts = n
	}

	var drain bool
	if _, ok := options.Context.Value(drainSubscriptionKey{}).(bool); ok {
		drain = true
	}

	// Conn.Options have higher priority than nats.Options
	// only if Addrs, Secure or TLSConfig were not set through a Conn.Option
	// we read them from nats.Option
	if len(options.Addrs) == 0 {
		options.Addrs = natsOpts.Servers
	}

	if !options.Secure {
		options.Secure = natsOpts.Secure
	}

	if options.TLSConfig == nil {
		options.TLSConfig = natsOpts.TLSConfig
	}
	if options.TLSConnectInformation != nil {
		rootCAErr := nats.RootCAs(options.TLSConnectInformation.RootCAs)(&natsOpts)
		if rootCAErr != nil {
			fmt.Printf("RootCAs init error: %s \n", rootCAErr.Error())
		}

		certErr := nats.ClientCert(options.TLSConnectInformation.Cert, options.TLSConnectInformation.Key)(&natsOpts)
		if certErr != nil {
			fmt.Printf("Certification init error: %s \n", certErr.Error())
		}
	}

	nb := &nConn{
		opts:  options,
		nopts: natsOpts,
		addrs: setAddrs(options.Addrs),
		drain: drain,
		trace: options.Trace,
	}

	return nb
}
