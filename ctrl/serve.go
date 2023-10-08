package ctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"

	"git.xswitch.cn/xswitch/proto/xctrl/client"
	"git.xswitch.cn/xswitch/proto/xctrl/util/log"

	"git.xswitch.cn/xswitch/proto/go/proto/cman"
	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/bus"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
)

// register 注册node节点
func (h *Ctrl) register(n *xctrl.Node) {
	h.nodes.Store(n.Name, n)
}

// deRegister 取消节点注册
func (h *Ctrl) deRegister(n *xctrl.Node) {
	h.nodes.Delete(n.Name)
}

// handleNode 节点事件响应
func (h *Ctrl) handleNode(natsEvent nats.Event) error {
	var event Request
	err := json.Unmarshal(natsEvent.Message().Body, &event)
	if err != nil {
		return fmt.Errorf("event parse error: %v", err)
	}
	node := &xctrl.Node{}
	err = json.Unmarshal(*event.Params, node)
	if err != nil {
		return fmt.Errorf("jsonrpc parse error: %v", err)
	}
	log.Tracef("Node Register: %s, %s, %s, %s, %d",
		node.GetName(),
		node.GetVersion(),
		node.GetUuid(),
		node.GetAddress(),
		node.GetRack())
	isMethodForNode := true
	switch event.Method {
	case "Event.NodeRegister":
		h.register(node)
	case "Event.NodeUnregister":
		h.deRegister(node)
	case "Event.NodeUpdate":
		h.register(node)
	default:
		isMethodForNode = false
		log.Warnf("Received unsupported event: %s\n", event.Method)
	}
	if isMethodForNode && h.nodeCallback != nil {
		node := new(xctrl.Node)
		err = json.Unmarshal(*event.Params, node)
		if err != nil {
			log.Error(err)
		} else {
			h.nodeCallback(node, event.Method)
		}
	}
	return nil
}

// nodeUpdate 节点状态更新
func (h *Ctrl) nodeUpdate(ctx context.Context, frame *json.RawMessage) error {
	n := &xctrl.Node{}
	err := json.Unmarshal(*frame, n)
	if err != nil {
		fmt.Errorf("jsonrpc parse error:%v", err)
		return nil
	}
	log.Tracef("Node Status: %s, %s, %s, %s, %d",
		n.GetName(),
		n.GetVersion(),
		n.GetUuid(),
		n.GetAddress(),
		n.GetRack())
	h.register(n)
	return nil
}

// newNodeService 创建 XNodeService
func (h *Ctrl) newNodeService() xctrl.XNodeService {
	c := newClient(h.conn, false)
	c.Init(client.Selector())
	return xctrl.NewXNodeService(h.serviceName, c)
}

// newAsyncService 异步调用 XNodeService
func (h *Ctrl) newAsyncService() xctrl.XNodeService {
	c := newClient(h.conn, true)
	c.Init(client.Selector())
	return xctrl.NewXNodeService(h.serviceName, c)
}

// newAService 同步调用 XNodeService，但可以使用context控制timeout
func (h *Ctrl) newAService() xctrl.XNodeService {
	c := newClient(h.conn, false)
	c.Init(client.Selector())
	c.SetAService()
	return xctrl.NewXNodeService(h.serviceName, c)
}

// NewCManService 创建 CManService
func (h *Ctrl) NewCManService(addr string) cman.CManService {
	c := newClient(h.conn, false, client.ServiceAddr(addr))
	c.Init(client.Selector())
	h.cmanService = cman.NewCManService(h.serviceName, c)
	return h.cmanService
}

// handleRpc rpc请求路由
func (h *Ctrl) handleRequest(handler RequestHandler, event nats.Event) error {
	message := &Request{}
	err := json.Unmarshal(event.Message().Body, message)
	if err != nil {
		fmt.Errorf("event parse error:%v", err)
		return nil
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Errorf("%s | %s ==> [painc] %v\n%s", event.Topic(), message.Method, err, string(debug.Stack()))
			}
		}()
		handler.Request(message, event)
	}()
	return nil
}

// handleChannel 处理channel 事件
func (h *Ctrl) handleChannel(handler AppHandler, message *Message, natsEvent nats.Event) error {
	if message.Params == nil {
		fmt.Print("recved nil param ", message.Params)
		return nil
	}
	channel := new(Channel)
	channel.userData = nil
	if message.Method == "Event.DTMF" {
		dtmfEvent := new(xctrl.DTMFEvent)
		err := json.Unmarshal(*message.Params, dtmfEvent)
		if err != nil {
			return fmt.Errorf("%s: application json parse error %+v", natsEvent.Topic(), natsEvent.Error())
		}
		channel.ChannelEvent = new(xctrl.ChannelEvent)
		channel.Dtmf = dtmfEvent.DtmfDigit
		channel.NodeUuid = dtmfEvent.NodeUuid
		channel.Uuid = dtmfEvent.Uuid
		channel.State = "DTMF"
	} else {
		err := json.Unmarshal(*message.Params, channel)
		if err != nil {
			return fmt.Errorf("%s: application json parse error %+v", natsEvent.Topic(), natsEvent.Error())
		}
	}
	if channel.ChannelEvent != nil {
		subject := natsEvent.Topic()
		tenantID := findTenantId(subject, h.fromPrefix)
		if tenantID != "" {
			channel.ChannelEvent.SetTenantID(tenantID)
		}
		if h.toPrefix != "" {
			channel.ChannelEvent.SetToPrefix(h.toPrefix)
		}
	}

	timeout := time.Duration(h.maxChannelLifeTime)*time.Hour + 10*time.Minute // make sure timeout is bigger than call duration
	switch channel.GetState() {
	case "START":
		log.Tracef("%s START", channel.GetUuid())
		var key ContextKey = "channel"
		ctx := context.WithValue(context.Background(), key, channel)
		bus.SubscribeWithExpire(channel.GetUuid(), channel.GetUuid(), timeout, func(ev *bus.Event) error {
			if ev.Params == nil {
				log.Error("ev.Params is nil")
				return nil
			}
			if natsEvent, ok := ev.Params.(nats.Event); ok {
				channelEvent := ev.Message.(*Channel)
				channelEvent.natsEvent = natsEvent
				handler.ChannelEvent(ctx, channelEvent)
				if ev.Flag == "DESTROY" || ev.Flag == "TIMEOUT" {
					bus.Unsubscribe(ev.Topic, ev.Queue)
				}
				return nil
			}
			return nil
		})
		log.Tracef("%s Subscribered", channel.GetUuid())

	case "CALLING":
		log.Tracef("%s CALLING", channel.GetUuid())
		var key ContextKey = "channel"
		ctx := context.WithValue(context.Background(), key, channel)
		bus.SubscribeWithExpire(channel.GetUuid(), channel.GetUuid(), timeout, func(ev *bus.Event) error {
			if ev.Params == nil {
				log.Error("ev.Params is nil")
				return nil
			}
			if natsEvent, ok := ev.Params.(nats.Event); ok {
				channelEvent := ev.Message.(*Channel)
				channelEvent.natsEvent = natsEvent
				handler.ChannelEvent(ctx, channelEvent)
				if ev.Flag == "DESTROY" || ev.Flag == "TIMEOUT" {
					bus.Unsubscribe(ev.Topic, ev.Queue)
				}
				return nil
			}
			return nil

		})
	default:
		log.Infof("Channel State %s %s", channel.GetUuid(), channel.GetState())
	}

	ev := bus.NewEvent(channel.GetState(), channel.GetUuid(), channel, natsEvent)
	bus.Publish(ev)
	return nil
}

// handleApp app路由
func (c *Ctrl) handleApp(handler AppHandler, natsEvent nats.Event) error {
	var message Message
	err := json.Unmarshal(natsEvent.Message().Body, &message)
	if err != nil {
		return fmt.Errorf("event parse error:%v", err)
	}

	switch message.Method {
	case "Event.Channel":
		c.handleChannel(handler, &message, natsEvent)
	default:
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("%s | %s ==> [painc] %v\n%s", natsEvent.Topic(), message.Method, err, string(debug.Stack()))
			}
		}()

		if c.forkDTMFEvent && message.Method == "Event.DTMF" {
			// fork this event and deliver to the Event Channel Thread
			c.handleChannel(handler, &message, natsEvent)
		}

		handler.Event(&message, natsEvent)
	}

	return nil
}

func processEvent(handler EventHandler, natsEvent nats.Event) {
	var message Request

	err := json.Unmarshal(natsEvent.Message().Body, &message)
	if err != nil {
		log.Errorf("event parse error:%v", err)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%s | %s ==> [painc] %v\n%s", natsEvent.Topic(), message.Method, err, string(debug.Stack()))
		}
	}()

	handler.Event(&message, natsEvent)
}

// handleEvent event路由
func (c *Ctrl) handleEvent(handler EventHandler, natsEvent nats.Event) error {
	if strings.HasPrefix(natsEvent.Topic(), "cn.xswitch.event.cdr") {
		processEvent(handler, natsEvent)
		return nil
	}
	go func() {
		processEvent(handler, natsEvent)
	}()
	return nil
}

// EnableApp APP事件
func (h *Ctrl) EnableApp(handler AppHandler, subject string, queue string) error {
	log.Infof("EnableApp subject=%s queue=%s", subject, queue)
	_, err := h.conn.Subscribe(subject, func(ev nats.Event) error {
		return h.handleApp(handler, ev)
	}, nats.Queue(queue))
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}

	mySubject := fmt.Sprintf(`%s.%s`, subject, h.uuid)
	log.Infof("EnableApp subscribe to subject=%s", mySubject)
	_, err = h.conn.Subscribe(mySubject, func(ev nats.Event) error {
		return h.handleApp(handler, ev)
	})
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}

	mySubject = fmt.Sprintf(`%s.%s`, "cn.xswitch.ctrl", h.uuid)
	log.Infof("EnableApp subscribe to subject=%s", mySubject)
	_, err = h.conn.Subscribe(mySubject, func(ev nats.Event) error {
		return h.handleApp(handler, ev)
	})
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}

	return nil
}

// EnableRequest 开启Request请求监听
func (h *Ctrl) EnableRequest(handler RequestHandler, subject string, queue string) error {
	// fetchXMl, Dialplan
	log.Infof("EnableRequest subject=%s queue=%s", subject, queue)
	_, err := h.conn.Subscribe(subject, func(ev nats.Event) error {
		return h.handleRequest(handler, ev)
	}, nats.Queue(queue))
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}
	return nil
}

// EnableEvent 开启事件监听
func (h *Ctrl) EnableEvent(handler EventHandler, subject string, queue string) error {
	// 例如
	// cn.xswitch.event.cdr
	// cn.xswitch.event.custom.sofia
	log.Infof("EnableEvent subject=%s queue=%s", subject, queue)
	_, err := h.conn.Subscribe(subject, func(ev nats.Event) error {
		return h.handleEvent(handler, ev)
	}, nats.Queue(queue))
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}

	// // uncomment to test a slow consumer which is default to 65535, 65535 * 1024
	// if err := sub.SetPendingLimits(10, 1024*10); err != nil {
	// 	log.Fatalf("Unable to set pending limits: %v", err)
	// }

	return nil
}

// EnbaleNodeStatus 开启节点监听
func (h *Ctrl) EnbaleNodeStatus(subject string) error {
	if subject == "" {
		subject = "cn.xswitch.status.node"
	}
	log.Infof("EnableNodeStatus subject=%s", subject)
	_, err := h.conn.Subscribe(subject, func(ev nats.Event) error {
		return h.handleNode(ev)
	})
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}
	h.enableNodeStatus = true
	return nil
}

func (h *Ctrl) OnEvicted(f func(string, interface{})) {

	h.nodes.store.OnEvicted(f)

}

// ForkDTMFEventToChannelEventThread
func (h *Ctrl) ForkDTMFEventToChannelEventThread() error {
	h.forkDTMFEvent = true
	return nil
}

func initCtrl(trace bool, addrs ...string) (*Ctrl, error) {
	c := &Ctrl{
		conn:               nats.NewConn(nats.Addrs(addrs...), nats.Trace(trace)),
		uuid:               uuid.New().String(),
		serviceName:        "cn.xswitch.nodes",
		enableNodeStatus:   false,
		channelHub:         map[string]*Channel{},
		resultCallbacks:    map[string]*AsyncCallOption{},
		maxChannelLifeTime: 4,
	}

	// 连接NATS消息队列
	if err := c.conn.Connect(); err != nil {
		return nil, err
	}

	// 同步调用 xswitch
	c.service = c.newNodeService()
	// 异步调用 xswitch
	c.asyncService = c.newAsyncService()
	// 同步调用 xswitch, 使用nats的RequestWithContext, 可以返回结果，可以中途取消
	c.aService = c.newAService()

	c.nodes = InitCtrlNodes()
	return c, nil
}

// 订阅消息
func (h *Ctrl) Subscribe(topic string, cb nats.EventCallback, queue string) (nats.Subscriber, error) {
	sub, err := h.conn.Subscribe(topic, func(ev nats.Event) error {
		return cb(context.Background(), ev)
	}, nats.Queue(queue))
	if err != nil {
		return nil, fmt.Errorf("topic %s subscribe error: %+v", topic, err.Error())
	}
	return sub, err
}

type NodeHashFun func(node *xctrl.Node, method string)

// RegisterHashNodeFun 注册hash节点事件
// nodeCallbackFunc 节点事件方法
func (h *Ctrl) registerHashNodeFun(nodeCallbackFunc NodeHashFun) {
	h.nodeCallback = nodeCallbackFunc
}

func before(str1, str2 string) string {
	index := strings.Index(str1, str2)
	if index == -1 {
		return ""
	}
	return str1[:index]
}

func findTenantId(str string, fromPrefix string) string {
	s := before(str, ".cn.xswitch")
	if s == "" {
		return ""
	}
	if fromPrefix == "" {
		return s
	}
	if strings.HasPrefix(s, fromPrefix) {
		return s[len(fromPrefix):]
	}
	return ""
}
