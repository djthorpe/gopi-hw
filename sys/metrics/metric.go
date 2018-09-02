/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package metrics

import (
	// Frameworks
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Metric implements the gopi.Metric interface
type metric struct {
	name        string
	metric_rate gopi.MetricRate
	metric_type gopi.MetricType
	buckets     map[int]*value
	current     int
}

type value struct {
	uint_sum  uint64
	float_sum float64
	samples   uint
}

////////////////////////////////////////////////////////////////////////////////
// NEW METRIC

func NewMetric(name string, metric_rate gopi.MetricRate, metric_type gopi.MetricType) *metric {
	// check for bad parameters
	if name == "" || metric_rate == gopi.METRIC_RATE_NONE || metric_type == gopi.METRIC_TYPE_NONE {
		return nil
	}
	// return new metric
	return &metric{
		name:        name,
		metric_rate: metric_rate,
		metric_type: metric_type,
		buckets:     make(map[int]*value, 60),
		current:     -1,
	}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE IMPLEMENTATION

func (this *metric) Name() string {
	return this.name
}
func (this *metric) Rate() gopi.MetricRate {
	return this.metric_rate
}
func (this *metric) Type() gopi.MetricType {
	return this.metric_type
}

func (this *metric) Unit() string {
	switch this.metric_type {
	case gopi.METRIC_TYPE_CELCIUS:
		return "°C"
	default:
		return ""
	}
}

func (this *metric) UintValue() uint {
	if this.current == -1 {
		// No current bucket, return 0
		return 0
	} else {
		return uint(math.Round(this.buckets[this.current].Mean()))
	}
}

func (this *metric) FloatValue() float64 {
	if this.current == -1 {
		// No current bucket, return NaN
		return math.NaN()
	} else {
		return float64(this.buckets[this.current].Mean())
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *metric) RecordValue(v reflect.Value, ts time.Time) error {
	// Determine the bucket (minute, hour or day)
	if i := this.bucket(ts); i < 0 {
		return gopi.ErrBadParameter
	} else {
		if i != this.current && this.current != -1 {
			// if 'this' bucket exists then delete it
			if _, exists := this.buckets[i]; exists {
				delete(this.buckets, i)
			}
		}
		// If the bucket doesn't exist then create it
		if _, exists := this.buckets[i]; exists == false {
			this.buckets[i] = &value{}
		}
		bucket := this.buckets[i]
		switch v.Kind() {
		case reflect.Uint:
			bucket.uint_sum += v.Uint()
		case reflect.Float64:
			bucket.float_sum += v.Float()
		default:
			return gopi.ErrBadParameter
		}
		// Increment number of samples captured in this bucket
		bucket.samples += 1
		this.current = i
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *metric) bucket(ts time.Time) int {
	switch this.metric_rate {
	case gopi.METRIC_RATE_MINUTE:
		// return second
		return ts.Second()
	case gopi.METRIC_RATE_HOUR:
		// return minute
		return ts.Minute()
	case gopi.METRIC_RATE_DAY:
		// return hour
		return ts.Hour()
	default:
		return -1
	}
}

func (this *value) Mean() float64 {
	if this.samples == 0 {
		return math.NaN()
	} else {
		return (float64(this.uint_sum) + this.float_sum) / float64(this.samples)
	}
}

func (this *value) Sum() float64 {
	if this.samples == 0 {
		return math.NaN()
	} else {
		return (float64(this.uint_sum) + this.float_sum)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *metric) String() string {
	return fmt.Sprintf("<metric>{ name='%v' rate=%v type=%v mean=%v%v }", this.name, this.metric_rate, this.metric_type, this.FloatValue(), this.Unit())
}
