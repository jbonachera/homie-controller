package messaging

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"time"
)

type CallbackHandler func(message homieMessage.HomieMessage)

type messagingBroker struct {
	c         MQTT.Client
	baseTopic string
	connected bool
}

var broker messagingBroker
var publishChan chan mqttMessage
var subChan chan subMessage
var unsubChan chan string
var done chan bool
var routines int

type mqttMessage struct {
	Topic   string
	Qos     byte
	Retain  bool
	Payload interface{}
}

type subMessage struct {
	Topic    string
	Qos      byte
	Callback MQTT.MessageHandler
}

func Stop() {
	for i := 0; i < routines; i++ {
		done <- true
	}
}

func Start(brokerHost string, client_id string, mqttBaseTopic string) {
	publishChan = make(chan mqttMessage, 20)
	subChan = make(chan subMessage, 20)
	unsubChan = make(chan string, 20)
	done = make(chan bool, 1)
	routines = 0

	opts := MQTT.NewClientOptions().AddBroker("tcp://" + brokerHost + ":1883")
	opts.SetClientID(client_id)
	broker = messagingBroker{MQTT.NewClient(opts), mqttBaseTopic, false}
	for !broker.connected {
		if token := broker.c.Connect(); token.Wait() && token.Error() != nil {
			log.Error("could not connect to MQTT.")
			time.Sleep(5 * time.Second)
		} else {
			log.Info("connected to MQTT broker")
			broker.connected = true
		}
	}
	go unqueuePublish()
	routines++
	go unqueueSub()
	routines++
	go unqueueUnsub()
	routines++

	routines++

}

func mqttPublish(msg mqttMessage) {
	publishChan <- msg
}
func PublishFile(topic string, payload interface{}) {
	mqttPublish(mqttMessage{topic, 1, false, payload})
}
func PublishMessage(message homieMessage.HomieMessage) {
	mqttPublish(mqttMessage{message.Topic, 1, false, message.Payload})
}
func PublishState(message homieMessage.HomieMessage) {
	mqttPublish(mqttMessage{message.Topic, 1, true, message.Payload})
}

func unqueueSub() {
	log.Debug("Starting MQTT subscribing queue")
	run := true
	for run {
		select {
		case msg := <-subChan:
			broker.c.Subscribe(msg.Topic, msg.Qos, msg.Callback)
		case <-done:
			run = false
			break
		}
	}
	log.Debug("Finished MQTT subscribing queue")
}

func unqueueUnsub() {
	log.Debug("Starting MQTT unsubscribing queue")
	run := true
	for run {
		select {
		case msg := <-unsubChan:
			broker.c.Unsubscribe(msg)
		case <-done:
			run = false
			break
		}
	}
	log.Debug("Finished MQTT unsubscribing queue")

}

func unqueuePublish() {
	log.Debug("Starting MQTT publishing queue")
	run := true
	for run {
		select {
		case msg := <-publishChan:
			log.Info("publishing to " + msg.Topic)
			broker.c.Publish(msg.Topic, msg.Qos, msg.Retain, msg.Payload)
		case <-done:
			run = false
			break
		}
	}
	log.Debug("Finished MQTT publishing queue")

}
func DelSubscription(topic string) {
	log.Debug("Unsubscribing to " + topic)
	unsubChan <- topic
}
func AddSubscription(topic string, qos byte, callback MQTT.MessageHandler) {
	log.Debug("Subscribing to " + topic)
	subChan <- subMessage{Topic: topic, Qos: qos, Callback: callback}
}

func AddHandler(topic string, callback CallbackHandler) {
	AddSubscription(topic, 0, func(mqttClient MQTT.Client, mqttMessage MQTT.Message) {
		message, err := homieMessage.Extract(mqttMessage, broker.baseTopic)
		if err != nil {
			log.Error("malformed message")
		} else {
			callback(message)
		}
	})
}
