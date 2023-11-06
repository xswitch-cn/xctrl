package fsds

import (
	"errors"
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
	callString := params.String()
	fmt.Println(callString)
}

func TestUser(t *testing.T) {
	u := &User{
		Number: "1000",
		Params: map[string]string{
			"a": "a",
		},
	}
	if u.String() != "{a=a}/user/1000" {
		t.Errorf("user.String() = %v, want %v", u.String(), "user/1000")
	}
	u.Domain = "test.test"
	if u.String() != "{a=a}/user/1000@"+u.Domain {
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
				Params: map[string]string{
					"png_ms": "20000",
				},
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

func TestDialString(t *testing.T) {
	cases := []Dial{
		{
			LocalExtensionNum: 1000,
			IP:                IPParam{},
			Gateway:           GatewayParam{},
		},
		{
			LocalExtensionNum: 0,
			IP: IPParam{
				Num:       1001,
				IP:        "127.0.0.1",
				Port:      5090,
				Transport: 0,
			},
			Gateway: GatewayParam{},
		},
		{
			LocalExtensionNum: 0,
			IP: IPParam{
				Num:       1001,
				IP:        "127.0.0.1",
				Port:      5090,
				Transport: TCP,
			},
			Gateway: GatewayParam{},
		},
		{
			LocalExtensionNum: 0,
			IP: IPParam{
				Num:       1001,
				IP:        "127.0.0.1",
				Port:      0,
				Transport: TLS,
			},
			Gateway: GatewayParam{},
		},
	}

	wanted := []DialCase{
		{
			gottenString: "user/1000",
			err:          nil,
		},
		{
			gottenString: "sofia/public/1001@127.0.0.1:5090",
			err:          nil,
		},
		{
			gottenString: "sofia/public/1001@127.0.0.1:5090;transport=tcp",
			err:          nil,
		},
		{
			gottenString: "sofia/public/1001@127.0.0.1;transport=tls",
			err:          nil,
		},
	}

	for index, signalCase := range cases {

		signalCaseString, err := signalCase.String()
		if wanted[index].gottenString != signalCaseString || !errors.Is(err, wanted[index].err) {
			t.Fail()
		}
	}
}

type DialCase struct {
	gottenString string
	err          error
}

func TestAgoraString(t *testing.T) {
	agora := Agora{}
	expectedResult := ""
	expectedError := fmt.Errorf("agora token and agora appid is nil")

	result, err := agora.String()
	if result != expectedResult || err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Test case 1 failed. Expected Result: '%s', Expected Error: '%s', Actual Result: '%s', Actual Error: '%s'", expectedResult, expectedError.Error(), result, err)
	}

	agora = Agora{
		Channel:    "channel1",
		DestNumber: "number1",
	}
	expectedResult = ""
	expectedError = fmt.Errorf("agora token and agora appid is nil")

	result, err = agora.String()
	if result != expectedResult || err.Error() != expectedError.Error() {
		t.Error(err.Error(), expectedError.Error())
		t.Errorf("Test case 2 failed. Expected Result: '%s', Expected Error: '%s', Actual Result: '%s', Actual Error: '%s'", expectedResult, expectedError, result, err)
	}

	agora = Agora{
		Token:      "token1",
		DestNumber: "number1",
	}
	expectedResult = ""
	expectedError = fmt.Errorf("agora channel is nil")

	result, err = agora.String()
	if result != expectedResult || err.Error() != expectedError.Error() {
		t.Errorf("Test case 3 failed. Expected Result: '%s', Expected Error: '%s', Actual Result: '%s', Actual Error: '%s'", expectedResult, expectedError, result, err)
	}

	agora = Agora{
		APPID:      "appid1",
		DestNumber: "number1",
	}
	expectedResult = ""
	expectedError = fmt.Errorf("agora channel is nil")

	result, err = agora.String()
	if result != expectedResult || err.Error() != expectedError.Error() {
		t.Errorf("Test case 4 failed. Expected Result: '%s', Expected Error: '%s', Actual Result: '%s', Actual Error: '%s'", expectedResult, expectedError, result, err)
	}

	agora = Agora{
		Token:      "token1",
		Channel:    "channel1",
		DestNumber: "number1",
		Params:     map[string]string{"a": "a"},
	}
	expectedResult = "{a=a}/agora/token1/channel1/number1"
	expectedError = nil

	result, err = agora.String()
	if result != expectedResult || err != expectedError {
		t.Errorf("Test case 5 failed. Expected Result: '%s', Expected Error: '%s', Actual Result: '%s', Actual Error: '%s'", expectedResult, expectedError, result, err)
	}
}

func TestXRTC(t *testing.T) {
	xrtc := XRTC{
		VideoUseAudioIce:    "video_use_audio_ice_value",
		RtpPayloadSpace:     "rtp_payload_space_value",
		AbsoluteCodecString: "absolute_codec_string_value",
		Url:                 "url_value",
		Params:              map[string]string{"a": "a"},
	}

	expected := "{a=a,video_use_audio_ice=video_use_audio_ice_value,rtp_payload_space=rtp_payload_space_value,absolute_codec_string=absolute_codec_string_value,url=url_value}"
	result, err := xrtc.String()

	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestTRTC(t *testing.T) {
	trtc := TRTC{
		UserId:     "trtc_user_id_value",
		UserSig:    "trtc_user_sig_value",
		AppId:      "trtc_app_id_value",
		RoomId:     "trtc_room_id_value",
		DestNumber: "trtc_dest_number_value",
		Params:     map[string]string{"a": "a"},
	}

	expected := "{a=a,trtc_user_id=trtc_user_id_value,trtc_user_sig=trtc_user_sig_value}/trtc/trtc_app_id_value/trtc_room_id_value/trtc_dest_number_value"
	result, err := trtc.String()

	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
