#include "processing_server.h"

using namespace my_namespace;
using namespace utility;

// Entry point of the application
int main() {
    // Print the current working directory
    printCurrentWorkingDirectory();

    // Load configuration settings from a JSON file
    nlohmann::json config = loadConfiguration();

    // Create Kafka producer and MinIO uploader instances
    kafka::KafkaProducer producer(config["kafka"]["broker"]);
    sender::MinIOUploader minioUploader(config["minio"]["endpoint"], config["minio"]["bucketName"],
                                        config["minio"]["keyId"], config["minio"]["keySecret"]);

    // Lambda function to process Kafka messages
    auto processMessage = [&](RdKafka::Message &message) {
        processKafkaMessage(message, producer, minioUploader, config);
    };

    // Create a KafkaConsumer instance with the processMessage function
    kafka::KafkaConsumer kafkaConsumer(config["kafka"]["broker"], config["kafka"]["group_id"],
                                       config["kafka"]["topic_frame_data"], processMessage);
    Logger::getInstance() << LogLevel::INFO << "Consuming" << std::endl;

    // Start consuming messages
    kafkaConsumer.startConsuming();

    return 0;
}

// Parse Kafka message into timestamp, camera ID, and image buffer
void parseMessage(const RdKafka::Message &message, google::protobuf::Timestamp &timestamp, std::string &cam_id,
                  std::vector<uchar> &imgBuffer) {
    // Deserialize the received Kafka message into a protobuf object
    auto payload = static_cast<char *>(message.payload());
    Logger::getInstance() << LogLevel::INFO << "Received message with size: " << message.len() << " bytes" << std::endl;
//    std::cout << "Received payload with size: " << strlen(payload) << " bytes" << std::endl;

    // Create and populate a FrameData protobuf object from the JSON payload
    message::FrameData frameData;
    google::protobuf::util::JsonParseOptions jsonParseOptions;
    google::protobuf::util::JsonStringToMessage(payload, &frameData, jsonParseOptions).ok();

    // Extract data from the protobuf message
    timestamp = frameData.timestamp();
    cam_id = frameData.cam_id();

    // Extract the frame data (image) from the protobuf object
    std::string receivedString = frameData.frame_data();
    const auto *buffer = reinterpret_cast<const uchar *>(receivedString.c_str());
    imgBuffer.assign(buffer, buffer + receivedString.size());
}

// Send results to Kafka topics based on processed data
void sendResultToServices(const kafka::KafkaProducer &producer, int64 timestamp, const std::string &cam_id,
                     bool detected, const nlohmann::json &config) {
    // Create and populate a FrameInfo protobuf object
    message::FrameInfo frameInfo;
    frameInfo.set_cam_id(cam_id);
    frameInfo.set_timestamp(timestamp);
    frameInfo.set_persondetected(detected);

    // Serialize the protobuf object to a string and send it to a Kafka topic
    std::string jsonOutput = frameInfo.SerializeAsString();
    producer.sendMessage(config["kafka"]["topic_frame_info"], jsonOutput);
}

// Format a string with an integer and another string
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

// Print the current working directory
void printCurrentWorkingDirectory() {
    std::filesystem::path currentPath = std::filesystem::current_path();
    Logger::getInstance() << LogLevel::INFO << "Current Working Directory: " << currentPath << std::endl;
}

// Load configuration settings from a JSON file
nlohmann::json loadConfiguration() {
    try {
        return loadConfigFromFile(configPath);
    } catch (const std::exception &e) {
        Logger::getInstance() << LogLevel::ERROR << "Error loading configuration: " << e.what() << std::endl;
        throw; // Rethrow the exception for the caller to handle
    }
}

// Process Kafka message, handle exceptions and log errors
void processKafkaMessage(const RdKafka::Message &message, const kafka::KafkaProducer &producer,
                         const sender::MinIOUploader &minioUploader, const nlohmann::json &config) {
    // Variables to store parsed message data
    google::protobuf::Timestamp timestamp;
    std::string cam_id;
    std::vector<uchar> imgBuffer;

    try {
        // Parse the Kafka message
        parseMessage(message, timestamp, cam_id, imgBuffer);

        Logger::getInstance() << LogLevel::INFO << "Timestamp: " << timestamp.ByteSizeLong() << std::endl;

        // Calculate Unix time in nanoseconds from timestamp
        int64_t unixTimeNanos = timestamp.seconds() * 1e9 + timestamp.nanos();

        // Generate image name and apply detection algorithm
        std::string image_name = formatString(unixTimeNanos, cam_id);
        bool detected = video::convertDetect(imgBuffer);

        // Store marked image in MinIO
        minioUploader.uploadImage(image_name, imgBuffer);

        // Send results to Kafka topics
        sendResultToServices(producer, unixTimeNanos, cam_id, detected, config);
    } catch (const std::exception &e) {
        // Log errors if an exception occurs during message processing
        Logger::getInstance() << LogLevel::ERROR << "Error processing Kafka message: " << e.what() << std::endl;
    }
}
