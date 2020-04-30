package gobootconfig

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_parseYaml(t *testing.T) {
	_clear()
	type args struct {
		read    io.Reader
		configs map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"None", args{bytes.NewBufferString(""), createMap()}, createMap()},
		{"Single", args{bytes.NewBufferString("go:\n\r  string: myValue"), createMap()}, createMap("go.string", "myValue")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseYaml(tt.args.read, tt.args.configs)
			if !reflect.DeepEqual(tt.args.configs, tt.want) {
				t.Errorf("parseYaml.getValue() = %v, want %v", tt.args.configs, tt.want)
			}
		})
	}
}

func Test_toMapRecursively(t *testing.T) {
	type args struct {
		base string
		m    interface{}
		s    map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"None", args{"", createInterfacegMap(), createMap()}, createMap()},
		{"Simple", args{"", createInterfacegMap("go", createInterfacegMap("string", "myValue")), createMap()}, createMap("go.string", "myValue")},
		{"Slice Key", args{"", createInterfacegMap("go", []interface{}{createInterfacegMap("string", "myValue")}), createMap()}, createMap("go.string", "myValue")},
		{"Slice Value", args{"", createInterfacegMap("go", createInterfacegMap("string", []interface{}{"myValue", "myValue2"})), createMap()}, createMapKeyValue("go.string", []interface{}{"myValue", "myValue2"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toMapRecursively(tt.args.base, tt.args.m, tt.args.s)
			if !reflect.DeepEqual(tt.args.s, tt.want) {
				t.Errorf("parseYaml.getValue() = %v, want %v", tt.args.s, tt.want)
			}
		})
	}
}
