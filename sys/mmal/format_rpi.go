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

	// Framerworks
	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *format) String() string {
	return fmt.Sprintf("<sys.hw.mmal.format>{ type=%v }", this.Type())
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *format) Type() hw.MMALFormatType {
	switch rpi.MMALStreamFormatType(this.handle) {
	case rpi.MMAL_STREAM_TYPE_UNKNOWN:
		return hw.MMAL_FORMAT_UNKNOWN
	case rpi.MMAL_STREAM_TYPE_CONTROL:
		return hw.MMAL_FORMAT_CONTROL
	case rpi.MMAL_STREAM_TYPE_AUDIO:
		return hw.MMAL_FORMAT_AUDIO
	case rpi.MMAL_STREAM_TYPE_VIDEO:
		return hw.MMAL_FORMAT_VIDEO
	case rpi.MMAL_STREAM_TYPE_SUBPICTURE:
		return hw.MMAL_FORMAT_SUBPICTURE
	default:
		return hw.MMAL_FORMAT_UNKNOWN
	}
}
