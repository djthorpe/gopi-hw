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

////////////////////////////////////////////////////////////////////////////////
// TEST EGL QUERY

func TestQuery_001(t *testing.T) { 
	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGL_GetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == nil {
		t.Error("EGL_GetDisplay returned nil")
	} else if major, minor, err := egl.EGL_Initialize(handle); err != nil {
		t.Error(err)
	} else if vendor := egl.EGL_QueryString(handle,egl.EGL_QUERY_VENDOR); vendor == "" {
		t.Error("Empty value returned for EGL_QUERY_VENDOR")
	} else if version := egl.EGL_QueryString(handle,egl.EGL_QUERY_VERSION); version == "" {
		t.Error("Empty value returned for EGL_QUERY_VERSION")
	} else if extensions := egl.EGL_QueryString(handle,egl.EGL_QUERY_EXTENSIONS); extensions == "" {
		t.Error("Empty value returned for EGL_QUERY_EXTENSIONS")
	} else if apis := egl.EGL_QueryString(handle,egl.EGL_QUERY_CLIENT_APIS); extensions == "" {
		t.Error("Empty value returned for EGL_QUERY_CLIENT_APIS")
	} else if err := egl.EGL_Terminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DX_DisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("EGL_QUERY_VENDOR= %v", vendor)
		t.Logf("EGL_QUERY_VERSION= %v", version)
		t.Logf("EGL_QUERY_EXTENSIONS= %v", extensions)
		t.Logf("EGL_QUERY_CLIENT_APIS= %v", apis)
	}
}


////////////////////////////////////////////////////////////////////////////////
// TEST EGL CONFIGS

func TestConfigs_001(t *testing.T) { 
	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGL_GetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == nil {
		t.Error("EGL_GetDisplay returned nil")
	} else if major, minor, err := egl.EGL_Initialize(handle); err != nil {
		t.Error(err)
	} else if configs,err := egl.EGL_GetConfigs(handle); err != nil {
		t.Error(err)
		} else if len(configs) == 0  {
			t.Error("Expected at least one config")
		} else if err := egl.EGL_Terminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DX_DisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("configs= %v", configs)
	}
}

func TestConfigs_002(t *testing.T) { 
	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGL_GetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == nil {
		t.Error("EGL_GetDisplay returned nil")
	} else if major, minor, err := egl.EGL_Initialize(handle); err != nil {
		t.Error(err)
	} else if configs,err := egl.EGL_GetConfigs(handle); err != nil {
		t.Error(err)
	} else if len(configs) == 0  {
		t.Error("Expected at least one config")
	} else if attributes, err := egl.EGL_GetConfigAttribs(handle,configs[0]); err != nil {
		t.Error(err)
	} else if err := egl.EGL_Terminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DX_DisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("attributes[0]= %v", attributes)
	}
}

func TestConfigs_003(t *testing.T) { 
	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGL_GetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == nil {
		t.Error("EGL_GetDisplay returned nil")
	} else if major, minor, err := egl.EGL_Initialize(handle); err != nil {
		t.Error(err)
	} else if config, err := egl.EGL_ChooseConfig(handle,8,8,0,0); err != nil {
		t.Error(err)
	} else if attributes, err := egl.EGL_GetConfigAttribs(handle,config); err != nil {
		t.Error(err)
	} else if err := egl.EGL_Terminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DX_DisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("attributes= %v", attributes)
	}
}
