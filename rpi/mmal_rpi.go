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

	"github.com/djthorpe/gopi"

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
	MMAL_CameraInfo          (*C.MMAL_PARAMETER_CAMERA_INFO_T)
	MMAL_CameraFlash         (C.MMAL_PARAMETER_CAMERA_INFO_FLASH_TYPE_T)
	MMAL_Camera              (C.MMAL_PARAMETER_CAMERA_INFO_CAMERA_T)
	MMAL_CameraAnnotation    (*C.MMAL_PARAMETER_CAMERA_ANNOTATE_V4_T)
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
		fmt.Println("ERROR")
		return nil, MMAL_Status(err)
	} else {
		fmt.Println(param)
		return param, nil
	}
}

func MMALPortParameterAllocFree(handle MMAL_ParameterHandle) {
	C.mmal_port_parameter_free(handle)
}

func MMALPortParameterSetBool(handle MMAL_PortHandle, name MMAL_ParameterType, value bool) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_boolean(handle, C.uint32_t(name), mmal_to_bool(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetBool(handle MMAL_PortHandle, name MMAL_ParameterType) (bool, error) {
	var value C.MMAL_BOOL_T
	if status := MMAL_Status(C.mmal_port_parameter_get_boolean(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return value != C.MMAL_BOOL_T(0), nil
	} else {
		return false, status
	}
}

func MMALPortParameterSetUint64(handle MMAL_PortHandle, name MMAL_ParameterType, value uint64) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_uint64(handle, C.uint32_t(name), C.uint64_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetUint64(handle MMAL_PortHandle, name MMAL_ParameterType) (uint64, error) {
	var value C.uint64_t
	if status := MMAL_Status(C.mmal_port_parameter_get_uint64(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return uint64(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetInt64(handle MMAL_PortHandle, name MMAL_ParameterType, value int64) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_int64(handle, C.uint32_t(name), C.int64_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetInt64(handle MMAL_PortHandle, name MMAL_ParameterType) (int64, error) {
	var value C.int64_t
	if status := MMAL_Status(C.mmal_port_parameter_get_int64(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return int64(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetUint32(handle MMAL_PortHandle, name MMAL_ParameterType, value uint32) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_uint32(handle, C.uint32_t(name), C.uint32_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetUint32(handle MMAL_PortHandle, name MMAL_ParameterType) (uint32, error) {
	var value C.uint32_t
	if status := MMAL_Status(C.mmal_port_parameter_get_uint32(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return uint32(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetInt32(handle MMAL_PortHandle, name MMAL_ParameterType, value int32) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_int32(handle, C.uint32_t(name), C.int32_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetInt32(handle MMAL_PortHandle, name MMAL_ParameterType) (int32, error) {
	var value C.int32_t
	if status := MMAL_Status(C.mmal_port_parameter_get_int32(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return int32(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetString(handle MMAL_PortHandle, name MMAL_ParameterType, value string) error {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	if status := MMAL_Status(C.mmal_port_parameter_set_string(handle, C.uint32_t(name), cValue)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterSetBytes(handle MMAL_PortHandle, name MMAL_ParameterType, value []byte) error {
	ptr := (*C.uint8_t)(unsafe.Pointer(&value[0]))
	len := len(value)
	if status := MMAL_Status(C.mmal_port_parameter_set_bytes(handle, C.uint32_t(name), ptr, C.uint(len))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterSetRational(handle MMAL_PortHandle, name MMAL_ParameterType, value hw.MMALRationalNum) error {
	value_ := C.MMAL_RATIONAL_T{C.int32_t(value.Num), C.int32_t(value.Den)}
	if status := MMAL_Status(C.mmal_port_parameter_set_rational(handle, C.uint32_t(name), value_)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetRational(handle MMAL_PortHandle, name MMAL_ParameterType) (hw.MMALRationalNum, error) {
	var value C.MMAL_RATIONAL_T
	if status := MMAL_Status(C.mmal_port_parameter_get_rational(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return hw.MMALRationalNum{int32(value.num), int32(value.den)}, nil
	} else {
		return hw.MMALRationalNum{}, status
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

func MMALPortParameterSetDisplayRegion(handle MMAL_PortHandle, name MMAL_ParameterType, value MMAL_DisplayRegion) error {
	if status := MMAL_Status(C.mmal_util_set_display_region(handle, value)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetDisplayRegion(handle MMAL_PortHandle, name MMAL_ParameterType) (MMAL_DisplayRegion, error) {
	var value (C.MMAL_DISPLAYREGION_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_DISPLAYREGION_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		display_region := MMAL_DisplayRegion(&value)
		display_region.set = MMAL_DISPLAY_SET_NONE
		return display_region, nil
	} else {
		return nil, status
	}
}

func MMALPortParameterGetVideoProfile(handle MMAL_PortHandle, name MMAL_ParameterType) (hw.MMALVideoProfile, error) {
	var value (C.MMAL_PARAMETER_VIDEO_PROFILE_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_VIDEO_PROFILE_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		return hw.MMALVideoProfile{hw.MMALVideoEncProfile(value.profile[0].profile), hw.MMALVideoEncLevel(value.profile[0].level)}, nil
	} else {
		return hw.MMALVideoProfile{}, status
	}
}

func MMALPortParameterSetVideoProfile(handle MMAL_PortHandle, name MMAL_ParameterType, value hw.MMALVideoProfile) error {
	var value_ (C.MMAL_PARAMETER_VIDEO_PROFILE_T)
	value_.hdr.id = C.uint32_t(name)
	value_.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_VIDEO_PROFILE_T{}))
	value_.profile[0].profile = C.MMAL_VIDEO_PROFILE_T(value.Profile)
	value_.profile[0].level = C.MMAL_VIDEO_LEVEL_T(value.Level)
	if status := MMAL_Status(C.mmal_port_parameter_set(handle, &value_.hdr)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALParamGetArrayVideoProfile(handle MMAL_ParameterHandle) []hw.MMALVideoProfile {
	fmt.Println("<TODO> SIZE=", handle.size)
	return []hw.MMALVideoProfile{}
}

func MMALPortParameterGetCameraInfo(handle MMAL_PortHandle, name MMAL_ParameterType) (MMAL_CameraInfo, error) {
	var value (C.MMAL_PARAMETER_CAMERA_INFO_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_CAMERA_INFO_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		return MMAL_CameraInfo(&value), nil
	} else {
		return nil, status
	}
}

func MMALPortParameterGetCameraMeteringMode(handle MMAL_PortHandle, name MMAL_ParameterType) (hw.MMALCameraMeteringMode, error) {
	return 0, gopi.ErrNotImplemented
}

func MMALPortParameterSetCameraMeteringMode(handle MMAL_PortHandle, name MMAL_ParameterType, value hw.MMALCameraMeteringMode) error {
	return gopi.ErrNotImplemented
}

func MMALPortParameterGetCameraExposureMode(handle MMAL_PortHandle, name MMAL_ParameterType) (hw.MMALCameraExposureMode, error) {
	var value (C.MMAL_PARAMETER_EXPOSUREMODE_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_EXPOSUREMODE_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		return hw.MMALCameraExposureMode(value.value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetCameraExposureMode(handle MMAL_PortHandle, name MMAL_ParameterType, value hw.MMALCameraExposureMode) error {
	var value_ (C.MMAL_PARAMETER_EXPOSUREMODE_T)
	value_.hdr.id = C.uint32_t(name)
	value_.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_EXPOSUREMODE_T{}))
	value_.value = C.MMAL_PARAM_EXPOSUREMODE_T(value)
	if status := MMAL_Status(C.mmal_port_parameter_set(handle, &value_.hdr)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetCameraAnnotation(handle MMAL_PortHandle, name MMAL_ParameterType) (MMAL_CameraAnnotation, error) {
	var value (C.MMAL_PARAMETER_CAMERA_ANNOTATE_V4_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_CAMERA_ANNOTATE_V4_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		return MMAL_CameraAnnotation(&value), nil
	} else {
		return nil, status
	}
}

func MMALPortParameterSetCameraAnnotation(handle MMAL_PortHandle, name MMAL_ParameterType, value MMAL_CameraAnnotation) error {
	var value_ (C.MMAL_PARAMETER_CAMERA_ANNOTATE_V4_T)
	value_.hdr.id = C.uint32_t(name)
	value_.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_CAMERA_ANNOTATE_V4_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_set(handle, &value_.hdr)); status == MMAL_SUCCESS {
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
// CAMERA ANNOTATION

func MMALCameraAnnotationEnabled(handle MMAL_CameraAnnotation) bool {
	return handle.enable == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetEnabled(handle MMAL_CameraAnnotation, value bool) {
	handle.enable = mmal_to_bool(value)
}

func MMALCameraAnnotationShowShutter(handle MMAL_CameraAnnotation) bool {
	return handle.show_shutter == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowShutter(handle MMAL_CameraAnnotation, value bool) {
	handle.show_shutter = mmal_to_bool(value)
}

func MMALCameraAnnotationShowAnalogGain(handle MMAL_CameraAnnotation) bool {
	return handle.show_analog_gain == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowAnalogGain(handle MMAL_CameraAnnotation, value bool) {
	handle.show_analog_gain = mmal_to_bool(value)
}

func MMALCameraAnnotationShowLens(handle MMAL_CameraAnnotation) bool {
	return handle.show_lens == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowLens(handle MMAL_CameraAnnotation, value bool) {
	handle.show_lens = mmal_to_bool(value)
}

func MMALCameraAnnotationShowCAF(handle MMAL_CameraAnnotation) bool {
	return handle.show_caf == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowCAF(handle MMAL_CameraAnnotation, value bool) {
	handle.show_caf = mmal_to_bool(value)
}

func MMALCameraAnnotationShowMotion(handle MMAL_CameraAnnotation) bool {
	return handle.show_motion == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowMotion(handle MMAL_CameraAnnotation, value bool) {
	handle.show_motion = mmal_to_bool(value)
}

func MMALCameraAnnotationShowFrameNum(handle MMAL_CameraAnnotation) bool {
	return handle.show_frame_num == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowFrameNum(handle MMAL_CameraAnnotation, value bool) {
	handle.show_frame_num = mmal_to_bool(value)
}

func MMALCameraAnnotationShowTextBackground(handle MMAL_CameraAnnotation) bool {
	return handle.enable_text_background == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowTextBackground(handle MMAL_CameraAnnotation, value bool) {
	handle.enable_text_background = mmal_to_bool(value)
}

func MMALCameraAnnotationUseCustomBackgroundColor(handle MMAL_CameraAnnotation) bool {
	return handle.custom_background_colour == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetUseCustomBackgroundColor(handle MMAL_CameraAnnotation, value bool) {
	handle.custom_background_colour = mmal_to_bool(value)
}

func MMALCameraAnnotationUseCustomColor(handle MMAL_CameraAnnotation) bool {
	return handle.custom_text_colour == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetUseCustomColor(handle MMAL_CameraAnnotation, value bool) {
	handle.custom_text_colour = mmal_to_bool(value)
}

func MMALCameraAnnotationText(handle MMAL_CameraAnnotation) string {
	return C.GoString(&handle.text[0])
}

func MMALCameraAnnotationSetText(handle MMAL_CameraAnnotation, value string) {
	cstr := C.CString(value)
	defer C.free(unsafe.Pointer(cstr))
	C.strncpy(&handle.text[0], cstr, C.uint(len(value)))
}

func MMALCameraAnnotationTextSize(handle MMAL_CameraAnnotation) uint8 {
	return uint8(handle.text_size)
}

func MMALCameraAnnotationSetTextSize(handle MMAL_CameraAnnotation, value uint8) {
	handle.text_size = C.uint8_t(value)
}

////////////////////////////////////////////////////////////////////////////////
// CAMERA INFO

func MMALCameraInfoGetCamerasNum(handle MMAL_CameraInfo) uint32 {
	return uint32(handle.num_cameras)
}

func MMALCameraInfoGetFlashesNum(handle MMAL_CameraInfo) uint32 {
	return uint32(handle.num_flashes)
}

func MMALCameraInfoGetCameras(handle MMAL_CameraInfo) []MMAL_Camera {
	cameras := make([]MMAL_Camera, int(handle.num_cameras))
	for i := 0; i < len(cameras); i++ {
		cameras[i] = MMAL_Camera(handle.cameras[i])
	}
	return cameras
}

func MMALCameraInfoGetFlashes(handle MMAL_CameraInfo) []hw.MMALCameraFlashType {
	flashes := make([]hw.MMALCameraFlashType, int(handle.num_flashes))
	for i := 0; i < len(flashes); i++ {
		flashes[i] = hw.MMALCameraFlashType(handle.flashes[i].flash_type)
	}
	return flashes
}

func MMALCameraInfoGetCameraId(handle MMAL_Camera) uint32 {
	return uint32(handle.port_id)
}

func MMALCameraInfoGetCameraName(handle MMAL_Camera) string {
	return C.GoString(&handle.camera_name[0])
}

func MMALCameraInfoGetCameraMaxWidth(handle MMAL_Camera) uint32 {
	return uint32(handle.max_width)
}
func MMALCameraInfoGetCameraMaxHeight(handle MMAL_Camera) uint32 {
	return uint32(handle.max_height)
}

func MMALCameraInfoGetCameraLensPresent(handle MMAL_Camera) bool {
	return handle.lens_present == MMAL_BOOL_TRUE
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
