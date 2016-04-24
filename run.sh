#!/bin/sh -e
# This is a simple shell script to start the program. The main reason is because
# it allows to pass along some parameters in environment variables.
exec /go/bin/app -mqtt $MQTT_HOST -influxdb $INFLUXDB_HOST -database $DATABASE_NAME
