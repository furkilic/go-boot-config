package gobootconfig

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strconv"
	"strings"
)

func parseAndAdd(configs map[string]interface{}, entries []string) {
	for _, e := range entries {
		keyVal := strings.SplitN(e, "=", 2)
		ksyS := sanitize(keyVal[0])
		var value interface{}
		value = keyVal[1]
		if reflect.TypeOf(value).Kind().String() == "string" && strings.Contains(fmt.Sprintf("%s", value), ",") {
			split := strings.Split(value.(string), ",")
			valArr := make([]interface{}, len(split))
			for i, v := range split {
				valArr[i] = strings.TrimSpace(v)
			}
			value = valArr
		}
		if val, ok := configs[ksyS]; ok {
			if reflect.TypeOf(val).Kind().String() != "slice" {
				configs[ksyS] = []interface{}{val}
			}
			configs[ksyS] = append(configs[ksyS].([]interface{}), value)
		} else {
			configs[ksyS] = value
		}
	}
}

func toMap(m map[string]interface{}, s string, val interface{}) {
	dd := strings.SplitN(s, ".", 2)
	if len(dd) > 1 {
		if m[dd[0]] == nil {
			m[dd[0]] = make(map[string]interface{})
		}
		toMap(m[dd[0]].(map[string]interface{}), dd[1], val)
	} else {
		m[dd[0]] = val
	}
}

func _getValueWithMap(key string, visited map[string]bool) (interface{}, error) {
	indexStart := strings.Index(key, "[")
	indexEnd := strings.Index(key, "]")
	arrayIndex := int64(-1)
	if indexStart != -1 && indexEnd != -1 {
		index := string(key[indexStart+1 : indexEnd])
		arrayIndexT, err := strconv.ParseInt(index, 10, 64)
		if err != nil {
			return index, err
		}
		arrayIndex = arrayIndexT
	}
	ps, err := _getPropertySource(key)
	if err == nil {
		var value interface{}
		value, err = expand(key, ps.getValue(), visited)
		if err != nil {
			return nil, err
		}
		if arrayIndex != -1 {
			if reflect.TypeOf(value).Kind().String() == "slice" {
				return value.([]interface{})[arrayIndex], nil
			} else {
				return value, errors.New(fmt.Sprintf("'%s' expected array but was not", key))
			}
		}
		return value, nil
	}
	return nil, err
}

func expand(key string, value interface{}, visited map[string]bool) (interface{}, error) {
	if reflect.TypeOf(value).Kind().String() == "slice" {
		vSlice := make([]interface{}, len(value.([]interface{})))
		for i, v := range value.([]interface{}) {
			var err error
			vSlice[i], err = expand(key, v, visited)
			if err != nil {
				return nil, err
			}
		}
		return vSlice, nil
	} else {
		r := fmt.Sprintf("%v", value)
		indexStart := strings.Index(r, "${")
		indexEnd := strings.Index(r, "}")
		var t interface{}
		t = value
		if indexStart != -1 && indexEnd != -1 {
			b := indexStart == 0 && indexEnd == (len(r)-1)
			saniKey := sanitize(key)
			if _, ok := visited[saniKey]; ok {
				return nil, errors.New(fmt.Sprintf("'%s' had been call in cyclic", key))
			}
			visited[saniKey] = true
			var err error
			sanKey := strings.TrimSpace(string(r[indexStart+2 : indexEnd]))
			split := strings.Split(sanKey, ":-")
			t, err = _getValueWithMap(split[0], visited)
			if err != nil {
				if len(split) == 1 {
					return t, err
				} else {
					t = split[1]
				}
			}
			if !b {
				t = fmt.Sprintf("%s%s%s", string(r[0:indexStart]), t, string(r[indexEnd+1:]))
			}
			t, err = expand(split[0], t, visited)
			delete(visited, saniKey)
			if err != nil {
				return nil, err
			}
		}
		return t, nil
	}
}

func buildObjectFromMap(m map[string]interface{}, objPtr interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           objPtr,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(m)
}

func buildMapFromKey(key string) map[string]interface{} {
	m := make(map[string]interface{})
	starKey := fmt.Sprintf("%s.", sanitize(key))
	for k, _ := range _environment.propertySources {
		if len(key) > 0 {
			if strings.HasPrefix(k, starKey) {
				realKey := strings.TrimPrefix(k, starKey)
				value, _ := _getValue(k)
				toMap(m, realKey, value)
			}
		} else {
			value, _ := _getValue(k)
			toMap(m, k, value)
		}
	}
	return m
}

func _getPropertySource(key string) (propertySource, error) {
	sanitizedKey := sanitize(key)
	if _, ok := _environment.propertySources[sanitizedKey]; ok {
		return _environment.propertySources[sanitizedKey][0], nil
	}
	return nil, errors.New(fmt.Sprintf("'%s' key not found", key))
}

func sanitize(key string) string {
	indexStart := strings.Index(key, "[")
	indexEnd := strings.Index(key, "]")
	s := strings.ToLower(key)
	if indexStart != -1 && indexEnd != -1 {
		s = fmt.Sprintf("%s%s", s[0:indexStart], s[indexEnd+1:])
	}
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	return s
}
