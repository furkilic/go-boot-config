package gobootconfig

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"reflect"
)

func parseYaml(read io.Reader, configs map[string]interface{}) {
	b, _ := ioutil.ReadAll(read)
	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal(b, &m)
	if err != nil {

	}
	toMapRecursively("", m, configs)
}

func toMapRecursively(base string, m interface{}, s map[string]interface{}) {
	switch reflect.TypeOf(m).Kind().String() {
	case "map":
		for k, v := range m.(map[interface{}]interface{}) {
			sss := ""
			if base != sss {
				sss = base + "."
			}
			toMapRecursively(sss+k.(string), v, s)
		}
		break
	case "slice":
		for _, v := range m.([]interface{}) {
			toMapRecursively(base, v, s)
		}
		break
	default:
		if val, ok := s[base]; ok {
			if reflect.TypeOf(val).Kind().String() != "slice" {
				s[base] = []interface{}{val}
			}
			s[base] = append(s[base].([]interface{}), m)
		} else {
			s[base] = m
		}
	}
}
