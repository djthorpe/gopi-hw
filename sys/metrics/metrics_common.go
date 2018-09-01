/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package metrics

import (
	"fmt"
	"reflect"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Metrics struct{}

type metrics struct {
	log     gopi.Logger
	metrics []*gopi.Metric
	cases   []reflect.SelectCase
	changed chan struct{}
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// Timestamp for module creation
	ts = time.Now()
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open creates a new metrics object, returns error if not possible
func (config Metrics) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<metrics.linux>Open{}")

	// create new driver
	this := new(metrics)
	this.log = log
	this.metrics = make([]*gopi.Metric, 0)

	// The array of channels on which we can
	// accept metrics
	this.cases = make([]reflect.SelectCase, 1)
	this.changed = make(chan struct{})
	this.cases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(this.changed),
	}

	// return driver
	return this, nil
}

// Close connection
func (this *metrics) Close() error {
	this.log.Debug("<metrics.linux>Close{}")

	// Release resources
	this.metrics = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RETURN UPTIME

func (this *metrics) UptimeApp() time.Duration {
	return time.Since(ts)
}

////////////////////////////////////////////////////////////////////////////////
// METRIC METHODS

// Return metric channel, which when you send a value on it will store the metric
func (this *metrics) NewMetricUint(metric_type gopi.MetricType, metric_rate gopi.MetricRate, name string) (chan<- uint, error) {
	// Check incoming parameters
	if metric_type == gopi.METRIC_TYPE_NONE || name == "" {
		return nil, gopi.ErrBadParameter
	}

	// Create a new metric structure and append
	metric := &gopi.Metric{
		Rate: metric_rate,
		Type: metric_type,
		Name: name,
	}
	this.metrics = append(this.metrics, metric)

	// Return channel for metrics
	return nil, nil
}

// Metrics returns all metrics of a particular type, or METRIC_TYPE_NONE for all metrics
func (this *metrics) Metrics(metric_type gopi.MetricType) []*gopi.Metric {
	metrics := make([]*gopi.Metric, 0, len(this.metrics))
	for _, metric := range this.metrics {
		if metric_type == gopi.METRIC_TYPE_NONE || metric_type == metric.Type {
			metrics = append(metrics, metric)
		}
	}
	return metrics
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *metrics) String() string {
	var l [3]float64
	l[0], l[1], l[2] = this.LoadAverage()
	return fmt.Sprintf("<metrics>{ uptime_host=%v uptime_app=%v load_average=%v metrics=%v }", this.UptimeHost(), this.UptimeHost(), l, this.metrics)
}
