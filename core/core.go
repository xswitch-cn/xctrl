package core

import (
	"encoding/json"

	"git.xswitch.cn/xswitch/xctrl/xctrl/server"
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
