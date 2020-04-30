package main

import (
	gobootconfig "github.com/furkilic/go-boot-config/pkg/go-boot-config"
	"os"
	"testing"
)

func ExampleMain_all() {
	os.Args = []string{"cmd", "--go.cmdArg", "cmd"}
	os.Setenv("GO_ENVPARAM", "env")
	os.Setenv("GO_CONFIG_LOCATION", "../../configs")
	main()
	// Output:
	// go.cmdArg  =>  cmd
	// go.envParam  =>  env
	// go.profile.test.properties  =>  Does not exists
	// go.profile.test.yaml  =>  Does not exists
	// go.profile.conf.properties  =>  profile-default-properties
	// go.profile.conf.yaml  =>  profile-default-yaml
	// go.conf.properties  =>  default-properties
	// go.conf.yaml  =>  default-yaml
	// go => {cmd env {{ } {profile-default-properties profile-default-yaml}} {default-properties default-yaml}}

}

func ExampleMain_profile() {
	os.Args = []string{"cmd", "--go.cmdArg", "cmd", "--go.profiles", "test"}
	os.Setenv("GO_ENVPARAM", "env")
	os.Setenv("GO_CONFIG_LOCATION", "../../configs")
	main()
	// Output:
	// go.cmdArg  =>  cmd
	// go.envParam  =>  env
	// go.profile.test.properties  =>  profile-test-properties
	// go.profile.test.yaml  =>  profile-test-yaml
	// go.profile.conf.properties  =>  Does not exists
	// go.profile.conf.yaml  =>  Does not exists
	// go.conf.properties  =>  default-properties
	// go.conf.yaml  =>  default-yaml
	// go => {cmd env {{profile-test-properties profile-test-yaml} { }} {default-properties default-yaml}}

}

func BenchmarkReload(b *testing.B) {
	os.Args = []string{"cmd", "--go.cmdArg", "cmd"}
	os.Setenv("GO_ENVPARAM", "env")
	os.Setenv("GO_CONFIG_LOCATION", "../../configs")
	for i := 0; i < b.N; i++ {
		gobootconfig.Reload()
	}
}
