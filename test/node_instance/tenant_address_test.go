package node_instance

import (
	"os"
	"testing"

	"git.xswitch.cn/xswitch/xctrl/ctrl"
)

func init() {
	natsURL := os.Getenv("NATS_ADDRESS")

	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	ctrl.Init(true, natsURL)
}

func TestTenantAddress(t *testing.T) {

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ctrl.TenantNodeAddress(tt.args.tenant, tt.args.nodeUUID); got != tt.want {
				t.Errorf("NodeAddress() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestTenantAddressWithPrefix(t *testing.T) {
	prefix := "to-"
	ctrl.SetToPrefix(prefix)
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
			want: prefix + "foo.cn.xswitch.node",
		},
		{
			name: "test with tenant and short uuid",
			args: args{
				tenant:   "foo",
				nodeUUID: "test.test",
			},
			want: prefix + "foo.cn.xswitch.node.test.test",
		},
		{
			name: "test with tenant and full uuid",
			args: args{
				tenant:   "foo",
				nodeUUID: "cn.xswitch.node.test.test",
			},
			want: prefix + "foo.cn.xswitch.node.test.test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ctrl.TenantNodeAddress(tt.args.tenant, tt.args.nodeUUID); got != tt.want {
				t.Errorf("NodeAddress() = %v, want %v", got, tt.want)
			} else {
				t.Logf("NodeAddress() = %v, want %v", got, tt.want)
			}
		})
	}

}
