package fsds

import (
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
	MS    string
	DText string
}

func (f *File) String() string {
	var sb strings.Builder

	// Append the file parameters
	if len(f.Params) > 0 {
		sb.WriteString("{")
		for key, value := range f.Params {
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(value)
			sb.WriteString(",")
		}
		// Remove the trailing comma
		str := sb.String()
		if len(str) > 0 {
			sb.Reset()
			sb.WriteString(str[:len(str)-1])
		}
		sb.WriteString("}")
	}

	// Append the file name
	sb.WriteString(f.Path)

	return sb.String()
}
