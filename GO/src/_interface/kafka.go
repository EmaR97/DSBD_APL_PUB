package _interface

type KafkaProducer interface {
	Close()
	Publish(
		topic string,
		message string,
	) error
}
