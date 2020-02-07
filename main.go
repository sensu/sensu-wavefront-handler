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
	Host                 string
	MetricsPort          int
	FlushIntervalSeconds int
	Prefix               string
	Tags                 map[string]string
}

const (
	host   = "host"
	port   = "metrics-port"
	flush  = "flush-interval-seconds"
	prefix = "prefix"
	tags   = "tags"
)

var (
	handlerConfig = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-wavefront-handler",
			Short:    "sends metrics to a wavefront proxy using the wavefront data format",
			Keyspace: "sensu.io/plugins/sensu-wavefront-handler/config",
		},
	}

	opts = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      host,
			Env:       "WAVEFRONT_HOST",
			Argument:  host,
			Shorthand: "",
			Default:   "127.0.0.1",
			Usage:     "the host of the wavefront proxy",
			Value:     &handlerConfig.Host,
		},
		&sensu.PluginConfigOption{
			Path:      port,
			Env:       "WAVEFRONT_METRICS_PORT",
			Argument:  port,
			Shorthand: "m",
			Default:   2878,
			Usage:     "the port of the wavefront proxy",
			Value:     &handlerConfig.MetricsPort,
		},
		&sensu.PluginConfigOption{
			Path:      flush,
			Env:       "WAVEFRONT_FLUSH_INTERVAL_SECONDS",
			Argument:  flush,
			Shorthand: "f",
			Default:   1,
			Usage:     "the flush interval of the wavefront proxy (in seconds)",
			Value:     &handlerConfig.FlushIntervalSeconds,
		},
		&sensu.PluginConfigOption{
			Path:      prefix,
			Env:       "WAVEFRONT_PREFIX",
			Argument:  prefix,
			Shorthand: "p",
			Default:   "",
			Usage:     "the string to prepended to the metric name",
			Value:     &handlerConfig.Prefix,
		},
		&sensu.PluginConfigOption{
			Path:      tags,
			Env:       "WAVEFRONT_TAGS",
			Argument:  tags,
			Shorthand: "t",
			Default:   nil,
			Usage:     "the additional tags to merge with the metric tags",
			Value:     &handlerConfig.Tags,
		},
	}
)

func main() {
	handler := sensu.NewGoHandler(&handlerConfig.PluginConfig, opts, checkArgs, executeHandler)
	handler.Execute()
}

func checkArgs(event *corev2.Event) error {
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

	proxyCfg := &wavefront.ProxyConfiguration{
		Host:                 handlerConfig.Host,
		MetricsPort:          handlerConfig.MetricsPort,
		FlushIntervalSeconds: handlerConfig.FlushIntervalSeconds,
	}

	sender, err := wavefront.NewProxySender(proxyCfg)
	if err != nil {
		return err
	}

	for _, point := range event.Metrics.Points {
		tags := make(map[string]string)
		for _, tag := range point.Tags {
			tags[tag.Name] = tag.Value
		}

		// merge tags if provided as config option
		if handlerConfig.Tags != nil {
			for k, v := range handlerConfig.Tags {
				tags[k] = v
			}
		}

		// prefix metric name if provided as config option
		name := point.Name
		if handlerConfig.Prefix != "" {
			name = fmt.Sprintf("%s.%s", handlerConfig.Prefix, name)
		}

		sender.SendMetric(name, point.Value, point.Timestamp, event.Entity.Name, tags)
	}

	log.Printf("sent %d metric points with %d failures", len(event.Metrics.Points), sender.GetFailureCount())
	sender.Flush()
	sender.Close()
	return nil
}
