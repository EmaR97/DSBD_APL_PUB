#include "utility.h"

using namespace std::chrono;
using namespace google::protobuf;
namespace my_namespace::utility {

    Timestamp timestamp_ns() {
        // Get the current time point
        const auto current_time = system_clock::now();

        // Convert the time point to nanoseconds since the epoch
        const long nano_seconds = time_point_cast<nanoseconds>(current_time).time_since_epoch().count();

        // Convert nanoseconds to Google Protocol Buffers Timestamp
        return util::TimeUtil::NanosecondsToTimestamp(nano_seconds);
    }

} // namespace camcontroller::utility
