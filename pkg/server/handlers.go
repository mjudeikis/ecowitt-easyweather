package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) ingest(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	s.log.Info("ingest")

	s.log.Info(r.RequestURI)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	data := map[string]interface{}{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		s.log.Info("Error unmarshalling body", zap.Error(err))
		return
	}

}

func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	s.log.Info("metrics")
	// Parse the request body into a struct.
	//var body api.WeatherData

	values := r.URL.Query()
	for k, v := range values {
		fmt.Println(k, " => ", v)
	}

}
