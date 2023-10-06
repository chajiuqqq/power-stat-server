package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"mqtt-wx-forward/types"
)

func main() {
	// MQTT连接参数
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://192.168.1.18:1883")
	opts.SetClientID("mqtt-client")

	// 连接到MQTT代理
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

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

		pushMsg := types.PushMsgData{
			Title:       "电量统计",
			Description: fmt.Sprintf("%s 总用电量%.2f度，今日用电%.2f度，昨日用电%.2f度", d.Time, d.Energy.Total, d.Energy.Today, d.Energy.Yesterday),
			Content:     "",
			Channel:     "",
			Token:       "",
			//To:          "@all",
			To: "CaiJiaChen",
		}

		// 发送POST请求
		err = sendPostRequest("http://www.chajiuqqq.cn:3000/push/chajiuqqq", pushMsg)
		if err != nil {
			log.Println("Failed to send POST request:", err)
			return
		}

		log.Println("POST request sent successfully!")
	}
	// 定义消息处理函数
	energyTodayHandler := func(client mqtt.Client, msg mqtt.Message) {
		log.Println("receive: ", string(msg.Payload()))
		// 解码消息为结构体
		var d = new(types.EnergyTodayData)
		var pushMsg = &types.PushMsgData{
			Title:       "独立电量统计",
			Description: string(msg.Payload()),
			Content:     "",
			Channel:     "",
			Token:       "",
			//To:          "@all",
			To: "CaiJiaChen",
		}
		err := json.Unmarshal(msg.Payload(), d)
		if err != nil {
			log.Println("Failed to decode message:", err, "send raw json:", string(msg.Payload()))
		} else {
			pushMsg.Description = fmt.Sprintf("总用电量%.2f度，今日用电%.2f度，昨日用电%.2f度", d.E.Total, d.E.Today, d.E.Yesterday)
		}

		// 发送POST请求
		err = sendPostRequest("http://www.chajiuqqq.cn:3000/push/chajiuqqq", *pushMsg)
		if err != nil {
			log.Println("Failed to send POST request:", err)
			return
		}

		log.Println("POST request sent successfully!")
	}

	// 订阅消息
	if token := client.Subscribe("tele/generic/SENSOR", 0, teleSensorHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if token := client.Subscribe("stat/generic/RESULT", 0, energyTodayHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Println("Subscribe success")

	// 持续运行
	for {
		time.Sleep(1 * time.Second)
	}
}

// 发送POST请求
func sendPostRequest(url string, d types.PushMsgData) error {
	payloadBytes, _ := json.Marshal(d)
	payload := strings.NewReader(string(payloadBytes))
	log.Println("post data:", string(payloadBytes))
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("POST request failed with status: %s", resp.Status)
	}

	return nil
}
