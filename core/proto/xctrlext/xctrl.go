package xctrlext

type NativeJsData struct {
	Command string      `json:"command,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
type NativeJSRequest struct {
	CtrlUuid string       `json:"ctrl_uuid,omitempty"`
	Data     NativeJsData `json:"data,omitempty"`
}

type NativeJSResponse struct {
	Code     int32  `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	NodeUuid string `json:"node_uuid,omitempty"`
	// optional
	Seq  string      `json:"seq,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
