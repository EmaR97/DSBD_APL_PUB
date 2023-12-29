// frame_uploader.cpp
#include "HttpHandler.h"
#include <iostream>

using my_namespace::utility::LogLevel;
namespace my_namespace::sender {
    utility::Logger &logger = utility::Logger::getInstance();

// Initialize the HttpHandler with necessary parameters and CURL setup
    void HttpHandler::initialize(const std::string &serverUrl, long timeoutMs, const std::string &username,
                                 const std::string &password, const std::string &loginEndpoint) {

        // Store configuration parameters
        username_.assign(username);
        password_.assign(password);
        timeoutMs_ = timeoutMs;

        // Initialize CURL
        curl_global_init(CURL_GLOBAL_DEFAULT);
        curl = curl_easy_init();
        if (!curl) {
            logger << LogLevel::ERROR << "Error initializing CURL" << std::endl;
            std::exit(-1);
        }

        // Construct the login and upload URLs
        loginUrl = serverUrl + loginEndpoint;

        // Additional setup for the CURL handle can go here
    }

// Destructor to clean up resources
    HttpHandler::~HttpHandler() {
        curl_easy_cleanup(curl);
        curl_global_cleanup();
    }

// Callback function for writing received data during CURL requests
    size_t WriteCallback(const char *contents, size_t size, size_t nmemb, void *userp) {
        static_cast<std::string *>(userp)->append(contents, size * nmemb);
        return size * nmemb;
    }

// Perform login and return the assigned camId
    std::basic_string<char> HttpHandler::performLogin(const std::string &camId_) {
        std::string response;
        camId = camId_;

        // Set CURL options for the login request
        curl_easy_setopt(curl, CURLOPT_URL, loginUrl.c_str());
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);
        curl_easy_setopt(curl, CURLOPT_COOKIESESSION, 1);
        curl_easy_setopt(curl, CURLOPT_COOKIEFILE, "");
        curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, timeoutMs_);
        curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, timeoutMs_);
        curl_easy_setopt(curl, CURLOPT_TCP_KEEPALIVE, 1L);
        curl_easy_setopt(curl, CURLOPT_TIMEOUT, 30L);

        // Set up the login form data
        std::string postData = "cam_id=" + camId_ + "&password=" + password_;
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, postData.c_str());

        // Perform the CURL login request
        CURLcode res = curl_easy_perform(curl);

        if (res != CURLE_OK) {
            logger << LogLevel::ERROR << std::string("Login failed: ") + curl_easy_strerror(res) << std::endl;
            exit(1);
        }

        int code;
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
        if (code != 200) {
            logger << LogLevel::ERROR << "Login failed, code: " << code << " response: " << response << std::endl;
            exit(1);
        }
        logger << LogLevel::INFO << "Login successful." << std::endl;

        if (!response.empty() && code == 202) {
            camId.assign(response);
            return camId;
        }
        // Extract session (cookie) information from response headers
        res = curl_easy_getinfo(curl, CURLINFO_COOKIELIST, &cookies);
        if (res != CURLE_OK) {
            logger << LogLevel::ERROR << std::string("Failed to get cookies: ") + curl_easy_strerror(res) << std::endl;
        }
        struct curl_slist *cookieItem = cookies;
        while (cookieItem != nullptr) {
            logger << LogLevel::INFO << "Cookie: " << cookieItem->data << std::endl;
            cookieItem = cookieItem->next;
        }
        return camId;

    }
}