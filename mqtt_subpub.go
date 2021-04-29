package yu

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MqttSubPub MQTT发布订阅
type MqttSubPub struct {
	clientID string
	client   mqtt.Client
	opts     *mqtt.ClientOptions
	filters  map[string]byte
	subDeal  func(mqtt.Message)
	mq       *QueueManager
}

type pubArgs struct {
	topic   string
	payload []byte
}

// NewMqttSubPub 创建MqttSubPub
func NewMqttSubPub(id, addr, user, pass string, topics ...string) *MqttSubPub {
	m := &MqttSubPub{
		clientID: id,
		mq:       NewQueueManager(),
	}
	m.mq.AddQueue(NewQueue(1024, 1, m.subscribe))
	m.mq.AddQueue(NewQueue(1024, 1, m.publish))
	m.setFliters(topics)
	m.opts = mqtt.NewClientOptions().
		AddBroker(addr).
		SetClientID(id).
		SetUsername(user).
		SetPassword(pass).
		SetCleanSession(false).
		SetAutoReconnect(false).
		SetKeepAlive(5 * time.Second).
		SetConnectTimeout(5 * time.Second).
		SetWriteTimeout(5 * time.Second).
		SetOnConnectHandler(m.onConnect).
		SetConnectionLostHandler(m.connectLost)
	m.client = mqtt.NewClient(m.opts)
	return m
}

// SetSubDeal 设置订阅处理回调
func (m *MqttSubPub) SetSubDeal(deal func(mqtt.Message)) {
	m.subDeal = deal
}

// OnStart 实现ServerInterface接口
func (m *MqttSubPub) OnStart() bool {
	// 启动队列
	m.mq.StartQueue()
	if err := m.connect(); err != nil {
		return false
	}
	return true
}

// OnStop 实现ServerInterface接口
func (m *MqttSubPub) OnStop() {
	if m.client == nil {
		return
	}
	if m.mq != nil {
		m.mq.StopQueue()
	}
	if len(m.filters) > 0 {
		m.client.Unsubscribe(m.getTopics()...)
	}
	if m.client.IsConnected() {
		m.client.Disconnect(250)
	}
	m.client = nil
}

// Info 实现ServerInterface接口
func (m *MqttSubPub) Info() string {
	return m.clientID
}

func (m *MqttSubPub) connect() (err error) {
	if m.client == nil {
		err = errors.New("client is nil")
		return
	}
	if m.client.IsConnected() {
		m.client.Disconnect(250)
	}
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		LogErr(token.Error(), m.clientID+" connect fail")
		return token.Error()
	}
	Logf("%s connect %s success", m.clientID, m.opts.Servers)
	return nil
}

func (m *MqttSubPub) getTopics() []string {
	if m.filters == nil {
		return []string{}
	}
	topics := make([]string, 0)
	for top := range m.filters {
		topics = append(topics, top)
	}
	return topics
}

func (m *MqttSubPub) setFliters(topics []string) {
	if m.filters == nil {
		m.filters = make(map[string]byte)
	}
	for _, top := range topics {
		m.filters[top] = 2
	}
}

func (m *MqttSubPub) onConnect(c mqtt.Client) {
	if token := c.SubscribeMultiple(m.filters, func(cli mqtt.Client, msg mqtt.Message) {
		// 将订阅消息送入处理对列
		m.mq.PushQueue(0, msg, true)
	}); token.Wait() && token.Error() != nil {
		LogErr(token.Error(), "onConnect SubscribeMultiple")
		return
	}
	Logf("%s subscribe topics: %v", m.clientID, m.getTopics())
}

func (m *MqttSubPub) connectLost(c mqtt.Client, err error) {
	LogErr(err, m.clientID+" lost")
	m.connect()
}

// Pub 发布消息到MQTT
func (m *MqttSubPub) Pub(topic string, payload []byte) {
	// 放到队列中执行
	m.mq.PushQueue(1, &pubArgs{topic: topic, payload: payload}, true)
}

func (m *MqttSubPub) subscribe(args interface{}) {
	msg, ok := args.(mqtt.Message)
	if !ok {
		Errf("subscribe args type is not mqtt.Message")
		return
	}
	if m.subDeal == nil {
		Errf("%s subDeal is nil", m.clientID)
		return
	}
	m.subDeal(msg)
}

func (m *MqttSubPub) publish(args interface{}) {
	pa, ok := args.(*pubArgs)
	if !ok {
		Errf("publish args type is not *pubArgs")
		return
	}
	if token := m.client.Publish(pa.topic, 2, false, pa.payload); token.Wait() && token.Error() != nil {
		LogErr(token.Error(), "pubMsg")
	}
}
