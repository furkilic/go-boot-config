package gobootconfig

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

func Test_filePropertySource_getSource(t *testing.T) {
	type fields struct {
		source   string
		location string
		name     string
		value    interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Simple", fields{"mySource", "myLocation", "myKey", "myValue"}, "mySource"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := filePropertySource{
				source:   tt.fields.source,
				location: tt.fields.location,
				name:     tt.fields.name,
				value:    tt.fields.value,
			}
			if got := sps.getSource(); got != tt.want {
				t.Errorf("filePropertySource.getSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filePropertySource_getLocation(t *testing.T) {
	type fields struct {
		source   string
		location string
		name     string
		value    interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Simple", fields{"mySource", "myLocation", "myKey", "myValue"}, "myLocation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := filePropertySource{
				source:   tt.fields.source,
				location: tt.fields.location,
				name:     tt.fields.name,
				value:    tt.fields.value,
			}
			if got := sps.getLocation(); got != tt.want {
				t.Errorf("filePropertySource.getLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filePropertySource_getName(t *testing.T) {
	type fields struct {
		source   string
		location string
		name     string
		value    interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Simple", fields{"mySource", "myLocation", "myKey", "myValue"}, "myKey"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := filePropertySource{
				source:   tt.fields.source,
				location: tt.fields.location,
				name:     tt.fields.name,
				value:    tt.fields.value,
			}
			if got := sps.getName(); got != tt.want {
				t.Errorf("filePropertySource.getName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filePropertySource_getValue(t *testing.T) {
	type fields struct {
		source   string
		location string
		name     string
		value    interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"Simple", fields{"mySource", "myLocation", "myKey", "myValue"}, "myValue"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sps := filePropertySource{
				source:   tt.fields.source,
				location: tt.fields.location,
				name:     tt.fields.name,
				value:    tt.fields.value,
			}
			if got := sps.getValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filePropertySource.getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFromProfileFile(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]interface{}
		want    interface{}
		wantErr bool
	}{
		{"Nothing", createMap(), nil, true},
		{"Properties Profile Default", createMap(goConfigLocation, "../../test/profile-default-properties"), "profile-default-properties", false},
		{"Yaml Profile Default", createMap(goConfigLocation, "../../test/profile-default-yaml"), "profile-default-yaml", false},
		{"Properties Default", createMap(goConfigLocation, "../../test/default-properties"), nil, true},
		{"Yaml Default", createMap(goConfigLocation, "../../test/default-yaml"), nil, true},
		{"Properties Profile Test", createMap(goConfigLocation, "../../test/profile-default-properties", goProfiles, "test"), "profile-test-properties", false},
		{"Yaml Profile Test", createMap(goConfigLocation, "../../test/profile-default-yaml", goProfiles, "test"), "profile-test-yaml", false},
		{"Properties Profile Test Not Active", createMap(goConfigLocation, "../../test/profile-default-properties", goProfiles, "test", goProfilesActive, "none"), nil, true},
		{"Yaml Profile Test Not Active", createMap(goConfigLocation, "../../test/profile-default-yaml", goProfiles, "test", goProfilesActive, "none"), nil, true},
		{"Properties file from location", createMap(goConfigLocation, "../../test/default-properties/application.properties"), nil, true},
		{"Yaml file from location", createMap(goConfigLocation, "../../test/default-yaml/application.yml"), nil, true},
	}
	for _, tt := range tests {
		_clear()
		for k, v := range tt.env {
			addInEnv(k, v)
		}
		t.Run(tt.name, func(t *testing.T) {
			loadFromProfileFile()
			got, err := _getValue("go.string")
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromProfileFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFromProfileFile.getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFromFile(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]interface{}
		want    interface{}
		wantErr bool
	}{
		{"Nothing", createMap(), nil, true},
		{"Properties Profile Default", createMap(goConfigLocation, "../../test/profile-default-properties"), "default-properties", false},
		{"Yaml Profile Default", createMap(goConfigLocation, "../../test/profile-default-yaml"), "default-properties", false},
		{"Properties Default", createMap(goConfigLocation, "../../test/default-properties"), "default-properties", false},
		{"Yaml Default", createMap(goConfigLocation, "../../test/default-yaml"), "default-yaml", false},
		{"Properties Profile Test", createMap(goConfigLocation, "../../test/profile-default-properties", goProfiles, "test"), "default-properties", false},
		{"Yaml Profile Test", createMap(goConfigLocation, "../../test/profile-default-yaml", goProfiles, "test"), "default-properties", false},
		{"Properties Profile Test Not Active", createMap(goConfigLocation, "../../test/profile-default-properties", goProfiles, "test", goProfilesActive, "none"), "default-properties", false},
		{"Yaml Profile Test Not Active", createMap(goConfigLocation, "../../test/profile-default-yaml", goProfiles, "test", goProfilesActive, "none"), "default-properties", false},
		{"Properties file from location", createMap(goConfigLocation, "../../test/default-properties/application.properties"), "default-properties", false},
		{"Yaml file from location", createMap(goConfigLocation, "../../test/default-yaml/application.yml"), "default-yaml", false},
	}
	for _, tt := range tests {
		_clear()
		for k, v := range tt.env {
			addInEnv(k, v)
		}
		t.Run(tt.name, func(t *testing.T) {
			loadFromFile()
			got, err := _getValue("go.string")
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFromFile.getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadAllFrom(t *testing.T) {
	go serveHttpFolder()
	type args struct {
		profiles   []string
		source     string
		extensions []string
		forProfile bool
		parser     func(io.Reader, map[string]interface{})
	}
	tests := []struct {
		name    string
		env     map[string]interface{}
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Nothing", createMap(), args{createArr(), "test", createArr(), true, parseProperties}, nil, true},
		{"Properties From File Default", createMap(goConfigLocation, "../../test/default-properties/application.properties"), args{createArr(), "test", createArr("properties"), false, parseProperties}, "default-properties", false},
		//{"Properties From HTTP", createMap(goConfigLocation, "http://localhost:3000/application.properties"), args{createArr(), "test", createArr("properties"), false, parseProperties}, "http-properties", false},
	}
	for _, tt := range tests {
		_clear()
		for k, v := range tt.env {
			addInEnv(k, v)
		}
		t.Run(tt.name, func(t *testing.T) {
			loadAllFrom(tt.args.profiles, tt.args.source, tt.args.extensions, tt.args.forProfile, tt.args.parser)
			got, err := _getValue("go.string")
			if (err != nil) != tt.wantErr {
				t.Errorf("loadAllFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadAllFrom.getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_retrieveProfiles(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]interface{}
		want []string
	}{
		{"None", createMap(), createArr("-default")},
		{"With Profile", createMap(goProfiles, "test"), createArr("-test")},
		{"With Active Profile", createMap(goProfiles, "test", goProfilesActive, "test"), createArr("-test")},
		{"Not Intersecting Active Profile", createMap(goProfiles, "test", goProfilesActive, "dev"), createArr()},
	}
	for _, tt := range tests {
		_clear()
		for k, v := range tt.env {
			addInEnv(k, v)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := retrieveProfiles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("retrieveProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filerActiveProfiles(t *testing.T) {
	type args struct {
		activeProfiles []string
		profiles       []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"None", args{createArr(), createArr()}, createArr()},
		{"All Match", args{createArr("test", "dev"), createArr("test", "dev")}, createArr("test", "dev")},
		{"Single Match", args{createArr("test"), createArr("test", "dev")}, createArr("test")},
		{"More Active Match", args{createArr("test", "dev"), createArr("test")}, createArr("test")},
		{"No Match", args{createArr("test"), createArr("dev")}, createArr()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filerActiveProfiles(tt.args.activeProfiles, tt.args.profiles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filerActiveProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFrom(t *testing.T) {
	f := func(m map[string]interface{}) func(string, func(io.Reader, map[string]interface{})) map[string]interface{} {
		return func(s string, i func(io.Reader, map[string]interface{})) map[string]interface{} {
			return m
		}
	}
	type args struct {
		path   string
		source string
		loader func(string, func(io.Reader, map[string]interface{})) map[string]interface{}
		parser func(io.Reader, map[string]interface{})
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"None", args{"myPath", "test", f(createMap()), nil}, nil, true},
		{"Simple", args{"myPath", "test", f(createMap("go.string", "myString")), nil}, "myString", false},
	}
	for _, tt := range tests {
		_clear()
		t.Run(tt.name, func(t *testing.T) {
			loadFrom(tt.args.path, tt.args.source, tt.args.loader, tt.args.parser)
			got, err := _getValue("go.string")
			if (err != nil) != tt.wantErr {
				t.Errorf("loadAllFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadAllFrom.getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFile(t *testing.T) {
	type args struct {
		filePath string
		read     func(io.Reader, map[string]interface{})
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"None", args{"none.properties", parseProperties}, createMap()},
		{"Simple", args{"../../test/default-properties/application.properties", parseProperties}, createMap("go.string", "default-properties")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadFile(tt.args.filePath, tt.args.read); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadHttp(t *testing.T) {
	go serveHttpFolder()
	type args struct {
		httpPath string
		read     func(io.Reader, map[string]interface{})
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"None", args{"http://localhost:3000/mamamia.properties", parseProperties}, createMap()},
		{"Simple", args{"http://localhost:3000//application.properties", parseProperties}, createMap("go.string", "http-properties")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadHttp(tt.args.httpPath, tt.args.read); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadHttp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_retrieveConfPath(t *testing.T) {
	go serveHttpFolder()
	type args struct {
		profiles   []string
		extensions []string
		forProfile bool
	}
	tests := []struct {
		name             string
		goConfigLocation string
		args             args
		want             []string
	}{
		{"None", "", args{createArr(""), createArr("properties"), false}, createArr()},
		{"Default Properties", filepath.Join("..", "..", "test", "profile-default-properties"), args{createArr(""), createArr("properties", "prop"), false}, createArr(filepath.Join("..", "..", "test", "profile-default-properties", "application.properties"))},
		{"Single file Properties", filepath.Join("..", "..", "test", "profile-default-properties", "application.properties"), args{createArr(""), createArr("properties", "prop"), false}, createArr(filepath.Join("..", "..", "test", "profile-default-properties", "application.properties"))},
		{"HTTP file Properties", "http://localhost:3000/application.properties", args{createArr(""), createArr("properties", "prop"), false}, createArr("http://localhost:3000/application.properties")},
	}
	for _, tt := range tests {
		_clear()
		if tt.goConfigLocation != "" {
			addInEnv(goConfigLocation, tt.goConfigLocation)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := retrieveConfPath(tt.args.profiles, tt.args.extensions, tt.args.forProfile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("retrieveConfPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isHttp(t *testing.T) {
	type args struct {
		location string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"File", args{"file://toto"}, false},
		{"HTTP", args{"http://toto"}, true},
		{"HTTPS", args{"https://toto"}, true},
		{"FTP", args{"ftp://toto"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHttp(tt.args.location); got != tt.want {
				t.Errorf("isHttp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathExists(t *testing.T) {
	type args struct {
		filePath string
		isDir    bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"None File", args{"none", false}, false, true},
		{"Existing File", args{"../../test/default-properties/application.properties", false}, true, false},
		{"Existing Not File", args{"../../test/default-properties", false}, false, false},

		{"None Directory", args{"none", true}, false, true},
		{"Existing Directory", args{"../../test/default-properties", true}, true, false},
		{"Existing Not Directory", args{"../../test/default-properties/application.properties", true}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pathExists(tt.args.filePath, tt.args.isDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("pathExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

var onceTest sync.Once

func serveHttpFolder() {
	onceTest.Do(func() {
		http.Handle("/", http.FileServer(http.Dir("../../test/http")))
		fmt.Println(http.ListenAndServe("0.0.0.0:3000", nil))
	})
}
