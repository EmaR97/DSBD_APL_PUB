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

static const char *const configPath = "config_ps.json";


void send_result_to_services(my_namespace::kafka::KafkaProducer &producer, int64 timestamp, const std::string &cam_id,
                             bool detected);

void parse_message(const RdKafka::Message &message, google::protobuf::Timestamp &timestamp, std::string &cam_id,
                   std::vector<uchar> &imgBuffer);


std::string formatString(int64_t value1, const std::string &str);

#endif //CAMCONTROLLER_PROCESSING_SERVER_H
