package comunication

import (
	"CamMonitoring/src/utility"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
)

type KafkaConsumer struct {
	consumer       *kafka.Consumer
	topicCallbacks map[string]func(string2 string) error
	stopChan       chan struct{}
	wg             sync.WaitGroup
}

func NewKafkaConsumer(brokers, groupID string, topicCallbacks map[string]func(string2 string) error) *KafkaConsumer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		utility.ErrorLog().Fatalf("Error creating consumer: %v", err)
	}

	// Subscribe to topics
	topics := make([]string, 0, len(topicCallbacks))
	for topic := range topicCallbacks {
		topics = append(topics, topic)
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		utility.ErrorLog().Fatalf("Error subscribing to topics: %v", err)
	}

	return &KafkaConsumer{
		consumer:       consumer,
		topicCallbacks: topicCallbacks,
		stopChan:       make(chan struct{}),
	}
}

func (c *KafkaConsumer) Close() {
	err := c.consumer.Close()
	if err != nil {
		return
	}
}

func (c *KafkaConsumer) Start() {
	utility.Start(c.stopChan, &c.wg, c.Consume, "Consuming process stopped")
}
func (c *KafkaConsumer) Stop() {
	utility.Stop(c.stopChan, &c.wg)
}
func (c *KafkaConsumer) Consume() {
	ev := c.consumer.Poll(10)
	switch e := ev.(type) {
	case *kafka.Message:
		callback, exists := c.topicCallbacks[*e.TopicPartition.Topic]
		if !exists || callback == nil {
			utility.InfoLog().Printf("No callback found for topic: %s", *e.TopicPartition.Topic)
		}
		err := callback(string(e.Value))
		if err != nil {
			// Log an error if there is an issue processing and storing the frame
			utility.ErrorLog().Println("Error processing message:", err)
			// Consider whether to continue processing or return an error to the caller
		}
	case kafka.PartitionEOF:
		utility.InfoLog().Printf("Reached %v", e)
	case kafka.Error:
		utility.ErrorLog().Printf("Error consuming message:  %v", e)
	default:
	}
}
