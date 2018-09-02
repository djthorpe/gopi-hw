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
	metrics []*metric
	cases   []reflect.SelectCase
	changed chan struct{}
	done    chan struct{}
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
	log.Debug("<metrics>Open{}")

	// create new driver
	this := new(metrics)
	this.log = log

	// The array of channels on which we can
	// accept metrics
	this.metrics = make([]*metric, 0)
	this.cases = make([]reflect.SelectCase, 1)
	this.changed = make(chan struct{})
	this.done = make(chan struct{})
	this.cases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(this.changed),
	}

	// Start goroutine to accept incoming metrics
	go this.goGatherMetrics()

	// return driver
	return this, nil
}

// Close connection
func (this *metrics) Close() error {
	this.log.Debug("<metrics>Close{}")

	// Close changed channel - which ends goGatherMetrics
	if this.done != nil {
		for _, c := range this.cases {
			c.Chan.Close()
		}
		<-this.done
	}

	// Release resources
	this.metrics = nil
	this.cases = nil
	this.changed = nil
	this.done = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RETURN UPTIME FOR THE APPLICATION

func (this *metrics) UptimeApp() time.Duration {
	return time.Since(ts)
}

////////////////////////////////////////////////////////////////////////////////
// METRIC METHODS

// NewMetricUint returns metric channel, which when you send a value on it will store the metric
func (this *metrics) NewMetricUint(metric_type gopi.MetricType, metric_rate gopi.MetricRate, name string) (chan<- uint, error) {
	// Create a new metric structure and append to list of metrics
	if m := NewMetric(name, metric_rate, metric_type); m == nil {
		return nil, gopi.ErrBadParameter
	} else {
		this.metrics = append(this.metrics, m)
	}

	// Create channel for metrics
	c := make(chan uint)
	this.cases = append(this.cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	})
	this.changed <- gopi.DONE

	// Return channel for metrics
	return c, nil
}

// NewMetricFloat64 returns metric channel, which when you send a value on it will store the metric
func (this *metrics) NewMetricFloat64(metric_type gopi.MetricType, metric_rate gopi.MetricRate, name string) (chan<- float64, error) {
	// Create a new metric structure and append to list of metrics
	if m := NewMetric(name, metric_rate, metric_type); m == nil {
		return nil, gopi.ErrBadParameter
	} else {
		this.metrics = append(this.metrics, m)
	}

	// Create channel for metrics
	c := make(chan float64)
	this.cases = append(this.cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	})
	this.changed <- gopi.DONE

	// Return channel for metrics
	return c, nil
}

// Metrics returns all metrics of a particular type, or METRIC_TYPE_NONE for all metrics
func (this *metrics) Metrics(metric_type gopi.MetricType) []gopi.Metric {
	metrics := make([]gopi.Metric, 0, len(this.metrics))
	for _, metric := range this.metrics {
		if metric_type == gopi.METRIC_TYPE_NONE || metric.Type() == metric_type {
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

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// goGatherMetrics accepts incoming metric values
func (this *metrics) goGatherMetrics() {
	this.log.Debug2("<metrics>goGatherMetrics started")
FOR_LOOP:
	for {
		// select cases
		i, v, ok := reflect.Select(this.cases)
		if i == 0 && ok == true {
			// We need to reload the cases
		} else if i == 0 && ok == false {
			// Closed changed channel, so end
			break FOR_LOOP
		} else if ok == true {
			if err := this.recordMetric(i, v); err != nil {
				this.log.Warn("<metrics>goGatherMetrics: %v", err)
			}
		}
	}
	this.log.Debug2("<metrics>goGatherMetrics ended")
	close(this.done)
}

// recordMetric records a metric and returns an error if the
// value could not be recorded
func (this *metrics) recordMetric(i int, v reflect.Value) error {
	if i < 1 || i > len(this.metrics) {
		return gopi.ErrBadParameter
	}
	if err := this.metrics[i-1].RecordValue(v, time.Now()); err != nil {
		return err
	} else {
		this.log.Debug("RecordValue: %v: value=%v mean=%v", this.metrics[i-1].Name(), v, this.metrics[i-1].FloatValue())
		return nil
	}
}
