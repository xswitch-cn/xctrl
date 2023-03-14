package ctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"git.xswitch.cn/xswitch/xctrl/core/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/stack/client"
	"git.xswitch.cn/xswitch/xctrl/stack/errors"
	"github.com/google/uuid"

	"git.xswitch.cn/xswitch/xctrl/core/proto/xctrl"
)

// Channel call channel
type Channel struct {
	xctrl.ChannelEvent
	lock     sync.RWMutex
	CtrlUuid string
	subs     []nats.Subscriber
}

// only call at the first time
func NewChannel(channel_uuid string) *Channel {
	if channel_uuid == "" {
		channel_uuid = uuid.New().String()
	}
	channel := &Channel{}
	channel.CtrlUuid = fmt.Sprintf("channel.%s", channel_uuid)
	channel.Uuid = channel.CtrlUuid
	return channel.Save()
}

// WithAddress 创建NODE地址
func WithAddressNot() client.CallOption {
	return client.WithAddress("cn.xswitch.node")
}

// WithAddress 创建NODE地址
func NodeAddress(nodeUUID string) string {
	if nodeUUID == "" {
		return "cn.xswitch.node"
	}

	return fmt.Sprintf("cn.xswitch.node.%s", nodeUUID)
}

// WithAddress 创建NODE地址
func WithAddress(nodeUUID string) client.CallOption {
	if nodeUUID == "" {
		return client.WithAddress("cn.xswitch.node")
	}

	return client.WithAddress("cn.xswitch.node." + nodeUUID)
}

func (channel *Channel) Save() *Channel {
	return WriteChannel(channel.GetUuid(), channel)
}

func (channel *Channel) FullCtrlUuid() string {
	return fmt.Sprintf("cn.xswitch.ctrl.%s", channel.CtrlUuid)
}

// GetChannelEvent .
func (channel *Channel) GetChannelEvent() *xctrl.ChannelEvent {
	return &channel.ChannelEvent
}

// NodeAddress 生成NODE地址
func (channel *Channel) NodeAddress() client.CallOption {
	return client.WithAddress("cn.xswitch.node." + channel.GetNodeUuid())
}

// Answer 应答
func (channel *Channel) Answer() *xctrl.Response {
	response, err := Service().Answer(context.TODO(), &xctrl.Request{
		CtrlUuid: UUID(),
		Uuid:     channel.GetUuid(),
	}, channel.NodeAddress())

	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}

	return response
}

// Accept 接管
func (channel *Channel) Accept(takeover bool) *xctrl.Response {
	response, err := Service().Accept(context.TODO(), &xctrl.AcceptRequest{
		CtrlUuid: UUID(),
		Uuid:     channel.GetUuid(),
		Takeover: takeover,
	}, channel.NodeAddress())

	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}

	return response
}

// Ready 判断通道是否正常状态
func (channel *Channel) Ready() bool {
	if channel == nil {
		return false
	}
	c, err := ReadChannel(channel.GetUuid())
	if err != nil || c == nil {
		return false
	}

	if channel.State == `DESTROY` {
		return false
	}
	return true
	// response := c.GetStates()
	// return response.GetReady()
}

// GetVariable 获取通道变量
func (channel *Channel) GetVariable(key string) string {
	if channel == nil {
		return ""
	}
	channel.lock.RLock()
	if channel.Params != nil {
		if v, ok := channel.Params[key]; ok {
			channel.lock.RUnlock()
			return v
		}
	}
	channel.lock.RUnlock()
	return ""
}

// SetVariable 保存通道变量
func (channel *Channel) SetVariable(key, value string) error {
	if channel == nil {
		return fmt.Errorf("Unable to locate Channel")
	}
	channel.lock.Lock()
	if channel.Params == nil {
		channel.Params = make(map[string]string)
	}
	channel.Params[key] = value
	channel.lock.Unlock()

	response := channel.SetVar(&xctrl.SetVarRequest{
		Uuid: channel.GetUuid(),
		Data: map[string]string{
			key: value,
		},
	})
	if response.GetCode() != 200 {
		return fmt.Errorf("[%d]%s", response.GetCode(), response.GetMessage())
	}
	return nil
}

// SetVariables 保存多个通道变量
func (channel *Channel) SetVariables(varKv map[string]string) error {
	if channel == nil {
		return fmt.Errorf("Unable to locate Channel")
	}
	channel.lock.Lock()
	data := make(map[string]string)
	if channel.Params == nil {
		channel.Params = make(map[string]string)
	}
	for k, v := range varKv {
		channel.Params[k] = v
		data[k] = v
	}

	channel.lock.Unlock()

	response := channel.SetVar(&xctrl.SetVarRequest{
		Uuid: channel.GetUuid(),
		Data: data,
	})
	if response.GetCode() != 200 {
		return fmt.Errorf("[%d]%s", response.GetCode(), response.GetMessage())
	}
	return nil
}

// Play 播放一个文件，默认超时时间1小时
func (channel *Channel) Play(req *xctrl.PlayRequest) *xctrl.Response {
	return channel.PlayWithTimeout(req, 1*time.Hour)
}

// PlayWithTimeout 播放一个文件，可传入超时时间
func (channel *Channel) PlayWithTimeout(req *xctrl.PlayRequest, timeout time.Duration) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().Play(context.Background(), req, client.WithRequestTimeout(timeout), channel.NodeAddress())
	if err != nil {
		responseErr := errors.Parse(err.Error())
		if response == nil {
			response = &xctrl.Response{
				Code:    responseErr.Code,
				Message: responseErr.Detail,
			}
		} else {
			response.Code = responseErr.Code
			response.Message = responseErr.Detail
		}
	}

	if response.Code < 200 || response.Code > 300 {
		if response.Code < 500 {
			fmt.Printf("%s Play %s error: %d %s", channel.GetUuid(), req.Media.Data, response.Code, response.GetMessage())
		} else {
			fmt.Errorf("%s Play %s error: %d %s", channel.GetUuid(), req.Media.Data, response.Code, response.GetMessage())
		}
	}
	return response
}

// Broadcast 播放多个文件
func (channel *Channel) Broadcast(req *xctrl.BroadcastRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().Broadcast(
		context.Background(),
		req,
		channel.NodeAddress(),
		client.WithRequestTimeout(60*time.Minute),
	)

	if err != nil {
		responseErr := errors.Parse(err.Error())
		response.Code = responseErr.Code
		response.Message = responseErr.Detail
	}
	return response
}

// Stop 停止当前正在执行的API
func (channel *Channel) Stop() *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	req := &xctrl.StopRequest{Uuid: channel.GetUuid()}
	response, err := Service().Stop(context.Background(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// Hangup 挂机
func (channel *Channel) Hangup(cause string, flag xctrl.HangupRequest_HangupFlag) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}
	req := &xctrl.HangupRequest{
		Uuid:  channel.GetUuid(),
		Cause: cause,
		Flag:  flag,
	}

	response, err := Service().Hangup(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}

	return response
}

// UnBridge 停止bridge
func (channel *Channel) UnBridge(req *xctrl.Request) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().UnBridge(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}

	return response
}

// UnBridge2 断开桥接，park channel
func (channel *Channel) UnBridge2(req *xctrl.Request) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().UnBridge2(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}

	return response
}

// Hold 呼叫保持
func (channel *Channel) Hold(req *xctrl.HoldRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().Hold(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// Transfer 转移（待定）
func (channel *Channel) Transfer(req *xctrl.TransferRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().Transfer(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// ThreeWay 三方通话
func (channel *Channel) ThreeWay(req *xctrl.ThreeWayRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().ThreeWay(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// Intercept 强插
func (channel *Channel) Intercept(req *xctrl.InterceptRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().Intercept(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// Consult 协商转移
func (channel *Channel) Consult(req *xctrl.ConsultRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().Consult(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// NativeAPI native Api
func (channel *Channel) NativeAPI(req *xctrl.NativeRequest) *xctrl.NativeResponse {
	if channel == nil {
		return &xctrl.NativeResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().NativeAPI(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// NativeApp native App
func (channel *Channel) NativeApp(req *xctrl.NativeRequest) *xctrl.NativeResponse {
	if channel == nil {
		return &xctrl.NativeResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}
	response, err := Service().NativeApp(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.NativeResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// NativeJSAPI Native js api
func (channel *Channel) NativeJSAPI(req *xctrl.NativeJSRequest) *xctrl.NativeJSResponse {
	if channel == nil {
		return &xctrl.NativeJSResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	response, err := Service().NativeJSAPI(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.NativeJSResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// GetChannelData 获取通道数据
func (channel *Channel) GetChannelData(req *xctrl.GetChannelDataRequest) *xctrl.ChannelDataResponse {
	if channel == nil {
		return &xctrl.ChannelDataResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().GetChannelData(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.ChannelDataResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// ReadDTMF 读取DTMF按键
func (channel *Channel) ReadDTMF(req *xctrl.DTMFRequest) *xctrl.DTMFResponse {
	if channel == nil {
		return &xctrl.DTMFResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	req.Uuid = channel.GetUuid()
	response, err := Service().ReadDTMF(context.Background(), req,
		channel.NodeAddress(), client.WithRequestTimeout(60*time.Minute))
	if err != nil {
		response = new(xctrl.DTMFResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// ReadDTMF 收集DTMF按键，支持正则和错误音提示
func (channel *Channel) ReadDigits(req *xctrl.DigitsRequest) *xctrl.DigitsResponse {
	if channel == nil {
		return &xctrl.DigitsResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	req.Uuid = channel.GetUuid()
	response, err := Service().ReadDigits(context.Background(), req,
		channel.NodeAddress(), client.WithRequestTimeout(60*time.Minute))
	if err != nil {
		response = new(xctrl.DigitsResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// 设置通道变量
func (channel *Channel) SetVar(req *xctrl.SetVarRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().SetVar(context.Background(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// GetVar 获取通道变量
func (channel *Channel) GetVar(req *xctrl.GetVarRequest) *xctrl.VarResponse {
	if channel == nil {
		return &xctrl.VarResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().GetVar(context.TODO(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.VarResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// GetStates 获取通道状态
func (channel *Channel) GetStates() *xctrl.StateResponse {
	if channel == nil {
		return &xctrl.StateResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	response, err := Service().GetState(context.Background(), &xctrl.GetStateRequest{
		Uuid: channel.GetUuid(),
	}, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.StateResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// Dial 外呼
func (channel *Channel) Dial(req *xctrl.DialRequest) *xctrl.DialResponse {
	if channel == nil {
		return &xctrl.DialResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}
	if req.GetCtrlUuid() == "" {
		req.CtrlUuid = UUID()
	}
	response, err := Service().Dial(context.TODO(), req, channel.NodeAddress(), client.WithRequestTimeout(2*time.Minute))
	if err != nil {
		response = new(xctrl.DialResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

// Bridge 在把当前呼叫桥接（发起）另一个呼叫
func (channel *Channel) Bridge(req *xctrl.BridgeRequest, async bool) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	if !async {
		response, err := Service().Bridge(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))
		if err != nil {
			response = new(xctrl.Response)
			e := errors.Parse(err.Error())
			response.Code = e.Code
			response.Message = e.Detail
		}
		return response
	}
	response, err := AsyncService().Bridge(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))
	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response

}

// ChannelBridge 桥接两个呼叫
func (channel *Channel) ChannelBridge(req *xctrl.ChannelBridgeRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := AsyncService().ChannelBridge(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))

	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

func (channel *Channel) SetMute(req *xctrl.MuteRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetCtrlUuid() == "" {
		req.Uuid = channel.GetUuid()
	}
	response, err := Service().Mute(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))

	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}

	return response
}

// Record 录音
func (channel *Channel) Record(req *xctrl.RecordRequest) *xctrl.RecordResponse {
	if channel == nil {
		return &xctrl.RecordResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	// Record
	response, err := Service().Record(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))

	if err != nil {
		response = new(xctrl.RecordResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
		return response
	}
	return response
}

// SendDTMF 发送DTMF
func (channel *Channel) SendDTMF(dtmf string) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}
	// Record
	response, err := Service().SendDTMF(context.Background(), &xctrl.SendDTMFRequest{
		Uuid: channel.GetUuid(),
		Dtmf: dtmf,
	}, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))

	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
		return response
	}
	return response
}

// DetectSpeech 语音识别
func (channel *Channel) DetectSpeech(req *xctrl.DetectRequest, async bool) *xctrl.DetectResponse {
	if channel == nil {
		return &xctrl.DetectResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	var response *xctrl.DetectResponse
	var err error

	if !async {
		response, err = Service().DetectSpeech(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))
	} else {
		response, err = AsyncService().DetectSpeech(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))
	}

	if err != nil {
		response = new(xctrl.DetectResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
		return response
	}
	return response
}

// RingBackDetection 回铃音检测
func (channel *Channel) RingBackDetection(req *xctrl.RingBackDetectionRequest, async bool) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	var response *xctrl.Response
	var err error

	if !async {
		response, err = Service().RingBackDetection(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))
	} else {
		response, err = AsyncService().RingBackDetection(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))
	}

	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
		return response
	}
	return response
}

// DetectFace 人脸识别
func (channel *Channel) DetectFace(req *xctrl.DetectFaceRequest) *xctrl.Response {
	if channel == nil {
		return &xctrl.Response{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}

	if req.GetUuid() == "" {
		req.Uuid = channel.GetUuid()
	}

	response, err := Service().DetectFace(context.Background(), req, channel.NodeAddress(), client.WithRequestTimeout(24*time.Hour))

	if err != nil {
		response = new(xctrl.Response)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
		return response
	}
	return response
}

// String marshalIndent
func (channel *Channel) String() string {
	if channel == nil {
		return ""
	}
	channel.lock.RLock()
	jsonBytes, _ := json.MarshalIndent(channel, "", "  ")
	channel.lock.RUnlock()
	return string(jsonBytes)
}

// Marshal marshal to JSON
func (channel *Channel) Marshal() []byte {
	if channel == nil {
		return []byte{}
	}
	channel.lock.RLock()
	jsonBytes, _ := json.Marshal(channel)
	channel.lock.RUnlock()
	return jsonBytes
}

func (channel *Channel) Subscribe(topic string, cb nats.EventCallback, queue string) (nats.Subscriber, error) {
	if globalCtrl == nil {
		return nil, fmt.Errorf("ctrl uninitialized")
	}

	if topic == "" {
		topic = "cn.xswitch.ctrl." + channel.CtrlUuid
	}

	sub, err := globalCtrl.Subscribe(topic, cb, queue)

	if err != nil {
		return nil, err
	}

	channel.subs = append(channel.subs, sub)

	return sub, nil
}

//FIFO 呼叫中心FIFO队列（先入先出）
func (channel *Channel) FIFO(req *xctrl.FIFORequest) *xctrl.FIFOResponse {
	if channel == nil {
		return &xctrl.FIFOResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}
	response, err := Service().FIFO(context.Background(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.FIFOResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

//Callcenter 呼叫中心Callcenter
func (channel *Channel) Callcenter(req *xctrl.CallcenterRequest) *xctrl.CallcenterResponse {
	if channel == nil {
		return &xctrl.CallcenterResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}
	response, err := Service().Callcenter(context.Background(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.CallcenterResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

//Conference 会议
func (channel *Channel) Conference(req *xctrl.ConferenceRequest) *xctrl.ConferenceResponse {
	if channel == nil {
		return &xctrl.ConferenceResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}
	response, err := Service().Conference(context.Background(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.ConferenceResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

//AI
func (channel *Channel) AI(req *xctrl.AIRequest) *xctrl.AIResponse {
	if channel == nil {
		return &xctrl.AIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}
	response, err := Service().AI(context.Background(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.AIResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}

func (channel *Channel) HttAPI(req *xctrl.HttAPIRequest) *xctrl.HttAPIResponse {
	if channel == nil {
		return &xctrl.HttAPIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unable to locate Channel",
		}
	}
	response, err := Service().HttAPI(context.Background(), req, channel.NodeAddress())
	if err != nil {
		response = new(xctrl.HttAPIResponse)
		e := errors.Parse(err.Error())
		response.Code = e.Code
		response.Message = e.Detail
	}
	return response
}
