package ctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"

	"git.xswitch.cn/xswitch/xctrl/xctrl/client"
	"git.xswitch.cn/xswitch/xctrl/xctrl/util/log"

	"git.xswitch.cn/xswitch/xctrl/ctrl/bus"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/proto/cman"
	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
)

// register 注册node节点
func (h *Ctrl) register(n *xctrl.Node) {
	// 节点注册
	nodes.Store(n.Name, n)
}

// deRegister 取消节点注册
func (h *Ctrl) deRegister(n *xctrl.Node) {
	// 节点离线
	nodes.Delete(n.Name)
}

// NodeRegister 节点注册
func (h *Ctrl) nodeRegister(ctx context.Context, frame *json.RawMessage) error {
	n := &xctrl.Node{}
	err := json.Unmarshal(*frame, n)
	if err != nil {
		fmt.Errorf("jsonrpc parse error:%s", err)
	}

	log.Tracef("Node Register: %s, %s, %s, %s, %d",
		n.GetName(),
		n.GetVersion(),
		n.GetUuid(),
		n.GetAddress(),
		n.GetRack())

	h.register(n)
	return nil
}

// nodeUnregister 节点离线
func (h *Ctrl) nodeUnregister(ctx context.Context, frame *json.RawMessage) error {
	n := &xctrl.Node{}
	err := json.Unmarshal(*frame, n)
	if err != nil {
		fmt.Errorf("jsonrpc parse error:%v", err)
		return nil
	}

	log.Tracef("Node Unregister: %s, %s, %s, %s, %d",
		n.GetName(),
		n.GetVersion(),
		n.GetUuid(),
		n.GetAddress(),
		n.GetRack())

	h.deRegister(n)
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

// handleNode 节点事件响应
func (h *Ctrl) handleNode(ctx context.Context, data nats.Event) error {
	var message Request
	err := json.Unmarshal(data.Message().Body, &message)
	if err != nil {
		fmt.Errorf("event parse error:%v", err)
		return nil
	}
	switch message.Method {
	case "Event.NodeRegister":
		h.nodeRegister(ctx, message.Params)
	case "Event.NodeUnregister":
		h.nodeUnregister(ctx, message.Params)
	case "Event.NodeUpdate":
		h.nodeUpdate(ctx, message.Params)
	default:
		fmt.Printf("Received unsupported event: %s\n", message.Method)
	}

	if h.enableNodeStatus && h.handler != nil {
		h.handler.Event(ctx, `cn.xswitch.node.status`, &message)
	}
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

// 监听节点注册事件
func (h *Ctrl) bindNodeStatus(subject string) error {
	_, err := h.conn.Subscribe(subject, func(ev nats.Event) error {
		return h.handleNode(context.Background(), ev)
	})
	if err != nil {
		fmt.Errorf("topic subscribe error: %v", err.Error())
		return err
	}
	return nil
}

// handleRpc rpc请求路由
func (h *Ctrl) handleRequest(ctx context.Context, data nats.Event) error {
	message := &Request{}
	err := json.Unmarshal(data.Message().Body, message)
	if err != nil {
		fmt.Errorf("event parse error:%v", err)
		return nil
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Errorf("%s | %s ==> [painc] %v\n%s", data.Topic(), message.Method, err, string(debug.Stack()))
			}
		}()

		if h.handler != nil {
			h.handler.Request(ctx, data.Topic(), data.Reply(), message)
		}
	}()
	return nil
}

// handleChannel 处理channel 事件
func (h *Ctrl) handleChannel(ctx context.Context, data nats.Event, message *Message) error {
	if message.Params == nil {
		fmt.Print("recved nil param ", message.Params)
		return nil
	}
	// channel
	channel := new(Channel)
	err := json.Unmarshal(*message.Params, channel)
	if err != nil {
		fmt.Errorf("%s: application json parse error %+v", data.Topic(), err.Error())
		return nil
	}

	timeout := 4*time.Hour + 10*time.Minute // make sure timeout is bigger than call duration
	switch channel.GetState() {
	case "START":
		fmt.Printf("%s START", channel.GetUuid())
		bus.SubscribeWithExpire(channel.GetUuid(), channel.GetUuid(), timeout, func(ev *bus.Event) error {
			if h.handler != nil {
				data := ev.Params.(nats.Event)
				message := ev.Message.(*Message)
				h.handler.App(ctx, data.Topic(), data.Reply(), message)
			}
			if ev.Flag == "DESTROY" || ev.Flag == "TIMEOUT" {
				bus.Unsubscribe(ev.Topic, ev.Queue)
			}
			return nil
		})
		fmt.Printf("%s Subscribered", channel.GetUuid())
	case "CALLING":
		bus.SubscribeWithExpire(channel.GetUuid(), channel.GetUuid(), timeout, func(ev *bus.Event) error {
			if h.handler != nil {
				data := ev.Params.(nats.Event)
				message := ev.Message.(*Message)
				h.handler.App(ctx, data.Topic(), data.Reply(), message)
			}
			if ev.Flag == "DESTROY" || ev.Flag == "TIMEOUT" {
				bus.Unsubscribe(ev.Topic, ev.Queue)
			}
			return nil
		})
	}

	ev := bus.NewEvent(channel.GetState(), channel.GetUuid(), message, data)
	bus.Publish(ev)
	return nil
}

// handleApp app路由
func (h *Ctrl) handleApp(ctx context.Context, data nats.Event) error {
	var message Message
	err := json.Unmarshal(data.Message().Body, &message)
	if err != nil {
		fmt.Errorf("event parse error:%v", err)
		return nil
	}

	switch message.Method {
	case "Event.Channel":
		h.handleChannel(ctx, data, &message)
	default:
		defer func() {
			if err := recover(); err != nil {
				fmt.Errorf("%s | %s ==> [painc] %v\n%s", data.Topic(), message.Method, err, string(debug.Stack()))
			}
		}()

		if h.forkDTMFEvent && message.Method == "Event.DTMF" {
			// fork this event and deliver to the Event Channel Thread
			h.handleChannel(ctx, data, &message)
		}

		if h.handler != nil {
			h.handler.App(ctx, data.Topic(), data.Reply(), &message)
		}
	}

	return nil
}

// handleApp result 路由
func (h *Ctrl) handleResult(ctx context.Context, data nats.Event) error {
	go func() {
		var result Result
		err := json.Unmarshal(data.Message().Body, &result)
		if err != nil {
			fmt.Errorf("event parse error:%v", err)
			return
		}

		defer func() {
			if err := recover(); err != nil {
				fmt.Errorf("%s | %s ==> [painc] %v\n%s", data.Topic(), *result.ID, err, string(debug.Stack()))
			}
		}()

		if h.handler != nil {
			h.handler.Result(ctx, data.Topic(), &result)
		}
	}()
	return nil
}

func processEvent(c *Ctrl, ctx context.Context, data nats.Event) {
	var message Request

	err := json.Unmarshal(data.Message().Body, &message)
	if err != nil {
		fmt.Errorf("event parse error:%v", err)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("%s | %s ==> [painc] %v\n%s", data.Topic(), message.Method, err, string(debug.Stack()))
		}
	}()

	if c.handler != nil {
		c.handler.Event(ctx, data.Topic(), &message)
	}
}

// handleEvent event路由
func (c *Ctrl) handleEvent(ctx context.Context, data nats.Event) error {
	if strings.HasPrefix(data.Topic(), "cn.xswitch.event.cdr") {
		processEvent(c, ctx, data)
		return nil
	}
	go func() {
		processEvent(c, ctx, data)
	}()
	return nil
}

// EnableApp APP事件
func (h *Ctrl) EnableApp(topic string) error {
	// 呼叫控制
	// cn.xswitch.app.callcenter 呼叫队列
	// cn.xswitch.app.autodialer 预测外呼
	_, err := h.conn.Subscribe(topic, func(ev nats.Event) error {
		return h.handleApp(context.Background(), ev)
	}, nats.Queue(`cn.xswitch.app`))
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}
	_, err = h.conn.Subscribe(fmt.Sprintf(`%s.%s`, topic, h.uuid), func(ev nats.Event) error {
		return h.handleApp(context.Background(), ev)
	}, nats.Queue(`cn.xswitch.app`))
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}
	return nil
}

// EnableApp Result 事件
func (h *Ctrl) EnableResult(topic string) error {
	_, err := h.conn.Subscribe(fmt.Sprintf(`%s.%s`, topic, h.uuid), func(ev nats.Event) error {
		return h.handleResult(context.Background(), ev)
	}, nats.Queue(`cn.xswitch.result`))
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}
	return nil
}

// EnableRequest 开启Request请求监听
func (h *Ctrl) EnableRequest(topic string) error {
	// fetchXMl, Dialplan
	_, err := h.conn.Subscribe(topic, func(ev nats.Event) error {
		return h.handleRequest(context.Background(), ev)
	}, nats.Queue(`cn.xswitch.request`))
	if err != nil {
		log.Errorf("topic subscribe error: %s", err.Error())
		return err
	}
	return nil
}

// EnableEvent 开启事件监听
func (h *Ctrl) EnableEvent(topic string, queue string) error {
	// 例如
	// cn.xswitch.event.cdr
	// cn.xswitch.event.custom.sofia
	if queue == "" {
		queue = "cn.xswitch.event"
	}
	_, err := h.conn.Subscribe(topic, func(ev nats.Event) error {
		return h.handleEvent(context.Background(), ev)
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
func (h *Ctrl) EnbaleNodeStatus() error {
	// 例如
	// cn.xswitch.node.status
	h.enableNodeStatus = true
	return nil
}

// ForkDTMFEventToChannelEventThread
func (h *Ctrl) ForkDTMFEventToChannelEventThread() error {
	h.forkDTMFEvent = true
	return nil
}

func initCtrl(handler Handler, trace bool, subject string, addrs ...string) (*Ctrl, error) {
	h := &Ctrl{
		conn:             nats.NewConn(nats.Addrs(addrs...), nats.Trace(trace)),
		uuid:             uuid.New().String(),
		serviceName:      "cn.xswitch.nodes",
		handler:          handler,
		enableNodeStatus: false,
		channelHub:       map[string]*Channel{},
		resultCallbacks:  map[string]*AsyncCallOption{},
	}

	// 连接消息队列
	if err := h.conn.Connect(); err != nil {
		return nil, err
	}

	// 监听节点状态事件
	if err := h.bindNodeStatus(subject); err != nil {
		return nil, err
	}

	// 同步调用 xswitch
	h.service = h.newNodeService()
	// 异步调用 xswitch
	h.asyncService = h.newAsyncService()
	// 同步调用 xswitch, 使用nats的RequestWithContext, 可以返回结果，可以中途取消
	h.aService = h.newAService()
	return h, nil
}

// 订阅消息
func (h *Ctrl) Subscribe(topic string, cb nats.EventCallback, queue string) (nats.Subscriber, error) {
	sub, err := h.conn.Subscribe(topic, func(ev nats.Event) error {
		return cb(context.Background(), ev)
	}, nats.Queue(queue))
	if err != nil {
		fmt.Errorf("topic %s subscribe error: %+v", topic, err.Error())
	}
	return sub, err
}
