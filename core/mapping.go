package core

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
)

/** event map item */
type EventMapItem struct {
	Topic      string `json:"topic"`
	RetryCount int    `json:"retryCount"`
}

/** event mapping collection */
type EventMappings map[string]EventMapItem

func readEventFile() []byte {
	mappingFile := flag.String("mappings", "./eventMappings.json", "Event mapping jsonfile")
	data, err := ioutil.ReadFile(*mappingFile)
	if err != nil {
		panic(err)
	}
	return data
}

func GetEventInfo(eventKey string) (*EventMapItem, error) {
	data := readEventFile()
	var items EventMappings

	err := json.Unmarshal(data, &items)
	if err != nil {
		return nil, errors.New("unable to parse event mapping information")
	}
	item := items[eventKey]

	return &item, nil
}
