#include <iostream>
#include <string>
#include <cstdlib>
#include <librdkafka/rdkafkacpp.h>
#include "../utility/logger.h"
#include "KafkaProducer.h"

using my_namespace::utility::LogLevel;
namespace my_namespace::kafka {
    utility::Logger &logger_ = utility::Logger::getInstance();

    // Callback function for delivery reports
    void ExampleDeliveryReportCb::dr_cb(RdKafka::Message &message) {
        if (message.err())
            logger_ << utility::LogLevel::ERROR << "Message delivery failed: " << message.errstr() << std::endl;
        else
            logger_ << utility::LogLevel::INFO << "Message delivered to topic " << message.topic_name() << " ["
                   << message.partition() << "] at offset " << message.offset() << std::endl;
    }

    // Constructor for KafkaProducer
    KafkaProducer::KafkaProducer(const std::string &brokers) {
        // Create a Kafka configuration object
        RdKafka::Conf *conf = RdKafka::Conf::create(RdKafka::Conf::CONF_GLOBAL);
        std::string errstr;

        // Set the bootstrap servers for the broker
        if (conf->set("bootstrap.servers", brokers, errstr) != RdKafka::Conf::CONF_OK) {
            logger_ << utility::LogLevel::ERROR << errstr << std::endl;
            exit(1);
        }

        // Set the delivery report callback
        if (conf->set("dr_cb", &delivery_report_cb, errstr) != RdKafka::Conf::CONF_OK) {
            logger_ << utility::LogLevel::ERROR << errstr << std::endl;
            exit(1);
        }

        // Create a Kafka producer instance
        producer = RdKafka::Producer::create(conf, errstr);
        if (!producer) {
            logger_ << utility::LogLevel::ERROR << "Failed to create producer: " << errstr << std::endl;
            exit(1);
        }

        // Configuration object is no longer needed after creating the producer
        delete conf;
    }

    // Destructor for KafkaProducer
    KafkaProducer::~KafkaProducer() {
        // Destroy the Kafka producer instance
        delete producer;
    }

    // Send a message to a specified topic
    void KafkaProducer::sendMessage(const std::string &topic, const std::string &message) {
        RdKafka::ErrorCode err = producer->produce(topic, RdKafka::Topic::PARTITION_UA, RdKafka::Producer::RK_MSG_COPY,
                                                   const_cast<char *>(message.c_str()), message.size(), nullptr, 0, 0,
                                                   nullptr, nullptr);

        if (err != RdKafka::ERR_NO_ERROR) {
            logger_ << utility::LogLevel::ERROR << "Failed to produce to topic " << topic << ": " << RdKafka::err2str(err) << std::endl;

            if (err == RdKafka::ERR__QUEUE_FULL) {
                // If the internal queue is full, wait for messages to be delivered and then retry
                producer->poll(1000);
                sendMessage(topic, message);  // Retry
            }
        } else {
            logger_ << utility::LogLevel::INFO << "Enqueued message (" << message.size() << " bytes) " << "for topic " << topic << std::endl;
        }

        // Poll to serve the delivery report queue
        producer->poll(0);
    }

    // Wait for pending message deliveries
    void KafkaProducer::wait_for_pending_delivery(int timeoutMs) {
        producer->flush(timeoutMs);
    }

    // Get the length of the producer's outgoing message queue
    int KafkaProducer::outq_len() {
        return producer->outq_len();
    }
}  // namespace my_namespace::kafka
