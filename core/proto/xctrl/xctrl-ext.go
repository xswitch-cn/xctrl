package xctrl

type XNativeJSData struct {
	Command string      `json:"command,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
type XNativeJSRequest struct {
	CtrlUuid string        `json:"ctrl_uuid,omitempty"`
	Data     XNativeJSData `json:"data,omitempty"`
}

type XNativeJSResponse struct {
	Code     int32  `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	NodeUuid string `json:"node_uuid,omitempty"`
	// optional
	Seq  string      `json:"seq,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
