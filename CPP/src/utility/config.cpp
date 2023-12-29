//
// Created by emanu on 06/12/2023.
//
#include "json.hpp"
#include <iostream>
#include <fstream>

#include "config.h"
#include <unistd.h>

namespace my_namespace::utility {
    nlohmann::basic_json<> loadConfigFromFile(const char *configPath) {
        std::ifstream config_file(configPath);
        if (config_file.is_open()) {
            try {
                nlohmann::json json_data;
                config_file >> json_data;
                return json_data;
            } catch (const std::exception &e) {
                std::cerr << "Error loading configuration: " << e.what() << std::endl;
                exit(1);
            }
        } else {
            std::cerr << "Unable to open config file." << std::endl;
            exit(1);
        }
    }


    void saveConfigToFile(const nlohmann::json &json_data, const char *configPath) {
        std::ofstream config_file(configPath);
        if (config_file.is_open()) {
            config_file << std::setw(4) << json_data;  // Pretty print with indentation
            std::cout << "Configuration saved to file." << std::endl;
        } else {
            std::cerr << "Unable to open config file for writing." << std::endl;
        }
    }
}