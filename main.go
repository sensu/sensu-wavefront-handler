package main

import (
	"fmt"
	"log"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	wavefront "github.com/wavefronthq/wavefront-sdk-go/senders"
)

// Config represents a handler config for the wavefront handler.
type Config struct {
	sensu.PluginConfig
	Server string
	Token  string
}

// Options represents the config options for the wavefront handler.
type Options struct {
	Example sensu.PluginConfigOption
}

const (
	server = "server"
	token  = "token"
)

var (
	handlerConfig = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-wavefront-handler",
			Short:    "a wavefront metrics handler built for use with sensu",
			Timeout:  10,
			Keyspace: "sensu.io/plugins/sensu-wavefront-handler/config",
		},
	}

	opts = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      server,
			Env:       "WAVEFRONT_SERVER",
			Argument:  server,
			Shorthand: "s",
			Default:   "",
			Usage:     "the address of the wavefront server",
			Value:     &handlerConfig.Server,
		},
		&sensu.PluginConfigOption{
			Path:      token,
			Env:       "WAVEFRONT_TOKEN",
			Argument:  token,
			Shorthand: "t",
			Default:   "",
			Usage:     "the API token for the wavefront server",
			Value:     &handlerConfig.Token,
		},
	}
)

func main() {
	handler := sensu.NewGoHandler(&handlerConfig.PluginConfig, opts, checkArgs, executeHandler)
	handler.Execute()
}

func checkArgs(event *corev2.Event) error {
	if handlerConfig.Server == "" {
		return fmt.Errorf("--server or WAVEFRONT_SERVER environment variable is required")
	}
	if handlerConfig.Token == "" {
		return fmt.Errorf("--token or WAVEFRONT_TOKEN environment variable is required")
	}
	if !event.HasMetrics() {
		return fmt.Errorf("event does not contain metrics")
	}
	return nil
}

func executeHandler(event *corev2.Event) error {
	if len(event.Metrics.Points) == 0 {
		log.Println("event does not contain metric points")
		return nil
	}

	directCfg := &wavefront.DirectConfiguration{
		Server:               handlerConfig.Server,
		Token:                handlerConfig.Token,
		BatchSize:            10000,
		MaxBufferSize:        50000,
		FlushIntervalSeconds: 1,
	}

	sender, err := wavefront.NewDirectSender(directCfg)
	if err != nil {
		return err
	}

	for _, point := range event.Metrics.Points {
		tags := make(map[string]string)
		for _, tag := range point.Tags {
			tags[tag.Name] = tag.Value
		}
		sender.SendMetric(point.Name, point.Value, point.Timestamp, event.Entity.Name, tags)
	}

	log.Printf("sent %d metric points with %d failures", len(event.Metrics.Points), sender.GetFailureCount())
	sender.Flush()
	sender.Close()
	return nil
}
