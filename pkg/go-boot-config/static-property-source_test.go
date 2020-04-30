package gobootconfig

import (
	"os"
	"reflect"
	"testing"
)

func Test_staticPropertySource_getSource(t *testing.T) {
	type fields struct {
		source string
		name   string
		value  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Simple", fields{"mySource", "go.string", "myValue"}, "mySource"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := staticPropertySource{
				source: tt.fields.source,
				name:   tt.fields.name,
				value:  tt.fields.value,
			}
			if got := sps.getSource(); got != tt.want {
				t.Errorf("staticPropertySource.getSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_staticPropertySource_getName(t *testing.T) {
	type fields struct {
		source string
		name   string
		value  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Simple", fields{"mySource", "go.string", "myValue"}, "go.string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := staticPropertySource{
				source: tt.fields.source,
				name:   tt.fields.name,
				value:  tt.fields.value,
			}
			if got := sps.getName(); got != tt.want {
				t.Errorf("staticPropertySource.getName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_staticPropertySource_getValue(t *testing.T) {
	type fields struct {
		source string
		name   string
		value  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"Simple", fields{"mySource", "go.string", "myValue"}, "myValue"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := staticPropertySource{
				source: tt.fields.source,
				name:   tt.fields.name,
				value:  tt.fields.value,
			}
			if got := sps.getValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("staticPropertySource.getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFromCmdLine(t *testing.T) {
	_clear()
	tests := []struct {
		name string
		args []string
		want interface{}
	}{
		{"None", createArr("cmd"), nil},
		{"Simple", createArr("cmd", "--go.string=myValue"), "myValue"},
	}
	for _, tt := range tests {
		os.Args = append([]string{"cmd"}, tt.args...)
		t.Run(tt.name, func(t *testing.T) {
			loadFromCmdLine()
			got, e := _getValue("go.string")
			if e == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFromCmdLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFromEnvironment(t *testing.T) {
	_clear()
	tests := []struct {
		name string
		envs map[string]string
		want interface{}
	}{
		{"None", createStringMap(), nil},
		{"Simple", createStringMap("GO_STRING", "myValue"), "myValue"},
	}
	for _, tt := range tests {
		os.Clearenv()
		for k, v := range tt.envs {
			os.Setenv(k, v)
		}
		t.Run(tt.name, func(t *testing.T) {
			loadFromEnvironment()
			got, e := _getValue("go.string")
			if e == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFromCmdLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
