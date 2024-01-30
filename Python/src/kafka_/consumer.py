import asyncio
import logging

from confluent_kafka import Consumer, KafkaError


class NotificationConsumer:
    consumer: Consumer = None

    def __init__(self, bootstrap_servers, group_id, auto_offset_reset, topic, message_handler):
        self.bootstrap_servers = bootstrap_servers
        self.group_id = group_id
        self.auto_offset_reset = auto_offset_reset
        self.topic = topic
        self.message_handler = message_handler
        self.received_messages = []

    def configure_consumer(self):
        consumer_conf = {'bootstrap.servers': self.bootstrap_servers, 'group.id': self.group_id,
                         'auto.offset.reset': self.auto_offset_reset}
        self.consumer = Consumer(consumer_conf)
        self.consumer.subscribe([self.topic])

    def consume_message(self):
        try:
            while True:
                msg = self.consumer.consume(1, timeout=1)

                if len(msg) == 0:
                    logging.debug('timeout topic: %s', self.topic)
                    continue
                msg = msg[0]
                if msg.error():
                    if msg.error().code() == KafkaError._PARTITION_EOF:
                        continue
                    else:
                        logging.error('Error: {}'.format(msg.error()))
                        break
                # Execute the provided message handler function
                logging.debug('Message: %s', msg.value())
                asyncio.run(self.message_handler(msg))

        except KeyboardInterrupt:
            pass

        finally:
            self.consumer.close()  # await self.print_received_messages()
