#include "processing_server.h"

using namespace my_namespace;
using namespace utility;
using namespace std::chrono;


// Entry point of the application
int main() {
    // Log the start of the application
    Logger::getInstance() << LogLevel::INFO << "Starting the processing server..." << std::endl;
    // Print the current working directory
    printCurrentWorkingDirectory();

    // Load configuration settings from a JSON file
    nlohmann::json config;
    try {
        config = loadConfiguration();
        // Log successful configuration loading
        Logger::getInstance() << LogLevel::INFO << "Configuration loaded successfully." << std::endl;
    } catch (const std::exception &e) {
        // Log errors if configuration loading fails
        Logger::getInstance() << LogLevel::ERROR << "Error loading configuration: " << e.what() << std::endl;
        return 1; // Exit the application in case of a configuration error
    }

    // Create Kafka producer and MinIO uploader instances
    kafka::KafkaProducer producer(config["kafka"]["broker"]);
    sender::MinIOUploader minioUploader(config["minio"]["endpoint"], config["minio"]["bucketName"],
                                        config["minio"]["keyId"], config["minio"]["keySecret"]);
    // Log successful creation of Kafka producer and MinIO uploader
    Logger::getInstance() << LogLevel::INFO << "Kafka producer and MinIO uploader created successfully." << std::endl;

    Chronometer chronometer;
    // Create Prometheus exposer on a specified port
    prometheus::Exposer exposer{std::string(config["prometheus"]["bind_address"])};
    // Create a Prometheus registry
    auto registry = std::make_shared<prometheus::Registry>();
    // Ask the exposer to scrape the registry on incoming HTTP requests
    exposer.RegisterCollectable(registry);
    // Add metrics definition
    auto &counter = prometheus::BuildCounter().Name("processed_messages").Help(
            "Total number of processed messages").Register(*registry).Add({{"metric", "total"}});
    auto &gauge = prometheus::BuildGauge().Name("working_time").Help(
            "Working time as a percentage of total time").Register(*registry).Add({{"metric", "gauge"}});
    auto &gauge_2 = prometheus::BuildGauge().Name("time_since_message_creation").Help(
            "Elapsed time from message creation to completing processing").Register(*registry).Add(
            {{"metric", "gauge"}});

    // Log successful creation of Prometheus exposer and registry
    Logger::getInstance() << LogLevel::INFO << "Prometheus exposer and registry created successfully." << std::endl;


    // Lambda function to process Kafka messages
    auto processMessage = [&](RdKafka::Message &message) {
        try {
            chronometer.startWorking();
            int64_t elapsedTimeInNanos;
            elapsedTimeInNanos = processKafkaMessage(message, producer, minioUploader, config);
            chronometer.stopWorking();
            gauge_2.Set(elapsedTimeInNanos);
            counter.Increment();
        } catch (const std::exception &e) {
            // Log errors if an exception occurs during message processing
            Logger::getInstance() << LogLevel::ERROR << "Error processing Kafka message: " << e.what() << std::endl;
        }

    };

    // Create a KafkaConsumer instance with the processMessage function
    kafka::KafkaConsumer kafkaConsumer(config["kafka"]["broker"], config["kafka"]["group_id"],
                                       config["kafka"]["topic_frame_data"], processMessage);

    // Log the successful start of message consumption
    Logger::getInstance() << LogLevel::INFO << "Message consumption started successfully." << std::endl;

    // Thread to periodically check and reset Cronometers
    std::thread checkThread([&chronometer, &gauge, &config]() {
        while (true) {
            std::this_thread::sleep_for(
                    std::chrono::seconds(config["prometheus"]["interval"])); // Adjust the interval as needed
            double workingTimePercentage = chronometer.checkAndReset();
            Logger::getInstance() << LogLevel::INFO << "WorkingTimePercentage: " << workingTimePercentage << std::endl;
            gauge.Set(workingTimePercentage);
        }
    });

    // Start consuming messages
    kafkaConsumer.startConsuming();
    checkThread.join();
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
void
sendResultToServices(const kafka::KafkaProducer &producer, int64 timestamp, const std::string &cam_id, bool detected,
                     const nlohmann::json &config) {
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
long processKafkaMessage(const RdKafka::Message &message, const kafka::KafkaProducer &producer,
                         const sender::MinIOUploader &minioUploader, const nlohmann::json &config) {
    const auto now = std::chrono::high_resolution_clock::now();
    // Variables to store parsed message data
    google::protobuf::Timestamp timestamp;
    std::string cam_id;
    std::vector<uchar> imgBuffer;
    int64_t unixTimeNanos;

    // Parse the Kafka message
    parseMessage(message, timestamp, cam_id, imgBuffer);

    Logger::getInstance() << LogLevel::INFO << "Timestamp: " << timestamp.ByteSizeLong() << std::endl;

    // Calculate Unix time in nanoseconds from timestamp
    unixTimeNanos = timestamp.seconds() * 1e9 + timestamp.nanos();

    // Generate image name and apply detection algorithm
    std::string image_name = formatString(unixTimeNanos, cam_id);
    bool detected = video::convertDetect(imgBuffer);

    // Store marked image in MinIO
    minioUploader.uploadImage(image_name, imgBuffer);

    // Send results to Kafka topics
    sendResultToServices(producer, unixTimeNanos, cam_id, detected, config);


    long elapsedTime = time_point_cast<nanoseconds>(now).time_since_epoch().count() - unixTimeNanos;
    return elapsedTime;
}
