package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"mqtt-wx-forward/service"
	"mqtt-wx-forward/types"
	"os"
	"os/signal"
	"syscall"
)

var sv *service.Service

func main() {
	conf := types.NewConfig()
	opt := &service.ServiceOption{
		MqttBroker: fmt.Sprintf("tcp://%s:%s", conf.BrokerIp, conf.Port),
		ClientID:   "",
	}
	logger := log.Default()
	sv = service.New(conf, logger, opt)
	ctx := context.Background()

	// 订阅消息
	if token := sv.Mqtt.Subscribe(types.TeleSensorTopic, 0, teleSensorHandler(ctx)); token.Wait() && token.Error() != nil {
		logger.Fatal(token.Error())
	}

	if token := sv.Mqtt.Subscribe(types.ResultTopic, 0, energyTodayHandler()); token.Wait() && token.Error() != nil {
		logger.Fatal(token.Error())
	}
	logger.Printf("parser started, conf: %+v", conf)

	// cron
	cronOpts := cron.WithSeconds()
	c := cron.New(cronOpts)
	_, err := c.AddFunc(conf.Cron, func() {
		tele, err := sv.GetTopTeleMsg(ctx)
		if err != nil {
			logger.Println("fail to get top tele msg:", err)
			return
		}
		err = sv.PushMsg(types.PushUrl, *tele)
		if err != nil {
			logger.Println("fail to push tele msg:", err)
		} else {
			logger.Println("push tele msg success")
		}
	})
	if err != nil {
		logger.Fatal(err)
	}
	c.Start()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	c.Stop()
	logger.Println("parser stopped")
}
