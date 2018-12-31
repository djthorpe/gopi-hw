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
	"github.com/djthorpe/gopi"

	// Frameworks

	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

// PARAMETER_CAMERA_INFO
func (this *port) CameraInfo() (hw.MMALCameraInfo, error) {
	if value, err := rpi.MMALPortParameterGetCameraInfo(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_INFO); err != nil {
		return nil, err
	} else {
		return &camerainfo{value}, nil
	}
}

// MMAL_PARAMETER_CAPTURE_EXPOSURE_COMP
func (this *port) CaptureExposureComp() (int32, error) {
	return rpi.MMALPortParameterGetInt32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAPTURE_EXPOSURE_COMP)
}

func (this *port) SetCaptureExposureComp(value int32) error {
	return rpi.MMALPortParameterSetInt32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAPTURE_EXPOSURE_COMP, value)
}

// MMAL_PARAMETER_OUTPUT_SHIFT
func (this *port) OutputShift() (int32, error) {
	return rpi.MMALPortParameterGetInt32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_OUTPUT_SHIFT)
}

func (this *port) SetOutputShift(value int32) error {
	return rpi.MMALPortParameterSetInt32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_OUTPUT_SHIFT, value)
}

// MMAL_PARAMETER_CCM_SHIFT
func (this *port) CCMShift() (int32, error) {
	return rpi.MMALPortParameterGetInt32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CCM_SHIFT)
}
func (this *port) SetCCMShift(value int32) error {
	return rpi.MMALPortParameterSetInt32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CCM_SHIFT, value)
}

// MMAL_PARAMETER_ROTATION
func (this *port) Rotation() (int32, error) {
	return rpi.MMALPortParameterGetInt32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ROTATION)
}
func (this *port) SetRotation(value int32) error {
	return rpi.MMALPortParameterSetInt32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ROTATION, value)
}

// MMAL_PARAMETER_CAMERA_NUM
func (this *port) CameraNum() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_NUM)
}
func (this *port) SetCameraNum(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_NUM, value)
}

// MMAL_PARAMETER_JPEG_Q_FACTOR
func (this *port) JPEGQFactor() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_JPEG_Q_FACTOR)
}
func (this *port) SetJPEGQFactor(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_JPEG_Q_FACTOR, value)
}

// MMAL_PARAMETER_ISO
func (this *port) ISO() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ISO)
}
func (this *port) SetISO(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ISO, value)
}

// MMAL_PARAMETER_CAMERA_MIN_ISO
func (this *port) MinISO() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_MIN_ISO)
}
func (this *port) SetMinISO(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_MIN_ISO, value)
}

// MMAL_PARAMETER_CAMERA_CUSTOM_SENSOR_CONFIG
func (this *port) CustomSensorConfig() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_CUSTOM_SENSOR_CONFIG)
}
func (this *port) SetCustomSensorConfig(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_CUSTOM_SENSOR_CONFIG, value)
}

// MMAL_PARAMETER_SHUTTER_SPEED
func (this *port) ShutterSpeed() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SHUTTER_SPEED)
}
func (this *port) SetShutterSpeed(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SHUTTER_SPEED, value)
}

// MMAL_PARAMETER_DPF_CONFIG
func (this *port) DPFConfig() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_DPF_CONFIG)
}
func (this *port) SetDPFConfig(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_DPF_CONFIG, value)
}

// MMAL_PARAMETER_JPEG_RESTART_INTERVAL
func (this *port) JPEGRestartInterval() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_JPEG_RESTART_INTERVAL)
}
func (this *port) SetJPEGRestartInterval(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_JPEG_RESTART_INTERVAL, value)
}

// MMAL_PARAMETER_CAMERA_ISP_BLOCK_OVERRIDE
func (this *port) CameraISPBlockOverride() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_ISP_BLOCK_OVERRIDE)
}
func (this *port) SetCameraISPBlockOverride(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_ISP_BLOCK_OVERRIDE, value)
}

// MMAL_PARAMETER_BLACK_LEVEL
func (this *port) BlackLevel() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_BLACK_LEVEL)
}
func (this *port) SetBlackLevel(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_BLACK_LEVEL, value)
}

// MMAL_PARAMETER_EXIF_DISABLE
func (this *port) EXIFDisable() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_EXIF_DISABLE)
}
func (this *port) SetEXIFDisable(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_EXIF_DISABLE, value)
}

// MMAL_PARAMETER_CAPTURE
func (this *port) Capture() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAPTURE)
}
func (this *port) SetCapture(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAPTURE, value)
}

// MMAL_PARAMETER_DRAW_BOX_FACES_AND_FOCUS
func (this *port) DrawBoxFacesAndFocus() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_DRAW_BOX_FACES_AND_FOCUS)
}
func (this *port) SetDrawBoxFacesAndFocus(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_DRAW_BOX_FACES_AND_FOCUS, value)
}

// MMAL_PARAMETER_VIDEO_STABILISATION
func (this *port) VideoStabilisation() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_VIDEO_STABILISATION)
}
func (this *port) SetVideoStabilisation(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_VIDEO_STABILISATION, value)
}

// MMAL_PARAMETER_ENABLE_RAW_CAPTURE
func (this *port) EnableRAWCapture() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ENABLE_RAW_CAPTURE)
}
func (this *port) SetEnableRAWCapture(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ENABLE_RAW_CAPTURE, value)
}

// MMAL_PARAMETER_ENABLE_DPF_FILE
func (this *port) EnableDPFFile() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ENABLE_DPF_FILE)
}
func (this *port) SetEnableDPFFile(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ENABLE_DPF_FILE, value)
}

// MMAL_PARAMETER_DPF_FAIL_IS_FATAL
func (this *port) DPFFailIsFatal() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_DPF_FAIL_IS_FATAL)
}
func (this *port) SetDPFFailIsFatal(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_DPF_FAIL_IS_FATAL, value)
}

// MMAL_PARAMETER_HIGH_DYNAMIC_RANGE
func (this *port) HighDynamicRange() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_HIGH_DYNAMIC_RANGE)
}
func (this *port) SetHighDynamicRange(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_HIGH_DYNAMIC_RANGE, value)
}

// MMAL_PARAMETER_ANTISHAKE
func (this *port) Antishake() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ANTISHAKE)
}
func (this *port) SetAntishake(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ANTISHAKE, value)
}

// MMAL_PARAMETER_CAMERA_BURST_CAPTURE
func (this *port) BurstCapture() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_BURST_CAPTURE)
}
func (this *port) SetBurstCapture(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAMERA_BURST_CAPTURE, value)
}

// MMAL_PARAMETER_CAPTURE_STATS_PASS
func (this *port) CaptureStatsPass() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAPTURE_STATS_PASS)
}
func (this *port) SetCaptureStatsPass(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CAPTURE_STATS_PASS, value)
}

// MMAL_PARAMETER_ENABLE_REGISTER_FILE
func (this *port) EnableRegisterFile() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ENABLE_REGISTER_FILE)
}
func (this *port) SetEnableRegisterFile(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ENABLE_REGISTER_FILE, value)
}

// MMAL_PARAMETER_REGISTER_FAIL_IS_FATAL
func (this *port) RegisterFailIsFatal() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_REGISTER_FAIL_IS_FATAL)
}
func (this *port) SetRegisterFailIsFatal(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_REGISTER_FAIL_IS_FATAL, value)
}

// MMAL_PARAMETER_JPEG_ATTACH_LOG
func (this *port) JPEGAttachLog() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_JPEG_ATTACH_LOG)
}
func (this *port) SetJPEGAttachLog(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_JPEG_ATTACH_LOG, value)
}

// MMAL_PARAMETER_SW_SHARPEN_DISABLE
func (this *port) SWSharpenDisable() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SW_SHARPEN_DISABLE)
}
func (this *port) SetSWSharpenDisable(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SW_SHARPEN_DISABLE, value)
}

// MMAL_PARAMETER_FLASH_REQUIRED
func (this *port) FlashRequired() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_FLASH_REQUIRED)
}
func (this *port) SetFlashRequired(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_FLASH_REQUIRED, value)
}

// MMAL_PARAMETER_SW_SATURATION_DISABLE
func (this *port) SWSaturationDisable() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SW_SATURATION_DISABLE)
}
func (this *port) SetSWSaturationDisable(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SW_SATURATION_DISABLE, value)
}

// MMAL_PARAMETER_VIDEO_DENOISE
func (this *port) VideoDenoise() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_VIDEO_DENOISE)
}
func (this *port) SetVideoDenoise(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_VIDEO_DENOISE, value)
}

// MMAL_PARAMETER_STILLS_DENOISE
func (this *port) StillsDenoise() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_STILLS_DENOISE)
}
func (this *port) SetStillsDenoise(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_STILLS_DENOISE, value)
}

// MMAL_PARAMETER_EXPOSURE_COMP
func (this *port) ExposureComp() (hw.MMALRationalNum, error) {
	return rpi.MMALPortParameterGetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_EXPOSURE_COMP)
}
func (this *port) SetExposureComp(value hw.MMALRationalNum) error {
	return rpi.MMALPortParameterSetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_EXPOSURE_COMP, value)
}

// MMAL_PARAMETER_SHARPNESS
func (this *port) Sharpness() (hw.MMALRationalNum, error) {
	return rpi.MMALPortParameterGetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SHARPNESS)
}
func (this *port) SetSharpness(value hw.MMALRationalNum) error {
	return rpi.MMALPortParameterSetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SHARPNESS, value)
}

// MMAL_PARAMETER_CONTRAST
func (this *port) Contrast() (hw.MMALRationalNum, error) {
	return rpi.MMALPortParameterGetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CONTRAST)

}
func (this *port) SetContrast(value hw.MMALRationalNum) error {
	return rpi.MMALPortParameterSetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_CONTRAST, value)

}

// MMAL_PARAMETER_BRIGHTNESS
func (this *port) Brightness() (hw.MMALRationalNum, error) {
	return rpi.MMALPortParameterGetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_BRIGHTNESS)
}
func (this *port) SetBrightness(value hw.MMALRationalNum) error {
	return rpi.MMALPortParameterSetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_BRIGHTNESS, value)

}

// MMAL_PARAMETER_SATURATION
func (this *port) Saturation() (hw.MMALRationalNum, error) {
	return rpi.MMALPortParameterGetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SATURATION)
}

func (this *port) SetSaturation(value hw.MMALRationalNum) error {
	return rpi.MMALPortParameterSetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_SATURATION, value)

}

// MMAL_PARAMETER_ANALOG_GAIN
func (this *port) AnalogGain() (hw.MMALRationalNum, error) {
	return rpi.MMALPortParameterGetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ANALOG_GAIN)
}
func (this *port) SetAnalogGain(value hw.MMALRationalNum) error {
	return rpi.MMALPortParameterSetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ANALOG_GAIN, value)

}

// MMAL_PARAMETER_DIGITAL_GAIN
func (this *port) DigitalGain() (hw.MMALRationalNum, error) {
	return rpi.MMALPortParameterGetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_DIGITAL_GAIN)
}
func (this *port) SetDigitalGain(value hw.MMALRationalNum) error {
	return rpi.MMALPortParameterSetRational(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_DIGITAL_GAIN, value)

}

// MMAL_PARAMETER_EXP_METERING_MODE
func (this *port) MeteringMode() (hw.MMALCameraMeteringMode, error) {
	return rpi.MMALPortParameterGetCameraMeteringMode(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_EXP_METERING_MODE)
}

func (this *port) SetMeteringMode(value hw.MMALCameraMeteringMode) error {
	return rpi.MMALPortParameterSetCameraMeteringMode(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_EXP_METERING_MODE, value)
}

// MMAL_PARAMETER_EXPOSURE_MODE
func (this *port) ExposureMode() (hw.MMALCameraExposureMode, error) {
	return rpi.MMALPortParameterGetCameraExposureMode(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_EXPOSURE_MODE)
}

func (this *port) SetExposureMode(value hw.MMALCameraExposureMode) error {
	return rpi.MMALPortParameterSetCameraExposureMode(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_EXPOSURE_MODE, value)
}

// MMAL_PARAMETER_ANNOTATE
func (this *port) Annotation() (hw.MMALCameraAnnotation, error) {
	if handle, err := rpi.MMALPortParameterGetCameraAnnotation(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ANNOTATE); err != nil {
		return nil, err
	} else {
		return &annotation{handle}, nil
	}
}

func (this *port) SetAnnotation(value hw.MMALCameraAnnotation) error {
	if value_, ok := value.(*annotation); ok == false {
		return gopi.ErrBadParameter
	} else {
		return rpi.MMALPortParameterSetCameraAnnotation(this.handle, rpi.MMAL_PARAMETER_GROUP_CAMERA|rpi.MMAL_PARAMETER_ANNOTATE, value_.handle)
	}
}
