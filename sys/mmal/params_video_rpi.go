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

	// Frameworks
	"github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// DISPLAY REGION

func (this *port) GetDisplayRegion() (hw.MMALDisplayRegion, error) {
	if value, err := rpi.MMALPortParameterGetDisplayRegion(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_DISPLAYREGION); err != nil {
		return nil, err
	} else {
		return &displayregion{value}, nil
	}
}

func (this *port) SetDisplayRegion(value hw.MMALDisplayRegion) error {
	if value_, ok := value.(*displayregion); ok == false {
		return gopi.ErrBadParameter
	} else {
		return rpi.MMALPortParameterSetDisplayRegion(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_DISPLAYREGION, value_.handle)
	}
}
