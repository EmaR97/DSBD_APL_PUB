//
// Created by emanuele on 18/12/23.
//

#ifndef CAMCONTROLLER_DETECTION_H
#define CAMCONTROLLER_DETECTION_H

#include <opencv2/opencv.hpp>
#include <opencv2/opencv.hpp>

namespace my_namespace::video {

    bool detectPeople(const cv::Mat &inputImage, cv::Mat &markedImage, double detectionThreshold = 0.0);

    bool convertDetect(std::vector<uchar> &imageData);
}
#endif //CAMCONTROLLER_DETECTION_H
