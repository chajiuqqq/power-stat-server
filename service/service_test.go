package service

import (
	"context"
	"fmt"
	"gotest.tools/v3/assert"
	"log"
	"mqtt-wx-forward/model/report"
	"mqtt-wx-forward/types"
	"testing"
	"time"
)

var sv *Service

func TestMain(m *testing.M) {
	conf := types.NewConfig()
	opt := &ServiceOption{
		MqttBroker: fmt.Sprintf("tcp://%s:%s", conf.BrokerIp, conf.Port),
		ClientID:   "",
	}
	logger := log.Default()
	sv = New(conf, logger, opt)
	m.Run()
	log.Println("test end")
}

func Test_config(t *testing.T) {
	assert.Equal(t, sv.Config.BrokerIp, "100.127.185.26")
	assert.Equal(t, sv.Config.Port, "1883")
	assert.Equal(t, sv.Config.Cron, "0 52 12 * * *")
	assert.Equal(t, sv.Config.Profile, "dev")
}

func TestService_SaveTeleLog(t *testing.T) {
	ctx := context.Background()
	err := sv.SaveTeleLog(ctx, types.SensorData{Time: "2023-10-08T09:21:03", Energy: types.SensorEnergy{
		TotalStartTime: "2023-10-08T09:21:03",
		Total:          100.001,
	}})
	assert.NilError(t, err)
	res, err := sv.ReportModel.FindLatest(ctx)
	assert.NilError(t, err)
	assert.Equal(t, res.Total, float64(100.001))
}

func TestModel(t *testing.T) {
	ctx := context.Background()
	res, err := sv.ReportModel.Insert(ctx, &report.Report{
		Time:           time.Now(),
		TotalStartTime: time.Now(),
		Total:          100.001,
		Yesterday:      90.001,
		Today:          10,
		Period:         0,
		Power:          0,
		ApparentPower:  0,
		ReactivePower:  0,
		Factor:         0,
		Frequency:      0,
		Voltage:        0,
		Current:        0,
	})
	assert.NilError(t, err)
	id, err := res.LastInsertId()
	assert.NilError(t, err)
	r, err := sv.ReportModel.FindOne(ctx, id)
	assert.NilError(t, err)
	assert.Equal(t, r.Id, id)
	assert.Equal(t, r.Total, float64(100.001))
	assert.Equal(t, r.Yesterday, float64(90.001))
	assert.Equal(t, r.Today, float64(10))
}
