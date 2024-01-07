cmake_minimum_required(VERSION 3.27)
project(CamController)

set(CMAKE_CXX_STANDARD 17)


# Find required packages
find_package(AWSSDK CONFIG COMPONENTS core dynamodb kinesis s3 REQUIRED)
find_package(Protobuf REQUIRED)
find_package(Threads REQUIRED)
find_package(OpenCV REQUIRED)
find_package(CURL REQUIRED)
find_package(PahoMqttCpp REQUIRED)
find_package(OpenSSL REQUIRED)
find_package(prometheus-cpp CONFIG REQUIRED)

# Include directories
include_directories(${OPENSSL_INCLUDE_DIR})

# Source files for CamController
set(CAM_CONTROLLER_SOURCES
        ../src/message/framedata.pb.h
        ../src/message/framedata.pb.cc
        ../src/video/frame_processor.cpp
        ../src/video/frame_processor.h
        ../src/http/HttpHandler.cpp
        ../src/http/HttpHandler.h
        ../src/utility/logger.h
        ../src/utility/utility.cpp
        ../src/utility/utility.h
        ../src/mqtt/MqttHandler.cpp
        ../src/mqtt/MqttHandler.h
        ../src/utility/config.cpp
        ../src/utility/config.h
        ../src/utility/json.hpp
        ../src/kafka/KafkaProducer.cpp
        ../src/kafka/KafkaProducer.h
        ../src/cam_controller.cpp
        ../src/cam_controller.h
)

# CamController executable
add_executable(CamController ${CAM_CONTROLLER_SOURCES})
target_link_libraries(CamController PRIVATE rdkafka++ ${PROTOBUF_LIBRARIES} ${OpenCV_LIBS} ${CURL_LIBRARIES} paho-mqttpp3 paho-mqtt3as)

# ProcessingServer executable
set(PROCESSING_SERVER_SOURCES
        ../src/message/framedata.pb.h
        ../src/message/framedata.pb.cc
        ../src/message/notification.pb.h
        ../src/message/notification.pb.cc
        ../src/message/frameinfo.pb.h
        ../src/message/frameinfo.pb.cc
        ../src/utility/logger.h
        ../src/utility/utility.cpp
        ../src/utility/utility.h
        ../src/utility/config.cpp
        ../src/utility/config.h
        ../src/utility/json.hpp
        ../src/kafka/KafkaProducer.cpp
        ../src/kafka/KafkaProducer.h
        ../src/kafka/KafkaConsumer.cpp
        ../src/kafka/KafkaConsumer.h
        ../src/processing_server.cpp
        ../src/processing_server.h
        ../src/http/MinIOUploader.cpp
        ../src/http/MinIOUploader.h
        ../src/video/detection.cpp
        ../src/video/detection.h
        ../src/utility/chronometer.cpp
        ../src/utility/chronometer.h
)

# ProcessingServer executable
add_executable(ProcessingServer ${PROCESSING_SERVER_SOURCES})
target_include_directories(ProcessingServer PRIVATE ${AWSSDK_INCLUDE_DIRS})
target_link_libraries(ProcessingServer PRIVATE ${AWSSDK_LIBRARIES} rdkafka++ ${PROTOBUF_LIBRARIES} ${OpenCV_LIBS} ${CURL_LIBRARIES} ${OPENSSL_LIBRARIES} prometheus-cpp::core prometheus-cpp::pull)
