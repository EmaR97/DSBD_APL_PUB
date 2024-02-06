import logging

from confluent_kafka import Producer


class MyProducer:
    producer: Producer = None

    def __init__(self, bootstrap_servers):
        self.bootstrap_servers = bootstrap_servers
        self.configure_producer()

    def configure_producer(self):
        producer_conf = {'bootstrap.servers': self.bootstrap_servers}
        self.producer = Producer(producer_conf)

    def produce_message(self, topic, message, key=None):
        # Produce message with optional key
        self.producer.produce(topic, value=message, key=key)
        self.producer.flush()
        logging.info(f'Sent on topic "{topic}" message: {message} with key: {key}')
