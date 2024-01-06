#include "processing_server.h"


using namespace my_namespace;


int main() {

    printCurrentWorkingDirectory();

    nlohmann::json config = loadConfiguration();

    kafka::KafkaProducer producer(config["kafka"]["broker"]);
    sender::MinIOUploader minioUploader(config["minio"]["endpoint"], config["minio"]["bucketName"],
                                        config["minio"]["keyId"], config["minio"]["keySecret"]);

    // Lambda function to process Kafka messages
    auto processMessage = [&](RdKafka::Message &message) {
        processKafkaMessage(message, producer, minioUploader, config);
    };

    // Create a KafkaConsumer instance with the processMessage function
    my_namespace::kafka::KafkaConsumer kafkaConsumer(config["kafka"]["broker"], config["kafka"]["group_id"],
                                                     config["kafka"]["topic_frame_data"], processMessage);
    utility::Logger::getInstance() << utility::LogLevel::INFO << "Consuming" << std::endl;

    // Start consuming messages
    kafkaConsumer.startConsuming();

    return 0;
}

void parseMessage(const RdKafka::Message &message, google::protobuf::Timestamp &timestamp, std::string &cam_id,
                  std::vector<uchar> &imgBuffer) {
    // Deserialize the received Kafka message into a protobuf object
    auto payload = static_cast<char *>(message.payload());
    utility::Logger::getInstance() << utility::LogLevel::INFO << "Received message with size: " << message.len()
                                   << " bytes" << std::endl;
//    std::cout << "Received payload with size: " << strlen(payload) << " bytes" << std::endl;
    message::FrameData frameData;
    google::protobuf::util::JsonParseOptions jsonParseOptions;
    google::protobuf::util::JsonStringToMessage(payload, &frameData, jsonParseOptions).ok();
    // Extract data from message
    timestamp = frameData.timestamp();
    cam_id = frameData.cam_id();
    // Extract the frame data (image) from the protobuf object
    std::string receivedString = frameData.frame_data();
    const auto *buffer = reinterpret_cast<const uchar *>(receivedString.c_str());
    imgBuffer.assign(buffer, buffer + receivedString.size());
}

void sendResultToServices(const kafka::KafkaProducer &producer, int64 timestamp, const std::string &cam_id, bool detected,
                          const nlohmann::json &config) {
    message::FrameInfo frameInfo;
    frameInfo.set_cam_id(cam_id);
    frameInfo.set_timestamp(timestamp);
    frameInfo.set_persondetected(detected);
    std::string jsonOutput = frameInfo.SerializeAsString();    // Send on needed topics
    producer.sendMessage(config["kafka"]["topic_frame_info"], jsonOutput);
}

std::string formatString(int64_t value1, const std::string &str) {
    // Calculate the size of the buffer needed for the formatted string
    size_t bufferSize = snprintf(nullptr, 0, "%s/%ld.jpg", str.c_str(), value1) + 1;

    // Create a buffer of the required size
    char buffer[bufferSize];

    // Use snprintf to format the string into the buffer
    snprintf(buffer, bufferSize, "%s/%ld.jpg", str.c_str(), value1);
    std::string name = std::string(buffer);
    // Return the formatted string
    return name;
}


void printCurrentWorkingDirectory() {
    std::filesystem::path currentPath = std::filesystem::current_path();
    utility::Logger::getInstance() << utility::LogLevel::INFO << "Current Working Directory: " << currentPath
                                   << std::endl;
}

nlohmann::json loadConfiguration() {
    try {
        return utility::loadConfigFromFile(configPath);
    } catch (const std::exception &e) {
        utility::Logger::getInstance() << utility::LogLevel::ERROR << "Error loading configuration: " << e.what()
                                       << std::endl;
        throw; // Rethrow the exception for the caller to handle
    }
}


void
processKafkaMessage(const RdKafka::Message &message, const kafka::KafkaProducer &producer, const sender::MinIOUploader &minioUploader,
                    const nlohmann::json &config) {
    google::protobuf::Timestamp timestamp;
    std::string cam_id;
    std::vector<uchar> imgBuffer;

    try {
        parseMessage(message, timestamp, cam_id, imgBuffer);

        utility::Logger::getInstance() << utility::LogLevel::INFO << "Timestamp: " << timestamp.ByteSizeLong()
                                       << std::endl;

        int64_t unixTimeNanos = timestamp.seconds() * 1e9 + timestamp.nanos();
        std::string image_name = formatString(unixTimeNanos, cam_id);
        //apply detection algorithm
        bool detected = video::convertDetect(imgBuffer);
        //store marked image
        minioUploader.uploadImage(image_name, imgBuffer);

        sendResultToServices(producer, unixTimeNanos, cam_id, detected, config);
    } catch (const std::exception &e) {
        utility::Logger::getInstance() << utility::LogLevel::ERROR << "Error processing Kafka message: " << e.what()
                                       << std::endl;
    }
}
