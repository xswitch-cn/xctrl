package ctrl

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"sync/atomic"
	"time"

	"git.xswitch.cn/xswitch/xctrl/core/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/core/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/xctrl/client"
	"git.xswitch.cn/xswitch/xctrl/xctrl/codec"
	"git.xswitch.cn/xswitch/xctrl/xctrl/errors"
	"git.xswitch.cn/xswitch/xctrl/xctrl/metadata"
	"git.xswitch.cn/xswitch/xctrl/xctrl/util/log"
)

const defaultTimeout = 60 * time.Second
const defaultDialTimeout = 5 * time.Second

func newOptions(options ...client.Option) client.Options {
	opts := client.Options{
		Codecs: make(map[string]codec.NewCodec),
		CallOptions: client.CallOptions{
			Backoff:        client.DefaultBackoff,
			Retry:          client.DefaultRetry,
			Retries:        client.DefaultRetries,
			RequestTimeout: client.DefaultRequestTimeout,
			DialTimeout:    defaultDialTimeout,
		},
		PoolSize: client.DefaultPoolSize,
		PoolTTL:  client.DefaultPoolTTL,
	}

	for _, o := range options {
		o(&opts)
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	return opts
}

type rpcRequest struct {
	service     string
	method      string
	endpoint    string
	contentType string
	codec       codec.Codec
	body        interface{}
	opts        client.RequestOptions
}

func newRequest(service, endpoint string, request interface{}, contentType string, reqOpts ...client.RequestOption) client.Request {
	var opts client.RequestOptions

	for _, o := range reqOpts {
		o(&opts)
	}

	// set the content-type specified
	if len(opts.ContentType) > 0 {
		contentType = opts.ContentType
	}

	return &rpcRequest{
		service:     service,
		method:      endpoint,
		endpoint:    endpoint,
		body:        request,
		contentType: contentType,
		opts:        opts,
	}
}

func (r *rpcRequest) ContentType() string {
	return r.contentType
}

func (r *rpcRequest) Service() string {
	return r.service
}

func (r *rpcRequest) Method() string {
	return r.method
}

func (r *rpcRequest) Endpoint() string {
	return r.endpoint
}

func (r *rpcRequest) Body() interface{} {
	return r.body
}

func (r *rpcRequest) Codec() codec.Writer {
	return r.codec
}

func (r *rpcRequest) Stream() bool {
	return r.opts.Stream
}

type ctrlClient struct {
	conn     nats.Conn
	opts     client.Options
	seq      uint64
	async    bool
	aservice bool
}

// NewClient New node client
func newClient(conn nats.Conn, async bool, opt ...client.Option) client.Client {
	opts := newOptions(opt...)
	return &ctrlClient{
		conn:  conn,
		opts:  opts,
		async: async,
	}
}

func (r *ctrlClient) Init(opts ...client.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	return nil
}

func (r *ctrlClient) SetAService() error {
	r.aservice = true
	return nil
}

func (r *ctrlClient) Options() client.Options {
	return r.opts
}

func (r *ctrlClient) Call(ctx context.Context, request client.Request, response interface{}, opts ...client.CallOption) error {

	// make a copy of call opts
	callOpts := r.opts.CallOptions
	// todo check this
	for _, opt := range opts {
		opt(&callOpts)
	}
	// TODO 默认设置消息最长响应时间为24小时
	if callOpts.RequestTimeout <= 0 {
		callOpts.RequestTimeout = 24 * time.Hour
	}

	// check if we already have a deadline
	d, ok := ctx.Deadline()
	if !ok {
		// no deadline so we create a new one
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, callOpts.RequestTimeout)
		defer cancel()
	} else {
		// got a deadline so no need to setup context
		// but we need to set the timeout we pass along
		opt := client.WithRequestTimeout(d.Sub(time.Now()))
		opt(&callOpts)
	}

	// should we noop right here?
	select {
	case <-ctx.Done():
		return errors.Timeout("nats.jsonrpc.client", fmt.Sprintf("%v", ctx.Err()))
	default:
	}

	// make copy of call method
	rcall := r.call

	// wrap the call in reverse
	for i := len(callOpts.CallWrappers); i > 0; i-- {
		rcall = callOpts.CallWrappers[i-1](rcall)
	}

	// make the call
	err := rcall(ctx, request, response, callOpts)
	return err
}

func (r *ctrlClient) call(ctx context.Context, req client.Request, resp interface{}, opts client.CallOptions) error {
	var err error
	request := new(Request)
	requestID, ok := metadata.Get(ctx, "request-seq-id")
	if !ok || requestID == "" {
		requestID = fmt.Sprintf(`"%d"`, atomic.AddUint64(&r.seq, 1))
	}
	id := json.RawMessage(requestID)
	request.ID = &id
	request.Method = req.Method()
	request.Version = "2.0"
	data, _ := json.Marshal(req.Body())
	raw := json.RawMessage(data)
	request.Params = &raw
	body, _ := json.MarshalIndent(request, "", "  ")

	address := ""
	if len(opts.Address) > 0 {
		address = opts.Address[0]
	}
	if address == "" {
		return errors.BadRequest("nats.jsonrpc.client", "The address cannot be empty")
	}

	requestTimeout := opts.RequestTimeout
	if requestTimeout == 0 {
		requestTimeout = defaultTimeout
	}

	if r.async {
		err = r.conn.Publish(address, body)
		if err != nil {
			log.Errorf("err : %v", err)
			return errors.Timeout("nats.jsonrpc.client", fmt.Sprintf("%v", err))
		}

		response := &xctrl.Response{
			Code:    http.StatusCreated,
			Message: requestID,
		}

		b, _ := json.MarshalIndent(response, "", " ")
		json.Unmarshal(b, resp)
		return nil
	}

	var msg *nats.Message
	var rsp Response

	if r.aservice {
		msg, err = r.conn.RequestWithContext(ctx, address, body)
	} else {
		msg, err = r.conn.Request(address, body, requestTimeout)
	}

	if err != nil {
		if err.Error() == "context canceled" {
			return errors.Canceled("nats.jsonrpc.client", "%v", err)
		}

		if err.Error() == "nats: timeout" {
			return errors.Timeout("nats.jsonrpc.client", "%v", err)
		}

		if err.Error() == "context deadline exceeded" {
			return errors.Timeout("nats.jsonrpc.client", "%v", err)
		}

		return errors.InternalServerError("nats.jsonrpc.client", "%v", err)
	}

	rsp.Result = resp
	err = json.Unmarshal(msg.Body, &rsp)

	if err != nil || rsp.getError() != nil {
		log.Errorf("nats.jsonrpc.client %v %v %v", err, resp, rsp)
		return errors.InternalServerError("nats.jsonrpc.client", "%v", err)
	}
	return nil
}

func (r *ctrlClient) Stream(ctx context.Context, request client.Request, opts ...client.CallOption) (client.Stream, error) {
	return nil, errServer
}

func (r *ctrlClient) Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	return errServer
}

func (r *ctrlClient) NewMessage(topic string, message interface{}, opts ...client.MessageOption) client.Message {
	return nil
}

func (r *ctrlClient) NewRequest(service, method string, request interface{}, reqOpts ...client.RequestOption) client.Request {
	//change part of the request's method into NativeJSAPI
	method = TranslateMethod(method)
	return newRequest(service, method, request, r.opts.ContentType, reqOpts...)
}

func (r *ctrlClient) String() string {
	return "ctrl"
}
