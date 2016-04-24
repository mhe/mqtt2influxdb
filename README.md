# mqtt2influxdb #

mqtt2influxdb is a small and simple program that publishes messages received on certain mqtt topics to an InfluxDB database. More documentation will come.

## Running in Docker ##

mqtt2influxdb is specifically designed to run in a container. A Dockerfile is provided. A number of things can be specified using environment variables:

- MQTT_HOST
- INFLUXDB_HOST
- DATABASE_NAME

