package gobootconfig

import (
	"os"
	"strings"
)

func parseEnv() map[string]interface{} {
	configs := make(map[string]interface{})
	environ := os.Environ()
	for i, env := range environ {
		environ[i] = strings.ReplaceAll(env, "_", ".")
	}
	parseAndAdd(configs, environ)
	return configs
}
