package gobootconfig

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

type randomPropertySource struct {
	source string
	name   string
}

func (sps randomPropertySource) getSource() string {
	return sps.source
}
func (sps randomPropertySource) getName() string {
	return sps.name
}
func (sps randomPropertySource) getValue() interface{} {
	rand.Seed(time.Now().UnixNano())
	switch sps.name {
	case "random.int":
		return seededRand.Int63()
	case "random.float":
		return seededRand.Float64()
	case "random.uuid":
		return uuid.New().String()
	}
	b := make([]byte, 100)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func loadFromRandom() {
	_addPropertySource("random.int", randomPropertySource{"random", "random.int"})
	_addPropertySource("random.float", randomPropertySource{"random", "random.float"})
	_addPropertySource("random.uuid", randomPropertySource{"random", "random.uuid"})
	_addPropertySource("random.value", randomPropertySource{"random", "random.value"})
}
