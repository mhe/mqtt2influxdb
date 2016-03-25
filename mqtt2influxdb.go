// mqtt2influxdb is a small program that subscribes to mqtt channels and
// publishes the data to an InfluxDB timeseries database.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ugorji/go/codec"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"
)

var configfile = flag.String("config", "mqtt2influxdb.toml", "Configuration file")

type Config struct {
	DefaultEncoding string

	Mqtt struct {
		Host     string
		ClientID string
	}
	InfluxDB struct {
		Host     string
		Database string
	}
	Mappings []*struct {
		Topic    string
		Template string
		Encoding string
	}
}

func getConfig(filename string) Config {
	var conf Config
	if _, err := toml.DecodeFile(filename, &conf); err != nil {
		log.Fatal(err)
	}

	// Set some defaults
	for _, mapping := range conf.Mappings {
		if mapping.Encoding == "" {
			mapping.Encoding = conf.DefaultEncoding
		}
	}

	return conf
}

func main() {
	fmt.Println("mqtt2influxdb")

	flag.Parse()

	conf := getConfig(*configfile)

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  conf.Mqtt.Host,
		ClientID: []byte(conf.Mqtt.ClientID),
	})
	if err != nil {
		panic(err)
	}

	// Make sure the new msgpack spec is used
	customMsgPackHandle := new(codec.MsgpackHandle)
	customMsgPackHandle.RawToString = true

	encodingMap := map[string]codec.Handle{
		"json":    new(codec.JsonHandle),
		"msgpack": customMsgPackHandle,
		"binc":    new(codec.BincHandle),
		"cbor":    new(codec.CborHandle),
	}

	influxDBWriteURL := conf.InfluxDB.Host + "/write?db=" + conf.InfluxDB.Database

	// Create topic subscriptions
	subscriptions := make([]*client.SubReq, len(conf.Mappings))
	for i, mapping := range conf.Mappings {
		// Setup the template
		topicTemplate := template.Must(template.New(mapping.Topic).Parse(mapping.Template))

		// Create a buffer to send the output of the template to the http post body
		buffer := new(bytes.Buffer)
		h := encodingMap[mapping.Encoding]

		subscriptions[i] = &client.SubReq{
			TopicFilter: []byte(mapping.Topic),
			QoS:         mqtt.QoS1,
			Handler: func(topicName, message []byte) {
				// Unmarshal the data into a interface{} object. Probably not the
				// fastest approach, but works for now.
				var f interface{}
				dec := codec.NewDecoderBytes(message, h)
				err := dec.Decode(&f)

				// Execute the template
				err = topicTemplate.Execute(buffer, f)
				if err != nil {
					log.Fatal(err)
				}

				// And send the result off
				resp, err := http.Post(influxDBWriteURL, "text/plain", buffer)
				if err != nil {
					log.Println("Error submitting data to InfluxDB data base: ", err)
				}
				// Cleanup the response
				resp.Body.Close()
			},
		}
	}

	// Actually subscribe to topics.
	err = cli.Subscribe(&client.SubscribeOptions{SubReqs: subscriptions})
	if err != nil {
		panic(err)
	}

	// Now simply wait until we get a signal to stop
	<-sigc

	// Disconnect the Network Connection.
	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}
