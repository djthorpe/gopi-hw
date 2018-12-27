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
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -I/opt/vc/include
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util.h>
#include <interface/mmal/util/mmal_util_params.h>

// Callback Functions
void mmal_port_callback(MMAL_PORT_T* port, MMAL_BUFFER_HEADER_T* buffer);
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - PORTS

func MMALPortEnable(handle MMAL_PortHandle) error {
	if status := MMAL_Status(C.mmal_port_enable(handle, C.MMAL_PORT_BH_CB_T(C.mmal_port_callback))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortDisable(handle MMAL_PortHandle) error {
	if status := MMAL_Status(C.mmal_port_disable(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortFlush(handle MMAL_PortHandle) error {
	if status := MMAL_Status(C.mmal_port_flush(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortName(handle MMAL_PortHandle) string {
	return C.GoString(handle.name)
}

func MMALPortType(handle MMAL_PortHandle) MMAL_PortType {
	return MMAL_PortType(handle._type)
}

func MMALPortIndex(handle MMAL_PortHandle) uint {
	return uint(handle.index)
}

func MMALPortIsEnabled(handle MMAL_PortHandle) bool {
	return (handle.is_enabled != 0)
}

func MMALPortCapabilities(handle MMAL_PortHandle) MMAL_PortCapability {
	return MMAL_PortCapability(handle.capabilities)
}

func MMALPortDisconnect(handle MMAL_PortHandle) error {
	if status := MMAL_Status(C.mmal_port_disconnect(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnect(this, other MMAL_PortHandle) error {
	if status := MMAL_Status(C.mmal_port_connect(this, other)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortFormatCommit(handle MMAL_PortHandle) error {
	if status := MMAL_Status(C.mmal_port_format_commit(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortFormat(handle MMAL_PortHandle) MMAL_StreamFormat {
	return handle.format
}

func MMALPortComponent(handle MMAL_PortHandle) MMAL_ComponentHandle {
	return handle.component
}

func MMALPortBufferNum(handle MMAL_PortHandle) (uint32, uint32) {
	// Minimum & recommended number of buffers the port requires
	// A value of zero for recommendation means no special recommendation
	return uint32(handle.buffer_num_min), uint32(handle.buffer_num_recommended)
}

func MMALPortBufferSize(handle MMAL_PortHandle) (uint32, uint32) {
	// Minimum & recommended size of buffers the port requires
	// A value of zero means no special recommendation
	return uint32(handle.buffer_size_min), uint32(handle.buffer_size_recommended)
}

func MMALPortBufferAlignment(handle MMAL_PortHandle) uint32 {
	// Minimum alignment requirement for the buffers. A value of zero
	// means no special alignment requirements.
	return uint32(handle.buffer_alignment_min)
}

func MMALPortSetURI(handle MMAL_PortHandle, value string) error {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	if status := MMAL_Status(C.mmal_util_port_set_uri(handle, cValue)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortSetDisplayRegion(handle MMAL_PortHandle, value MMAL_DisplayRegion) error {
	if status := MMAL_Status(C.mmal_util_set_display_region(handle, value)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - CALLBACKS

//export mmal_port_callback
func mmal_port_callback(port *C.MMAL_PORT_T, buffer *C.MMAL_BUFFER_HEADER_T) {
	fmt.Printf("TODO: mmal_port_callback port=%v buffer=%v\n", port, buffer)
}