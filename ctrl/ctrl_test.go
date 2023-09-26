package ctrl

import "testing"

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
