package bus

import (
	"fmt"
	"runtime"
	"time"
)

const (
	queueBufferSize = 10240
	busBufferSize   = 102400
)

// Subscriber  订阅者
type Subscriber struct {
	// 订阅的主题
	subject string
	// 队列
	queue string
	// 自动过期
	expire time.Duration
	// 消息Handler
	handler Handler
	// action register|unregister|event
	action string
	// event, only valid if action == event
	event *Event
}

func (s *Subscriber) runWithRecovery(e *Event) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
		}
	}()
	if s.handler != nil {
		s.handler(e)
	}
}

// newSubscriber 创建订阅者
func newSubscriber(subject string, queue string, h Handler) *Subscriber {
	return &Subscriber{
		subject: subject,
		queue:   queue,
		handler: h,
	}
}

// subject 事件订阅者组
type subscriberGroup struct {
	subject     string
	subscribers map[string]*Subscriber
}

// add 添加订阅者
func (g *subscriberGroup) Add(sub *Subscriber) {
	g.subscribers[sub.queue] = sub
}

// add 删除订阅者
func (g *subscriberGroup) Del(sub *Subscriber) {
	delete(g.subscribers, sub.queue)
}

// newSubscriberGroup 创建订阅组
func newSubscriberGroup(subject string) *subscriberGroup {
	return &subscriberGroup{
		subject:     subject,
		subscribers: make(map[string]*Subscriber, 0),
	}
}

type eventBus struct {
	subscribers map[string]*subscriberGroup
	queues      map[string]*queue
	ch          chan *Subscriber
}

var bus *eventBus

func init() {
	bus = newBus()
	bus.start()
}

func newBus() *eventBus {
	b := &eventBus{
		subscribers: make(map[string]*subscriberGroup, 0),
		queues:      make(map[string]*queue, 0),
		ch:          make(chan *Subscriber, busBufferSize),
	}

	return b
}

func (h *eventBus) mainLoop() {
	for {
		// xlog.Error("select")
		select {
		case subscriber, ok := <-h.ch:
			if !ok {
				return
			}

			switch subscriber.action {
			case "register":
				subgroup, ok := h.subscribers[subscriber.subject]
				if !ok {
					subgroup = newSubscriberGroup(subscriber.subject)
				}
				subgroup.Add(subscriber)
				h.subscribers[subscriber.subject] = subgroup

				if subscriber.queue != "" {
					_queue, _ok := h.queues[subscriber.queue]
					if _ok {
						_queue.addRef()
					} else {
						h.queues[subscriber.queue] = newQueue(subscriber.queue, subscriber.expire)
					}
				}
			case "unregister":

				subgroup, ok := h.subscribers[subscriber.subject]
				if ok {
					subgroup.Del(subscriber)
					if len(subgroup.subscribers) == 0 {
						delete(h.subscribers, subgroup.subject)
					} else {
						h.subscribers[subgroup.subject] = subgroup
					}
				}
				if subscriber.queue != "" {
					_queue, ok := h.queues[subscriber.queue]
					if ok {
						if _queue.release() < 1 {
							delete(h.queues, subscriber.queue)
						}
					}
				}
			case "event":
				ev := subscriber.event
				if subgroup, ok := h.subscribers[ev.Topic]; ok {
					for _, sub := range subgroup.subscribers {
						if sub.queue != "" {
							bus.publishToQueue(sub, ev)
						} else {
							go sub.runWithRecovery(ev)
						}
					}
				}
			}
		}
	}
}

func (h *eventBus) publishToQueue(s *Subscriber, ev *Event) {
	if q, ok := h.queues[s.queue]; ok {
		newEv := &Event{
			Flag:    ev.Flag,
			Topic:   s.subject,
			Message: ev.Message,
			Params:  ev.Params,
			Queue:   s.queue,
			handler: s.handler,
		}
		select {
		case q.inbound <- newEv:
		default:
		}
	}
}

func (h *eventBus) start() {
	go h.mainLoop()
}

// Event is given to a subscription handler for processing
type Event struct {
	Flag    string      `json:"flag,omitempty"`
	Topic   string      `json:"topic,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Queue   string      `json:"queue,omitempty"`
	handler Handler     `json:"-"`
}

func NewEvent(flag string, topic string, data interface{}, params interface{}) *Event {
	e := &Event{
		Flag:    flag,
		Message: data,
		Topic:   topic,
		Params:  params,
	}
	return e
}

// Publish 发布事件
func Publish(ev *Event) {
	s := &Subscriber{
		action: "event",
		event:  ev,
	}

	select {
	case bus.ch <- s:
	default:
		fmt.Errorf("event inbound chan block drop event")
	}
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(*Event) error

// Subscribe 订阅事件
func Subscribe(topic string, queue string, h Handler) error {
	if h != nil {
		s := &Subscriber{
			subject: topic,
			queue:   queue,
			handler: h,
			action:  "register",
		}
		bus.ch <- s
		return nil
	}
	return fmt.Errorf("handler must not be nil")
}

func SubscribeWithExpire(topic string, queue string, expire time.Duration, h Handler) error {
	if h != nil {
		s := &Subscriber{
			subject: topic,
			queue:   queue,
			expire:  expire,
			handler: h,
			action:  "register",
		}
		bus.ch <- s
		return nil
	}
	return fmt.Errorf("handler must not be nil")
}

// Unsubscribe 取消订阅
func Unsubscribe(topic string, queue string) {
	s := &Subscriber{
		subject: topic,
		queue:   queue,
		action:  "unregister",
	}
	bus.ch <- s
}
