//
// Created by emanuele on 18/12/23.
//

#ifndef CAMCONTROLLER_MINIOUPLOADER_H
#define CAMCONTROLLER_MINIOUPLOADER_H

#include <iostream>
#include <curl/curl.h>
#include <opencv2/opencv.hpp>
#include <vector>
#include <openssl/hmac.h>
#include <openssl/sha.h>
#include <aws/core/Aws.h>
#include <aws/s3/S3Client.h>
#include <aws/s3/model/PutObjectRequest.h>
#include <aws/core/auth/AWSCredentials.h>
#include <aws/core/utils/memory/stl/AWSString.h>
#include <curl/curl.h>
#include <iostream>
#include <aws/core/http/HttpClient.h>
#include <opencv2/opencv.hpp>

namespace my_namespace::sender {

    class MinIOUploader {
    public:
        MinIOUploader(const std::string &minioEndpoint, std::string bucketName, std::string keyId,
                      std::string keySecret);

        ~MinIOUploader();

        void uploadImage(const std::basic_string<char> &objectName, std::vector<uchar> &imageData);

    private:
        static size_t readCallback(void *ptr, size_t size, size_t nmemb, void *userdata);

        Aws::SDKOptions options;
        std::shared_ptr<Aws::S3::S3Client > s3_client;
        CURL *curl;
        std::basic_string<char> bucketName;

        void initilizeMinioClient(const char *keyId, const char *keySecret, const char  *endpoint);
    };

}
#endif //CAMCONTROLLER_MINIOUPLOADER_H
