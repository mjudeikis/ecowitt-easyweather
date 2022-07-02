package config

import (
	"net/url"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerURI            string   `envconfig:"SERVER_URI" default:":8080"`
	LogLevel             string   `yaml:"logLevel,omitempty" envconfig:"LOG_LEVEL" default:"info"`
	ServerAllowedOrigins []string `envconfig:"SERVER_ALLOWED_ORIGIN" default:""`
}

func Load() (*Config, error) {
	c := &Config{}
	// 1. Load .env file
	godotenv.Load()

	// 2. Load ENV and defaults
	err := envconfig.Process("", c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (c *Config) GetAllowedOrigins() ([]url.URL, error) {
	var allowedOriginURLs []url.URL
	// split , separated and remove spaces
	for _, origin := range c.ServerAllowedOrigins {
		originURL, err := url.Parse(origin)
		if err != nil {
			return nil, err
		}
		allowedOriginURLs = append(allowedOriginURLs, *originURL)
	}
	return allowedOriginURLs, nil
}
