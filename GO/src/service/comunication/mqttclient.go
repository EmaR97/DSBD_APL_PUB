package comunication

import (
	"CamMonitoring/src/utility"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type MqttClient struct {
	Client MQTT.Client
}

func NewMqttClient(brokerURL, clientID, username, password string) (*MqttClient, error) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetAutoReconnect(true)
	client := MQTT.NewClient(opts)
	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			utility.ErrorLog().Printf("Error initializing MQTT client: %v. Retring...", token.Error())
			time.Sleep(1 * time.Minute)
			//return nil, token.Error()
			continue
		}
		break
	}
	return &MqttClient{Client: client}, nil
}

func (m *MqttClient) PublishMessage(topic, payload string) error {
	token := m.Client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

func (m *MqttClient) Disconnect() {
	m.Client.Disconnect(250)
}
