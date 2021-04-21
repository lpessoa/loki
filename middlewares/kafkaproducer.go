package middlewares

import (
	"context"
	"net/http"

	"br.com.telecine/loki/kafka"
)

type customContextValuekey string

const (
	KafkaProducerInstance customContextValuekey = "kafkaProducer"
)

type KafkaProducerMiddleware struct {
	producer *kafka.Producer
}

// Initialize it somewhere
func (kpm *KafkaProducerMiddleware) Setup() {
	kpm.producer = kafka.CreateProducer()
}

// Middleware function, which will be called for each request
func (amw *KafkaProducerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), KafkaProducerInstance, amw.producer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
