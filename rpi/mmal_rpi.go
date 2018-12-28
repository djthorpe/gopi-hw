//+build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

// Frameworks

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -I/opt/vc/include
#cgo LDFLAGS: -L/opt/vc/lib -lmmal -lmmal_core -lmmal_util
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_connection.h>
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

type MMAL_PortCallback func(port MMAL_PortHandle, buffer MMAL_Buffer)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - OTHER

func mmal_to_bool(value bool) C.MMAL_BOOL_T {
	if value {
		return MMAL_BOOL_TRUE
	} else {
		return MMAL_BOOL_FALSE
	}
}
