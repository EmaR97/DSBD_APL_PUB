#include <iostream>
#include <utility>
#include "KafkaConsumer.h"
#include "../utility/logger.h"
#include <librdkafka/rdkafkacpp.h>

namespace my_namespace::kafka {
    utility::Logger &logger = utility::Logger::getInstance();

    KafkaConsumer::KafkaConsumer(std::string brokers, std::string group_id, std::string topic_name,
                                 MessageProcessor messageProcessor) : brokers_(std::move(brokers)),
                                                                      group_id_(std::move(group_id)),
                                                                      topic_name_(std::move(topic_name)),
                                                                      messageProcessor(std::move(messageProcessor)) {
        // Set up Kafka configuration
        conf_ = RdKafka::Conf::create(RdKafka::Conf::CONF_GLOBAL);
        conf_->set("metadata.broker.list", brokers_, errstr);
        // Set group id and other configurations
        conf_->set("group.id", group_id_, errstr);
        // Create Kafka consumer
        consumer_ = RdKafka::KafkaConsumer::create(conf_, errstr);
        if (!consumer_) {
            logger << utility::LogLevel::ERROR << "Failed to create consumer: " << errstr << std::endl;
            std::exit(-1);
        }

        // Subscribe to the topic
        std::vector<std::string> topics = {topic_name_};
        consumer_->subscribe(topics);
        delete conf_;

    }

    KafkaConsumer::~KafkaConsumer() {
        // Unsubscribe and clean up
        consumer_->unsubscribe();
        delete consumer_;
        delete conf_;
    }

    void KafkaConsumer::startConsuming() {
        // Start consuming messages

        while (true) {
            RdKafka::Message *message;
            // Poll for messages with a timeout

            message = consumer_->consume(1000); // 1 second timeout
            logger << utility::LogLevel::INFO << "Time out" << std::endl;

            switch (message->err()) {
                case RdKafka::ERR__TIMED_OUT:
                    break;
                case RdKafka::ERR_NO_ERROR:
                    logger << utility::LogLevel::INFO << "Received message with size: " << message->len() << " bytes"
                           << ", offset: " << message->offset() << std::endl;
                    messageProcessor(*message);
                    consumer_->commitSync();
                    break;
                default:
                    logger << utility::LogLevel::ERROR << "Error while consuming message: " << message->errstr()
                           << std::endl;
            }
            delete message;

        }
    }

};
