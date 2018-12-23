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
	return fmt.Sprintf("<sys.hw.mmal.port>{ name='%v' type=%v enabled=%v capabilities=%v format=%v }", this.Name(), rpi.MMALPortType(this.handle), this.Enabled(), rpi.MMALPortCapabilities(this.handle), this.Format())
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
	if value {
		return rpi.MMALPortEnable(this.handle)
	} else {
		return rpi.MMALPortDisable(this.handle)
	}
}

func (this *port) Flush() error {
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
	return rpi.MMALPortFormatCommit(this.handle)
}

func (this *port) Connect(other hw.MMALPort) error {
	if other_, ok := other.(*port); ok == false {
		return gopi.ErrBadParameter
	} else {
		return rpi.MMALPortConnect(this.handle, other_.handle)
	}
}

func (this *port) Disconnect() error {
	return rpi.MMALPortDisconnect(this.handle)
}

func (this *port) Format() hw.MMALFormat {
	return &format{rpi.MMALPortFormat(this.handle)}
}
