package nats

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
)

var addrTestCases = []struct {
	name        string
	description string
	addrs       map[string]string // expected address : set address
}{
	{
		"connOpts",
		"set broker addresses through a Option in constructor",
		map[string]string{
			"nats://192.168.10.1:5222": "192.168.10.1:5222",
			"nats://10.20.10.0:4222":   "10.20.10.0:4222"},
	},
	{
		"connInit",
		"set broker addresses through a Option in broker.Init()",
		map[string]string{
			"nats://192.168.10.1:5222": "192.168.10.1:5222",
			"nats://10.20.10.0:4222":   "10.20.10.0:4222"},
	},
	{
		"default",
		"check if default Address is set correctly",
		map[string]string{
			"nats://127.0.0.1:4222": "",
		},
	},
}

// TestInitAddrs tests issue #100. Ensures that if the addrs is set by an option in init it will be used.
func TestInitAddrs(t *testing.T) {

	for _, tc := range addrTestCases {
		t.Run(fmt.Sprintf("%s: %s", tc.name, tc.description), func(t *testing.T) {

			var c Conn
			var addrs []string

			for _, addr := range tc.addrs {
				addrs = append(addrs, addr)
			}

			switch tc.name {
			case "connOpts":
				// we know that there are just two addrs in the dict
				c = NewConn(Addrs(addrs[0], addrs[1]))
				c.Init()
			case "connInit":
				c = NewConn()
				// we know that there are just two addrs in the dict
				c.Init(Addrs(addrs[0], addrs[1]))
			case "default":
				c = NewConn()
				c.Init()
			}

			natsBroker, ok := c.(*nConn)
			if !ok {
				t.Fatal("Expected broker to be of types *nbroker")
			}
			// check if the same amount of addrs we set has actually been set, default
			// have only 1 address nats://127.0.0.1:4222 (current nats code) or
			// nats://localhost:4222 (older code version)
			if len(natsBroker.addrs) != len(tc.addrs) && tc.name != "default" {
				t.Errorf("Expected Addr count = %d, Actual Addr count = %d",
					len(natsBroker.addrs), len(tc.addrs))
			}

			for _, addr := range natsBroker.addrs {
				_, ok := tc.addrs[addr]
				if !ok {
					t.Errorf("Expected '%s' has not been set", addr)
				}
			}
		})

	}
}

func TestRequest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	url := os.Getenv("NATS_ADDRESS")
	if url == "" {
		url = nats.DefaultURL
	}
	nc, err := nats.Connect(url)

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	// Request with context
	msg, err := nc.RequestWithContext(ctx, "foo", []byte("bar"))

	if err != nil {
		if err.Error() != "context canceled" {
			t.Error(err)
		}
	}

	fmt.Print(msg)
}
