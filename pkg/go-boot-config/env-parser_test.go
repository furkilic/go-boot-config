package gobootconfig

import (
	"os"
	"reflect"
	"testing"
)

func Test_parseEnv(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		want map[string]interface{}
	}{
		{"No args", createStringMap(), createMap()},
		{"Simple", createStringMap("SIMPLE", "mySimple"), createMap("simple", "mySimple")},
		{"Nested Variable", createStringMap("PARENT_CHILD", "myVal"), createMap("parent.child", "myVal")},
	}
	for _, tt := range tests {
		os.Clearenv()
		for k, v := range tt.envs {
			os.Setenv(k, v)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := parseEnv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
