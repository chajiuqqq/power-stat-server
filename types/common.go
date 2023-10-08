package types

import "time"

const (
	TeleSensorTopic = "tele/generic/SENSOR"
	ResultTopic     = "stat/generic/RESULT"
	PushUrl         = "http://www.chajiuqqq.cn:3000/push/chajiuqqq"
	ProfileDev      = "dev"
	ProfileTest     = "test"
)

// 定义消息结构体
type SensorEnergy struct {
	TotalStartTime string
	Total          float64
	Yesterday      float64
	Today          float64
	Period         int64
	Power          int64
	ApparentPower  int64
	ReactivePower  int64
	Factor         float64
	Frequency      int64
	Voltage        int64
	Current        float64
}
type SensorData struct {
	Time   string       `json:"Time"`
	Energy SensorEnergy `json:"ENERGY"`
}

type EnergyToday struct {
	Total     float64
	Yesterday float64
	Today     float64
}

type EnergyTodayData struct {
	E    EnergyToday `json:"EnergyToday"`
	Time time.Time
}

type PushMsgData struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
	Channel     string `json:"channel,omitempty"`
	Token       string `json:"token,omitempty"`
	To          string `json:"to,omitempty"`
}
