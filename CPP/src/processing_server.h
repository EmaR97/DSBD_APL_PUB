//
// Created by emanuele on 17/12/23.
//

#ifndef CAMCONTROLLER_PROCESSING_SERVER_H
#define CAMCONTROLLER_PROCESSING_SERVER_H

#include <string>
#include <iostream>
#include <filesystem>
#include "kafka/KafkaConsumer.h"
#include "utility/config.h"
#include "message/framedata.pb.h"
#include "message/notification.pb.h"
#include "message/frameinfo.pb.h"
#include <opencv2/opencv.hpp>
#include <google/protobuf/util/json_util.h>
#include "video/frame_processor.h"
#include "video/detection.h"
#include "kafka/KafkaProducer.h"
#include "http/MinIOUploader.h"
#include <prometheus/exposer.h>
#include <prometheus/registry.h>
#include <prometheus/counter.h>
#include <prometheus/gauge.h>
#include "utility/chronometer.h"

static const char *const configPath = "config_ps.json";


void
sendResultToServices(const my_namespace::kafka::KafkaProducer &producer, int64 timestamp, const std::string &cam_id,
                     bool detected, const nlohmann::json &config);

void parseMessage(const RdKafka::Message &message, google::protobuf::Timestamp &timestamp, std::string &cam_id,
                  std::vector<uchar> &imgBuffer);


long processKafkaMessage(const RdKafka::Message &message, const my_namespace::kafka::KafkaProducer &producer,
                         const my_namespace::sender::MinIOUploader &minioUploader, const nlohmann::json &config);

std::string formatString(int64_t value1, const std::string &str);


void printCurrentWorkingDirectory();

nlohmann::json loadConfiguration();

#endif //CAMCONTROLLER_PROCESSING_SERVER_H
