package middlewares

import (
	"context"
	"net/http"

	"br.com.telecine/loki/kafka"
)

type customContextValuekey string

const (
	KafkaProducerInstance customContextValuekey = "kafkaProducer"
	TopicMappingFile      customContextValuekey = "mappingFile"
)

type KafkaProducerMiddleware struct {
	producer    kafka.IProducer
	mappingFile *string
}

// Initialize it somewhere
func (kpm *KafkaProducerMiddleware) Setup(mappingFile *string) {
	kpm.producer = kafka.NewProducer()
	kpm.mappingFile = mappingFile
}

// Middleware function, which will be called for each request
func (amw *KafkaProducerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), KafkaProducerInstance, amw.producer)
		ctx = context.WithValue(ctx, TopicMappingFile, amw.mappingFile)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
