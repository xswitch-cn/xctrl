package wrapper

import (
	"context"
	"reflect"
	"testing"

	"git.xswitch.cn/xswitch/xctrl/stack/client"
	"git.xswitch.cn/xswitch/xctrl/stack/metadata"
	"git.xswitch.cn/xswitch/xctrl/stack/server"
)

func TestWrapper(t *testing.T) {
	testData := []struct {
		existing  metadata.Metadata
		headers   metadata.Metadata
		overwrite bool
	}{
		{
			existing: metadata.Metadata{},
			headers: metadata.Metadata{
				"Foo": "bar",
			},
			overwrite: true,
		},
		{
			existing: metadata.Metadata{
				"Foo": "bar",
			},
			headers: metadata.Metadata{
				"Foo": "baz",
			},
			overwrite: false,
		},
	}

	for _, d := range testData {
		c := &fromServiceWrapper{
			headers: d.headers,
		}

		ctx := metadata.NewContext(context.Background(), d.existing)
		ctx = c.setHeaders(ctx)
		md, _ := metadata.FromContext(ctx)

		for k, v := range d.headers {
			if d.overwrite && md[k] != v {
				t.Fatalf("Expected %s=%s got %s=%s", k, v, k, md[k])
			}
			if !d.overwrite && md[k] != d.existing[k] {
				t.Fatalf("Expected %s=%s got %s=%s", k, d.existing[k], k, md[k])
			}
		}
	}
}

type testRequest struct {
	service  string
	endpoint string

	server.Request
}

func (r testRequest) Service() string {
	return r.service
}

func (r testRequest) Endpoint() string {
	return r.endpoint
}

type testClient struct {
	callCount int
	callRsp   interface{}
	client.Client
}

func (c *testClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	c.callCount++

	if c.callRsp != nil {
		val := reflect.ValueOf(rsp).Elem()
		val.Set(reflect.ValueOf(c.callRsp).Elem())
	}

	return nil
}

type testRsp struct {
	value string
}
