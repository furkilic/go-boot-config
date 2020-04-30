package gobootconfig

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_parseProperties(t *testing.T) {
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
		{"Single", args{bytes.NewBufferString("go.string=myValue"), createMap()}, createMap("go.string", "myValue")},
		{"Comment", args{bytes.NewBufferString(" # go.string=myValue"), createMap()}, createMap()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseProperties(tt.args.read, tt.args.configs)
			if !reflect.DeepEqual(tt.args.configs, tt.want) {
				t.Errorf("parseProperties.getValue() = %v, want %v", tt.args.configs, tt.want)
			}
		})
	}
}
