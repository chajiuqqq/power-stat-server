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
	Wx       WxConfig
}

type WxConfig struct {
	Token          string `json:"token"`
	ReceiverId     string `json:"receiverId"`
	EncodingAeskey string `json:"encodingAesKey"`
}

func NewConfig() *Config {

	res := &Config{
		BrokerIp: os.Getenv("BROKER_IP"),
		Port:     os.Getenv("PORT"),
		Profile:  os.Getenv("PROFILE"),
		Cron:     os.Getenv("CRON"),
		Wx: WxConfig{
			Token:          os.Getenv("WX_TOKEN"),
			ReceiverId:     os.Getenv("WX_RECEIVER_ID"),
			EncodingAeskey: os.Getenv("WX_ENCODING_AES_KEY"),
		},
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
	if res.Wx.Token == "" {
		log.Fatal("Please set WX_TOKEN")
	}
	if res.Wx.ReceiverId == "" {
		log.Fatal("Please set WX_RECEIVER_ID")
	}
	if res.Wx.EncodingAeskey == "" {
		log.Fatal("Please set WX_ENCODING_AES_KEY")
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
