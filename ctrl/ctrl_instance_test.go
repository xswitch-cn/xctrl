package ctrl

import "testing"

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
