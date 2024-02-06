package service

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	URL      string `mapstructure:"url"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// MQTTConfig holds MQTT-related configuration
type MQTTConfig struct {
	BrokerURL string `mapstructure:"brokerURL"`
	ClientID  string `mapstructure:"clientID"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
}

// MQTTManagementConfig holds MQTT-related configuration
type MQTTManagementConfig struct {
	Url      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	HttpPort        string        `mapstructure:"httpPort"`
	GrpcPort        string        `mapstructure:"grpcPort"`
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
	AllowedIPs      []string      `mapstructure:"allowedIPs"`
	TemplateDir     string        `mapstructure:"templateDir"`
	AuthEndpoint    string        `mapstructure:"authEndpoint"`
}

// MinioConfig holds server-related configuration
type MinioConfig struct {
	Endpoint         string  `mapstructure:"endpoint"`
	Username         string  `mapstructure:"username"`
	Password         string  `mapstructure:"password"`
	BucketName       string  `mapstructure:"bucketName"`
	CleanUpInterval  float64 `mapstructure:"cleanUpInterval"`
	CleanUpOlderThan float64 `mapstructure:"cleanUpOlderThan"`
}

// KafkaConfig holds server-related configuration
type KafkaConfig struct {
	Brokers                   string `mapstructure:"brokers"`
	GroupID                   string `mapstructure:"groupID"`
	TopicGetCamIds            string `mapstructure:"topic_in_cam_ids"`
	TopicProcessAndStoreFrame string `mapstructure:"topic_in_frame_info"`
	TopicOut                  string `mapstructure:"topic_out"`
}

// Config struct holds all configuration parameters
type Config struct {
	Database DatabaseConfig       `mapstructure:"database"`
	MQTT     MQTTConfig           `mapstructure:"mqtt"`
	Broker   MQTTManagementConfig `mapstructure:"mqtt-management"`
	Server   ServerConfig         `mapstructure:"server"`
	Minio    MinioConfig          `mapstructure:"minio"`
	Kafka    KafkaConfig          `mapstructure:"kafka"`
}

// LoadConfig loads the configuration from the specified path or environment variable.
func LoadConfig() Config {
	configPath := os.Getenv("APP_CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json" // Default configuration file path
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Failed to read configuration from %s: %v\n", configPath, err)
		// Handle the error gracefully, e.g., return a default configuration
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Failed to unmarshal configuration: %v\n", err)
		// Handle the error gracefully, e.g., return a default configuration
	}

	return config
}
