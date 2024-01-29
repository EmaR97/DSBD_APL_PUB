#### RUN THE CAM CONTROLLER LOCALLY FOR IT TO ACCESS THE CAM

Update and install required packages

```bash
apt-get update && \
    apt-get install -y make cmake gdb g++ ffmpeg libcurl4-openssl-dev libopencv-dev \
                       libprotobuf-dev protobuf-compiler libpaho-mqtt-dev git build-essential libssl-dev \
                       librdkafka-dev curl zip unzip tar
```

Install vcpkg

```bash
git clone https://github.com/microsoft/vcpkg.git /vcpkg
cd /vcpkg
./bootstrap-vcpkg.sh
./vcpkg integrate install
```

Install dependencies using vcpkg

```bash
./vcpkg install aws-sdk-cpp
```

Explicitly set CMAKE_TOOLCHAIN_FILE for vcpkg

```bash
export CMAKE_TOOLCHAIN_FILE=/vcpkg/scripts/buildsystems/vcpkg.cmake
```

Clone and build paho.mqtt.cpp library

```bash
git clone https://github.com/eclipse/paho.mqtt.cpp.git /paho.mqtt.cpp
cd /paho.mqtt.cpp
mkdir build
cd build
cmake ..
make -j8
make install
```

Clone and build prometheus client library

```bash
git clone --recursive https://github.com/jupp0r/prometheus-cpp.git /prometheus/prometheus-cpp
mkdir /prometheus/prometheus-cpp/build
cd /prometheus/prometheus-cpp/build
cmake ..
make -j8
make install
```

Move to the project root dir

```bash
cd /project/root #to replace with the correct path
```

Compile Protobuf files

```bash
protoc -I=./proto --cpp_out=./src/message framedata.proto
protoc -I=./proto --cpp_out=./src/message frameinfo.proto
```

Move to the CPP build folder

```bash
cd CPP/build
```

Run CMake to configure the build

```bash
cmake ..
```

Compile the C++ code

```bash
make
```

Run the CamController

```bash
CamController
```

#### DEPLOY THE REST OF THE SISTEM ON KIND

```bash
kind create cluster --config K8s/cluster.yaml
kubectl apply -k K8s
```

#### DEPLOY THE REST OF THE SISTEM ON DOCKER

```bash
docker compose up
```