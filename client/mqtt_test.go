package client

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewMqttClient(t *testing.T) {
	cli := NewMqttClient(&MQTTConfig{
		ClientID:             "test",
		Addr:                 "tcp://127.0.0.1:1883",
		Username:             "admin",
		Password:             "public",
		AutoReconnect:        true,
		MaxReconnectInterval: time.Second * 5,
	}, log.DefaultLogger)
	cli.Start(context.Background())

	cliPub := NewMqttClient(&MQTTConfig{
		ClientID:             "public",
		Addr:                 "tcp://127.0.0.1:1883",
		Username:             "admin",
		Password:             "public",
		AutoReconnect:        true,
		MaxReconnectInterval: time.Second * 5,
	}, log.DefaultLogger)
	cliPub.Start(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)
	sub1 := func(client mqtt.Client, message mqtt.Message) {
		println("callback", message.Topic(), string(message.Payload()))
		assert.Equal(t, message.Payload(), []byte("test"))
		wg.Done()
	}

	cli.Subscribe("/test", 1, sub1)

	for i := 0; i < 10; i++ {
		cliPub.Client.Publish("/test", 1, false, "test")
	}
	wg.Wait()
	cli.Client.Unsubscribe("/test")

	cli.Stop(context.Background())
	cliPub.Stop(context.Background())
}
