package gobootconfig

func GetValue(key string) (interface{}, error) {
	startListener()
	request := valRequest{key, make(chan valResponse)}
	_getValCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetValueWithDefault(key string, def interface{}) interface{} {
	value, err := GetValue(key)
	if err == nil {
		return value
	}
	return def
}

func GetString(key string) (string, error) {
	startListener()
	request := stringRequest{key, make(chan stringResponse)}
	_getStringCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetStringWithDefault(key string, def string) string {
	value, err := GetString(key)
	if err == nil {
		return value
	}
	return def
}

func GetBool(key string) (bool, error) {
	startListener()
	request := boolRequest{key, make(chan boolResponse)}
	_getBoolCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetBoolWithDefault(key string, def bool) bool {
	value, err := GetBool(key)
	if err == nil {
		return value
	}
	return def
}

func GetInt(key string) (int64, error) {
	startListener()
	request := intRequest{key, make(chan intResponse)}
	_getIntCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetIntWithDefault(key string, def int64) int64 {
	value, err := GetInt(key)
	if err == nil {
		return value
	}
	return def
}

func GetFloat(key string) (float64, error) {
	startListener()
	request := floatRequest{key, make(chan floatResponse)}
	_getFloatCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetFloatWithDefault(key string, def float64) float64 {
	value, err := GetFloat(key)
	if err == nil {
		return value
	}
	return def
}

func GetSlice(key string) ([]interface{}, error) {
	startListener()
	request := sliceRequest{key, make(chan sliceResponse)}
	_getSliceCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetStringSlice(key string) ([]string, error) {
	startListener()
	request := stringSliceRequest{key, make(chan stringSliceResponse)}
	_getStringSliceCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetStringSliceWithDefault(key string, def []string) []string {
	value, err := GetStringSlice(key)
	if err == nil {
		return value
	}
	return def
}

func GetIntSlice(key string) ([]int64, error) {
	startListener()
	request := intSliceRequest{key, make(chan intSliceResponse)}
	_getIntSliceCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetIntSliceWithDefault(key string, def []int64) []int64 {
	value, err := GetIntSlice(key)
	if err == nil {
		return value
	}
	return def
}

func GetFloatSlice(key string) ([]float64, error) {
	startListener()
	request := floatSliceRequest{key, make(chan floatSliceResponse)}
	_getFloatSliceCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetFloatSliceWithDefault(key string, def []float64) []float64 {
	value, err := GetFloatSlice(key)
	if err == nil {
		return value
	}
	return def
}

func GetBoolSlice(key string) ([]bool, error) {
	startListener()
	request := boolSliceRequest{key, make(chan boolSliceResponse)}
	_getBoolSliceCh <- request
	resp := <-request.callback
	return resp.value, resp.err
}

func GetBoolSliceWithDefault(key string, def []bool) []bool {
	value, err := GetBoolSlice(key)
	if err == nil {
		return value
	}
	return def
}

func GetObject(key string, objPtr interface{}) error {
	startListener()
	request := objRequest{key, objPtr, make(chan objResponse)}
	_getObjCh <- request
	resp := <-request.callback
	return resp.err
}
