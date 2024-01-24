import json
import logging
import os

from mongoengine import connect

from prometheus.query import PrometheusInterface
from sla_manager.sla_server import SlaManger


def load_configuration():
    config_path = os.environ.get('CONFIG_PATH', 'config.json')
    with open(config_path, 'r') as file:
        config = json.load(file)
    return config


def main() -> None:
    # Load configuration from a JSON file
    config = load_configuration()
    # Enable logging
    logging.basicConfig(format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.INFO)
    logging.getLogger("httpx").setLevel(logging.WARNING)
    logging.info(f"Loaded configuration: {config}")

    # Database connection
    logging.info("Connecting to the database...")
    connect(config["mongo"]["dbname"], host=config["mongo"]["host"], port=int(config["mongo"]["port"]),
            username=config["mongo"]["username"], password=config["mongo"]["password"])
    logging.info("Database connection established.")

    logging.info("Starting server thread...")
    SlaManger.prometheus = PrometheusInterface(config["prometheus"])
    SlaManger.app.run(debug=True, host='0.0.0.0', port=13000)


if __name__ == '__main__':
    main()
