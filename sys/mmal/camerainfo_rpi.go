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

	// Frameworks
	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

type camerainfo struct {
	handle rpi.MMAL_CameraInfo
}

type camera struct {
	handle rpi.MMAL_Camera
}

func (this *camerainfo) String() string {
	parts := ""
	parts += fmt.Sprintf("num_cameras=%v ", rpi.MMALCameraInfoGetCamerasNum(this.handle))
	parts += fmt.Sprintf("num_flashes=%v ", rpi.MMALCameraInfoGetFlashesNum(this.handle))
	parts += fmt.Sprintf("cameras=%v ", this.Cameras())
	parts += fmt.Sprintf("flashes=%v ", this.Flashes())
	return fmt.Sprintf("<sys.hw.mmal.camerainfo>{ %v }", strings.TrimSpace(parts))
}

func (this *camera) String() string {
	parts := ""
	parts += fmt.Sprintf("id=%v ", rpi.MMALCameraInfoGetCameraId(this.handle))
	parts += fmt.Sprintf("name=\"%v\" ", rpi.MMALCameraInfoGetCameraName(this.handle))
	parts += fmt.Sprintf("size={ %v, %v } ", rpi.MMALCameraInfoGetCameraMaxWidth(this.handle), rpi.MMALCameraInfoGetCameraMaxHeight(this.handle))
	parts += fmt.Sprintf("lens_present=%v ", rpi.MMALCameraInfoGetCameraLensPresent(this.handle))
	return fmt.Sprintf("<sys.hw.mmal.camera>{ %v }", strings.TrimSpace(parts))
}

func (this *camerainfo) Flashes() []hw.MMALCameraFlashType {
	return rpi.MMALCameraInfoGetFlashes(this.handle)
}

func (this *camerainfo) Cameras() []hw.MMALCamera {
	cameras := rpi.MMALCameraInfoGetCameras(this.handle)
	cameras_ := make([]hw.MMALCamera, len(cameras))
	for i := 0; i < len(cameras_); i++ {
		cameras_[i] = &camera{cameras[i]}
	}
	return cameras_
}

func (this *camera) Id() uint32 {
	return rpi.MMALCameraInfoGetCameraId(this.handle)
}

func (this *camera) Name() string {
	return rpi.MMALCameraInfoGetCameraName(this.handle)
}

func (this *camera) Size() (uint32, uint32) {
	return rpi.MMALCameraInfoGetCameraMaxWidth(this.handle), rpi.MMALCameraInfoGetCameraMaxHeight(this.handle)
}

func (this *camera) LensPresent() bool {
	return rpi.MMALCameraInfoGetCameraLensPresent(this.handle)
}
