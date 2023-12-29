//
// Created by emanuele on 17/12/23.
//

#include "cam_controller.h"

using namespace my_namespace;
using utility::LogLevel;

// Signal handler for graceful shutdown
void handleShutdown(int signal) {
    // Shutdown command listener and video stream processor
    commandListener->shutdown();
    videoStreamProcessor->shutdown();
}

int main() {
    // Set up signal handler for graceful shutdown
    std::signal(SIGINT, handleShutdown);

    // Initialize logger
    utility::Logger &logger = utility::Logger::getInstance();

    // Get the current working directory
    std::filesystem::path currentPath = std::filesystem::current_path();

    // Print the current working directory
    logger << LogLevel::INFO << "Current Working Directory: " << currentPath << std::endl;

    // Load configuration from a file
    config = utility::loadConfigFromFile(configPath);

    // Initialize video stream processor
    videoStreamProcessor = new video::FrameProcessor();

    // Initialize frame uploader with configuration parameters
    frameUploader.initialize(config["auth"]["address"], 1000 / std::stof(to_string(config["FRAMERATE"])),
                             config["auth"]["username"], config["auth"]["password"], config["auth"]["endpoint"]);

//    // Perform login and update configuration with the obtained token
//    auto token = frameUploader.performLogin(config["CAM_ID"]);
//    if (token.empty()) {
//        exit(0);
//    }
//    config["TOKEN"] = token;
//    //TODO: implement dynamic creation of Kafka credential for cam access on login

    // Initialize MQTT handler for receiving commands
    commandListener = new mqtt_::MqttHandler(config["mqtt"]["address"], config["mqtt"]["client_id"]);
    // Connect to MQTT server using provided credentials
    commandListener->connect(config["mqtt"]["username"], config["mqtt"]["password"]);

    // Define base topic for MQTT commands
    std::string baseTopic = config["mqtt"]["topic"];

    // Define lambda function as callback for executing commands received from MQTT
    auto commandExecuter = [&](const std::string &payload) {
        videoStreamProcessor->startAndStopProcessing(payload);
    };

    // Define lambda function to subscribe to MQTT topic and set callback
    auto commandListenerFunction = [&]() {
        // Subscribe to specific MQTT topic for this camera with defined callback
        commandListener->subscribe(baseTopic + static_cast<std::string>(config["CAM_ID"]), commandExecuter);
    };

    // Create thread for MQTT subscriber
    std::thread subscriber(commandListenerFunction);

    // Uncomment the following lines if you want to set the video source index
    int source = 0;
    videoStreamProcessor->setSourceIndex(&source);

    // Set alternate source for video frames
    // std::string source = config["SOURCE_PATH"];
    // videoStreamProcessor.setAlternateSource(&source);

    // Initialize Kafka producer
    kafka::KafkaProducer producer(config["kafka"]["address"]);
    // Define lambda function for sending video frames using HttpHandler
    auto videoSendingFunction = [&](const std::vector<uchar> &buffer) {
        my_namespace::message::FrameData protobufData;
        fillProtobufData(protobufData, buffer, config["CAM_ID"]);

        std::string jsonOutput;
        google::protobuf::util::JsonPrintOptions options;
        options.add_whitespace = true;  // Add whitespace to make the output more readable
        google::protobuf::util::MessageToJsonString(protobufData, &jsonOutput, options).ok();
        logger << LogLevel::INFO << "Sending frame" << std::endl;
        producer.sendMessage(config["kafka"]["topic_frame_data"], jsonOutput);
    };

    // Define lambda function for processing video frames and applying the callback
    auto videoProcessingFunction = [&]() {
        videoStreamProcessor->processFramesAndApplyCallback(config["FRAMERATE"], config["FORMAT"],
                                                            videoSendingFunction);
    };

    // Create thread for video frame processor and publisher
    std::thread publisher(videoProcessingFunction);

    // Wait for publisher thread to finish
    publisher.join();

    // Wait for subscriber thread to finish
    subscriber.join();
    delete commandListener;
    delete videoStreamProcessor;
    // Return 0 to indicate successful program execution
    return 0;
}

// Function to fill protobuf data
void
fillProtobufData(message::FrameData &protobufData, const std::vector<uchar> &frameData, const std::string &cam_id) {
    protobufData.mutable_timestamp()->set_seconds(my_namespace::utility::timestamp_ns().seconds());
    protobufData.mutable_timestamp()->set_nanos(my_namespace::utility::timestamp_ns().nanos());
    protobufData.set_frame_data(frameData.data(), frameData.size());
    protobufData.set_cam_id(cam_id);
}
