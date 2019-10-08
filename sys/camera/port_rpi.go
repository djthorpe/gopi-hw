// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package camera

import (
	hw "github.com/djthorpe/gopi-hw"
)

////////////////////////////////////////////////////////////////////////////////

const (
	CAMERA_PREVIEW_OUTPORT = 0
	CAMERA_VIDEO_OUTPORT   = 1
	CAMERA_IMAGE_OUTPORT   = 2
	PREVIEW_INPORT         = 0
	ENCODER_INPORT         = 0
	ENCODER_OUTPORT        = 0
)

////////////////////////////////////////////////////////////////////////////////

func (this *camera) CameraPreviewOutPort() hw.MMALPort {
	return this.camera.Outputs()[CAMERA_PREVIEW_OUTPORT]
}

func (this *camera) CameraImageOutPort() hw.MMALPort {
	return this.camera.Outputs()[CAMERA_IMAGE_OUTPORT]
}

func (this *camera) CameraVideoOutPort() hw.MMALPort {
	return this.camera.Outputs()[CAMERA_VIDEO_OUTPORT]
}

func (this *camera) CameraControlPort() hw.MMALPort {
	return this.camera.Control()
}

func (this *camera) PreviewInPort() hw.MMALPort {
	return this.renderer.Inputs()[PREVIEW_INPORT]
}

func (this *camera) ImageEncoderInPort() hw.MMALPort {
	return this.image_encoder.Inputs()[ENCODER_INPORT]
}

func (this *camera) ImageEncoderOutPort() hw.MMALPort {
	return this.image_encoder.Outputs()[ENCODER_OUTPORT]
}

func (this *camera) VideoEncoderInPort() hw.MMALPort {
	return this.video_encoder.Inputs()[ENCODER_INPORT]
}

func (this *camera) VideoEncoderOutPort() hw.MMALPort {
	return this.video_encoder.Outputs()[ENCODER_OUTPORT]
}
