// mqtt.cpp
#include "MqttHandler.h"

#include <utility>

using my_namespace::utility::LogLevel;
namespace my_namespace::mqtt_ {
    auto &logger = utility::Logger::getInstance();

    // Implementation of MqttCallback methods

    // Called when connected to the MQTT broker
    void MqttCallback::connected(const mqtt::string &cause) {
        logger << LogLevel::INFO << clientId << ": Connected to MQTT broker" << std::endl;
    }

    // Called when the connection to the MQTT broker is lost
    void MqttCallback::connection_lost(const std::string &cause) {
        logger << LogLevel::INFO << clientId << ": Connection lost: " << cause << std::endl;
    }

    // Called when a new MQTT message is received
    void MqttCallback::message_arrived(mqtt::const_message_ptr msg) {
        std::string payload = msg->get_payload_str();
        logger << LogLevel::INFO << clientId << ": Received message: " << payload << std::endl;

        // Call the user-defined callback function if set
        if (messageHandler) {
            messageHandler(payload);
        } else {
            logger << LogLevel::WARNING << clientId << ": message handler function not set." << std::endl;
        }
    }

    // Called when an MQTT message delivery is complete
    void MqttCallback::delivery_complete(mqtt::delivery_token_ptr token) {
        logger << LogLevel::INFO << clientId << ": Message delivered" << std::endl;
    }

    // Setter for the user-defined callback function
    void MqttCallback::setMessageHandler(const std::function<void(std::string)> &callback) {
        messageHandler = callback;
    }

    // Setter for the MQTT client ID
    void MqttCallback::setClientId(const std::string &id) {
        MqttCallback::clientId = id;
    }

    // Implementation of MqttHandler methods

    // Constructor for MqttHandler
    MqttHandler::MqttHandler(const std::string &serverAddress, const std::string &clientId, std::string username,
                             std::string password) : client(serverAddress, clientId), username_(std::move(username)),
                                                     password_(std::move(password)) {
    }

    // Destructor for MqttHandler
    MqttHandler::~MqttHandler() {
        // Disconnect MQTT client during destruction
        client.disconnect();
    }

    // Connect to the MQTT broker with specified username and password
    void MqttHandler::connect() {
        mqtt::connect_options conn_opts;
        conn_opts.set_automatic_reconnect(10, 40);
        conn_opts.set_keep_alive_interval(20);
        conn_opts.set_clean_session(true);
        conn_opts.set_user_name(username_);
        conn_opts.set_password(password_);
        cb.setClientId(client.get_client_id());
        while (true)
            try {
                // Attempt to connect to the MQTT broker
                logger << LogLevel::INFO << client.get_client_id() << ": Connecting to MQTT broker for ..."
                       << std::endl;
                client.connect(conn_opts)->wait();
                return;
            } catch (const mqtt::exception &exc) {
                logger << LogLevel::ERROR << client.get_client_id() << ": Unable to connect to MQTT broker: "
                       << exc.what() << std::endl;
                logger << LogLevel::INFO << client.get_client_id() << ": Trying again " << 10 << " in second  ..."
                       << std::endl;
                std::this_thread::sleep_for(std::chrono::seconds(10));
            }
    }

    // Subscribe to an MQTT topic with the specified callback and QoS level
    void MqttHandler::subscribe(const std::string &topic, const std::function<void(std::string)> &callback, int qos) {
        cb.setMessageHandler(callback);
        client.set_callback(cb);
        topic_ = topic;
        // Subscribe to the specified topic
        logger << LogLevel::INFO << client.get_client_id() << ": Subscribing to " << topic << std::endl;
        client.subscribe(topic, qos);

        logger << LogLevel::INFO << client.get_client_id() << ": Listening ... " << std::endl;
    }

    // Publish a message to an MQTT topic
    [[maybe_unused]] void MqttHandler::publish(const std::string &topic, const std::string &message) {
        client.set_callback(cb);

        // Publish the message to the specified topic
        logger << LogLevel::INFO << client.get_client_id() << ": Publishing to " << topic << std::endl;
        client.publish(topic, message.c_str(), message.length(), 0, false)->wait();
    }

    // Shutdown the MQTT handler
    void MqttHandler::shutdown() {
        // Log an exit message
        logger << LogLevel::INFO << client.get_client_id() << ": Exiting ... " << std::endl;
        // Perform any cleanup or additional shutdown steps if needed
        client.unsubscribe(topic_);
        client.disconnect();
        logger << LogLevel::INFO << client.get_client_id() << ": Exiting Completed " << std::endl;
    }
}
