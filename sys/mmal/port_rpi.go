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
// CLOSE

func (this *port) Close() error {
	this.log.Debug("<sys.hw.mmal.port>Close{ name='%v' }", this.Name())

	err := new(errors.CompoundError)

	// Check for already closed
	if this.handle == nil {
		return gopi.ErrOutOfOrder
	}

	// Close the channel
	close(this.lock)

	// Disable port
	if rpi.MMALPortIsEnabled(this.handle) {
		if err_ := rpi.MMALPortDisable(this.handle); err_ != nil {
			err.Add(err_)
		}
	}

	// Deregister callback
	rpi.MMALPortDeregisterCallback(this.handle)

	// Destroy pool and queue
	if this.pool != nil {
		rpi.MMALPortPoolDestroy(this.handle, this.pool)
		this.pool = nil
	}
	if this.queue != nil {
		rpi.MMALQueueDestroy(this.queue)
		this.queue = nil
	}

	// Release resources
	this.handle = nil

	// Return errors
	return err.ErrorOrSelf()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *port) String() string {
	if this.handle == nil {
		return fmt.Sprintf("<sys.hw.mmal.port>{ nil }")
	} else {
		parts := ""
		parts += fmt.Sprintf("name='%v' type=%v enabled=%v ", this.Name(), rpi.MMALPortType(this.handle), this.Enabled())
		if cap := rpi.MMALPortCapabilities(this.handle); cap != 0 {
			parts += fmt.Sprintf("capabilities=%v ", cap)
		}
		parts += fmt.Sprintf("format=%v pool=%v", this.Format(), rpi.MMALPoolString(this.pool))
		return fmt.Sprintf("<sys.hw.mmal.port>{ %v}", parts)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *port) Name() string {
	return rpi.MMALPortName(this.handle)
}

func (this *port) Enabled() bool {
	return rpi.MMALPortIsEnabled(this.handle)
}

func (this *port) SetEnabled(value bool) error {
	this.log.Debug2("<sys.hw.mmal.port>SetEnabled{ name='%v' value=%v }", this.Name(), value)

	if value {
		// Resize the pool of buffers
		if this.pool != nil {
			buffer_size := uint32(0)
			if this.CapabilityAllocation() {
				buffer_size, _, _ = rpi.MMALPortBufferSize(this.handle)
			}
			buffer_num, _, _ := rpi.MMALPortBufferNum(this.handle)
			this.log.Debug2("<sys.hw.mmal.port>MMALPoolResize{ num=%v size=%v }", buffer_num, buffer_size)
			if err := rpi.MMALPoolResize(this.pool, buffer_num, buffer_size); err != nil {
				return err
			}
		}
		if err := rpi.MMALPortEnable(this.handle); err != nil {
			return err
		}
	} else {
		if err := rpi.MMALPortDisable(this.handle); err != nil {
			return err
		} else if rpi.MMALPortType(this.handle) == rpi.MMAL_PORT_TYPE_OUTPUT {
			// Flush output port
			for {
				if buffer := rpi.MMALQueueGet(this.queue); buffer == nil {
					break
				} else {
					rpi.MMALBufferRelease(buffer)
				}
			}
		}
	}

	// return success
	return nil
}

func (this *port) Flush() error {
	this.log.Debug2("<sys.hw.mmal.port>Flush{ name='%v' }", this.Name())

	return rpi.MMALPortFlush(this.handle)
}

func (this *port) CapabilityPassthrough() bool {
	return rpi.MMALPortCapabilities(this.handle)&rpi.MMAL_PORT_CAPABILITY_PASSTHROUGH != 0
}

func (this *port) CapabilityAllocation() bool {
	return rpi.MMALPortCapabilities(this.handle)&rpi.MMAL_PORT_CAPABILITY_ALLOCATION != 0
}

func (this *port) CapabilitySupportsEventFormatChange() bool {
	return rpi.MMALPortCapabilities(this.handle)&rpi.MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE != 0
}

func (this *port) CopyFormat(other hw.MMALFormat) error {
	this.log.Debug2("<sys.hw.mmal.port>CopyFormat{ name='%v' src_format=%v }", this.Name(), other)
	if other_, ok := other.(*format); ok == false {
		return gopi.ErrBadParameter
	} else if err := rpi.MMALStreamFormatFullCopy(rpi.MMALPortFormat(this.handle), other_.handle); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *port) CommitFormatChange() error {
	this.log.Debug2("<sys.hw.mmal.port>CommitFormatChange{ name='%v' }", this.Name())
	if err := rpi.MMALPortFormatCommit(this.handle); err != nil {
		return err
	}

	// Change buffer parameters
	_, buffer_num_min, buffer_num_recommended := rpi.MMALPortBufferNum(this.handle)
	_, buffer_size_min, buffer_size_recommended := rpi.MMALPortBufferSize(this.handle)

	// Determine number
	if buffer_num_recommended == 0 {
		buffer_num_recommended = buffer_num_min
	}
	if buffer_num_recommended == 0 {
		buffer_num_recommended = 1
	}

	// Determine size
	if buffer_size_recommended == 0 {
		buffer_size_recommended = buffer_size_min
	}

	// Set buffer parameters
	this.log.Debug2("<sys.hw.mmal.port>CommitFormatChange{ buffer_num=%v buffer_size=%v }", buffer_num_recommended, buffer_size_recommended)
	rpi.MMALPortBufferSet(this.handle, buffer_num_recommended, buffer_size_recommended)

	// Success
	return nil
}

func (this *port) Connect(other hw.MMALPort) error {
	this.log.Debug2("<sys.hw.mmal.port>Connect{ name='%v' other='%v' }", this.Name(), other.Name())
	if other_, ok := other.(*port); ok == false {
		return gopi.ErrBadParameter
	} else {
		return rpi.MMALPortConnect(this.handle, other_.handle)
	}
}

func (this *port) Disconnect() error {
	this.log.Debug2("<sys.hw.mmal.port>Disconnect{ name='%v' }", this.Name())
	return rpi.MMALPortDisconnect(this.handle)
}

func (this *port) Format() hw.MMALFormat {
	return this.NewFormat()
}

func (this *port) VideoFormat() hw.MMALVideoFormat {
	format := this.NewFormat()
	if format.Type() != hw.MMAL_FORMAT_VIDEO {
		return nil
	} else {
		return format
	}
}

func (this *port) AudioFormat() hw.MMALAudioFormat {
	format := this.NewFormat()
	if format.Type() != hw.MMAL_FORMAT_AUDIO {
		return nil
	} else {
		return format
	}
}

func (this *port) SubpictureFormat() hw.MMALSubpictureFormat {
	format := this.NewFormat()
	if format.Type() != hw.MMAL_FORMAT_SUBPICTURE {
		return nil
	} else {
		return format
	}
}

func (this *port) Send(b hw.MMALBuffer) error {
	if buffer_, ok := b.(*buffer); ok == false {
		return gopi.ErrBadParameter
	} else {
		this.log.Debug2("<sys.hw.mmal.port>Send{ name='%v' buffer=%v }", this.Name(), rpi.MMALBufferString(buffer_.handle))
		return rpi.MMALPortSendBuffer(this.handle, buffer_.handle)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *port) NewFormat() *format {
	return &format{this.log, rpi.MMALPortFormat(this.handle)}
}
