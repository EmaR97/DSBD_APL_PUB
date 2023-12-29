#ifndef HTTP_HANDLER_H
#define HTTP_HANDLER_H

#include <string>
#include <vector>
#include <opencv2/core/hal/interface.h>
#include <curl/curl.h>
#include "../utility/logger.h"
#include "../message/framedata.pb.h"

namespace my_namespace::sender {

    class HttpHandler {
    public:
        void initialize(const std::string &serverUrl, long timeoutMs, const std::string &username,
                        const std::string &password, const std::string &login_endpoint);

        ~HttpHandler();

        std::basic_string<char> performLogin(const std::string &camId);

    private:
        CURL *curl;
        my_namespace::utility::Logger *logger_;
        long timeoutMs_;
        std::string username_;
        std::string password_;
        std::string loginUrl;
        std::string camId;
        curl_slist *cookies = nullptr;


    };
}
#endif // HTTP_HANDLER_H

