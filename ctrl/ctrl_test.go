package ctrl

import (
	"os"
	"reflect"
	"testing"

	"git.xswitch.cn/xswitch/proto/xctrl/client"
)

func TestNodeAddress(t *testing.T) {
	type args struct {
		nodeUUID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				nodeUUID: "test.test",
			},
			want: "cn.xswitch.node.test.test",
		},
		{
			name: "test1",
			args: args{
				nodeUUID: "cn.xswitch.node.test.test",
			},
			want: "cn.xswitch.node.test.test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NodeAddress(tt.args.nodeUUID); got != tt.want {
				t.Errorf("NodeAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
func init() {
	natsURL := os.Getenv("NATS_ADDRESS")

	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	Init(true, natsURL)

}

func TestWithTenantAddress(t *testing.T) {
	type args struct {
		tenant   string
		nodeUUID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test with empty tenant and uuid",
			args: args{
				tenant:   "",
				nodeUUID: "",
			},
			want: "cn.xswitch.node",
		},
		{
			name: "test with empty tenant but short uuid",
			args: args{
				tenant:   "",
				nodeUUID: "test.test",
			},
			want: "cn.xswitch.node.test.test",
		},
		{
			name: "test with empty tenant and full uuid",
			args: args{
				tenant:   "",
				nodeUUID: "cn.xswitch.node.test.test",
			},
			want: "cn.xswitch.node.test.test",
		},
		{
			name: "test with tenant and empty uuid",
			args: args{
				tenant:   "foo",
				nodeUUID: "",
			},
			want: "foo.cn.xswitch.node",
		},
		{
			name: "test with tenant and short uuid",
			args: args{
				tenant:   "foo",
				nodeUUID: "test.test",
			},
			want: "foo.cn.xswitch.node.test.test",
		},
		{
			name: "test with tenant and full uuid",
			args: args{
				tenant:   "foo",
				nodeUUID: "cn.xswitch.node.test.test",
			},
			want: "foo.cn.xswitch.node.test.test",
		},
	}
	addrString := func(opt client.CallOption) string {
		opts := client.CallOptions{}
		opt(&opts)
		return opts.Address[0]
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addrString(WithTenantAddress(tt.args.tenant, tt.args.nodeUUID)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTenantAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
