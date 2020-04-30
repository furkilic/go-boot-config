package gobootconfig

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Test_parseCmdLine(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want map[string]interface{}
	}{
		{"No args", createArr(), createMap()},
		{"Bool args", createArr("--bool"), createMap("bool", "true")},
		{"Bool With Equal args", createArr("--bool=false"), createMap("bool", "false")},
		{"Value args", createArr("--value", "myVal"), createMap("value", "myVal")},
		{"Value With Equal args", createArr("--value=myVal"), createMap("value", "myVal")},
		{"Multiple args", createArr("--val1", "--val2", "myVal", "--val3"), createMap("val1", "true", "val2", "myVal", "val3", "true")},
	}
	for _, tt := range tests {
		os.Args = append([]string{"cmd"}, tt.args...)
		t.Run(tt.name, func(t *testing.T) {
			if got := parseCmdLine(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCmdLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createArr(values ...string) []string {
	sl := make([]string, len(values))
	for i, v := range values {
		sl[i] = v
	}
	return sl
}

func createMap(keyValues ...string) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(keyValues); i += 2 {
		m[keyValues[i]] = keyValues[i+1]
	}
	return m
}
func createStringMap(keyValues ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i < len(keyValues); i += 2 {
		m[keyValues[i]] = keyValues[i+1]
	}
	return m
}
func createInterfacegMap(keyValues ...interface{}) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for i := 0; i < len(keyValues); i += 2 {
		m[keyValues[i]] = keyValues[i+1]
	}
	return m
}

func createMapKeyValue(keyValues ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(keyValues); i += 2 {
		m[fmt.Sprintf("%s", keyValues[i])] = keyValues[i+1]
	}
	return m
}
