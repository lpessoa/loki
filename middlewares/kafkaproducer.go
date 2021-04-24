package middlewares

import (
	"context"
	"net/http"

	"br.com.telecine/loki/core"
	"br.com.telecine/loki/kafka"
)

type customContextValuekey string

const (
	KafkaProducerInstance customContextValuekey = "kafkaProducer"
	EventProvider         customContextValuekey = "eventProvider"
)

type KafkaProducerMiddleware struct {
	producer      kafka.IProducer
	eventProvider *core.EventInfoProvider
}

// Initialize it somewhere
func (kpm *KafkaProducerMiddleware) Setup(mappingFile *string, yaml bool) {
	kpm.producer = kafka.NewProducer()
	kpm.eventProvider = core.NewEventProvider(mappingFile, yaml)
}

// Middleware function, which will be called for each request
func (kpm *KafkaProducerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), KafkaProducerInstance, kpm.producer)
		ctx = context.WithValue(ctx, EventProvider, kpm.eventProvider)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
