package comunication

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
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

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
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
