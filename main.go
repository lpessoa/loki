package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"br.com.telecine/loki/core"
	lokiHandlers "br.com.telecine/loki/handlers"
	"br.com.telecine/loki/middlewares"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	mappingFile := flag.String("mappings", "./eventMappings.json", "Event mappings file")
	yaml := flag.Bool("yaml", false, "Mapping file is a YAML")
	flag.Parse()

	m := &middlewares.KafkaProducerMiddleware{}
	m.Setup(mappingFile, *yaml)

	r := mux.NewRouter()
	r.Use(m.Middleware)
	r.Handle("/events/{topic}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(lokiHandlers.TopicHandler)))

	listenAddr := fmt.Sprintf(":%v", core.GetEnv("PORT", "8080"))
	log.Fatal(http.ListenAndServe(listenAddr, r))
}
