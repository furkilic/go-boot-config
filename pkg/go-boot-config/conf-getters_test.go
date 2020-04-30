package gobootconfig

import (
	"reflect"
	"testing"
)

func addInEnv(key string, val interface{}) {
	_environment.propertySources[key] =
		append([]propertySource{}, staticPropertySource{"test", key, val})
}

func TestGetValue(t *testing.T) {
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
			got, err := GetValue(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValueWithDefault(t *testing.T) {
	addInEnv("simple", "myValue")
	type args struct {
		key string
		def interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"None", args{"none", "myDefault"}, "myDefault"},
		{"Simple", args{"simple", "myDefault"}, "myValue"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValueWithDefault(tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("GetValueWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetString(t *testing.T) {
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
			got, err := GetString(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringWithDefault(t *testing.T) {
	addInEnv("simple", "myValue")
	type args struct {
		key string
		def string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"None", args{"none", "myDefault"}, "myDefault"},
		{"Simple", args{"simple", "myDefault"}, "myValue"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringWithDefault(tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("GetStringWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
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
			got, err := GetBool(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBoolWithDefault(t *testing.T) {
	addInEnv("simple", "true")
	type args struct {
		key string
		def bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"None", args{"none", true}, true},
		{"Simple", args{"simple", false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBoolWithDefault(tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("GetBoolWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
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
			got, err := GetInt(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntWithDefault(t *testing.T) {
	addInEnv("simple", "123456789")
	type args struct {
		key string
		def int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"None", args{"none", 111111}, 111111},
		{"Simple", args{"simple", 111111}, 123456789},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIntWithDefault(tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("GetIntWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFloat(t *testing.T) {
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
			got, err := GetFloat(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFloatWithDefault(t *testing.T) {
	addInEnv("simple", "12345.6789")
	type args struct {
		key string
		def float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"None", args{"none", 111.111}, 111.111},
		{"Simple", args{"simple", 111.111}, 12345.6789},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFloatWithDefault(tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("GetFloatWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSlice(t *testing.T) {
	addInEnv("simple", []interface{}{"aa", "bb"})
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringSlice(t *testing.T) {
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
			got, err := GetStringSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStringSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringSliceWithDefault(t *testing.T) {
	addInEnv("simple", []interface{}{"aa", "bb"})
	def := []string{"def"}
	type args struct {
		key string
		def []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"None", args{"none", def}, def},
		{"Simple", args{"simple", def}, []string{"aa", "bb"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringSliceWithDefault(tt.args.key, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStringSliceWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntSlice(t *testing.T) {
	addInEnv("simple", []interface{}{1, 2, 3})
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIntSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIntSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIntSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntSliceWithDefault(t *testing.T) {
	addInEnv("simple", []interface{}{1, 2, 3})
	def := []int64{1}
	type args struct {
		key string
		def []int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{"None", args{"none", def}, def},
		{"Simple", args{"simple", def}, []int64{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIntSliceWithDefault(tt.args.key, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIntSliceWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFloatSlice(t *testing.T) {
	addInEnv("simple", []interface{}{1.1, 2.2, 3.3})
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFloatSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFloatSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFloatSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFloatSliceWithDefault(t *testing.T) {
	addInEnv("simple", []interface{}{1.1, 2.2, 3.3})
	def := []float64{1.2}
	type args struct {
		key string
		def []float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{"None", args{"none", def}, def},
		{"Simple", args{"simple", def}, []float64{1.1, 2.2, 3.3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFloatSliceWithDefault(tt.args.key, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFloatSliceWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBoolSlice(t *testing.T) {
	addInEnv("simple", []interface{}{true, false, true})
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBoolSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBoolSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBoolSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBoolSliceWithDefault(t *testing.T) {
	addInEnv("simple", []interface{}{true, false, true})
	def := []bool{true}
	type args struct {
		key string
		def []bool
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{"None", args{"none", def}, def},
		{"Simple", args{"simple", def}, []bool{true, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBoolSliceWithDefault(tt.args.key, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBoolSliceWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testObjStruct struct {
	String      string
	Bool        bool
	Int         int64
	Float       float64
	StringSlice []string
	BoolSlice   []bool
	IntSlice    []int64
	FloatSlice  []float64
}

type testObjParentStruct struct {
	Child testObjStruct
}

func TestGetObject(t *testing.T) {
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
			if err := GetObject(tt.args.key, tt.args.objPtr); (err != nil) != tt.wantErr {
				t.Errorf("GetObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.objPtr, tt.want) {
				t.Errorf("GetObject() = %v, want %v", tt.args.objPtr, tt.want)
			}
		})
	}
}
