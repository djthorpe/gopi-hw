package mmal_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/mmal"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////
// TEST MMAL

func TestMMAL_000(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal"); mmal == nil {
		t.Fatal("Missing mmal module")
	} else {
		t.Log(mmal)
	}
}

func TestMMAL_001(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		t.Fatal("Invalid mmal module")
	} else if _, err := mmal.CameraInfoComponent(); err != nil {
		t.Error(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST CAMERA INFO PARAMETERS

func TestCameraInfo_000(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		t.Fatal("Invalid mmal module")
	} else if camera_info, err := mmal.CameraInfoComponent(); err != nil {
		t.Error(err)
	} else {
		t.Log(camera_info)
	}
}

func TestCameraInfo_001(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		t.Fatal("Invalid mmal module")
	} else if camera_info, err := mmal.CameraInfoComponent(); err != nil {
		t.Error(err)
	} else if value, err := camera_info.Control().CameraInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(value)
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST RATIONAL PARAMETERS

func TestCamera_001(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		t.Fatal("Invalid mmal module")
	} else if camera, err := mmal.CameraComponent(); err != nil {
		t.Error(err)
	} else if brightness, err := camera.Control().Brightness(); err != nil {
		t.Error(err)
	} else if sharpness, err := camera.Control().Sharpness(); err != nil {
		t.Error(err)
	} else if contrast, err := camera.Control().Contrast(); err != nil {
		t.Error(err)
	} else if saturation, err := camera.Control().Saturation(); err != nil {
		t.Error(err)
	} else if analog_gain, err := camera.Control().AnalogGain(); err != nil {
		t.Error(err)
	} else if digital_gain, err := camera.Control().DigitalGain(); err != nil {
		t.Error(err)
	} else {
		t.Log("brightness", brightness)
		t.Log("sharpness", sharpness)
		t.Log("contrast", contrast)
		t.Log("saturation", saturation)
		t.Log("analog_gain", analog_gain)
		t.Log("digital_gain", digital_gain)
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST CAMERA ANNOTAION
func TestCamera_002(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		t.Fatal("Invalid mmal module")
	} else if camera, err := mmal.CameraComponent(); err != nil {
		t.Error(err)
	} else if annotation, err := camera.Control().Annotation(); err != nil {
		t.Error(err)
	} else {
		t.Log(annotation)
	}
}
