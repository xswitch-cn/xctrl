package client

import (
	"context"
	"testing"
	"time"

	"git.xswitch.cn/xswitch/xctrl/xctrl/codec"
)

func TestBackoff(t *testing.T) {
	results := []time.Duration{
		0 * time.Second,
		10 * time.Millisecond,
		100 * time.Millisecond,
		1000 * time.Millisecond,
		10000 * time.Millisecond,
		100000 * time.Millisecond,
	}

	r := &testRequest{
		service: "test",
		method:  "test",
	}

	for i := 0; i <= 5; i++ {
		d, err := exponentialBackoff(context.TODO(), r, i)
		if err != nil {
			t.Fatal(err)
		}

		if d != results[i] {
			t.Fatalf("Expected equal than %v, got %v", results[i], d)
		}
	}
}

type testRequest struct {
	service     string
	method      string
	endpoint    string
	contentType string
	codec       codec.Codec
	body        interface{}
	opts        RequestOptions
}

func (r *testRequest) ContentType() string {
	return r.contentType
}

func (r *testRequest) Service() string {
	return r.service
}

func (r *testRequest) Method() string {
	return r.method
}

func (r *testRequest) Endpoint() string {
	return r.endpoint
}

func (r *testRequest) Body() interface{} {
	return r.body
}

func (r *testRequest) Codec() codec.Writer {
	return r.codec
}
