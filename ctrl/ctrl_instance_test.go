package ctrl

import (
	"testing"
)

func init() {

}

func Test_getTenancyTopicAndUser(t *testing.T) {
	type args struct {
		rawTopic string
	}
	tests := []struct {
		name      string
		args      args
		wantUser  string
		wantTopic string
	}{
		// TODO: Add test cases.
		{"from-prifx-xyt", args{rawTopic: "from-xyt.cn.xswitch.test"}, "xyt", "cn.xswitch.test"},
		{"from-prifx-cherry", args{rawTopic: "from-cherry.cn.xswitch.test"}, "cherry", "cn.xswitch.test"},
		{"no-prifx", args{rawTopic: "cn.xswitch.test"}, "", "cn.xswitch.test"},
		{"no-str", args{rawTopic: ""}, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotTopic := GetTenancyTopicAndUser(tt.args.rawTopic)
			if gotUser != tt.wantUser {
				t.Errorf("tenancyTopicParser() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
			if gotTopic != tt.wantTopic {
				t.Errorf("tenancyTopicParser() gotTopic = %v, want %v", gotTopic, tt.wantTopic)
			}
		})
	}
}

func TestCtrl_GetTenancyTopicAddress(t *testing.T) {
	SetToPrefix("to")
	type args struct {
		userPrefix string
		topic      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"user prefix and to prefix is full", args{
			userPrefix: "foobar",
			topic:      "cn.xswitch.test",
		},
			"to-foobar.cn.xswitch.test",
		}, {"user prefix and to prefix is full", args{
			userPrefix: "",
			topic:      "cn.xswitch.test",
		},
			"cn.xswitch.test",
		}, {"user prefix and to prefix is full", args{
			userPrefix: "foobar",
			topic:      "",
		},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTenancyTopicAddress(tt.args.userPrefix, tt.args.topic); got != tt.want {
				t.Errorf("GetTenancyTopicAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtrl_GetTenancyTopicAndUser(t *testing.T) {
	globalCtrl = &Ctrl{}
	globalCtrl.SetFromPrefix("from-")
	type args struct {
		rawTopic string
	}
	tests := []struct {
		name      string
		args      args
		wantUser  string
		wantTopic string
	}{
		{name: "test empty", args: args{
			rawTopic: "cn.xswitch",
		}, wantUser: "", wantTopic: "cn.xswitch"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotUser, gotTopic := GetTenancyTopicAndUser(tt.args.rawTopic)
			if gotUser != tt.wantUser {
				t.Errorf("GetTenancyTopicAndUser() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
			if gotTopic != tt.wantTopic {
				t.Errorf("GetTenancyTopicAndUser() gotTopic = %v, want %v", gotTopic, tt.wantTopic)
			}
		})
	}
}
