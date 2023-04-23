# XSwitch XCC Proto Buffer 协议参考文档

<a name="top"></a>
<a name="user-content-top"></a>

这是[XCC API文档](https://docs.xswitch.cn/xcc-api/)的协议参考，使用[Google Protocol Buffers](https://protobuf.dev/)描述。

本文档只是对具体协议数据格式及类型的参考说明，详细的字段说明和用法请参考[XCC API列表](https://docs.xswitch.cn/xcc-api/api/)，原始的`.proto`文件可以在[proto](../)相关目录中找到。

## 目录

- [core/proto/xctrl/xctrl.proto](#core_proto_xctrl_xctrl-proto)
  - [AIRequest](#xctrl-AIRequest)
  - [AIRequest.DataEntry](#xctrl-AIRequest-DataEntry)
  - [AIResponse](#xctrl-AIResponse)
  - [AcceptRequest](#xctrl-AcceptRequest)
  - [Action](#xctrl-Action)
  - [Application](#xctrl-Application)
  - [BridgeRequest](#xctrl-BridgeRequest)
  - [BroadcastRequest](#xctrl-BroadcastRequest)
  - [CallParam](#xctrl-CallParam)
  - [CallParam.ParamsEntry](#xctrl-CallParam-ParamsEntry)
  - [CallcenterRequest](#xctrl-CallcenterRequest)
  - [CallcenterResponse](#xctrl-CallcenterResponse)
  - [ChannelBridgeRequest](#xctrl-ChannelBridgeRequest)
  - [ChannelData](#xctrl-ChannelData)
  - [ChannelDataResponse](#xctrl-ChannelDataResponse)
  - [ChannelEvent](#xctrl-ChannelEvent)
  - [ChannelEvent.ParamsEntry](#xctrl-ChannelEvent-ParamsEntry)
  - [ConferenceInfoRequest](#xctrl-ConferenceInfoRequest)
  - [ConferenceInfoRequestData](#xctrl-ConferenceInfoRequestData)
  - [ConferenceInfoRequestDataData](#xctrl-ConferenceInfoRequestDataData)
  - [ConferenceInfoRequestDataData.MemberFiltersEntry](#xctrl-ConferenceInfoRequestDataData-MemberFiltersEntry)
  - [ConferenceInfoResponse](#xctrl-ConferenceInfoResponse)
  - [ConferenceInfoResponseConference](#xctrl-ConferenceInfoResponseConference)
  - [ConferenceInfoResponseData](#xctrl-ConferenceInfoResponseData)
  - [ConferenceInfoResponseFlags](#xctrl-ConferenceInfoResponseFlags)
  - [ConferenceInfoResponseMembers](#xctrl-ConferenceInfoResponseMembers)
  - [ConferenceInfoResponseVariables](#xctrl-ConferenceInfoResponseVariables)
  - [ConferenceRequest](#xctrl-ConferenceRequest)
  - [ConferenceResponse](#xctrl-ConferenceResponse)
  - [ConsultRequest](#xctrl-ConsultRequest)
  - [Ctrl](#xctrl-Ctrl)
  - [DTMFRequest](#xctrl-DTMFRequest)
  - [DTMFResponse](#xctrl-DTMFResponse)
  - [Destination](#xctrl-Destination)
  - [Destination.GlobalParamsEntry](#xctrl-Destination-GlobalParamsEntry)
  - [DetectFaceRequest](#xctrl-DetectFaceRequest)
  - [DetectRequest](#xctrl-DetectRequest)
  - [DetectResponse](#xctrl-DetectResponse)
  - [DetectedData](#xctrl-DetectedData)
  - [DetectedFaceEvent](#xctrl-DetectedFaceEvent)
  - [DialRequest](#xctrl-DialRequest)
  - [DialResponse](#xctrl-DialResponse)
  - [DigitsRequest](#xctrl-DigitsRequest)
  - [DigitsResponse](#xctrl-DigitsResponse)
  - [Echo2Request](#xctrl-Echo2Request)
  - [EngineData](#xctrl-EngineData)
  - [FIFORequest](#xctrl-FIFORequest)
  - [FIFOResponse](#xctrl-FIFOResponse)
  - [GetChannelDataRequest](#xctrl-GetChannelDataRequest)
  - [GetStateRequest](#xctrl-GetStateRequest)
  - [GetVarRequest](#xctrl-GetVarRequest)
  - [HangupRequest](#xctrl-HangupRequest)
  - [Header](#xctrl-Header)
  - [HoldRequest](#xctrl-HoldRequest)
  - [HttAPIRequest](#xctrl-HttAPIRequest)
  - [HttAPIRequest.DataEntry](#xctrl-HttAPIRequest-DataEntry)
  - [HttAPIResponse](#xctrl-HttAPIResponse)
  - [InterceptRequest](#xctrl-InterceptRequest)
  - [JStatusIdleCPU](#xctrl-JStatusIdleCPU)
  - [JStatusRequest](#xctrl-JStatusRequest)
  - [JStatusRequest.JStatusData](#xctrl-JStatusRequest-JStatusData)
  - [JStatusResponse](#xctrl-JStatusResponse)
  - [JStatusResponseData](#xctrl-JStatusResponseData)
  - [JStatusSessions](#xctrl-JStatusSessions)
  - [JStatusSessionsCount](#xctrl-JStatusSessionsCount)
  - [JStatusSessionsRate](#xctrl-JStatusSessionsRate)
  - [JStatusStackSize](#xctrl-JStatusStackSize)
  - [JStatusUptime](#xctrl-JStatusUptime)
  - [Media](#xctrl-Media)
  - [MuteRequest](#xctrl-MuteRequest)
  - [NativeJSRequest](#xctrl-NativeJSRequest)
  - [NativeJSRequestData](#xctrl-NativeJSRequestData)
  - [NativeJSResponse](#xctrl-NativeJSResponse)
  - [NativeRequest](#xctrl-NativeRequest)
  - [NativeResponse](#xctrl-NativeResponse)
  - [Node](#xctrl-Node)
  - [Payload](#xctrl-Payload)
  - [PlayRequest](#xctrl-PlayRequest)
  - [RecordEvent](#xctrl-RecordEvent)
  - [RecordRequest](#xctrl-RecordRequest)
  - [RecordResponse](#xctrl-RecordResponse)
  - [Request](#xctrl-Request)
  - [Response](#xctrl-Response)
  - [RingBackDetectionRequest](#xctrl-RingBackDetectionRequest)
  - [SendDTMFRequest](#xctrl-SendDTMFRequest)
  - [SendINFORequest](#xctrl-SendINFORequest)
  - [SetVarRequest](#xctrl-SetVarRequest)
  - [SetVarRequest.DataEntry](#xctrl-SetVarRequest-DataEntry)
  - [SpeechRequest](#xctrl-SpeechRequest)
  - [SpeechRequest.ParamsEntry](#xctrl-SpeechRequest-ParamsEntry)
  - [StashResult](#xctrl-StashResult)
  - [StateResponse](#xctrl-StateResponse)
  - [StopDetectRequest](#xctrl-StopDetectRequest)
  - [StopRequest](#xctrl-StopRequest)
  - [ThreeWayRequest](#xctrl-ThreeWayRequest)
  - [TransferRequest](#xctrl-TransferRequest)
  - [VarResponse](#xctrl-VarResponse)
  - [VarResponse.DataEntry](#xctrl-VarResponse-DataEntry)
  - [VideoResizeEvent](#xctrl-VideoResizeEvent)

  - [HangupRequest.HangupFlag](#xctrl-HangupRequest-HangupFlag)
  - [MediaType](#xctrl-MediaType)
  - [RecordRequest.RecordAction](#xctrl-RecordRequest-RecordAction)

  - [XNode](#xctrl-XNode)

- [Scalar Value Types](#scalar-value-types)



<a name="core_proto_xctrl_xctrl-proto"></a>
<a name="user-content-core_proto_xctrl_xctrl-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## core/proto/xctrl/xctrl.proto



<a name="xctrl-AIRequest"></a>
<a name="user-content-xctrl-AIRequest"></a>

### AIRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| url | [string](#string) |  |  |
| data | [map<string,string>](#map<string,string>) |  |  |












<a name="xctrl-AIResponse"></a>
<a name="user-content-xctrl-AIResponse"></a>

### AIResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-AcceptRequest"></a>
<a name="user-content-xctrl-AcceptRequest"></a>

### AcceptRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| takeover | [bool](#bool) |  |  |






<a name="xctrl-Action"></a>
<a name="user-content-xctrl-Action"></a>

### Action



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| owner_uid | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| param | [CallParam](#xctrl-CallParam) |  |  |






<a name="xctrl-Application"></a>
<a name="user-content-xctrl-Application"></a>

### Application



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| app | [string](#string) |  |  |
| data | [string](#string) |  |  |






<a name="xctrl-BridgeRequest"></a>
<a name="user-content-xctrl-BridgeRequest"></a>

### BridgeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| ringback | [string](#string) |  |  |
| flow_control | [string](#string) |  | NONE | CALLER | CALLEE | ANY |
| continue_on_fail | [string](#string) |  | true | false | comma separated freeswitch causes e.g. USER_BUSY,NO_ANSWER |
| destination | [Destination](#xctrl-Destination) |  |  |






<a name="xctrl-BroadcastRequest"></a>
<a name="user-content-xctrl-BroadcastRequest"></a>

### BroadcastRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| media | [Media](#xctrl-Media) |  |  |
| option | [string](#string) |  | BOTH, ALEG, BLEG, AHOLDB, BHOLDA |






<a name="xctrl-CallParam"></a>
<a name="user-content-xctrl-CallParam"></a>

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
| params | [map<string,string>](#map<string,string>) | repeated |  |












<a name="xctrl-CallcenterRequest"></a>
<a name="user-content-xctrl-CallcenterRequest"></a>

### CallcenterRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |






<a name="xctrl-CallcenterResponse"></a>
<a name="user-content-xctrl-CallcenterResponse"></a>

### CallcenterResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-ChannelBridgeRequest"></a>
<a name="user-content-xctrl-ChannelBridgeRequest"></a>

### ChannelBridgeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| peer_uuid | [string](#string) |  |  |
| bridge_delay | [int32](#int32) |  |  |
| flow_control | [string](#string) |  | NONE | CALLER | CALLEE | ANY |






<a name="xctrl-ChannelData"></a>
<a name="user-content-xctrl-ChannelData"></a>

### ChannelData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| variable_cc_queue | [string](#string) |  |  |
| variable_cc_queue_name | [string](#string) |  |  |
| variable_cc_agent_session_uuid | [string](#string) |  |  |
| variable_cc_member_uuid | [string](#string) |  |  |
| variable_xcc_origin_dest_number | [string](#string) |  |  |






<a name="xctrl-ChannelDataResponse"></a>
<a name="user-content-xctrl-ChannelDataResponse"></a>

### ChannelDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| format | [string](#string) |  |  |
| data | [ChannelData](#xctrl-ChannelData) |  | when format == JSON |






<a name="xctrl-ChannelEvent"></a>
<a name="user-content-xctrl-ChannelEvent"></a>

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
| params | [map<string,string>](#map<string,string>) |  |  |
| billsec | [string](#string) |  |  |
| duration | [string](#string) |  |  |
| cause | [string](#string) |  |  |
| answered | [bool](#bool) |  |  |
| node_ip | [string](#string) |  |  |
| domain | [string](#string) |  |  |












<a name="xctrl-ConferenceInfoRequest"></a>
<a name="user-content-xctrl-ConferenceInfoRequest"></a>

### ConferenceInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| data | [ConferenceInfoRequestData](#xctrl-ConferenceInfoRequestData) |  |  |






<a name="xctrl-ConferenceInfoRequestData"></a>
<a name="user-content-xctrl-ConferenceInfoRequestData"></a>

### ConferenceInfoRequestData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| command | [string](#string) |  |  |
| data | [ConferenceInfoRequestDataData](#xctrl-ConferenceInfoRequestDataData) |  |  |






<a name="xctrl-ConferenceInfoRequestDataData"></a>
<a name="user-content-xctrl-ConferenceInfoRequestDataData"></a>

### ConferenceInfoRequestDataData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| conferenceName | [string](#string) |  |  |
| showMembers | [bool](#bool) |  |  |
| memberFilters | [map<string,string>](#map<string,string>) |  |  |













<a name="xctrl-ConferenceInfoResponse"></a>
<a name="user-content-xctrl-ConferenceInfoResponse"></a>

### ConferenceInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| seq | [string](#string) |  | optional |
| data | [ConferenceInfoResponseData](#xctrl-ConferenceInfoResponseData) |  |  |






<a name="xctrl-ConferenceInfoResponseConference"></a>
<a name="user-content-xctrl-ConferenceInfoResponseConference"></a>

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
| variables | [ConferenceInfoResponseVariables](#xctrl-ConferenceInfoResponseVariables) |  |  |
| members | [ConferenceInfoResponseMembers](#xctrl-ConferenceInfoResponseMembers) | repeated |  |






<a name="xctrl-ConferenceInfoResponseData"></a>
<a name="user-content-xctrl-ConferenceInfoResponseData"></a>

### ConferenceInfoResponseData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| conference | [ConferenceInfoResponseConference](#xctrl-ConferenceInfoResponseConference) |  |  |






<a name="xctrl-ConferenceInfoResponseFlags"></a>
<a name="user-content-xctrl-ConferenceInfoResponseFlags"></a>

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






<a name="xctrl-ConferenceInfoResponseMembers"></a>
<a name="user-content-xctrl-ConferenceInfoResponseMembers"></a>

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
| flags | [ConferenceInfoResponseFlags](#xctrl-ConferenceInfoResponseFlags) |  |  |






<a name="xctrl-ConferenceInfoResponseVariables"></a>
<a name="user-content-xctrl-ConferenceInfoResponseVariables"></a>

### ConferenceInfoResponseVariables







<a name="xctrl-ConferenceRequest"></a>
<a name="user-content-xctrl-ConferenceRequest"></a>

### ConferenceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| profile | [string](#string) |  |  |
| flags | [string](#string) | repeated |  |






<a name="xctrl-ConferenceResponse"></a>
<a name="user-content-xctrl-ConferenceResponse"></a>

### ConferenceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-ConsultRequest"></a>
<a name="user-content-xctrl-ConsultRequest"></a>

### ConsultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| destination | [Destination](#xctrl-Destination) |  |  |






<a name="xctrl-Ctrl"></a>
<a name="user-content-xctrl-Ctrl"></a>

### Ctrl



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| ip | [string](#string) |  |  |
| version | [string](#string) |  |  |
| rack | [uint32](#uint32) |  |  |
| address | [string](#string) |  |  |






<a name="xctrl-DTMFRequest"></a>
<a name="user-content-xctrl-DTMFRequest"></a>

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
| media | [Media](#xctrl-Media) |  | play or tts |
| max_tries | [uint32](#uint32) |  | not implemented yet |
| regex | [string](#string) |  |  |
| media_invalid | [Media](#xctrl-Media) |  | invalid meida |
| play_last_invalid_prompt | [bool](#bool) |  | default false |






<a name="xctrl-DTMFResponse"></a>
<a name="user-content-xctrl-DTMFResponse"></a>

### DTMFResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| dtmf | [string](#string) |  |  |
| terminator | [string](#string) |  | optional |






<a name="xctrl-Destination"></a>
<a name="user-content-xctrl-Destination"></a>

### Destination



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ringall | [bool](#bool) |  |  |
| global_params | [map<string,string>](#map<string,string>) | repeated |  |
| call_params | [CallParam](#xctrl-CallParam) | repeated |  |
| channel_params | [string](#string) | repeated |  |












<a name="xctrl-DetectFaceRequest"></a>
<a name="user-content-xctrl-DetectFaceRequest"></a>

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






<a name="xctrl-DetectRequest"></a>
<a name="user-content-xctrl-DetectRequest"></a>

### DetectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| media | [Media](#xctrl-Media) |  | oneof play or tts |
| dtmf | [DTMFRequest](#xctrl-DTMFRequest) |  | detect dtmf too, optional |
| speech | [SpeechRequest](#xctrl-SpeechRequest) |  | speech params, mandatory |






<a name="xctrl-DetectResponse"></a>
<a name="user-content-xctrl-DetectResponse"></a>

### DetectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| data | [DetectedData](#xctrl-DetectedData) |  |  |






<a name="xctrl-DetectedData"></a>
<a name="user-content-xctrl-DetectedData"></a>

### DetectedData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dtmf | [string](#string) |  |  |
| terminator | [string](#string) |  |  |
| text | [string](#string) |  |  |
| confidence | [double](#double) |  |  |
| is_final | [bool](#bool) |  | final or partial result |
| error | [string](#string) |  | when error |
| type | [string](#string) |  | DTMF Speech.Begin Speech.Partial Speech.End ERROR */ |
| engine | [string](#string) |  | the ASR engine |
| engine_data | [EngineData](#xctrl-EngineData) |  | string or JSON Struct, detailed object returned from ASR engine |
| offset | [uint32](#uint32) |  |  |






<a name="xctrl-DetectedFaceEvent"></a>
<a name="user-content-xctrl-DetectedFaceEvent"></a>

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






<a name="xctrl-DialRequest"></a>
<a name="user-content-xctrl-DialRequest"></a>

### DialRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| destination | [Destination](#xctrl-Destination) |  |  |
| apps | [Application](#xctrl-Application) | repeated |  |






<a name="xctrl-DialResponse"></a>
<a name="user-content-xctrl-DialResponse"></a>

### DialResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| cause | [string](#string) |  |  |






<a name="xctrl-DigitsRequest"></a>
<a name="user-content-xctrl-DigitsRequest"></a>

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
| media | [Media](#xctrl-Media) |  | play or tts |
| max_tries | [uint32](#uint32) |  | not implemented yet |
| regex | [string](#string) |  |  |
| media_invalid | [Media](#xctrl-Media) |  | invalid meida |
| play_last_invalid_prompt | [bool](#bool) |  | default false |






<a name="xctrl-DigitsResponse"></a>
<a name="user-content-xctrl-DigitsResponse"></a>

### DigitsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| dtmf | [string](#string) |  |  |
| terminator | [string](#string) |  | optional |






<a name="xctrl-Echo2Request"></a>
<a name="user-content-xctrl-Echo2Request"></a>

### Echo2Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| action | [string](#string) |  | START | STOP |
| direction | [string](#string) |  | SELF | OTHER |






<a name="xctrl-EngineData"></a>
<a name="user-content-xctrl-EngineData"></a>

### EngineData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [Header](#xctrl-Header) |  |  |
| payload | [Payload](#xctrl-Payload) |  |  |






<a name="xctrl-FIFORequest"></a>
<a name="user-content-xctrl-FIFORequest"></a>

### FIFORequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| inout | [string](#string) |  |  |
| wait_music | [string](#string) |  |  |
| exit_announce | [string](#string) |  |  |
| priority | [int32](#int32) |  |  |






<a name="xctrl-FIFOResponse"></a>
<a name="user-content-xctrl-FIFOResponse"></a>

### FIFOResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-GetChannelDataRequest"></a>
<a name="user-content-xctrl-GetChannelDataRequest"></a>

### GetChannelDataRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| format | [string](#string) |  | optional JSON(default) JSONSTR XML TXT LIST |






<a name="xctrl-GetStateRequest"></a>
<a name="user-content-xctrl-GetStateRequest"></a>

### GetStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-GetVarRequest"></a>
<a name="user-content-xctrl-GetVarRequest"></a>

### GetVarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| data | [string](#string) | repeated |  |






<a name="xctrl-HangupRequest"></a>
<a name="user-content-xctrl-HangupRequest"></a>

### HangupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| cause | [string](#string) |  | NORMAL_CLEARING USER_BUSY CALL_REJECTED ... |
| flag | [HangupRequest.HangupFlag](#xctrl-HangupRequest-HangupFlag) |  |  |






<a name="xctrl-Header"></a>
<a name="user-content-xctrl-Header"></a>

### Header



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| namespace | [string](#string) |  |  |
| name | [string](#string) |  |  |
| status | [double](#double) |  |  |
| message_id | [string](#string) |  |  |
| task_id | [string](#string) |  |  |
| status_text | [string](#string) |  |  |






<a name="xctrl-HoldRequest"></a>
<a name="user-content-xctrl-HoldRequest"></a>

### HoldRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| action | [string](#string) |  | ON OFF TOGGLE |
| display | [string](#string) |  | OPTIONAL only supported by some phones |






<a name="xctrl-HttAPIRequest"></a>
<a name="user-content-xctrl-HttAPIRequest"></a>

### HttAPIRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| url | [string](#string) |  |  |
| data | [map<string,string>](#map<string,string>) |  |  |












<a name="xctrl-HttAPIResponse"></a>
<a name="user-content-xctrl-HttAPIResponse"></a>

### HttAPIResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-InterceptRequest"></a>
<a name="user-content-xctrl-InterceptRequest"></a>

### InterceptRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| target_uuid | [string](#string) |  |  |






<a name="xctrl-JStatusIdleCPU"></a>
<a name="user-content-xctrl-JStatusIdleCPU"></a>

### JStatusIdleCPU



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| used | [float](#float) |  |  |
| allowed | [float](#float) |  |  |






<a name="xctrl-JStatusRequest"></a>
<a name="user-content-xctrl-JStatusRequest"></a>

### JStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| data | [JStatusRequest.JStatusData](#xctrl-JStatusRequest-JStatusData) |  |  |






<a name="xctrl-JStatusRequest-JStatusData"></a>
<a name="user-content-xctrl-JStatusRequest-JStatusData"></a>

### JStatusRequest.JStatusData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| command | [string](#string) |  |  |
| data | [string](#string) |  |  |






<a name="xctrl-JStatusResponse"></a>
<a name="user-content-xctrl-JStatusResponse"></a>

### JStatusResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| seq | [string](#string) |  |  |
| data | [JStatusResponseData](#xctrl-JStatusResponseData) |  |  |






<a name="xctrl-JStatusResponseData"></a>
<a name="user-content-xctrl-JStatusResponseData"></a>

### JStatusResponseData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| systemStatus | [string](#string) |  |  |
| version | [string](#string) |  |  |
| uptime | [JStatusUptime](#xctrl-JStatusUptime) |  |  |
| sessions | [JStatusSessions](#xctrl-JStatusSessions) |  |  |
| idleCPU | [JStatusIdleCPU](#xctrl-JStatusIdleCPU) |  |  |
| stackSizeKB | [JStatusStackSize](#xctrl-JStatusStackSize) |  |  |






<a name="xctrl-JStatusSessions"></a>
<a name="user-content-xctrl-JStatusSessions"></a>

### JStatusSessions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| count | [JStatusSessionsCount](#xctrl-JStatusSessionsCount) |  |  |
| rate | [JStatusSessionsRate](#xctrl-JStatusSessionsRate) |  |  |






<a name="xctrl-JStatusSessionsCount"></a>
<a name="user-content-xctrl-JStatusSessionsCount"></a>

### JStatusSessionsCount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int32](#int32) |  |  |
| active | [int32](#int32) |  |  |
| peak | [int32](#int32) |  |  |
| peak5Min | [int32](#int32) |  |  |
| limit | [int32](#int32) |  |  |






<a name="xctrl-JStatusSessionsRate"></a>
<a name="user-content-xctrl-JStatusSessionsRate"></a>

### JStatusSessionsRate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| current | [int32](#int32) |  |  |
| max | [int32](#int32) |  |  |
| peak | [int32](#int32) |  |  |
| peak5Min | [int32](#int32) |  |  |






<a name="xctrl-JStatusStackSize"></a>
<a name="user-content-xctrl-JStatusStackSize"></a>

### JStatusStackSize



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| current | [int32](#int32) |  |  |
| max | [int32](#int32) |  |  |






<a name="xctrl-JStatusUptime"></a>
<a name="user-content-xctrl-JStatusUptime"></a>

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






<a name="xctrl-Media"></a>
<a name="user-content-xctrl-Media"></a>

### Media



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  | FILE TEXT SSML |
| data | [string](#string) |  |  |
| engine | [string](#string) |  |  |
| voice | [string](#string) |  |  |
| loop | [uint32](#uint32) |  |  |
| offset | [uint32](#uint32) |  |  |






<a name="xctrl-MuteRequest"></a>
<a name="user-content-xctrl-MuteRequest"></a>

### MuteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| direction | [string](#string) |  | WRITE, READ, BOTH |
| level | [int32](#int32) |  |  |
| flag | [string](#string) |  | FIRST, LAST |






<a name="xctrl-NativeJSRequest"></a>
<a name="user-content-xctrl-NativeJSRequest"></a>

### NativeJSRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| data | [NativeJSRequestData](#xctrl-NativeJSRequestData) |  |  |






<a name="xctrl-NativeJSRequestData"></a>
<a name="user-content-xctrl-NativeJSRequestData"></a>

### NativeJSRequestData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| command | [string](#string) |  |  |
| data | [string](#string) |  | a string or a native JSON struct to google.protobuf.Any or .Struct |






<a name="xctrl-NativeJSResponse"></a>
<a name="user-content-xctrl-NativeJSResponse"></a>

### NativeJSResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| seq | [string](#string) |  |  |
| data | [string](#string) |  |  |






<a name="xctrl-NativeRequest"></a>
<a name="user-content-xctrl-NativeRequest"></a>

### NativeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| cmd | [string](#string) |  |  |
| args | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-NativeResponse"></a>
<a name="user-content-xctrl-NativeResponse"></a>

### NativeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| data | [string](#string) |  |  |
| result | [string](#string) |  |  |






<a name="xctrl-Node"></a>
<a name="user-content-xctrl-Node"></a>

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






<a name="xctrl-Payload"></a>
<a name="user-content-xctrl-Payload"></a>

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
| stash_result | [StashResult](#xctrl-StashResult) |  |  |
| audio_extra_info | [string](#string) |  |  |
| sentence_id | [string](#string) |  |  |
| gender_score | [double](#double) |  |  |






<a name="xctrl-PlayRequest"></a>
<a name="user-content-xctrl-PlayRequest"></a>

### PlayRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| media | [Media](#xctrl-Media) |  |  |






<a name="xctrl-RecordEvent"></a>
<a name="user-content-xctrl-RecordEvent"></a>

### RecordEvent
enum Method {
Invalid = 0;
// 节点注册
NodeRegister = 1;
// 节点离线
NodeUnregister = 2;
// 节点数据更新
NodeUpdate = 3;

// 通道
Channel = 4;
// FreeSWITCH原生消息
Native = 5;

// old VCC event
Vcc = 6;
// 如果API请求有后续事件
Result = 7;
//　异步Dial结果
DialResult = 8;
// NativeAPI结果
NativeResult = 9;

// 获取配置信息
FetchXML = 10;
// 获取Dialplan
Dialplan = 11;

// 话机消息
NativeEvent = 12;

// 话单
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






<a name="xctrl-RecordRequest"></a>
<a name="user-content-xctrl-RecordRequest"></a>

### RecordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| path | [string](#string) |  |  |
| action | [string](#string) |  | all params are optional |
| limit | [uint32](#uint32) |  |  |
| beep | [string](#string) |  | play a beep before record "default" or TGML https://freeswitch.org/confluence/display/FREESWITCH/TGML |
| terminators | [string](#string) |  |  |
| silence_seconds | [uint32](#uint32) |  |  |
| thresh | [uint32](#uint32) |  | VAD threshold, 0 is disabled, 1~10000 |
| rate | [uint32](#uint32) |  | valid rates are 8000, 16000, 22050, 24000, 32000, 44100, 48000 |






<a name="xctrl-RecordResponse"></a>
<a name="user-content-xctrl-RecordResponse"></a>

### RecordResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| terminator | [string](#string) |  | if terminated by DTMF |
| path | [string](#string) |  | mirror back of the path |






<a name="xctrl-Request"></a>
<a name="user-content-xctrl-Request"></a>

### Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-Response"></a>
<a name="user-content-xctrl-Response"></a>

### Response



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |






<a name="xctrl-RingBackDetectionRequest"></a>
<a name="user-content-xctrl-RingBackDetectionRequest"></a>

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






<a name="xctrl-SendDTMFRequest"></a>
<a name="user-content-xctrl-SendDTMFRequest"></a>

### SendDTMFRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| dtmf | [string](#string) |  |  |






<a name="xctrl-SendINFORequest"></a>
<a name="user-content-xctrl-SendINFORequest"></a>

### SendINFORequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| content_type | [string](#string) |  |  |
| data | [string](#string) |  |  |






<a name="xctrl-SetVarRequest"></a>
<a name="user-content-xctrl-SetVarRequest"></a>

### SetVarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| data | [map<string,string>](#map<string,string>) |  |  |












<a name="xctrl-SpeechRequest"></a>
<a name="user-content-xctrl-SpeechRequest"></a>

### SpeechRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| engine | [string](#string) |  |  |
| no_input_timeout | [uint32](#uint32) |  |  |
| speech_timeout | [uint32](#uint32) |  |  |
| partial_events | [bool](#bool) |  |  |
| disable_detected_data_event | [bool](#bool) |  |  |
| params | [map<string,string>](#map<string,string>) |  |  |
| grammar | [string](#string) |  |  |
| max_speech_timeout | [uint32](#uint32) |  |  |












<a name="xctrl-StashResult"></a>
<a name="user-content-xctrl-StashResult"></a>

### StashResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sentenceId | [double](#double) |  |  |
| beginTime | [double](#double) |  |  |
| text | [string](#string) |  |  |
| currentTime | [double](#double) |  |  |
| words | [string](#string) | repeated |  |






<a name="xctrl-StateResponse"></a>
<a name="user-content-xctrl-StateResponse"></a>

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






<a name="xctrl-StopDetectRequest"></a>
<a name="user-content-xctrl-StopDetectRequest"></a>

### StopDetectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-StopRequest"></a>
<a name="user-content-xctrl-StopRequest"></a>

### StopRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |






<a name="xctrl-ThreeWayRequest"></a>
<a name="user-content-xctrl-ThreeWayRequest"></a>

### ThreeWayRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| target_uuid | [string](#string) |  |  |
| direction | [string](#string) |  | LISTEN ABC AC BC TOA TOB STOP |






<a name="xctrl-TransferRequest"></a>
<a name="user-content-xctrl-TransferRequest"></a>

### TransferRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctrl_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| extension | [string](#string) |  |  |
| dialplan | [string](#string) |  |  |
| context | [string](#string) |  |  |






<a name="xctrl-VarResponse"></a>
<a name="user-content-xctrl-VarResponse"></a>

### VarResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |
| node_uuid | [string](#string) |  |  |
| uuid | [string](#string) |  | optional |
| data | [map<string,string>](#map<string,string>) |  |  |











<a name="xctrl-VideoResizeEvent"></a>
<a name="user-content-xctrl-VideoResizeEvent"></a>

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


<a name="xctrl-HangupRequest-HangupFlag"></a>
<a name="user-content-xctrl-HangupRequest-HangupFlag"></a>

### HangupRequest.HangupFlag


| Name | Number | Description |
| ---- | ------ | ----------- |
| SELF | 0 |  |
| PEER | 1 |  |
| BOTH | 2 |  |



<a name="xctrl-MediaType"></a>
<a name="user-content-xctrl-MediaType"></a>

### MediaType


| Name | Number | Description |
| ---- | ------ | ----------- |
| FILE | 0 |  |
| TEXT | 1 |  |
| SSML | 2 |  |



<a name="xctrl-RecordRequest-RecordAction"></a>
<a name="user-content-xctrl-RecordRequest-RecordAction"></a>

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


<a name="xctrl-XNode"></a>
<a name="user-content-xctrl-XNode"></a>

### XNode


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Dial | [DialRequest](#xctrl-DialRequest) | [DialResponse](#xctrl-DialResponse) | 外呼 |
| Answer | [Request](#xctrl-Request) | [Response](#xctrl-Response) | 应答 |
| Accept | [AcceptRequest](#xctrl-AcceptRequest) | [Response](#xctrl-Response) | 接管呼叫，示接管的呼叫将会在10s后挂断，其它所有API都隐含接管 |
| Play | [PlayRequest](#xctrl-PlayRequest) | [Response](#xctrl-Response) | 播放一个文件或TTS |
| Stop | [StopRequest](#xctrl-StopRequest) | [Response](#xctrl-Response) | 停止当前正在执行的API |
| Broadcast | [BroadcastRequest](#xctrl-BroadcastRequest) | [Response](#xctrl-Response) | 广播 |
| Mute | [MuteRequest](#xctrl-MuteRequest) | [Response](#xctrl-Response) | 设置静音 |
| Record | [RecordRequest](#xctrl-RecordRequest) | [RecordResponse](#xctrl-RecordResponse) | 录音 |
| Hangup | [HangupRequest](#xctrl-HangupRequest) | [Response](#xctrl-Response) | 挂断当前UUID |
| Bridge | [BridgeRequest](#xctrl-BridgeRequest) | [Response](#xctrl-Response) | 在把当前呼叫桥接（发起）另一个呼叫 |
| ChannelBridge | [ChannelBridgeRequest](#xctrl-ChannelBridgeRequest) | [Response](#xctrl-Response) | 桥接两个呼叫 |
| UnBridge | [Request](#xctrl-Request) | [Response](#xctrl-Response) | 将桥接的呼叫分开 |
| UnBridge2 | [Request](#xctrl-Request) | [Response](#xctrl-Response) | 将桥接的呼叫分开 |
| Hold | [HoldRequest](#xctrl-HoldRequest) | [Response](#xctrl-Response) | 呼叫保持/取消保持 |
| Transfer | [TransferRequest](#xctrl-TransferRequest) | [Response](#xctrl-Response) | 转移（待定） |
| ThreeWay | [ThreeWayRequest](#xctrl-ThreeWayRequest) | [Response](#xctrl-Response) | 三方通话 |
| Echo2 | [Echo2Request](#xctrl-Echo2Request) | [Response](#xctrl-Response) | 回声，说话者可以听到自己的声音 |
| Intercept | [InterceptRequest](#xctrl-InterceptRequest) | [Response](#xctrl-Response) | 强插 |
| Consult | [ConsultRequest](#xctrl-ConsultRequest) | [Response](#xctrl-Response) | 协商转移 |
| SetVar | [SetVarRequest](#xctrl-SetVarRequest) | [Response](#xctrl-Response) | 设置通道变量 |
| GetVar | [GetVarRequest](#xctrl-GetVarRequest) | [VarResponse](#xctrl-VarResponse) | 获取通道变量 |
| GetState | [GetStateRequest](#xctrl-GetStateRequest) | [StateResponse](#xctrl-StateResponse) | 获取通道状态 |
| GetChannelData | [GetChannelDataRequest](#xctrl-GetChannelDataRequest) | [ChannelDataResponse](#xctrl-ChannelDataResponse) | 获取通道数据 |
| ReadDTMF | [DTMFRequest](#xctrl-DTMFRequest) | [DTMFResponse](#xctrl-DTMFResponse) | 读取DTMF按键 |
| ReadDigits | [DigitsRequest](#xctrl-DigitsRequest) | [DigitsResponse](#xctrl-DigitsResponse) | 读取DTMF按键 |
| DetectSpeech | [DetectRequest](#xctrl-DetectRequest) | [DetectResponse](#xctrl-DetectResponse) | 语音识别 |
| StopDetectSpeech | [StopDetectRequest](#xctrl-StopDetectRequest) | [Response](#xctrl-Response) | 停止语音识别 |
| RingBackDetection | [RingBackDetectionRequest](#xctrl-RingBackDetectionRequest) | [Response](#xctrl-Response) | 回铃音检测 |
| DetectFace | [DetectFaceRequest](#xctrl-DetectFaceRequest) | [Response](#xctrl-Response) | 人脸识别 |
| SendDTMF | [SendDTMFRequest](#xctrl-SendDTMFRequest) | [Response](#xctrl-Response) | 发送DTMF |
| SendINFO | [SendINFORequest](#xctrl-SendINFORequest) | [Response](#xctrl-Response) | 发送SIP INFO |
| NativeApp | [NativeRequest](#xctrl-NativeRequest) | [NativeResponse](#xctrl-NativeResponse) | 执行原生APP |
| NativeAPI | [NativeRequest](#xctrl-NativeRequest) | [NativeResponse](#xctrl-NativeResponse) | 执行原生API |
| NativeJSAPI | [NativeJSRequest](#xctrl-NativeJSRequest) | [NativeJSResponse](#xctrl-NativeJSResponse) | 执行原生JSAPI |
| JStatus | [JStatusRequest](#xctrl-JStatusRequest) | [JStatusResponse](#xctrl-JStatusResponse) | 状态 |
| ConferenceInfo | [ConferenceInfoRequest](#xctrl-ConferenceInfoRequest) | [ConferenceInfoResponse](#xctrl-ConferenceInfoResponse) | 获取会议信息 |
| FIFO | [FIFORequest](#xctrl-FIFORequest) | [FIFOResponse](#xctrl-FIFOResponse) | 呼叫中心FIFO队列（先入先出） |
| Callcenter | [CallcenterRequest](#xctrl-CallcenterRequest) | [CallcenterResponse](#xctrl-CallcenterResponse) | 呼叫中心Callcenter |
| Conference | [ConferenceRequest](#xctrl-ConferenceRequest) | [ConferenceResponse](#xctrl-ConferenceResponse) | 会议Conference |
| AI | [AIRequest](#xctrl-AIRequest) | [AIResponse](#xctrl-AIResponse) | 会议AI |
| HttAPI | [HttAPIRequest](#xctrl-HttAPIRequest) | [HttAPIResponse](#xctrl-HttAPIResponse) | HttAPI |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /><a name="user-content-double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /><a name="user-content-float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /><a name="user-content-int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /><a name="user-content-int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /><a name="user-content-uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /><a name="user-content-uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /><a name="user-content-sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /><a name="user-content-sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /><a name="user-content-fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /><a name="user-content-fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /><a name="user-content-sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /><a name="user-content-sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /><a name="user-content-bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /><a name="user-content-string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /><a name="user-content-bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

