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
	// Frameworks
	"fmt"
	"reflect"
	"unsafe"

	hw "github.com/djthorpe/gopi-hw"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -I/opt/vc/include
#include <interface/mmal/mmal.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func MMALBufferCommand(handle MMAL_Buffer) hw.MMALEncodingType {
	return hw.MMALEncodingType(handle.cmd)
}

func MMALBufferBytes(handle MMAL_Buffer) []byte {
	var value []byte
	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
	sliceHeader.Cap = int(handle.alloc_size)
	sliceHeader.Len = int(handle.length)
	sliceHeader.Data = uintptr(unsafe.Pointer(handle.data))
	// Return data
	return value
}

func MMALBufferFlags(handle MMAL_Buffer) hw.MMALBufferFlag {
	return hw.MMALBufferFlag(handle.flags)
}

func MMALBufferLength(handle MMAL_Buffer) uint32 {
	return uint32(handle.length)
}

func MMALBufferOffset(handle MMAL_Buffer) uint32 {
	return uint32(handle.offset)
}

func MMALBufferString(handle MMAL_Buffer) string {
	if handle == nil {
		return fmt.Sprintf("<MMAL_Buffer>{ nil }")
	} else {
		return fmt.Sprintf("<MMAL_Buffer>{ cmd=%v length=%v offset=%v flags=%v }", hw.MMALEncodingType(handle.cmd), uint32(handle.length), uint32(handle.offset), hw.MMALBufferFlag(handle.flags))
	}
}
