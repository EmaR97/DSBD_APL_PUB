#ifndef LOGGER_H
#define LOGGER_H

#include <iostream>
#include <fstream>
#include <ctime>
#include <iomanip> // Added for std::put_time
#include <chrono>

namespace my_namespace::utility {

    // Enum to represent different log levels
    enum class LogLevel {
        ERROR, WARNING, INFO
    };

    // Logger class for logging messages to a file and console
    class Logger {
    public:
        static Logger& getInstance() {
            static Logger instance; // Guaranteed to be initialized only once.
            return instance;
        }
        // Constructor: opens the log file
        Logger() {
            openLogFile();
        }

        // Destructor: closes the log file
        ~Logger() {
            closeLogFile();
        }

        // Overloaded << operator to write data to both log file and console
        template<typename T>
        Logger &operator<<(const T &data) {
            writeToLog(data);
            writeToConsole(data);
            return *this;
        }

        // Overloaded << operator to handle manipulators (e.g., std::endl)
        Logger &operator<<(std::ostream &(*manipulator)(std::ostream &)) {
            writeToLog(manipulator);
            writeToConsole(manipulator);
            return *this;
        }

        // Overloaded << operator to print log level headers
        Logger &operator<<(const LogLevel &level) {
            printLogHeader(level);
            return *this;
        }

    private:
        // Log file stream
        std::ofstream logFile;

        // Opens the log file for appending
        void openLogFile() {
            logFile.open("logfile.txt", std::ios::app);
            if (!logFile.is_open()) {
                std::cerr << "Unable to open log file for writing." << std::endl;
            }
        }

        // Closes the log file
        void closeLogFile() {
            if (logFile.is_open()) {
                logFile.close();
            }
        }

        // Writes data to the log file
        template<typename T>
        void writeToLog(const T &data) {
            if (logFile.is_open()) {
                logFile << data;
            }
        }

        // Writes manipulators to the log file
        void writeToLog(std::ostream &(*manipulator)(std::ostream &)) {
            if (logFile.is_open()) {
                manipulator(logFile);
            }
        }

        // Writes data to the console
        template<typename T>
        void writeToConsole(const T &data) {
            std::cout << data;
        }

        // Writes manipulators to the console
        void writeToConsole(std::ostream &(*manipulator)(std::ostream &)) {
            manipulator(std::cout);
        }

        // Prints log header with timestamp and log level
        void printLogHeader(const LogLevel &level) {
            printTimestamp();
            printLogLevel(level);
        }

        // Prints timestamp in the format: [YYYY/MM/DD HH:MM:SS.xxx]
        void printTimestamp() {
            auto now = std::chrono::system_clock::now();
            auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(now.time_since_epoch()) % 1000;

            std::time_t currentTime = std::chrono::system_clock::to_time_t(now);
            std::tm localTime = *std::localtime(&currentTime);

            char buffer[24]; // Increased the size to accommodate milliseconds
            std::strftime(buffer, sizeof(buffer), "%Y/%m/%d %H:%M:%S", &localTime);

            *this << "[" << buffer << "." << std::setfill('0') << std::setw(3) << ms.count() << "]";
        }

        // Prints log level based on LogLevel enum
        void printLogLevel(const LogLevel &level) {
            switch (level) {
                case LogLevel::ERROR:
                    *this << "[ERROR] ";
                    break;
                case LogLevel::WARNING:
                    *this << "[WARNING] ";
                    break;
                case LogLevel::INFO:
                    *this << "[INFO] ";
                    break;
                default:
                    *this << "[UNKNOWN] ";
            }
        }
    };
}

#endif // LOGGER_H
