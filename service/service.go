package service

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"log"
	"mqtt-wx-forward/model/report"
	"mqtt-wx-forward/tools/wxbizmsgcrypt"
	"mqtt-wx-forward/types"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	Mqtt        mqtt.Client
	Http        *http.Client
	Config      *types.Config
	Logger      *log.Logger
	ReportModel report.ReportModel
	wxcpt       *wxbizmsgcrypt.WXBizMsgCrypt
}
type ServiceOption struct {
	MqttBroker string
	ClientID   string
}

func New(conf *types.Config, logger *log.Logger, opt *ServiceOption) *Service {
	if conf == nil {
		conf = types.NewConfig()
	}
	// MQTT连接参数
	if opt.ClientID == "" {
		opt.ClientID = "mqtt-client-" + time.Now().Format("20060102150405")
	}
	var mqttc mqtt.Client
	if opt.MqttBroker != "" {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(opt.MqttBroker)
		opts.SetClientID(opt.ClientID)

		// 连接到MQTT代理
		mqttc = mqtt.NewClient(opts)
		if conf.Profile != types.ProfileTest {
			if token := mqttc.Connect(); token.Wait() && token.Error() != nil {
				logger.Fatal(token.Error())
			}
			logger.Println("Connect to broker successfully:", opt.MqttBroker)
		} else {
			logger.Println("Testing env,skip connect to mqtt broker")
		}
	}

	// http
	httpc := &http.Client{}

	// mysql dsn
	conn := sqlx.NewMysql(fmt.Sprintf("root:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.DBPass, conf.DBHost, conf.DB))
	return &Service{
		Mqtt:        mqttc,
		Http:        httpc,
		Config:      conf,
		Logger:      logger,
		ReportModel: report.NewReportModel(conn),
		wxcpt:       wxbizmsgcrypt.NewWXBizMsgCrypt(conf.Wx.Token, conf.Wx.EncodingAeskey, conf.Wx.ReceiverId, wxbizmsgcrypt.XmlType),
	}
}

func (s *Service) PushMsg(url string, d types.PushMsgData) error {
	payloadBytes, _ := json.Marshal(d)
	payload := strings.NewReader(string(payloadBytes))
	s.Logger.Println("push data:", string(payloadBytes))

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

func (s *Service) SaveTeleLog(ctx context.Context, t types.SensorData) error {
	myTime, err := time.Parse("2006-01-02T15:04:05", t.Time)
	if err != nil {
		return err
	}
	startTime, err := time.Parse("2006-01-02T15:04:05", t.Energy.TotalStartTime)
	if err != nil {
		return err
	}
	_, err = s.ReportModel.Insert(ctx, &report.Report{
		Time:           myTime,
		TotalStartTime: startTime,
		Total:          t.Energy.Total,
		Yesterday:      t.Energy.Yesterday,
		Today:          t.Energy.Today,
		Period:         t.Energy.Period,
		Power:          t.Energy.Power,
		ApparentPower:  t.Energy.ApparentPower,
		ReactivePower:  t.Energy.ReactivePower,
		Factor:         t.Energy.Factor,
		Frequency:      t.Energy.Frequency,
		Voltage:        t.Energy.Voltage,
		Current:        t.Energy.Current,
	})
	return err
}
func (s *Service) SaveEnergyLog(t types.EnergyTodayData) error {

	return nil
}
func (s *Service) GetTopTeleMsg(ctx context.Context) (*types.PushMsgData, error) {
	d, err := s.ReportModel.FindLatest(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetTopTeleMsg error:%w", err)
	}
	res := &types.PushMsgData{
		Title:       fmt.Sprintf("%s 电量统计", d.Time.Format(time.DateOnly)),
		Description: fmt.Sprintf("总用电量%.2f度，今日用电%.2f度，昨日用电%.2f度", d.Total, d.Today, d.Yesterday),
		Content:     "",
		Channel:     "",
		Token:       "",
	}
	if s.Config.IsDev() {
		res.To = s.Config.DevGroup()
	} else {
		res.To = s.Config.ProdGroup()
	}
	return res, nil
}

// func (s *Service) GetTopEnergyMsg() *types.PushMsgData {
// 	lens := len(s.energyLogs)
// 	if lens == 0 {
// 		return nil
// 	}
// 	d := s.energyLogs[lens-1]
// 	res := &types.PushMsgData{
// 		Title:       fmt.Sprintf("%s 测试电量统计", d.Time.Format(time.DateOnly)),
// 		Description: fmt.Sprintf("总用电量%.2f度，今日用电%.2f度，昨日用电%.2f度", d.E.Total, d.E.Today, d.E.Yesterday),
// 		Content:     "",
// 		Channel:     "",
// 		Token:       "",
// 	}
//
// 	if s.Config.Profile == "" || s.Config.Profile == types.ProfileDev {
// 		res.To = "CaiJiaChen"
// 	} else {
// 		res.To = "@all"
// 	}
// 	return res
// }
