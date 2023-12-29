#include "processing_server.h"


using namespace my_namespace;

nlohmann::basic_json<> config_2;
utility::Logger &logger = utility::Logger::getInstance();

int main() {

    // Get the current working directory
    std::filesystem::path currentPath = std::filesystem::current_path();

    // Print the current working directory
    logger << utility::LogLevel::INFO << "Current Working Directory: " << currentPath << std::endl;

    // Load configuration from a file
    config_2 = my_namespace::utility::loadConfigFromFile(configPath);
    kafka::KafkaProducer producer(config_2["kafka"]["broker"]);


    sender::MinIOUploader minioUploader(config_2["minio"]["endpoint"], config_2["minio"]["bucketName"],
                                        config_2["minio"]["keyId"], config_2["minio"]["keySecret"]);
    // Lambda function to process Kafka messages
    auto processMessage = [&](RdKafka::Message &message) {
        google::protobuf::Timestamp timestamp;
        std::string cam_id;
        std::vector<uchar> imgBuffer;
        parse_message(message, timestamp, cam_id, imgBuffer);
        logger << utility::LogLevel::INFO << "Timestamp: " << timestamp.ByteSizeLong() << std::endl;
        int64_t unixTimeNanos = timestamp.seconds() * 1e9 + timestamp.nanos();

        std::string image_name = formatString(unixTimeNanos, cam_id);
        // apply detection
        bool detected = video::convert_detect(imgBuffer);
        // store the marked image
        minioUploader.uploadImage(image_name, imgBuffer);
        send_result_to_services(producer, unixTimeNanos, cam_id, detected);
    };

    // Create a KafkaConsumer instance with the processMessage function
    my_namespace::kafka::KafkaConsumer kafkaConsumer(config_2["kafka"]["broker"], config_2["kafka"]["group_id"],
                                                     config_2["kafka"]["topic_frame_data"], processMessage);
    logger << utility::LogLevel::INFO << "Consuming" << std::endl;

    // Start consuming messages
    kafkaConsumer.startConsuming();

    return 0;
}

void parse_message(const RdKafka::Message &message, google::protobuf::Timestamp &timestamp, std::string &cam_id,
                   std::vector<uchar> &imgBuffer) {
    // Deserialize the received Kafka message into a protobuf object
    auto payload = static_cast<char *>(message.payload());
    logger << utility::LogLevel::INFO << "Received message with size: " << message.len() << " bytes" << std::endl;
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

void
send_result_to_services(kafka::KafkaProducer &producer, int64 timestamp, const std::string &cam_id, bool detected) {
    message::FrameInfo frameInfo;
    frameInfo.set_cam_id(cam_id);
    frameInfo.set_timestamp(timestamp);
    frameInfo.set_persondetected(detected);
    std::string jsonOutput = frameInfo.SerializeAsString();    // Send on needed topics
    producer.sendMessage(config_2["kafka"]["topic_frame_info"], jsonOutput);
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