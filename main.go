package main

import (
	"flag"
	"fmt"
	"io"
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

	var handlerWriter io.Writer = os.Stdout
	logOutput := core.GetLogger()
	if logOutput != nil {
		handlerWriter = logOutput
		defer logOutput.Close()
	}

	m := &middlewares.KafkaProducerMiddleware{}
	m.Setup(mappingFile, *yaml)

	r := mux.NewRouter()
	r.Use(m.Middleware)
	r.Handle("/events/{topic}", handlers.LoggingHandler(handlerWriter, http.HandlerFunc(lokiHandlers.TopicHandler)))

	listenAddr := fmt.Sprintf(":%v", core.GetEnv("PORT", "8080"))
	log.Fatal(http.ListenAndServe(listenAddr, r))
}
