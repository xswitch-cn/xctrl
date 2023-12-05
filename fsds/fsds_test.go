package fsds

import (
	"testing"
)

func TestEndpointString(t *testing.T) {
	endpoint := Endpoint{
		FSDS: &FSDS{
			Params:         map[string]string{"param1": "value1", "param2": "value2"},
			CallerIDName:   "test1",
			CallerIDNumber: "test2",
		},
		Type: "test",
		Dest: "1234",
	}
	expectedResult := "{param1=value1,param2=value2,caller_id_name=test1,caller_id_number=test2}test/1234"

	result := endpoint.String()

	if result != expectedResult {
		t.Errorf("Endpoint.String() = %s; want %s", result, expectedResult)
	}
}

func TestIPEndpointString(t *testing.T) {
	ipEndpoint := IP{
		Endpoint: &Endpoint{
			FSDS: &FSDS{
				Params: map[string]string{},
			},
			Type:    "test",
			Dest:    "1234",
			Profile: "public",
		},
		IP:        "127.0.0.1",
		Port:      "5060",
		Transport: "tcp",
	}
	expectedResult := "test/public/1234@127.0.0.1:5060;transport=tcp"

	result := ipEndpoint.String()

	if result != expectedResult {
		t.Errorf("IPEndpoint.String() = %s; want %s", result, expectedResult)
	}
}

func TestGatewayParamString(t *testing.T) {
	gatewayEndpoint := Gateway{
		Endpoint: &Endpoint{
			FSDS: &FSDS{
				Params:       map[string]string{"param1": "value1", "param2": "value2"},
				CallerIDName: "test1",
			},
			Type:    "test",
			Profile: "gateway",
			Dest:    "1234",
		},
		GatewayName: "gatewayname",
	}
	expectedResult := "{param1=value1,param2=value2,caller_id_name=test1}test/gateway/gatewayname/1234"

	result := gatewayEndpoint.String()

	if result != expectedResult {
		t.Errorf("gatewayEndpoint.String() = %s; want %s", result, expectedResult)
	}
}

func TestUserEndpoint_String(t *testing.T) {
	endpoint := User{
		Endpoint: &Endpoint{
			FSDS: &FSDS{
				Params:         map[string]string{"param1": "value1", "param2": "value2"},
				CallerIDName:   "test1",
				CallerIDNumber: "test2",
			},
			Type: "test",
			Dest: "1234",
		},
		Domain: "domain",
	}

	expected := "{param1=value1,param2=value2,caller_id_name=test1,caller_id_number=test2}test/1234@domain"
	result := endpoint.String()

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestFile(t *testing.T) {
	file_string := "{png_ms=20000,dtext=请输入会议号}tmp/test.jpg"
	file := &File{
		FSDS: &FSDS{
			Params: map[string]string{
				"png_ms": "20000",
				"dtext":  "请输入会议号",
			},
		},
		Path: "tmp/test.jpg",
	}
	if file.String() != file_string {
		t.Errorf("file.String() = %v, want %v", file.String(), file_string)
	}
}

func TestPNGFile(t *testing.T) {
	Cases := []PNGFile{
		{
			File: &File{
				FSDS: &FSDS{
					Params: nil,
				},
				Path: "/tmp",
				Name: "test.jpg",
			},
			MS:    "20000",
			DText: "请输入会议号",
		},
		{
			File: &File{
				Path: "tmp",
				Name: "test.jpg",
				FSDS: &FSDS{
					Params: nil,
				},
			},
			MS:    "20000",
			DText: "请输入会议号",
		},
		{
			File: &File{
				Path: "tmp/",
				Name: "/test.jpg",
				FSDS: &FSDS{
					Params: map[string]string{
						"png_ms": "20000",
					},
				},
			},
			DText: "请输入会议号",
		},
		{
			File: &File{
				Path: "tmp/",
				Name: "",
				FSDS: &FSDS{
					Params: map[string]string{
						"png_ms": "20000",
					},
				},
			},
			DText: "请输入会议号",
		},
		{
			File: &File{
				Path: "tmp",
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
			"{png_ms=20000,dtext=请输入会议号}tmp/test.jpg", false,
		},
		{
			"{png_ms=20000,dtext=请输入会议号}tmp/test.jpg", false,
		},
		{
			"", true,
		},
		{
			"{png_ms=1000,dtext=dtext,png_fps=10,bg=bg,fg=fg,text=text,tts_engine=ttsengine,tts_voice=ttsvoice,alpha=true,size=size,scale_w=scalew,scale_h=scaleh}tmp/test.png", false,
		},
	}
	for index, signalCase := range Cases {
		stringGot, err := signalCase.String()
		if !(err == nil) != wanted[index].gotError || stringGot != wanted[index].gotString {
			t.Error("index:", index, "error", stringGot, wanted[index].gotString)
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

func TestAgoraString(t *testing.T) {
	agora := Agora{
		Endpoint: &Endpoint{
			FSDS: &FSDS{
				Params: map[string]string{"a": "a"},
			},
			Type: "agora",
			Dest: "number1",
		},
		Token:   "token1",
		Channel: "channel1",
	}
	expectedResult := "{a=a}agora/token1/channel1/number1"

	result, _ := agora.String()
	if result != expectedResult {
		t.Errorf("Test case 5 failed. Expected Result: '%s', Actual Result: '%s'", expectedResult, result)
	}
}

func TestXRTC(t *testing.T) {
	xrtc := XRTC{
		Endpoint: &Endpoint{
			FSDS: &FSDS{
				Params:       map[string]string{"a": "a"},
				CallerIDName: "test",
			},
		},
		VideoUseAudioIce:    "video_use_audio_ice_value",
		RtpPayloadSpace:     "rtp_payload_space_value",
		AbsoluteCodecString: "absolute_codec_string_value",
		Url:                 "url_value",
	}

	expected := "{a=a,caller_id_name=test,video_use_audio_ice=video_use_audio_ice_value,rtp_payload_space=rtp_payload_space_value,absolute_codec_string=absolute_codec_string_value,url=url_value}"
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
		Endpoint: &Endpoint{
			FSDS: &FSDS{
				Params: map[string]string{
					"a": "a",
				},
				CallerIDName:   "",
				CallerIDNumber: "",
			},
			Type: "trtc",
			Dest: "trtc_dest_number_value",
		},
		UserId:  "trtc_user_id_value",
		UserSig: "trtc_user_sig_value",
		AppId:   "trtc_app_id_value",
		RoomId:  "trtc_room_id_value",
	}

	expected := "{a=a,trtc_user_id=trtc_user_id_value,trtc_user_sig=trtc_user_sig_value}trtc/trtc_app_id_value/trtc_room_id_value/trtc_dest_number_value"
	result, err := trtc.String()

	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestString(t *testing.T) {
	vvFile := VVFile{
		Endpoint: &Endpoint{
			FSDS: &FSDS{
				Params: map[string]string{
					"param": "test",
				},
			},
		},
		VVMs:   "some-value",
		Engine: "some-engine",
		Voice:  "some-voice",
		Text:   "some-text",
	}

	expectedResult := "{param=test,vv_ms=some-value}vv://tts://some-engine|some-voice|some-text" // 设置预期结果

	actualResult, err := vvFile.String()

	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	if actualResult != expectedResult {
		t.Errorf("Expected %s, but got %s", expectedResult, actualResult)
	}
}
