package statsd

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
	"time"
)

const standardPort = "8125"

var client *statsd.Client

func init() {
	Disable()
}

func Configure(host, namespace, env, component string) error {
	c, err := statsd.New(host + ":" + standardPort)
	if err != nil {
		return err
	}

	c.Namespace = fmt.Sprintf("%s.", namespace)

	if env != "" {
		c.Tags = append(c.Tags, "env:"+env)
	}

	if component != "" {
		c.Tags = append(c.Tags, "component:"+component)
	}

	client = c
	return nil
}

func Disable() error {
	c, err := statsd.NewWithWriter(SilentStatsdWriter{})
	if err != nil {
		return err
	}
	client = c
	return nil
}

func Count(name string, value int64, tags []string, rate float64) {
	if err := client.Count(name, value, tags, rate); err != nil {
		logMetricError(err)
	}
}

func Decr(name string, tags []string, rate float64) {
	if err := client.Decr(name, tags, rate); err != nil {
		logMetricError(err)
	}
}

func Distribution(name string, value float64, tags []string, rate float64) {
	if err := client.Distribution(name, value, tags, rate); err != nil {
		logMetricError(err)
	}
}

func Gauge(name string, value float64, tags []string, rate float64) {
	if err := client.Gauge(name, value, tags, rate); err != nil {
		logMetricError(err)
	}
}

func Histogram(name string, value float64, tags []string, rate float64) {
	if err := client.Histogram(name, value, tags, rate); err != nil {
		logMetricError(err)
	}
}

func Incr(name string, tags []string, rate float64) {
	if err := client.Incr(name, tags, rate); err != nil {
		logMetricError(err)
	}
}

// MeasureExecTime returns the execution time in milliseconds since the provided start time. It's
// intended to be used with a defer block.
func MeasureExecTime(name string, tags []string, rate float64, start time.Time) time.Duration {
	elapsed := time.Since(start)
	Histogram(name, float64(elapsed/1e6), tags, rate)
	return elapsed
}

// SimpleMeasureExecTime is a convenience wrapper around MeasureExecTime and returns the execution
// time in milliseconds since the provided start time. It's intended to be used with a defer block.
func SimpleMeasureExecTime(name string, start time.Time) time.Duration {
	return MeasureExecTime(name, []string{}, 1.0, start)
}

func logMetricError(err error) {
	log.Printf("error reporting metric: %v", err)
}

type SilentStatsdWriter struct{}

func (ssw SilentStatsdWriter) Close() error                         { return nil }
func (ssw SilentStatsdWriter) Write(data []byte) (n int, err error) { return 0, nil }
func (ssw SilentStatsdWriter) SetWriteTimeout(time.Duration) error  { return nil }
