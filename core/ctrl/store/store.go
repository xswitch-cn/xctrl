package store

import (
	"encoding/json"
	"fmt"
	"strings"

	"git.xswitch.cn/xswitch/xctrl/stack/store"
	"git.xswitch.cn/xswitch/xctrl/stack/store/redis"
)

var storage store.Store

// Setting .
type Setting struct {
	Address  string `json:"address"`
	Auth     bool   `json:"auth,string"`
	Password string `json:"password"`
}

// Init .
func Init(opts ...store.Option) {
	storage = redis.NewStore(opts...)
}

// Key .
func Key(args ...string) string {
	return strings.Join(args, ".")
}

// Read .
func Read(key string, valPointer interface{}) error {
	return read(key, valPointer)
}

// ReadInt .
func ReadInt(key string) (result int, err error) {
	return result, read(key, &result)
}

// ReadBool .
func ReadBool(key string) (result bool, err error) {
	return result, read(key, &result)
}

// ReadString .
func ReadString(key string) (result string, err error) {
	return result, read(key, &result)
}

// read .
func read(key string, val interface{}) error {
	rows, err := storage.Read(key)
	if err != nil {
		return fmt.Errorf("redis: read [%s] error %v", key, err)
	}
	for _, row := range rows {
		if err := json.Unmarshal(row.Value, val); err != nil {
			return fmt.Errorf("redis: read [%s] error %v", key, err)
		}
		return nil
	}
	return nil
}

// Write .
func Write(key string, val interface{}, opts ...store.WriteOption) error {
	value, _ := json.Marshal(val)
	record := new(store.Record)
	record.Key = key
	record.Value = value
	if err := storage.Write(record, opts...); err != nil {
		return fmt.Errorf(`redis: write [%s] error %v`, key, err)
	}
	return nil
}

// Delete .
func Delete(key string) error {
	if err := storage.Delete(key); err != nil {
		return fmt.Errorf(`redis: delete [%s] error %v`, key, err)
	}
	return nil
}
