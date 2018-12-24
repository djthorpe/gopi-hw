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
	"github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
	rpi "github.com/djthorpe/gopi-hw/rpi"
	"github.com/djthorpe/gopi/util/errors"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (this *component) Close() error {
	this.log.Debug("<sys.hw.mmal.component>Close{}")

	err := new(errors.CompoundError)

	if this.handle == nil {
		// Already closed
		return gopi.ErrOutOfOrder
	}

	// Disable input and output ports, destroy pools
	for _, port := range this.input {
		if rpi.MMALPortIsEnabled(port.handle) {
			if err_ := rpi.MMALPortDisable(port.handle); err_ != nil {
				err.Add(err_)
			}
		}
		if err_ := rpi.MMALPortPoolDestroy(port.handle, port.pool); err_ != nil {
			err.Add(err_)
		}
	}
	for _, port := range this.output {
		if rpi.MMALPortIsEnabled(port.handle) {
			if err_ := rpi.MMALPortDisable(port.handle); err_ != nil {
				err.Add(err_)
			}
		}
		if err_ := rpi.MMALPortPoolDestroy(port.handle, port.pool); err_ != nil {
			err.Add(err_)
		}
	}

	if err_ := rpi.MMALComponentDestroy(this.handle); err_ != nil {
		err.Add(err_)
	}

	this.handle = nil

	return err.ErrorOrSelf()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *component) String() string {
	return fmt.Sprintf("<sys.hw.mmal.component>{ name='%v' id=%08X enabled=%v control=%v input=%v output=%v clock=%v }", this.Name(), this.Id(), this.Enabled(), this.control, this.input, this.output, this.clock)
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *component) Name() string {
	if this.handle == nil {
		return ""
	} else {
		return rpi.MMALComponentName(this.handle)
	}
}

func (this *component) Id() uint32 {
	if this.handle == nil {
		return 0xFFFFFFFF
	} else {
		return rpi.MMALComponentId(this.handle)
	}
}

func (this *component) Enabled() bool {
	if this.handle == nil {
		return false
	} else {
		return rpi.MMALComponentIsEnabled(this.handle)
	}
}

func (this *component) SetEnabled(value bool) error {
	if this.handle == nil {
		// Component is not open
		return gopi.ErrOutOfOrder
	} else if value {
		return rpi.MMALComponentEnable(this.handle)
	} else {
		return rpi.MMALComponentDisable(this.handle)
	}
}

func (this *component) Acquire() error {
	if this.handle == nil {
		// Component is not open
		return gopi.ErrOutOfOrder
	} else {
		return rpi.MMALComponentAcquire(this.handle)
	}
}

func (this *component) Release() error {
	if this.handle == nil {
		// Component is not open
		return gopi.ErrOutOfOrder
	} else {
		return rpi.MMALComponentRelease(this.handle)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PORTS

func (this *component) Control() hw.MMALPort {
	return this.control
}

func (this *component) Input() []hw.MMALPort {
	ports := make([]hw.MMALPort, len(this.input))
	for i, port := range this.input {
		ports[i] = port
	}
	return ports
}

func (this *component) Output() []hw.MMALPort {
	ports := make([]hw.MMALPort, len(this.output))
	for i, port := range this.output {
		ports[i] = port
	}
	return ports
}

func (this *component) Clock() []hw.MMALPort {
	ports := make([]hw.MMALPort, len(this.clock))
	for i, port := range this.clock {
		ports[i] = port
	}
	return ports
}
