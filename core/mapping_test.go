package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mappingFile  string             = "./eventMappingsFixture.json"
	jsonProvider *EventInfoProvider = NewEventProvider(&mappingFile, false)
)

func TestGetEventInfo(t *testing.T) {
	info, _ := jsonProvider.GetEventInfo("test-1")
	assert.Equal(t, info.Topic, "something")
}

func TestGetMissingEventInfo(t *testing.T) {
	info, _ := jsonProvider.GetEventInfo("missing-topic")
	assert.Nil(t, info)
}
