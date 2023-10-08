package main

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"mqtt-wx-forward/types"
)

func teleSensorHandler(ctx context.Context) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		// 解码消息为结构体
		var d types.SensorData
		err := json.Unmarshal(msg.Payload(), &d)
		if err != nil {
			sv.Logger.Println("Failed to decode message:", err)
			return
		}

		sv.Logger.Println("receive: ", string(msg.Payload()))

		err = sv.SaveTeleLog(ctx, d)
		if err != nil {
			sv.Logger.Println("fail to save tele log:", err, "raw:", string(msg.Payload()))
			return
		}
	}
}
func energyTodayHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		// sv.Logger.Println("receive: ", string(msg.Payload()))
		// // 解码消息为结构体
		// var d types.EnergyTodayData
		// err := json.Unmarshal(msg.Payload(), &d)
		// if err != nil {
		// 	sv.Logger.Println("Failed to decode message:", err)
		// }
		// d.Time = time.Now()
		// err = sv.SaveEnergyLog(d)
		// if err != nil {
		// 	sv.Logger.Println("fail to save energy log:", err, "raw:", msg.Payload())
		// 	return
		// }
		//
		// energy := sv.GetTopEnergyMsg()
		// if energy != nil {
		// 	err := sv.PushMsg(types.PushUrl, *energy)
		// 	if err != nil {
		// 		sv.Logger.Println("fail to push energy msg:", err)
		// 	} else {
		// 		sv.Logger.Println("push energy msg success")
		// 	}
		// }
	}
}
