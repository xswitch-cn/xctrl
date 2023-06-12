
## bus

`bus`是一个消息总线，相当于一个内部消息队列，支持Pub/Sub模式。

```go
bus.Subscribe("topic", "queue", func(ev *Event) error {
})

ev := NewEvent("Flag", "test-topic", "message", "data")
bus.Publish(ev)
```

`Publish`用于异步地往消息队列中发送一个消息。消息会发送到一个`chan`缓冲队列中，如果队列中未消费的消息达到最大值，`Publish`操作将会被阻塞。默认的最大值为：`inboundBufferSize = 10240000`。

`Subscribe`用于订阅一个主题（`toipc`），收到消息后会回调一个回调函数。如果`queue`参数为空字符串，则回调函数会在一个新的Go Routine中回调，因此可能无法保证顺序。

如果`queue`非空，则为对于每一个订阅者而言，每一个`queue`生成一个Go Routine，所有发送到该`queue`的消息将会被顺序调用，因此应该保证`queue`的粒度，在回调函数中不要过度阻塞。

`queue`的典型应用是针对在FreeSWITCH中的一路通话，每一个Channel UUID都可以作为一个独立的`queue`进行订阅，这样，即使消息回调函数发生阻塞，也只影响这一路通话。

如果Event的`Flag`参数为`DESTROY`，则Go Routine将会终止，并自动取消订阅。

### 过期

在异常情况下，可能由于收不到`DESTROY`相关的消息，导致Go Routine无法正常终止，相关的资源也无法释放。使用`SubscribeWithExpire`，可以在极端情况下保证资源释放。需要检查回调中的`Flag`是否为`TIMEOUT`，如：

```go
bus.SubscribeWithExpire("topic", "queue", time.Hour, func(ev *Event) error {
	if ev.Flag == "TIMEOUT" {
		bus.Unsubscribe("topic", "queue")
	}
})
```

### 多次订阅相同的`topic`和相同的`queue`

在实际生产中会有很多个订阅者同时订阅相同的`topic`和相同的`queue`，多个订阅者是竞争关系，即对于一个特定的消息，有且只有一个订阅者能接收到消息。这一点跟NATS的Queue订阅类似。

### 多次订阅相同的`topic`和不同的`queue`

多个订阅者订阅相同的`topic`和不同的`queue`，对于一条特定的消息，多个订阅者都能收到。跟NATS类似。

### 多次订阅不同的`topic`和相同的`queue`

在实际生产中会有很多个订阅者订阅不同的`topic`和相同的`queue`，`queue`之间没有必然的联系，因为订阅者首先是以Topic区分的。
