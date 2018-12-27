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
	this.log.Debug2("<sys.hw.mmal.format>SetBitrate{ value=%v }", value)
	rpi.MMALStreamFormatSetBitrate(this.handle, value)
}

func (this *format) Encoding() (hw.MMALEncodingType, hw.MMALEncodingType) {
	return rpi.MMALStreamFormatEncoding(this.handle)
}

func (this *format) SetEncoding(value hw.MMALEncodingType) {
	this.log.Debug2("<sys.hw.mmal.format>SetEncoding{ value=%v }", value)
	rpi.MMALStreamFormatSetEncoding(this.handle, value, 0)
}

func (this *format) SetEncodingVariant(value, variant hw.MMALEncodingType) {
	this.log.Debug2("<sys.hw.mmal.format>SetEncodingVariant{ value=%v variant=%v }", value, variant)
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
	this.log.Debug2("<sys.hw.mmal.format>SetWidthHeight{ w=%v h=%v }", w, h)
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		rpi.MMALStreamFormatVideoSetWidthHeight(this.handle, w, h)
	}
}

func (this *format) Crop() hw.MMALRect {
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		return rpi.MMALStreamFormatVideoCrop(this.handle)
	} else {
		return hw.MMALRect{}
	}
}
func (this *format) SetCrop(value hw.MMALRect) {
	this.log.Debug2("<sys.hw.mmal.format>SetCrop{ value=%v }", value)
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		rpi.MMALStreamFormatVideoSetCrop(this.handle, value)
	}
}

func (this *format) FrameRate() hw.MMALRationalNum {
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		return rpi.MMALStreamFormatVideoFrameRate(this.handle)
	} else {
		return hw.MMALRationalNum{}
	}
}
func (this *format) SetFrameRate(value hw.MMALRationalNum) {
	this.log.Debug2("<sys.hw.mmal.format>SetFrameRate{ value=%v }", value)
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		rpi.MMALStreamFormatVideoSetFrameRate(this.handle, value)
	}
}

func (this *format) PixelAspectRatio() hw.MMALRationalNum {
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		return rpi.MMALStreamFormatVideoPixelAspectRatio(this.handle)
	} else {
		return hw.MMALRationalNum{}
	}
}

func (this *format) SetPixelAspectRatio(value hw.MMALRationalNum) {
	this.log.Debug2("<sys.hw.mmal.format>SetPixelAspectRatio{ value=%v }", value)
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		rpi.MMALStreamFormatVideoSetPixelAspectRatio(this.handle, value)
	}
}

func (this *format) ColorSpace() hw.MMALEncodingType {
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		return rpi.MMALStreamFormatVideoColorSpace(this.handle)
	} else {
		return 0
	}
}

func (this *format) SetColorSpace(value hw.MMALEncodingType) {
	this.log.Debug2("<sys.hw.mmal.format>SetColorSpace{ value=%v }", value)
	if rpi.MMALStreamFormatType(this.handle) == rpi.MMAL_STREAM_TYPE_VIDEO {
		rpi.MMALStreamFormatVideoSetColorSpace(this.handle, value)
	}
}
