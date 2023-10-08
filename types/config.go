package types

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	BrokerIp string
	Port     string
	Profile  string
	Cron     string
}

func NewConfig() *Config {

	res := &Config{
		BrokerIp: os.Getenv("BROKER_IP"),
		Port:     os.Getenv("PORT"),
		Profile:  os.Getenv("PROFILE"),
		Cron:     os.Getenv("CRON"),
	}
	if res.BrokerIp == "" {
		log.Fatal("Please set BROKER_IP")
	}
	if res.Port == "" {
		log.Fatal("Please set PORT")
	}
	if res.Profile == "" {
		log.Fatal("Please set PROFILE")
	}
	if res.Cron == "" {
		log.Fatal("Please set CRON")
	}
	return res
}
func (c *Config) IsDev() bool {
	return c.Profile == ProfileDev
}
func (c *Config) IsTest() bool {
	return c.Profile == ProfileTest
}
func (c *Config) DevGroup() string {
	devg := []string{"CaiJiaChen"}
	return strings.Join(devg, ",")
}
func (c *Config) ProdGroup() string {
	return "@all"
}
