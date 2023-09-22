package ctrl

import "testing"

func Test_findTenantId(t *testing.T) {
	type args struct {
		str        string
		fromPrefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no tenant id no prefix",
			args: args{
				str:        "cn.xswitch.ctrl",
				fromPrefix: "",
			},
			want: "",
		},
		{
			name: "no tenant id with prefix",
			args: args{
				str:        "cn.xswitch.ctrl",
				fromPrefix: "from-",
			},
			want: "",
		},
		{
			name: "has tenant id has prefix",
			args: args{
				str:        "from-foo.cn.xswitch.ctrl",
				fromPrefix: "from-",
			},
			want: "foo",
		},
		{
			name: "has tenant id no prefix",
			args: args{
				str:        "foo.cn.xswitch.ctrl",
				fromPrefix: "",
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findTenantId(tt.args.str, tt.args.fromPrefix); got != tt.want {
				t.Errorf("findTenantId() = %v, want %v", got, tt.want)
			}
		})
	}
}
