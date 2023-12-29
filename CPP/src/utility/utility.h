#ifndef CAMCONTROLLER_UTILITY_H
#define CAMCONTROLLER_UTILITY_H

#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/util/time_util.h>

namespace my_namespace::utility {
    google::protobuf::Timestamp timestamp_ns();
} // namespace camcontroller::utility

#endif // CAMCONTROLLER_UTILITY_H
