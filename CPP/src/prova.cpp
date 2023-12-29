#include <opencv2/opencv.hpp>
#include "http/MinIOUploader.h"

void test(const char *src, const char *dest, const MinIOUploader &minioUploader);

int main() {
    const char *minioEndpoint = "http://localhost:9000";
    const char *bucketName = "your-bucket";
    const char *src = "your_image.jpg";
    MinIOUploader minioUploader(minioEndpoint, bucketName);

    test(src, "a/image1.jpg", minioUploader);
    test(src, "c/image1.jpg", minioUploader);

    return 0;
}

void test(const char *src, const char *dest, const MinIOUploader &minioUploader) {
    std::vector<uchar> imageData;
    if (!cv::imencode(".jpg",  cv::imread(src), imageData)) {
        std::cerr << "Error encoding frame" << std::endl;
        exit(EXIT_FAILURE);
    }
    // Example usage for uploading multiple images
    minioUploader.uploadImage(dest, imageData);
    // Add more images as needed
}
