//
// Created by emanuele on 17/12/23.
//

#include "cam_controller.h"

using namespace my_namespace;
using namespace utility;

int main() {
    // Set up a signal handler for graceful shutdown when receiving SIGINT (Ctrl+C)
    std::signal(SIGINT, handleShutdown);

    // Initialize the logger instance
    Logger &logger = Logger::getInstance();

    // Get the current working directory
    std::filesystem::path currentPath = std::filesystem::current_path();

    // Print the current working directory to the logger
    logger << LogLevel::INFO << "Current Working Directory: " << currentPath << std::endl;

    // Load configuration from a file
    try {
        config = loadConfigFromFile(configPath);
    } catch (const std::exception &e) {
        // Log an error message if configuration loading fails
        logger << LogLevel::ERROR << "Failed to load configuration: " << e.what() << std::endl;
        return 1; // Return a non-zero value to indicate an error
    }

//    // Perform login and update configuration with the obtained token
//    auto token = frameUploader.performLogin(config["CAM_ID"]);
//    if (token.empty()) {
//        exit(0);
//    }
//    config["TOKEN"] = token;
//TODO: implement dynamic creation of Kafka credential for cam access on login


    // Create threads for MQTT subscriber and video frame processor/publisher
    std::thread commandListener(runCommandListener);
    // Create thread for video frame processor and publisher
    std::thread frameSender(processAndSenFrame);

    // Wait for the publisher thread to finish
    frameSender.join();
    // Wait for the subscriber thread to finish
    commandListener.join();

    return 0;
}

void processAndSenFrame() {
    // Initialize the frame uploader with configuration parameters
    frameUploader.initialize(config["auth"]["address"], 1000 / std::stof(to_string(config["FRAMERATE"])),
                             config["auth"]["username"], config["auth"]["password"], config["auth"]["endpoint"]);

    int source = 0;
    videoStreamProcessor.setSourceIndex(&source);
    Logger &logger = Logger::getInstance();

    // Set alternate source for video frames
    // std::string source = config["SOURCE_PATH"];
    // videoStreamProcessor.setAlternateSource(&source);

    // Initialize Kafka producer
    kafka::KafkaProducer producer(config["kafka"]["address"]);

    // Define lambda function for sending video frames using HttpHandler
    auto videoSendingFunction = [&](const std::vector<uchar> &buffer) {
        sendFrame(logger, producer, buffer);
    };

    // Process frames and apply the callback function
    videoStreamProcessor.processFramesAndApplyCallback(config["FRAMERATE"], config["FORMAT"], videoSendingFunction);
}

void runCommandListener() {
    // Initialize MQTT handler for receiving commands
    mqtt_::MqttHandler commandListener(config["mqtt"]["address"], config["mqtt"]["client_id"],
                                       config["mqtt"]["username"], config["mqtt"]["password"]);

    std::string baseTopic = config["mqtt"]["topic"];

    // Define lambda function as a callback for executing commands received from MQTT
    auto commandExecuter = [&](const std::string &payload) {
        videoStreamProcessor.startAndStopProcessing(payload);
    };

    // Connect to MQTT server using provided credentials
    commandListener.connect();

    // Subscribe to a specific MQTT topic for this camera with the defined callback
    commandListener.subscribe(baseTopic + static_cast<std::string>(config["CAM_ID"]), commandExecuter);

    // Wait until notified to exit
    std::unique_lock<std::mutex> lock(exitMutex);
    exitCondition.wait(lock);

    // Shutdown the command listener
    commandListener.shutdown();
}

void sendFrame(Logger &logger, const kafka::KafkaProducer &producer, const std::vector<uchar> &buffer) {
    // Prepare and send the frame data using Kafka
    message::FrameData protobufData;
    fillProtobufData(protobufData, buffer, config["CAM_ID"]);

    std::string jsonOutput;
    google::protobuf::util::JsonPrintOptions options;
    options.add_whitespace = true;  // Add whitespace to make the output more readable
    google::protobuf::util::MessageToJsonString(protobufData, &jsonOutput, options).ok();
    logger << LogLevel::INFO << "Sending frame" << std::endl;
    producer.sendMessage(config["kafka"]["topic_frame_data"], jsonOutput);
}

// Function to fill protobuf data
void
fillProtobufData(message::FrameData &protobufData, const std::vector<uchar> &frameData, const std::string &cam_id) {
    protobufData.mutable_timestamp()->set_seconds(timestamp_ns().seconds());
    protobufData.mutable_timestamp()->set_nanos(timestamp_ns().nanos());
    protobufData.set_frame_data(frameData.data(), frameData.size());
    protobufData.set_cam_id(cam_id);
}

// Signal handler for graceful shutdown
void handleShutdown(int signal) {
    // Shutdown command listener and video stream processor
    videoStreamProcessor.shutdown();

    // Notify waiting threads to exit
    std::unique_lock<std::mutex> lock(exitMutex);
    exitCondition.notify_all();
}
