package comunication

import (
	"CamMonitoring/src/utility"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(brokers string) *KafkaProducer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		utility.ErrorLog().Fatalf("Error creating producer: %v", err)
	}
	return &KafkaProducer{
		producer: producer,
	}
}
func (p KafkaProducer) Close() {
	p.producer.Close()

}

func (p KafkaProducer) Publish(topic string, message string, key string) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Key:   []byte(key),
	}

	err := p.producer.Produce(msg, nil)
	if err != nil {
		utility.ErrorLog().Printf("Failed to produce message to topic %s: %v", topic, err)
		return err
	}

	utility.InfoLog().Printf("Produced message to topic %s: %s", topic, message)
	return nil
}
