package egl_test

import (
	"testing"
	"fmt"
	"strings"

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

////////////////////////////////////////////////////////////////////////////////
// TEST API

func TestAPI_000(t *testing.T) {
	for api := egl.EGL_API_MIN; api <= egl.EGL_API_MAX; api++ {
		api_string := fmt.Sprint(api)
		if strings.HasPrefix(api_string,"EGL_API_") == true {
			t.Logf("%v => %v",api,api_string)
		} else {
			t.Errorf("Error for %v => %v",api,api_string)
		}
	}
}

func TestAPI_001(t *testing.T) {
	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGL_GetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == nil {
		t.Error("EGL_GetDisplay returned nil")
	} else if _, _, err := egl.EGL_Initialize(handle); err != nil {
		t.Error(err)
	} else {
		types := strings.Split(egl.EGL_QueryString(handle,egl.EGL_QUERY_CLIENT_APIS)," ")
		for _,api_string := range types {
			if surface_type, exists  := egl.EGL_SurfaceTypeMap[api_string]; exists == false {
				t.Error("Does not exist in EGL_SurfaceTypeMap:",api_string)
			} else if api, exists := egl.EGL_APIMap[surface_type]; exists == false {
				t.Error("Does not exist in EGL_APIMap:",api_string)
			} else if renderable, exists := egl.EGL_RenderableMap[surface_type]; exists == false {
				t.Error("Does not exist in EGL_Renderable_Map:",api_string)
			} else if err := egl.EGL_BindAPI(api); err != nil {
				t.Error("Error in EGL_BindAPI:",err)
			} else if api_, err := egl.EGL_QueryAPI(); err != nil {
				t.Error(err)
			} else if api != api_ {
				t.Error("Unexpected mismatch",api,api_)
			} else {
				t.Logf("%v => %v => %v => %v, %v",api_string,surface_type,api,api_,renderable)
			}
		}
		if err := egl.EGL_Terminate(handle); err != nil {
			t.Error(err)
		} else if err := rpi.DX_DisplayClose(display); err != nil {
			t.Error(err)
		}
	}
}
