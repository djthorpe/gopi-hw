/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package openvg

// Frameworks

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: vg egl
#include <VG/openvg.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	VG_Error uint
	VG_Int   C.VGint
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_ERROR_NONE                     VG_Error = 0
	VG_ERROR_BAD_HANDLE               VG_Error = 0x1000
	VG_ERROR_ILLEGAL_ARGUMENT         VG_Error = 0x1001
	VG_ERROR_OUT_OF_MEMORY            VG_Error = 0x1002
	VG_ERROR_PATH_CAPABILITY          VG_Error = 0x1003
	VG_ERROR_UNSUPPORTED_IMAGE_FORMAT VG_Error = 0x1004
	VG_ERROR_UNSUPPORTED_PATH_FORMAT  VG_Error = 0x1005
	VG_ERROR_IMAGE_IN_USE             VG_Error = 0x1006
	VG_ERROR_NO_CONTEXT               VG_Error = 0x1007
	VG_ERROR_MIN                               = VG_ERROR_BAD_HANDLE
	VG_ERROR_MAX                               = VG_ERROR_NO_CONTEXT
)

////////////////////////////////////////////////////////////////////////////////
// ERRORS

func VG_GetError() error {
	if ret := VG_Error(C.vgGetError()); ret == VG_ERROR_NONE {
		return nil
	} else {
		return ret
	}
}

////////////////////////////////////////////////////////////////////////////////
// FLUSH/FINISH

func VG_Flush() error {
	C.vgFlush()
	return VG_GetError()
}

func VG_Finish() error {
	C.vgFinish()
	return VG_GetError()
}

////////////////////////////////////////////////////////////////////////////////
// CLEAR

func VG_Clear(x, y, width, height VG_Int) error {
	C.vgClear(C.VGint(x), C.VGint(y), C.VGint(width), C.VGint(height))
	return VG_GetError()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e VG_Error) Error() string {
	switch e {
	case VG_ERROR_NONE:
		return "VG_ERROR_NONE"
	case VG_ERROR_BAD_HANDLE:
		return "VG_ERROR_BAD_HANDLE"
	case VG_ERROR_ILLEGAL_ARGUMENT:
		return "VG_ERROR_ILLEGAL_ARGUMENT"
	case VG_ERROR_OUT_OF_MEMORY:
		return "VG_ERROR_OUT_OF_MEMORY"
	case VG_ERROR_PATH_CAPABILITY:
		return "VG_ERROR_PATH_CAPABILITY"
	case VG_ERROR_UNSUPPORTED_IMAGE_FORMAT:
		return "VG_ERROR_UNSUPPORTED_IMAGE_FORMAT"
	case VG_ERROR_UNSUPPORTED_PATH_FORMAT:
		return "VG_ERROR_UNSUPPORTED_PATH_FORMAT"
	case VG_ERROR_IMAGE_IN_USE:
		return "VG_ERROR_IMAGE_IN_USE"
	case VG_ERROR_NO_CONTEXT:
		return "VG_ERROR_NO_CONTEXT"
	default:
		return "[?? Invalid VG_Error value]"
	}
}
