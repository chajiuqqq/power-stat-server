package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/robfig/cron/v3"
	"log"
	"mqtt-wx-forward/service"
	"mqtt-wx-forward/types"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf := types.NewConfig()
	opt := &service.ServiceOption{
		MqttBroker: fmt.Sprintf("tcp://%s:%s", conf.BrokerIp, conf.Port),
		ClientID:   "",
	}
	logger := log.Default()
	sv := service.New(conf, logger, opt)

	// 定义消息处理函数
	teleSensorHandler := func(client mqtt.Client, msg mqtt.Message) {
		// 解码消息为结构体
		var d types.SensorData
		err := json.Unmarshal(msg.Payload(), &d)
		if err != nil {
			logger.Println("Failed to decode message:", err)
			return
		}

		logger.Println("receive: ", string(msg.Payload()))

		err = sv.SaveTeleLog(d)
		if err != nil {
			logger.Println("fail to save tele log:", err, "raw:", msg.Payload())
			return
		}
	}
	// 定义消息处理函数
	energyTodayHandler := func(client mqtt.Client, msg mqtt.Message) {
		logger.Println("receive: ", string(msg.Payload()))
		// 解码消息为结构体
		var d types.EnergyTodayData
		err := json.Unmarshal(msg.Payload(), &d)
		if err != nil {
			logger.Println("Failed to decode message:", err)
		}
		d.Time = time.Now()
		err = sv.SaveEnergyLog(d)
		if err != nil {
			logger.Println("fail to save energy log:", err, "raw:", msg.Payload())
			return
		}

		energy := sv.GetTopEnergyMsg()
		if energy != nil {
			err := sv.PushMsg(types.PushUrl, *energy)
			if err != nil {
				logger.Println("fail to push energy msg:", err)
			} else {
				logger.Println("push energy msg success")
			}
		}
	}

	// 订阅消息
	if token := sv.Mqtt.Subscribe(types.TeleSensorTopic, 0, teleSensorHandler); token.Wait() && token.Error() != nil {
		logger.Fatal(token.Error())
	}

	if token := sv.Mqtt.Subscribe(types.ResultTopic, 0, energyTodayHandler); token.Wait() && token.Error() != nil {
		logger.Fatal(token.Error())
	}
	logger.Println("Subscribe success")

	cronOpts := cron.WithSeconds()
	c := cron.New(cronOpts)
	_, err := c.AddFunc(conf.Cron, func() {
		tele := sv.GetTopTeleMsg()
		if tele != nil {
			err := sv.PushMsg(types.PushUrl, *tele)
			if err != nil {
				logger.Println("fail to push tele msg:", err)
			} else {
				logger.Println("push tele msg success")
			}
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
	logger.Println("server scheduler stopped")
}
