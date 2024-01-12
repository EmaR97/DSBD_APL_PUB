//
// Created by emanuele on 17/12/23.
//

#ifndef CAMCONTROLLER_CAM_CONTROLLER_H
#define CAMCONTROLLER_CAM_CONTROLLER_H


#include "video/frame_processor.h"
#include "mqtt/MqttHandler.h"
#include "utility/config.h"
#include "utility/json.hpp"
#include <string>
#include <thread>
#include <mutex>
#include "kafka/KafkaProducer.h"
#include <csignal>
#include "utility/utility.h"
#include <google/protobuf/util/json_util.h>

static const char *const configPath = "config_cc.json";

std::condition_variable exitCondition;

std::mutex exitMutex;

auto videoStreamProcessor = my_namespace::video::FrameProcessor();

nlohmann::basic_json<> config;

my_namespace::sender::HttpHandler frameUploader;

void processAndSenFrame();

void runCommandListener();

void fillProtobufData(my_namespace::message::FrameData &protobufData, const std::vector<uchar> &frameData,
                      const std::string &cam_id);

void sendFrame(my_namespace::utility::Logger &logger, const my_namespace::kafka::KafkaProducer &producer,
               const std::vector<uchar> &buffer);

void handleShutdown(int signal);

#endif //CAMCONTROLLER_CAM_CONTROLLER_H
