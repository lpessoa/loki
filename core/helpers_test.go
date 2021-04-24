package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	fallback := "fallbackvalue"
	fromEnv := GetEnv("NOT_DEFINED_VAR", fallback)
	assert.Equal(t, fromEnv, fallback)
}
