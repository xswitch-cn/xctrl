package bus

import (
	"fmt"
	"runtime"
	"sync"
	"time"

)

// queueEvent is given to a subscription handler for processing
type queueEvent struct {
	handler Handler
	ev      *Event
}

type queue struct {
	name     string
	members  map[string]chan *Event
	refCount int
	inbound  chan *Event
	done     chan bool
	lock     sync.Mutex
	expire   time.Duration
}

func newQueue(name string, expire time.Duration) *queue {
	q := &queue{
		inbound:  make(chan *Event, queueBufferSize),
		done:     make(chan bool, 1),
		name:     name,
		refCount: 1,
		expire:   expire,
	}
	q.start()
	return q
}

func (q *queue) runWithRecovery(e *Event) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Errorf("cron: panic running job: %v\n%s", r, string(buf))
		}
	}()
	if e.handler != nil {
		e.handler(e)
	}
}

func (q *queue) addRef() {
	q.lock.Lock()
	q.refCount = q.refCount + 1
	q.lock.Unlock()
}

func (q *queue) release() int {
	q.lock.Lock()
	q.refCount = q.refCount - 1
	if q.refCount < 1 {
		q.done <- true
	}
	q.lock.Unlock()
	return q.refCount
}

func (q *queue) start() {
	go func() {
		//@TODO, this may cause cli-ecc crash, we disable this now, need todo
		// defer func() {
		// 	close(q.inbound)
		// 	close(q.done)
		// }()
		running := true
		fmt.Printf("Queue %s started", q.name)
		if q.expire > 0 {
			var handler Handler
			for running {
				// xlog.Infof("Queue %s running", q.name)
				select {
				case e, ok := <-q.inbound:
					if !ok {
						fmt.Errorf("error read from inbound chan")
						continue
					}
					fmt.Printf("%s delivered to handler", e.Topic)
					// cache the last Handler
					handler = e.handler
					q.runWithRecovery(e)
				case <-q.done:
					running = false
				case <-time.After(q.expire):
					// sigh, we don't have a handler here ?
					q.runWithRecovery(&Event{Flag: "TIMEOUT", handler: handler})
					running = false
					fmt.Printf("Queue %s timeout %d", q.name, q.expire)
				}
			}
		} else {
			for running {
				// xlog.Infof("Queue %s running", q.name)
				select {
				case e, ok := <-q.inbound:
					//xlog.Debugf("queue inbound %v\n", e.Queue)
					if ok {
						q.runWithRecovery(e)
					}
				case <-q.done:
					running = false
				}
			}

		}
		fmt.Printf("Queue %s done", q.name)
	}()
}
