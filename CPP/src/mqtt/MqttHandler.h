//
// Created by emanu on 02/12/2023.
//

#ifndef CAMCONTROLLER_MQTTHANDLER_H
#define CAMCONTROLLER_MQTTHANDLER_H

#include <iostream>
#include <mqtt/async_client.h>
#include <functional>
#include <mutex>
#include <condition_variable>
#include "../utility/logger.h"

namespace my_namespace::mqtt_ {
    class MqttCallback : public virtual mqtt::callback {
    public:
        void connected(const mqtt::string &cause) override;

        void connection_lost(const std::string &cause) override;

        void message_arrived(mqtt::const_message_ptr msg) override;

        void delivery_complete(mqtt::delivery_token_ptr token) override;

        void setCallbackFunction(const std::function<void(std::string)> &callback);

        void setClientId(const std::string &id);


    private:
        std::string clientId;
        std::function<void(std::string)> callbackFunction = nullptr;
    };

    class MqttHandler {
    public:

        MqttHandler(const std::string &serverAddress, const std::string &clientId);

        ~MqttHandler();

        void connect(const std::string &username, const std::string &password);

        void subscribe(const std::string &topic, const std::function<void(std::string)> &callback, int qos = 0);

        void publish(const std::string &topic, const std::string &message);

        void shutdown();


    private:
        mqtt::async_client client;
        MqttCallback cb;
        std::condition_variable exitCondition;
        std::mutex exitMutex;

    };
}
#endif //CAMCONTROLLER_MQTTHANDLER_H
