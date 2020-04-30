package gobootconfig

import (
	"reflect"
	"testing"
)

func Test_parseAndAdd(t *testing.T) {
	type args struct {
		configs map[string]interface{}
		entries []string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"Simple", args{make(map[string]interface{}), createArr("simple=myString")}, createMapKeyValue("simple", "myString")},
		{"Simple to Array", args{make(map[string]interface{}), createArr("simple=myString", "simple=myString2")}, createMapKeyValue("simple", []interface{}{"myString", "myString2"})},
		{"Array", args{make(map[string]interface{}), createArr("simple=myString, myString2")}, createMapKeyValue("simple", []interface{}{"myString", "myString2"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseAndAdd(tt.args.configs, tt.args.entries)
			if !reflect.DeepEqual(tt.args.configs, tt.want) {
				t.Errorf("parseAndAdd() = %v, want %v", tt.args.configs, tt.want)
			}
		})
	}
}

func Test_toMap(t *testing.T) {
	type args struct {
		m   map[string]interface{}
		s   string
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"Simple", args{make(map[string]interface{}), "simple", "myString"}, createMapKeyValue("simple", "myString")},
		{"Nested Object", args{make(map[string]interface{}), "simple.nested", "myString"}, createMapKeyValue("simple", createMapKeyValue("nested", "myString"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toMap(tt.args.m, tt.args.s, tt.args.val)
			if !reflect.DeepEqual(tt.args.m, tt.want) {
				t.Errorf("toMap() = %v, want %v", tt.args.m, tt.want)
			}
		})
	}
}

func Test_getValueWithMap(t *testing.T) {
	_clear()
	addInEnv("simple", "myString")
	addInEnv("array", []interface{}{"arr0", "arr1", "arr2"})
	addInEnv("expand", "${simple}")
	addInEnv("notexpand", "${none}")
	addInEnv("expanddefault", "${none:-defvalue}")
	type args struct {
		key     string
		visited map[string]bool
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Not Exist", args{"none", make(map[string]bool)}, nil, true},
		{"Simple", args{"simple", make(map[string]bool)}, "myString", false},
		{"Array", args{"array[2]", make(map[string]bool)}, "arr2", false},
		{"Not Array", args{"simple[2]", make(map[string]bool)}, "myString", true},
		{"Not Array Index", args{"array[aa]", make(map[string]bool)}, "aa", true},
		{"Expandable", args{"expand", make(map[string]bool)}, "myString", false},
		{"Not Expandable", args{"notexpand", make(map[string]bool)}, nil, true},
		{"Expandable with default", args{"expanddefault", make(map[string]bool)}, "defvalue", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getValueWithMap(tt.args.key, tt.args.visited)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getValueWithMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_getValueWithMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_expand(t *testing.T) {
	_clear()
	addInEnv("simple", "myString")
	addInEnv("expand", "${simple}")
	addInEnv("expanddefault", "${none:-defvalue}")
	addInEnv("cyclic1", "${cyclic2}")
	addInEnv("cyclic2", "${cyclic1}")
	type args struct {
		key     string
		value   interface{}
		visited map[string]bool
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Simple", args{"simple", "myString", make(map[string]bool)}, "myString", false},
		{"Simple Expand", args{"expand", "${simple}", make(map[string]bool)}, "myString", false},
		{"Not Expandable", args{"notexpand", "${nonde}", make(map[string]bool)}, nil, true},
		{"Expandable with default", args{"expanddefault", "${none:-defvalue}", make(map[string]bool)}, "defvalue", false},
		{"Extra Char Expand", args{"extraexpand", "AA${simple}AA", make(map[string]bool)}, "AAmyStringAA", false},
		{"Cyclic", args{"cyclic1", "${cyclic2}", make(map[string]bool)}, nil, true},
		{"Multiple Expand", args{"multiple", "${simple} ${simple}", make(map[string]bool)}, "myString myString", false},
		{"Multiple Expand Error", args{"multiple", "${simple} ${none}", make(map[string]bool)}, nil, true},
		{"Array Expand", args{"array", []interface{}{"${simple}", "${simple}"}, make(map[string]bool)}, []interface{}{"myString", "myString"}, false},
		{"Array Expand With Error", args{"array", []interface{}{"${simple}", "${none}"}, make(map[string]bool)}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expand(tt.args.key, tt.args.value, tt.args.visited)
			if (err != nil) != tt.wantErr {
				t.Errorf("expand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildObjectFromMap(t *testing.T) {
	_clear()

	objStruct := testObjStruct{"myString", true, 1, 1.1, []string{"my", "string"}, []bool{true, false, true}, []int64{1, 2, 3}, []float64{1.1, 2.2, 3.3}}
	objParentStruct := testObjParentStruct{objStruct}

	type args struct {
		m      map[string]interface{}
		objPtr interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Simple", args{createMapKeyValue("string", "myString", "bool", "true", "int", 1, "float", 1.1, "stringslice", []interface{}{"my", "string"}, "boolslice", []interface{}{true, false, true}, "intslice", []interface{}{1, 2, 3}, "floatslice", []interface{}{1.1, 2.2, 3.3}), &testObjStruct{}}, &objStruct, false},
		{"Nested", args{createMapKeyValue("child", createMapKeyValue("string", "myString", "bool", "true", "int", 1, "float", 1.1, "stringslice", []interface{}{"my", "string"}, "boolslice", []interface{}{true, false, true}, "intslice", []interface{}{1, 2, 3}, "floatslice", []interface{}{1.1, 2.2, 3.3})), &testObjParentStruct{}}, &objParentStruct, false},
		{"Error", args{createMap(), nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := buildObjectFromMap(tt.args.m, tt.args.objPtr); (err != nil) != tt.wantErr {
				t.Errorf("buildObjectFromMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.objPtr, tt.want) {
				t.Errorf("buildObjectFromMap() = %v, want %v", tt.args.objPtr, tt.want)
			}
		})
	}
}

func Test_buildMapFromKey(t *testing.T) {
	_clear()
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
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"All", args{""}, createMapKeyValue("string", "myString", "bool", "true", "int", 1, "float", 1.1, "stringslice", []interface{}{"my", "string"}, "boolslice", []interface{}{true, false, true}, "intslice", []interface{}{1, 2, 3}, "floatslice", []interface{}{1.1, 2.2, 3.3}, "child", createMapKeyValue("string", "myString", "bool", "true", "int", 1, "float", 1.1, "stringslice", []interface{}{"my", "string"}, "boolslice", []interface{}{true, false, true}, "intslice", []interface{}{1, 2, 3}, "floatslice", []interface{}{1.1, 2.2, 3.3}))},
		{"Nested", args{"child"}, createMapKeyValue("string", "myString", "bool", "true", "int", 1, "float", 1.1, "stringslice", []interface{}{"my", "string"}, "boolslice", []interface{}{true, false, true}, "intslice", []interface{}{1, 2, 3}, "floatslice", []interface{}{1.1, 2.2, 3.3})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildMapFromKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildMapFromKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPropertySource(t *testing.T) {
	_clear()
	addInEnv("simple", "myString")
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    propertySource
		wantErr bool
	}{
		{"Simple", args{"simple"}, staticPropertySource{"test", "simple", "myString"}, false},
		{"Not Exists", args{"none"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _getPropertySource(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("_getPropertySource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_getPropertySource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sanitize(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Simple", args{"simple"}, "simple"},
		{"Nested", args{"nested.simple"}, "nested.simple"},
		{"Kebab Case", args{"child.kebab-case"}, "child.kebabcase"},
		{"Camel Case", args{"child.camelCase"}, "child.camelcase"},
		{"Snake Case", args{"child.snake_case"}, "child.snakecase"},
		{"Array", args{"child.arr[1]"}, "child.arr"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitize(tt.args.key); got != tt.want {
				t.Errorf("sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}
