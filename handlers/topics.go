package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"br.com.telecine/loki/kafka"
	"br.com.telecine/loki/middlewares"
	"github.com/gorilla/mux"
)

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	ctxProducer := r.Context().Value(middlewares.KafkaProducerInstance)
	producer := ctxProducer.(*kafka.Producer)

	vars := mux.Vars(r)
	topic := vars["topic"]

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	data, _ := ioutil.ReadAll(r.Body)
	msgId, _ := producer.PublishMessageToTopic(topic, &data)

	w.Header().Add("content-type", "application/json")
	body := map[string]interface{}{
		"messageId": msgId,
	}
	reply, _ := json.Marshal(body)

	w.WriteHeader(201)
	w.Write(reply)
}
