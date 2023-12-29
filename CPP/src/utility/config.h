
#ifndef CAMCONTROLLER_CONFIG_H
#define CAMCONTROLLER_CONFIG_H
#include "json.hpp"
#include <string>
namespace my_namespace::utility {
    nlohmann::basic_json<> loadConfigFromFile(const char *configPath);

    void saveConfigToFile(const nlohmann::json& json_data, const char *configPath);
}
#endif //CAMCONTROLLER_CONFIG_H
