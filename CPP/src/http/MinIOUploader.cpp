#include "MinIOUploader.h"
#include "../utility/logger.h"

#include <utility>


using namespace my_namespace;
namespace my_namespace::sender {

    utility::Logger &logger = utility::Logger::getInstance();

    MinIOUploader::MinIOUploader(const std::string &minioEndpoint, std::string bucketName, const std::string& keyId,
                                 const std::string& keySecret) : bucketName(std::move(bucketName)) {
        curl = curl_easy_init();
        if (!curl) {
            logger << utility::LogLevel::INFO << "Failed to initialize curl!" << std::endl;
            exit(EXIT_FAILURE);
        }

        initilizeMinioClient(keyId.c_str(), keySecret.c_str(), minioEndpoint.c_str());
    }

    void MinIOUploader::initilizeMinioClient(const char *keyId, const char *keySecret, const char *endpoint) {
        Aws::InitAPI(options);
        Aws::Client::ClientConfiguration clientConfig;
        logger << utility::LogLevel::INFO << "endpoint: " << endpoint << std::endl;

        clientConfig.endpointOverride = Aws::String(endpoint);
        clientConfig.scheme = Aws::Http::Scheme::HTTP;

        Aws::Auth::AWSCredentials credentials;
        credentials.SetAWSAccessKeyId(Aws::String(keyId));
        credentials.SetAWSSecretKey(Aws::String(keySecret));
        s3_client = static_cast<const std::shared_ptr<Aws::S3::S3Client>>(new Aws::S3::S3Client(credentials,clientConfig, Aws::Client::AWSAuthV4Signer::PayloadSigningPolicy(),false));
    }

    MinIOUploader::~MinIOUploader() {
        Aws::ShutdownAPI(options);
        if (curl) {
            curl_easy_cleanup(curl);
        }
    }


    void MinIOUploader::uploadImage(const std::basic_string<char> &objectName, std::vector<uchar> &imageData) const {
        //TODO create bucket if missing

        // Construct the upload URL
        Aws::String presignedUrlGet = s3_client->GeneratePresignedUrl(bucketName, objectName,
                                                                      Aws::Http::HttpMethod::HTTP_PUT, 300);
        logger << utility::LogLevel::INFO << "presignedUrlGet: " << presignedUrlGet << std::endl;
        // Set the upload URL and enable uploading
        curl_easy_setopt(curl, CURLOPT_URL, presignedUrlGet.c_str());
        curl_easy_setopt(curl, CURLOPT_UPLOAD, 1L);

        // Set the callback function for reading data during upload
        curl_easy_setopt(curl, CURLOPT_READFUNCTION, readCallback);
        curl_easy_setopt(curl, CURLOPT_READDATA, &imageData);
        curl_easy_setopt(curl, CURLOPT_INFILESIZE_LARGE, static_cast<curl_off_t>(imageData.size()));

        // Perform the upload
        CURLcode result = curl_easy_perform(curl);
        if (result != CURLE_OK) {
            logger << utility::LogLevel::ERROR << "curl_easy_perform() failed: " << curl_easy_strerror(result)
                   << std::endl;
            exit(EXIT_FAILURE);
        }
    }

    size_t MinIOUploader::readCallback(void *ptr, size_t size, size_t nmemb, void *userdata) {
        auto *data = static_cast<std::vector<uchar> *>(userdata);
        size_t dataSize = size * nmemb;
        ptrdiff_t bytesToCopy = static_cast<std::vector<int>::difference_type>(std::min(dataSize, data->size()));

        // Copy data to the buffer and remove copied data from the vector
        std::copy(data->begin(), data->begin() + bytesToCopy, static_cast<uchar *>(ptr));
        data->erase(data->begin(), data->begin() + bytesToCopy);

        return bytesToCopy;
    }
}