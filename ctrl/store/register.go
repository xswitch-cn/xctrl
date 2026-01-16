package store

import (
	"time"

	"github.com/xswitch-cn/proto/xctrl/store"
)

// Register .
type Register struct {
	Domain    string
	Extension string
}

// NewRegister .
func NewRegister(domain string, extension string) *Register {
	return &Register{Domain: domain, Extension: extension}
}

// key .
func (r *Register) key() string {
	return Key(r.Domain, "extension", r.Extension)
}

// Read .
func (r *Register) Read() (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := read(r.key(), &data); err != nil {
		return nil, err
	}
	return data, nil
}

// Write .
func (r *Register) Write(data map[string]interface{}) error {
	if data != nil {
		return Write(r.key(), &data, store.WriteTTL(time.Minute*5))
	}
	return Write(r.key(), nil, store.WriteTTL(time.Minute*5))
}

// Delete .
func (r *Register) Delete() error {
	return Delete(r.key())
}
