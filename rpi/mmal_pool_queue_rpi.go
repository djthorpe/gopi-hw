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
		fmt.Println("Setting pool callback", pool, C.mmal_bh_release_callback)
		C.mmal_pool_callback_set(pool, C.MMAL_POOL_BH_CB_T(C.mmal_bh_release_callback), nil)
		return pool, nil
	}
}

func MMALPortPoolDestroy(handle MMAL_PortHandle, pool MMAL_Pool) error {
	C.mmal_port_pool_destroy(handle, pool)
	return nil
}

func MMALPoolGetBuffer(pool MMAL_Pool) MMAL_Buffer {
	fmt.Println("MMALPoolGetBuffer", pool, pool.queue)
	return MMAL_Buffer(C.mmal_queue_get(pool.queue))
}

func MMALPoolResize(handle MMAL_Pool, num, payload_size uint32) error {
	if status := MMAL_Status(C.mmal_pool_resize(handle, C.uint32_t(num), C.uint32_t(payload_size))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPoolString(pool MMAL_Pool) string {
	if pool == nil {
		return "<MMAL_Pool>{ nil }"
	} else {
		buffers := make([]string, pool.headers_num)
		for i := range buffers {
			buffers[i] = "<TODO>"
		}
		return fmt.Sprintf("<MMAL_Pool>{ queue=%v buffers=%v }", MMALQueueString(pool.queue), buffers)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - QUEUES

func MMALQueueCreate() MMAL_Queue {
	return MMAL_Queue(C.mmal_queue_create())
}

func MMALQueueDestroy(handle MMAL_Queue) {
	C.mmal_queue_destroy(handle)
}

func MMALQueueString(handle MMAL_Queue) string {
	if handle == nil {
		return "<MMAL_Queue>{ nil }"
	} else {
		return fmt.Sprintf("<MMAL_Queue>{ %v }", handle)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - CALLBACKS

//export mmal_bh_release_callback
func mmal_bh_release_callback(pool *C.MMAL_POOL_T, buffer *C.MMAL_BUFFER_HEADER_T, userdata unsafe.Pointer) C.MMAL_BOOL_T {
	fmt.Printf("TODO: mmal_bh_release_callback pool=%v buffer=%v userdata=%v\n", pool, buffer, userdata)
	return MMAL_BOOL_FALSE
}
