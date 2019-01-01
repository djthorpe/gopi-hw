package egl_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi-hw/egl"
	"github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TEST EGL INIT

func TestEGL_000(t *testing.T) {
	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if err := rpi.DX_DisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log(display)
	}
}
func TestEGL_001(t *testing.T) {
	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if major, minor, err := egl.EGL_Initialize(egl.EGL_GetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD))); err != nil {
		t.Error(err)
	} else if err := egl.EGL_Terminate(egl.EGL_GetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD))); err != nil {
		t.Error(err)
	} else if err := rpi.DX_DisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
	}
}
