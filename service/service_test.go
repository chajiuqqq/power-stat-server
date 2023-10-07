package service

import (
	"fmt"
	"gotest.tools/v3/assert"
	"log"
	"mqtt-wx-forward/types"
	"strconv"
	"testing"
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
	for i := 1; i <= 100; i++ {
		sv.SaveTeleLog(types.SensorData{Time: strconv.Itoa(i)})
	}
	assert.Equal(t, len(sv.teleSensorLogs), 100)
	assert.Equal(t, sv.teleSensorLogs[len(sv.teleSensorLogs)-1].Time, "100")
	sv.SaveTeleLog(types.SensorData{Time: "101"})
	assert.Equal(t, len(sv.teleSensorLogs), 1)
	assert.Equal(t, sv.teleSensorLogs[0].Time, "101")
}
