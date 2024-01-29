# Camera Controller Application

The Camera Controller Application is a C++ program designed to manage and control a camera stream. It includes
components for video processing, MQTT command reception, and frame uploading to a server.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Components](#components)
- [Signal Handling](#signal-handling)

## Features

- **Video Processing:** The application uses the OpenCV library to process video frames, detect people, and execute
  commands based on MQTT messages.

- **MQTT Command Reception:** An MQTT handler listens for commands related to camera control. The application subscribes
  to a specific MQTT topic to receive commands.

- **Frame Uploading:** Processed video frames are uploaded to a server using a HTTPHandler component. The uploader
  uses a Protobuf message format for efficient data transmission.

- **Signal Handling:** The application handles the SIGINT signal (Ctrl+C) for graceful shutdown.

## Prerequisites

- C++ Compiler (Supporting C++17)
- CMake (Version 3.10 or higher)
- OpenCV Library
- libcurl Library
- nlohmann/json Library
- Protobuf Library

### Installing Prerequisites on Ubuntu

1. Update and upgrade existing packages
    ```bash
    sudo apt update && sudo apt upgrade -y
   
2. Install essential development tools
    ```bash
    sudo apt-get install -y cmake gdb g++

3. Install libraries and dependencies for multimedia and development
    ```bash
    sudo apt-get install -y ffmpeg libprotobuf-dev protobuf-compiler libopencv-dev libcurl4-openssl-dev

4. Install Linux tools and hardware data
    ```bash
    sudo apt install -y linux-tools-virtual hwdata

5. Configure USB/IP and update alternatives
    ```bash
    sudo update-alternatives --install /usr/local/bin/usbip usbip "$(ls /usr/lib/linux-tools/*/usbip | tail -n1)" 20

6. Install MQTT Library (Paho C++)
    ```bash
    git clone https://github.com/eclipse/paho.mqtt.cpp.git
    sudo apt-get install -y libssl-dev libpaho-mqtt-dev
    cd paho.mqtt.cpp
    mkdir build
    cd build
    cmake ..
    make
    sudo make install

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/camera-controller.git
   cd camera-controller
   
2. Build the application using CMake:

    ```bash
    mkdir build
    cd build
    cmake ..
    make
   
3. Run the executable:

    ```bash
    ./camera-controller

## Usage

1. Start the application, and it will load configuration settings from the config.json file.

2. The application connects to the MQTT server and subscribes to the specified topic.

3. Video frames are processed in real-time, and detected people are logged.

4. Commands received through MQTT (e.g., "START" or "STOP") control the video processing.

5. Processed video frames are uploaded to the server specified in the configuration.

6. Gracefully shutdown the application by pressing Ctrl+C.

## Configuration

The application uses a configuration file (`config.json`) to set parameters such as MQTT server details, video source,
and server endpoints.

**Sample `config.json`:**

    {
      "MQTT_SERVER_ADDRESS": "mqtt://mqtt-broker.com",
      "CLIENT_ID_SUB": "camera-client-sub",
      "TOPIC": "camera/",
      "USERNAME_MQTT": "mqtt-username",
      "PASSWORD_MQTT": "mqtt-password",
      "GOLANG_SERVER_ADDRESS": "http://golang-server.com",
      "USERNAME_UPLOAD": "upload-username",
      "PASSWORD_UPLOAD": "upload-password",
      "LOGIN_ENDPOINT": "/login",
      "UPLOAD_ENDPOINT": "/upload",
      "CAM_ID": "unique-camera-id",
      "SOURCE_PATH": "path/to/video/source.mp4",
      "FRAMERATE": 30,
      "FORMAT": "JPEG"
    }

Adjust the configuration parameters according to your setup.

## Components

1. FrameProcessor
   The FrameProcessor class processes video frames, detects people, and controls video streaming based on commands
   received through MQTT.

2. MqttHandler
   The MqttHandler class manages MQTT communication, including connecting to the MQTT server, subscribing to topics, and
   handling incoming messages.

3. HTTPHandler
   The HTTPHandler class handles the uploading of processed video frames to a specified server. It uses Protobuf for
   message serialization.

## Signal Handling

The application handles the SIGINT signal (Ctrl+C) to perform a graceful shutdown. Upon receiving the signal, the MQTT
handler and video stream processor are shutdown before exiting.

## VLC Configuration for Video Capture

The Camera Controller Application is configured to work with VLC running in capture mode with the following options:

   ```bash
   sout=#transcode{vcodec=h264,scale=Auto,acodec=mpga,ab=128,channels=2,samplerate=44100}:http{mux=ts,dst=:8555/}
   ```
Make sure to set up VLC with this configuration before running the Camera Controller Application. This VLC configuration ensures proper transcoding and streaming of video frames for the Camera Controller to process.

If you have a different configuration or want to customize the VLC settings, you may need to update the SOURCE_PATH in the config.json file accordingly.

## Build and Run on Docker container
   ```bash
  rm /app/cmake-build-debug/CMakeCache.txt
  /usr/bin/cmake -DCMAKE_BUILD_TYPE=Debug -G "CodeBlocks - Unix Makefiles" -S /app -B /app/cmake-build-debug   
  /usr/bin/cmake --build /app/cmake-build-debug --target CPP -- -j 14
     cd /app && /app/cmake-build-debug/CPP
   
   ```

# Camera Controller Overview

The Camera Controller is a C++ application responsible for managing video streams, processing frames, and controlling the system based on commands received via MQTT. The following provides a summary of its functionality based on the provided code snippets:

## `main.cpp` (Cam Controller Entry Point):

1. **Initialization:**
   - Initializes the logger and a video stream processor.
   - Retrieves the current working directory and prints it.

2. **Signal Handling:**
   - Registers a signal handler for handling SIGINT (e.g., Ctrl+C) to gracefully shutdown the application.

3. **Configuration Loading:**
   - Loads configuration settings from a file using a utility function.

4. **HTTP Handler Initialization:**
   - Initializes an HTTP handler (`frameUploader`) with authentication parameters.

5. **MQTT Handler Initialization:**
   - Initializes an MQTT handler (`commandListener`) with connection parameters.
   - Defines a callback function for executing commands received from MQTT.
   - Subscribes to a specific MQTT topic with the defined callback.

6. **Video Stream Processor Initialization:**
   - Sets the video source index.
   - Initializes a Kafka producer for sending video frames.

7. **Thread Creation:**
   - Creates threads for the MQTT subscriber and the video frame processor.
   - Optionally sets an alternate source for video frames.

8. **Thread Joining:**
   - Waits for the threads to finish.

## `frame_processor.cpp` (Video Frame Processor):

1. **Video Input Opening:**
   - Opens the video input based on specified sources (source index or alternate source).
   - Sets the video codec to MJPEG.

2. **Frame Processing and Callback:**
   - Continuously reads frames from the video input.
   - Applies a callback function to encode and process video frames.
   - Checks if publishing is enabled and if the specified frame rate has passed.

3. **Control Mechanism:**
   - Listens for control commands (START/STOP) from MQTT to control frame processing.

4. **Source Configuration:**
   - Allows setting the video source index or an alternate source.

5. **Shutdown:**
   - Gracefully shuts down the frame processor.

## `mqtt.cpp` (MQTT Handler):

1. **MQTT Callback Implementation:**
   - Implements callback methods for MQTT connection events, message arrival, and message delivery completion.

2. **MQTT Handler Implementation:**
   - Connects to an MQTT broker with specified credentials.
   - Subscribes to an MQTT topic with a specified callback function.
   - Publishes messages to an MQTT topic.
   - Gracefully shuts down the MQTT handler.

## `frame_uploader.cpp` (HTTP Handler):

1. **HTTP Handler Implementation:**
   - Initializes an HTTP handler with necessary parameters and CURL setup.
   - Performs user login with a provided camId and credentials.
   - Cleans up resources during destruction.

## `KafkaProducer.cpp` (Kafka Producer Implementation):

1. **Kafka Producer Implementation:**
   - Creates a Kafka producer with specified brokers.
   - Sets up delivery report callback and bootstrap servers.
   - Sends messages to a specified Kafka topic.
   - Waits for pending message deliveries and gets the length of the outgoing message queue.
   - Cleans up resources during destruction.

## Summary:

The Camera Controller orchestrates the capture, processing, and distribution of video frames. It interacts with external systems via MQTT, sends frames to Kafka, and performs user authentication via HTTP. The system is designed to be modular, allowing easy integration of new components for video processing and communication.


# Processing Server Overview

The processing server is a C++ application designed for handling video frame data within a larger system. The following summarizes its functionality based on the provided code snippets:

## `main.cpp` (Processing Server Entry Point):

1. **Working Directory and Configuration:**
   - Retrieves and prints the current working directory.
   - Loads configuration settings from a file using `my_namespace::utility::loadConfigFromFile`.
   - Initializes a Kafka producer with broker information from the configuration.

2. **MinIOUploader and ProcessMessage Lambda:**
   - Creates a `MinIOUploader` instance with MinIO endpoint and bucket name from the configuration.
   - Defines a lambda function `processMessage` for processing Kafka messages:
      - Parses the received Kafka message using `parseMessage`.
      - Converts the timestamp to Unix time in nanoseconds.
      - Generates an image name based on timestamp and camera ID.
      - Applies detection using `video::convertDetect`.
      - Uploads the marked image to MinIO using `minioUploader.uploadImage`.
      - Sends detection results to other services using `sendResultToServices`.

3. **Kafka Consumer Initialization:**
   - Creates a Kafka consumer with broker, group ID, and topic information from the configuration.
   - Starts consuming messages and invokes the `processMessage` lambda for each received message.

## `KafkaProducer` (Kafka Producer Implementation):

- Initializes a Kafka producer using librdkafka.
- Provides a method `sendMessage` to send messages to a specified Kafka topic.
- Utilizes a delivery report callback to log the success or failure of message delivery.
- Includes a method `wait_for_pending_delivery` to flush pending deliveries.

## `MqttHandler` (MQTT Handler Implementation):

- Implements an MQTT client using the Eclipse Paho C++ library.
- Connects to an MQTT broker, subscribes to topics, and publishes messages.
- Provides callback functions for connection events, message reception, and message delivery.

## `HttpHandler` (HTTP Handler Implementation):

- Implements an HTTP handler using libcurl.
- Initializes libcurl with server URL, timeout, and authentication parameters.
- Provides a callback function for writing received data during CURL requests.
- Performs login to the server using a specified endpoint, username, and password.

## `video::convertDetect` (Video Processing):

- Converts received image data into a cv::Mat using OpenCV.
- Applies a people detection algorithm using Histogram of Oriented Gradients (HOG) with OpenCV.
- Marks detected people in the image and encodes the marked image back to image data.

## Summary:

The processing server consumes video frame data from a Kafka topic, processes the frames (including people detection), uploads marked images to MinIO, and sends detection results to other services using Kafka. It integrates with MQTT and HTTP for additional communication. The code is modular and follows a structured approach to handle different aspects of the video processing pipeline.
