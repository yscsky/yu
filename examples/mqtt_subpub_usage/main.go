package main

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/yscsky/yu"
)

var (
	ginSvr *yu.GinServer
	mqsp1  *yu.MqttSubPub
	mqsp2  *yu.MqttSubPub
)

func main() {
	ginSvr = yu.NewGinServer("MqttSubPubUsage", ":8080", gin.DebugMode)
	mqsp1 = yu.NewMqttSubPub("mqsp1", ":1883", "", "", "mqttsubpub/1/#")
	mqsp2 = yu.NewMqttSubPub("mqsp2", ":1883", "", "", "mqttsubpub/2/#")
	yu.Run(yu.NewApp("MqttSubPubUsage", start, func() {}, ginSvr, mqsp1, mqsp2))
}

func start() bool {
	ginSvr.Engine.GET("talk", yu.NoCache(), yu.LogControl(true, []string{}), talk)
	mqsp1.SetSubDeal(deal1)
	mqsp2.SetSubDeal(deal2)
	return true
}

func talk(c *gin.Context) {
	mqsp1.Pub("mqttsubpub/2/0", []byte("hello, I'm mqsp1"))
	yu.JsonOK(c, nil)
}

func deal1(msg mqtt.Message) {
	yu.Logf("deal1 topic: %s, payload: %s", msg.Topic(), msg.Payload())
	time.Sleep(time.Second)
	mqsp1.Pub("mqttsubpub/2/0", []byte("hello, I'm mqsp1"))
}

func deal2(msg mqtt.Message) {
	yu.Logf("deal2 topic: %s, payload: %s", msg.Topic(), msg.Payload())
	time.Sleep(time.Second)
	mqsp2.Pub("mqttsubpub/1/0", []byte("hello, I'm mqsp2"))
}
