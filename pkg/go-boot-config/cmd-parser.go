package gobootconfig

import (
	"os"
	"strings"
)

func parseCmdLine() map[string]interface{} {
	configs := make(map[string]interface{})
	args := os.Args
	l := len(args)
	i := 1
	cmdLines := make([]string, 0)
	for i < l {
		s := os.Args[i]
		if len(s) > 2 && strings.Contains(s, "--") {
			key := strings.TrimLeft(s, "--")
			if strings.Contains(key, "=") {
				cmdLines = append(cmdLines, key)
			} else if i < len(args)-1 {
				v := os.Args[i+1]
				if !strings.Contains(v, "--") {
					cmdLines = append(cmdLines, key+"="+v)
					i++
				} else {
					cmdLines = append(cmdLines, key+"=true")
				}
			} else {
				cmdLines = append(cmdLines, key+"=true")
			}
		}
		i++
	}
	parseAndAdd(configs, cmdLines)
	return configs
}
