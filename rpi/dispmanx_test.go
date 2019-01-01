package rpi_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TEST DISPLAY OPEN

func TestDisplay_000(t *testing.T) {
	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if err := rpi.DX_DisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log(display)
	}
}

func TestDisplay_001(t *testing.T) {
	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if info, err := rpi.DX_DisplayGetInfo(display); err != nil {
		t.Error(err)
	} else if err := rpi.DX_DisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log(info)
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST RECT

func TestRect_000(t *testing.T) {
	r := rpi.DX_NewRect(0, 0, 0, 0)
	if size := rpi.DX_RectSize(r); size.W != 0 || size.H != 0 {
		t.Error("Unexpected values for rect size")
	} else if origin := rpi.DX_RectOrigin(r); origin.X != 0 || origin.Y != 0 {
		t.Error("Unexpected values for rect size")
	} else {
		t.Log("rect", rpi.DX_RectString(r))
		t.Log("size", size)
		t.Log("origin", origin)
	}
}
func TestRect_001(t *testing.T) {
	r := rpi.DX_NewRect(-100, -99, 100, 99)
	if size := rpi.DX_RectSize(r); size.W != 100 || size.H != 99 {
		t.Error("Unexpected values for rect size")
	} else if origin := rpi.DX_RectOrigin(r); origin.X != -100 || origin.Y != -99 {
		t.Error("Unexpected values for rect size")
	} else {
		t.Log("rect", rpi.DX_RectString(r))
		t.Log("size", size)
		t.Log("origin", origin)
	}
}

func TestRect_002(t *testing.T) {
	r := rpi.DX_NewRect(0, 0, 0, 0)
	if err := rpi.DX_RectSet(r, -100, -99, 100, 99); err != nil {
		t.Error(err)
	} else if size := rpi.DX_RectSize(r); size.W != 100 || size.H != 99 {
		t.Error("Unexpected values for rect size")
	} else if origin := rpi.DX_RectOrigin(r); origin.X != -100 || origin.Y != -99 {
		t.Error("Unexpected values for rect size")
	} else {
		t.Log("rect", rpi.DX_RectString(r))
		t.Log("size", size)
		t.Log("origin", origin)
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST RESOURCES

func TestResources_001(t *testing.T) {
	rpi.DX_Init()
	if resource, err := rpi.DX_ResourceCreate(rpi.DX_IMAGE_TYPE_RGBA32, rpi.DX_Size{100, 100}); err != nil {
		t.Error(err)
	} else if err := rpi.DX_ResourceDelete(resource); err != nil {
		t.Error(err)
	} else {
		t.Log(resource)
	}
}
