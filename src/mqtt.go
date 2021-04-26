package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
)

const (
	broker = "tcp://172.31.120.1:1883"
	qos    = 0
)

var client MQTT.Client

func MQTTInit() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(dogName)
	opts.SetCleanSession(false)
	//opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
	//	mqttChannel <- [2]string{msg.Topic(), string(msg.Payload())}
	//})
	client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func publishMQTT(topic string, payload string) {
	INFO.Print("Publish MQTT Message:")
	INFO.Printf("\ttopic: %s\n", topic)
	INFO.Printf("\tmessage: %s\n", payload)
	token := client.Publish(topic, byte(qos), false, payload)
	token.Wait()
	INFO.Print("Publish Finished")
}

func subscribeMQTT(topic string) {
	if token := client.Subscribe(topic, byte(qos), func(client MQTT.Client, msg MQTT.Message) {
		mqttChannel <- [2]string{msg.Topic(), string(msg.Payload())}
	}); token.Wait() && token.Error() != nil {
		CRITICAL.Print(token.Error())
		os.Exit(1)
	}
	INFO.Printf("Subscribe Topic: %s", topic)
}
