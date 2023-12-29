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

func (p KafkaProducer) Publish(topic string, message string) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
	}

	err := p.producer.Produce(msg, nil)
	if err != nil {
		utility.ErrorLog().Printf("Failed to produce message to topic %s: %v", topic, err)
		return err
	}

	utility.InfoLog().Printf("Produced message to topic %s: %s", topic, message)
	return nil
}

type KafkaConsumer struct {
	consumer *kafka.Consumer
}

func NewKafkaConsumer(brokers, groupID, topic string) *KafkaConsumer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		utility.ErrorLog().Fatalf("Error creating consumer: %v", err)
	}

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		utility.ErrorLog().Fatalf("Error subscribing to topic %s: %v", topic, err)
	}

	return &KafkaConsumer{
		consumer: consumer,
	}
}

func (c KafkaConsumer) Close() {
	err := c.consumer.Close()
	if err != nil {
		return
	}
}

func (c KafkaConsumer) Consume() (string, error) {
	ev := c.consumer.Poll(10)
	switch e := ev.(type) {
	case *kafka.Message:
		_, err := c.consumer.CommitMessage(e)
		if err == nil {
			utility.InfoLog().Printf("Consumed message from topic %s: %s", *e.TopicPartition.Topic, string(e.Value))
			return string(e.Value), nil
		}
	case kafka.PartitionEOF:
		utility.InfoLog().Printf("Reached %v", e)
	case kafka.Error:
		return "", e
	default:
	}
	return "", nil

}
