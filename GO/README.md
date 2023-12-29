# Video Monitoring Application Overview

## Application Overview:

1. **Purpose:**
    - The application is designed for video monitoring and surveillance.

2. **Components:**
    - The application is divided into several packages and components, including `handler`, `comunication`, `middleware`, and `storage`.

3. **Backend Services:**
    - The backend services are implemented using the Go programming language.
    - Key backend services include handling video feeds, managing cameras, and processing video frames.

4. **Communication:**
    - The application uses Kafka for message communication between different components.
    - It utilizes gRPC for communication related to camera subscriptions.

5. **User Management:**
    - There's a user management system with authentication and authorization features.
    - User authentication is done with tokens, and there's a middleware for handling authentication.

6. **Web Server:**
    - The application includes a web server using the Gin framework.
    - It serves HTML pages, such as video feeds and index pages.

7. **Configuration:**
    - Configuration parameters for the application are loaded using the Viper library from a JSON configuration file.

8. **External Services:**
    - The application interacts with external services such as Kafka, Minio (for storage), and MongoDB (for database).

9. **Cleanup Service:**
    - There's a periodic cleanup service responsible for managing the cleanup of camera images.

10. **Integration:**
    - The application integrates various services, such as video feed handling, camera management, user authentication, and storage.

## Key Components:

- **Video Feed Handling:**
    - The `VideoFeedHandler` is responsible for rendering video feeds and managing the index page with user's cameras.

- **Camera Management:**
    - The application has a `CameraHandler` for handling camera-related operations, including creating, logging in, and processing camera commands.

- **Message Communication:**
    - The `comunication` package includes components for interacting with Kafka, handling messages, and managing subscriptions.

- **Frame Storage:**
    - The `FrameHandler` is responsible for storing frame information, processing messages, and interacting with the camera repository.

- **User Authentication and Middleware:**
    - Middleware functions, such as authentication and IP authorization, are used to enhance the security of the application.

## Technologies Used:

- **Programming Language:**
    - The backend is implemented in Go.

- **Frameworks:**
    - The Gin framework is used for web server functionalities.

- **Message Brokers:**
    - Kafka is used for asynchronous message communication.

- **Data Storage:**
    - MongoDB is used for camera-related data storage.

- **Configuration Management:**
    - Viper is used for loading and managing configuration parameters.

## Overall Flow:

1. Users log in to the application, and their authentication is validated using tokens.

2. Cameras are managed, and the application interacts with external services for message communication (Kafka), storage (Minio, MongoDB), and user management.

3. Video feeds are rendered on web pages, and periodic cleanup services manage the removal of old camera images.

4. The application is designed to handle asynchronous communication, enabling real-time updates and notifications related to camera events.

Please note that this is a high-level overview, and the actual behavior and functionality of the application might require more detailed information from the complete codebase.
