//+build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -I/opt/vc/include
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util.h>

// Callback Functions
MMAL_BOOL_T mmal_bh_release_callback(MMAL_POOL_T* pool, MMAL_BUFFER_HEADER_T* buffer,void *userdata);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - POOLS

func MMALPortPoolCreate(handle MMAL_PortHandle, num, payload_size uint32) (MMAL_Pool, error) {
	if pool := C.mmal_port_pool_create(handle, C.uint32_t(num), C.uint32_t(payload_size)); pool == nil {
		return nil, MMAL_EINVAL
	} else {
		C.mmal_pool_callback_set(pool, C.MMAL_POOL_BH_CB_T(C.mmal_bh_release_callback), nil)
		return pool, nil
	}
}

func MMALPortPoolDestroy(handle MMAL_PortHandle, pool MMAL_Pool) error {
	C.mmal_port_pool_destroy(handle, pool)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - CALLBACKS

//export mmal_bh_release_callback
func mmal_bh_release_callback(pool *C.MMAL_POOL_T, buffer *C.MMAL_BUFFER_HEADER_T, userdata unsafe.Pointer) C.MMAL_BOOL_T {
	fmt.Printf("TODO: mmal_bh_release_callback pool=%v buffer=%v userdata=%v\n", pool, buffer, userdata)
	return MMAL_BOOL_FALSE
}
