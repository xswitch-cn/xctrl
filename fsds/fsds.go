package fsds

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
)

const (
	TCP string = "tcp"
	TLS string = "tls"
)

type FSDS struct {
	Params         map[string]string
	CallerIDName   string
	CallerIDNumber string
}

type Endpoint struct {
	*FSDS
	Type    string
	Profile string
	Dest    string
}

type IP struct {
	*Endpoint
	IP        string
	Port      string
	Transport string
}

type Gateway struct {
	*Endpoint
	GatewayName string
}

type User struct {
	*Endpoint
	Domain string
}

func (fsds FSDS) String() string {
	var sb strings.Builder
	if len(fsds.Params) > 0 || fsds.CallerIDName != "" || fsds.CallerIDNumber != "" {
		sb.WriteString("{")
		comma := ""
		for key, value := range fsds.Params {
			sb.WriteString(comma)
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(quote(value))
			comma = ","
		}
		if fsds.CallerIDName != "" {
			sb.WriteString(comma)
			sb.WriteString("caller_id_name")
			sb.WriteString("=")
			sb.WriteString(quote(fsds.CallerIDName))
			comma = ","
		}
		if fsds.CallerIDNumber != "" {
			sb.WriteString(comma)
			sb.WriteString("caller_id_number")
			sb.WriteString("=")
			sb.WriteString(quote(fsds.CallerIDNumber))
			comma = ","
		}
		sb.WriteString("}")
	} else {
		return ""
	}
	return sb.String()
}

func (e Endpoint) String() string {
	var sb strings.Builder
	if e.FSDS != nil {
		sb.WriteString(e.FSDS.String())
	}
	sb.WriteString(e.Type)
	sb.WriteString("/")
	if e.Profile != "" {
		sb.WriteString(e.Profile)
		sb.WriteString("/")
	}
	sb.WriteString(e.Dest)
	return sb.String()
}

func (e IP) String() string {
	var sb strings.Builder
	sb.WriteString(e.Endpoint.String())
	if e.IP != "" {
		sb.WriteString("@")
		sb.WriteString(e.IP)
		if e.Port != "" {
			sb.WriteString(":")
			sb.WriteString(e.Port)
		}
	}
	if e.Transport != "" {
		sb.WriteString(";")
		sb.WriteString("transport=")
		sb.WriteString(e.Transport)
	}
	return sb.String()
}

func (e Gateway) String() string {
	var sb strings.Builder
	if e.FSDS != nil {
		sb.WriteString(e.FSDS.String())
	}
	sb.WriteString(e.Endpoint.Type)
	sb.WriteString("/")
	sb.WriteString(e.Endpoint.Profile)
	sb.WriteString("/")
	sb.WriteString(e.GatewayName)
	sb.WriteString("/")
	sb.WriteString(e.Endpoint.Dest)
	return sb.String()
}

func (e User) String() string {
	var sb strings.Builder
	sb.WriteString(e.Endpoint.String())
	if e.Domain != "" {
		sb.WriteString("@")
		sb.WriteString(e.Domain)
	}
	return sb.String()
}

type File struct {
	*FSDS
	Path string
	Name string
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
	if f.FSDS != nil {
		sb.WriteString(f.FSDS.String())
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
	if f.FSDS != nil && f.FSDS.String() != "" {
		sb.WriteString(f.FSDS.String()[:len(f.FSDS.String())-1])
		sb.WriteString(",")
	} else {
		sb.WriteString("{")
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

type Agora struct {
	*Endpoint
	APPID   string
	Token   string
	Channel string
}

func (agora Agora) String() (string, error) {
	var sb strings.Builder
	if agora.FSDS != nil {
		sb.WriteString(agora.FSDS.String())
	}
	if agora.Endpoint.Type != "" {
		sb.WriteString(agora.Endpoint.Type)
		sb.WriteString("/")
	}
	if agora.Endpoint.Profile != "" {
		sb.WriteString(agora.Endpoint.Profile)
		sb.WriteString("/")
	}
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
	if agora.Endpoint.Dest != "" {
		sb.WriteString(agora.Endpoint.Dest)
	}
	return sb.String(), nil
}

type XRTC struct {
	*Endpoint
	VideoUseAudioIce    string
	RtpPayloadSpace     string
	AbsoluteCodecString string
	Url                 string
}

func (xrtc XRTC) String() (string, error) {
	var sb strings.Builder
	if xrtc.FSDS != nil && xrtc.FSDS.String() != "" {
		sb.WriteString(xrtc.FSDS.String()[:len(xrtc.FSDS.String())-1])
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
	if xrtc.Endpoint.Type != "" {
		sb.WriteString(xrtc.Endpoint.Type)
		sb.WriteString("/")
	}
	if xrtc.Endpoint.Profile != "" {
		sb.WriteString(xrtc.Endpoint.Profile)
		sb.WriteString("/")
	}
	return sb.String(), nil
}

type TRTC struct {
	*Endpoint
	AppId   string
	RoomId  string
	UserId  string
	UserSig string
}

func (trtc TRTC) String() (string, error) {
	var sb strings.Builder
	if trtc.FSDS != nil && trtc.FSDS.String() != "" {
		sb.WriteString(trtc.FSDS.String()[:len(trtc.FSDS.String())-1])
		sb.WriteString(",")
	}
	sb.WriteString("trtc_user_id=")
	sb.WriteString(trtc.UserId)
	sb.WriteString(",")
	sb.WriteString("trtc_user_sig=")
	sb.WriteString(trtc.UserSig)
	sb.WriteString("}")
	sb.WriteString(trtc.Endpoint.Type)
	sb.WriteString("/")
	sb.WriteString(trtc.AppId)
	sb.WriteString("/")
	sb.WriteString(trtc.RoomId)
	sb.WriteString("/")
	sb.WriteString(trtc.Endpoint.Dest)
	return sb.String(), nil
}
