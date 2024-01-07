#include "chronometer.h"

Chronometer::Chronometer() : totalDuration(0), workingDuration(0), working(false),
                             lastCheck(std::chrono::high_resolution_clock::now()) {}

void Chronometer::startWorking() {
    std::lock_guard<std::mutex> lock(mutex);
    if (!working) {
        working = true;
        workingStart = std::chrono::high_resolution_clock::now();
    }
}

void Chronometer::stopWorking() {
    std::lock_guard<std::mutex> lock(mutex);
    if (working) {
        auto now = std::chrono::high_resolution_clock::now();
        workingDuration += std::chrono::duration_cast<std::chrono::nanoseconds>(now - workingStart).count();
        working = false;
    }
}

double Chronometer::checkAndReset() {
    std::lock_guard<std::mutex> lock(mutex);

    auto now = std::chrono::high_resolution_clock::now();

    totalDuration += std::chrono::duration_cast<std::chrono::nanoseconds>(now - lastCheck).count();

    if (working) {
        workingDuration += std::chrono::duration_cast<std::chrono::nanoseconds>(now - workingStart).count();
        workingStart = now;
    }

    double workingTimePercentage = 0;
    if (totalDuration > 0) {
        workingTimePercentage = (static_cast<double>(workingDuration) / static_cast<double>(totalDuration)) * 100;
    }

    workingDuration = 0;
    totalDuration = 0;
    lastCheck = now;
    return workingTimePercentage;
}


