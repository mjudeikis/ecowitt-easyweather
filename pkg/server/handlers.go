package server

import (
	"fmt"
	"net/http"
)

func (s *Server) ingest(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	s.log.Info("ingest")

	// Parse the request body into a struct.
	//var body api.WeatherData

	values := r.URL.Query()
	for k, v := range values {
		fmt.Println(k, " => ", v)
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
