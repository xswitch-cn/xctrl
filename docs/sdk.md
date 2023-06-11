# xctrl SDK使用

- 不使用本SDK：直接根据原始的XCC API收发JSON数据即可。
- 仅使用本SDK的数据结构，使用`xctrl`包。
- 使用本SDK提供的`Service()`相关函数，可以执行同步和异步的XCC API调用。推荐使用。
- 使用自动生成的`channel`接口，可以使用比较少的参数调用函数。推荐使用。
- 使用手动补充的`channel`接口，这部分接口可能会有变化，暂不推荐使用。

下面以`Answer`和`Play`为例说明。

## `Service()`相关函数

收到`Event.Channel`的`state = START`事件后，一般都会先应答，即执行`Answer()`函数。

### Answer

对通话进行应答。

函数原型：

```go
func (c *xNodeService) Answer(ctx context.Context, in *AnswerRequest, opts ...client.CallOption) (*Response, error)
```

返回值：

- `Response`：返回值，包含`Code`和`Message`等，一般只需要判断`Code`是否为200即可。
- `error`：错误信息。

使用示例：

```go
ctrl.Service().Answer(context.TODO(), &xctrl.AnswerRequest{
    Uuid: channel.Uuid,
}, ctrl.WithTimeout(time.Second*5), ctrl.WithAddress(channel.NodeUuid()))
```

- `ctrl.Service()`：一个单例，代表当前连接的NATS服务
- `ctx`：一般用`context.TODO()`即可，只有在`AService()`调用中可以使用`context.WithCancel()`，用于取消异步调用。
- `in`：是一个`xctrl.AnswerRequest`结构体，需要包含`Uuid`，来自`ChannelEvent`的`Uuid`值。
- `opts`：可选参数，可以设置超时时间和目标节点等。
    - `ctrl.WithAddress(channel.NodeUuid())`：必填，设置目标节点为`channel.NodeUuid()`，以便NATS请求信息能发送到来源节点上。
    - `ctrl.WithTimeout(time.Second*5)`：设置超时时间为5秒。

### Play

播放媒体文件或TTS。


```go
func (c *xNodeService) Play(ctx context.Context, in *PlayRequest, opts ...client.CallOption) (*Response, error)
```

参数含义大致与`Answer`相同，只是请求时需要填充`PlayRequest`结构体。

示例

```go
media := &xctrl.Media{
    Data: file,
}
req := &xctrl.PlayRequest{
    CtrlUuid: ctrl.UUID(),
    Uuid:     channel.Uuid,
    Media:    media,
}
ctrl.Service().Play(context.Background(), req, client.WithRequestTimeout(30*time.Second), ctrl.WithAddress(channel.NodeUuid))
```

## channel接口

xctrl本身的函数参数较多，写起来比较复杂。为了简化使用，xctrl提供了`channel`接口，可以使用较少的参数调用函数。

### 自动生成部分

下面的函数也是系统根据`.proto`文件自动生成的，使用方法比较统一。这些函数实际上也是一个语法糖，写起来比较简洁，实际上还是调用了上面的对应的`Service()`相关函数。

#### Answer

函数原型：

```go
func (c *ChannelEvent) Answer(in *AnswerRequest, opts ...client.CallOption) *Response {
```

示例：

```go
channel.Answer(&xctrl.AnswerRequest{})
```

在上述示例中，只需要一个空的`xctrl.AnswerRequest`结构体即可，`NodeUuid`和`Uuid`等参数都会自动填充。

#### Play

函数原型：

```go
func (c *ChannelEvent) Play(in *PlayRequest, opts ...client.CallOption) *Response {
```

示例：

```go
media := &xctrl.Media{
    Data: file,
}
req := &xctrl.PlayRequest{
    CtrlUuid: ctrl.UUID(),
    Uuid:     channel.Uuid,
    Media:    media,
}
channel.Play(req)
```

在上述示例中，只需要填充`PlayRequest`结构体即可，`NodeUuid`和`CtrlUuid`等参数都会自动填充。后面有可选参数，比如指定不同的超时时间：

```go
channel.Play(req, ctrl.WithTimeout(10*time.Second))
```

#### 超时

在自动生成这些函数时，针对不同的函数设置了默认的超时值。如果在使用时需要不同的超时值，可以使用`ctrl.WithTimeout()`函数来设置，覆盖掉内部的默认值。如：

```go
channel.Play(req, ctrl.WithTimeout(10*time.Second))
```

默认超时值列表如下（单位为秒）：

```go
var channelMethodTimeout = map[string]int{
	"Answer":            5,
	"Accept":            5,
	"Play":              180,
	"Stop":              5,
	"Broadcast":         5,
	"Mute":              5,
	"Record":            5,
	"Hangup":            5,
	"Bridge":            3600,
	"ChannelBridge":     5,
	"Unbrdige":          5,
	"Unbrdige2":         5,
	"Hold":              5,
	"Transfer":          5,
	"ThreeWay":          3600,
	"Echo2":             3600,
	"Intercept":         3600,
	"Consult":           3600,
	"SetVar":            5,
	"GetVar":            5,
	"GetState":          5,
	"GetChannelData":    5,
	"ReadDTMF":          60,
	"ReadDigits":        60,
	"DetectSpeech":      180,
	"StopDetectSpeech":  5,
	"RingBackDetection": 60,
	"DetectFace":        180,
	"SendDTMF":          5,
	"SendInfo":          5,
	"NativeApp":         5,
	"FIFO":              1200,
	"Callcenter":        1200,
	"Conference":        3600,
	"AI":                3600,
	"HttAPI":            3600,
	"Lua":               3600,
}
```

考虑到代码可以会有变动，具体可以参见代码中的`channelMethodTimeout` map中的值。

注意，超时有可能是对方服务无法及时响应，也有可能是对方发生了故障、崩溃等。但一些阻塞的请求，如`Play`等，需要等到文件播放完成才能返回结果，所以，如果播放一个`10`秒的文件，超时时间应该大于`10`秒，如`12`秒或`15`秒。在系统负载较大时，响应的返回也可能较慢，因此最好设置一些余量。但超时间隔太大的话，会在对方服务故障时无法及时检测到。

### 手动生成部分

在1.0版本中有很多手动生成的channel相关函数，部分保留了下来。为了区分不同的参数，将原函数名称改为`0`结尾，比如`Answer0`。

```go
channel.Answer0()
channel.Play0(req)
```

这些函数参数更少，用起来更简单，参数均使用系统默认值。如果需要自定义参数，那就使用上面的自动生成的部分。

这部分函数在不同的版本可能会有变化，不推荐使用。调用者也可以参见这些函数的实现源码自行封装。
