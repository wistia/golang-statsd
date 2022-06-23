package statsd

import (
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

const standardPort = "8125"

var client *statsd.Client
var log Logger

func init() {
	Disable()
}

func Configure(host, namespace, env, component string, logger Logger) error {
	log = logger
	var tags []string

	if env != "" {
		tags = append(tags, "env:"+env)
	}
	if component != "" {
		tags = append(tags, "component:"+component)
	}

	c, err := statsd.New(host+":"+standardPort,
		statsd.WithNamespace(namespace),
		statsd.WithTags(tags))
	if err != nil {
		return err
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

// IncrByOne convenience function that increments the tags by exactly 1
func IncrByOne(name string, tags ...string) {
	IncrByOne(name, tags...)
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

func Flush() {
	if err := client.Flush(); err != nil {
		logMetricError(err)
	}
}

func logMetricError(err error) {
	log.Printf("error reporting metric: %v", err)
}

type Logger interface {
	Printf(format string, args ...interface{})
}

type SilentStatsdWriter struct{}

func (ssw SilentStatsdWriter) Close() error                         { return nil }
func (ssw SilentStatsdWriter) Write(data []byte) (n int, err error) { return 0, nil }
func (ssw SilentStatsdWriter) SetWriteTimeout(time.Duration) error  { return nil }
