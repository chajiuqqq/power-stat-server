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
	config := types.NewConfig()
	opt := &service.ServiceOption{
		MqttBroker: fmt.Sprintf("tcp://%s:%s", config.BrokerIp, config.Port),
		ClientID:   "mqtt-client",
		Config:     config,
	}
	sv := service.New(opt)
	log.Println("Connect to broker successfully:", opt.MqttBroker)

	// 定义消息处理函数
	teleSensorHandler := func(client mqtt.Client, msg mqtt.Message) {
		// 解码消息为结构体
		var d types.SensorData
		err := json.Unmarshal(msg.Payload(), &d)
		if err != nil {
			log.Println("Failed to decode message:", err)
			return
		}

		log.Println("receive: ", string(msg.Payload()))

		err = sv.SaveTeleLog(d)
		if err != nil {
			log.Println("fail to save tele log:", err, "raw:", msg.Payload())
			return
		}
	}
	// 定义消息处理函数
	energyTodayHandler := func(client mqtt.Client, msg mqtt.Message) {
		log.Println("receive: ", string(msg.Payload()))
		// 解码消息为结构体
		var d types.EnergyTodayData
		err := json.Unmarshal(msg.Payload(), &d)
		if err != nil {
			log.Println("Failed to decode message:", err)
		}
		d.Time = time.Now()
		err = sv.SaveEnergyLog(d)
		if err != nil {
			log.Println("fail to save energy log:", err, "raw:", msg.Payload())
			return
		}

		energy := sv.GetTopEnergyMsg()
		if energy != nil {
			err := sv.PushMsg(types.PushUrl, *energy)
			if err != nil {
				log.Println("fail to push energy msg:", err)
			} else {
				log.Println("push energy msg success")
			}
		}
	}

	// 订阅消息
	if token := sv.Mqtt.Subscribe(types.TeleSensorTopic, 0, teleSensorHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if token := sv.Mqtt.Subscribe(types.ResultTopic, 0, energyTodayHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Println("Subscribe success")

	c := cron.New()
	c.AddFunc("0 0 9 * * *", func() {
		tele := sv.GetTopTeleMsg()
		if tele != nil {
			err := sv.PushMsg(types.PushUrl, *tele)
			if err != nil {
				log.Println("fail to push tele msg:", err)
			} else {
				log.Println("push tele msg success")
			}
		}
	})
	c.Start()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	c.Stop()
	log.Println("server scheduler stopped")
}
