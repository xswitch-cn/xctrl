package store

import (
	"fmt"
	"time"

	"git.xswitch.cn/xswitch/proto/xctrl/store"
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
