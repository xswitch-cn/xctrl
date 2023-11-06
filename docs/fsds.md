# FSDS字符串生成

FSDS指FreeSWITCH Dial String，亦指File String and Dial String，用于格式化生成文件字符串和叫叫字符串。

使用方法：

```go
import "git.xswitch.cn/xswitch/xctrl/fsds"
var params = &fsds.CallParams{
	...
}
callString := params.String()
```

## FSDS

FSDS结构体输入示例。

```go
var FSDS = &fsds.FSDS{
	Endpoint:  "sofia",
	Profile:   "public",
	DestNum:   "number",
	IP:        "ip",
	Port:      "port",
	Transport: "tcp",
	Params: map[string]string{
		"a": "1",
		"b": "2",
	},
}
callString := FSDS.String()
```

FSDS输出。

```sh
{a=1,b=2}sofia/public/number@ip:port;transport=tcp
```

## Endpoint接口

### user

user结构体输入。

```go
u := &fsds.User{
	Number: "1000",
	Domain: "test.test",//可选参数
	Params: map[string]string{
		"a": "1",
		"b": "2",
	},
}
callString := u.String()
```

user输出。

```sh
{a=1,b=2}/user/1000@test.test
```

user结构体输入，domain为空时。

```go
u := &User{
	Number: "1000",
	Params: map[string]string{
		"a": "1",
		"b": "2",
	},
}
callString := u.String()
```

user输出。

```shell
user/1000
```


### dial

dial结构体输入。

```go
Dial := fsds.Dial {
	LocalExtensionNum: 0,
	IP: IPParam{
		Num:       1001,
		IP:        "127.0.0.1",
		Port:      5090,//可选参数
		Transport: 0,
	},
	Gateway: GatewayParam{},
	Params: map[string]string{
		"a": "1",
		"b": "2",
	},
}
signalCaseString, err := Dial.String()
```

Transport：如果使用TCP或TLS链路：

- 0：空
- 1：TCP
- 2：TLS

dial输出。

```
{A=1,B=2}/sofia/public/1001@127.0.0.1:5090
```

### Agora

Agora结构体输入。

```go
agora := fsds.Agora{
	AppID:      "appid1",
	Token:      "token1",
	Channel:    "channel1",
	DestNumber: "number1",
	Params: map[string]string{
		"a": "1",
		"b": "2",
	},
}
expectedResult = "agora/token1/channel1/number1"
result, err = agora.String()
```

- AppID：Agora App ID，字符串，同一个ID下的媒体才有可能互通。
- Token：App ID泄漏后后果严重，因此使用Token来代替App ID。
- Channel：Agora Channel（频道），字符串，加入同一个Channel中的所有成员都可以媒体互通。

Agora输出。

```sh
{a=1,b=2}/agora/token1/channel1/number1
```

### XRTC

xrtc结构体输入。

```go
xrtc := fsds.XRTC{
	VideoUseAudioIce:    "video_use_audio_ice_value",
	RtpPayloadSpace:     "rtp_payload_space_value",
	AbsoluteCodecString: "absolute_codec_string_value",
	Url:                 "url_value",
		Params: map[string]string{
		"a": "1",
		"b": "2",
	},
}
result, err := xrtc.String()
```

- video_use_audio_ice：在Bundle模式，视频使用音频的ICE绑定信息。
- rtp_payload_space：指定视频的RTP载荷值，不指定有时候会有冲突。
- absolute_codec_string：音视频编码，一般为OPUS,H264，在呼叫字符串中使用时其中的逗号要使用\转义。
url：SRS 推拉流 URL。

XRTC输出。

```
{a=1,b=2,video_use_audio_ice=video_use_audio_ice_value,rtp_payload_space=rtp_payload_space_value,absolute_codec_string=absolute_codec_string_value,url=url_value}

```


### TRTC

trtc结构体输入。

```go
trtc := fsds.TRTC{
	UserId:     "trtc_user_id_value",
	UserSig:    "trtc_user_sig_value",
	AppId:      "trtc_app_id_value",
	RoomId:     "trtc_room_id_value",
	DestNumber: "trtc_dest_number_value",
		Params: map[string]string{
		"a": "1",
		"b": "2",
	},
	
}
result, err := trtc.String()
```

- UserId：SDK App ID，字符串，不同AppID之间的数据不互通
- RoomId：房间ID，字符串，加入同一个room的成员可以媒体互通
- UserId：用户ID，字符串，TRTC不支持同一个UserID （除非 SDKAppID 不同）在两个设备同时使用
- UserSig：动态签名，详见 https://github.com/tencentyun/tls-sig-api下面的README.md

TRTC输出。

```sh
{a=1,b=2,trtc_user_id=trtc_user_id_value,trtc_user_sig=trtc_user_sig_value}/trtc/trtc_app_id_value/trtc_room_id_value/trtc_dest_number_value"
```

## File接口

###  通过File接口

File结构体输入。

```go
file := &fsds.File{
	Path: "/tmp/test.jpg", //支持 .mp4 .vv .png .jpg文件
	Params: map[string]string{
		"png_ms": "20000",
		"dtext":  "请输入会议号",
	},
}
fileString := file.String()
```

file输出。

```shell
{png_ms=20000,dtext=请输入会议号}/tmp/test.jpg
```


## PNGFile

PNGFile输入。

```go
pngFile := fsds.PNGFile{
	File: &File{
		Path: "tmp/",
		Name: "test.png",
		Params: map[string]string{
		"a": "1",
		"b": "2",
	},
	MS:        "1000",//可选参数
	Alpha:     true,//可选参数
	PNGFPS:    10,//可选参数
	Text:      "text",//可选参数
	TTSEngine: "ttsengine",//可选参数
	TTSVoice:  "ttsvoice",//可选参数
	DText:     "dtext",//可选参数
	FG:        "fg",//可选参数
	BG:        "bg",//可选参数
	Size:      "size",//可选参数
	ScaleW:    "scalew",//可选参数
	ScaleH:    "scaleh",//可选参数
}
pngFileString, _ := pngFile.String()
```

- MS：时长，毫秒
- PNGFPS：帧率，默认为5
- Alpha：是否支持Alpha通道
- Text：文本，可以以TTS方式播放，但需要提供下列参数
- TTSEngine：TTS引擎
- TTSVoice：TTS发音人
- DText：显示文本
- FG：显示文本前景色，Web格式，如#FFFFFF。
- BG：显示文件背景色，Web格式，如#000000，支持透明度（Alpha Channel），如#00000020。
- Size：字体大小，像素值，如24。也支持相对大小，如5vw、5vh，等，其中一个vw或vh分别为图像宽度和高度的百分之一。
- ScaleW：缩放图像宽度，像素
- ScaleH：缩放图像高度，像素

pngfile输出。

```sh
{png_ms=1000,dtext=dtext,png_fps=10,bg=bg,fg=fg,text=text,tts_engine=ttsengine,tts_voice=ttsvoice,alpha=true,size=size,scale_w=scalew,scale_h=scaleh}/tmp/test.png
```

详细参数可参见：<https://git.xswitch.cn/xswitch/vip/src/branch/main/docs/freeswitch/png.md>，仅限VIP用户访问。

### VVFile

Todo.

## AV

Todo.
