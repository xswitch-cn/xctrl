package fsds

import (
	"errors"
	"fmt"
	"net"
	"path"
	"strconv"
	"strings"

	log "git.xswitch.cn/xswitch/proto/xctrl/logger"
)

type CallParams struct {
	Endpoint  string
	Profile   string
	DestNum   string
	IP        string
	Port      string
	Transport string
	Params    map[string]string
}

func (params *CallParams) String() string {
	var sb strings.Builder

	// Append the call parameters
	if len(params.Params) > 0 {
		comma := ""
		sb.WriteString("{")
		for key, value := range params.Params {
			sb.WriteString(comma)
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(value)
			comma = ","
		}
		sb.WriteString("}")
	}
	// Append the endpoint, profile, destination number, IP, port, and transport
	sb.WriteString(params.Endpoint)
	sb.WriteString("/")
	sb.WriteString(params.Profile)
	sb.WriteString("/")
	sb.WriteString(params.DestNum)
	sb.WriteString("@")
	sb.WriteString(params.IP)
	sb.WriteString(":")
	sb.WriteString(params.Port)
	if params.Transport == "udp" || params.Transport == "tcp" || params.Transport == "tls" {
		sb.WriteString(";transport=")
		sb.WriteString(params.Transport)
	} else {
		log.Warn("Transport in params is wrong or not found")
	}

	return sb.String()
}

type User struct {
	Domain string
	Number string
	IP     string
	Params map[string]string
}

func (u *User) String() string {
	var sb strings.Builder
	if len(u.Params) > 0 {
		sb.WriteString("{")
		comma := ""
		for key, value := range u.Params {
			sb.WriteString(comma)
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(quote(value))
			comma = ","
		}
		sb.WriteString("}")
		sb.WriteString("/")
	}
	s := "user/" + u.Number
	if u.Domain != "" {
		s += "@" + u.Domain
	}
	sb.WriteString(s)
	return sb.String()
}

type File struct {
	Path   string
	Name   string
	Params map[string]string
}

type PNGFile struct {
	*File
	MS        string // 图片显示时长
	Alpha     bool   // 是否支持Alpha通道
	PNGFPS    int    // png_fps:帧率，默认为5
	Text      string // 文本，可以以TTS方式播放，但需要提供下列参数
	TTSEngine string // TTS引擎
	TTSVoice  string // TTS发音人
	DText     string // 文字
	FG        string // 显示文本前景色，Web格式，如#FFFFFF
	BG        string // 显示文件背景色，Web格式，如#000000，支持透明度(Alpha Channel)，如#00000020
	Size      string // 字体大小，像素值，如24。也支持相对大小，如5vw、5vh，等，其中一个vw或vh分别为图像宽度和高度的百分之一
	ScaleW    string // 缩放图像宽度，像素
	ScaleH    string // 缩放图像高度，像素

}

func quote(str string) string {
	q := false
	if strings.Contains(str, "'") {
		str = strings.Replace(str, "'", "\\'", -1)
		q = true
	}
	if strings.Contains(str, ",") {
		// str = strings.Replace(str, ",", "\\,", -1)
		q = true
	}
	if q {
		str = "'" + str + "'"
	}
	return str
}

func (f *File) String() string {
	var sb strings.Builder

	// Append the file parameters
	if len(f.Params) > 0 {
		sb.WriteString("{")
		comma := ""
		for key, value := range f.Params {
			sb.WriteString(comma)
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(quote(value))
			comma = ","
		}
		sb.WriteString("}")
	}

	// Append the file name
	sb.WriteString(f.Path)

	return sb.String()
}

func (f PNGFile) String() (string, error) {
	if f.File.Name == "" {
		return "", errors.New("the name parameter is not set")
	}
	var sb strings.Builder
	// Append the file parameters
	sb.WriteString("{")
	if len(f.Params) > 0 {
		for key, value := range f.Params {
			writeString(&sb, key, value)
		}
	}
	if f.MS != "" {
		writeString(&sb, "png_ms", f.MS)
	}
	if f.DText != "" {
		writeString(&sb, "dtext", f.DText)
	}
	if f.PNGFPS != 0 {
		writeString(&sb, "png_fps", strconv.Itoa(f.PNGFPS))
	}
	if f.BG != "" {
		writeString(&sb, "bg", f.BG)
	}
	if f.FG != "" {
		writeString(&sb, "fg", f.FG)
	}
	if f.Text != "" {
		writeString(&sb, "text", f.Text)
	}
	if f.TTSEngine != "" {
		writeString(&sb, "tts_engine", f.TTSEngine)
	}
	if f.TTSVoice != "" {
		writeString(&sb, "tts_voice", f.TTSVoice)
	}
	if f.Alpha != false {
		writeString(&sb, "alpha", "true")
	}
	if f.Size != "" {
		writeString(&sb, "size", f.Size)
	}
	if f.ScaleW != "" {
		writeString(&sb, "scale_w", f.ScaleW)
	}
	if f.ScaleH != "" {
		writeString(&sb, "scale_h", f.ScaleH)
	}

	ParamsString := strings.TrimRight(sb.String(), ",")
	sb.Reset()
	sb.WriteString(ParamsString)
	sb.WriteString("}")
	sb.WriteString(path.Join("/", f.Path, f.Name))
	return sb.String(), nil
}

func writeString(sb *strings.Builder, key string, value string) {
	sb.WriteString(key + "=" + quote(value) + ",")
}

type Dial struct {
	LocalExtensionNum int          // 通过分机号码呼出
	IP                IPParam      // 通过IP呼出
	Gateway           GatewayParam // 通过网关呼出
	Params            map[string]string
}

type IPParam struct {
	Num       int
	IP        string
	Port      int
	Transport TransportProtocol // TCP = 1; TLS = 2
}

type TransportProtocol int

const (
	NULL TransportProtocol = iota
	TCP
	TLS
)

type GatewayParam struct {
	Realm    string
	Username string
	Password string
}

func (d Dial) String() (string, error) {
	var sb strings.Builder
	comma := ""
	for key, value := range d.Params {
		sb.WriteString(comma)
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(quote(value))
		comma = ","
	}
	if d.LocalExtensionNum != 0 {
		sb.WriteString("user/" + strconv.Itoa(d.LocalExtensionNum))
		return sb.String(), nil
	} else if d.Gateway != (GatewayParam{}) {
		// todo
		return "", errors.New("dial by gateway is not valid")
	} else if d.IP != (IPParam{}) {
		ip := net.ParseIP(d.IP.IP).To4()
		if ip == nil {
			return "", errors.New("IP is not invalid")
		}
		sb.WriteString("sofia/public/" + strconv.Itoa(d.IP.Num) + "@" + d.IP.IP)
		if d.IP.Port != 0 {
			if d.IP.Port < 0 || d.IP.Port > 65535 {
				return "", errors.New("port is not valid")
			} else {
				sb.WriteString(":" + strconv.Itoa(d.IP.Port))
			}
		}
		if d.IP.Transport != NULL {
			if d.IP.Transport == TCP {
				sb.WriteString(";transport=tcp")
			} else if d.IP.Transport == TLS {
				sb.WriteString(";transport=tls")
			} else {
				return "", errors.New("transport protocol is not valid")
			}
		}
		return sb.String(), nil
	} else {
		return "", errors.New("input nothing")
	}
}

type Agora struct {
	APPID      string
	Token      string
	Channel    string
	DestNumber string
	Params     map[string]string
}

func (agora Agora) String() (string, error) {
	var sb strings.Builder
	if len(agora.Params) > 0 {
		sb.WriteString("{")
		comma := ""
		for key, value := range agora.Params {
			sb.WriteString(comma)
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(quote(value))
			comma = ","
		}
		sb.WriteString("}")
		sb.WriteString("/")
	}
	sb.WriteString("agora")
	sb.WriteString("/")
	if agora.Token != "" {
		sb.WriteString(agora.Token)
		sb.WriteString("/")
	} else if agora.APPID != "" {
		sb.WriteString(agora.APPID)
		sb.WriteString("/")
	} else {
		return "", fmt.Errorf("agora token and agora appid is nil")
	}
	if agora.Channel != "" {
		sb.WriteString(agora.Channel)
		sb.WriteString("/")
	} else {
		return "", fmt.Errorf("agora channel is nil")
	}
	if agora.DestNumber != "" {
		sb.WriteString(agora.DestNumber)
	}
	return sb.String(), nil
}

type XRTC struct {
	VideoUseAudioIce    string
	RtpPayloadSpace     string
	AbsoluteCodecString string
	Url                 string
	Params              map[string]string
}

func (xrtc XRTC) String() (string, error) {
	var sb strings.Builder
	sb.WriteString("{")
	for key, value := range xrtc.Params {
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(quote(value))
		sb.WriteString(",")
	}
	sb.WriteString("video_use_audio_ice=")
	sb.WriteString(xrtc.VideoUseAudioIce)
	sb.WriteString(",")
	sb.WriteString("rtp_payload_space=")
	sb.WriteString(xrtc.RtpPayloadSpace)
	sb.WriteString(",")
	sb.WriteString("absolute_codec_string=")
	sb.WriteString(xrtc.AbsoluteCodecString)
	sb.WriteString(",")
	sb.WriteString("url=")
	sb.WriteString(xrtc.Url)
	sb.WriteString("}")
	return sb.String(), nil
}

type TRTC struct {
	AppId      string
	RoomId     string
	UserId     string
	UserSig    string
	DestNumber string
	Params     map[string]string
}

func (trtc TRTC) String() (string, error) {
	var sb strings.Builder
	sb.WriteString("{")
	for key, value := range trtc.Params {
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(quote(value))
		sb.WriteString(",")
	}
	sb.WriteString("trtc_user_id=")
	sb.WriteString(trtc.UserId)
	sb.WriteString(",")
	sb.WriteString("trtc_user_sig=")
	sb.WriteString(trtc.UserSig)
	sb.WriteString("}")
	sb.WriteString("/")
	sb.WriteString("trtc")
	sb.WriteString("/")
	sb.WriteString(trtc.AppId)
	sb.WriteString("/")
	sb.WriteString(trtc.RoomId)
	sb.WriteString("/")
	sb.WriteString(trtc.DestNumber)
	return sb.String(), nil
}
