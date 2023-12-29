//
// Created by emanuele on 18/12/23.
//


#include <miniocpp/client.h>
#include <opencv2/core/mat.hpp>
#include <opencv2/imgcodecs.hpp>

// Function to upload image to MinIO
int uploadImageToMinIO(const cv::Mat &image, const std::string &bucketName, const std::string &objectName,
                       minio::s3::Client minioClient) {
// Initialize MinIO client

// Convert cv::Mat to std::vector<uint8_t>
    std::vector<uint8_t> imageData;
    cv::imencode(".png", image, imageData);

// Upload image to MinIO


    minio::s3::UploadObjectArgs args;
    args.bucket = bucketName;
    args.object = reinterpret_cast<const char *>(imageData.data());
    args.object = imageData.size();
    args.filename = objectName;
    args.content_type = "image/png";
    minio::s3::UploadObjectResponse resp = minioClient.UploadObject(args);
    if (!resp) {
        std::cout << "unable to upload object; " << resp.Error() << std::endl;
        return 0;
    }
    return 1;
}


int main() {


    minio::s3::BaseUrl base_url("YOUR_MINIO_ENDPOINT");
    minio::creds::StaticProvider provider("YOUR_ACCESS_KEY", "YOUR_SECRET_KEY");
    minio::s3::Client minioClient(base_url, &provider);
// store the marked image
    std::string bucketName = "your_bucket_name";
    std::string objectName = "path/to/store/markedImage.png";
    cv::Mat markedImage;
    uploadImageToMinIO(markedImage, bucketName, objectName, minioClient);
    return 0;


};


