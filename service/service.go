package service

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"mqtt-wx-forward/types"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	Mqtt           mqtt.Client
	Http           *http.Client
	teleSensorLogs []types.SensorData
	energyLogs     []types.EnergyTodayData
	logLensLimit   int
	Config         *types.Config
}
type ServiceOption struct {
	MqttBroker string
	ClientID   string
	Config     *types.Config
}

func New(opt *ServiceOption) *Service {
	if opt.Config == nil {
		opt.Config = types.NewConfig()
	}
	// MQTT连接参数
	opts := mqtt.NewClientOptions()
	opts.AddBroker(opt.MqttBroker)
	opts.SetClientID(opt.ClientID)
	if opt.ClientID == "" {
		opts.SetClientID("mqtt-client")
	}

	// 连接到MQTT代理
	mqttc := mqtt.NewClient(opts)
	if token := mqttc.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// http
	httpc := &http.Client{}
	return &Service{
		Mqtt:           mqttc,
		Http:           httpc,
		teleSensorLogs: make([]types.SensorData, 0),
		energyLogs:     make([]types.EnergyTodayData, 0),
		logLensLimit:   1000,
		Config:         opt.Config,
	}
}

func (s *Service) PushMsg(url string, d types.PushMsgData) error {
	payloadBytes, _ := json.Marshal(d)
	payload := strings.NewReader(string(payloadBytes))
	log.Println("post data:", string(payloadBytes))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.Http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("POST request failed with status: %s", resp.Status)
	}

	return nil
}

func (s *Service) SaveTeleLog(t types.SensorData) error {
	s.clearLogs()
	s.teleSensorLogs = append(s.teleSensorLogs, t)
	return nil
}
func (s *Service) SaveEnergyLog(t types.EnergyTodayData) error {
	s.clearLogs()
	s.energyLogs = append(s.energyLogs, t)
	return nil
}
func (s *Service) GetTopTeleMsg() *types.PushMsgData {
	lens := len(s.teleSensorLogs)
	if lens == 0 {
		return nil
	}
	d := s.teleSensorLogs[lens-1]
	t, err := time.Parse(d.Time, "2006-01-02T15:04:05")
	if err != nil {
		log.Println("fail to parse tele time:", d.Time)
		t = time.Now()
	}
	res := &types.PushMsgData{
		Title:       fmt.Sprintf("%s 电量统计", t.Format(time.DateOnly)),
		Description: fmt.Sprintf("总用电量%.2f度，今日用电%.2f度，昨日用电%.2f度", d.Energy.Total, d.Energy.Today, d.Energy.Yesterday),
		Content:     "",
		Channel:     "",
		Token:       "",
	}
	if s.Config.Profile == "" || s.Config.Profile == types.ProfileDev {
		res.To = "CaiJiaChen"
	} else {
		res.To = "@all"
	}
	return res
}
func (s *Service) GetTopEnergyMsg() *types.PushMsgData {
	lens := len(s.energyLogs)
	if lens == 0 {
		return nil
	}
	d := s.energyLogs[lens-1]
	res := &types.PushMsgData{
		Title:       fmt.Sprintf("%s 测试电量统计", d.Time.Format(time.DateOnly)),
		Description: fmt.Sprintf("总用电量%.2f度，今日用电%.2f度，昨日用电%.2f度", d.E.Total, d.E.Today, d.E.Yesterday),
		Content:     "",
		Channel:     "",
		Token:       "",
	}

	if s.Config.Profile == "" || s.Config.Profile == types.ProfileDev {
		res.To = "CaiJiaChen"
	} else {
		res.To = "@all"
	}
	return res
}
func (s *Service) clearLogs() {
	if len(s.teleSensorLogs) > s.logLensLimit {
		s.teleSensorLogs = make([]types.SensorData, 0)
	}
	if len(s.energyLogs) > s.logLensLimit {
		s.energyLogs = make([]types.EnergyTodayData, 0)
	}
}
