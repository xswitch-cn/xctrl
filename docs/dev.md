# SDK开发

本章主要针对SDK开发。

## 增加Channel函数

可以在`channel.go`中增加方便易用的常用函数，参见`PlayFile`、`PlayTTS`等。

如果在`.proto`里增加了Channel相关的函数，需要在 xctrl/cmd/protoc-gen-xctrl/plugin/xctrl.go中的`channelMethodTimeout`变量中增加对应的函数名和超时时间，以便能自动生成对应的Channel函数。

## 提PR

欢迎提PR，提PR前自己建一个账号，然后给我们联系授权即可。
