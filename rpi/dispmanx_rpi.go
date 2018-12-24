//+build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"fmt"
	"unsafe"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host
    #cgo LDFLAGS:  -L/opt/vc/lib -lbcm_host
	#include "vc_dispmanx.h"
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type DXDisplayId uint16
type DXDisplayHandle uint32
type DXInputFormat uint32
type DXTransform int

type DXDisplayModeInfo struct {
	Size        DXSize
	Transform   DXTransform
	InputFormat DXInputFormat
	Handle      DXDisplayHandle
}

type DXSize struct {
	Width  uint32
	Height uint32
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DX_DISPLAY_NONE DXDisplayHandle = 0
)

const (
	/* Success and failure conditions */
	DX_SUCCESS   = 0
	DX_INVALID   = -1
	DX_NO_HANDLE = 0
)

const (
	// dxTransform values
	DX_NO_ROTATE DXTransform = iota
	DX_ROTATE_90
	DX_ROTATE_180
	DX_ROTATE_270
)

const (
	// dxInputFormat values
	DX_INPUT_FORMAT_INVALID DXInputFormat = iota
	DX_INPUT_FORMAT_RGB888
	DX_INPUT_FORMAT_RGB565
)

const (
	DX_ID_MAIN_LCD DXDisplayId = iota
	DX_ID_AUX_LCD
	DX_ID_HDMI
	DX_ID_SDTV
	DX_ID_FORCE_LCD
	DX_ID_FORCE_TV
	DX_ID_FORCE_OTHER
	DX_ID_MAX = DX_ID_FORCE_OTHER
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func DXDisplayOpen(display uint) (DXDisplayHandle, error) {
	if handle := DXDisplayHandle(C.vc_dispmanx_display_open(C.uint32_t(display))); handle != DX_DISPLAY_NONE {
		return handle, nil
	} else {
		return DX_DISPLAY_NONE, gopi.ErrUnexpectedResponse
	}
}

func DXDisplayClose(display DXDisplayHandle) error {
	if C.vc_dispmanx_display_close(C.DISPMANX_DISPLAY_HANDLE_T(display)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrUnexpectedResponse
	}
}

func DXDisplayGetInfo(display DXDisplayHandle) (*DXDisplayModeInfo, error) {
	fmt.Println("get info")
	info := &DXDisplayModeInfo{}
	if C.vc_dispmanx_display_get_info(C.DISPMANX_DISPLAY_HANDLE_T(display), (*C.DISPMANX_MODEINFO_T)(unsafe.Pointer(info))) == DX_SUCCESS {
		return info, nil
	} else {
		return nil, gopi.ErrUnexpectedResponse
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *DXDisplayModeInfo) String() string {
	return fmt.Sprintf("DXDisplayModeInfo{ size=%v transform=%v input_format=%v }", this.Size, this.Transform, this.InputFormat)
}

func (size DXSize) String() string {
	return fmt.Sprintf("DXSize{%v,%v}", size.Width, size.Height)
}

func (t DXTransform) String() string {
	switch t {
	case DX_NO_ROTATE:
		return "DX_NO_ROTATE"
	case DX_ROTATE_90:
		return "DX_ROTATE_90"
	case DX_ROTATE_180:
		return "DX_ROTATE_180"
	case DX_ROTATE_270:
		return "DX_ROTATE_270"
	default:
		return "[?? Invalid DXTransform value]"
	}
}

func (f DXInputFormat) String() string {
	switch f {
	case DX_INPUT_FORMAT_RGB888:
		return "DX_INPUT_FORMAT_RGB888"
	case DX_INPUT_FORMAT_RGB565:
		return "DX_INPUT_FORMAT_RGB565"
	default:
		return "DX_INPUT_FORMAT_INVALID"
	}
}

func (d DXDisplayId) String() string {
	switch d {
	case DX_ID_MAIN_LCD:
		return "DX_ID_MAIN_LCD"
	case DX_ID_AUX_LCD:
		return "DX_ID_AUX_LCD"
	case DX_ID_HDMI:
		return "DX_ID_HDMI"
	case DX_ID_SDTV:
		return "DX_ID_SDTV"
	case DX_ID_FORCE_LCD:
		return "DX_ID_FORCE_LCD"
	case DX_ID_FORCE_TV:
		return "DX_ID_FORCE_TV"
	case DX_ID_FORCE_OTHER:
		return "DX_ID_FORCE_OTHER"
	default:
		return "[?? Invalid DXDisplayId value]"
	}
}
