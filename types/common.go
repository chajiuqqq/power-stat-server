package types

// 定义消息结构体
type SensorEnergy struct {
	TotalStartTime string
	Total          float64
	Yesterday      float64
	Today          float64
	Power          int64
	ApparentPower  int64
	ReactivePower  int64
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
	E EnergyToday `json:"EnergyToday"`
}

type PushMsgData struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
	Channel     string `json:"channel,omitempty"`
	Token       string `json:"token,omitempty"`
	To          string `json:"to,omitempty"`
}
