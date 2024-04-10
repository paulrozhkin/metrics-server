package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const DefaultRequestTimeoutSeconds = 60

type Configuration struct {
	Server     ServerConfigurations
	Production bool `yaml:"production"`
}

type ServerConfigurations struct {
	Host                  string `yaml:"host"`
	Port                  int    `yaml:"port"`
	RequestTimeoutSeconds int    `yaml:"requestTimeoutSeconds"`
	DisableCORS           bool   `yaml:"disableCORS"`
}

func LoadConfigurations() (*Configuration, error) {
	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("METRICS_SERVER")

	_ = viper.BindEnv("production", "METRICS_SERVER_PRODUCTION")
	_ = viper.BindEnv("server.port", "METRICS_SERVER_PORT")
	_ = viper.BindEnv("server.host", "METRICS_SERVER_HOST")
	_ = viper.BindEnv("server.requestTimeoutSeconds", "METRICS_SERVER_TIMEOUT")
	_ = viper.BindEnv("server.disableCORS", "METRICS_SERVER_DISABLE_CORS")

	configuration := &Configuration{
		Server: ServerConfigurations{
			RequestTimeoutSeconds: DefaultRequestTimeoutSeconds,
		},
	}

	var err error
	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	err = viper.Unmarshal(configuration)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	validationError := configuration.validate()
	if validationError != nil {
		return nil, validationError
	}
	return configuration, nil
}

func (c *Configuration) validate() error {
	return nil
}
