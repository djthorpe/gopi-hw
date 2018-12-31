// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package mmal

import (
	"fmt"

	// Frameworks

	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *connection) Input() hw.MMALPort {
	return this.input
}

func (this *connection) Output() hw.MMALPort {
	return this.output
}

func (this *connection) Close() error {
	this.log.Debug2("<sys.hw.mmal.connection>Close{ input='%v' output='%v' }", this.input, this.output)

	if err := rpi.MMALPortConnectionDisable(this.handle); err != nil {
		return err
	} else if err := rpi.MMALPortConnectionDestroy(this.handle); err != nil {
		return err
	} else {
		this.handle = nil
		this.input = nil
		this.output = nil
		return nil
	}
}

func (this *connection) Acquire() error {
	this.log.Debug2("<sys.hw.mmal.connection>Acquire{ input='%v' output='%v' }", this.input, this.output)
	return rpi.MMALPortConnectionAcquire(this.handle)
}

func (this *connection) Release() error {
	this.log.Debug2("<sys.hw.mmal.connection>Release{ input='%v' output='%v' }", this.input, this.output)
	return rpi.MMALPortConnectionRelease(this.handle)
}

func (this *connection) Enabled() bool {
	return rpi.MMALPortConnectionEnabled(this.handle)
}

func (this *connection) SetEnabled(value bool) error {
	this.log.Debug2("<sys.hw.mmal.connection>SetEnabled{ input='%v' output='%v' enable=%v }", this.input, this.output, value)
	if value {
		return rpi.MMALPortConnectionEnable(this.handle)
	} else {
		return rpi.MMALPortConnectionDisable(this.handle)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *connection) String() string {
	return fmt.Sprintf("<sys.hw.mmal.connection>{ input='%v' output='%v' }", this.input, this.output)
}
