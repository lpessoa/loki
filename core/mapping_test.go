package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mappingFile string = "./eventMappingsFixture.json"
)

func TestGetEventInfo(t *testing.T) {
	info, _ := GetEventInfo("test-1", &mappingFile)
	assert.Equal(t, info.Topic, "something")
}

func TestGetMissingEventInfo(t *testing.T) {
	info, _ := GetEventInfo("missing-topic", &mappingFile)
	assert.Nil(t, info)
}
