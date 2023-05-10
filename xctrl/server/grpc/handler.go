package grpc

import (
	"git.xswitch.cn/xswitch/xctrl/xctrl/server"
)

type rpcHandler struct {
	name    string
	handler interface{}
	opts    server.HandlerOptions
}

func (r *rpcHandler) Name() string {
	return r.name
}

func (r *rpcHandler) Handler() interface{} {
	return r.handler
}

func (r *rpcHandler) Options() server.HandlerOptions {
	return r.opts
}
