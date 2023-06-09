package fsds

import (
	"errors"
	"path"
	"strconv"
	"strings"

	log "git.xswitch.cn/xswitch/xctrl/xctrl/logger"
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

func FSDS(params *CallParams) string {
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
}

func (u *User) String() string {
	s := "user/" + u.Number
	if u.Domain != "" {
		s += "@" + u.Domain
	}
	return s
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
