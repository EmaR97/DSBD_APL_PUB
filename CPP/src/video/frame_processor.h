// frame_processor.h

#ifndef FRAME_PROCESSOR_H
#define FRAME_PROCESSOR_H

#include <filesystem>
#include <opencv2/opencv.hpp>
#include "../http/HttpHandler.h"
#include "../utility/logger.h"
#include <thread>
#include "../mqtt/MqttHandler.h"

namespace my_namespace::video {


    class FrameProcessor {
    public:
        // Constructor with member initialization
        FrameProcessor() : videoCapture_(new cv::VideoCapture()), publishEnabled(true) {}

        // Destructor to release the allocated videoCapture_ object
        ~FrameProcessor() {
            delete videoCapture_;
        }

        using FrameCallback = std::function<void(const std::vector<uchar> &)>;

        void processFramesAndApplyCallback(int frameRate, const cv::String &imageFormat, const FrameCallback &callback);

        void startAndStopProcessing(const std::string &payload);


        void setSourceIndex(int *sourceIndex);

        [[maybe_unused]] void setAlternateSource(std::string *alternateSource, bool isFile = false);

        void shutdown();

    private:
        cv::VideoCapture *videoCapture_;
        int *sourceIndex_ = nullptr;
        std::mutex exitMutex;
        std::condition_variable exitCondition;
        std::string *alternateSource_ = nullptr;
        bool publishEnabled;
        std::mutex controlMutex;
        std::condition_variable controlCondition;
        bool shutdownRequested = false;

        bool shouldPublish();

        static bool hasElapsedTimePassed(long &previousTime, double rate);

        void openVideoInput();

    };

} // namespace cam_controller

#endif // FRAME_PROCESSOR_H
