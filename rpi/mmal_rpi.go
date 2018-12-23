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
#cgo LDFLAGS: -L/opt/vc/lib -lmmal -lmmal_core -lmmal_util
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util.h>
#include <interface/mmal/util/mmal_util_params.h>

void mmal_port_callback(MMAL_PORT_T* port, MMAL_BUFFER_HEADER_T* buffer);
*/
import "C"

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMAL_Status             (C.MMAL_STATUS_T)
	MMAL_ComponentHandle    (*C.MMAL_COMPONENT_T)
	MMAL_PortType           (C.MMAL_PORT_TYPE_T)
	MMAL_PortHandle         (*C.MMAL_PORT_T)
	MMAL_PortCapability     (C.uint32_t)
	MMAL_ParameterHandle    (*C.MMAL_PARAMETER_HEADER_T)
	MMAL_ParameterType      uint
	MMAL_DisplayRegion      (*C.MMAL_DISPLAYREGION_T)
	MMAL_Rational           (C.MMAL_RATIONAL_T)
	MMAL_StreamType         (C.MMAL_ES_TYPE_T)
	MMAL_StreamFormat       (*C.MMAL_ES_FORMAT_T)
	MMAL_StreamCompareFlags (C.uint32_t)
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_BOOL_FALSE = C.MMAL_BOOL_T(0)
	MMAL_BOOL_TRUE  = C.MMAL_BOOL_T(1)
)

const (
	MMAL_SUCCESS   MMAL_Status = iota
	MMAL_ENOMEM                // Out of memory
	MMAL_ENOSPC                // Out of resources (other than memory)
	MMAL_EINVAL                // Argument is invalid
	MMAL_ENOSYS                // Function not implemented
	MMAL_ENOENT                // No such file or directory
	MMAL_ENXIO                 // No such device or address
	MMAL_EIO                   // I/O error
	MMAL_ESPIPE                // Illegal seek
	MMAL_ECORRUPT              // Data is corrupt
	MMAL_ENOTREADY             // Component is not ready
	MMAL_ECONFIG               // Component is not configured
	MMAL_EISCONN               // Port is already connected
	MMAL_ENOTCONN              // Port is disconnected
	MMAL_EAGAIN                // Resource temporarily unavailable. Try again later
	MMAL_EFAULT                // Bad address
	MMAL_MAX       = MMAL_EFAULT
)

const (
	MMAL_PORT_TYPE_UNKNOWN MMAL_PortType = iota
	MMAL_PORT_TYPE_CONTROL               // Control port
	MMAL_PORT_TYPE_INPUT                 // Input port
	MMAL_PORT_TYPE_OUTPUT                // Output port
	MMAL_PORT_TYPE_CLOCK                 // Clock port
	MMAL_PORT_TYPE_MAX     = MMAL_PORT_TYPE_CLOCK
	MMAL_PORT_TYPE_NONE    = MMAL_PORT_TYPE_UNKNOWN
)

const (
	MMAL_PORT_CAPABILITY_PASSTHROUGH                  MMAL_PortCapability = 0x01
	MMAL_PORT_CAPABILITY_ALLOCATION                   MMAL_PortCapability = 0x02
	MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE MMAL_PortCapability = 0x04
	MMAL_PORT_CAPABILITY_MAX                                              = MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE
	MMAL_PORT_CAPABILITY_MIN                                              = MMAL_PORT_CAPABILITY_PASSTHROUGH
)

const (
	MMAL_COMPONENT_DEFAULT_VIDEO_DECODER   = "vc.ril.video_decode"
	MMAL_COMPONENT_DEFAULT_VIDEO_ENCODER   = "vc.ril.video_encode"
	MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER  = "vc.ril.video_render"
	MMAL_COMPONENT_DEFAULT_IMAGE_DECODER   = "vc.ril.image_decode"
	MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER   = "vc.ril.image_encode"
	MMAL_COMPONENT_DEFAULT_CAMERA          = "vc.ril.camera"
	MMAL_COMPONENT_DEFAULT_VIDEO_CONVERTER = "vc.video_convert"
	MMAL_COMPONENT_DEFAULT_SPLITTER        = "vc.splitter"
	MMAL_COMPONENT_DEFAULT_SCHEDULER       = "vc.scheduler"
	MMAL_COMPONENT_DEFAULT_VIDEO_INJECTER  = "vc.video_inject"
	MMAL_COMPONENT_DEFAULT_VIDEO_SPLITTER  = "vc.ril.video_splitter"
	MMAL_COMPONENT_DEFAULT_AUDIO_DECODER   = "none"
	MMAL_COMPONENT_DEFAULT_AUDIO_RENDERER  = "vc.ril.audio_render"
	MMAL_COMPONENT_DEFAULT_MIRACAST        = "vc.miracast"
	MMAL_COMPONENT_DEFAULT_CLOCK           = "vc.clock"
	MMAL_COMPONENT_DEFAULT_CAMERA_INFO     = "vc.camera_info"
)

const (
	MMAL_PARAMETER_GROUP_COMMON   MMAL_ParameterType = (iota << 16)
	MMAL_PARAMETER_GROUP_CAMERA   MMAL_ParameterType = (iota << 16) // Camera-specific parameter ID group
	MMAL_PARAMETER_GROUP_VIDEO    MMAL_ParameterType = (iota << 16) // Video-specific parameter ID group
	MMAL_PARAMETER_GROUP_AUDIO    MMAL_ParameterType = (iota << 16) // Audio-specific parameter ID group
	MMAL_PARAMETER_GROUP_CLOCK    MMAL_ParameterType = (iota << 16) // Clock-specific parameter ID group
	MMAL_PARAMETER_GROUP_MIRACAST MMAL_ParameterType = (iota << 16) // Miracast-specific parameter ID group
	MMAL_PARAMETER_GROUP_MAX                         = MMAL_PARAMETER_GROUP_MIRACAST
	MMAL_PARAMETER_GROUP_MIN                         = MMAL_PARAMETER_GROUP_COMMON
)

const (
	// MMAL_PARAMETER_GROUP_COMMON
	_                                   = iota
	MMAL_PARAMETER_SUPPORTED_ENCODINGS  // Takes a MMAL_PARAMETER_ENCODING_T
	MMAL_PARAMETER_URI                  // Takes a MMAL_PARAMETER_URI_T
	MMAL_PARAMETER_CHANGE_EVENT_REQUEST // Takes a MMAL_PARAMETER_CHANGE_EVENT_REQUEST_T
	MMAL_PARAMETER_ZERO_COPY            // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_BUFFER_REQUIREMENTS  // Takes a MMAL_PARAMETER_BUFFER_REQUIREMENTS_T
	MMAL_PARAMETER_STATISTICS           // Takes a MMAL_PARAMETER_STATISTICS_T
	MMAL_PARAMETER_CORE_STATISTICS      // Takes a MMAL_PARAMETER_CORE_STATISTICS_T
	MMAL_PARAMETER_MEM_USAGE            // Takes a MMAL_PARAMETER_MEM_USAGE_T
	MMAL_PARAMETER_BUFFER_FLAG_FILTER   // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_SEEK                 // Takes a MMAL_PARAMETER_SEEK_T
	MMAL_PARAMETER_POWERMON_ENABLE      // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_LOGGING              // Takes a MMAL_PARAMETER_LOGGING_T
	MMAL_PARAMETER_SYSTEM_TIME          // Takes a MMAL_PARAMETER_UINT64_T
	MMAL_PARAMETER_NO_IMAGE_PADDING     // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_LOCKSTEP_ENABLE      // Takes a MMAL_PARAMETER_BOOLEAN_T
)

const (
	// MMAL_ES_TYPE_T
	MMAL_STREAM_TYPE_UNKNOWN    MMAL_StreamType = iota // Unknown elementary stream type
	MMAL_STREAM_TYPE_CONTROL                           // Elementary stream of control commands
	MMAL_STREAM_TYPE_AUDIO                             // Audio elementary stream
	MMAL_STREAM_TYPE_VIDEO                             //  Video elementary stream
	MMAL_STREAM_TYPE_SUBPICTURE                        // Sub-picture elementary stream (e.g. subtitles, overlays)
	MMAL_STREAM_TYPE_MIN        = MMAL_STREAM_TYPE_UNKNOWN
	MMAL_STREAM_TYPE_MAX        = MMAL_STREAM_TYPE_SUBPICTURE
)

const (
	// MMAL_StreamCompareFlags
	MMAL_STREAM_COMPARE_FLAG_TYPE               MMAL_StreamCompareFlags = 0x0001 // The type is different
	MMAL_STREAM_COMPARE_FLAG_ENCODING           MMAL_StreamCompareFlags = 0x0002 // The encoding is different
	MMAL_STREAM_COMPARE_FLAG_BITRATE            MMAL_StreamCompareFlags = 0x0004 // The bitrate is different
	MMAL_STREAM_COMPARE_FLAG_FLAGS              MMAL_StreamCompareFlags = 0x0008 // The flags are different
	MMAL_STREAM_COMPARE_FLAG_EXTRADATA          MMAL_StreamCompareFlags = 0x0010 // The extradata is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_RESOLUTION   MMAL_StreamCompareFlags = 0x0100 // The video resolution is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_CROPPING     MMAL_StreamCompareFlags = 0x0200 // The video cropping is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_FRAME_RATE   MMAL_StreamCompareFlags = 0x0400 // The video frame rate is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_ASPECT_RATIO MMAL_StreamCompareFlags = 0x0800 // The video aspect ratio is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_COLOR_SPACE  MMAL_StreamCompareFlags = 0x1000 // The video color space is different
	MMAL_STREAM_COMPARE_FLAG_MIN                                        = MMAL_STREAM_COMPARE_FLAG_TYPE
	MMAL_STREAM_COMPARE_FLAG_MAX                                        = MMAL_STREAM_COMPARE_FLAG_VIDEO_COLOR_SPACE
)

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
// STRINGIFY

func (s MMAL_Status) String() string {
	return C.GoString(C.mmal_status_to_string(C.MMAL_STATUS_T(s)))
}

func (s MMAL_Status) Error() string {
	switch s {
	case MMAL_SUCCESS:
		return "MMAL_SUCCESS"
	case MMAL_ENOMEM:
		return "MMAL_ENOMEM: Out of memory"
	case MMAL_ENOSPC:
		return "MMAL_ENOSPC: Out of resources (other than memory)"
	case MMAL_EINVAL:
		return "MMAL_EINVAL: Argument is invalid"
	case MMAL_ENOSYS:
		return "MMAL_ENOSYS: Function not implemented"
	case MMAL_ENOENT:
		return "MMAL_ENOENT: No such file or directory"
	case MMAL_ENXIO:
		return "MMAL_ENXIO: No such device or address"
	case MMAL_EIO:
		return "MMAL_EIO: I/O error"
	case MMAL_ESPIPE:
		return "MMAL_ESPIPE: Illegal seek"
	case MMAL_ECORRUPT:
		return "MMAL_ECORRUPT: Data is corrupt"
	case MMAL_ENOTREADY:
		return "MMAL_ENOTREADY: Component is not ready"
	case MMAL_ECONFIG:
		return "MMAL_ECONFIG: Component is not configured"
	case MMAL_EISCONN:
		return "MMAL_EISCONN: Port is already connected"
	case MMAL_ENOTCONN:
		return "MMAL_ENOTCONN: Port is disconnected"
	case MMAL_EAGAIN:
		return "MMAL_EAGAIN: Resource temporarily unavailable. Try again later"
	case MMAL_EFAULT:
		return "MMAL_EFAULT: Bad address"
	default:
		return "[?? Invalid MMAL_StatusType value]"
	}
}

func (p MMAL_PortType) String() string {
	switch p {
	case MMAL_PORT_TYPE_UNKNOWN:
		return "MMAL_PORT_TYPE_UNKNOWN"
	case MMAL_PORT_TYPE_CONTROL:
		return "MMAL_PORT_TYPE_CONTROL"
	case MMAL_PORT_TYPE_INPUT:
		return "MMAL_PORT_TYPE_INPUT"
	case MMAL_PORT_TYPE_OUTPUT:
		return "MMAL_PORT_TYPE_OUTPUT"
	case MMAL_PORT_TYPE_CLOCK:
		return "MMAL_PORT_TYPE_CLOCK"
	default:
		return "[?? Invalid MMAL_PortType value]"
	}
}

func (c MMAL_PortCapability) String() string {
	parts := ""
	for flag := MMAL_PORT_CAPABILITY_MIN; flag <= MMAL_PORT_CAPABILITY_MAX; flag <<= 1 {
		if c&flag == 0 {
			continue
		}
		switch flag {
		case MMAL_PORT_CAPABILITY_PASSTHROUGH:
			parts += "|" + "MMAL_PORT_CAPABILITY_PASSTHROUGH"
		case MMAL_PORT_CAPABILITY_ALLOCATION:
			parts += "|" + "MMAL_PORT_CAPABILITY_ALLOCATION"
		case MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE:
			parts += "|" + "MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE"
		default:
			parts += "|" + "[?? Invalid MMAL_PortCapability value]"
		}
	}
	return strings.Trim(parts, "|")
}

func (s MMAL_StreamType) String() string {
	switch s {
	case MMAL_STREAM_TYPE_UNKNOWN:
		return "MMAL_STREAM_TYPE_UNKNOWN"
	case MMAL_STREAM_TYPE_CONTROL:
		return "MMAL_STREAM_TYPE_CONTROL"
	case MMAL_STREAM_TYPE_AUDIO:
		return "MMAL_STREAM_TYPE_AUDIO"
	case MMAL_STREAM_TYPE_VIDEO:
		return "MMAL_STREAM_TYPE_VIDEO"
	case MMAL_STREAM_TYPE_SUBPICTURE:
		return "MMAL_STREAM_TYPE_SUBPICTURE"
	default:
		return "[?? Invalid MMAL_StreamType value]"
	}
}

func (f MMAL_StreamCompareFlags) String() string {
	parts := ""
	for flag := MMAL_STREAM_COMPARE_FLAG_MIN; flag <= MMAL_STREAM_COMPARE_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case MMAL_STREAM_COMPARE_FLAG_TYPE:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_TYPE"
		case MMAL_STREAM_COMPARE_FLAG_ENCODING:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_ENCODING"
		case MMAL_STREAM_COMPARE_FLAG_BITRATE:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_BITRATE"
		case MMAL_STREAM_COMPARE_FLAG_FLAGS:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_FLAGS"
		case MMAL_STREAM_COMPARE_FLAG_EXTRADATA:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_EXTRADATA"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_RESOLUTION:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_RESOLUTION"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_CROPPING:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_CROPPING"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_FRAME_RATE:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_FRAME_RATE"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_ASPECT_RATIO:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_ASPECT_RATIO"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_COLOR_SPACE:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_COLOR_SPACE"
		default:
			parts += "|" + "[?? Invalid MMAL_StreamCompareFlags value]"
		}
	}
	return strings.Trim(parts, "|")
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - CALLBACKS

//export mmal_port_callback
func mmal_port_callback(port *C.MMAL_PORT_T, buffer *C.MMAL_BUFFER_HEADER_T) {
	fmt.Printf("mmal_port_callback port=%v buffer=%v\n", port, buffer)
}

func mmal_to_bool(value bool) C.MMAL_BOOL_T {
	if value {
		return MMAL_BOOL_TRUE
	} else {
		return MMAL_BOOL_FALSE
	}
}
