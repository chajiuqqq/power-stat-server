package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"mqtt-wx-forward/service"
	"mqtt-wx-forward/types"
	"time"
)

func main() {
	opt := &service.ServiceOption{
		MqttBroker: "tcp://192.168.1.18:1883",
		ClientID:   "mqtt-client",
	}
	sv := service.New(opt)

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
	}

	// 订阅消息
	if token := sv.Mqtt.Subscribe(types.TeleSensorTopic, 0, teleSensorHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if token := sv.Mqtt.Subscribe(types.ResultTopic, 0, energyTodayHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Println("Subscribe success")

	// 持续运行
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			tele := sv.GetTopTeleMsg()
			if tele != nil {
				err := sv.PushMsg(types.PushUrl, *tele)
				if err != nil {
					log.Println("fail to push tele msg:", err)
				}
			}
			energy := sv.GetTopEnergyMsg()
			if energy != nil {
				err := sv.PushMsg(types.PushUrl, *energy)
				if err != nil {
					log.Println("fail to push energy msg:", err)
				}
			}
			log.Println("push msg success")
		case <-time.Tick(1 * time.Second):
			log.Println(time.Now().Format(time.DateTime), " running...")
		default:

		}
	}
}
