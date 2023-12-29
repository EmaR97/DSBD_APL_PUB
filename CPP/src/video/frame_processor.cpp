// frame_processor.cpp
#include "frame_processor.h"

using namespace std;
using namespace cv;
using namespace ml;
using my_namespace::utility::LogLevel;

namespace my_namespace::video {
    auto &logger = utility::Logger::getInstance();

    // Open the video input based on the specified sources
    void FrameProcessor::openVideoInput() {
        bool open = false;

        if (sourceIndex_ != nullptr)
            open = videoCapture_->open(*sourceIndex_, CAP_V4L2);
        else if (alternateSource_ != nullptr)
            open = videoCapture_->open(*alternateSource_, CAP_FFMPEG);
        else
            logger << LogLevel::ERROR << "No source video defined" << endl;

        // Set video codec to MJPEG
        videoCapture_->set(CAP_PROP_FOURCC, VideoWriter::fourcc('M', 'J', 'P', 'G'));

        if (!open) {
            logger << LogLevel::ERROR << "Error opening camera" << endl;
            exit(-1);
        }
        logger << LogLevel::INFO << "Camera opened correctly" << endl;
    }

    // Process frames and apply a callback function
    void FrameProcessor::processFramesAndApplyCallback(int frameRate, const String &imageFormat,
                                                       const FrameCallback &callback) {
        while (true) {
            openVideoInput();
            Mat currentFrame;

            long previousTime = 0;
            logger << LogLevel::INFO << "Start Sending" << endl;

            while (true) {
                {
                    lock_guard<mutex> lock(exitMutex);
                    // Check if shutdown is requested
                    if (shutdownRequested) {
                        // Shutdown requested, perform cleanup and exit
                        logger << LogLevel::INFO << "Thread is shutting down." << endl;
                        return;
                    }
                }

                // Read a frame from the video input
                if (!videoCapture_->read(currentFrame)) {
                    logger << LogLevel::ERROR << "Error reading frame" << endl;
                    continue;
                }


                // Check if the specified frame rate has passed
                if (!hasElapsedTimePassed(previousTime, frameRate)) {
                    continue;
                }

                // Check if publishing is enabled
                if (!shouldPublish())
                    break;

                // Encode the marked frame and apply the callback
                vector<uchar> buffer;
                if (!imencode(imageFormat, currentFrame, buffer)) {
                    logger << LogLevel::ERROR << "Error encoding frame" << endl;
                    continue;
                }
                callback(buffer);
            }
            videoCapture_->release();

            {
                unique_lock<mutex> lock(controlMutex);
                // Wait for the control condition (START/STOP) to be signaled
                controlCondition.wait(lock, [this] { return publishEnabled; });
            }
        }
    }

    // Check if a specified time has elapsed based on a frame rate
    bool FrameProcessor::hasElapsedTimePassed(long &previousTime, double rate) {
        const long now = chrono::duration_cast<chrono::milliseconds>(
                chrono::system_clock::now().time_since_epoch()).count();
        if (const auto timeElapsed = static_cast<double>(now - previousTime); timeElapsed > (1000.0 / rate)) {
            previousTime = now;
            return true;
        }
        return false;
    }

    // Start or stop processing based on a received payload
    void FrameProcessor::startAndStopProcessing(const string &payload) {
        lock_guard<mutex> lock(controlMutex);
        if (payload == "STOP") {
            publishEnabled = false;
        } else if (payload == "START") {
            publishEnabled = true;
            // Notify the processing loop to continue
            controlCondition.notify_one();
        }
    }

    // Check if publishing is enabled
    bool FrameProcessor::shouldPublish() {
        lock_guard<mutex> lock(controlMutex);
        return publishEnabled;
    }

  
    // Set the source index for video input
    void FrameProcessor::setSourceIndex(int *sourceIndex) {
        sourceIndex_ = sourceIndex;
    }

    // Set the alternate source for video input (file or stream)
    [[maybe_unused]] void FrameProcessor::setAlternateSource(string *alternateSource, bool isFile) {
        if (isFile) {
            if (filesystem::exists(*alternateSource)) {
                logger << LogLevel::INFO << "File exists: " + *alternateSource << endl;
            } else {
                logger << LogLevel::ERROR << "File does not exist: " + *alternateSource << endl;
                exit(-1);
            }
        }
        alternateSource_ = alternateSource;
    }


    // Shutdown the frame processor
    void FrameProcessor::shutdown() {
        // Perform any cleanup or additional shutdown steps if needed
        lock_guard<mutex> lock(exitMutex);
        shutdownRequested = true; // Set the shutdown flag
    }

} // namespace my_namespace::video
