package tboy

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xswitch-cn/proto/go/proto/xctrl"
	"github.com/xswitch-cn/xctrl/ctrl"
	"github.com/xswitch-cn/xctrl/ctrl/nats"
)

func init() {
	log = logrus.New()
	log.SetReportCaller(true)
}

// FakeFakeIvrChannel 模拟IVR通道
type FakeFakeIvrChannel struct {
	CtrlUuid string
	UUID     string
	Lock     sync.RWMutex
	Data     *xctrl.ChannelEvent
	Context  context.Context
	Cancel   context.CancelFunc
	DTMF     string // 存储接收到的DTMF按键
}

// FakeFakeIvrServer IVR服务器
type FakeFakeIvrServer struct {
	NodeUUID string
	Domain   string
	Channels map[string]*FakeFakeIvrChannel
	Lock     sync.RWMutex
}

// NewFakeFakeIvrServer 创建新的IVR服务器
func NewFakeFakeIvrServer(nodeUUID, domain string) *FakeFakeIvrServer {
	return &FakeFakeIvrServer{
		NodeUUID: nodeUUID,
		Domain:   domain,
		Channels: make(map[string]*FakeFakeIvrChannel),
	}
}

// CacheChannel 缓存通道
func (server *FakeFakeIvrServer) CacheChannel(uuid string, channel *FakeFakeIvrChannel) {
	server.Lock.Lock()
	defer server.Lock.Unlock()
	server.Channels[uuid] = channel
}

// GetChannel 获取通道
func (server *FakeFakeIvrServer) GetChannel(uuid string) (*FakeFakeIvrChannel, bool) {
	server.Lock.RLock()
	defer server.Lock.RUnlock()
	channel, ok := server.Channels[uuid]
	return channel, ok
}

// DeleteChannel 删除通道
func (server *FakeFakeIvrServer) DeleteChannel(uuid string) {
	server.Lock.Lock()
	defer server.Lock.Unlock()
	delete(server.Channels, uuid)
}

// OK 发送成功响应
func (server *FakeFakeIvrServer) OK(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)
	if err != nil {
		return
	}

	res := xctrl.Response{
		Code:     200,
		Message:  "OK",
		NodeUuid: server.NodeUUID,
		Uuid:     request.Uuid,
	}

	resBytes, _ := json.Marshal(res)
	raw := json.RawMessage(resBytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	resultBytes, _ := json.Marshal(result)

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, resultBytes)
}

// Error 发送错误响应
func (server *FakeFakeIvrServer) Error(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)
	if err != nil {
		return
	}

	res := xctrl.Response{
		Code:     404,
		Message:  "Unsupported Method " + msg.Method,
		NodeUuid: server.NodeUUID,
		Uuid:     request.Uuid,
	}

	resBytes, _ := json.Marshal(res)
	raw := json.RawMessage(resBytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	resultBytes, _ := json.Marshal(result)

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, resultBytes)
}

// Answer 处理应答请求
func (server *FakeFakeIvrServer) Answer(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request
	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	channel, ok := server.GetChannel(request.Uuid)
	if ok {
		channel.CtrlUuid = request.CtrlUuid
	}

	server.OK(ctx, msg, reply)

	// 发送ANSWERED事件
	if channel != nil {
		controller := "cn.xswitch.ctrl." + channel.CtrlUuid
		channel.Data.State = "ANSWERED"
		channel.Data.AnswerEpoch = uint32(time.Now().Unix())

		eventReq := ctrl.Request{
			Version: "2.0",
			Method:  "Event.Channel",
			Params:  ctrl.ToRawMessage(channel.Data),
		}

		reqStr, _ := json.Marshal(eventReq)
		ctrl.Publish(controller, reqStr)
	}
}

// Hangup 处理挂机请求
func (server *FakeFakeIvrServer) Hangup(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request
	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	channel, ok := server.GetChannel(request.Uuid)
	if !ok {
		server.Error(ctx, msg, reply)
		return
	}

	if channel.Cancel != nil {
		channel.Cancel()
	}

	server.OK(ctx, msg, reply)

	// 发送DESTROY事件
	controller := "cn.xswitch.ctrl." + channel.CtrlUuid
	channel.Data.State = "DESTROY"
	channel.Data.Cause = "NORMAL_CLEARING"

	eventReq := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channel.Data),
	}

	reqStr, _ := json.Marshal(eventReq)
	ctrl.Publish(controller, reqStr)

	// 清理通道
	server.DeleteChannel(request.Uuid)
}

// Play 处理播放请求
func (server *FakeFakeIvrServer) Play(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.PlayRequest
	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	_, ok := server.GetChannel(request.Uuid)
	if !ok {
		server.Error(ctx, msg, reply)
		return
	}

	// 模拟播放完成
	go func() {
		time.Sleep(2 * time.Second) // 模拟播放时间
		server.OK(ctx, msg, reply)
	}()
}

// ReadDTMF 处理DTMF读取请求
func (server *FakeFakeIvrServer) ReadDTMF(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.DTMFRequest
	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	channel, ok := server.GetChannel(request.Uuid)
	if !ok {
		server.Error(ctx, msg, reply)
		return
	}

	// 模拟DTMF输入
	go func() {
		time.Sleep(1 * time.Second) // 等待用户输入

		// 设置默认DTMF为"1"
		dtmf := "1"
		if channel.DTMF != "" {
			dtmf = channel.DTMF
		}

		res := xctrl.DTMFResponse{
			Code:       200,
			Message:    "OK",
			NodeUuid:   server.NodeUUID,
			Uuid:       request.Uuid,
			Dtmf:       dtmf,
			Terminator: "#",
		}

		resBytes, _ := json.Marshal(res)
		raw := json.RawMessage(resBytes)

		result := ctrl.Result{
			Version: "2.0",
			ID:      msg.ID,
			Result:  &raw,
		}

		resultBytes, _ := json.Marshal(result)

		if reply == "" {
			reply = "cn.xswitch.ctrl." + request.CtrlUuid
		}

		ctrl.Publish(reply, resultBytes)
	}()
}

// SetVar 处理设置变量请求
func (server *FakeFakeIvrServer) SetVar(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.SetVarRequest
	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	channel, ok := server.GetChannel(request.Uuid)
	if !ok {
		server.Error(ctx, msg, reply)
		return
	}

	channel.Lock.Lock()
	if channel.Data.Params == nil {
		channel.Data.Params = make(map[string]string)
	}
	for k, v := range request.Data {
		channel.Data.Params[k] = v
	}
	channel.Lock.Unlock()

	server.OK(ctx, msg, reply)
}

// CreateTestCall 创建测试通话
func (server *FakeFakeIvrServer) CreateTestCall(ctrlUUID, caller, callee string) {
	callUUID := uuid.New().String()

	channelEvent := xctrl.ChannelEvent{
		NodeUuid:    server.NodeUUID,
		Uuid:        callUUID,
		Direction:   "inbound",
		State:       "START",
		CidName:     caller,
		CidNumber:   caller,
		DestNumber:  callee,
		AnswerEpoch: uint32(time.Now().Unix()),
		Answered:    false,
		CreateEpoch: uint32(time.Now().Unix()),
		Params: map[string]string{
			"xcc_session": callUUID,
		},
	}

	channel := &FakeFakeIvrChannel{
		CtrlUuid: ctrlUUID,
		UUID:     callUUID,
		Data:     &channelEvent,
		DTMF:     "1", // 默认DTMF输入
	}
	channel.Context, channel.Cancel = context.WithCancel(context.Background())

	server.CacheChannel(callUUID, channel)

	// 发送START事件
	controller := "cn.xswitch.ctrl." + ctrlUUID
	eventReq := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channelEvent),
	}

	reqStr, _ := json.Marshal(eventReq)
	ctrl.Publish(controller, reqStr)

	log.Infof("Created test call: %s -> %s, UUID: %s", caller, callee, callUUID)
}

// SetDTMF 设置DTMF输入（用于测试）
func (server *FakeFakeIvrServer) SetDTMF(callUUID, dtmf string) bool {
	channel, ok := server.GetChannel(callUUID)
	if !ok {
		return false
	}

	channel.DTMF = dtmf
	return true
}

// Event 处理事件
func (server *FakeFakeIvrServer) Event(msg *ctrl.Message, natsEvent nats.Event) {
	topic := natsEvent.Topic()
	reply := natsEvent.Reply()
	log.Infof("Received event: %s %s", topic, msg.Method)

	if msg.Method == "" && msg.Result != nil {
		log.Infof("Got a response: %s", msg.ID)
		return
	}

	ctx := context.Background()

	switch msg.Method {
	case "XNode.Answer":
		server.Answer(ctx, msg, reply)
	case "XNode.Hangup":
		server.Hangup(ctx, msg, reply)
	case "XNode.Play":
		server.Play(ctx, msg, reply)
	case "XNode.ReadDTMF":
		server.ReadDTMF(ctx, msg, reply)
	case "XNode.SetVar":
		server.SetVar(ctx, msg, reply)
	default:
		log.Errorf("Unsupported Method: %s", msg.Method)
		server.Error(ctx, msg, reply)
	}
}
