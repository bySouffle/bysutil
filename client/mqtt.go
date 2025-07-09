package client

import (
	//"InspectionRobot/app/task_core_module/internal/conf"
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/container/gmap"
	"sync"
	"time"
)

type MQTTConfig struct {
	Network  string
	Addr     string
	ClientID string
	Username string
	Password string

	AutoReconnect        bool
	MaxReconnectInterval time.Duration
}

type PayloadFunc struct {
	Qos      int
	function mqtt.MessageHandler
}

type Mqtt struct {
	Client    mqtt.Client
	log       *log.Helper
	mutex     sync.Mutex
	topicList *gmap.Map
}

func (m *Mqtt) DefaultReceiveCallBack(client mqtt.Client, msg mqtt.Message) {
	//m.log.Infof("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
func (m *Mqtt) DefaultConnectCallBack(client mqtt.Client) {
	rOps := client.OptionsReader()
	cid := rOps.ClientID()
	m.log.Infof("Connected: %v", cid)
	m.reSubscribe()
	topics := m.topicList.Keys()
	for _, v := range topics {
		m.log.Infof("Sub: %v: %v", v, m.topicList.Get(v))
	}
}
func (m *Mqtt) DefaultCloseCallBack(client mqtt.Client, err error) {
	m.log.Infof("Connect lost: %v", err)
}

func (m *Mqtt) DefaultReConnectsCallBack(client mqtt.Client, err error) {
	m.log.Infof("Connect lost: %v, reconnect...", err)
}

func NewMqttClient(server *MQTTConfig, logger log.Logger) *Mqtt {
	mqClient := Mqtt{
		log:       log.NewHelper(logger),
		topicList: gmap.New(true),
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server.Addr)
	opts.SetClientID(server.ClientID)
	opts.SetUsername(server.Username)
	opts.SetPassword(server.Password)
	opts.SetDefaultPublishHandler(mqClient.DefaultReceiveCallBack)
	opts.OnConnect = mqClient.DefaultConnectCallBack
	opts.OnConnectionLost = mqClient.DefaultCloseCallBack

	opts.SetConnectionLostHandler(mqClient.DefaultReConnectsCallBack)
	opts.SetAutoReconnect(server.AutoReconnect)               // 启用自动重连
	opts.SetMaxReconnectInterval(server.MaxReconnectInterval) // 设置最大重连间隔
	//opts.SetCleanSession(false)
	mqClient.Client = mqtt.NewClient(opts)
	if token := mqClient.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return &mqClient
}

func (m *Mqtt) Start(ctx context.Context) error {
	rOps := m.Client.OptionsReader()
	s := rOps.Servers()
	log.Infof("[mqtt] client listening on: %v", s)
	return nil
}

func (m *Mqtt) Stop(ctx context.Context) error {
	log.Info("[mqtt] client stopping")
	topics := m.topicList.Keys()
	for _, topic := range topics {
		m.Client.Unsubscribe(topic.(string))
	}
	m.Client.Disconnect(250)
	return nil
}

func (m *Mqtt) Subscribe(topic string, qos int, callback mqtt.MessageHandler) error {
	if token := m.Client.Subscribe(topic, byte(qos), callback); token.Wait() && token.Error() != nil {
		m.log.Error(token.Error())
		return token.Error()
	}
	log.Infof("[mqtt] client subscribe: %v\t[qos: %v]", topic, qos)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	payloadFunction := PayloadFunc{
		Qos:      qos,
		function: callback,
	}
	m.topicList.Set(topic, payloadFunction)
	return nil
}

func (m *Mqtt) Unsubscribe(topic string) error {
	if token := m.Client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		m.log.Error(token.Error())
	}
	log.Infof("[mqtt] client unsubscribe: %v", topic)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.topicList.Remove(topic)

	return nil
}

func (m *Mqtt) reSubscribe() {
	topics := m.topicList.Keys()
	for _, topic := range topics {
		payloadFunction := m.topicList.Get(topic).(PayloadFunc)
		if token := m.Client.Subscribe(topic.(string), byte(payloadFunction.Qos), payloadFunction.function); token.Wait() && token.Error() != nil {
			m.log.Error(token.Error())
		}
	}
}
