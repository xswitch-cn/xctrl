package xctrl

import "encoding/json"

type XNativeJSRequestData struct {
	Command string          `json:"command,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type XNativeJSRequest struct {
	CtrlUuid string                `json:"ctrl_uuid,omitempty"`
	Data     *XNativeJSRequestData `json:"data,omitempty"`
}

type XNativeJSResponse struct {
	Code     int32  `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	NodeUuid string `json:"node_uuid,omitempty"`
	// optional
	Seq  string           `json:"seq,omitempty"`
	Data *json.RawMessage `json:"data,omitempty"`
}

var service *XNodeService

func Service() XNodeService {
	return *service
}

func SetService(s *XNodeService) {
	service = s
}
