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

func TestCommon_000(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal"); mmal == nil {
		t.Fatal("Missing mmal module")
	} else {
		t.Log(mmal)
	}
}

func TestCommon_001(t *testing.T) {
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
// TEST COMMON PARAMETERS
func TestCommonSupportedEncodings_000(t *testing.T) {
	config := gopi.NewAppConfig("hw/mmal")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		t.Fatal("Invalid mmal module")
	} else if encoder, err := mmal.VideoEncoderComponent(); err != nil {
		t.Error(err)
	} else if encodings, err := encoder.Outputs()[0].SupportedEncodings(); err != nil {
		t.Error(err)
	} else if len(encodings) == 0 {
		t.Error("Expected encoder to support some encodings")
	} else {
		t.Log(encodings)
	}
}
