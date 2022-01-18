package core

import (
	"encoding/json"

	"git.xswitch.cn/xswitch/xctrl/stack/server"
)

// Response .
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ParseRequest .
func ParseRequest(req server.Request) []byte {
	data, _ := json.Marshal(req.Body())
	return data
}
