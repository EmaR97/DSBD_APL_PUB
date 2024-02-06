import json
import logging
import os

from mongoengine import connect


def enable_log_load_conf():
    # Enable logging
    logging.basicConfig(format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.DEBUG)
    logging.getLogger("httpx").setLevel(logging.DEBUG)

    logging.info("Loading configuration from 'config.json'...")
    config_path = os.environ.get('CONFIG_PATH', 'config.json')
    with open(config_path, 'r') as file:
        config = json.load(file)
    return config


def connect_to_mongo_db(dbname_, host_, port, username_, password_):
    # Database connection
    logging.info("Connecting to the database...")
    connect(dbname_, host=host_, port=port, username=username_, password=password_)
    logging.info("Database connection established.")
