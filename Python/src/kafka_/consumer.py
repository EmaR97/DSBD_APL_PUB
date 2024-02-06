import logging
import time

from confluent_kafka import Consumer, KafkaError


class MyConsumer:
    consumer: Consumer = None

    def __init__(self, bootstrap_servers, topic_handlers, group_id, auto_offset_reset):
        self.bootstrap_servers = bootstrap_servers
        self.group_id = group_id
        self.auto_offset_reset = auto_offset_reset
        self.topic_handlers = topic_handlers
        self.received_messages = []
        self.configure_consumer()
        self.stop_requested = False

    def configure_consumer(self):
        consumer_conf = {'bootstrap.servers': self.bootstrap_servers, 'group.id': self.group_id,
                         'auto.offset.reset': self.auto_offset_reset}
        self.consumer = Consumer(consumer_conf)
        topics = list(self.topic_handlers.keys())
        logging.info('Listening on topics ' + str(topics))

        self.consumer.subscribe(topics)

    def close(self):
        self.consumer.close()

    def stop(self):
        self.stop_requested = True

    def loop_consume_message(self):
        try:
            while not self.stop_requested:
                msg = self.consume()
                if not msg:
                    continue
                # Execute the provided message handler function based on the topic
                handler = self.topic_handlers.get(msg.topic())
                if handler:
                    handler(msg.value())
        except KeyboardInterrupt:
            pass
        finally:
            self.consumer.close()

    def listen_for_response(self, timeout):
        start_time = time.time()
        while time.time() - start_time < timeout:
            msg = self.consume()
            if not msg:
                continue
            return msg.value()  # Return the response message if received
        logging.info("Timeout reached. No response received within the timeout interval.")
        return None

    def consume(self, num_messages=1, timeout=1):
        msg = self.consumer.consume(num_messages, timeout=timeout)
        if len(msg) == 0:
            logging.debug('Timeout consuming')
            return
        msg = msg[0]
        if msg.error():
            if msg.error().code() == KafkaError._PARTITION_EOF:
                return
            logging.error(f'Error: {msg.error()}')
            raise Exception()
        m = {'topic': msg.topic(), 'key': msg.key(), 'value': msg.value(), 'partition': msg.partition(),
             'offset': msg.offset()}
        logging.debug(f'Received message: {m}')
        return msg
