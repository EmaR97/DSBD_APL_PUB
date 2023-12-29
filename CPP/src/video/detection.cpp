//
// Created by emanuele on 18/12/23.
//

#include "detection.h"
#include "../utility/logger.h"

using namespace cv;
using namespace ml;
using namespace std;
namespace my_namespace::video {
    auto &logger = utility::Logger::getInstance();

    bool convert_detect(std::vector<uchar> &imageData) {
        // Decode the image data into a cv::Mat using OpenCV
        cv::Mat receivedImage = cv::imdecode(imageData, cv::IMREAD_UNCHANGED);
        imageData.clear();
        // Additional processing or handling of the received image can be done here
        cv::Mat markedFrame;
        bool detected = detectPeople(receivedImage, markedFrame, 0);
        if (!cv::imencode(".jpg", markedFrame, imageData)) {
            logger << utility::LogLevel::ERROR  << "Error encoding frame" << std::endl;
            exit(EXIT_FAILURE);
        }
        return detected;
    }

// Detect people in the input image and mark the detections
    bool detectPeople(const Mat &inputImage, Mat &markedImage, double detectionThreshold) {
        HOGDescriptor hog;
        hog.setSVMDetector(HOGDescriptor::getDefaultPeopleDetector());

        // Resize the input image if needed
        Mat img = inputImage.clone();
        resize(img, img, Size(img.cols * 2, img.rows * 2));

        vector<Rect> found;
        vector<double> weights;

        // Use HOG detector to find people in the image
        hog.detectMultiScale(img, found, weights, detectionThreshold, Size(8, 8), Size(32, 32), 1.05, 2);

        if (found.empty()) {
            // No people found
            markedImage = inputImage.clone();
            return false;
        }

        // People found, mark the detections on the image
        markedImage = inputImage.clone();
        for (size_t i = 0; i < found.size(); i++) {
            Rect r = found[i];
            rectangle(markedImage, found[i], Scalar(0, 0, 255), 3);
            stringstream temp;
            temp << weights[i];
            putText(markedImage, temp.str(), Point(found[i].x, found[i].y + 50), FONT_HERSHEY_SIMPLEX, 1,
                    Scalar(0, 0, 255));
        }

        return true;
    }
}