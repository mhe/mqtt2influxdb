# mqtt2influxdb #

mqtt2influxdb is a small and simple program that subscribes to mqtt topics and publishes messages received to an InfluxDB database. 

## Building & installing ##

First make sure you have setup your [Go](https://golang.org) environment. After having setup your $GOPATH you can do a

    go get github.com/mhe/mqtt2influxdb

to install a mqtt2influxdb binary in your $GOPATH/bin.

Alternatively you can clone this repository somewhere (for example in your $GOPATH/src) and install the dependencies

    cd mqtt2influxdb 
    go get ...

Then build it:

    go build .

Which should result in a binary in your working directory.

## Running ##

mqtt2influxdb has a number of commandline options. Invoking `mqtt2influxdb -h` will print:

```
mqtt2influxdb
Usage of ./mqtt2influxdb:
  -clientid string
    	ClientID to use when connecting to mqtt broker. (default "mqtt2influxdb")
  -config string
    	Configuration file with mappings. (default "mqtt2influxdb.toml")
  -database string
    	Name of the InfluxDB database to use. (default "mqtt")
  -influxdb string
    	InfluxDB host address. Should include both protocol (http or https) and port number. (default "http://localhost:8086")
  -mqtt string
    	Mqtt host (including port number). (default "localhost:1883")
  -test
    	Print InfluxDB insert lines to stdout instead of actually submitting data.
```

## Configuration ##
How to map messages from an mqtt bus to points in a influxdb database is specified in a configuration file (using the [TOML](https://github.com/toml-lang/toml) format).

Each mapping consists of three items:

- `topic`: the mqtt topic to subscribe to. It can contain wildcards (e.g., `+`). 
- `template`: a template (using [Go's text/template](https://golang.org/pkg/text/template/)) to build a line according to [InfluxDB's line protocol](https://docs.influxdata.com/influxdb/v0.12/write_protocols/write_syntax/). 
- `encoding`: specifies the encoding of the message, can be one of four: [json](http://json.org/), [msgpack](https://github.com/msgpack/msgpack), [binc](http://github.com/ugorji/binc), or [cbor](http://cbor.io/]. It can be omitted if a default encoding is specified (the `defaultEncoding` top-level entry).

See the provided example configuration file for more information.

## Running in Docker ##

A basic Dockerfile is provided to create a docker image to run mqtt2influxdb in. Note that it is relatively unoptimized in terms of size. There are a couple of configuration parameters that can be specified using environment variables: 

- `MQTT_HOST`
- `INFLUXDB_HOST`
- `DATABASE_NAME`

