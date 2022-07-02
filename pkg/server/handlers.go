package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (s *Server) ingest(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	s.log.Info("ingest")
	start := time.Now()

	s.log.Info(r.RequestURI)
	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
		time.Since(start),
	)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	log.Println((string(body)))

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
