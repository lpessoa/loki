package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mappingFile  string             = "./eventMappingsFixture.json"
	jsonProvider *EventInfoProvider = NewEventProvider(&mappingFile, false)
	yamlProvider *EventInfoProvider = NewEventProvider(&mappingFile, true)
)

func TestGetEventInfoFromJSON(t *testing.T) {
	info, _ := jsonProvider.GetEventInfo("test-1")
	assert.Equal(t, info.Topic, "something")
}

func TestGetMissingEventInfoFromJSON(t *testing.T) {
	info, _ := jsonProvider.GetEventInfo("missing-topic")
	assert.Nil(t, info)
}

func TestGetEventInfoFromYAML(t *testing.T) {
	info, _ := yamlProvider.GetEventInfo("test-1")
	assert.Equal(t, info.Topic, "something")
}

func TestGetMissingEventInfoFromYAML(t *testing.T) {
	info, _ := yamlProvider.GetEventInfo("missing-topic")
	assert.Nil(t, info)
}
