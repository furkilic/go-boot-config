package gobootconfig

import (
	"bufio"
	"io"
	"strings"
)

func parseProperties(read io.Reader, configs map[string]interface{}) {
	scanner := bufio.NewScanner(read)
	scanner.Split(bufio.ScanLines)
	var properties []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] != '#' && strings.Contains(line, "=") {
			properties = append(properties, line)
		}
	}
	parseAndAdd(configs, properties)
}
