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
	this.log.Debug("<sys.hw.mmal.component>Close{ name='%v' }", this.Name())

	err := new(errors.CompoundError)

	if this.handle == nil {
		// Already closed
		return gopi.ErrOutOfOrder
	}

	// Disable ports
	for _, port := range this.input {
		if err_ := port.Close(); err_ != nil {
			err.Add(err_)
		}
	}
	for _, port := range this.output {
		if err_ := port.Close(); err_ != nil {
			err.Add(err_)
		}
	}
	for _, port := range this.clock {
		if err_ := port.Close(); err_ != nil {
			err.Add(err_)
		}
	}
	if err_ := this.control.Close(); err_ != nil {
		err.Add(err_)
	}

	// Destroy component
	if err_ := rpi.MMALComponentDestroy(this.handle); err_ != nil {
		err.Add(err_)
	}

	// Release resources
	this.handle = nil
	this.control = nil
	this.input = nil
	this.output = nil
	this.clock = nil
	this.port_map = nil

	return err.ErrorOrSelf()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *component) String() string {
	if this.handle == nil {
		return fmt.Sprintf("<sys.hw.mmal.component>{ nil }")
	} else {
		return fmt.Sprintf("<sys.hw.mmal.component>{ name='%v' id=%08X enabled=%v control=%v input=%v output=%v clock=%v }", this.Name(), this.Id(), this.Enabled(), this.control, this.input, this.output, this.clock)
	}
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
	this.log.Debug2("<sys.hw.mmal.component>SetEnabled{ name='%v' value=%v }", this.Name(), value)

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
	this.log.Debug2("<sys.hw.mmal.component>Acquire{ name='%v' }", this.Name())

	if this.handle == nil {
		// Component is not open
		return gopi.ErrOutOfOrder
	} else {
		return rpi.MMALComponentAcquire(this.handle)
	}
}

func (this *component) Release() error {
	this.log.Debug2("<sys.hw.mmal.component>Release{ name='%v' }", this.Name())

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

func (this *component) Inputs() []hw.MMALPort {
	ports := make([]hw.MMALPort, len(this.input))
	for i, port := range this.input {
		ports[i] = port
	}
	return ports
}

func (this *component) Outputs() []hw.MMALPort {
	ports := make([]hw.MMALPort, len(this.output))
	for i, port := range this.output {
		ports[i] = port
	}
	return ports
}

func (this *component) Clocks() []hw.MMALPort {
	ports := make([]hw.MMALPort, len(this.clock))
	for i, port := range this.clock {
		ports[i] = port
	}
	return ports
}

////////////////////////////////////////////////////////////////////////////////
// BUFFERS

func (this *component) GetEmptyBufferOnPort(p hw.MMALPort, blocking bool) (hw.MMALBuffer, error) {
	this.log.Debug2("<sys.hw.mmal.component>GetEmptyBufferOnPort{ name='%v' port=%v blocking=%v }", this.Name(), p, blocking)
	if port_, ok := p.(*port); ok == false {
		return nil, gopi.ErrBadParameter
	} else if _, exists := this.port_map[port_.handle]; exists == false {
		return nil, fmt.Errorf("Port is invalid: %v", p)
	} else if pool := port_.pool; pool == nil {
		return nil, fmt.Errorf("Pool is invalid on port: %v", p)
	} else {
		// Get a buffer from the pool queue
		for {
			if buffer_handle := rpi.MMALPoolGetBuffer(pool); buffer_handle != nil {
				return &buffer{this.log, buffer_handle}, nil
			} else if blocking == false {
				return nil, nil
			} else {
				this.log.Debug2("GetEmptyBufferOnPort: Waiting for empty buffer to become available")
				<-port_.lock
			}
		}
	}
}

func (this *component) GetFullBufferOnPort(p hw.MMALPort, blocking bool) (hw.MMALBuffer, error) {
	this.log.Debug2("<sys.hw.mmal.component>GetFullBufferOnPort{ name='%v' port=%v blocking=%v }", this.Name(), p, blocking)
	if port_, ok := p.(*port); ok == false {
		return nil, gopi.ErrBadParameter
	} else if _, exists := this.port_map[port_.handle]; exists == false {
		return nil, fmt.Errorf("Port is invalid: %v", p)
	} else if queue := port_.queue; queue == nil {
		return nil, fmt.Errorf("Pool is invalid: %v", p)
	} else {
		// Get a buffer from the 'full queue'
		for {
			if buffer_handle := rpi.MMALQueueGet(queue); buffer_handle != nil {
				this.log.Debug2("GetFullBufferOnPort: got %v", rpi.MMALBufferString(buffer_handle))
				return &buffer{this.log, buffer_handle}, nil
			} else if blocking == false {
				return nil, nil
			} else {
				this.log.Debug2("GetFullBufferOnPort: Waiting for full buffer to become available")
				<-port_.lock
			}
		}
	}
}
