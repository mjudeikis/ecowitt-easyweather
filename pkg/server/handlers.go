package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/mjudeikis/ecowitt-easyweather/pkg/api"
	promutil "github.com/mjudeikis/ecowitt-easyweather/pkg/utils/prometheus"
	"go.uber.org/zap"
)

func (s *Server) ingest(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	s.log.Info("ingest")

	w.WriteHeader(http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Error("Error reading body", zap.Error(err))
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	data := string(body)
	metrics := &api.WeatherData{}
	parts := strings.Split(data, "&")
	for _, part := range parts {
		p := strings.Split(part, "=")
		setField(metrics, p[0], p[1])
	}

	s.setMetrics(metrics)

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

func (s *Server) setMetrics(m *api.WeatherData) {
	var v float64
	var err error
	v, err = strconv.ParseFloat(m.Runtime, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationRuntime.WithLabelValues(m.StationType, m.Model).Set(v)

	//v, err = strconv.ParseFloat(m.DateUTC, 64)
	//if err != nil {
	//	s.log.Error("Error parsing value", zap.Error(err))
	//}
	//promutil.WeatherStationDateUTC.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.TempInf, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationTemeratureFahrenheit.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.HumidityIn, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationHumidityIn.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.BaromRelIn, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationBaromRelIn.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.BaromAbsIn, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationBaromAbsIn.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.TempF, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationTempF.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.Humidity, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationHumidity.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.WindDir, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationWindDir.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.WindSpeedMph, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationWindSpeedMph.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.WindGustMph, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationWindGustMph.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.MaxDailyGust, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationMaxDailyGust.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.SolarRadiation, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationSolarRadiation.WithLabelValues(m.StationType, m.Model).Set(v)
	promutil.WeatherStationSolarRadiationTotal.WithLabelValues(m.StationType, m.Model).Add(v)

	v, err = strconv.ParseFloat(m.UV, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationUV.WithLabelValues(m.StationType, m.Model).Set(v)
	promutil.WeatherStationUVTotal.WithLabelValues(m.StationType, m.Model).Add(v)

	v, err = strconv.ParseFloat(m.RainRateIn, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationRainRateIn.Reset()
	promutil.WeatherStationRainRateIn.WithLabelValues(m.StationType, m.Model).Add(v)

	v, err = strconv.ParseFloat(m.HourlyRainIn, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationHourlyRainIn.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.DailyRainIn, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationDailyRainIn.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.MonthlyRainIn, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationMonthlyRainIn.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.YearlyRainIn, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationYearlyRainIn.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.TotalRainIn, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationTotalRainIn.WithLabelValues(m.StationType, m.Model).Set(v)

	v, err = strconv.ParseFloat(m.WH65Batt, 64)
	if err != nil {
		s.log.Error("Error parsing value", zap.Error(err))
	}
	promutil.WeatherStationWH65Batt.WithLabelValues(m.StationType, m.Model).Set(v)

}
