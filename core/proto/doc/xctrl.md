# XSwitch XCC Proto Buffer 协议参考文档

<a name="top"></a>
<a name="user-content-top"></a>

这是[XCC API文档](https://docs.xswitch.cn/xcc-api/)的协议参考，使用[Google Protocol Buffers](https://protobuf.dev/)描述。

本文档只是对具体协议数据格式及类型的参考说明，详细的字段说明和用法请参考[XCC API列表](https://docs.xswitch.cn/xcc-api/api/)，原始的`.proto`文件可以在[proto](../)相关目录中找到。

## 目录

- [xctrl.proto](#xctrl.proto)
  - [AIRequest](#xctrl.AIRequest)
  - [AIResponse](#xctrl.AIResponse)
  - [AcceptRequest](#xctrl.AcceptRequest)
  - [Action](#xctrl.Action)
  - [Application](#xctrl.Application)
  - [BridgeRequest](#xctrl.BridgeRequest)
  - [BroadcastRequest](#xctrl.BroadcastRequest)
  - [CallParam](#xctrl.CallParam)
  - [CallcenterRequest](#xctrl.CallcenterRequest)
  - [CallcenterResponse](#xctrl.CallcenterResponse)
  - [ChannelBridgeRequest](#xctrl.ChannelBridgeRequest)
  - [ChannelData](#xctrl.ChannelData)
  - [ChannelDataResponse](#xctrl.ChannelDataResponse)
  - [ChannelEvent](#xctrl.ChannelEvent)
  - [ConferenceInfoRequest](#xctrl.ConferenceInfoRequest)
  - [ConferenceInfoRequestData](#xctrl.ConferenceInfoRequestData)
  - [ConferenceInfoRequestDataData](#xctrl.ConferenceInfoRequestDataData)
  - [ConferenceInfoResponse](#xctrl.ConferenceInfoResponse)
  - [ConferenceInfoResponseConference](#xctrl.ConferenceInfoResponseConference)
  - [ConferenceInfoResponseData](#xctrl.ConferenceInfoResponseData)
  - [ConferenceInfoResponseFlags](#xctrl.ConferenceInfoResponseFlags)
  - [ConferenceInfoResponseMembers](#xctrl.ConferenceInfoResponseMembers)
  - [ConferenceInfoResponseVariables](#xctrl.ConferenceInfoResponseVariables)
  - [ConferenceRequest](#xctrl.ConferenceRequest)
  - [ConferenceResponse](#xctrl.ConferenceResponse)
  - [ConsultRequest](#xctrl.ConsultRequest)
  - [Ctrl](#xctrl.Ctrl)
  - [DTMFRequest](#xctrl.DTMFRequest)
  - [DTMFResponse](#xctrl.DTMFResponse)
  - [Destination](#xctrl.Destination)
  - [DetectFaceRequest](#xctrl.DetectFaceRequest)
  - [DetectRequest](#xctrl.DetectRequest)
  - [DetectResponse](#xctrl.DetectResponse)
  - [DetectedData](#xctrl.DetectedData)
  - [DetectedFaceEvent](#xctrl.DetectedFaceEvent)
  - [DialRequest](#xctrl.DialRequest)
  - [DialResponse](#xctrl.DialResponse)
  - [DigitsRequest](#xctrl.DigitsRequest)
  - [DigitsResponse](#xctrl.DigitsResponse)
  - [Echo2Request](#xctrl.Echo2Request)
  - [EngineData](#xctrl.EngineData)
  - [FIFORequest](#xctrl.FIFORequest)
  - [FIFOResponse](#xctrl.FIFOResponse)
  - [GetChannelDataRequest](#xctrl.GetChannelDataRequest)
  - [GetStateRequest](#xctrl.GetStateRequest)
  - [GetVarRequest](#xctrl.GetVarRequest)
  - [HangupRequest](#xctrl.HangupRequest)
  - [Header](#xctrl.Header)
  - [HoldRequest](#xctrl.HoldRequest)
  - [HttAPIRequest](#xctrl.HttAPIRequest)
  - [HttAPIResponse](#xctrl.HttAPIResponse)
  - [InterceptRequest](#xctrl.InterceptRequest)
  - [JStatusIdleCPU](#xctrl.JStatusIdleCPU)
  - [JStatusRequest](#xctrl.JStatusRequest)
  - [JStatusRequest.JStatusData](#xctrl.JStatusRequest.JStatusData)
  - [JStatusResponse](#xctrl.JStatusResponse)
  - [JStatusResponseData](#xctrl.JStatusResponseData)
  - [JStatusSessions](#xctrl.JStatusSessions)
  - [JStatusSessionsCount](#xctrl.JStatusSessionsCount)
  - [JStatusSessionsRate](#xctrl.JStatusSessionsRate)
  - [JStatusStackSize](#xctrl.JStatusStackSize)
  - [JStatusUptime](#xctrl.JStatusUptime)
  - [LuaRequest](#xctrl.LuaRequest)
  - [LuaResponse](#xctrl.LuaResponse)
  - [Media](#xctrl.Media)
  - [MuteRequest](#xctrl.MuteRequest)
  - [NativeJSRequest](#xctrl.NativeJSRequest)
  - [NativeJSResponse](#xctrl.NativeJSResponse)
  - [NativeRequest](#xctrl.NativeRequest)
  - [NativeResponse](#xctrl.NativeResponse)
  - [Node](#xctrl.Node)
  - [Payload](#xctrl.Payload)
  - [PlayRequest](#xctrl.PlayRequest)
  - [RecordEvent](#xctrl.RecordEvent)
  - [RecordRequest](#xctrl.RecordRequest)
  - [RecordResponse](#xctrl.RecordResponse)
  - [Request](#xctrl.Request)
  - [Response](#xctrl.Response)
  - [RingBackDetectionRequest](#xctrl.RingBackDetectionRequest)
  - [SendDTMFRequest](#xctrl.SendDTMFRequest)
  - [SendINFORequest](#xctrl.SendINFORequest)
  - [SetVarRequest](#xctrl.SetVarRequest)
  - [SpeechRequest](#xctrl.SpeechRequest)
  - [StashResult](#xctrl.StashResult)
  - [StateResponse](#xctrl.StateResponse)
  - [StopDetectRequest](#xctrl.StopDetectRequest)
  - [StopRequest](#xctrl.StopRequest)
  - [ThreeWayRequest](#xctrl.ThreeWayRequest)
  - [TransferRequest](#xctrl.TransferRequest)
  - [VarResponse](#xctrl.VarResponse)
  - [VideoResizeEvent](#xctrl.VideoResizeEvent)

  - [HangupRequest.HangupFlag](#xctrl.HangupRequest.HangupFlag)
  - [MediaType](#xctrl.MediaType)
  - [RecordRequest.RecordAction](#xctrl.RecordRequest.RecordAction)


  - [XNode](#xctrl.XNode)


- [Scalar Value Types](#scalar-value-types)



<a name="user-content-xctrl.proto"/>
<a name="xctrl.proto"/>
<p align="right"><a href="#top">Top</a></p>

## xctrl.proto



<a name="user-content-xctrl.AIRequest"/>
<a name="xctrl.AIRequest"/>

### AIRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| url | [string](#string) |  |  |
| data | [map<string, string>](#map-string-string) |  |  |






<a name="user-content-xctrl.AIResponse"/>
<a name="xctrl.AIResponse"/>

### AIResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.AcceptRequest"/>
<a name="xctrl.AcceptRequest"/>

### AcceptRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  | Controller UUID |
| uuid | [string](#string) |  | optional, Channel UUID |
| takeover | [bool](#bool) |  | optional, default to false.when true, all subsequest events will be delivered to the new controller if already controlled by other controller, otherwise it will fail |






<a name="user-content-xctrl.Action"/>
<a name="xctrl.Action"/>

### Action



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| owner_uid | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| param | [CallParam](#xctrl.CallParam) |  |  |






<a name="user-content-xctrl.Application"/>
<a name="xctrl.Application"/>

### Application



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| app | [string](#string) |  |  |
| data | [string](#string) |  |  |






<a name="user-content-xctrl.BridgeRequest"/>
<a name="xctrl.BridgeRequest"/>

### BridgeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| ringback | [string](#string) |  |  |
| flow_control | [string](#string) |  | NONE | CALLER | CALLEE | ANY |
| continue_on_fail | [string](#string) |  | true | false | comma separated freeswitch causes e.g. USER_BUSY,NO_ANSWER |
| destination | [Destination](#xctrl.Destination) |  |  |






<a name="user-content-xctrl.BroadcastRequest"/>
<a name="xctrl.BroadcastRequest"/>

### BroadcastRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| media | [Media](#xctrl.Media) |  |  |
| option | [string](#string) |  | BOTH, ALEG, BLEG, AHOLDB, BHOLDA |






<a name="user-content-xctrl.CallParam"/>
<a name="xctrl.CallParam"/>

### CallParam



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| cid_name | [string](#string) |  |  |
| cid_number | [string](#string) |  |  |
| dest_number | [string](#string) |  |  |
| dial_string | [string](#string) |  |  |
| timeout | [uint32](#uint32) |  |  |
| max_duration | [uint32](#uint32) |  |  |
| params | [map<string, string>](#map-string-string) |  |  |






<a name="user-content-xctrl.CallcenterRequest"/>
<a name="xctrl.CallcenterRequest"/>

### CallcenterRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |






<a name="user-content-xctrl.CallcenterResponse"/>
<a name="xctrl.CallcenterResponse"/>

### CallcenterResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.ChannelBridgeRequest"/>
<a name="xctrl.ChannelBridgeRequest"/>

### ChannelBridgeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| peer_uuid | [string](#string) |  |  |
| bridge_delay | [int32](#int32) |  |  |
| flow_control | [string](#string) |  | NONE | CALLER | CALLEE | ANY |






<a name="user-content-xctrl.ChannelData"/>
<a name="xctrl.ChannelData"/>

### ChannelData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| variable_cc_queue | [string](#string) |  |  |
| variable_cc_queue_name | [string](#string) |  |  |
| variable_cc_agent_session_uuid | [string](#string) |  |  |
| variable_cc_member_uuid | [string](#string) |  |  |
| variable_xcc_origin_dest_number | [string](#string) |  |  |






<a name="user-content-xctrl.ChannelDataResponse"/>
<a name="xctrl.ChannelDataResponse"/>

### ChannelDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| format | [string](#string) |  |  |
| data | [ChannelData](#xctrl.ChannelData) |  | when format == JSON |






<a name="user-content-xctrl.ChannelEvent"/>
<a name="xctrl.ChannelEvent"/>

### ChannelEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| peer_uuid | [string](#string) |  |  |
| direction | [string](#string) |  | call direction inbound |outbound |
| state | [string](#string) |  | START RINGING ANSWERED ACTIVE DESTROY READY ... |
| cid_name | [string](#string) |  |  |
| cid_number | [string](#string) |  |  |
| dest_number | [string](#string) |  |  |
| create_epoch | [uint32](#uint32) |  |  |
| ring_epoch | [uint32](#uint32) |  |  |
| answer_epoch | [uint32](#uint32) |  |  |
| hangup_epoch | [uint32](#uint32) |  |  |
| peers | [string](#string) | repeated | list of uuids |
| params | [map<string, string>](#map-string-string) |  |  |
| billsec | [string](#string) |  |  |
| duration | [string](#string) |  |  |
| cause | [string](#string) |  |  |
| answered | [bool](#bool) |  |  |
| node_ip | [string](#string) |  |  |
| domain | [string](#string) |  |  |






<a name="user-content-xctrl.ConferenceInfoRequest"/>
<a name="xctrl.ConferenceInfoRequest"/>

### ConferenceInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| data | [ConferenceInfoRequestData](#xctrl.ConferenceInfoRequestData) |  |  |






<a name="user-content-xctrl.ConferenceInfoRequestData"/>
<a name="xctrl.ConferenceInfoRequestData"/>

### ConferenceInfoRequestData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| command | [string](#string) |  |  |
| data | [ConferenceInfoRequestDataData](#xctrl.ConferenceInfoRequestDataData) |  |  |






<a name="user-content-xctrl.ConferenceInfoRequestDataData"/>
<a name="xctrl.ConferenceInfoRequestDataData"/>

### ConferenceInfoRequestDataData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| conferenceName | [string](#string) |  |  |
| showMembers | [bool](#bool) |  |  |
| memberFilters | [map<string, string>](#map-string-string) |  |  |






<a name="user-content-xctrl.ConferenceInfoResponse"/>
<a name="xctrl.ConferenceInfoResponse"/>

### ConferenceInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| seq | [string](#string) |  | optional |
| data | [ConferenceInfoResponseData](#xctrl.ConferenceInfoResponseData) |  |  |






<a name="user-content-xctrl.ConferenceInfoResponseConference"/>
<a name="xctrl.ConferenceInfoResponseConference"/>

### ConferenceInfoResponseConference



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| conference_name | [string](#string) |  |  |
| member_count | [int32](#int32) |  |  |
| ghost_count | [int32](#int32) |  |  |
| rate | [int32](#int32) |  |  |
| run_time | [int32](#int32) |  |  |
| conference_uuid | [string](#string) |  |  |
| canvas_count | [int32](#int32) |  |  |
| max_bw_in | [int32](#int32) |  |  |
| force_bw_in | [int32](#int32) |  |  |
| video_floor_packets | [int32](#int32) |  |  |
| locked | [bool](#bool) |  |  |
| destruct | [bool](#bool) |  |  |
| wait_mod | [bool](#bool) |  |  |
| audio_always | [bool](#bool) |  |  |
| running | [bool](#bool) |  |  |
| answered | [bool](#bool) |  |  |
| enforce_min | [bool](#bool) |  |  |
| bridge_to | [bool](#bool) |  |  |
| dynamic | [bool](#bool) |  |  |
| exit_sound | [bool](#bool) |  |  |
| enter_sound | [bool](#bool) |  |  |
| recording | [bool](#bool) |  |  |
| video_bridge | [bool](#bool) |  |  |
| video_floor_only | [bool](#bool) |  |  |
| video_rfc4579 | [bool](#bool) |  |  |
| variables | [ConferenceInfoResponseVariables](#xctrl.ConferenceInfoResponseVariables) |  |  |
| members | [ConferenceInfoResponseMembers](#xctrl.ConferenceInfoResponseMembers) | repeated |  |






<a name="user-content-xctrl.ConferenceInfoResponseData"/>
<a name="xctrl.ConferenceInfoResponseData"/>

### ConferenceInfoResponseData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| conference | [ConferenceInfoResponseConference](#xctrl.ConferenceInfoResponseConference) |  |  |






<a name="user-content-xctrl.ConferenceInfoResponseFlags"/>
<a name="xctrl.ConferenceInfoResponseFlags"/>

### ConferenceInfoResponseFlags



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| can_hear | [bool](#bool) |  |  |
| can_see | [bool](#bool) |  |  |
| can_speak | [bool](#bool) |  |  |
| hold | [bool](#bool) |  |  |
| mute_detect | [bool](#bool) |  |  |
| talking | [bool](#bool) |  |  |
| has_video | [bool](#bool) |  |  |
| video_bridge | [bool](#bool) |  |  |
| has_floor | [bool](#bool) |  |  |
| is_moderator | [bool](#bool) |  |  |
| end_conference | [bool](#bool) |  |  |






<a name="user-content-xctrl.ConferenceInfoResponseMembers"/>
<a name="xctrl.ConferenceInfoResponseMembers"/>

### ConferenceInfoResponseMembers



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  |  |
| id | [int32](#int32) |  |  |
| uuid | [string](#string) |  |  |
| caller_id_name | [string](#string) |  |  |
| caller_id_number | [string](#string) |  |  |
| join_time | [int32](#int32) |  |  |
| last_talking | [int32](#int32) |  |  |
| energy | [int32](#int32) |  |  |
| volume_in | [int32](#int32) |  |  |
| volume_out | [int32](#int32) |  |  |
| output_volume | [int32](#int32) |  |  |
| input_volume | [int32](#int32) |  |  |
| flags | [ConferenceInfoResponseFlags](#xctrl.ConferenceInfoResponseFlags) |  |  |






<a name="user-content-xctrl.ConferenceInfoResponseVariables"/>
<a name="xctrl.ConferenceInfoResponseVariables"/>

### ConferenceInfoResponseVariables







<a name="user-content-xctrl.ConferenceRequest"/>
<a name="xctrl.ConferenceRequest"/>

### ConferenceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| profile | [string](#string) |  |  |
| flags | [string](#string) | repeated |  |






<a name="user-content-xctrl.ConferenceResponse"/>
<a name="xctrl.ConferenceResponse"/>

### ConferenceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.ConsultRequest"/>
<a name="xctrl.ConsultRequest"/>

### ConsultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| destination | [Destination](#xctrl.Destination) |  |  |






<a name="user-content-xctrl.Ctrl"/>
<a name="xctrl.Ctrl"/>

### Ctrl



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| ip | [string](#string) |  |  |
| version | [string](#string) |  |  |
| rack | [uint32](#uint32) |  |  |
| address | [string](#string) |  |  |






<a name="user-content-xctrl.DTMFRequest"/>
<a name="xctrl.DTMFRequest"/>

### DTMFRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| min_digits | [uint32](#uint32) |  | optional default = 1 |
| max_digits | [uint32](#uint32) |  | optional default = 1 |
| timeout | [uint32](#uint32) |  | optiona default = 5000ms |
| digit_timeout | [uint32](#uint32) |  | optional default = 2000ms |
| terminators | [string](#string) |  | optional default none, can be 0-9,*,# |
| media | [Media](#xctrl.Media) |  | play or tts |
| max_tries | [uint32](#uint32) |  | not implemented yet |
| regex | [string](#string) |  |  |
| media_invalid | [Media](#xctrl.Media) |  | Media to playback when received DTMF doesn't match the regex |
| play_last_invalid_prompt | [bool](#bool) |  | default false |






<a name="user-content-xctrl.DTMFResponse"/>
<a name="xctrl.DTMFResponse"/>

### DTMFResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| dtmf | [string](#string) |  |  |
| terminator | [string](#string) |  | optional |






<a name="user-content-xctrl.Destination"/>
<a name="xctrl.Destination"/>

### Destination



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ringall | [bool](#bool) |  |  |
| global_params | [map<string, string>](#map-string-string) |  |  |
| call_params | [CallParam](#xctrl.CallParam) | repeated |  |
| channel_params | [string](#string) | repeated |  |






<a name="user-content-xctrl.DetectFaceRequest"/>
<a name="xctrl.DetectFaceRequest"/>

### DetectFaceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| mask | [string](#string) |  |  |
| action | [string](#string) |  | START STOP TEXT CLEAR |
| text | [string](#string) |  | if action == "TEXT" |
| font | [string](#string) |  |  |
| font_size | [string](#string) |  |  |
| fg_color | [string](#string) |  |  |
| bg_color | [string](#string) |  |  |






<a name="user-content-xctrl.DetectRequest"/>
<a name="xctrl.DetectRequest"/>

### DetectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| media | [Media](#xctrl.Media) |  | oneof play or tts |
| dtmf | [DTMFRequest](#xctrl.DTMFRequest) |  | detect dtmf too, optional |
| speech | [SpeechRequest](#xctrl.SpeechRequest) |  | speech params, mandatory |






<a name="user-content-xctrl.DetectResponse"/>
<a name="xctrl.DetectResponse"/>

### DetectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| data | [DetectedData](#xctrl.DetectedData) |  |  |






<a name="user-content-xctrl.DetectedData"/>
<a name="xctrl.DetectedData"/>

### DetectedData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dtmf | [string](#string) |  |  |
| terminator | [string](#string) |  |  |
| text | [string](#string) |  |  |
| confidence | [double](#double) |  |  |
| is_final | [bool](#bool) |  | final or partial result |
| error | [string](#string) |  | when error |
| type | [string](#string) |  | DTMF Speech.Begin Speech.Partial Speech.End ERROR |
| engine | [string](#string) |  | the ASR engine |
| engine_data | [EngineData](#xctrl.EngineData) |  | string or JSON Struct, detailed object returned from ASR engine |
| offset | [uint32](#uint32) |  |  |






<a name="user-content-xctrl.DetectedFaceEvent"/>
<a name="xctrl.DetectedFaceEvent"/>

### DetectedFaceEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| picture | [string](#string) |  |  |
| width | [int32](#int32) |  |  |
| height | [int32](#int32) |  |  |






<a name="user-content-xctrl.DialRequest"/>
<a name="xctrl.DialRequest"/>

### DialRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| destination | [Destination](#xctrl.Destination) |  |  |
| apps | [Application](#xctrl.Application) | repeated |  |






<a name="user-content-xctrl.DialResponse"/>
<a name="xctrl.DialResponse"/>

### DialResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| cause | [string](#string) |  |  |






<a name="user-content-xctrl.DigitsRequest"/>
<a name="xctrl.DigitsRequest"/>

### DigitsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| min_digits | [uint32](#uint32) |  | optional default = 1 |
| max_digits | [uint32](#uint32) |  | optional default = 1 |
| timeout | [uint32](#uint32) |  | optiona default = 5000ms |
| digit_timeout | [uint32](#uint32) |  | optional default = 2000ms |
| terminators | [string](#string) |  | optional default none, can be 0-9,*,# |
| media | [Media](#xctrl.Media) |  | play or tts |
| max_tries | [uint32](#uint32) |  | not implemented yet |
| regex | [string](#string) |  |  |
| media_invalid | [Media](#xctrl.Media) |  | invalid  meida |
| play_last_invalid_prompt | [bool](#bool) |  | default false |






<a name="user-content-xctrl.DigitsResponse"/>
<a name="xctrl.DigitsResponse"/>

### DigitsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| dtmf | [string](#string) |  |  |
| terminator | [string](#string) |  | optional |






<a name="user-content-xctrl.Echo2Request"/>
<a name="xctrl.Echo2Request"/>

### Echo2Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| action | [string](#string) |  | START | STOP |
| direction | [string](#string) |  | SELF | OTHER |






<a name="user-content-xctrl.EngineData"/>
<a name="xctrl.EngineData"/>

### EngineData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [Header](#xctrl.Header) |  |  |
| payload | [Payload](#xctrl.Payload) |  |  |






<a name="user-content-xctrl.FIFORequest"/>
<a name="xctrl.FIFORequest"/>

### FIFORequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| inout | [string](#string) |  |  |
| wait_music | [string](#string) |  |  |
| exit_announce | [string](#string) |  |  |
| priority | [int32](#int32) |  |  |






<a name="user-content-xctrl.FIFOResponse"/>
<a name="xctrl.FIFOResponse"/>

### FIFOResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.GetChannelDataRequest"/>
<a name="xctrl.GetChannelDataRequest"/>

### GetChannelDataRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| format | [string](#string) |  | optional JSON(default) JSONSTR XML TXT LIST |






<a name="user-content-xctrl.GetStateRequest"/>
<a name="xctrl.GetStateRequest"/>

### GetStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.GetVarRequest"/>
<a name="xctrl.GetVarRequest"/>

### GetVarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| data | [string](#string) | repeated |  |






<a name="user-content-xctrl.HangupRequest"/>
<a name="xctrl.HangupRequest"/>

### HangupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| cause | [string](#string) |  | NORMAL_CLEARING USER_BUSY CALL_REJECTED ... |
| flag | [HangupRequest.HangupFlag](#xctrl.HangupRequest.HangupFlag) |  |  |






<a name="user-content-xctrl.Header"/>
<a name="xctrl.Header"/>

### Header



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| namespace | [string](#string) |  |  |
| name | [string](#string) |  |  |
| status | [double](#double) |  |  |
| message_id | [string](#string) |  |  |
| task_id | [string](#string) |  |  |
| status_text | [string](#string) |  |  |






<a name="user-content-xctrl.HoldRequest"/>
<a name="xctrl.HoldRequest"/>

### HoldRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| action | [string](#string) |  | ON OFF TOGGLE |
| display | [string](#string) |  | OPTIONAL only supported by some phones |






<a name="user-content-xctrl.HttAPIRequest"/>
<a name="xctrl.HttAPIRequest"/>

### HttAPIRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| url | [string](#string) |  |  |
| data | [map<string, string>](#map-string-string) |  |  |






<a name="user-content-xctrl.HttAPIResponse"/>
<a name="xctrl.HttAPIResponse"/>

### HttAPIResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.InterceptRequest"/>
<a name="xctrl.InterceptRequest"/>

### InterceptRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| target_uuid | [string](#string) |  |  |






<a name="user-content-xctrl.JStatusIdleCPU"/>
<a name="xctrl.JStatusIdleCPU"/>

### JStatusIdleCPU



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| used | [float](#float) |  |  |
| allowed | [float](#float) |  |  |






<a name="user-content-xctrl.JStatusRequest"/>
<a name="xctrl.JStatusRequest"/>

### JStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| data | [JStatusRequest.JStatusData](#xctrl.JStatusRequest.JStatusData) |  |  |






<a name="user-content-xctrl.JStatusRequest.JStatusData"/>
<a name="xctrl.JStatusRequest.JStatusData"/>

### JStatusRequest.JStatusData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| command | [string](#string) |  |  |
| data | [string](#string) |  |  |






<a name="user-content-xctrl.JStatusResponse"/>
<a name="xctrl.JStatusResponse"/>

### JStatusResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| seq | [string](#string) |  |  |
| data | [JStatusResponseData](#xctrl.JStatusResponseData) |  |  |






<a name="user-content-xctrl.JStatusResponseData"/>
<a name="xctrl.JStatusResponseData"/>

### JStatusResponseData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| systemStatus | [string](#string) |  |  |
| version | [string](#string) |  |  |
| uptime | [JStatusUptime](#xctrl.JStatusUptime) |  |  |
| sessions | [JStatusSessions](#xctrl.JStatusSessions) |  |  |
| idleCPU | [JStatusIdleCPU](#xctrl.JStatusIdleCPU) |  |  |
| stackSizeKB | [JStatusStackSize](#xctrl.JStatusStackSize) |  |  |






<a name="user-content-xctrl.JStatusSessions"/>
<a name="xctrl.JStatusSessions"/>

### JStatusSessions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| count | [JStatusSessionsCount](#xctrl.JStatusSessionsCount) |  |  |
| rate | [JStatusSessionsRate](#xctrl.JStatusSessionsRate) |  |  |






<a name="user-content-xctrl.JStatusSessionsCount"/>
<a name="xctrl.JStatusSessionsCount"/>

### JStatusSessionsCount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int32](#int32) |  |  |
| active | [int32](#int32) |  |  |
| peak | [int32](#int32) |  |  |
| peak5Min | [int32](#int32) |  |  |
| limit | [int32](#int32) |  |  |






<a name="user-content-xctrl.JStatusSessionsRate"/>
<a name="xctrl.JStatusSessionsRate"/>

### JStatusSessionsRate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| current | [int32](#int32) |  |  |
| max | [int32](#int32) |  |  |
| peak | [int32](#int32) |  |  |
| peak5Min | [int32](#int32) |  |  |






<a name="user-content-xctrl.JStatusStackSize"/>
<a name="xctrl.JStatusStackSize"/>

### JStatusStackSize



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| current | [int32](#int32) |  |  |
| max | [int32](#int32) |  |  |






<a name="user-content-xctrl.JStatusUptime"/>
<a name="xctrl.JStatusUptime"/>

### JStatusUptime



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| years | [int32](#int32) |  |  |
| days | [int32](#int32) |  |  |
| hours | [int32](#int32) |  |  |
| minutes | [int32](#int32) |  |  |
| seconds | [int32](#int32) |  |  |
| milliseconds | [int32](#int32) |  |  |
| microseconds | [int32](#int32) |  |  |






<a name="user-content-xctrl.LuaRequest"/>
<a name="xctrl.LuaRequest"/>

### LuaRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| script | [string](#string) |  |  |






<a name="user-content-xctrl.LuaResponse"/>
<a name="xctrl.LuaResponse"/>

### LuaResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.Media"/>
<a name="xctrl.Media"/>

### Media



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  | FILE TEXT SSML |
| data | [string](#string) |  |  |
| engine | [string](#string) |  |  |
| voice | [string](#string) |  |  |
| loop | [uint32](#uint32) |  |  |
| offset | [uint32](#uint32) |  |  |






<a name="user-content-xctrl.MuteRequest"/>
<a name="xctrl.MuteRequest"/>

### MuteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| direction | [string](#string) |  | WRITE, READ, BOTH |
| level | [int32](#int32) |  |  |
| flag | [string](#string) |  | FIRST, LAST |






<a name="user-content-xctrl.NativeJSRequest"/>
<a name="xctrl.NativeJSRequest"/>

### NativeJSRequest
placeholer type, do not use it, use XNativeJSRequest instead






<a name="user-content-xctrl.NativeJSResponse"/>
<a name="xctrl.NativeJSResponse"/>

### NativeJSResponse
placeholer type, do not use it, use XNativeJSResponse instead






<a name="user-content-xctrl.NativeRequest"/>
<a name="xctrl.NativeRequest"/>

### NativeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| cmd | [string](#string) |  |  |
| args | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.NativeResponse"/>
<a name="xctrl.NativeResponse"/>

### NativeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| data | [string](#string) |  |  |
| result | [string](#string) |  |  |






<a name="user-content-xctrl.Node"/>
<a name="xctrl.Node"/>

### Node



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| ip | [string](#string) |  |  |
| version | [string](#string) |  |  |
| rack | [uint32](#uint32) |  |  |
| address | [string](#string) |  |  |
| uptime | [uint32](#uint32) |  | 启动以来秒数 |
| sessions | [uint32](#uint32) |  | 当前Session数 |
| sessions_max | [uint32](#uint32) |  | Session最大阈值 |
| sps_max | [uint32](#uint32) |  | 每秒Session最大阈值 |
| sps_last | [uint32](#uint32) |  | 最后一秒的Session数 |
| sps_last_5min | [uint32](#uint32) |  | 最后5分钟每秒的Session均值 |
| sessions_since_startup | [uint32](#uint32) |  | 开机以来的Session数 |
| session_peak_5min | [uint32](#uint32) |  | 5分钟Session最大值 |
| session_peak_max | [uint32](#uint32) |  | 历史Session最大值 |






<a name="user-content-xctrl.Payload"/>
<a name="xctrl.Payload"/>

### Payload



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| index | [double](#double) |  |  |
| time | [double](#double) |  |  |
| result | [string](#string) |  |  |
| confidence | [double](#double) |  |  |
| words | [string](#string) | repeated |  |
| status | [double](#double) |  |  |
| gender | [string](#string) |  |  |
| begin_time | [double](#double) |  |  |
| stash_result | [StashResult](#xctrl.StashResult) |  |  |
| audio_extra_info | [string](#string) |  |  |
| sentence_id | [string](#string) |  |  |
| gender_score | [double](#double) |  |  |






<a name="user-content-xctrl.PlayRequest"/>
<a name="xctrl.PlayRequest"/>

### PlayRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| media | [Media](#xctrl.Media) |  |  |






<a name="user-content-xctrl.RecordEvent"/>
<a name="xctrl.RecordEvent"/>

### RecordEvent
enum Method {
Invalid = 0;
节点注册
NodeRegister = 1;
节点离线
NodeUnregister = 2;
节点数据更新
NodeUpdate = 3;

通道
Channel = 4;
FreeSWITCH原生消息
Native = 5;

old VCC event
Vcc = 6;
如果API请求有后续事件
Result = 7;
　异步Dial结果
DialResult = 8;
NativeAPI结果
NativeResult = 9;

获取配置信息
FetchXML = 10;
获取Dialplan
Dialplan = 11;

话机消息
NativeEvent = 12;

话单
CDR = 13;
}


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| action | [string](#string) |  | START STOP |
| path | [string](#string) |  |  |
| size | [uint32](#uint32) |  |  |
| samples | [uint32](#uint32) |  |  |
| record_ms | [uint32](#uint32) |  |  |
| completion_cause | [string](#string) |  | success-silence |






<a name="user-content-xctrl.RecordRequest"/>
<a name="xctrl.RecordRequest"/>

### RecordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| path | [string](#string) |  |  |
| action | [string](#string) |  | all params are optional |
| limit | [uint32](#uint32) |  |  |
| beep | [string](#string) |  | play a beep before record"default" or TGML https://freeswitch.org/confluence/display/FREESWITCH/TGML |
| terminators | [string](#string) |  |  |
| silence_seconds | [uint32](#uint32) |  |  |
| thresh | [uint32](#uint32) |  | VAD threshold, 0 is disabled, 1~10000 |
| rate | [uint32](#uint32) |  | valid rates are 8000, 16000, 22050, 24000, 32000, 44100, 48000 |






<a name="user-content-xctrl.RecordResponse"/>
<a name="xctrl.RecordResponse"/>

### RecordResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| terminator | [string](#string) |  | if terminated by DTMF |
| path | [string](#string) |  | mirror back of the path |






<a name="user-content-xctrl.Request"/>
<a name="xctrl.Request"/>

### Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.Response"/>
<a name="xctrl.Response"/>

### Response



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |






<a name="user-content-xctrl.RingBackDetectionRequest"/>
<a name="xctrl.RingBackDetectionRequest"/>

### RingBackDetectionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| stop_tone | [string](#string) |  | optional |
| ignore_samples | [string](#string) |  | optional |
| auto_hangup | [bool](#bool) |  | optional default = true |
| answer_auto_stop | [bool](#bool) |  | optional default = true |
| max_detect_time | [uint32](#uint32) |  | optional default = 60 |






<a name="user-content-xctrl.SendDTMFRequest"/>
<a name="xctrl.SendDTMFRequest"/>

### SendDTMFRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| dtmf | [string](#string) |  |  |






<a name="user-content-xctrl.SendINFORequest"/>
<a name="xctrl.SendINFORequest"/>

### SendINFORequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| content_type | [string](#string) |  |  |
| data | [string](#string) |  |  |






<a name="user-content-xctrl.SetVarRequest"/>
<a name="xctrl.SetVarRequest"/>

### SetVarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| data | [map<string, string>](#map-string-string) |  |  |






<a name="user-content-xctrl.SpeechRequest"/>
<a name="xctrl.SpeechRequest"/>

### SpeechRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| engine | [string](#string) |  |  |
| no_input_timeout | [uint32](#uint32) |  |  |
| speech_timeout | [uint32](#uint32) |  |  |
| partial_events | [bool](#bool) |  |  |
| disable_detected_data_event | [bool](#bool) |  |  |
| params | [map<string, string>](#map-string-string) |  |  |
| grammar | [string](#string) |  |  |
| max_speech_timeout | [uint32](#uint32) |  |  |






<a name="user-content-xctrl.StashResult"/>
<a name="xctrl.StashResult"/>

### StashResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sentenceId | [double](#double) |  |  |
| beginTime | [double](#double) |  |  |
| text | [string](#string) |  |  |
| currentTime | [double](#double) |  |  |
| words | [string](#string) | repeated |  |






<a name="user-content-xctrl.StateResponse"/>
<a name="xctrl.StateResponse"/>

### StateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| channel_state | [string](#string) |  |  |
| call_state | [string](#string) |  |  |
| answer_state | [string](#string) |  |  |
| bridged | [bool](#bool) |  |  |
| answered | [bool](#bool) |  |  |
| hold | [bool](#bool) |  |  |
| video | [bool](#bool) |  |  |
| video_ready | [bool](#bool) |  |  |
| controlled | [bool](#bool) |  |  |
| ready | [bool](#bool) |  |  |
| up | [bool](#bool) |  |  |






<a name="user-content-xctrl.StopDetectRequest"/>
<a name="xctrl.StopDetectRequest"/>

### StopDetectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.StopRequest"/>
<a name="xctrl.StopRequest"/>

### StopRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="user-content-xctrl.ThreeWayRequest"/>
<a name="xctrl.ThreeWayRequest"/>

### ThreeWayRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| target_uuid | [string](#string) |  |  |
| direction | [string](#string) |  | LISTEN ABC AC BC TOA TOB STOP |






<a name="user-content-xctrl.TransferRequest"/>
<a name="xctrl.TransferRequest"/>

### TransferRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| extension | [string](#string) |  |  |
| dialplan | [string](#string) |  |  |
| context | [string](#string) |  |  |






<a name="user-content-xctrl.VarResponse"/>
<a name="xctrl.VarResponse"/>

### VarResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| data | [map<string, string>](#map-string-string) |  |  |






<a name="user-content-xctrl.VideoResizeEvent"/>
<a name="xctrl.VideoResizeEvent"/>

### VideoResizeEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| old_width | [uint32](#uint32) |  |  |
| old_height | [uint32](#uint32) |  |  |
| new_width | [uint32](#uint32) |  |  |
| new_height | [uint32](#uint32) |  |  |





 <!-- end messages -->


<a name="user-content-xctrl.HangupRequest.HangupFlag"/>
<a name="xctrl.HangupRequest.HangupFlag"/>

### HangupRequest.HangupFlag


| Name | Number | Description |
| ---- | ------ | ----------- |
| SELF | 0 |  |
| PEER | 1 |  |
| BOTH | 2 |  |



<a name="user-content-xctrl.MediaType"/>
<a name="xctrl.MediaType"/>

### MediaType


| Name | Number | Description |
| ---- | ------ | ----------- |
| FILE | 0 |  |
| TEXT | 1 |  |
| SSML | 2 |  |



<a name="user-content-xctrl.RecordRequest.RecordAction"/>
<a name="xctrl.RecordRequest.RecordAction"/>

### RecordRequest.RecordAction


| Name | Number | Description |
| ---- | ------ | ----------- |
| RECORD | 0 | block sync recording |
| START | 1 | unblock async recording |
| STOP | 2 |  |
| MASK | 3 |  |
| UNMASK | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="user-content-xctrl.XNode"/>
<a name="xctrl.XNode"/>

### XNode


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Dial | [DialRequest](#xctrl.DialRequest) | [DialResponse](#xctrl.DialRequest) | 外呼 |
| Answer | [Request](#xctrl.Request) | [Response](#xctrl.Request) | 应答 |
| Accept | [AcceptRequest](#xctrl.AcceptRequest) | [Response](#xctrl.AcceptRequest) | 接管呼叫，示接管的呼叫将会在10s后挂断，其它所有API都隐含接管 |
| Play | [PlayRequest](#xctrl.PlayRequest) | [Response](#xctrl.PlayRequest) | 播放一个文件或TTS |
| Stop | [StopRequest](#xctrl.StopRequest) | [Response](#xctrl.StopRequest) | 停止当前正在执行的API |
| Broadcast | [BroadcastRequest](#xctrl.BroadcastRequest) | [Response](#xctrl.BroadcastRequest) | 广播 |
| Mute | [MuteRequest](#xctrl.MuteRequest) | [Response](#xctrl.MuteRequest) | 设置静音 |
| Record | [RecordRequest](#xctrl.RecordRequest) | [RecordResponse](#xctrl.RecordRequest) | 录音 |
| Hangup | [HangupRequest](#xctrl.HangupRequest) | [Response](#xctrl.HangupRequest) | 挂断当前UUID |
| Bridge | [BridgeRequest](#xctrl.BridgeRequest) | [Response](#xctrl.BridgeRequest) | 在把当前呼叫桥接（发起）另一个呼叫 |
| ChannelBridge | [ChannelBridgeRequest](#xctrl.ChannelBridgeRequest) | [Response](#xctrl.ChannelBridgeRequest) | 桥接两个呼叫 |
| UnBridge | [Request](#xctrl.Request) | [Response](#xctrl.Request) | 将桥接的呼叫分开 |
| UnBridge2 | [Request](#xctrl.Request) | [Response](#xctrl.Request) | 将桥接的呼叫分开 |
| Hold | [HoldRequest](#xctrl.HoldRequest) | [Response](#xctrl.HoldRequest) | 呼叫保持/取消保持 |
| Transfer | [TransferRequest](#xctrl.TransferRequest) | [Response](#xctrl.TransferRequest) | 转移（待定） |
| ThreeWay | [ThreeWayRequest](#xctrl.ThreeWayRequest) | [Response](#xctrl.ThreeWayRequest) | 三方通话 |
| Echo2 | [Echo2Request](#xctrl.Echo2Request) | [Response](#xctrl.Echo2Request) | 回声，说话者可以听到自己的声音 |
| Intercept | [InterceptRequest](#xctrl.InterceptRequest) | [Response](#xctrl.InterceptRequest) | 强插 |
| Consult | [ConsultRequest](#xctrl.ConsultRequest) | [Response](#xctrl.ConsultRequest) | 协商转移 |
| SetVar | [SetVarRequest](#xctrl.SetVarRequest) | [Response](#xctrl.SetVarRequest) | 设置通道变量 |
| GetVar | [GetVarRequest](#xctrl.GetVarRequest) | [VarResponse](#xctrl.GetVarRequest) | 获取通道变量 |
| GetState | [GetStateRequest](#xctrl.GetStateRequest) | [StateResponse](#xctrl.GetStateRequest) | 获取通道状态 |
| GetChannelData | [GetChannelDataRequest](#xctrl.GetChannelDataRequest) | [ChannelDataResponse](#xctrl.GetChannelDataRequest) | 获取通道数据 |
| ReadDTMF | [DTMFRequest](#xctrl.DTMFRequest) | [DTMFResponse](#xctrl.DTMFRequest) | 读取DTMF按键 |
| ReadDigits | [DigitsRequest](#xctrl.DigitsRequest) | [DigitsResponse](#xctrl.DigitsRequest) | 读取DTMF按键 |
| DetectSpeech | [DetectRequest](#xctrl.DetectRequest) | [DetectResponse](#xctrl.DetectRequest) | 语音识别 |
| StopDetectSpeech | [StopDetectRequest](#xctrl.StopDetectRequest) | [Response](#xctrl.StopDetectRequest) | 停止语音识别 |
| RingBackDetection | [RingBackDetectionRequest](#xctrl.RingBackDetectionRequest) | [Response](#xctrl.RingBackDetectionRequest) | 回铃音检测 |
| DetectFace | [DetectFaceRequest](#xctrl.DetectFaceRequest) | [Response](#xctrl.DetectFaceRequest) | 人脸识别 |
| SendDTMF | [SendDTMFRequest](#xctrl.SendDTMFRequest) | [Response](#xctrl.SendDTMFRequest) | 发送DTMF |
| SendINFO | [SendINFORequest](#xctrl.SendINFORequest) | [Response](#xctrl.SendINFORequest) | 发送SIP INFO |
| NativeApp | [NativeRequest](#xctrl.NativeRequest) | [NativeResponse](#xctrl.NativeRequest) | 执行原生APP |
| NativeAPI | [NativeRequest](#xctrl.NativeRequest) | [NativeResponse](#xctrl.NativeRequest) | 执行原生API |
| NativeJSAPI | [NativeJSRequest](#xctrl.NativeJSRequest) | [NativeJSResponse](#xctrl.NativeJSRequest) | 执行原生JSAPI |
| JStatus | [JStatusRequest](#xctrl.JStatusRequest) | [JStatusResponse](#xctrl.JStatusRequest) | 状态 |
| ConferenceInfo | [ConferenceInfoRequest](#xctrl.ConferenceInfoRequest) | [ConferenceInfoResponse](#xctrl.ConferenceInfoRequest) | 获取会议信息 |
| FIFO | [FIFORequest](#xctrl.FIFORequest) | [FIFOResponse](#xctrl.FIFORequest) | 呼叫中心FIFO队列（先入先出） |
| Callcenter | [CallcenterRequest](#xctrl.CallcenterRequest) | [CallcenterResponse](#xctrl.CallcenterRequest) | 呼叫中心Callcenter |
| Conference | [ConferenceRequest](#xctrl.ConferenceRequest) | [ConferenceResponse](#xctrl.ConferenceRequest) | 会议Conference |
| AI | [AIRequest](#xctrl.AIRequest) | [AIResponse](#xctrl.AIRequest) | 会议AI |
| HttAPI | [HttAPIRequest](#xctrl.HttAPIRequest) | [HttAPIResponse](#xctrl.HttAPIRequest) | HttAPI |
| Lua | [LuaRequest](#xctrl.LuaRequest) | [LuaResponse](#xctrl.LuaRequest) | Lua |
| Register | [Request](#xctrl.Request) | [Response](#xctrl.Request) | Node Register |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="user-content-double" /><a name="double" /> double |  | double | double | float |
| <a name="user-content-float" /><a name="float" /> float |  | float | float | float |
| <a name="user-content-int32" /><a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="user-content-int64" /><a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="user-content-uint32" /><a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="user-content-uint64" /><a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="user-content-sint32" /><a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="user-content-sint64" /><a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="user-content-fixed32" /><a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="user-content-fixed64" /><a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="user-content-sfixed32" /><a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="user-content-sfixed64" /><a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="user-content-bool" /><a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="user-content-string" /><a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="user-content-bytes" /><a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |


<a name="user-content-map-string-string" />

map&lt;[string](#string), [string](#string)&gt;
