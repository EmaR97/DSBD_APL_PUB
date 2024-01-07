#ifndef CRONOMETER_H
#define CRONOMETER_H

#include <iostream>
#include <chrono>
#include <mutex>

class Chronometer {
public:
    Chronometer();

    void startWorking();

    void stopWorking();

    double checkAndReset();

private:
    std::mutex mutex;
    std::chrono::high_resolution_clock::time_point lastCheck;
    std::chrono::high_resolution_clock::time_point workingStart;
    long long totalDuration;
    long long workingDuration;
    bool working;
};

#endif // CRONOMETER_H
