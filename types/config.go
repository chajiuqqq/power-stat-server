package types

import "os"

type Config struct {
	BrokerIp string
	Port     string
	Profile  string
}

func NewConfig() *Config {
	return &Config{
		BrokerIp: os.Getenv("BROKER_IP"),
		Port:     os.Getenv("PORT"),
		Profile:  os.Getenv("PROFILE"),
	}
}
