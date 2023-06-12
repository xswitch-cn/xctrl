package ctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/xctrl/client"
	"git.xswitch.cn/xswitch/xctrl/xctrl/errors"
	"github.com/google/uuid"

	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
)

// Channel call channel
type Channel struct {
	*xctrl.ChannelEvent                   // the parent ChannelEvent
	CtrlUuid            string            // the Controller UUID
	lock                sync.RWMutex      // a Mutex to protect the Channel, internal use only
	subs                []nats.Subscriber // todo
	natsEvent           nats.Event        // the original natsEvent received
	userData            interface{}       // store private userData from the higher level Application
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

func NewChannelEvent() *Channel {
	return &Channel{
		ChannelEvent: &xctrl.ChannelEvent{},
	}
}

func (channel *Channel) GetNatsEvent() nats.Event {
	return channel.natsEvent
}

func (channel *Channel) Save() *Channel {
	return WriteChannel(channel.GetUuid(), channel)
}

func (channel *Channel) FullCtrlUuid() string {
	return fmt.Sprintf("cn.xswitch.ctrl.%s", channel.CtrlUuid)
}

// GetChannelEvent .
func (channel *Channel) GetChannelEvent() *xctrl.ChannelEvent {
	return channel.ChannelEvent
}

// NodeAddress 生成NODE地址
func (channel *Channel) NodeAddress() client.CallOption {
	return client.WithAddress("cn.xswitch.node." + channel.GetNodeUuid())
}

// Answer 应答
func (channel *Channel) Answer0(opts ...client.CallOption) *xctrl.Response {
	response := channel.Answer(&xctrl.AnswerRequest{
		CtrlUuid: UUID(),
		Uuid:     channel.GetUuid(),
	}, opts...)

	return response
}

// Answer 应答
func (channel *Channel) AnswerWithChannelParams(channel_params []string, opts ...client.CallOption) *xctrl.Response {
	response := channel.Answer(&xctrl.AnswerRequest{
		CtrlUuid:      UUID(),
		Uuid:          channel.GetUuid(),
		ChannelParams: channel_params,
	}, opts...)

	return response
}

// Accept 接管
func (channel *Channel) Accept0(opts ...client.CallOption) *xctrl.Response {
	response := channel.Accept(&xctrl.AcceptRequest{
		CtrlUuid: UUID(),
		Uuid:     channel.GetUuid(),
	}, opts...)

	return response
}

// Accept 接管
func (channel *Channel) AcceptAndTakeOver(opts ...client.CallOption) *xctrl.Response {
	response := channel.Accept(&xctrl.AcceptRequest{
		CtrlUuid: UUID(),
		Uuid:     channel.GetUuid(),
		Takeover: true,
	}, opts...)

	return response
}

// Accept 接管
func (channel *Channel) AcceptWithChannelParams(channel_params []string, opts ...client.CallOption) *xctrl.Response {
	response := channel.Accept(&xctrl.AcceptRequest{
		CtrlUuid:      UUID(),
		Uuid:          channel.GetUuid(),
		ChannelParams: channel_params,
	}, opts...)

	return response
}

// Accept 接管
func (channel *Channel) AcceptWithChannelParamsAndTakeOver(channel_params []string, opts ...client.CallOption) *xctrl.Response {
	response := channel.Accept(&xctrl.AcceptRequest{
		CtrlUuid:      UUID(),
		Uuid:          channel.GetUuid(),
		ChannelParams: channel_params,
		Takeover:      true,
	}, opts...)

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
func (channel *Channel) GetVariable0(key string) string {
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
func (channel *Channel) SetVariable0(key, value string) error {
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
func (channel *Channel) SetVariables(vars map[string]string) error {
	if channel == nil {
		return fmt.Errorf("Unable to locate Channel")
	}
	channel.lock.Lock()
	data := make(map[string]string)
	if channel.Params == nil {
		channel.Params = make(map[string]string)
	}
	for k, v := range vars {
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

func (channel *Channel) PlayFile(file string, opts ...client.CallOption) *xctrl.Response {
	media := &xctrl.Media{
		Data: file,
	}
	req := &xctrl.PlayRequest{
		Media: media,
	}
	return channel.Play(req, opts...)
}

func (channel *Channel) PlayTTS(engine string, voice string, text string, opts ...client.CallOption) *xctrl.Response {
	media := &xctrl.Media{
		Type:   "TEXT",
		Data:   text,
		Engine: engine,
		Voice:  voice,
	}
	req := &xctrl.PlayRequest{
		Uuid:  channel.Uuid,
		Media: media,
	}
	return channel.Play(req, opts...)
}

// Play 播放一个文件，默认超时时间1小时
func (channel *Channel) Play0(req *xctrl.PlayRequest) *xctrl.Response {
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

// Stop 停止当前正在执行的API
func (channel *Channel) Stop0() *xctrl.Response {
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
func (channel *Channel) Hangup0(cause string, flag xctrl.HangupRequest_HangupFlag) *xctrl.Response {
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

// Bridge 在把当前呼叫桥接（发起）另一个呼叫
func (channel *Channel) Bridge0(req *xctrl.BridgeRequest, async bool) *xctrl.Response {
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

// SendDTMF 发送DTMF
func (channel *Channel) SendDTMF0(dtmf string) *xctrl.Response {
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
func (channel *Channel) DetectSpeech0(req *xctrl.DetectRequest, async bool) *xctrl.DetectResponse {
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
func (channel *Channel) RingBackDetection0(req *xctrl.RingBackDetectionRequest, async bool) *xctrl.Response {
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

func (channel *Channel) Subscribe(subject string, cb nats.EventCallback, queue string) (nats.Subscriber, error) {
	if globalCtrl == nil {
		return nil, fmt.Errorf("ctrl uninitialized")
	}

	if subject == "" {
		return nil, fmt.Errorf("no subject specified")
	}

	sub, err := globalCtrl.Subscribe(subject, cb, queue)

	if err != nil {
		return nil, err
	}

	channel.subs = append(channel.subs, sub)

	return sub, nil
}

func (channel *Channel) SetUserData(userData interface{}) {
	if channel == nil {
		return
	}
	channel.userData = userData
}

func (channel *Channel) GetUserData() interface{} {
	return channel.userData
}
