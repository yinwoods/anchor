from configparser import ConfigParser


config = ConfigParser()
config.read('./config.ini')

servers = config["test"]["servers"].split(",")
