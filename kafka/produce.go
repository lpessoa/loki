package kafka

import (
	"encoding/binary"
	"errors"
	"fmt"

	"br.com.telecine/loki/core"

	"github.com/google/uuid"

	"github.com/riferrei/srclient"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var flushTimeout = 50

type IProducer interface {
	PublishMessageToTopic(eventItem *core.EventMapItem, payload *[]byte) (*string, error)
}

type Producer struct {
	kProducer    *kafka.Producer
	schemaClient *srclient.SchemaRegistryClient
}

func NewProducer() IProducer {
	schemaRegistryClient := srclient.CreateSchemaRegistryClient("http://schema-registry:8081")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "bitnami-docker-kafka_kafka_1:9092"})

	if err != nil {
		panic(err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &Producer{
		kProducer:    p,
		schemaClient: schemaRegistryClient,
	}
}

func (producer *Producer) getSchemaId(subject string) (*srclient.Schema, error) {
	return producer.schemaClient.GetLatestSchema(subject, false)
}

func (producer *Producer) makePayload(schema *srclient.Schema, data *[]byte) ([]byte, error) {
	/* setup schema ID */
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))
	/* set payload */
	native, _, err := schema.Codec().NativeFromTextual(*data)

	if err != nil {
		return nil, err
	}

	valueBytes, err := schema.Codec().BinaryFromNative(nil, native)

	if err != nil {
		return nil, err
	}

	var record []byte
	record = append(record, byte(0))
	record = append(record, schemaIDBytes...)
	record = append(record, valueBytes...)

	return record, nil
}

func (producer *Producer) getHeaders(eventItem *core.EventMapItem) []kafka.Header {
	messageId := uuid.NewString()
	return []kafka.Header{
		{
			Key:   "x-message-id",
			Value: []byte(messageId),
		},
		{
			Key:   "x-klzr-retry-count",
			Value: []byte(fmt.Sprintf("%v", eventItem.RetryCount)),
		},
	}
}

func (producer *Producer) PublishMessageToTopic(eventItem *core.EventMapItem, payload *[]byte) (*string, error) {
	schema, err := producer.getSchemaId(eventItem.Topic)
	if err != nil {
		panic(err)
	}
	pl, err := producer.makePayload(schema, payload)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("unable to create payload with provided schema")
	}

	messageId := uuid.NewString()
	headers := producer.getHeaders(eventItem)

	err = producer.kProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &eventItem.Topic, Partition: kafka.PartitionAny},
		Headers:        headers,
		Value:          pl,
	}, nil)

	go func() {
		producer.kProducer.Flush(flushTimeout)
	}()

	if err != nil {
		return nil, err
	}

	return &messageId, nil
}

func (producer *Producer) Close() {
	producer.kProducer.Close()
}
