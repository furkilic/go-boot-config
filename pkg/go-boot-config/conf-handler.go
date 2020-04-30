package gobootconfig

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var _environment = environment{time.Now(), make(map[string][]propertySource)}
var _loaded = false

var once sync.Once

var _loadCh = make(chan chan error)
var _reloadCh = make(chan chan error)

type valRequest struct {
	key      string
	callback chan valResponse
}
type valResponse struct {
	value interface{}
	err   error
}
type stringRequest struct {
	key      string
	callback chan stringResponse
}
type stringResponse struct {
	value string
	err   error
}

type boolRequest struct {
	key      string
	callback chan boolResponse
}
type boolResponse struct {
	value bool
	err   error
}

type intRequest struct {
	key      string
	callback chan intResponse
}
type intResponse struct {
	value int64
	err   error
}

type floatRequest struct {
	key      string
	callback chan floatResponse
}
type floatResponse struct {
	value float64
	err   error
}

type sliceRequest struct {
	key      string
	callback chan sliceResponse
}
type sliceResponse struct {
	value []interface{}
	err   error
}
type stringSliceRequest struct {
	key      string
	callback chan stringSliceResponse
}
type stringSliceResponse struct {
	value []string
	err   error
}

type boolSliceRequest struct {
	key      string
	callback chan boolSliceResponse
}
type boolSliceResponse struct {
	value []bool
	err   error
}

type intSliceRequest struct {
	key      string
	callback chan intSliceResponse
}
type intSliceResponse struct {
	value []int64
	err   error
}

type floatSliceRequest struct {
	key      string
	callback chan floatSliceResponse
}
type floatSliceResponse struct {
	value []float64
	err   error
}

type objRequest struct {
	key      string
	obj      interface{}
	callback chan objResponse
}
type objResponse struct {
	err error
}

var _getValCh = make(chan valRequest)
var _getStringCh = make(chan stringRequest)
var _getBoolCh = make(chan boolRequest)
var _getIntCh = make(chan intRequest)
var _getFloatCh = make(chan floatRequest)
var _getSliceCh = make(chan sliceRequest)
var _getStringSliceCh = make(chan stringSliceRequest)
var _getBoolSliceCh = make(chan boolSliceRequest)
var _getIntSliceCh = make(chan intSliceRequest)
var _getFloatSliceCh = make(chan floatSliceRequest)
var _getObjCh = make(chan objRequest)

func _startListener() {
	for {
		select {
		case loadCallback := <-_loadCh:
			loadCallback <- _load()
			break
		case refreshCallback := <-_reloadCh:
			refreshCallback <- _reload()
			break
		case req := <-_getValCh:
			value, e := _getValue(req.key)
			req.callback <- valResponse{value, e}
			break
		case req := <-_getStringCh:
			value, e := _getString(req.key)
			req.callback <- stringResponse{value, e}
			break
		case req := <-_getBoolCh:
			value, e := _getBool(req.key)
			req.callback <- boolResponse{value, e}
			break
		case req := <-_getIntCh:
			value, e := _getInt(req.key)
			req.callback <- intResponse{value, e}
			break
		case req := <-_getFloatCh:
			value, e := _getFloat(req.key)
			req.callback <- floatResponse{value, e}
			break
		case req := <-_getSliceCh:
			value, e := _getSlice(req.key)
			req.callback <- sliceResponse{value, e}
			break
		case req := <-_getStringSliceCh:
			value, e := _getStringSlice(req.key)
			req.callback <- stringSliceResponse{value, e}
			break
		case req := <-_getBoolSliceCh:
			value, e := _getBoolSlice(req.key)
			req.callback <- boolSliceResponse{value, e}
			break
		case req := <-_getIntSliceCh:
			value, e := _getIntSlice(req.key)
			req.callback <- intSliceResponse{value, e}
			break
		case req := <-_getFloatSliceCh:
			value, e := _getFloatSlice(req.key)
			req.callback <- floatSliceResponse{value, e}
			break
		case req := <-_getObjCh:
			req.callback <- objResponse{_getObject(req.key, req.obj)}
			break
		}
	}
}

func startListener() {
	once.Do(func() { go _startListener() })
}

func Load() error {
	startListener()
	callback := make(chan error)
	_loadCh <- callback
	return <-callback
}

func Reload() error {
	startListener()
	callback := make(chan error)
	_reloadCh <- callback
	return <-callback
}

// Sequential
func _load() error {
	if _loaded == false {
		_loaded = true
		loadConfiguration()
		return nil
	}
	return errors.New("configuration already loaded, please use conf.Reload() instead")
}
func _reload() error {
	_clear()
	loadConfiguration()
	return nil
}

func loadConfiguration() {
	loadFromCmdLine()
	loadFromEnvironment()
	loadFromRandom()
	loadFromProfileFile()
	loadFromFile()
}
func _clear() {
	_environment = environment{time.Now(), make(map[string][]propertySource)}
}

func _addPropertySource(key string, ps propertySource) {
	s := sanitize(key)
	_environment.propertySources[s] =
		append(_environment.propertySources[s], ps)
}

func _getValue(key string) (interface{}, error) {
	return _getValueWithMap(key, make(map[string]bool))
}

func _getString(key string) (string, error) {
	value, err := _getValue(key)
	if err == nil {
		return fmt.Sprintf("%v", value), nil
	}
	return "", err
}

func _getBool(key string) (bool, error) {
	valS, e := _getString(key)
	if e == nil {
		return strconv.ParseBool(valS)
	}
	return false, e
}

func _getInt(key string) (int64, error) {
	valS, e := _getString(key)
	if e == nil {
		return strconv.ParseInt(valS, 10, 64)
	}
	return 0, e
}
func _getFloat(key string) (float64, error) {
	valS, e := _getString(key)
	if e == nil {
		return strconv.ParseFloat(valS, 64)
	}
	return 0, e
}

func _getSlice(key string) ([]interface{}, error) {
	valS, e := _getValue(key)
	if e == nil {
		if reflect.TypeOf(valS).Kind().String() == "slice" {
			return valS.([]interface{}), nil
		} else {
			return []interface{}{valS}, nil
		}
	}
	return nil, e
}

func _getStringSlice(key string) ([]string, error) {
	value, err := _getSlice(key)
	if err == nil {
		valArr := make([]string, len(value))
		for i, v := range value {
			valArr[i] = fmt.Sprintf("%v", v)
		}
		return valArr, nil
	}
	return nil, err
}

func _getIntSlice(key string) ([]int64, error) {
	value, err := _getStringSlice(key)
	if err == nil {
		valArr := make([]int64, len(value))
		for i, v := range value {
			valArr[i], err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
		}
		return valArr, nil
	}
	return nil, err
}

func _getFloatSlice(key string) ([]float64, error) {
	value, err := _getStringSlice(key)
	if err == nil {
		valArr := make([]float64, len(value))
		for i, v := range value {
			valArr[i], err = strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
		}
		return valArr, nil
	}
	return nil, err
}

func _getBoolSlice(key string) ([]bool, error) {
	value, err := _getStringSlice(key)
	if err == nil {
		valArr := make([]bool, len(value))
		for i, v := range value {
			valArr[i], err = strconv.ParseBool(v)
			if err != nil {
				return nil, err
			}
		}
		return valArr, nil
	}
	return nil, err
}

func _getObject(key string, objPtr interface{}) error {
	m := buildMapFromKey(key)
	return buildObjectFromMap(m, objPtr)
}
