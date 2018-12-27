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
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *port) String() string {
	return fmt.Sprintf("<sys.hw.mmal.port>{ name='%v' type=%v enabled=%v capabilities=%v format=%v pool=%v }", this.Name(), rpi.MMALPortType(this.handle), this.Enabled(), rpi.MMALPortCapabilities(this.handle), this.Format(), this.pool)
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
		return rpi.MMALPortEnable(this.handle)
	} else {
		return rpi.MMALPortDisable(this.handle)
	}
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

func (this *port) CommitFormatChange() error {
	this.log.Debug2("<sys.hw.mmal.port>CommitFormatChange{ name='%v' }", this.Name())
	return rpi.MMALPortFormatCommit(this.handle)
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

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *port) NewFormat() *format {
	return &format{this.log, rpi.MMALPortFormat(this.handle)}
}
