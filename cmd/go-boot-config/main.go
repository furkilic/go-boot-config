package main

import (
	"fmt"
	gobootconfig "go-boot-config/pkg/go-boot-config"
	"log"
)

type MyConfig struct {
	CmdArg   string
	EnvParam string
	Profile  MyProfile
	Conf     MyFiles
}

type MyProfile struct {
	Test MyFiles
	Conf MyFiles
}

type MyFiles struct {
	Properties string
	Yaml       string
}

func main() {
	err := gobootconfig.Reload()
	if err != nil {
		log.Fatalf("Error while loading the configuration: %v", err)
	}

	parameters := []string{
		"go.cmdArg",
		"go.envParam",
		"go.profile.test.properties",
		"go.profile.test.yaml",
		"go.profile.conf.properties",
		"go.profile.conf.yaml",
		"go.conf.properties",
		"go.conf.yaml",
	}

	for _, param := range parameters {
		fmt.Println(param, " => ", gobootconfig.GetStringWithDefault(param, "Does not exists"))
	}

	myConf := MyConfig{}
	err = gobootconfig.GetObject("go", &myConf)
	if err != nil {
		log.Fatalf("Error while loading go object: %v", err)
	}
	fmt.Println("go", "=>", myConf)

}
