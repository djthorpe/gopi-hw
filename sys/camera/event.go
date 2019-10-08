/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package camera

import (
	"fmt"

	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type camera_event struct {
	source gopi.Driver
	data   []byte
	flags  hw.CameraDataFlag
}

////////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

func NewEvent(source gopi.Driver, data []byte, flags hw.CameraDataFlag) hw.CameraEvent {
	return &camera_event{source, data, flags}
}

////////////////////////////////////////////////////////////////////////////////
// EVENT IMPLEMENTATION

func (this *camera_event) Name() string {
	return "CameraEvent"
}

func (this *camera_event) Source() gopi.Driver {
	return this.source
}

func (this *camera_event) Data() []byte {
	return this.data
}

func (this *camera_event) Flags() hw.CameraDataFlag {
	return this.flags
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *camera_event) String() string {
	return fmt.Sprintf("<media.camera.CameraEvent>{ data=%v bytes flags=%v }", len(this.data), this.flags)
}
