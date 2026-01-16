package store

import (
	"time"

	"github.com/xswitch-cn/proto/xctrl/store"
)

// Config .
type Config struct {
	Domain string
}

// NewConfig .
func NewConfig() *Config {
	return &Config{Domain: "system"}
}

// key .
func (c *Config) key(key string) string {
	return Key(c.Domain, "setting", key)
}

// Read .
func (c *Config) Read(key string) (string, error) {
	return ReadString(c.key(key))
}

// Write .
func (c *Config) Write(key string, val string, expiry time.Duration) error {
	return Write(c.key(key), val, store.WriteTTL(expiry))
}

// Delete .
func (c *Config) Delete(key string) error {
	return Delete(c.key(key))
}
