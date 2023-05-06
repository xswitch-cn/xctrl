package ds_builder

import (
	"fmt"
	"testing"
)

func TestDsBuilder(t *testing.T) {
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
	callString := DsBuilder(params)
	fmt.Println(callString)
}
