package ctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"git.xswitch.cn/xswitch/xctrl/core/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/core/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/stack/registry"
	"github.com/google/uuid"
)

// Ctrl 控制中心
type Ctrl struct {
	conn             nats.Conn
	uuid             string
	serviceName      string
	registry         registry.Registry
	service          xctrl.XNodeService // 同步调用
	asyncService     xctrl.XNodeService // 异步调用
	aService         xctrl.XNodeService // 异步调用2
	handler          Handler
	enableNodeStatus bool

	hubLock    sync.RWMutex
	channelHub map[string]*Channel

	seq             uint64
	cbLock          sync.RWMutex
	resultCallbacks map[string]*AsyncCallOption
}

type AsyncCallOption struct {
	id   string
	cb   ResultCallbackFunc
	data interface{}
	ts   time.Time
}

type ResultCallbackFunc func(msg *Message, data interface{})

// Handler Ctrl事件响应
type Handler interface {
	// ctx , topic, reply,Params
	Request(context.Context, string, string, *Request)
	// ctx , topic ,reply  Params
	App(context.Context, string, string, *Message)
	// ctx , topic , Params
	Event(context.Context, string, *Request)
	// ctx , topic , Params
	Result(context.Context, string, *Result)
}

var globalCtrl *Ctrl

// UUID get ctrl uuid
func UUID() string {
	return globalCtrl.uuid
}

// ServiceList 服务列表
func ServiceList() ([]*registry.Service, error) {
	return globalCtrl.registry.ListServices()
}

// Service 同步调用
func Service() xctrl.XNodeService {
	if globalCtrl == nil || globalCtrl.service == nil {
		return nil
	}
	return globalCtrl.service
}

// AsyncService 异步调用
func AsyncService() xctrl.XNodeService {
	if globalCtrl == nil || globalCtrl.asyncService == nil {
		return nil
	}
	return globalCtrl.asyncService
}

// AsyncService 异步调用2，使用Context
func AService() xctrl.XNodeService {
	if globalCtrl == nil || globalCtrl.asyncService == nil {
		return nil
	}
	return globalCtrl.aService
}

// Publish 发送消息
func Publish(topic string, msg []byte, opts ...nats.PublishOption) error {
	return globalCtrl.conn.Publish(topic, msg, opts...)
}

// PublishJSON 发送JSON消息
func PublishJSON(topic string, obj interface{}, opts ...nats.PublishOption) error {
	msg, _ := json.MarshalIndent(obj, "", "  ")
	return globalCtrl.conn.Publish(topic, msg, opts...)
}

func Transfer(ctrlID string, channel *xctrl.ChannelEvent) error {
	body, err := json.Marshal(channel)
	if err != nil {
		fmt.Errorf("err:%v", err)
		return err
	}
	channel.State = "START"
	request := Request{
		Version: "2.0",
		Method:  "XNode.Channel",
		Params:  RawMessage(body),
	}
	return globalCtrl.conn.Publish("cn.xswitch.ctrl."+ctrlID, request.Marshal())
}

// Call 发起 request 请求
func Call(topic string, req *Request, timeout time.Duration) (*nats.Message, error) {
	req.Version = "2.0"
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Errorf("execute native api error: %v", err)
		return nil, err
	}
	return globalCtrl.conn.Request(topic, body, timeout)
}

// Call 发起 request 请求
func XCall(topic string, method string, params interface{}, timeout time.Duration) (*nats.Message, error) {
	req := XRequest{
		Version: "2.0",
		Method:  method,
		ID:      "0",
		Params:  params,
	}
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Errorf("execute native api error: %v", err)
		return nil, err
	}
	return globalCtrl.conn.Request(topic, body, timeout)
}

// Respond 响应nats request 请求
func Respond(topic string, resp *Response, opts ...nats.PublishOption) error {
	resp.Version = "2.0"
	body, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Errorf("execute native api error: %v", err)
		return err
	}
	return globalCtrl.conn.Publish(topic, body)
}

// Init 初始化Ctrl global 是否接收全局事件， addrs nats消息队列连接地址
func Init(h Handler, trace bool, addrs string) error {
	c, err := initCtrl(h, trace, strings.Split(addrs, ",")...)
	if err != nil {
		return err
	}
	globalCtrl = c
	return err
}

// ExecAPI 执行原生 API
func ExecAPI(hostname string, cmd string, args ...string) (string, error) {
	node := Node(hostname)
	if node == nil {
		return "", fmt.Errorf("节点未注册: %s", hostname)
	}
	resp, err := Service().NativeAPI(context.Background(),
		&xctrl.NativeRequest{
			Cmd:  cmd,
			Args: strings.Join(args, " "),
		}, WithAddress(node.Uuid))
	if err != nil {
		fmt.Errorf("execute native api error: %v", err)
		return "", err
	}
	if (resp.GetCode() / 100) != 2 {
		fmt.Errorf("[%d] %s", resp.GetCode(), resp.GetMessage())
		return "", fmt.Errorf(resp.GetMessage())
	}
	return resp.GetData(), nil
}

// EnableEvent 开启事件监听
// cn.xswitch.event.cdr
// cn.xswitch.event.custom.sofia>
func EnableEvent(topic string, queue string) error {
	if globalCtrl != nil {
		return globalCtrl.EnableEvent(topic, queue)
	}
	return fmt.Errorf("ctrl uninitialized")
}

// EnableRequest 开启Request请求监听
// FetchXMl
// Dialplan
func EnableRequest(topic string) error {
	if globalCtrl != nil {
		return globalCtrl.EnableRequest(topic)
	}
	return fmt.Errorf("ctrl uninitialized")
}

// EnableApp APP事件
// cn.xswitch.app.callcenter 呼叫队列
// cn.xswitch.app.autodialer 预测外呼
func EnableApp(topic string) error {
	if globalCtrl != nil {
		return globalCtrl.EnableApp(topic)
	}
	return fmt.Errorf("ctrl uninitialized")
}

// EnableNodeStatus 启用节点状态事件
// cn.xswitch.node.status
func EnableNodeStatus() error {
	if globalCtrl != nil {
		return globalCtrl.EnbaleNodeStatus()
	}
	return fmt.Errorf("ctrl uninitialized")
}

func EnableResult(topic string) error {
	if globalCtrl != nil {
		return globalCtrl.EnableResult(topic)
	}
	return fmt.Errorf("ctrl uninitialized")
}

func Subscribe(topic string, cb nats.EventCallback, queue string) (nats.Subscriber, error) {
	if globalCtrl != nil {
		return globalCtrl.Subscribe(topic, cb, queue)
	}
	return nil, fmt.Errorf("ctrl uninitialized")
}

func ToRawMessage(vPoint interface{}) *json.RawMessage {
	d, _ := json.Marshal(vPoint)
	data := json.RawMessage(d)
	return &data
}

type EmptyHandler struct {
}

func (h *EmptyHandler) Request(context.Context, string, string, *Request) {}
func (h *EmptyHandler) App(context.Context, string, string, *Message)     {}
func (h *EmptyHandler) Event(context.Context, string, *Request)           {}
func (h *EmptyHandler) Result(context.Context, string, *Result)           {}

func ACallOption() *AsyncCallOption {
	return &AsyncCallOption{}
}

func (opt *AsyncCallOption) WithCallback(f ResultCallbackFunc) *AsyncCallOption {
	opt.cb = f
	opt.ts = time.Now()
	return opt
}

func (opt *AsyncCallOption) WithData(data interface{}) *AsyncCallOption {
	opt.data = data
	return opt
}

func ACall(subject string, method string, req interface{}, opts ...*AsyncCallOption) error {
	opt := ACallOption()
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.id == "" {
		opt.id = uuid.New().String()
	}

	rpc := Request{
		Version: "2.0",
		ID:      ToRawMessage(opt.id),
		Method:  fmt.Sprintf("XNode.%s", method),
		Params:  ToRawMessage(req),
	}

	bytes, err := json.Marshal(rpc)

	if err != nil {
		return err
	}

	if opt.cb != nil {
		globalCtrl.cbLock.Lock()
		globalCtrl.resultCallbacks[opt.id] = opt
		globalCtrl.cbLock.Unlock()
	}

	Publish(subject, bytes)

	return nil
}

func DoResultCallback(msg *Message) {
	globalCtrl.cbLock.RLock()

	id := ""
	err := json.Unmarshal(*msg.ID, &id)
	if err != nil {
		return
	}
	opt, found := globalCtrl.resultCallbacks[id]
	globalCtrl.cbLock.RUnlock()

	if found {
		globalCtrl.cbLock.Lock()
		delete(globalCtrl.resultCallbacks, id)
		globalCtrl.cbLock.Unlock()
		// xlog.Errorf("do callback %s", id)
		opt.cb(msg, opt.data)
	}
}
