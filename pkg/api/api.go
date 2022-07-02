package api

import "time"

//PASSKEY=xxxxxxxxxxxxxxxx&stationtype=EasyWeatherPro_V5.0.7&runtime=1&dateutc=2022-07-02+12:13:51&tempinf=78.1&humidityin=50&baromrelin=29.902&baromabsin=29.624&tempf=90.7&humidity=35&winddir=150&windspeedmph=1.34&windgustmph=2.24&maxdailygust=6.93&solarradiation=655.99&uv=6&rainratein=0.000&eventrainin=0.000&hourlyrainin=0.000&dailyrainin=0.000&weeklyrainin=0.000&monthlyrainin=0.000&yearlyrainin=0.000&totalrainin=0.000&wh65batt=0&freq=868M&model=WS2900_V2.01.18

type WeatherData struct {
	StationType    string `json:"stationtype"`
	Runtime        string `json:"runtime"`
	DateUTC        string `json:"dateutc"`
	TempInf        string `json:"tempinf"`
	HumidityIn     string `json:"humidityin"`
	BaromRelIn     string `json:"baromrelin"`
	BaromAbsIn     string `json:"baromabsin"`
	TempF          string `json:"tempf"`
	Humidity       string `json:"humidity"`
	WindDir        string `json:"winddir"`
	WindSpeedMph   string `json:"windspeedmph"`
	WindGustMph    string `json:"windgustmph"`
	MaxDailyGust   string `json:"maxdailygust"`
	SolarRadiation string `json:"solarradiation"`
	UV             string `json:"uv"`
	RainRateIn     string `json:"rainratein"`
	EventRainIn    string `json:"eventrainin"`
	HourlyRainIn   string `json:"hourlyrainin"`
	DailyRainIn    string `json:"dailyrainin"`
	WeeklyRainIn   string `json:"weeklyrainin"`
	MonthlyRainIn  string `json:"monthlyrainin"`
	YearlyRainIn   string `json:"yearlyrainin"`
	TotalRainIn    string `json:"totalrainin"`
	WH65Batt       string `json:"wh65batt"`
	Freq           string `json:"freq"`
	Model          string `json:"model"`
}

// CorrelationData represents any data, used for metrics or tracing.
type CorrelationData struct {
	// RequestID contains value of request id
	RequestID string `json:"requestID,omitempty"`

	// RequestTime is the time that the request was received
	RequestTime time.Time `json:"requestTime,omitempty"`
}
