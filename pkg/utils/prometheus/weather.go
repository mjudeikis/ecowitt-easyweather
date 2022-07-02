package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(WeatherStationRuntime)
	prometheus.MustRegister(WeatherStationDateUTC)
	prometheus.MustRegister(WeatherStationTemeratureFahrenheit)
	prometheus.MustRegister(WeatherStationHumidityIn)
	prometheus.MustRegister(WeatherStationBaromRelIn)
	prometheus.MustRegister(WeatherStationBaromAbsIn)
	prometheus.MustRegister(WeatherStationTempF)
	prometheus.MustRegister(WeatherStationHumidity)
	prometheus.MustRegister(WeatherStationWindDir)
	prometheus.MustRegister(WeatherStationWindSpeedMph)
	prometheus.MustRegister(WeatherStationWindGustMph)
	prometheus.MustRegister(WeatherStationMaxDailyGust)
	prometheus.MustRegister(WeatherStationSolarRadiation)
	prometheus.MustRegister(WeatherStationUV)
	prometheus.MustRegister(WeatherStationRainRateIn)
	prometheus.MustRegister(WeatherStationHourlyRainIn)
	prometheus.MustRegister(WeatherStationDailyRainIn)
	prometheus.MustRegister(WeatherStationMonthlyRainIn)
	prometheus.MustRegister(WeatherStationYearlyRainIn)
	prometheus.MustRegister(WeatherStationTotalRainIn)
	prometheus.MustRegister(WeatherStationWH65Batt)
}

var (
	WeatherStationRuntime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_station_runtime",
		Help: "Station runtime",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationDateUTC = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_station_report_time",
		Help: "Station last report time",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationTemeratureFahrenheit = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_tempinf",
		Help: "tempinf",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationHumidityIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_humidityin",
		Help: "humidityin",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationBaromRelIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_baromrelin",
		Help: "baromrelin",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationBaromAbsIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_baromabsin",
		Help: "baromabsin",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationTempF = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_tempf",
		Help: "tempf",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationHumidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_humidity",
		Help: "humidity",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationWindDir = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_winddir",
		Help: "winddir",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationWindSpeedMph = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_windspeedmph",
		Help: "windspeedmph",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationWindGustMph = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_windgustmph",
		Help: "windgustmph",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationMaxDailyGust = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_maxdailygust",
		Help: "maxdailygust",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationSolarRadiation = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_solarradiation",
		Help: "solarradiation",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationUV = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_uv",
		Help: "uv",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationRainRateIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_rainratein",
		Help: "rainratein",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationEventRainIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_eventrainin",
		Help: "eventrainin",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationHourlyRainIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_hourlyrainin",
		Help: "hourlyrainin",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationDailyRainIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_dailyrainin",
		Help: "dailyrainin",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationMonthlyRainIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_monthlyrainin",
		Help: "monthlyrainin",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationYearlyRainIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_yearlyrainin",
		Help: "yearlyrainin",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationTotalRainIn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_totalrainin",
		Help: "totalrainin",
	},
		[]string{"weather_station", "model"},
	)

	WeatherStationWH65Batt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_wh65battn",
		Help: "wh65batt",
	},
		[]string{"weather_station", "model"},
	)
)
