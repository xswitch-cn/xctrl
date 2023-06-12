/* TBoy 测试框架，模拟XSWITCH COPYRIGHT 2021 烟台小樱桃 ALL RIGHTS RESERVED x-y-t.cn
 */

package tboy

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetReportCaller(true)
}

// RequestParam .
type RequestParam struct {
	UUID     string `json:"uuid"`
	NodeIP   string `json:"node_ip"`
	NodeUUID string `json:"node_uuid"`

	CDR *CallDetailRecord `json:"cdr" validate:"required"`
}

// CallDetailRecord 话单详情
type CallDetailRecord struct {
	UUID                        string `json:"uuid" validate:"required"`       // 话单UUID
	Domain                      string `json:"xcc_domain"`                     // 域
	Mark                        string `json:"xcc_mark"`                       // 号码方向
	UID                         string `json:"xcc_uid"`                        // 用户UID
	Context                     string `json:"context"`                        // context
	Billsec                     int    `json:"billsec,string"`                 // 通话时长
	CallerIDName                string `json:"caller_id_name"`                 // 主叫名称
	CallerIDNumber              string `json:"caller_id_number"`               // 主叫号码
	DestinationNumber           string `json:"destination_number"`             // 被叫号码
	OriginCidNumber             string `json:"xcc_origin_cid_number"`          // 原始主叫号码
	OriginDestNumber            string `json:"xcc_origin_dest_number"`         // 原始被叫号码
	OriginOutboundCidNumber     string `json:"xcc_origin_outbound_cid_number"` // 原始外显号码
	StationType                 string `json:"xcc_stationtype"`                // SIP | PSTN
	Direction                   string `json:"direction"`                      // 逻辑方向
	Duration                    int    `json:"duration,string"`                // 花费时长
	HangupCause                 string `json:"hangup_cause"`                   // 挂机原因
	PeerUUID                    string `json:"peer_uuid"`                      // 对端UUID
	SipToHost                   string `json:"sip_to_host"`                    //
	SipFromHost                 string `json:"sip_from_host"`                  //
	SipDisplay                  string `json:"sip_from_display"`               //
	SipNetworkAddr              string `json:"sip_local_network_addr"`         //
	SipNetworkPort              int    `json:"sip_network_port,string"`        //
	SipHangupDisposition        string `json:"sip_hangup_disposition,omitempty"`
	SofiaProfileName            string `json:"sofia_profile_name"`               // default 分机号码, public 外线号码, interconnect 转移
	SofiaProfileURL             string `json:"sofia_profile_url"`                //
	StartStamp                  string `json:"start_stamp"`                      // 开始时间
	AnswerStamp                 string `json:"answer_stamp"`                     // 接听时间
	EndStamp                    string `json:"end_stamp"`                        // 结束时间
	Leg                         string `json:"leg"`                              //
	ServingSide                 string `json:"cc_side"`                          //
	ServingAgentUUID            string `json:"cc_agent"`                         // 坐席UUID
	ServingAgentName            string `json:"cc_agent_name"`                    // 坐席名称
	ServingAgentScore           int    `json:"cc_agent_rating_score,string"`     // 评分
	ServingAgentSession         string `json:"cc_agent_session_uuid"`            //
	ServingAgentEmployeeNumber  string `json:"cc_agent_employee_number"`         // 工号
	ServingQueueUUID            string `json:"cc_queue"`                         // 队列UUID
	ServingQueueName            string `json:"cc_queue_name"`                    // 队列名称
	ServingQueueJoinedEpoch     int64  `json:"cc_queue_joined_epoch,string"`     // 加入队列时间
	ServingQueueAnsweredEpoch   int64  `json:"cc_queue_answered_epoch,string"`   // 坐席应答时间
	ServingQueueTerminatedEpoch int64  `json:"cc_queue_terminated_epoch,string"` // 坐席挂机时间
	ServingContact              string `json:"cc_contact"`                       //
	MemberUUID                  string `json:"cc_member_uuid"`                   // 成员UUID
	MemberSession               string `json:"cc_member_session_uuid"`           //
	Session                     string `json:"xcc_session"`                      // session
}

// Dialplan .
type Dialplan struct {
	UID                string `json:"xcc_uid,omitempty"`
	UUID               string `json:"uuid,omitempty"`
	NodeUUID           string `json:"node_uuid,omitempty"`
	Domain             string `json:"xcc_domain,omitempty"`
	Context            string `json:"context,omitempty"`
	FromHost           string `json:"sip_from_host,omitempty"`
	ToHost             string `json:"sip_to_host,omitempty"`
	Network            string `json:"sip_network,omitempty"`
	Date               string `json:"date_local,omitempty"`
	Direction          string `json:"caller_direction,omitempty"`
	CallerName         string `json:"caller_name,omitempty"`
	CallerNumber       string `json:"caller_number,omitempty"`
	DestinationNumber  string `json:"destination_number,omitempty"`
	RoutingTag         string `json:"xcc_routing_tag,omitempty"`
	RouteUuid          string `json:"xcc_route_uuid,omitempty"`
	RedlistExecChecked string `json:"xcc_redlist_exec_checked,omitempty"` // 判断是否已经执行过红名单

	OriginCidNumber  string // 内部转发带不过来，需从 session 中取
	OriginDestNumber string // 内部转发带不过来，需从 session 中取

	Session string `json:"xcc_session,omitempty"`
}

// App .
type App struct {
	App  string `json:"app,omitempty"`
	Data string `json:"data,omitempty"`
}

const (
	send_cdr     = true
	cdr_template = `{
				 "hangup_cause": "NORMAL_CLEARING",
				 "caller_id_number":     "10000200",
				 "duration":     "29",
				 "xcc_uid":      "913240a1-e6af-4ae1-8116-a6aeb26efc37",
				 "context":      "default",
				 "cc_queue_joined_epoch":        "1619621138",
				 "uuid": "5c29b3fb-b75a-4dbd-a33c-e4cb8c9af829",
				 "cc_agent_session_uuid":        "16473c52-295f-4fb3-83ce-f40f03cc125e",
				 "sip_from_display":     "10016",
				 "start_stamp":  "2021-04-28 21:56:17",
				 "end_stamp":    "2021-04-28 21:56:46",
				 "sip_to_host":  "154.8.164.96",
				 "sip_term_status":      "200",
				 "sofia_profile_name":   "public",
				 "sip_local_network_addr":       "192.168.0.111",
				 "xcc_session":  "612ee4b4-4fb9-4797-9459-53840e4a99d0",
				 "sofia_profile_url":    "sip:mod_sofia@192.168.0.111:17080",
				 "sip_from_host":        "xcc.xswitch.cn",
				 "caller_id_name":       "10000200",
				 "cc_member_session_uuid":       "5c29b3fb-b75a-4dbd-a33c-e4cb8c9af829",
				 "cc_side":      "infopd",
				 "answer_stamp": "2021-04-28 21:56:17",
				 "billsec":      "29",
				 "destination_number":   "10000200",
				 "direction":    "outbound",
				 "sip_network_port":     "20003",
				 "cc_agent_found":       "true",
				 "xcc_domain":   "test.test",
				 "sip_hangup_disposition":       "send_bye",
				 "leg":  "a",
				 "logical_direction":    "outbound"
		 }`
)

type FakeChannel struct {
	CtrlUuid string
	PeerUuid string
	Lock     sync.RWMutex
	Data     *xctrl.ChannelEvent
	Context  context.Context
	Cancel   context.CancelFunc
}

type Options struct {
	PeerWait   int
	PeerAnswer bool
	PeerReject bool
	AcdAssign  bool
	ActualPlay bool
}

// TBoy .
type TBoy struct {
	NodeUUID string
	Domain   string
	*Options
}

type OptionFn func(*Options)

func OptionPeerAnswer(answer bool) OptionFn {
	return func(o *Options) {
		o.PeerAnswer = answer
	}
}

func OptionPeerWait(wait int) OptionFn {
	return func(o *Options) {
		o.PeerWait = wait
	}
}

func OptionPeerReject(reject bool) OptionFn {
	return func(o *Options) {
		o.PeerReject = reject
	}
}
func OptionAcdAssign(assgin bool) OptionFn {
	return func(o *Options) {
		o.AcdAssign = assgin
	}
}
func OptionActualPlay(play bool) OptionFn {
	return func(o *Options) {
		o.ActualPlay = play
	}
}

func (boy *TBoy) Init() {
	Channels = map[string]*FakeChannel{}
}

func (boy *TBoy) SetUUID(uuid string) {
	boy.NodeUUID = uuid
}

func (boy *TBoy) SetDomain(domain string) {
	boy.Domain = domain
}

var Channels map[string]*FakeChannel

func CacheChannel(uuid string, channel *FakeChannel) {
	Channels[uuid] = channel
}

func (boy *TBoy) OK(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	res := xctrl.Response{
		Code:     200,
		Message:  "OK",
		NodeUuid: boy.NodeUUID,
		Uuid:     request.Uuid,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) Error(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	res := xctrl.Response{
		Code:     404,
		Message:  "Unsupported Method " + msg.Method,
		NodeUuid: boy.NodeUUID,
		Uuid:     request.Uuid,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) DialError(ctx context.Context, msg *ctrl.Message, reply string, code int32) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	cause := "NORMAL_CLEARING"
	switch code {
	case 404:
		cause = "NO_ROUTE_DESTINATION"
	case 486:
		cause = "USER_BUSY"
	}
	res := xctrl.DialResponse{
		Code:     code,
		Message:  "Dial Error",
		NodeUuid: boy.NodeUUID,
		Uuid:     request.Uuid,
		Cause:    cause,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) Accept(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.AcceptRequest
	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	channel, ok := Channels[request.Uuid]

	if ok {
		channel.CtrlUuid = request.CtrlUuid
	}

	boy.OK(ctx, msg, reply)
}

func (boy *TBoy) Answer(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request
	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	channel, ok := Channels[request.Uuid]

	if ok {
		channel.CtrlUuid = request.CtrlUuid
	}

	boy.OK(ctx, msg, reply)
}

func (boy *TBoy) Stop(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	res := xctrl.Response{
		Code:     200,
		Message:  "OK",
		NodeUuid: boy.NodeUUID,
		Uuid:     request.Uuid,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")
	log.Infof("reply to %s %s", reply, result_bytes)
	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) SetVar(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.SetVarRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	channel, ok := Channels[request.Uuid]

	var res xctrl.Response
	if ok {
		channel.Lock.Lock()
		for k, v := range request.Data {
			channel.Data.Params[k] = v
		}
		channel.Lock.Unlock()
		res = xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: boy.NodeUUID,
			Uuid:     request.Uuid,
		}
	} else {
		res = xctrl.Response{
			Code:     http.StatusBadRequest,
			Message:  "can not locate session",
			NodeUuid: boy.NodeUUID,
			Uuid:     request.Uuid,
		}
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}
func (boy *TBoy) GetVar(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.SetVarRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	channel, ok := Channels[request.Uuid]

	resMap := make(map[string]string)
	var res xctrl.VarResponse
	if ok {
		channel.Lock.Lock()
		for k, _ := range request.Data {
			resMap[k] = channel.Data.Params[k]
		}
		channel.Lock.Unlock()
		res = xctrl.VarResponse{
			Code:     200,
			Message:  "OK",
			NodeUuid: boy.NodeUUID,
			Uuid:     request.Uuid,
			Data:     resMap,
		}
	}

	res = xctrl.VarResponse{
		Code:     http.StatusBadRequest,
		Message:  "can not locate session",
		NodeUuid: boy.NodeUUID,
		Uuid:     request.Uuid,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) ReadDTMF(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.SetVarRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	channel, ok := Channels[request.Uuid]

	if ok {
		channel.Lock.Lock()
		for k, v := range request.Data {
			channel.Data.Params[k] = v
		}
		channel.Lock.Unlock()
	}

	res := xctrl.DTMFResponse{
		Code:       200,
		Message:    "OK",
		NodeUuid:   boy.NodeUUID,
		Uuid:       request.Uuid,
		Dtmf:       "1",
		Terminator: "#",
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) Dial(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.DialRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	controller := "cn.xswitch.ctrl." + request.CtrlUuid

	uuid := request.Destination.CallParams[0].Uuid

	call_params_params := request.Destination.CallParams[0].Params

	if call_params_params == nil {
		call_params_params = make(map[string]string)
	}

	channelEvent := xctrl.ChannelEvent{
		NodeUuid:    boy.NodeUUID,
		Uuid:        uuid,
		Direction:   "outbound",
		State:       "CALLING",
		CidName:     "TEST",
		CidNumber:   request.Destination.CallParams[0].CidNumber,
		DestNumber:  request.Destination.CallParams[0].DestNumber,
		AnswerEpoch: uint32(time.Now().Unix()),
		Answered:    true,
		CreateEpoch: uint32(time.Now().Unix()),
		Params:      call_params_params,
	}

	channelEvent.Params["xcc_session"] = channelEvent.Params["sip_h_X-FS-Session"]
	channel := &FakeChannel{
		CtrlUuid: request.CtrlUuid,
		Data:     &channelEvent,
	}
	CacheChannel(uuid, channel)

	event_req := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channelEvent),
	}

	req_str, _ := json.MarshalIndent(event_req, "", "  ")
	ctrl.Publish(controller, req_str)

	channelEvent.State = "ANSWERED"
	event_req.Params = ctrl.ToRawMessage(channelEvent)
	req_str, _ = json.MarshalIndent(event_req, "", "  ")
	ctrl.Publish(controller, req_str)

	channelEvent.State = "READY"
	event_req.Params = ctrl.ToRawMessage(channelEvent)
	req_str, _ = json.MarshalIndent(event_req, "", "  ")
	ctrl.Publish(controller, req_str)

	channel = &FakeChannel{
		CtrlUuid: request.CtrlUuid,
		Data:     &channelEvent,
	}
	if channel.Context == nil {
		channel.Context, channel.Cancel = context.WithCancel(context.Background())
	}

	CacheChannel(uuid, channel)

	if channelEvent.GetParams()["xcc_action"] == "" {
		go func() {
			t := time.NewTimer(3 * time.Second)
			select {
			case <-channel.Context.Done():
				return
			case <-t.C:
				controller := "cn.xswitch.ctrl." + channel.CtrlUuid

				channelEvent.State = "DESTROY"
				channelEvent.Cause = "NORMAL_CLEARING"
				event_req.Params = ctrl.ToRawMessage(channelEvent)
				req_str, _ = json.MarshalIndent(event_req, "", "  ")
				ctrl.Publish(controller, req_str)

				cdr := new(CallDetailRecord)
				err = json.Unmarshal([]byte(cdr_template), cdr)

				if err != nil {
					log.Error(err)
				}

				if send_cdr {

					cdr.UUID = channel.Data.Uuid
					cdr.Leg = "a"
					cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					cdr.CallerIDName = channel.Data.CidName
					cdr.CallerIDNumber = channel.Data.CidNumber
					cdr.DestinationNumber = channel.Data.DestNumber
					cdr.Session = channel.Data.Params["sip_h_X-FS-Session"]

					params := &RequestParam{
						NodeUUID: "test-node",
						NodeIP:   "127.0.0.1",
						UUID:     channel.Data.Uuid,
						CDR:      cdr,
					}

					rpc := *&ctrl.Request{
						Version: "2.0",
						Method:  "Event.CDR",
						Params:  ctrl.ToRawMessage(params),
					}

					str, _ := json.MarshalIndent(rpc, "", "  ")
					ctrl.Publish("cn.xswitch.event.cdr", str)
				}
			}

		}()
	}

	// response
	res := xctrl.Response{
		Code:     200,
		Message:  "OK",
		NodeUuid: boy.NodeUUID,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) Bridge(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.BridgeRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	channel, ok := Channels[request.Uuid]

	if !ok {
		log.Error("Invalid Channel", request.Uuid)
	}

	ctrl_uuid := request.CtrlUuid

	if ctrl_uuid == "" {
		ctrl_uuid = channel.CtrlUuid
	}

	controller := "cn.xswitch.ctrl." + ctrl_uuid

	if len(request.Destination.CallParams) < 1 {
		log.Error(request)
		return
	}

	cparams := request.Destination.CallParams[0]
	call_params_params := request.Destination.CallParams[0].Params

	if call_params_params == nil {
		call_params_params = make(map[string]string)
	}

	bleg := xctrl.ChannelEvent{
		NodeUuid:    boy.NodeUUID,
		Uuid:        cparams.Uuid,
		Direction:   "outbound",
		State:       "CALLING",
		CidName:     "TEST",
		CidNumber:   cparams.CidNumber,
		DestNumber:  cparams.DestNumber,
		Answered:    false,
		CreateEpoch: uint32(time.Now().Unix()),
		Params:      call_params_params,
	}

	bleg.Params["xcc_session"] = bleg.Params["sip_h_X-FS-Session"]

	event_req := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(bleg),
	}

	req_str, _ := json.MarshalIndent(event_req, "", "  ")
	ctrl.Publish(controller, req_str)

	time.Sleep(time.Second * time.Duration(boy.PeerWait))
	//拒接
	if boy.PeerReject {
		aleg := channel.Data
		controller := "cn.xswitch.ctrl." + channel.CtrlUuid
		aleg.State = "DESTROY"
		aleg.Cause = "NORMAL_CLEARING"
		event_req.Params = ctrl.ToRawMessage(aleg)
		req_str, _ = json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)

		cdr := new(CallDetailRecord)
		err = json.Unmarshal([]byte(cdr_template), cdr)

		if err != nil {
			log.Error(err)
		}

		if send_cdr {
			cdr.UUID = aleg.Uuid
			cdr.Leg = "a"
			cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
			if aleg.AnswerEpoch != 0 {
				cdr.AnswerStamp = time.Unix(int64(aleg.AnswerEpoch), 0).Format("2006-01-02 15:03:04")
			}
			cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
			cdr.CallerIDName = aleg.CidName
			cdr.CallerIDNumber = aleg.CidNumber
			cdr.DestinationNumber = aleg.DestNumber
			cdr.Session = aleg.Params["sip_h_X-FS-Session"]

			params := &RequestParam{
				NodeUUID: "test-node",
				NodeIP:   "127.0.0.1",
				UUID:     aleg.Uuid,
				CDR:      cdr,
			}

			rpc := *&ctrl.Request{
				Version: "2.0",
				Method:  "Event.CDR",
				Params:  ctrl.ToRawMessage(params),
			}

			str, _ := json.MarshalIndent(rpc, "", "  ")
			ctrl.Publish("cn.xswitch.event.cdr", str)
		}
		go func() {

			acontroller := "cn.xswitch.ctrl." + request.CtrlUuid

			bleg.State = "DESTROY"
			bleg.Cause = "NORMAL_CLEARING"
			event_req.Params = ctrl.ToRawMessage(bleg)
			req_str, _ = json.MarshalIndent(event_req, "", "  ")
			ctrl.Publish(acontroller, req_str)

			if send_cdr {
				cdr := new(CallDetailRecord)
				err = json.Unmarshal([]byte(cdr_template), cdr)

				if err != nil {
					log.Error(err)
				}

				cdr.UUID = bleg.Uuid
				cdr.Leg = "b"
				cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
				cdr.AnswerStamp = ""
				cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
				cdr.CallerIDName = bleg.CidName
				cdr.CallerIDNumber = bleg.CidNumber
				cdr.DestinationNumber = bleg.DestNumber
				cdr.Session = bleg.Params["sip_h_X-FS-Session"]

				params := &RequestParam{
					NodeUUID: "test-node",
					NodeIP:   "127.0.0.1",
					UUID:     channel.Data.Uuid,
					CDR:      cdr,
				}

				rpc := *&ctrl.Request{
					Version: "2.0",
					Method:  "Event.CDR",
					Params:  ctrl.ToRawMessage(params),
				}

				str, _ := json.MarshalIndent(rpc, "", "  ")
				ctrl.Publish("cn.xswitch.event.cdr", str)
			}
		}()
	} else {
		//接听
		if boy.PeerAnswer {
			bleg.State = "ANSWERED"
			bleg.AnswerEpoch = uint32(time.Now().Unix())
			event_req.Params = ctrl.ToRawMessage(bleg)
			req_str, _ = json.MarshalIndent(event_req, "", "  ")
			ctrl.Publish(controller, req_str)

			fake_channel := &FakeChannel{
				CtrlUuid: request.CtrlUuid,
				Data:     &bleg,
			}

			CacheChannel(bleg.Uuid, fake_channel)

			aleg := channel.Data

			aleg.State = "BRIDGED"
			event_req.Params = ctrl.ToRawMessage(aleg)
			req_str, _ = json.MarshalIndent(event_req, "", "  ")
			ctrl.Publish(controller, req_str)

			bleg.State = "BRIDGED"
			event_req.Params = ctrl.ToRawMessage(bleg)
			req_str, _ = json.MarshalIndent(event_req, "", "  ")
			ctrl.Publish(controller, req_str)

			go func() {
				time.Sleep(3 * time.Second)

				controller := "cn.xswitch.ctrl." + channel.CtrlUuid

				aleg.State = "UNBRIDGE"
				event_req.Params = ctrl.ToRawMessage(aleg)
				req_str, _ = json.MarshalIndent(event_req, "", "  ")
				ctrl.Publish(controller, req_str)

				aleg.State = "DESTROY"
				aleg.Cause = "NORMAL_CLEARING"
				event_req.Params = ctrl.ToRawMessage(aleg)
				req_str, _ = json.MarshalIndent(event_req, "", "  ")
				ctrl.Publish(controller, req_str)

				cdr := new(CallDetailRecord)
				err = json.Unmarshal([]byte(cdr_template), cdr)

				if err != nil {
					log.Error(err)
				}

				if send_cdr {
					cdr.UUID = aleg.Uuid
					cdr.Leg = "a"
					cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					if aleg.AnswerEpoch != 0 {
						cdr.AnswerStamp = time.Unix(int64(aleg.AnswerEpoch), 0).Format("2006-01-02 15:03:04")
					}
					cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					cdr.CallerIDName = aleg.CidName
					cdr.CallerIDNumber = aleg.CidNumber
					cdr.DestinationNumber = aleg.DestNumber
					cdr.Session = aleg.Params["sip_h_X-FS-Session"]

					params := &RequestParam{
						NodeUUID: "test-node",
						NodeIP:   "127.0.0.1",
						UUID:     aleg.Uuid,
						CDR:      cdr,
					}

					rpc := *&ctrl.Request{
						Version: "2.0",
						Method:  "Event.CDR",
						Params:  ctrl.ToRawMessage(params),
					}

					str, _ := json.MarshalIndent(rpc, "", "  ")
					ctrl.Publish("cn.xswitch.event.cdr", str)
				}
			}()

			go func() {
				time.Sleep(3 * time.Second)

				acontroller := "cn.xswitch.ctrl." + request.CtrlUuid

				bleg.State = "UNBRIDGE"
				event_req.Params = ctrl.ToRawMessage(bleg)
				req_str, _ = json.MarshalIndent(event_req, "", "  ")
				ctrl.Publish(acontroller, req_str)

				bleg.State = "DESTROY"
				bleg.Cause = "NORMAL_CLEARING"
				event_req.Params = ctrl.ToRawMessage(bleg)
				req_str, _ = json.MarshalIndent(event_req, "", "  ")
				ctrl.Publish(acontroller, req_str)

				if send_cdr {
					cdr := new(CallDetailRecord)
					err = json.Unmarshal([]byte(cdr_template), cdr)

					if err != nil {
						log.Error(err)
					}

					cdr.UUID = bleg.Uuid
					cdr.Leg = "b"
					cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					if aleg.AnswerEpoch != 0 {
						cdr.AnswerStamp = time.Unix(int64(aleg.AnswerEpoch), 0).Format("2006-01-02 15:03:04")
					}
					cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					cdr.CallerIDName = bleg.CidName
					cdr.CallerIDNumber = bleg.CidNumber
					cdr.DestinationNumber = bleg.DestNumber
					cdr.Session = bleg.Params["sip_h_X-FS-Session"]

					params := &RequestParam{
						NodeUUID: "test-node",
						NodeIP:   "127.0.0.1",
						UUID:     channel.Data.Uuid,
						CDR:      cdr,
					}

					rpc := *&ctrl.Request{
						Version: "2.0",
						Method:  "Event.CDR",
						Params:  ctrl.ToRawMessage(params),
					}

					str, _ := json.MarshalIndent(rpc, "", "  ")
					ctrl.Publish("cn.xswitch.event.cdr", str)
				}
			}()
		} else {
			time.Sleep(time.Second * 60)

			aleg := channel.Data

			go func() {

				controller := "cn.xswitch.ctrl." + channel.CtrlUuid

				aleg.State = "DESTROY"
				aleg.Cause = "NORMAL_CLEARING"
				event_req.Params = ctrl.ToRawMessage(aleg)
				req_str, _ = json.MarshalIndent(event_req, "", "  ")
				ctrl.Publish(controller, req_str)

				cdr := new(CallDetailRecord)
				err = json.Unmarshal([]byte(cdr_template), cdr)

				if err != nil {
					log.Error(err)
				}

				if send_cdr {
					cdr.UUID = aleg.Uuid
					cdr.Leg = "a"
					cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					if aleg.AnswerEpoch != 0 {
						cdr.AnswerStamp = time.Unix(int64(aleg.AnswerEpoch), 0).Format("2006-01-02 15:03:04")
					}
					cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					cdr.CallerIDName = aleg.CidName
					cdr.CallerIDNumber = aleg.CidNumber
					cdr.DestinationNumber = aleg.DestNumber
					cdr.Session = aleg.Params["sip_h_X-FS-Session"]

					params := &RequestParam{
						NodeUUID: "test-node",
						NodeIP:   "127.0.0.1",
						UUID:     aleg.Uuid,
						CDR:      cdr,
					}

					rpc := *&ctrl.Request{
						Version: "2.0",
						Method:  "Event.CDR",
						Params:  ctrl.ToRawMessage(params),
					}

					str, _ := json.MarshalIndent(rpc, "", "  ")
					ctrl.Publish("cn.xswitch.event.cdr", str)
				}
			}()

			go func() {
				acontroller := "cn.xswitch.ctrl." + request.CtrlUuid

				bleg.State = "DESTROY"
				bleg.Cause = "NORMAL_CLEARING"
				event_req.Params = ctrl.ToRawMessage(bleg)
				req_str, _ = json.MarshalIndent(event_req, "", "  ")
				ctrl.Publish(acontroller, req_str)

				if send_cdr {
					cdr := new(CallDetailRecord)
					err = json.Unmarshal([]byte(cdr_template), cdr)

					if err != nil {
						log.Error(err)
					}

					cdr.UUID = bleg.Uuid
					cdr.Leg = "b"
					cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					cdr.AnswerStamp = ""
					cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
					cdr.CallerIDName = bleg.CidName
					cdr.CallerIDNumber = bleg.CidNumber
					cdr.DestinationNumber = bleg.DestNumber
					cdr.Session = bleg.Params["sip_h_X-FS-Session"]

					params := &RequestParam{
						NodeUUID: "test-node",
						NodeIP:   "127.0.0.1",
						UUID:     channel.Data.Uuid,
						CDR:      cdr,
					}

					rpc := *&ctrl.Request{
						Version: "2.0",
						Method:  "Event.CDR",
						Params:  ctrl.ToRawMessage(params),
					}

					str, _ := json.MarshalIndent(rpc, "", "  ")
					ctrl.Publish("cn.xswitch.event.cdr", str)
				}
			}()
		}

	}

	// response
	res := xctrl.Response{
		Code:     200,
		Message:  "OK",
		NodeUuid: boy.NodeUUID,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) ChannelBridge(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.ChannelBridgeRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}
	aleg, ok := Channels[request.Uuid]

	if !ok {
		log.Error("Invalid Channel", request.Uuid)
	}

	ctrl_uuid := request.CtrlUuid

	if ctrl_uuid == "" {
		ctrl_uuid = aleg.CtrlUuid
	}

	acontroller := "cn.xswitch.ctrl." + ctrl_uuid

	aleg.Data.PeerUuid = request.PeerUuid
	aleg.Data.State = "BRIDGE"

	event_req := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(aleg.Data),
	}

	req_str, _ := json.MarshalIndent(event_req, "", "  ")
	ctrl.Publish(acontroller, req_str)
	log.Infof("a leg:%s b leg:%s ", aleg.Data.Uuid, aleg.Data.PeerUuid)

	bleg, ok := Channels[aleg.Data.PeerUuid]
	if !ok {
		log.Info("Invalid b Channel", request.Uuid)
	} else {
		bcontroller := "cn.xswitch.ctrl." + bleg.CtrlUuid

		bleg.Data.PeerUuid = request.Uuid
		bleg.Data.State = "BRIDGE"
		event_req.Params = ctrl.ToRawMessage(bleg.Data)
		req_str, _ = json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(bcontroller, req_str)
	}

	go func() {
		time.Sleep(3 * time.Second)

		aleg.Data.State = "UNBRIDGE"
		event_req.Params = ctrl.ToRawMessage(aleg.Data)
		req_str, _ = json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(acontroller, req_str)

		aleg.Data.State = "DESTROY"
		aleg.Data.Cause = "NORMAL_CLEARING"
		event_req.Params = ctrl.ToRawMessage(aleg.Data)
		req_str, _ = json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(acontroller, req_str)

		if aleg.Cancel != nil {
			aleg.Cancel()
		}
		if ok {
			bleg.Data.State = "UNBRIDGE"
			event_req.Params = ctrl.ToRawMessage(bleg.Data)
			req_str, _ = json.MarshalIndent(event_req, "", "  ")
			ctrl.Publish(acontroller, req_str)

			bleg.Data.State = "DESTROY"
			bleg.Data.Cause = "NORMAL_CLEARING"
			event_req.Params = ctrl.ToRawMessage(bleg.Data)
			req_str, _ = json.MarshalIndent(event_req, "", "  ")
			ctrl.Publish(acontroller, req_str)
			if bleg.Cancel != nil {
				bleg.Cancel()
			}
		}
		if send_cdr {

			cdr := new(CallDetailRecord)
			err = json.Unmarshal([]byte(cdr_template), cdr)

			if err != nil {
				log.Error(err)
			}

			cdr.UUID = aleg.Data.Uuid
			cdr.PeerUUID = aleg.Data.PeerUuid
			cdr.Leg = "a"
			cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
			cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
			cdr.CallerIDName = aleg.Data.CidName
			cdr.CallerIDNumber = aleg.Data.CidNumber
			cdr.DestinationNumber = aleg.Data.DestNumber

			params := &RequestParam{
				NodeUUID: "test-node",
				NodeIP:   "127.0.0.1",
				UUID:     aleg.Data.Uuid,
				CDR:      cdr,
			}

			rpc := *&ctrl.Request{
				Version: "2.0",
				Method:  "Event.CDR",
				Params:  ctrl.ToRawMessage(params),
			}

			str, _ := json.MarshalIndent(rpc, "", "  ")
			ctrl.Publish("cn.xswitch.event.cdr", str)
			if ok {
				cdr := new(CallDetailRecord)
				err = json.Unmarshal([]byte(cdr_template), cdr)

				if err != nil {
					log.Error(err)
				}

				cdr.UUID = bleg.Data.Uuid
				cdr.PeerUUID = bleg.Data.PeerUuid
				cdr.Leg = "b"
				cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
				cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
				cdr.CallerIDName = bleg.Data.CidName
				cdr.CallerIDNumber = bleg.Data.CidNumber
				cdr.DestinationNumber = bleg.Data.DestNumber

				params := &RequestParam{
					NodeUUID: "test-node",
					NodeIP:   "127.0.0.1",
					UUID:     bleg.Data.Uuid,
					CDR:      cdr,
				}

				rpc := *&ctrl.Request{
					Version: "2.0",
					Method:  "Event.CDR",
					Params:  ctrl.ToRawMessage(params),
				}

				str, _ := json.MarshalIndent(rpc, "", "  ")
				ctrl.Publish("cn.xswitch.event.cdr", str)
			}
		}
	}()

	// response
	res := xctrl.Response{
		Code:     200,
		Message:  "OK",
		NodeUuid: boy.NodeUUID,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) Record(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.RecordRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}
	_, ok := Channels[request.Uuid]
	if !ok {
		res := xctrl.Response{
			Code:     http.StatusBadRequest,
			Message:  "can not locate session",
			NodeUuid: boy.NodeUUID,
			Uuid:     request.Uuid,
		}
		res_bytes, _ := json.MarshalIndent(res, "", "  ")
		raw := json.RawMessage(res_bytes)

		result := ctrl.Result{
			Version: "2.0",
			ID:      msg.ID,
			Result:  &raw,
		}

		result_bytes, _ := json.MarshalIndent(result, "", "  ")
		if reply == "" {
			reply = "cn.xswitch.ctrl." + request.CtrlUuid
		}
		ctrl.Publish(reply, result_bytes)
		return
	}

	res := xctrl.RecordEvent{
		NodeUuid: boy.NodeUUID,
		Uuid:     request.Uuid,
		Action:   request.Action,
		Path:     request.Path,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) Hangup(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}
	channel, ok := Channels[request.Uuid]
	if !ok {
		res := xctrl.Response{
			Code:     http.StatusBadRequest,
			Message:  "can not locate session",
			NodeUuid: boy.NodeUUID,
			Uuid:     request.Uuid,
		}
		res_bytes, _ := json.MarshalIndent(res, "", "  ")
		raw := json.RawMessage(res_bytes)

		result := ctrl.Result{
			Version: "2.0",
			ID:      msg.ID,
			Result:  &raw,
		}

		result_bytes, _ := json.MarshalIndent(result, "", "  ")
		if reply == "" {
			reply = "cn.xswitch.ctrl." + request.CtrlUuid
		}
		ctrl.Publish(reply, result_bytes)
		return
	}
	if channel.Cancel != nil {
		channel.Cancel()
	}

	boy.OK(ctx, msg, reply)

	controller := "cn.xswitch.ctrl." + channel.CtrlUuid

	channel.Data.State = "DESTROY"
	channel.Data.Cause = "NORMAL_CLEARING"
	event_req := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channel.Data),
	}
	req_str, _ := json.MarshalIndent(event_req, "", "  ")
	ctrl.Publish(controller, req_str)

	cdr := new(CallDetailRecord)
	err = json.Unmarshal([]byte(cdr_template), cdr)

	if err != nil {
		log.Error(err)
	}

	if send_cdr {

		cdr.UUID = channel.Data.Uuid
		cdr.Leg = "b"
		cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
		cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
		cdr.CallerIDName = channel.Data.CidName
		cdr.CallerIDNumber = channel.Data.CidNumber
		cdr.DestinationNumber = channel.Data.DestNumber
		cdr.Session = channel.Data.Params["sip_h_X-FS-Session"]

		params := &RequestParam{
			NodeUUID: "test-node",
			NodeIP:   "127.0.0.1",
			UUID:     channel.Data.Uuid,
			CDR:      cdr,
		}

		rpc := *&ctrl.Request{
			Version: "2.0",
			Method:  "Event.CDR",
			Params:  ctrl.ToRawMessage(params),
		}

		str, _ := json.MarshalIndent(rpc, "", "  ")
		ctrl.Publish("cn.xswitch.event.cdr", str)
	}
}

func (boy *TBoy) NativeAPI(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.NativeRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	switch request.Cmd {
	case "status":
		boy.NativeAPIStatus(ctx, msg, reply)
	default:
		boy.Error(ctx, msg, reply)
	}

}

func (boy *TBoy) NativeAPIStatus(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}
	data := "UP 0 years, 0 days, 0 hours, 9 minutes, 10 seconds, 674 milliseconds, 555 microseconds\nFreeSWITCH (Version 1.10.8-dev git 89ac59d 2022-03-30 04:05:36Z 64bit) is ready\n0 session(s) since startup\n0 session(s) - peak 0, last 5min 0 \n0 session(s) per Sec out of max 30, peak 0, last 5min 0 \n1000 session(s) max\nmin idle cpu 0.00/96.33\nCurrent Stack Size/Max 240K/8192K\n"

	res := xctrl.NativeResponse{
		Code:     200,
		Message:  "OK",
		NodeUuid: boy.NodeUUID,
		Data:     data,
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) NativeJSAPI(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.XNativeJSRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}
	log.Info(request)

	switch request.Data.Command {
	case "status":
		boy.NativeJSAPIStatus(ctx, msg, reply)
	case "sofia.status":
		log.Errorf("Unsupported Method: %s", msg.Method)
	default:
		boy.Error(ctx, msg, reply)
	}

}

func (boy *TBoy) NativeJSAPIStatus(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}
	dataMap := map[string]interface{}{
		"systemStatus": "ready",
		"uptime": map[string]int{
			"years":        0,
			"days":         0,
			"hours":        0,
			"minutes":      9,
			"seconds":      10,
			"milliseconds": 679,
			"microseconds": 389,
		},
		"version": "1.10.8-dev git 89ac59d 2022-03-30 04:05:36Z 64bit",
		"sessions": map[string]interface{}{
			"count": map[string]int{
				"total":    0,
				"active":   0,
				"peak":     0,
				"peak5Min": 0,
				"limit":    1000,
			},
			"rate": map[string]int{
				"current":  0,
				"max":      30,
				"peak":     0,
				"peak5Min": 0,
			},
		},
		"idleCPU": map[string]interface{}{
			"used":    0,
			"allowed": 96.333333333333329,
		},
		"stackSizeKB": map[string]int{
			"current": 240,
			"max":     8192,
		},
	}

	res := xctrl.XNativeJSResponse{
		Code:     200,
		Message:  "OK",
		NodeUuid: boy.NodeUUID,
		Data:     ctrl.ToRawMessage(dataMap),
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, result_bytes)
}

func (boy *TBoy) Event(msg *ctrl.Message, natsEvent nats.Event) {
	topic := natsEvent.Topic()
	reply := natsEvent.Reply()
	log.Infof("%s %s", topic, msg.Method)

	if msg.Method == "" && msg.Result != nil {
		log.Infof("Got a response: %s", msg.ID)
		return
	}

	ctx := context.Background()

	switch msg.Method {
	case "XNode.Accept":
		boy.Accept(ctx, msg, reply)
	case "XNode.SetVar":
		boy.SetVar(ctx, msg, reply)
	case "XNode.GetVar":
		boy.GetVar(ctx, msg, reply)
	case "XNode.Answer":
		boy.Answer(ctx, msg, reply)
	case "XNode.Hangup":
		boy.Hangup(ctx, msg, reply)

	case "XNode.NativeApp":
		go func() {
			time.Sleep(200 * time.Millisecond)
			boy.OK(ctx, msg, reply)
		}()
	case "XNode.NativeAPI":
		boy.NativeAPI(ctx, msg, reply)
	case "XNode.NativeJSAPI":
		boy.NativeJSAPI(ctx, msg, reply)
	case "XNode.Play":
		go func() {
			var playRequest xctrl.PlayRequest
			err := json.Unmarshal(*msg.Params, &playRequest)
			if boy.Options.ActualPlay && err == nil && runtime.GOOS == "darwin" {
				if playRequest.Media.Type == "TEXT" {
					cmd := exec.Command("say", "-v", "Ting-Ting", playRequest.Media.Data)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						log.Fatal(err)
					}
				} else if strings.HasPrefix(playRequest.Media.Data, "https://xswitch.cn") {
					cmd := exec.Command("wget", "--quiet", "-O", "/tmp/test.wav", playRequest.Media.Data)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						log.Fatal(err)
					}
					cmd = exec.Command("play", "/tmp/test.wav")
					err = cmd.Run()
					if err != nil {
						log.Fatal(err)
					}
				}
			} else {
				time.Sleep(900 * time.Millisecond)
			}
			boy.OK(ctx, msg, reply)
		}()
	case "XNode.Stop":
		boy.Stop(ctx, msg, reply)
	case "XNode.Record":
		boy.Record(ctx, msg, reply)
	case "XNode.ReadDTMF":
		boy.ReadDTMF(ctx, msg, reply)
	case "XNode.Dial":
		boy.Dial(ctx, msg, reply)
	case "XNode.Bridge":
		boy.Bridge(ctx, msg, reply)
	case "XNode.ChannelBridge":
		boy.ChannelBridge(ctx, msg, reply)
	default:
		log.Errorf("Unsupported Method: %s", msg.Method)
		boy.Error(ctx, msg, reply)
	}
}

// we don't use this method in node side
func (boy *TBoySimple) ChannelEvent(ctx context.Context, channel *ctrl.Channel) {
}

// trimSuffix 清除 topic 后缀
func trimSuffix(source string) string {
	return strings.TrimSuffix(source, `.`+ctrl.UUID())
}
