package gobootconfig

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	os.Args = []string{"cmd", "--test"}
	_loaded = false
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"First Load", false},
		{"Second Load", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			resp := GetBoolWithDefault("test", false)
			if resp != true {
				t.Errorf("Load() got = %v, want %v", resp, true)
			}
		})
	}
}

func TestReload(t *testing.T) {
	_clear()
	tests := []struct {
		name    string
		args    []string
		want    bool
		wantErr bool
	}{
		{"First ReLoad", []string{"cmd", "--test"}, true, false},
		{"Second ReLoad", []string{"cmd", "--test", "false"}, false, false},
	}
	for _, tt := range tests {
		os.Args = tt.args
		t.Run(tt.name, func(t *testing.T) {
			if err := Reload(); (err != nil) != tt.wantErr {
				t.Errorf("Reload() error = %v, wantErr %v", err, tt.wantErr)
			}
			resp := GetBoolWithDefault("test", false)
			if resp != tt.want {
				t.Errorf("Reload() got = %v, want %v", resp, tt.want)
			}
		})
	}
}

func Test_load(t *testing.T) {
	_clear()
	_loaded = false
	os.Args = []string{"cmd", "--test"}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"First Load", false},
		{"Second Load", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := _load(); (err != nil) != tt.wantErr {
				t.Errorf("_load() error = %v, wantErr %v", err, tt.wantErr)
			}
			resp := GetBoolWithDefault("test", false)
			if resp != true {
				t.Errorf("_load() got = %v, want %v", resp, true)
			}
		})
	}
}

func Test_reload(t *testing.T) {
	_clear()
	tests := []struct {
		name    string
		args    []string
		want    bool
		wantErr bool
	}{
		{"First ReLoad", []string{"cmd", "--test"}, true, false},
		{"Second ReLoad", []string{"cmd", "--test", "false"}, false, false},
	}
	for _, tt := range tests {
		os.Args = tt.args
		t.Run(tt.name, func(t *testing.T) {
			if err := _reload(); (err != nil) != tt.wantErr {
				t.Errorf("_reload() error = %v, wantErr %v", err, tt.wantErr)
			}
			resp := GetBoolWithDefault("test", false)
			if resp != tt.want {
				t.Errorf("_reload() got = %v, want %v", resp, tt.want)
			}
		})
	}
}

func Test_loadConfiguration(t *testing.T) {
	_clear()
	os.Args = []string{"cmd", "--test"}
	tests := []struct {
		name string
	}{
		{"First Load"},
		{"Second Load"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadConfiguration()
			resp := GetBoolWithDefault("test", false)
			if resp != true {
				t.Errorf("loadConfiguration() got = %v, want %v", resp, true)
			}
		})
	}
}

func Test_clear(t *testing.T) {
	_clear()
	lastTime := time.Now()
	tests := []struct {
		name string
	}{
		{"First Clean"},
		{"Second Clean"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_clear()
			if _environment.loadTime.Before(lastTime) {
				t.Errorf("clear() got loadTime= %v, wanted after %v", _environment.loadTime, lastTime)
			}
			lastTime = _environment.loadTime
			if len(_environment.propertySources) != 0 {
				t.Errorf("clear() got propertySources= %v, wanted empty", _environment.loadTime)
			}
		})
	}
}

func Test_addPropertySource(t *testing.T) {
	_clear()
	ps1 := staticPropertySource{"test", "key1", "val1"}
	ps2 := staticPropertySource{"test", "key2", "val2"}
	type args struct {
		key string
		ps  propertySource
	}
	tests := []struct {
		name string
		args args
		want []propertySource
	}{
		{"test1 key", args{"key1", ps1}, []propertySource{ps1}},
		{"test2 key", args{"key2", ps2}, []propertySource{ps2}},
		{"test1 multiple key", args{"key1", ps1}, []propertySource{ps1, ps1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_addPropertySource(tt.args.key, tt.args.ps)
			source := _environment.propertySources[sanitize(tt.args.key)]
			if !reflect.DeepEqual(source, tt.want) {
				t.Errorf("_addPropertySource() = %v, want %v", source, tt.want)
			}
		})
	}
}

func Test_getValue(t *testing.T) {
	addInEnv("simple", "myValue")
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"None", args{"none"}, nil, true},
		{"Simple", args{"simple"}, "myValue", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getValue(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getString(t *testing.T) {
	addInEnv("simple", "myValue")
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"None", args{"none"}, "", true},
		{"Simple", args{"simple"}, "myValue", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getString(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("_getString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBool(t *testing.T) {
	addInEnv("simple", "true")
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"None", args{"none"}, false, true},
		{"Simple", args{"simple"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getBool(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("_getBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getInt(t *testing.T) {
	addInEnv("simple", "123456789")
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"None", args{"none"}, 0, true},
		{"Simple", args{"simple"}, 123456789, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getInt(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("_getInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFloat(t *testing.T) {
	addInEnv("simple", "12345.6789")
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"None", args{"none"}, 0, true},
		{"Simple", args{"simple"}, 12345.6789, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getFloat(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("_getFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSlice(t *testing.T) {
	addInEnv("simple", []interface{}{"aa", "bb"})
	addInEnv("notarray", "aa")
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{"None", args{"none"}, nil, true},
		{"Simple", args{"simple"}, []interface{}{"aa", "bb"}, false},
		{"Not Array", args{"notarray"}, []interface{}{"aa"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_getSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringSlice(t *testing.T) {
	addInEnv("simple", []interface{}{"aa", "bb"})
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"None", args{"none"}, nil, true},
		{"Simple", args{"simple"}, []string{"aa", "bb"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getStringSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getStringSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_getStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIntSlice(t *testing.T) {
	addInEnv("simple", []interface{}{1, 2, 3})
	addInEnv("notparsable", []interface{}{1, "aa", 3})
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []int64
		wantErr bool
	}{
		{"None", args{"none"}, nil, true},
		{"Simple", args{"simple"}, []int64{1, 2, 3}, false},
		{"Not Parsable", args{"notparsable"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getIntSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getIntSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_getIntSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFloatSlice(t *testing.T) {
	addInEnv("simple", []interface{}{1.1, 2.2, 3.3})
	addInEnv("notparsable", []interface{}{1, "aa", 3})
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []float64
		wantErr bool
	}{
		{"None", args{"none"}, nil, true},
		{"Simple", args{"simple"}, []float64{1.1, 2.2, 3.3}, false},
		{"Not Parsable", args{"notparsable"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getFloatSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getFloatSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_getFloatSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBoolSlice(t *testing.T) {
	addInEnv("simple", []interface{}{true, false, true})
	addInEnv("notparsable", []interface{}{1, "aa", 3})
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []bool
		wantErr bool
	}{
		{"None", args{"none"}, nil, true},
		{"Simple", args{"simple"}, []bool{true, false, true}, false},
		{"Not Parsable", args{"notparsable"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getBoolSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getBoolSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_getBoolSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getObject(t *testing.T) {
	addInEnv("string", "myString")
	addInEnv("bool", "true")
	addInEnv("int", 1)
	addInEnv("float", 1.1)
	addInEnv("stringslice", []interface{}{"my", "string"})
	addInEnv("boolslice", []interface{}{true, false, true})
	addInEnv("intslice", []interface{}{1, 2, 3})
	addInEnv("floatslice", []interface{}{1.1, 2.2, 3.3})
	addInEnv("child.string", "myString")
	addInEnv("child.bool", "true")
	addInEnv("child.int", 1)
	addInEnv("child.float", 1.1)
	addInEnv("child.stringslice", []interface{}{"my", "string"})
	addInEnv("child.boolslice", []interface{}{true, false, true})
	addInEnv("child.intslice", []interface{}{1, 2, 3})
	addInEnv("child.floatslice", []interface{}{1.1, 2.2, 3.3})

	objStruct := testObjStruct{"myString", true, 1, 1.1, []string{"my", "string"}, []bool{true, false, true}, []int64{1, 2, 3}, []float64{1.1, 2.2, 3.3}}
	objParentStruct := testObjParentStruct{objStruct}

	type args struct {
		key    string
		objPtr interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"None", args{"none", &testObjStruct{}}, &testObjStruct{}, false},
		{"No Key", args{"", &testObjStruct{}}, &objStruct, false},
		{"No Key With Parent", args{"", &testObjParentStruct{}}, &objParentStruct, false},
		{"With Key", args{"child", &testObjStruct{}}, &objStruct, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := _getObject(tt.args.key, tt.args.objPtr); (err != nil) != tt.wantErr {
				t.Errorf("_getObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.objPtr, tt.want) {
				t.Errorf("_getObject() = %v, want %v", tt.args.objPtr, tt.want)
			}
		})
	}
}
