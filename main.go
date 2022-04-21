package main

import (
	"fmt"
	"log"
	"math"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
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

	opts = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      host,
			Env:       "WAVEFRONT_HOST",
			Argument:  host,
			Shorthand: "",
			Default:   "127.0.0.1",
			Usage:     "the host of the wavefront proxy",
			Value:     &handlerConfig.Host,
		},
		&sensu.PluginConfigOption[int]{
			Path:      port,
			Env:       "WAVEFRONT_METRICS_PORT",
			Argument:  port,
			Shorthand: "m",
			Default:   2878,
			Usage:     "the port of the wavefront proxy",
			Value:     &handlerConfig.MetricsPort,
		},
		&sensu.PluginConfigOption[int]{
			Path:      flush,
			Env:       "WAVEFRONT_FLUSH_INTERVAL_SECONDS",
			Argument:  flush,
			Shorthand: "f",
			Default:   1,
			Usage:     "the flush interval of the wavefront proxy (in seconds)",
			Value:     &handlerConfig.FlushIntervalSeconds,
		},
		&sensu.PluginConfigOption[string]{
			Path:      prefix,
			Env:       "WAVEFRONT_PREFIX",
			Argument:  prefix,
			Shorthand: "p",
			Usage:     "the string to be prepended to the metric name",
			Value:     &handlerConfig.Prefix,
		},
		&sensu.MapPluginConfigOption[string]{
			Path:      tags,
			Env:       "WAVEFRONT_TAGS",
			Argument:  tags,
			Shorthand: "t",
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

	defer sender.Close()

	for _, point := range event.Metrics.Points {
		tags := make(map[string]string)
		// merge tags if provided as config option
		if handlerConfig.Tags != nil {
			for k, v := range handlerConfig.Tags {
				tags[k] = v
			}
		}
		// overwrite tags with those from the original event
		for _, tag := range point.Tags {
			tags[tag.Name] = tag.Value
		}

		// prefix metric name if provided as config option
		name := point.Name
		if handlerConfig.Prefix != "" {
			name = fmt.Sprintf("%s.%s", handlerConfig.Prefix, name)
		}
		err := sender.SendMetric(name, point.Value, secTimestamp(point.Timestamp), event.Entity.Name, tags)
		if err != nil {
			log.Printf("error sending metric: %s", err)
		}
	}

	err = sender.Flush()
	log.Printf("sent %d metric points with %d failures", len(event.Metrics.Points), sender.GetFailureCount())
	return err
}

// msTimestamp auto-detection of metric point timestamp precision using a heuristic with a 250-ish year cutoff
func secTimestamp(ts int64) int64 {
	timestamp := ts
	switch ts := math.Log10(float64(timestamp)); {
	case ts < 10:
		// assume timestamp is seconds
	case ts < 13:
		// assume timestamp is milliseconds
		timestamp = (timestamp / 1e3)
	case ts < 16:
		// assume timestamp is microseconds
		timestamp = (timestamp / 1e6)
	default:
		// assume timestamp is nanoseconds
		timestamp = (timestamp / 1e9)
	}

	return timestamp
}
