package kafka

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/riferrei/srclient"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var flushTimeout = 50

type Producer struct {
	kProducer    *kafka.Producer
	schemaClient *srclient.SchemaRegistryClient
}

func CreateProducer() *Producer {
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

func (producer *Producer) makePayload(schema *srclient.Schema, data *interface{}) ([]byte, error) {
	/* setup schema ID */
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))
	/* set payload */
	value, _ := json.Marshal(data)
	native, _, _ := schema.Codec().NativeFromTextual(value)
	valueBytes, _ := schema.Codec().BinaryFromNative(nil, native)

	var record []byte
	record = append(record, byte(0))
	record = append(record, schemaIDBytes...)
	record = append(record, valueBytes...)

	return record, nil
}

func (producer *Producer) getHeaders() []kafka.Header {
	messageId := uuid.NewString()
	return []kafka.Header{
		{
			Key:   "x-message-id",
			Value: []byte(messageId),
		},
	}
}

func (producer *Producer) PublishMessageToTopic(topic string, payload interface{}) (string, error) {
	schema, err := producer.getSchemaId(topic)
	if err != nil {
		panic(err)
	}
	pl, err := producer.makePayload(schema, &payload)
	if err != nil {
		panic("unable to create payload with provided schema")
	}

	messageId := uuid.NewString()
	headers := producer.getHeaders()

	err = producer.kProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Headers:        headers,
		Value:          pl,
	}, nil)
	go func() {
		producer.kProducer.Flush(flushTimeout)
	}()

	if err != nil {
		panic(err)
	}

	return messageId, nil
}

func (producer *Producer) Close() {
	producer.kProducer.Close()
}
