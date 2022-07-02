package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		targetMux.ServeHTTP(w, r)

		// log request by who(IP address)
		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t\t%s\t\t%s\t\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		log.Println((string(body)))
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	html := ""
	w.Write([]byte(html))
}

func main() {

	// direct all log messages to webrequests.log
	log.SetOutput(os.Stdout)

	log.Println("Start logger")
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	http.ListenAndServe(":9080", RequestLogger(mux))
}
