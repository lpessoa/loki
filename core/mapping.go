package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

/** event map item */
type EventMapItem struct {
	Topic      string `json:"topic" yaml:"topic"`
	RetryCount int    `json:"retryCount" yaml:"retryCount"`
}

/** event mapping collection */
type EventMappings map[string]EventMapItem

func readEventFile(mappingFile *string) []byte {
	data, err := ioutil.ReadFile(*mappingFile)
	if err != nil {
		panic(err)
	}
	return data
}

type EventInfoProvider struct {
	yaml       bool
	eventItems EventMappings
}

func NewEventProvider(mappingFile *string, useYaml bool) *EventInfoProvider {
	provider := &EventInfoProvider{
		yaml: useYaml,
	}
	var items EventMappings
	data := readEventFile(mappingFile)

	if provider.yaml {
		err := yaml.Unmarshal(data, &items)
		if err != nil {
			panic("unable to parse yaml event mapping information")
		}
	} else {
		err := json.Unmarshal(data, &items)
		if err != nil {
			panic("unable to parse json event mapping information")
		}
	}
	provider.eventItems = items
	return provider
}

func (provider *EventInfoProvider) GetEventInfo(eventKey string) (*EventMapItem, error) {
	item, present := provider.eventItems[eventKey]

	if !present {
		return nil, errors.New("missing event information for provided topic")
	}

	return &item, nil
}
