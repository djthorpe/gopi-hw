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
	"strings"

	// Framerworks
	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *format) String() string {
	parts := ""
	parts += fmt.Sprintf("type=%v", this.Type())
	return fmt.Sprintf("<sys.hw.mmal.format>{ %v }", strings.TrimSpace(parts))
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION - STREAM FORMAT

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

func (this *format) Bitrate() uint32 {
	return rpi.MMALStreamFormatBitrate(this.handle)
}

func (this *format) SetBitrate(value uint32) {
	rpi.MMALStreamFormatSetBitrate(this.handle, value)
}

func (this *format) Encoding() (hw.MMALEncodingType, hw.MMALEncodingType) {
	return rpi.MMALStreamFormatEncoding(this.handle)
}

func (this *format) SetEncoding(value hw.MMALEncodingType) {
	rpi.MMALStreamFormatSetEncoding(this.handle, value, 0)
}

func (this *format) SetEncodingVariant(value, variant hw.MMALEncodingType) {
	rpi.MMALStreamFormatSetEncoding(this.handle, value, variant)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION - VIDEO STREAM FORMAT

func (this *format) WidthHeight() (uint32, uint32) {
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		return rpi.MMALStreamFormatVideoWidthHeight(this.handle)
	} else {
		return 0, 0
	}
}

func (this *format) SetWidthHeight(w, h uint32) {
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		rpi.MMALStreamFormatVideoSetWidthHeight(this.handle, w, h)
	}
}
