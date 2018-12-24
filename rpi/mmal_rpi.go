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
	"reflect"
	"unsafe"

	// Frameworks
	hw "github.com/djthorpe/gopi-hw"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -I/opt/vc/include
#cgo LDFLAGS: -L/opt/vc/lib -lmmal -lmmal_core -lmmal_util
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util.h>
#include <interface/mmal/util/mmal_util_params.h>
#include <interface/mmal/util/mmal_connection.h>

// Callback Functions
void mmal_port_callback(MMAL_PORT_T* port, MMAL_BUFFER_HEADER_T* buffer);
MMAL_BOOL_T mmal_bh_release_callback(MMAL_POOL_T* pool, MMAL_BUFFER_HEADER_T* buffer,void *userdata);
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMAL_Status              (C.MMAL_STATUS_T)
	MMAL_ComponentHandle     (*C.MMAL_COMPONENT_T)
	MMAL_PortHandle          (*C.MMAL_PORT_T)
	MMAL_PortConnection      (*C.MMAL_CONNECTION_T)
	MMAL_DisplayRegion       (*C.MMAL_DISPLAYREGION_T)
	MMAL_PortType            (C.MMAL_PORT_TYPE_T)
	MMAL_PortCapability      (C.uint32_t)
	MMAL_Rational            (C.MMAL_RATIONAL_T)
	MMAL_StreamType          (C.MMAL_ES_TYPE_T)
	MMAL_StreamFormat        (*C.MMAL_ES_FORMAT_T)
	MMAL_StreamCompareFlags  (C.uint32_t)
	MMAL_PortConnectionFlags (C.uint32_t)
	MMAL_Buffer              (*C.MMAL_BUFFER_HEADER_T)
	MMAL_Pool                (*C.MMAL_POOL_T)
	MMAL_Queue               (*C.MMAL_QUEUE_T)
	MMAL_ParameterHandle     (*C.MMAL_PARAMETER_HEADER_T)
	MMAL_ParameterType       uint
	MMAL_ParameterSeek       (C.MMAL_PARAMETER_SEEK_T)
)

type MMAL_Rect struct {
	X, Y int32
	W, H uint32
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - COMPONENTS

func MMALComponentCreate(name string, handle *MMAL_ComponentHandle) error {
	var cHandle (*C.MMAL_COMPONENT_T)

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	if status := MMAL_Status(C.mmal_component_create(cName, &cHandle)); status == MMAL_SUCCESS {
		*handle = MMAL_ComponentHandle(cHandle)
		return nil
	} else {
		return status
	}
}

func MMALComponentDestroy(handle MMAL_ComponentHandle) error {
	if status := MMAL_Status(C.mmal_component_destroy(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALComponentAcquire(handle MMAL_ComponentHandle) error {
	C.mmal_component_acquire(handle)
	return nil
}

func MMALComponentRelease(handle MMAL_ComponentHandle) error {
	if status := MMAL_Status(C.mmal_component_release(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALComponentEnable(handle MMAL_ComponentHandle) error {
	if status := MMAL_Status(C.mmal_component_enable(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALComponentDisable(handle MMAL_ComponentHandle) error {
	if status := MMAL_Status(C.mmal_component_disable(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALComponentName(handle MMAL_ComponentHandle) string {
	return C.GoString(handle.name)
}

func MMALComponentId(handle MMAL_ComponentHandle) uint32 {
	return uint32(handle.id)
}

func MMALComponentIsEnabled(handle MMAL_ComponentHandle) bool {
	return (handle.is_enabled != 0)
}

func MMALComponentControlPort(handle MMAL_ComponentHandle) MMAL_PortHandle {
	return handle.control
}

func MMALComponentInputPortNum(handle MMAL_ComponentHandle) uint {
	return uint(handle.input_num)
}

func MMALComponentInputPortAtIndex(handle MMAL_ComponentHandle, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.input, uint(handle.input_num), index)
}

func MMALComponentOutputPortNum(handle MMAL_ComponentHandle) uint {
	return uint(handle.output_num)
}

func MMALComponentOutputPortAtIndex(handle MMAL_ComponentHandle, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.output, uint(handle.output_num), index)
}

func MMALComponentClockPortNum(handle MMAL_ComponentHandle) uint {
	return uint(handle.clock_num)
}

func MMALComponentClockPortAtIndex(handle MMAL_ComponentHandle, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.clock, uint(handle.clock_num), index)
}

func MMALComponentPortNum(handle MMAL_ComponentHandle) uint {
	return uint(handle.port_num)
}

func MMALComponentPortAtIndex(handle MMAL_ComponentHandle, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.port, uint(handle.port_num), index)
}

func mmal_component_port_at_index(array **C.MMAL_PORT_T, num, index uint) MMAL_PortHandle {
	var handles []MMAL_PortHandle
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&handles)))
	sliceHeader.Cap = int(num)
	sliceHeader.Len = int(num)
	sliceHeader.Data = uintptr(unsafe.Pointer(array))
	return handles[index]
}

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
// PUBLIC METHODS - QUEUES
/*
func MMALPortQueueCreate(handle MMAL_PortHandle) (MMAL_Pool, error) {
	if pool := C.mmal_port_pool_create(handle, C.uint32_t(num), C.uint32_t(payload_size)); pool == nil {
		return nil, MMAL_EINVAL
	} else {
		C.mmal_pool_callback_set(pool, C.mmal_bh_release_callback, nil)
		return pool, nil
	}
}

func MMALPortQueueDestroy(handle MMAL_PortHandle, pool MMAL_Pool) error {
	C.mmal_port_pool_destroy(handle, pool)
	return nil
}
*/

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - CONNECTIONS

func MMALPortConnectionCreate(handle *MMAL_PortConnection, output_port, input_port MMAL_PortHandle, flags hw.MMALPortConnectionFlags) error {
	var cHandle (*C.MMAL_CONNECTION_T)
	if status := MMAL_Status(C.mmal_connection_create(&cHandle, output_port, input_port, C.uint(flags))); status == MMAL_SUCCESS {
		*handle = MMAL_PortConnection(cHandle)
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionAcquire(handle MMAL_PortConnection) error {
	C.mmal_connection_acquire(handle)
	return nil
}

func MMALPortConnectionRelease(handle MMAL_PortConnection) error {
	if status := MMAL_Status(C.mmal_connection_release(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionDestroy(handle MMAL_PortConnection) error {
	if status := MMAL_Status(C.mmal_connection_destroy(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionEnable(handle MMAL_PortConnection) error {
	if status := MMAL_Status(C.mmal_connection_enable(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionDisable(handle MMAL_PortConnection) error {
	if status := MMAL_Status(C.mmal_connection_disable(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionEventFormatChanged(handle MMAL_PortConnection, buffer MMAL_Buffer) error {
	if status := MMAL_Status(C.mmal_connection_event_format_changed(handle, buffer)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionEnabled(handle MMAL_PortConnection) bool {
	return (handle.is_enabled != 0)
}

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
// PUBLIC METHODS - PARAMETERS

func MMALPortParameterAllocGet(handle MMAL_PortHandle, name MMAL_ParameterType, size uint32) (MMAL_ParameterHandle, error) {
	var err C.MMAL_STATUS_T
	if param := C.mmal_port_parameter_alloc_get(handle, C.uint32_t(name), C.uint32_t(size), &err); MMAL_Status(err) != MMAL_SUCCESS {
		return nil, MMAL_Status(err)
	} else {
		return param, nil
	}
}

func MMALPortParameterAllocFree(handle MMAL_ParameterHandle) {
	C.mmal_port_parameter_free(handle)
}

func MMALPortParameterSetDisplayRegion(handle MMAL_PortHandle, name MMAL_ParameterType, value MMAL_DisplayRegion) error {
	if status := MMAL_Status(C.mmal_util_set_display_region(handle, value)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetDisplayRegion(handle MMAL_PortHandle, name MMAL_ParameterType) (MMAL_DisplayRegion, error) {
	var value (C.MMAL_DISPLAYREGION_T)
	value.hdr.id = C.uint(name)
	value.hdr.size = C.uint(unsafe.Sizeof(C.MMAL_DISPLAYREGION_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		display_region := MMAL_DisplayRegion(&value)
		display_region.set = MMAL_DISPLAY_SET_NONE
		return display_region, nil
	} else {
		return nil, status
	}
}

func MMALPortParameterSetBool(handle MMAL_PortHandle, name MMAL_ParameterType, value bool) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_boolean(handle, C.uint(name), mmal_to_bool(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetBool(handle MMAL_PortHandle, name MMAL_ParameterType) (bool, error) {
	var value C.MMAL_BOOL_T
	if status := MMAL_Status(C.mmal_port_parameter_get_boolean(handle, C.uint(name), &value)); status == MMAL_SUCCESS {
		return value != C.MMAL_BOOL_T(0), nil
	} else {
		return false, status
	}
}

func MMALPortParameterSetUint64(handle MMAL_PortHandle, name MMAL_ParameterType, value uint64) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_uint64(handle, C.uint(name), C.uint64_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetUint64(handle MMAL_PortHandle, name MMAL_ParameterType) (uint64, error) {
	var value C.uint64_t
	if status := MMAL_Status(C.mmal_port_parameter_get_uint64(handle, C.uint(name), &value)); status == MMAL_SUCCESS {
		return uint64(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetInt64(handle MMAL_PortHandle, name MMAL_ParameterType, value int64) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_int64(handle, C.uint(name), C.int64_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetInt64(handle MMAL_PortHandle, name MMAL_ParameterType) (int64, error) {
	var value C.int64_t
	if status := MMAL_Status(C.mmal_port_parameter_get_int64(handle, C.uint(name), &value)); status == MMAL_SUCCESS {
		return int64(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetUint32(handle MMAL_PortHandle, name MMAL_ParameterType, value uint32) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_uint32(handle, C.uint(name), C.uint32_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetUint32(handle MMAL_PortHandle, name MMAL_ParameterType) (uint32, error) {
	var value C.uint32_t
	if status := MMAL_Status(C.mmal_port_parameter_get_uint32(handle, C.uint(name), &value)); status == MMAL_SUCCESS {
		return uint32(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetInt32(handle MMAL_PortHandle, name MMAL_ParameterType, value int32) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_int32(handle, C.uint(name), C.int32_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetInt32(handle MMAL_PortHandle, name MMAL_ParameterType) (int32, error) {
	var value C.int32_t
	if status := MMAL_Status(C.mmal_port_parameter_get_int32(handle, C.uint(name), &value)); status == MMAL_SUCCESS {
		return int32(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetString(handle MMAL_PortHandle, name MMAL_ParameterType, value string) error {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	if status := MMAL_Status(C.mmal_port_parameter_set_string(handle, C.uint(name), cValue)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterSetBytes(handle MMAL_PortHandle, name MMAL_ParameterType, value []byte) error {
	ptr := (*C.uint8_t)(unsafe.Pointer(&value[0]))
	len := len(value)
	if status := MMAL_Status(C.mmal_port_parameter_set_bytes(handle, C.uint(name), ptr, C.uint(len))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterSetRational(handle MMAL_PortHandle, name MMAL_ParameterType, value MMAL_Rational) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_rational(handle, C.uint(name), C.MMAL_RATIONAL_T(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetRational(handle MMAL_PortHandle, name MMAL_ParameterType) (MMAL_Rational, error) {
	var value C.MMAL_RATIONAL_T
	if status := MMAL_Status(C.mmal_port_parameter_get_rational(handle, C.uint(name), &value)); status == MMAL_SUCCESS {
		return MMAL_Rational(value), nil
	} else {
		return MMAL_Rational(value), status
	}
}

func MMALPortParameterSetSeek(handle MMAL_PortHandle, name MMAL_ParameterType, value MMAL_ParameterSeek) error {
	value.hdr.id = C.uint32_t(name)
	if status := MMAL_Status(C.mmal_port_parameter_set(handle, (*C.MMAL_PARAMETER_HEADER_T)(&value.hdr))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALParamGetArrayUint32(handle MMAL_ParameterHandle) []uint32 {
	var array []uint32

	// Data and length of the array
	data := uintptr(unsafe.Pointer(handle)) + unsafe.Sizeof(*handle)
	len := (uintptr(handle.size) - unsafe.Sizeof(*handle)) / C.sizeof_uint32_t

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&array)))
	sliceHeader.Cap = int(len)
	sliceHeader.Len = int(len)
	sliceHeader.Data = data

	// Return the array
	return array
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - STREAM FORMATS

func MMALStreamFormatAlloc() MMAL_StreamFormat {
	return MMAL_StreamFormat(C.mmal_format_alloc())
}

func MMALStreamFormatFree(handle MMAL_StreamFormat) {
	C.mmal_format_free(handle)
}

func MMALStreamFormatExtraDataAlloc(handle MMAL_StreamFormat, size uint) error {
	if status := MMAL_Status(C.mmal_format_extradata_alloc(handle, C.uint(size))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALStreamFormatCopy(dest, src MMAL_StreamFormat) error {
	C.mmal_format_copy(dest, src)
	return nil
}

func MMALStreamFormatFullCopy(dest, src MMAL_StreamFormat) error {
	if status := MMAL_Status(C.mmal_format_full_copy(dest, src)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALStreamFormatCompare(dest, src MMAL_StreamFormat) MMAL_StreamCompareFlags {
	return MMAL_StreamCompareFlags(C.mmal_format_compare(dest, src))
}

func MMALStreamFormatType(handle MMAL_StreamFormat) MMAL_StreamType {
	return MMAL_StreamType(handle._type)
}

////////////////////////////////////////////////////////////////////////////////
// DISPLAY REGION

func MMALDisplayRegionGetDisplayNum(handle MMAL_DisplayRegion) uint32 {
	return uint32(handle.display_num)
}

func MMALDisplayRegionGetFullScreen(handle MMAL_DisplayRegion) bool {
	return handle.fullscreen != 0
}

func MMALDisplayRegionSetFullScreen(handle MMAL_DisplayRegion, value bool) {
	handle.fullscreen = mmal_to_bool(value)
	handle.set |= MMAL_DISPLAY_SET_FULLSCREEN
}

func MMALDisplayRegionGetLayer(handle MMAL_DisplayRegion) int32 {
	return int32(handle.layer)
}

func MMALDisplayRegionGetAlpha(handle MMAL_DisplayRegion) uint32 {
	return uint32(handle.alpha)
}

func MMALDisplayRegionSetLayer(handle MMAL_DisplayRegion, value int32) {
	handle.layer = C.int32_t(value)
	handle.set |= MMAL_DISPLAY_SET_LAYER
}

func MMALDisplayRegionSetAlpha(handle MMAL_DisplayRegion, value uint32) {
	handle.alpha = C.uint32_t(value)
	handle.set |= MMAL_DISPLAY_SET_ALPHA
}

func MMALDisplayRegionGetTransform(handle MMAL_DisplayRegion) hw.MMALDisplayTransform {
	return hw.MMALDisplayTransform(handle.transform)
}

func MMALDisplayRegionGetMode(handle MMAL_DisplayRegion) hw.MMALDisplayMode {
	return hw.MMALDisplayMode(handle.mode)
}

func MMALDisplayRegionSetTransform(handle MMAL_DisplayRegion, value hw.MMALDisplayTransform) {
	handle.transform = C.MMAL_DISPLAYTRANSFORM_T(value)
	handle.set |= MMAL_DISPLAY_SET_TRANSFORM
}

func MMALDisplayRegionSetMode(handle MMAL_DisplayRegion, value hw.MMALDisplayMode) {
	handle.mode = C.MMAL_DISPLAYMODE_T(value)
	handle.set |= MMAL_DISPLAY_SET_MODE
}

func MMALDisplayRegionGetNoAspect(handle MMAL_DisplayRegion) bool {
	return handle.noaspect != 0
}

func MMALDisplayRegionGetCopyProtect(handle MMAL_DisplayRegion) bool {
	return handle.copyprotect_required != 0
}

func MMALDisplayRegionSetNoAspect(handle MMAL_DisplayRegion, value bool) {
	handle.noaspect = mmal_to_bool(value)
	handle.set |= MMAL_DISPLAY_SET_NOASPECT
}

func MMALDisplayRegionSetCopyProtect(handle MMAL_DisplayRegion, value bool) {
	handle.copyprotect_required = mmal_to_bool(value)
	handle.set |= MMAL_DISPLAY_SET_COPYPROTECT
}

func MMALDisplayRegionGetDestRect(handle MMAL_DisplayRegion) MMAL_Rect {
	return MMAL_Rect{int32(handle.dest_rect.x), int32(handle.dest_rect.y), uint32(handle.dest_rect.width), uint32(handle.dest_rect.height)}
}

func MMALDisplayRegionGetSrcRect(handle MMAL_DisplayRegion) MMAL_Rect {
	return MMAL_Rect{int32(handle.src_rect.x), int32(handle.src_rect.y), uint32(handle.src_rect.width), uint32(handle.src_rect.height)}
}

func MMALDisplayRegionSetDestRect(handle MMAL_DisplayRegion, value MMAL_Rect) {
	handle.dest_rect = C.MMAL_RECT_T{C.int32_t(value.X), C.int32_t(value.Y), C.int32_t(value.W), C.int32_t(value.H)}
	handle.set |= MMAL_DISPLAY_SET_DEST_RECT
}

func MMALDisplayRegionSetSrcRect(handle MMAL_DisplayRegion, value MMAL_Rect) {
	handle.src_rect = C.MMAL_RECT_T{C.int32_t(value.X), C.int32_t(value.Y), C.int32_t(value.W), C.int32_t(value.H)}
	handle.set |= MMAL_DISPLAY_SET_SRC_RECT
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - SEEK

func (this *MMAL_ParameterSeek) SetOffset(value int64) {
	this.offset = C.int64_t(value)
}

func (this *MMAL_ParameterSeek) SetFlags(value uint32) {
	this.flags = C.uint32_t(value)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - CALLBACKS

//export mmal_port_callback
func mmal_port_callback(port *C.MMAL_PORT_T, buffer *C.MMAL_BUFFER_HEADER_T) {
	fmt.Printf("TODO: mmal_port_callback port=%v buffer=%v\n", port, buffer)
}

//export mmal_bh_release_callback
func mmal_bh_release_callback(pool *C.MMAL_POOL_T, buffer *C.MMAL_BUFFER_HEADER_T, userdata unsafe.Pointer) C.MMAL_BOOL_T {
	fmt.Printf("TODO: mmal_bh_release_callback pool=%v buffer=%v userdata=%v\n", pool, buffer, userdata)
	return MMAL_BOOL_FALSE
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - OTHER

func mmal_to_bool(value bool) C.MMAL_BOOL_T {
	if value {
		return MMAL_BOOL_TRUE
	} else {
		return MMAL_BOOL_FALSE
	}
}
