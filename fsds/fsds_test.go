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

func TestPNGFile(t *testing.T) {
	Cases := []PNGFile{
		{
			File: &File{
				Path:   "tmp",
				Name:   "test.jpg",
				Params: nil,
			},
			MS:    "20000",
			DText: "请输入会议号",
		},
		{
			File: &File{
				Path:   "/tmp",
				Name:   "/test.jpg",
				Params: nil,
			},
			MS:    "20000",
			DText: "请输入会议号",
		},
		{
			File: &File{
				Path: "tmp/",
				Name: "/test.jpg",
				Params: map[string]string{
					"png_ms": "20000",
				},
			},
			DText: "请输入会议号",
		},
		{
			File: &File{
				Path: "tmp/",
				Name: "",
			},
			DText: "请输入会议号",
		},
		{
			File: &File{
				Path: "tmp/",
				Name: "test.png",
			},
			MS:        "1000",
			Alpha:     true,
			PNGFPS:    10,
			Text:      "text",
			TTSEngine: "ttsengine",
			TTSVoice:  "ttsvoice",
			DText:     "dtext",
			FG:        "fg",
			BG:        "bg",
			Size:      "size",
			ScaleW:    "scalew",
			ScaleH:    "scaleh",
		},
	}
	wanted := []PNGFileCase{
		{
			"{png_ms=20000,dtext=请输入会议号}/tmp/test.jpg", false,
		},
		{
			"{png_ms=20000,dtext=请输入会议号}/tmp/test.jpg", false,
		},
		{
			"{png_ms=20000,dtext=请输入会议号}/tmp/test.jpg", false,
		},
		{
			"", true,
		},
		{
			"{png_ms=1000,dtext=dtext,png_fps=10,bg=bg,fg=fg,text=text,tts_engine=ttsengine,tts_voice=ttsvoice,alpha=true,size=size,scale_w=scalew,scale_h=scaleh}/tmp/test.png", false,
		},
	}
	for index, signalCase := range Cases {
		stringGot, err := signalCase.String()
		if !(err == nil) != wanted[index].gotError || stringGot != wanted[index].gotString {
			t.Error("index:", index, "error")
		}
	}
}

type PNGFileCase struct {
	gotString string
	gotError  bool
}

func TestQuote(t *testing.T) {
	Cases := []string{
		"'test'",
		"'test'demo",
		"'test',demo,example",
	}

	wanted := [][]string{
		{"'", "\\", "'", "t", "e", "s", "t", "\\", "'", "'"},
		{"'", "\\", "'", "t", "e", "s", "t", "\\", "'", "d", "e", "m", "o", "'"},
		{"'", "\\", "'", "t", "e", "s", "t", "\\", "'", ",", "d", "e", "m", "o", ",", "e", "x", "a", "m", "p", "l", "e", "'"},
	}
	for index, signalCase := range Cases {
		for j, char := range []rune(quote(signalCase)) {
			if string(char) != wanted[index][j] {
				t.Fail()
			}
		}
	}
}
