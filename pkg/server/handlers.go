package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/mjudeikis/weewx-easyweather/pkg/api"

	"github.com/davecgh/go-spew/spew"
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

	data := string(body)
	spew.Dump(data)
	metrics := api.WeatherData{}
	parts := strings.Split(data, "&")
	for _, part := range parts {
		p := strings.Split(part, "=")
		setField(&metrics, p[0], p[1])
	}

	spew.Dump(metrics)

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

func setField(item interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJsonName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("json"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", fmt.Errorf("tag provided does not define a json tag", fieldName)
	}
	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		jname, _ := findJsonName(tag)
		fieldNames[jname] = i
	}

	fieldNum, ok := fieldNames[fieldName]
	if !ok {
		return fmt.Errorf("field %s does not exist within the provided item", fieldName)
	}
	fieldVal := v.Field(fieldNum)
	fieldVal.Set(reflect.ValueOf(value))
	return nil
}
