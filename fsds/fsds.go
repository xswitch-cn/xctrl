package fsds

import (
	"errors"
	"path"
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
	MS    string //图片显示时长
	DText string //文字
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
	comma := ""
	if len(f.Params) > 0 {
		for key, value := range f.Params {
			sb.WriteString(comma)
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(quote(value))
			comma = ","
		}
	}
	if f.MS != "" {
		sb.WriteString(comma)
		if comma == "" {
			comma = ","
		}
		sb.WriteString("png_ms")
		sb.WriteString("=")
		sb.WriteString(f.MS)
	}
	if f.DText != "" {
		sb.WriteString(comma)
		sb.WriteString("dtext")
		sb.WriteString("=")
		sb.WriteString(quote(f.DText))
	}
	sb.WriteString("}")
	sb.WriteString(path.Join("/", f.Path, f.Name))
	return sb.String(), nil
}
