package store

import (
	"encoding/json"
	"fmt"
	"time"

	"git.xswitch.cn/xswitch/xctrl/core/proto/xctrl"

	"git.xswitch.cn/xswitch/xctrl/xctrl/store"
	"git.xswitch.cn/xswitch/xctrl/xctrl/store/redis"
	"github.com/google/uuid"
)

const (
	debug = false
)

// Session .
type Session struct {
	ID        string `json:"id"`
	Domain    string `json:"domain"`
	CreatedAt string `json:"created_at"`
}

// NewSession .
func NewSession(domain string) *Session {
	s := new(Session)

	s.ID = uuid.New().String()
	s.Domain = domain
	s.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	if debug {
		fmt.Printf("%s new session created, domain=%s", s.ID, domain)
	}

	s.Write("trace", make([]*xctrl.Trace, 0))
	s.Write("domain", s.Domain)
	s.Write("created_at", s.CreatedAt)
	return s
}

// NewSession .
func NewSessionById(id string, domain string) *Session {
	s := new(Session)

	s.ID = id
	s.Domain = domain
	s.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	if debug {
		fmt.Printf("%s new session created, domain=%s", s.ID, domain)
	}

	s.Write("trace", make([]*xctrl.Trace, 0))
	s.Write("domain", s.Domain)
	s.Write("created_at", s.CreatedAt)
	return s
}

// ReadSession .
func ReadSession(id string) *Session {
	s := new(Session)
	s.ID = id

	s.Domain = s.ReadString("domain")
	s.CreatedAt = s.ReadString("created_at")
	return s
}

// key .
func (s *Session) key(key string) string {
	return Key("session", s.ID, key)
}

// Write .
func (s *Session) Write(key string, val interface{}) error {
	if debug {
		fmt.Printf("%s session write domain=%s %s=%s", s.ID, s.Domain, key, val)
	}
	return Write(s.key(key), val, store.WriteTTL(time.Hour*6))
}

// ReadInt .
func (s *Session) ReadInt(key string) int {
	result, err := ReadInt(s.key(key))
	if err != nil {
		fmt.Printf(s.ID, err)
	}
	return result
}

// ReadBool .
func (s *Session) ReadBool(key string) bool {
	result, err := ReadBool(s.key(key))
	if err != nil {
		fmt.Errorf(s.ID, err)
	}
	return result
}

// ReadString .
func (s *Session) ReadString(key string) string {
	result, err := ReadString(s.key(key))
	if err != nil {
		fmt.Errorf(s.ID, err)
	}
	return result
}

// WriteTrace .
func (s *Session) WriteTrace(t *xctrl.Trace) error {
	return s.Write("trace", append(s.ReadTrace(), t))
}

// ReadTrace .
func (s *Session) ReadTrace() []*xctrl.Trace {
	trace := make([]*xctrl.Trace, 0)
	if err := Read(s.key("trace"), &trace); err != nil {
		fmt.Errorf("err = %s", err)
	}
	return trace
}

//AppendTrace
func (s *Session) AppendTrace(t *xctrl.Trace) error {
	value, _ := json.Marshal(t)
	record := new(store.Record)
	record.Key = s.key("trace")
	record.Value = value
	if err := redis.Rediskv.Append(record, store.WriteTTL(time.Hour*6)); err != nil {
		return fmt.Errorf(`redis: write [%s] error %v`, record.Key, err)
	}
	return nil

}
