# Define the Kafka topics creation commands
commands=(
    "kafka-topics --create --if-not-exists --bootstrap-server kafka:9092 --replication-factor 1 --partitions 100 --topic frame_data"
    "kafka-topics --create --if-not-exists --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1 --topic frame_info"
    "kafka-topics --create --if-not-exists --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1 --topic notification"
)

# Loop through the commands until a positive response is received
for cmd in "${commands[@]}"; do
    echo "Executing command: $cmd"
    # Run the command and check the exit code
    until $cmd; do
        echo "Command failed, retrying in 5 seconds..."
        sleep 5
    done
    echo "Command executed successfully."
done

echo "All Kafka topics created successfully."