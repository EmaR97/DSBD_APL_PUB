#ifndef KAFKA_CONSUMER_H
#define KAFKA_CONSUMER_H

#include <string>
#include <librdkafka/rdkafkacpp.h>
#include <functional>
#include <iostream>

namespace my_namespace::kafka {
    using MessageProcessor = std::function<void(RdKafka::Message &)>;

    class KafkaConsumer {
    public:

        KafkaConsumer(std::string brokers, std::string group_id, std::string topic_name,
                      MessageProcessor messageProcessor);

        ~KafkaConsumer();

        void startConsuming();

    private:
        std::string brokers_;
        std::string group_id_;
        std::string topic_name_;
        MessageProcessor messageProcessor;
        RdKafka::Conf *conf_{};
        RdKafka::KafkaConsumer *consumer_{};

        // Error string for librdkafka
        std::string errstr;

    };

}

#endif //KAFKA_CONSUMER_H
