package ctrl

import (
	"encoding/json"
	"strings"
	"time"

	"git.xswitch.cn/xswitch/proto/go/proto/cman"
	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/proto/xctrl/client"
	"git.xswitch.cn/xswitch/proto/xctrl/util/log"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	natsio "github.com/nats-io/nats.go"
)

func NewCtrlInstance(trace bool, addrs string) (*Ctrl, error) {
	log.Infof("ctrl starting with addrs=%s\n", addrs)
	c, err := initCtrl(trace, strings.Split(addrs, ",")...)
	if err != nil {
		return &Ctrl{}, err
	}
	c.nodes = InitCtrlNodes()
	return c, nil
}

func (c *Ctrl) SetInstanceName(instanceName string) {
	if instanceName != "" {
		c.instanceName = instanceName
	}
}

func (c *Ctrl) SetMaxChannelLifeTime(time uint) {
	c.maxChannelLifeTime = time
}

func (c *Ctrl) GetNATSConn() *natsio.Conn {
	return c.conn.GetConn()
}

func (c *Ctrl) GetInstanceName() string {
	return c.instanceName
}

func (c *Ctrl) GetNodeList() map[string]*xctrl.Node {
	return c.nodes.GetNodeList()
}

// UUID get ctrl uuid
func (c *Ctrl) UUID() string {
	return c.uuid
}

// Service 同步调用
func (c *Ctrl) Service() xctrl.XNodeService {
	if c.service == nil {
		return nil
	}
	return c.service
}

// AsyncService 异步调用，Depracated
func (c *Ctrl) AsyncService() xctrl.XNodeService {
	if c.asyncService == nil {
		return nil
	}
	log.Warn("AsyncService is deprecated, use Service with WithAsync option instead")
	return c.asyncService
}

func (c *Ctrl) AService() xctrl.XNodeService {
	if c.asyncService == nil {
		return nil
	}
	return c.aService
}

// CManService 同步调用
func (c *Ctrl) CManService() cman.CManService {
	if c.cmanService == nil {
		return nil
	}
	return c.cmanService
}

// Publish 发送消息
func (c *Ctrl) Publish(topic string, msg []byte, opts ...nats.PublishOption) error {
	return c.conn.Publish(topic, msg, opts...)
}

// PublishJSON 发送JSON消息
func (c *Ctrl) PublishJSON(topic string, obj interface{}, opts ...nats.PublishOption) error {
	msg, _ := json.MarshalIndent(obj, "", "  ")
	return c.Publish(topic, msg, opts...)
}

func (c *Ctrl) Transfer(ctrlID string, channel *xctrl.ChannelEvent) error {
	body, err := json.Marshal(channel)
	if err != nil {
		return err
	}
	channel.State = "START"
	request := Request{
		Version: "2.0",
		Method:  "XNode.Channel",
		Params:  RawMessage(body),
	}
	return c.Publish("cn.xswitch.ctrl."+ctrlID, request.Marshal())
}

func (c *Ctrl) CtrlStartUp(req *xctrl.CtrlStartUpRequest) error {
	request := Request{
		Version: "2.0",
		Method:  "XNode.CtrlStartUp",
		Params:  ToRawMessage(req),
	}
	return c.Publish("cn.xswitch.node", request.Marshal())
}

// Call 发起 request 请求
func (c *Ctrl) Call(topic string, req *Request, timeout time.Duration) (*nats.Message, error) {
	req.Version = "2.0"
	//change part of the request's method into NativeJSAPI
	req.Method = TranslateMethod(req.Method)

	body, err := json.Marshal(req)
	if err != nil {
		log.Errorf("execute native api error: %v", err)
		return nil, err
	}
	return c.conn.Request(topic, body, timeout)
}

// XCall 发起 request 请求
func (c *Ctrl) XCall(topic string, method string, params interface{}, timeout time.Duration) (*nats.Message, error) {
	//change part of the request's method into NativeJSAPI
	method = TranslateMethod(method)
	req := XRequest{
		Version: "2.0",
		Method:  method,
		ID:      "0",
		Params:  params,
	}
	body, err := json.Marshal(req)
	if err != nil {
		log.Errorf("execute native api error: %v", err)
		return nil, err
	}
	return c.conn.Request(topic, body, timeout)
}

// Respond 响应NATS Request 请求
func (c *Ctrl) Respond(topic string, resp *Response, opts ...nats.PublishOption) error {
	resp.Version = "2.0"
	body, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Errorf("execute native api error: %v", err)
		return err
	}
	return c.conn.Publish(topic, body)
}

func (c *Ctrl) SetFromPrefix(prefix string) {
	c.fromPrefix = prefix
}

func SetFromPrefix(prefix string) {
	if globalCtrl != nil {
		globalCtrl.SetFromPrefix(prefix)
	}
}

func (c *Ctrl) SetToPrefix(prefix string) {
	c.toPrefix = prefix
}

func SetToPrefix(prefix string) {
	if globalCtrl != nil {
		globalCtrl.SetToPrefix(prefix)
	}
}

func (c *Ctrl) GetTenantId(subject string) string {
	return findTenantId(subject, c.fromPrefix)
}

func (c *Ctrl) GetTenantID(subject string) string {
	return findTenantId(subject, c.fromPrefix)
}

func (c *Ctrl) WithTenantAddress(tenant string, nodeUUID string) client.CallOption {
	address := c.TenantNodeAddress(tenant, nodeUUID)
	if tenant == "" {
		return WithAddress(address)
	}
	return client.WithAddress(address)
}

func (c *Ctrl) TenantNodeAddress(tenant string, nodeUUID string) string {
	if tenant == "" {
		return NodeAddress(nodeUUID)
	}

	prefix := c.toPrefix + tenant + "."
	address := ""
	if nodeUUID == "" {
		address = prefix + "cn.xswitch.node"
	} else {
		if !strings.HasPrefix(nodeUUID, "cn.xswitch.") {
			address = prefix + "cn.xswitch.node." + nodeUUID
		} else {
			address = prefix + nodeUUID
		}
	}
	return address
}
