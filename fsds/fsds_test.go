package fsds

import (
	"fmt"
	"testing"
)

func TestFSDS(t *testing.T) {
	var params = &CallParams{
		Endpoint:  "sofia",
		Profile:   "public",
		DestNum:   "number",
		IP:        "ip",
		Port:      "port",
		Transport: "tcp",
		Params: map[string]string{
			"a": "1",
			"b": "2",
		},
	}
	callString := FSDS(params)
	fmt.Println(callString)
}

func TestUser(t *testing.T) {
	u := &User{
		Number: "1000",
	}
	if u.String() != "user/1000" {
		t.Errorf("user.String() = %v, want %v", u.String(), "user/1000")
	}
	u.Domain = "test.test"
	if u.String() != "user/1000@"+u.Domain {
		t.Errorf("user.String() = %v, want %v", u.String(), "user/1000@"+u.Domain)
	}
}

func TestFile(t *testing.T) {
	file_string := "{png_ms=20000,dtext=请输入会议号}/tmp/test.jpg"
	file := &File{
		Path: "/tmp/test.jpg",
		Params: map[string]string{
			"png_ms": "20000",
			"dtext":  "请输入会议号",
		},
	}
	if file.String() != file_string {
		t.Errorf("file.String() = %v, want %v", file.String(), file_string)
	}
}
