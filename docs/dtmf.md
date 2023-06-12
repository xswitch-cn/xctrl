# DTMF处理

## ReadDTMF

直接调用`channel.ReadDTMF`可以阻塞获取DTMF。

```go
channelReadDTMF(...)
```

这种方式的好处是它是阻塞的，获取比较简单，而且，状态机维护在XNode侧，只有匹配相应的正则表达式才会返回，推荐使用。

## 通过异步回调获取

如果自己实现DTMF状态机，则可以异步获取。

```go
type PlayHandler struct {
	*ctrl.Channel
}

func (h *PlayHandler) ChannelEvent(ctx context.Context, c *ctrl.Channel) {
}

func (h *PlayHandler) Event(msg *ctrl.Message, natsEvent nats.Event) {
	log.Println("Event: method = ", msg.Method)
	switch msg.Method {
	case "Event.DTMF":
		{
			dtmfEvent := new(xctrl.DTMFEvent)
			err := json.Unmarshal(*msg.Params, dtmfEvent)
			if err != nil {
				return
			}
            fmt.Printfln("DTMF: ", dtmfEvent.DtmfDigit, "uuid: ", dtmfEvent.Uuid)
		}
	}
}

ctrl.EnableApp(new(PlayHandler), "cn.xswitch.ctrl", "q")
```

在这种情况下，DTMF只与Uuid相关联，如果需要与当前的Channel关联，可以自己将相关数据存到在`channel.Uuid`为Key的map中。

## 将DTMF信息发到ChannelEvent协程处理

在上述示例中，通过将DTMF信息转到ChannelEvent线程处理，以便于当前的channel关联。

ChannelEvent是在单独的Go Routine中回调的，因而可以收到DTMF信息。由于ChannelEvent回调的参数是`ctrl.Channel`类型，因而，需要将DTMF事件进行转换，代码如下：

```go
func (h *PlayHandler) Event(msg *ctrl.Message, natsEvent nats.Event) {
	log.Println("Event: method = ", msg.Method)
	switch msg.Method {
	case "Event.DTMF":
		{
			dtmfEvent := new(xctrl.DTMFEvent)
			err := json.Unmarshal(*msg.Params, dtmfEvent)
			if err != nil {
				return
			}
			channel := ctrl.NewChannelEvent()
			channel.Dtmf = dtmfEvent.DtmfDigit
			channel.NodeUuid = dtmfEvent.NodeUuid
			channel.Uuid = dtmfEvent.Uuid
			channel.State = "DTMF"                 // 将DTMF事件放到State字段中，实际上Channel没有这个State
			ctrl.DeliverToChannelEventThread(channel, natsEvent)
		}
	}
}

func (h *PlayHandler) ChannelEvent(ctx context.Context, c *ctrl.Channel) {
	log.Println("App: method = Event.Channel, state = ", c.State)

    switch c.State {
	case "DTMF":
		{
			log.Println("App DTMF = ", c.Dtmf)
		}
    }
}
```

注意，在这种情况下，只有在Channel处理处理媒体状态时（如正在Play）才能收到DTMF，但同时，如果是正在阻塞的Play，则回调会阻塞。这时，如果想让`ChannelEvent`回调被执行，在上一步的Play可以放到Go Routine中执行，如：

```go
func (h *PlayHandler) ChannelEvent(ctx context.Context, c *ctrl.Channel) {
	log.Println("App: method = Event.Channel, state = ", c.State)

    switch c.State {
	case "ANSWERED":
        {
            go func() {
                channel.Play0("/tmp/welcome.wav")
            }()
        }
	case "DTMF":
		{
			log.Println("App DTMF = ", c.Dtmf)
		}
    }
}
```

## 自动转发DTMF到ChannelEvent协程

可以通过`ForkDTMFEventToChannelEventThread`函数在`xctrl` SDK中自动开启DTMF消息复制转发，它会将DTMF消息复制一份发送到`ChannelEvent`所在的线程，原来的`Event`回调也能收到DTMF。

```go
ctrl.EnableApp(new(PlayHandler), "cn.xswitch.ctrl", "q")
ctrl.ForkDTMFEventToChannelEventThread()
```

## 在收到DTMF时打断当前的播放

播放是阻塞的，通过`channel.Play()`或`ctrl.Service().Play()`函数播放时，想提前结束播放的话，如果直接通过`ctrl.WithTimeout`函数设置超时时间，只是在客户端侧（Go侧、ctrl侧）解除了函数阻塞，并不能中断XNode侧的播放。因此，要终止XNode侧的播放，可以发送`Stop()`或`break` Native App命令，如：

```go
func (h *PlayHandler) ChannelEvent(ctx context.Context, c *ctrl.Channel) {
	log.Println("App: method = Event.Channel, state = ", c.State)
	switch c.State {
	case "DTMF":
		{
			log.Println("App DTMF = ", c.Dtmf)
			if c.Dtmf == "0" {
                c.Stop() // stop playing
			}
		}
    }
}
```

或：

```go
func (h *PlayHandler) ChannelEvent(ctx context.Context, c *ctrl.Channel) {
	log.Println("App: method = Event.Channel, state = ", c.State)
	switch c.State {
	case "DTMF":
		{
			log.Println("App DTMF = ", c.Dtmf)
			if c.Dtmf == "0" {
				req := &xctrl.NativeRequest{
					Uuid: c.Uuid,
					Cmd:  "break",
					Sync: true,
				}
				c.NativeApp(req) // stop playing
			}
		}
    }
}
```

注意其中的`Sync: true`参数，该参数是必须的，它会让`break` App优先于当前的`Play` App执行，否则，`break` App会在`Play` App执行完后才执行，达不到预期效果。
