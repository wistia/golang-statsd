package statsd

import (
	"fmt"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

const standardPort = "8125"

var client *statsd.Client
var log Logger

func init() {
	_ = Disable()
}

// Configure sets up the internal client and logger and sets them up for use
func Configure(host, namespace, env, component string, logger Logger) error {
	log = logger
	var tags []string

	if env != "" {
		tags = append(tags, "env:"+env)
	}
	if component != "" {
		tags = append(tags, "component:"+component)
	}

	if IsIPv6(host) {
		host = fmt.Sprintf("[%s]", host)
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

// IsIPv6 checks whether the given address conforms to IPv6 notation
func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

// Disable stops metrics from being emitted
func Disable() error {
	c, err := statsd.NewWithWriter(SilentStatsdWriter{})
	if err != nil {
		return err
	}
	client = c
	return nil
}

// Count increments the given metric by the given value
func Count(name string, value int64, tags []string, rate float64) {
	if err := client.Count(name, value, tags, rate); err != nil {
		logMetricError(err)
	}
}

// Decr is just Count of -1
func Decr(name string, tags []string, rate float64) {
	if err := client.Decr(name, tags, rate); err != nil {
		logMetricError(err)
	}
}

// Distribution tracks the statistical distribution of a set of values across
// your infrastructure.
func Distribution(name string, value float64, tags []string, rate float64) {
	if err := client.Distribution(name, value, tags, rate); err != nil {
		logMetricError(err)
	}
}

// Gauge measures the value of a metric at a particular time.
func Gauge(name string, value float64, tags []string, rate float64) {
	if err := client.Gauge(name, value, tags, rate); err != nil {
		logMetricError(err)
	}
}

// Histogram tracks the statistical distribution of a set of values on each host.
func Histogram(name string, value float64, tags []string, rate float64) {
	if err := client.Histogram(name, value, tags, rate); err != nil {
		logMetricError(err)
	}
}

// Incr is just Count of 1
func Incr(name string, tags []string, rate float64) {
	if err := client.Incr(name, tags, rate); err != nil {
		logMetricError(err)
	}
}

// SimpleIncr is a convenience function that increments the metric by exactly 1
func SimpleIncr(name string, tags ...string) {
	Incr(name, tags, 1.0)
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

// Flush forces a flush of all the queued dogstatsd payloads. This method is
// blocking and will not return until everything is sent through the network.
func Flush() {
	if err := client.Flush(); err != nil {
		logMetricError(err)
	}
}

func logMetricError(err error) {
	log.Printf("error reporting metric: %v", err)
}

// Logger interface that represents a client logger
type Logger interface {
	Printf(format string, args ...interface{})
}

// SilentStatsdWriter creates a writer that does nothing, dropping all data
// coming in.
type SilentStatsdWriter struct{}

// Close implements the io.Closer interface
func (ssw SilentStatsdWriter) Close() error { return nil }

// Write implements the io.Writer interface
func (ssw SilentStatsdWriter) Write(_ []byte) (n int, err error) { return 0, nil }

// SetWriteTimeout pretends to set a timeout for the writer
func (ssw SilentStatsdWriter) SetWriteTimeout(time.Duration) error { return nil }
