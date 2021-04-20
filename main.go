package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"br.com.telecine/loki/kafka"
	"br.com.telecine/loki/middlewares"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type J struct {
	Name string `json:"name"`
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	ctxProducer := r.Context().Value(middlewares.KafkaProducerInstance)
	producer := ctxProducer.(*kafka.Producer)

	b := &J{
		Name: "go bananas",
	}
	topic := "test-1"
	msgId, _ := producer.PublishMessageToTopic(topic, b)

	w.Write([]byte(fmt.Sprintf("Message published to %v! with %v\n", topic, msgId)))
}

func main() {
	m := &middlewares.KafkaProducerMiddleware{}
	m.Setup()
	r := mux.NewRouter()
	r.Use(m.Middleware)
	// Routes consistof a path and a handler function.
	r.Handle("/events", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(YourHandler)))

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))
}
