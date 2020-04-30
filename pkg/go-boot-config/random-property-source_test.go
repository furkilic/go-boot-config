package gobootconfig

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"testing"
)

func Test_randomPropertySource_getSource(t *testing.T) {
	type fields struct {
		source string
		name   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Simple", fields{"mySource", "random.int"}, "mySource"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := randomPropertySource{
				source: tt.fields.source,
				name:   tt.fields.name,
			}
			if got := sps.getSource(); got != tt.want {
				t.Errorf("randomPropertySource.getSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randomPropertySource_getName(t *testing.T) {
	type fields struct {
		source string
		name   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Simple", fields{"mySource", "random.int"}, "random.int"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := randomPropertySource{
				source: tt.fields.source,
				name:   tt.fields.name,
			}
			if got := sps.getName(); got != tt.want {
				t.Errorf("randomPropertySource.getName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randomPropertySource_getValue(t *testing.T) {
	type fields struct {
		source string
		name   string
	}
	tests := []struct {
		name   string
		fields fields
		want   func(interface{}) bool
	}{
		{"Random Int", fields{"mySource", "random.int"},
			func(v interface{}) bool {
				if v == nil {
					return false
				}
				i, err := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
				return err == nil && i != 0
			},
		},
		{"Random Float", fields{"mySource", "random.float"},
			func(v interface{}) bool {
				if v == nil {
					return false
				}
				i, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
				return err == nil && i != 0
			}},
		{"Random UUID", fields{"mySource", "random.uuid"},
			func(v interface{}) bool {
				if v == nil {
					return false
				}
				_, err := uuid.Parse(fmt.Sprintf("%v", v))
				return err == nil
			}},
		{"Random Value", fields{"mySource", "random.value"},
			func(v interface{}) bool {
				return v != nil
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := randomPropertySource{
				source: tt.fields.source,
				name:   tt.fields.name,
			}
			if got := sps.getValue(); !tt.want(got) {
				t.Errorf("randomPropertySource.getValue() = %v", tt.want(got))
			}
		})
	}
}

func Test_loadFromRandom(t *testing.T) {
	_clear()
	tests := []struct {
		name string
	}{
		{"Load All"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadFromRandom()
			randoms := []string{"random.int", "random.float", "random.uuid", "random.value"}
			for _, random := range randoms {
				_, err := _getValue(random)
				if err != nil {
					t.Errorf("loadFromRandom %v : %v", random, err)
					return
				}

			}
		})
	}
}
