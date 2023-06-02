package ctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"git.xswitch.cn/xswitch/xctrl/ctrl/bus"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/proto/cman"
	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/xctrl/client"
	"git.xswitch.cn/xswitch/xctrl/xctrl/util/log"
	"github.com/google/uuid"
)

// Ctrl 控制中心
type Ctrl struct {
	conn             nats.Conn
	uuid             string
	serviceName      string
	service          xctrl.XNodeService // 同步调用
	asyncService     xctrl.XNodeService // 异步调用
	aService         xctrl.XNodeService // 异步调用2
	cmanService      cman.CManService   // CManService
	enableNodeStatus bool
	forkDTMFEvent    bool // fork and push Event.DTMF into Channel Event Thread Too

	hubLock    sync.RWMutex
	channelHub map[string]*Channel

	seq             uint64
	cbLock          sync.RWMutex
	resultCallbacks map[string]*AsyncCallOption
	nodeCallback    NodeHashFun
}

type AsyncCallOption struct {
	id   string
	cb   ResultCallbackFunc
	data interface{}
	ts   time.Time
}

type ResultCallbackFunc func(msg *Message, data interface{})

// Handler Ctrl事件响应
type EventHandler interface {
	Event(req *Request, natsEvent nats.Event)
}

type AppHandler interface {
	ChannelEvent(context.Context, *Channel)
	Event(msg *Message, natsEvent nats.Event)
}

type RequestHandler interface {
	Request(req *Request, natsEvent nats.Event)
}

type ContextKey string

type LogLevel int

const (
	LLFatal LogLevel = iota
	LLError
	LLWarn
	LLInfo
	LLDebug
	LLTrace
)

type Logger interface {
	Log(level int, v ...interface{})
	Logf(level int, format string, v ...interface{})
}

var globalCtrl *Ctrl

// UUID get ctrl uuid
func UUID() string {
	return globalCtrl.uuid
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

// CManService 同步调用
func CManService() cman.CManService {
	if globalCtrl == nil || globalCtrl.cmanService == nil {
		return nil
	}
	return globalCtrl.cmanService
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

func CtrlStartUp(req *xctrl.CtrlStartUpRequest) error {
	request := Request{
		Version: "2.0",
		Method:  "XNode.CtrlStartUp",
		Params:  ToRawMessage(req),
	}
	return globalCtrl.conn.Publish("cn.xswitch.node", request.Marshal())
}

// Call 发起 request 请求
func Call(topic string, req *Request, timeout time.Duration) (*nats.Message, error) {
	req.Version = "2.0"
	//change part of the request's method into NativeJSAPI
	req.Method = TranslateMethod(req.Method)

	body, err := json.Marshal(req)
	if err != nil {
		log.Errorf("execute native api error: %v", err)
		return nil, err
	}
	return globalCtrl.conn.Request(topic, body, timeout)
}

// XCall 发起 request 请求
func XCall(topic string, method string, params interface{}, timeout time.Duration) (*nats.Message, error) {
	//change part of the request's method into NativeJSAPI
	method = TranslateMethod(method)
	req := XRequest{
		Version: "2.0",
		Method:  method,
		ID:      "0",
		Params:  params,
	}
	body, err := json.Marshal(req)
	if err != nil {
		log.Errorf("execute native api error: %v", err)
		return nil, err
	}
	return globalCtrl.conn.Request(topic, body, timeout)
}

// Respond 响应NATS Request 请求
func Respond(topic string, resp *Response, opts ...nats.PublishOption) error {
	resp.Version = "2.0"
	body, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Errorf("execute native api error: %v", err)
		return err
	}
	return globalCtrl.conn.Publish(topic, body)
}

// Init 初始化Ctrl trace 是否开启NATS消息跟踪， addrs nats消息队列连接地址
func Init(trace bool, addrs string) error {
	log.Infof("ctrl starting with addrs=%s\n", addrs)
	c, err := initCtrl(trace, strings.Split(addrs, ",")...)
	if err != nil {
		return err
	}
	globalCtrl = c
	xctrl.SetService(&c.service)
	return err
}

func InitCManService(addr string) error {
	if globalCtrl != nil {
		globalCtrl.NewCManService(addr)
		return nil
	}
	return fmt.Errorf("ctrl uninitialized")
}

// EnableEvent 开启事件监听
// cn.xswitch.ctrl.event.cdr
// cn.xswitch.ctrl.event.custom.sofia>
func EnableEvent(handler EventHandler, subject string, queue string) error {
	if globalCtrl != nil {
		return globalCtrl.EnableEvent(handler, subject, queue)
	}
	return fmt.Errorf("ctrl uninitialized")
}

// EnableRequest 开启Request请求监听
// FetchXMl
// Dialplan
func EnableRequest(handler RequestHandler, subject string, queue string) error {
	if globalCtrl != nil {
		return globalCtrl.EnableRequest(handler, subject, queue)
	}
	return fmt.Errorf("ctrl uninitialized")
}

// EnableApp APP事件
func EnableApp(handler AppHandler, subject string, queue string) error {
	if globalCtrl != nil {
		return globalCtrl.EnableApp(handler, subject, queue)
	}
	return fmt.Errorf("ctrl uninitialized")
}

// EnableNodeStatus 启用节点状态事件
// cn.xswitch.node.status
func EnableNodeStatus(subject string) error {
	if globalCtrl != nil {
		return globalCtrl.EnbaleNodeStatus(subject)
	}
	return fmt.Errorf("ctrl uninitialized")
}

// ForkDTMFEventToChannelEvent 将DTMF事件放到ChannelEvent事件相同的线程处理
func ForkDTMFEventToChannelEventThread() error {
	if globalCtrl != nil {
		return globalCtrl.ForkDTMFEventToChannelEventThread()
	}
	return fmt.Errorf("ctrl uninitialized")
}

func DeliverToChannelEventThread(channel *Channel, natsEvent nats.Event) {
	ev := bus.NewEvent(channel.GetState(), channel.GetUuid(), channel, natsEvent)
	bus.Publish(ev)
}

func Subscribe(subject string, cb nats.EventCallback, queue string) (nats.Subscriber, error) {
	if globalCtrl != nil {
		return globalCtrl.Subscribe(subject, cb, queue)
	}
	return nil, fmt.Errorf("ctrl uninitialized")
}

func ToRawMessage(vPoint interface{}) *json.RawMessage {
	d, _ := json.Marshal(vPoint)
	data := json.RawMessage(d)
	return &data
}

type EmptyAppHandler struct{}
type EmptyEventHandler struct{}
type EmptyRequestHandler struct{}

func (h *EmptyAppHandler) ChannelEvent(context.Context, *Channel) {}
func (h *EmptyAppHandler) Event(*Message, nats.Event)             {}
func (h *EmptyEventHandler) Event(*Request)                       {}
func (h *EmptyRequestHandler) Request(*Request, nats.Event)       {}

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
		log.Tracef("do callback %s", id)
		opt.cb(msg, opt.data)
	}
}

// TranslateMethod change request method to NativeJSAPI
func TranslateMethod(method string) string {
	//XNode.JStatus->XNode.NativeJSAPI
	if method == "XNode.JStatus" {
		return "XNode.NativeJSAPI"
	}
	//XNode.ConferenceInfo->XNode.NativeJSAPI
	if method == "XNode.ConferenceInfo" {
		return "XNode.NativeJSAPI"
	}
	//XNode.ConferenceList->XNode.NativeJSAPI
	if method == "XNode.ConferenceList" {
		return "XNode.NativeJSAPI"
	}
	return method
}

func WithAddressDefault() client.CallOption {
	return client.WithAddress("cn.xswitch.node")
}

// Node Address 标准化Node地址
func NodeAddress(nodeUUID string) string {
	if nodeUUID == "" {
		return "cn.xswitch.node"
	}
	if !strings.HasPrefix(nodeUUID, "cn.xswitch.") {
		return "cn.xswitch.node." + nodeUUID
	}
	return nodeUUID
}

// WithAddress 创建Node地址
func WithAddress(nodeUUID string) client.CallOption {
	return client.WithAddress(NodeAddress(nodeUUID))
}

// NATS Request Timeout
func WithRequestTimeout(d time.Duration) client.CallOption {
	return client.WithRequestTimeout(d)
}

func WithTimeout(d time.Duration) client.CallOption {
	return client.WithRequestTimeout(d)
}

func SetLogLevel(level LogLevel) {
	log.SetLevel(log.Level(level))
}

func SetLogger(l Logger) {
	log.SetLogger(l)
}

func RegisterHashNodeFun(nodeCallbackFunc NodeHashFun) {
	if globalCtrl != nil {
		globalCtrl.registerHashNodeFun(nodeCallbackFunc)
	}
}
