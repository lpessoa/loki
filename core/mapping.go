package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

/** event map item */
type EventMapItem struct {
	Topic      string `json:"topic"`
	RetryCount int    `json:"retryCount"`
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

func GetEventInfo(eventKey string, mappingFile *string) (*EventMapItem, error) {
	data := readEventFile(mappingFile)
	var items EventMappings

	err := json.Unmarshal(data, &items)
	if err != nil {
		return nil, errors.New("unable to parse event mapping information")
	}

	item, present := items[eventKey]

	if !present {
		return nil, errors.New("missing event information for provided topic")
	}

	return &item, nil
}
