package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"br.com.telecine/loki/core"
	"br.com.telecine/loki/kafka"
	"br.com.telecine/loki/middlewares"
	"github.com/gorilla/mux"
)

type TopicRely map[string]interface{}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(422)
	errorBody, _ := json.Marshal(TopicRely{
		"error": err.Error(),
	})
	w.Write(errorBody)
}

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	ctxProducer := r.Context().Value(middlewares.KafkaProducerInstance)
	producer := ctxProducer.(*kafka.Producer)

	vars := mux.Vars(r)
	topic := vars["topic"]
	eventInfo, err := core.GetEventInfo(topic)

	if err != nil {
		writeError(w, err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	data, _ := ioutil.ReadAll(r.Body)
	msgId, err := producer.PublishMessageToTopic(eventInfo.Topic, &data)

	w.Header().Add("content-type", "application/json")

	if err != nil {
		writeError(w, err)
		return
	}

	body := map[string]interface{}{
		"messageId": msgId,
	}
	reply, _ := json.Marshal(body)

	w.WriteHeader(201)
	w.Write(reply)
}
