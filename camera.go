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

type Camera interface {
	gopi.Driver

	// Get Camera Properties
	CameraId() uint32
	CameraModel() string
	CameraFrameSize() gopi.Size

	// Read and write the camera configuration
	CameraConfig() (CameraConfig, error)
	SetCameraConfig(CameraConfig) error

	// Preview, image and video capture
	Preview() error
	ImageCapture() error
}

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

type CameraConfigFlag uint64

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
	FLAG_NONE     = 0
	FLAG_ROTATION = FLAG_IMAGE_ROTATION | FLAG_PREVIEW_ROTATION | FLAG_VIDEO_ROTATION
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// HasFlags returns true if all flags in argument are set
func (config CameraConfig) HasFlags(flag CameraConfigFlag) bool {
	return config.Flags&flag == flag
}
