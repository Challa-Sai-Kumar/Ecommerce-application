package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Kafka    Kafka    `mapstructure:"kafka"`
	Email    Email    `mapstructure:"email"`
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
	Logging  Logging  `mapstructure:"logging"`
}

type Kafka struct {
	BrokerList     []string          `mapstructure:"broker_list"`
	Topics         map[string]string `mapstructure:"topics"`
	ConsumerGroups map[string]string `mapstructure:"consumer_groups"`
}

type Email struct {
	SMTPHost    string `mapstructure:"smtp_host"`
	SMTPPort    int    `mapstructure:"smtp_port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	FromAddress string `mapstructure:"from_address"`
}

type Database struct {
	DataSourceName     string `mapstructure:"data_source_name"`
	DriverName         string `mapstructure:"driver_name"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
}

type Server struct {
	Port int `mapstructure:"port"`
}

type Logging struct {
	Level string `mapstructure:"level"`
}

func loadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
