#ifndef KAFKA_PRODUCER_H
#define KAFKA_PRODUCER_H

#include <iostream>
#include <string>
#include <cstdlib>
#include <csignal>
#include <librdkafka/rdkafkacpp.h>
#include "../utility/logger.h"


namespace my_namespace::kafka {



    class ExampleDeliveryReportCb : public RdKafka::DeliveryReportCb {
    public:
        void dr_cb(RdKafka::Message &message) override;
    };

    class KafkaProducer {
    private:
        RdKafka::Producer *producer;
        ExampleDeliveryReportCb delivery_report_cb;

    public:
        explicit KafkaProducer(const std::string &brokers);

        ~KafkaProducer();

        void sendMessage(const std::string &topic, const std::string &message);

        void wait_for_pending_delivery(int timeoutMs);

        int outq_len();
    };
}

#endif //KAFKA_PRODUCER_H
