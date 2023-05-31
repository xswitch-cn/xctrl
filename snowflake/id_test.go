package snowflake

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	err := Init()
	if err != nil {
		t.Error(err)
	}
	cases := []struct {
		function interface{}
		name     string
		wanted   string
	}{
		{New(), "New", "snowflake.ID"},
		{NewString(), "NewString", "string"},
		{NewInt64(), "NewInt64", "int64"},
		{NewBase64(), "NewBase64", "string"},
		{NewBase58(), "NewBase58", "string"},
		{NewBase36(), "NewBase36", "string"},
		{NewBase32(), "NewBase32", "string"},
	}

	for _, singleCase := range cases {
		if singleCase.wanted != typeof(singleCase.function) {
			t.Error("function:" + singleCase.name + "error")
		}
	}
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
