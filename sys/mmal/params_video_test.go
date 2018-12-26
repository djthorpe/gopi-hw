package mmal_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/metrics"
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
	} else if _, err := mmal.CameraComponent(); err != nil {
		t.Error(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST PROFILES
func TestVideoProfile_000(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		t.Fatal("Invalid mmal module")
	} else if encoder, err := mmal.VideoEncoderComponent(); err != nil {
		t.Error(err)
	} else if profile, err := encoder.Output()[0].VideoProfile(); err != nil {
		t.Error(err)
	} else {
		t.Log(profile)
	}
}

func TestVideoProfile_001(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		t.Fatal("Invalid mmal module")
	} else if encoder, err := mmal.VideoEncoderComponent(); err != nil {
		t.Error(err)
	} else if profile, err := encoder.Output()[0].VideoProfile(); err != nil {
		t.Error(err)
	} else if err := encoder.Output()[0].SetVideoProfile(profile); err != nil {
		t.Error(err)
	} else {
		t.Log(profile)
	}
}

func TestVideoProfile_002(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		t.Fatal("Invalid mmal module")
	} else if encoder, err := mmal.CameraComponent(); err != nil {
		t.Error(err)
	} else if profiles, err := encoder.Output()[1].SupportedVideoProfiles(); err != nil {
		t.Error(err)
	} else {
		t.Log(profiles)
	}
}
