package bus

import (
	"encoding/json"
	"fmt"
	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBusNoQueue(t *testing.T) {
	got := 0

	Subscribe("test-topic", "", func(ev *Event) error {
		got++
		t.Log(ev.Message)
		return nil
	})

	ev := NewEvent("", "test-topic", "message", "data")
	Publish(ev)

	ev1 := NewEvent("", "test-topic", "message2", "data")
	Publish(ev1)

	time.Sleep(100 * time.Millisecond)

	if got != 2 {
		t.Error(got)
	}

	if len(bus.subscribers) != 1 {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(bus.queues) != 0 {
		t.Errorf("%+v", bus.queues)
	}

	Unsubscribe("test-topic", "")
	time.Sleep(200 * time.Millisecond)

	if len(bus.subscribers) != 0 {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(bus.queues) != 0 {
		t.Errorf("%+v", bus.queues)
	}
}

func TestBusQueue(t *testing.T) {
	got := 0

	Subscribe("test-topic", "test-queue", func(ev *Event) error {
		got++
		return nil
	})

	ev := NewEvent("Flag", "test-topic", "message", "data")
	Publish(ev)
	Publish(ev)

	time.Sleep(100 * time.Millisecond)

	if got != 2 {
		t.Error(got)
	}

	if len(bus.subscribers) != 1 {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(bus.queues) != 1 {
		t.Errorf("%+v", bus.queues)
	}

	Unsubscribe("test-topic", "test-queue")
	time.Sleep(200 * time.Millisecond)

	if len(bus.subscribers) != 0 {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(bus.queues) != 0 {
		t.Errorf("%+v", bus.queues)
	}
}

func TestBusQueueExpire(t *testing.T) {
	got := 0
	expired := 0

	SubscribeWithExpire("test-topic", "test-queue", 80*time.Millisecond, func(ev *Event) error {
		t.Log(ev)
		if ev.Flag == "TIMEOUT" {
			t.Log("TIMEOUT")
			expired++
			Unsubscribe("test-topic", "test-queue")
			return nil
		}
		got++
		return nil
	})

	ev := NewEvent("Flag", "test-topic", "message", "data")
	Publish(ev)
	Publish(ev)

	time.Sleep(100 * time.Millisecond)

	if got != 2 {
		t.Error(got)
	}

	if expired != 1 {
		t.Error(expired)
	}

	if len(bus.subscribers) != 0 {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(bus.queues) != 0 {
		t.Errorf("%+v", bus.queues)
	}
}

func TestBus2QueueSub(t *testing.T) {
	got := 0

	Subscribe("test-topic", "test-queue", func(ev *Event) error {
		got++
		t.Logf("got %+v", ev.Message)
		return nil
	})

	got2 := 0

	Subscribe("test-topic", "test-queue", func(ev *Event) error {
		got2++
		t.Logf("got %+v", ev.Message)
		return nil
	})

	ev := NewEvent("Flag", "test-topic", "message", "data")
	Publish(ev)
	ev = NewEvent("Flag", "test-topic", "message2", "data")
	Publish(ev)

	time.Sleep(200 * time.Millisecond)

	if got+got2 != 2 {
		t.Error(got)
		t.Error(got2)
	}

	if len(bus.subscribers) != 1 {
		t.Errorf("%+v", bus.subscribers)
	}

	subs, found := bus.subscribers["test-topic"]

	if !found {
		t.Error(bus.subscribers)
	}

	if len(subs.subscribers) != 1 {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(bus.queues) != 1 {
		t.Errorf("%+v", bus.queues)
	}

	Unsubscribe("test-topic", "test-queue")
	Unsubscribe("test-topic", "test-queue")
	time.Sleep(500 * time.Millisecond)

	if len(bus.subscribers) != 0 {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(bus.queues) != 0 {
		t.Errorf("%+v", bus.queues)
	}
}

func TestBus2Topic(t *testing.T) {
	got := 0

	Subscribe("test-topic", "test-queue", func(ev *Event) error {
		got++
		t.Logf("got %+v", ev.Message)
		return nil
	})

	got2 := 0

	Subscribe("test-topic2", "test-queue", func(ev *Event) error {
		got2++
		t.Logf("got %+v", ev.Message)
		return nil
	})

	ev := NewEvent("Flag", "test-topic", "message", "data")
	Publish(ev)
	ev = NewEvent("Flag", "test-topic2", "message2", "data")
	Publish(ev)

	time.Sleep(100 * time.Millisecond)

	if got != 1 {
		t.Error(got)
	}

	if got2 != 1 {
		t.Error(got2)
	}

	if len(bus.subscribers) != 2 {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(bus.queues) != 1 {
		t.Errorf("%+v", bus.queues)
	}

	Unsubscribe("test-topic", "test-queue")
	Unsubscribe("test-topic2", "test-queue")
	time.Sleep(500 * time.Millisecond)

	if len(bus.subscribers) != 0 {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(bus.queues) != 0 {
		t.Errorf("%+v", bus.queues)
	}
}

func BenchmarkQueue(b *testing.B) {
	start := time.Now()
	N := 1 * 1000
	Z := 16

	var seq uint64
	done := make(chan bool, 16*N)
	b.ResetTimer()

	for i := 0; i < N; i++ {
		go func(i int) {
			queue := fmt.Sprintf("acd-%d", i)
			Subscribe("app.acd.dial", queue, func(e *Event) error {
				done <- true
				return nil
			})

			for i := 0; i < Z; i++ {
				requestID := fmt.Sprintf(`"%d"`, atomic.AddUint64(&seq, 1))
				id := json.RawMessage(requestID)
				channel := &xctrl.ChannelEvent{
					NodeUuid: "xcc-node-3",
					Uuid:     "1d67d67f-45bf-4377-a6c1-4e1e03185326",
					State:    "DESTROY",
					Params: map[string]string{
						"XCC-Control-UUID":       "app.dial.6eeb469e-234a-47e8-bd54-aa9e7291409b",
						"xcc_mark":               "Dialing",
						"hangup_cause":           "NORMAL_CLEARING",
						"duration":               "66",
						"context":                "xcc_bridge,park",
						"xcc_stationtype":        "SIP",
						"xcc_origin_cid_number":  "1003",
						"cc_agent":               "c296c363-c17d-4439-8f19-8d780ce9fdac",
						"xcc_direction":          "inbound",
						"cc_agent_session_uuid":  "2045dbbf-8ce5-42ea-ae41-1f6b0fc65fa0",
						"xcc_identity":           "callee",
						"xcc_origin_dest_number": "15223352450",
						"xcc_session":            "f52f710b-f36a-425b-bfa9-bb40c3ec4499",
						"cc_side":                "member",
						"billsec":                "60",
						"xcc_domain":             "dev.xswitch.cn",
						"sip_hangup_disposition": "recv_bye",
					},
				}

				b, _ := json.Marshal(channel)

				msg := Request{
					Version: "2.0",
					Method:  "Event.Channel",
					Params:  RawMessage(b),
					ID:      &id,
				}

				e := NewEvent("", "app.acd.dial",
					&nats.Message{
						Body: msg.Marshal(),
					}, nil,
				)
				Publish(e)
			}
		}(i)
	}

	b.Log("watiting for done")
	count := 0
	for {
		<-done
		count = count + 1
		if count >= Z*N {
			break
		}
	}
	end := time.Now().Sub(start)
	b.Logf("Benchmark in queue %d %v %.4f cps", N*Z, end, float64(N*Z)/float64(end.Seconds()))
	b.ReportAllocs()
	return

}

// Request RPC 请求对象
type Request struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
	ID      *json.RawMessage `json:"id"`
}

// RawMessage .
func RawMessage(data []byte) *json.RawMessage {
	raw := json.RawMessage(data)
	return &raw
}

func (r *Request) RawMessage() *json.RawMessage {
	b, _ := json.Marshal(r)
	raw := json.RawMessage(b)
	return &raw
}

func (r *Request) Marshal() []byte {
	b, _ := json.Marshal(r)
	return b
}

func BenchmarkGeneral(b *testing.B) {
	start := time.Now()
	N := 16 * 100000

	var seq uint64

	wg := &sync.WaitGroup{}

	Subscribe("app.acd.dial", "", func(e *Event) error {
		wg.Done()
		rand.Seed(time.Now().UnixNano())
		s := rand.Intn(10)
		time.Sleep(time.Millisecond * time.Duration(s))
		return nil
	})

	b.ResetTimer()
	for i := 0; i < N; i++ {
		wg.Add(1)
		requestID := fmt.Sprintf(`"%d"`, atomic.AddUint64(&seq, 1))
		id := json.RawMessage(requestID)

		channel := &xctrl.ChannelEvent{
			NodeUuid: "xcc-node-3",
			Uuid:     "1d67d67f-45bf-4377-a6c1-4e1e03185326",
			State:    "DESTROY",
			Params: map[string]string{
				"XCC-Control-UUID":       "app.dial.6eeb469e-234a-47e8-bd54-aa9e7291409b",
				"xcc_mark":               "Dialing",
				"hangup_cause":           "NORMAL_CLEARING",
				"duration":               "66",
				"context":                "xcc_bridge,park",
				"xcc_stationtype":        "SIP",
				"xcc_origin_cid_number":  "1003",
				"cc_agent":               "c296c363-c17d-4439-8f19-8d780ce9fdac",
				"xcc_direction":          "inbound",
				"cc_agent_session_uuid":  "2045dbbf-8ce5-42ea-ae41-1f6b0fc65fa0",
				"xcc_identity":           "callee",
				"xcc_origin_dest_number": "15223352450",
				"xcc_session":            "f52f710b-f36a-425b-bfa9-bb40c3ec4499",
				"cc_side":                "member",
				"billsec":                "60",
				"xcc_domain":             "dev.xswitch.cn",
				"sip_hangup_disposition": "recv_bye",
			},
		}

		b, _ := json.Marshal(channel)

		msg := Request{
			Version: "2.0",
			Method:  "Event.Channel",
			Params:  RawMessage(b),
			ID:      &id,
		}

		e := NewEvent("", "app.acd.dial",
			&nats.Message{
				Body: msg.Marshal(),
			}, nil,
		)
		go Publish(e)
	}
	wg.Wait()
	Unsubscribe("app.acd.dial", "")
	end := time.Now().Sub(start)
	b.Logf("Benchmark in queue %d %v %.4f cps", N*16, end, float64(N*16)/float64(end.Seconds()))
	b.ReportAllocs()
}

//订阅相同的topic和相同的队列

func TestMuchSubscribes(t *testing.T) {

	got1 := 0
	Subscribe("test-topic", "test-queue", func(event *Event) error {
		got1++
		t.Log(event.Message)
		return nil
	})

	got2 := 0
	Subscribe("test-topic", "test-queue", func(event *Event) error {
		got2++
		t.Log(event.Message)
		return nil
	})

	got3 := 0
	Subscribe("test-topic", "test-queue", func(event *Event) error {
		got3++
		t.Log(event.Message)
		return nil
	})

	got4 := 0
	Subscribe("test-topic", "test-queue", func(event *Event) error {
		got4++
		t.Log(event.Message)
		return nil
	})

	time.Sleep(100 * time.Millisecond)

	if len(bus.subscribers) != 1 {
		t.Errorf("%+v", len(bus.subscribers))
	}

	if len(bus.queues) != 1 {
		t.Errorf("%+v", len(bus.queues))
	}

	ev1 := NewEvent("", "test-topic", "{'name':'topic1'}", "data")
	Publish(ev1)
	Publish(ev1)
	Publish(ev1)
	Publish(ev1)

	time.Sleep(100 * time.Millisecond)

	if got1 != 0 {
		t.Errorf("%+v", 0)
	}

	if got2 != 0 {
		t.Errorf("%+v", got2)
	}

	if got3 != 0 {
		t.Errorf("%+v", got3)
	}

	if got4 != 4 {
		t.Errorf("%+v", got4)
	}

	got1, got2, got3, got4 = 0, 0, 0, 0

	subs, found := bus.subscribers["test-topic"]

	if !found {
		t.Errorf("%+v", bus.subscribers)
	}

	if len(subs.subscribers) != 1 {
		t.Errorf("%+v", subs.subscribers)
	}

	if len(bus.queues) != 1 {
		t.Errorf("%+v", bus.queues)
	}

	Unsubscribe("test-topic", "test-queue")

	ev2 := NewEvent("", "test-topic", "{'name':'topic2'}", "data")
	Publish(ev2)

	ev3 := NewEvent("", "test-topic", "{'name':'topic3'}", "data")
	Publish(ev3)

	time.Sleep(100 * time.Millisecond)
	if got1 != 0 {
		t.Errorf("%+v", got1)
	}

	if got2 != 0 {
		t.Errorf("%+v", got2)
	}

	if got3 != 0 {
		t.Errorf("%+v", got3)
	}

	if got4 != 0 {
		t.Errorf("%+v", got4)
	}

	fmt.Println("register channel buffer len:", len(bus.ch))
}

func BenchmarkMuchSubscribes(b *testing.B) {
	start := time.Now()

	for i := 0; i < b.N; i++ {
		Subscribe("test-topic", "test-queue", func(event *Event) error {
			b.Log(event.Message)
			return nil
		})
	}
	for i := 0; i < b.N; i++ {
		ev := NewEvent("", "test-topic", "{'name':'bar'}", "data")
		Publish(ev)
	}

	Unsubscribe("test-topic", "test-queue")
	end := time.Now().Sub(start)
	b.Logf("Benchmark in queue %d %v %.4f cps", b.N, end, float64(b.N)/float64(end.Seconds()))
	b.ReportAllocs()
}

func TestTheSameSubscribeDifferentQueue(t *testing.T) {
	bus.queues = map[string]*queue{}
	bus.subscribers = map[string]*subscriberGroup{}

	got1 := 0
	Subscribe("test-topic", "my-test-queue1", func(event *Event) error {
		got1++
		t.Log(event.Message)
		return nil
	})

	got2 := 0
	Subscribe("test-topic", "my-test-queue2", func(event *Event) error {
		got2++
		t.Log(event.Message)
		return nil
	})

	got3 := 0
	Subscribe("test-topic", "my-test-queue3", func(event *Event) error {
		got3++
		t.Log(event.Message)
		return nil
	})

	time.Sleep(100 * time.Millisecond)

	if len(bus.subscribers) != 1 {
		t.Errorf("%+v", len(bus.subscribers))
	}

	if len(bus.queues) != 3 {
		t.Errorf("%+v", len(bus.queues))
	}

	ev1 := NewEvent("", "test-topic", "{'name':'bar'}", "data")
	Publish(ev1)

	time.Sleep(100 * time.Millisecond)

	if got1 != 1 {
		t.Errorf("%+v", got1)
	}

	if got2 != 1 {
		t.Errorf("%+v", got2)
	}

	if got3 != 1 {
		t.Errorf("%+v", got3)
	}

	Unsubscribe("test-topic", "my-test-queue1")
	time.Sleep(100 * time.Millisecond)

	if len(bus.subscribers) != 1 {
		t.Errorf("%+v", len(bus.subscribers))
	}

	if len(bus.queues) != 2 {
		t.Errorf("%+v", len(bus.queues))
	}

	got1, got2, got3 = 0, 0, 0

	ev4 := NewEvent("", "test-topic", "{'name':'foo'}", "data")
	Publish(ev4)

	ev5 := NewEvent("", "test-topic", "{'name':'foo1'}", "data")
	Publish(ev5)

	time.Sleep(100 * time.Millisecond)

	if got1 != 0 {
		t.Errorf("%+v", got1)
	}

	if got2 != 2 {
		t.Errorf("%+v", got2)
	}

	if got3 != 2 {
		t.Errorf("%+v", got3)
	}
}

// 订阅相同的topic和不同的队列
func BenchmarkTheSameSubscribeDifferentQueue(b *testing.B) {
	start := time.Now()

	for i := 0; i < b.N; i++ {
		Subscribe("test-topic", "test-queue"+strconv.Itoa(i), func(event *Event) error {
			b.Log(event.Message)
			return nil
		})
	}
	for i := 0; i < b.N; i++ {
		ev := NewEvent("", "test-topic", "{'name':'msg4'}", "data")
		Publish(ev)
	}

	for i := 0; i < b.N; i++ {
		Unsubscribe("test-topic", "test-queue"+strconv.Itoa(i))
	}

	end := time.Now().Sub(start)
	b.Logf("Benchmark in queue %d %v %.4f cps", b.N, end, float64(b.N)/float64(end.Seconds()))
	b.ReportAllocs()
}

// 订阅不同的topic和相同的queue
func BenchmarkDifferentSubscribeTheSameQueue(b *testing.B) {
	start := time.Now()

	for i := 0; i < b.N; i++ {
		Subscribe("test-topic"+strconv.Itoa(i), "test-queue", func(event *Event) error {
			b.Log(event.Message)
			return nil
		})
	}
	for i := 0; i < b.N; i++ {
		ev := NewEvent("Flag", "test-topic", "{'name':'msg3'}", "data")
		Publish(ev)
	}

	for i := 0; i < b.N; i++ {
		Unsubscribe("test-topic"+strconv.Itoa(i), "test-queue")
	}
	end := time.Now().Sub(start)
	b.Logf("Benchmark in queue %d %v %.4f cps", b.N, end, float64(b.N)/float64(end.Seconds()))
	b.ReportAllocs()
}

// 订阅不同的topic和相同的queue
func TestDifferentSubscribeTheSameQueue(t *testing.T) {
	bus.queues = map[string]*queue{}
	bus.subscribers = map[string]*subscriberGroup{}
	got1 := 0
	Subscribe("test-topic1", "test-queue", func(event *Event) error {
		got1++
		t.Log(event.Message)
		return nil
	})

	got2 := 0
	Subscribe("test-topic2", "test-queue", func(event *Event) error {
		got2++
		t.Log(event.Message)
		return nil
	})

	got3 := 0
	Subscribe("test-topic3", "test-queue", func(event *Event) error {
		got3++
		t.Log(event.Message)
		return nil
	})

	time.Sleep(time.Second)

	if len(bus.subscribers) != 3 {
		t.Errorf("%+v", len(bus.subscribers))
	}

	if len(bus.queues) != 1 {
		t.Errorf("%+v", len(bus.queues))
	}

	ev := NewEvent("", "test-topic1", "{'name':'msg1'}", "data")
	Publish(ev)

	time.Sleep(100 * time.Millisecond)

	if got1 != 1 {
		t.Errorf("%+v", got1)
	}

	if got2 != 0 {
		t.Errorf("%+v", got2)
	}

	if got3 != 0 {
		t.Errorf("%+v", got3)
	}

	ev3 := NewEvent("", "test-topic3", "{'name':'msg2'}", "data")
	Publish(ev3)

	time.Sleep(100 * time.Millisecond)

	got1, got2 = 0, 0

	if got1 != 0 {
		t.Errorf("%+v", got1)
	}

	if got2 != 0 {
		t.Errorf("%+v", got2)
	}

	if got3 != 1 {
		t.Errorf("%+v", got3)
	}

}

func TestCancelSubscribe(t *testing.T) {
	bus.queues = map[string]*queue{}
	bus.subscribers = map[string]*subscriberGroup{}
	got1 := 0
	Subscribe("test-subscribe1", "my-test-queue", func(event *Event) error {
		got1++
		t.Log(event.Message)
		return nil
	})
	got2 := 0
	Subscribe("test-subscribe1", "my-test-queue1", func(event *Event) error {
		got2++
		t.Log(event.Message)
		return nil
	})

	time.Sleep(100 * time.Millisecond)

	ev := NewEvent("", "test-subscribe1", "data1111", "data")
	Publish(ev)

	//ev1 := NewEvent("", "test-subscribe", "data2222", "data")
	//Publish(ev1)

	time.Sleep(100 * time.Millisecond)

	if len(bus.queues) != 2 {
		t.Errorf("%+v", len(bus.queues))
	}

	if len(bus.subscribers) != 1 {
		t.Errorf("%+v", len(bus.subscribers))
	}

	if got1 != 1 {
		t.Errorf("%+v", got1)
	}

	if got2 != 1 {
		t.Errorf("%+v", got2)
	}

	Unsubscribe("test-subscribe1", "my-test-queue")
	time.Sleep(100 * time.Millisecond)

	if len(bus.queues) != 1 {
		t.Errorf("%+v", len(bus.queues))
	}

	if len(bus.subscribers) != 1 {
		t.Errorf("%+v", len(bus.subscribers))
	}

	ev2 := NewEvent("", "test-subscribe1", "", "")
	Publish(ev2)

	got1, got2 = 0, 0

	if got1 != 0 {
		t.Errorf("%+v", got1)
	}

	if got2 != 0 {
		t.Errorf("%+v", got2)
	}

}

func TestNoQueue(t *testing.T) {
	bus.queues = map[string]*queue{}
	bus.subscribers = map[string]*subscriberGroup{}
	got1 := 0
	Subscribe("test-subscribe2", "", func(event *Event) error {
		got1++
		t.Log(event.Message)
		return nil
	})

	got2 := 0
	Subscribe("test-subscribe2", "", func(event *Event) error {
		got2++
		t.Log(event.Message)
		return nil
	})

	time.Sleep(100 * time.Millisecond)

	ev := NewEvent("", "test-subscribe2", "data11", "data")
	Publish(ev)

	//ev1 := NewEvent("", "test-subscribe", "data22", "data")
	//Publish(ev1)

	time.Sleep(100 * time.Millisecond)

	if len(bus.queues) != 0 {
		t.Errorf("%+v", len(bus.queues))
	}

	if len(bus.subscribers) != 1 {
		t.Errorf("%+v", len(bus.subscribers))

	}

	if got1 != 0 {
		t.Errorf("%+v", got1)
	}

	if got2 != 1 {
		t.Errorf("%+v", got2)
	}
}
