package snowflake

import "testing"

func TestNew(t *testing.T) {
	err := Init()
	if err != nil {
		t.Error(err)
	}
	t.Log(NewString())
	t.Log("New:", New())
	t.Log("NewInt64", NewInt64())
	t.Log("NewBase64", NewBase64())
	t.Log("NewBase58", NewBase58())
	t.Log("NewBase36()", NewBase36())
	t.Log("NewBase32()", NewBase32())
	t.Error("OK")
}
