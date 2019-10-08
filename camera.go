/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2018-2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	"errors"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// ERRORS

var (
	ErrInvalidCameraId = errors.New("Invalid or missing camera")
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Camera represents the ability to capture images and video from a
// Raspberry Pi camera device
type Camera interface {
	gopi.Driver
	gopi.Publisher

	// CameraId returns the Camera ID, usually zero
	CameraId() uint32

	// CameraModel returns the model name of the camera
	CameraModel() string

	// CameraFrameSize returns the number of pixels a frame contains
	// from the camera
	CameraFrameSize() gopi.Size

	// SupportedImageEncodings returns an array of supported image capture
	// formats, or nil on error
	SupportedImageEncodings() []MMALEncodingType

	// SupportedVideoEncodings returns an array of supported video capture
	// formats, or nil on error
	SupportedVideoEncodings() []MMALEncodingType

	// CameraConfig reads the current camera configuration
	CameraConfig() (CameraConfig, error)

	// SetCameraConfig sets any changed parameters for the camera configuration
	SetCameraConfig(CameraConfig) error

	// Preview starts and stops the preview of the camera image. Setting argument
	// to true begins the preview, and setting to false stops the preview.
	Preview(bool) error

	// ImageCapture blocks whilst the camera image is captured
	ImageCapture() error
}

// CameraConfig represents all the configuration parameters for camera image
// and video capture. If you want to set any parameter, then you need to
// set a flag for that parameter before calling SetCameraConfig
type CameraConfig struct {
	PreviewRotation            int32
	ImageRotation              int32
	VideoRotation              int32
	ImageFormatEncoding        MMALEncodingType
	ImageFormatEncodingVariant MMALEncodingType
	ImageFrameSize             gopi.Size
	PreviewFrameSize           gopi.Size
	ImageJPEGQuality           uint32
	Flags                      CameraConfigFlag
}

type CameraEvent interface {
	gopi.Event

	// Data stream from the camera
	Data() []byte

	// Flags returns information about the event: Audio, image or
	// video data, and whether at the beginning or end of the stream
	Flags() CameraDataFlag
}

// CameraConfigFlag represents a set of flags of configuration parameters
// which can be changed
type CameraConfigFlag uint64

// CameraDataFlag represents the type of data stream from the camera, and
// whether this is start of stream or end of stream
type CameraDataFlag uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FLAG_IMAGE_ROTATION CameraConfigFlag = (1 << iota)
	FLAG_PREVIEW_ROTATION
	FLAG_VIDEO_ROTATION
	FLAG_IMAGE_ENCODING_FORMAT
	FLAG_IMAGE_FRAMESIZE
	FLAG_PREVIEW_FRAMESIZE
	FLAG_IMAGE_ENCODING_JPEGQUALITY
	FLAG_NONE     CameraConfigFlag = 0
	FLAG_ROTATION CameraConfigFlag = FLAG_IMAGE_ROTATION | FLAG_PREVIEW_ROTATION | FLAG_VIDEO_ROTATION
)

const (
	FLAG_DATA_IMAGE CameraDataFlag = (1 << iota)
	FLAG_DATA_VIDEO
	FLAG_DATA_AUDIO
	FLAG_DATA_STREAM_START
	FLAG_DATA_STREAM_END
	FLAG_DATA_NONE CameraDataFlag = 0
	FLAG_DATA_MIN                 = FLAG_DATA_IMAGE
	FLAG_DATA_MAX                 = FLAG_DATA_STREAM_END
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// HasFlags returns true if all flags in argument are set
func (config CameraConfig) HasFlags(flag CameraConfigFlag) bool {
	return config.Flags&flag == flag
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f CameraDataFlag) String() string {
	if f == FLAG_DATA_NONE {
		return "FLAG_DATA_NONE"
	}
	parts := ""
	for flag := FLAG_DATA_MIN; flag <= FLAG_DATA_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case FLAG_DATA_IMAGE:
			parts += "|" + "FLAG_DATA_IMAGE"
		case FLAG_DATA_VIDEO:
			parts += "|" + "FLAG_DATA_VIDEO"
		case FLAG_DATA_AUDIO:
			parts += "|" + "FLAG_DATA_AUDIO"
		case FLAG_DATA_STREAM_START:
			parts += "|" + "FLAG_DATA_STREAM_START"
		case FLAG_DATA_STREAM_END:
			parts += "|" + "FLAG_DATA_STREAM_END"
		default:
			parts += "|" + "[?? Invalid CameraDataFlag value]"
		}
	}
	return strings.Trim(parts, "|")
}
