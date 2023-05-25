package ctrl

import (
	"context"
	"encoding/json"
	"fmt"

	"git.xswitch.cn/xswitch/xctrl/xctrl/metadata"
)

// Request RPC 请求对象
type Request struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
	ID      *json.RawMessage `json:"id"`
}

// Request RPC 请求对象
type XRequest struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id"`
}

// Response RPC 返回对象
type Response struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Result  interface{}      `json:"result,omitempty"`
	Error   interface{}      `json:"error,omitempty"`
}

// Result RPC 异步返回对象
type Result struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *json.RawMessage `json:"error,omitempty"`
}

// Message Node异步请求消息
type Message struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	ID      *json.RawMessage `json:"id"`
	Params  *json.RawMessage `json:"params"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *json.RawMessage `json:"error,omitempty"`
}

func (m *Message) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func (r *Response) getError() error {
	if r.Error == nil {
		return nil
	}
	return fmt.Errorf("%v", r.Error)
}

func (r *Request) Marshal() []byte {
	b, _ := json.Marshal(r)
	return b
}

func (r *Request) RawMessage() *json.RawMessage {
	b, _ := json.Marshal(r)
	raw := json.RawMessage(b)
	return &raw
}

// ContextWithID 创建带请求ID的context
func ContextWithID(id string) context.Context {
	ctx := metadata.NewContext(context.Background(),
		metadata.Metadata{
			"request-seq-id": id,
		})
	return ctx
}

// RawMessage .
func RawMessage(data []byte) *json.RawMessage {
	raw := json.RawMessage(data)
	return &raw
}
