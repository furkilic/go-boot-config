package gobootconfig

import (
	"time"
)

type environment struct {
	loadTime        time.Time
	propertySources map[string][]propertySource
}

type propertySource interface {
	getSource() string
	getName() string
	getValue() interface{}
}
