/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	"encoding/binary"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALFormatType          uint
	MMALEncodingType        uint32
	MMALDisplayTransform    uint
	MMALDisplayMode         uint
	MMALPortConnectionFlags uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type MMAL interface {
	gopi.Driver

	// Return component
	ComponentWithName(name string) (MMALComponent, error)

	// Connect and disconnect component ports
	Connect(input, output MMALPort, flags MMALPortConnectionFlags) (MMALPortConnection, error)
	Disconnect(MMALPortConnection) error
}

type MMALComponent interface {
	Name() string
	Id() uint32

	// Enable and disable
	Enabled() bool
	SetEnabled(value bool) error

	// Acquire and release
	Acquire() error
	Release() error

	// Return ports
	Control() MMALPort
	Clock() []MMALPort
	Input() []MMALPort
	Output() []MMALPort
}

type MMALPort interface {
	Name() string
	CapabilityPassthrough() bool
	CapabilityAllocation() bool
	CapabilitySupportsEventFormatChange() bool

	// Enable and disable
	Enabled() bool
	SetEnabled(value bool) error

	// Connect and Disconnect & Flush
	Connect(other MMALPort) error
	Disconnect() error
	Flush() error

	// Formats
	Format() MMALFormat
	CommitFormatChange() error

	// Get/Set Port Parameters
	MMALCommonParameters
	MMALVideoParameters
}

type MMALPortConnection interface {
	// Input and Output ports
	Input() MMALPort
	Output() MMALPort

	// Enable and disable
	Enabled() bool
	SetEnabled(value bool) error

	// Acquire and release
	Acquire() error
	Release() error
}

type MMALCommonParameters interface {
	// Get Parameters
	SupportedEncodings() ([]MMALEncodingType, error)
	Uri() (string, error)
	ZeroCopy() (bool, error)
	NoImagePadding() (bool, error)
	LockstepEnable() (bool, error)
	PowermonEnable() (bool, error)
	BufferFlagFilter() (uint32, error)
	SystemTime() (uint64, error)

	// Set Parameters
	SetUri(value string) error
	SetZeroCopy(value bool) error
	SetNoImagePadding(value bool) error
	SetLockstepEnable(value bool) error
	SetPowermonEnable(value bool) error
	SetBufferFlagFilter(value uint32) error
}

type MMALVideoParameters interface {
	// Get Parameters
	GetDisplayRegion() (MMALDisplayRegion, error)

	// Set Parameters
	SetDisplayRegion(MMALDisplayRegion) error
}

type MMALFormat interface {
	Type() MMALFormatType
}

type MMALDisplayRegion interface {
	// Get properties
	Display() uint16
	FullScreen() bool
	Layer() int16
	Alpha() uint8
	Transform() MMALDisplayTransform
	NoAspect() bool
	Mode() MMALDisplayMode
	CopyProtect() bool
	DestRect() (int32, int32, uint32, uint32)
	SrcRect() (int32, int32, uint32, uint32)

	// Set properties
	SetFullScreen(bool)
	SetLayer(int16)
	SetAlpha(uint8)
	SetTransform(MMALDisplayTransform)
	SetNoAspect(bool)
	SetMode(MMALDisplayMode)
	SetCopyProtect(bool)
	SetDestRect(x, y int32, width, height uint32)
	SetSrcRect(x, y int32, width, height uint32)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_FORMAT_UNKNOWN MMALFormatType = iota
	MMAL_FORMAT_CONTROL
	MMAL_FORMAT_AUDIO
	MMAL_FORMAT_VIDEO
	MMAL_FORMAT_SUBPICTURE
	MMAL_FORMAT_MAX = MMAL_FORMAT_SUBPICTURE
)

const (
	MMAL_DISPLAY_TRANSFORM_NONE MMALDisplayTransform = iota
	MMAL_DISPLAY_TRANSFORM_MIRROR
	MMAL_DISPLAY_TRANSFORM_ROT180_MIRROR
	MMAL_DISPLAY_TRANSFORM_ROT180
	MMAL_DISPLAY_TRANSFORM_ROT90_MIRROR
	MMAL_DISPLAY_TRANSFORM_ROT270
	MMAL_DISPLAY_TRANSFORM_ROT90
	MMAL_DISPLAY_TRANSFORM_ROT270_MIRROR
	MMAL_DISPLAY_TRANSFORM_MAX = MMAL_DISPLAY_TRANSFORM_ROT270_MIRROR
)

const (
	MMAL_DISPLAY_MODE_FILL MMALDisplayMode = iota
	MMAL_DISPLAY_MODE_LETTERBOX
	MMAL_DISPLAY_MODE_STEREO_LEFT_TO_LEFT
	MMAL_DISPLAY_MODE_STEREO_TOP_TO_TOP
	MMAL_DISPLAY_MODE_STEREO_LEFT_TO_TOP
	MMAL_DISPLAY_MODE_STEREO_TOP_TO_LEFT
	MMAL_DISPLAY_MODE_MAX = MMAL_DISPLAY_MODE_STEREO_TOP_TO_LEFT
)

const (
	// MMAL_PortConnectionFlags
	MMAL_CONNECTION_FLAG_TUNNELLING               MMALPortConnectionFlags = 0x0001 // The connection is tunnelled. Buffer headers do not transit via the client but directly from the output port to the input port.
	MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT      MMALPortConnectionFlags = 0x0002 // Force the pool of buffer headers used by the connection to be allocated on the input port.
	MMAL_CONNECTION_FLAG_ALLOCATION_ON_OUTPUT     MMALPortConnectionFlags = 0x0004 // Force the pool of buffer headers used by the connection to be allocated on the output port.
	MMAL_CONNECTION_FLAG_KEEP_BUFFER_REQUIREMENTS MMALPortConnectionFlags = 0x0008 // Specify that the connection should not modify the buffer requirements.
	MMAL_CONNECTION_FLAG_DIRECT                   MMALPortConnectionFlags = 0x0010 // The connection is flagged as direct. This doesn't change the behaviour of the connection itself but is used by the the graph utility to specify that the buffer should be sent to the input port from with the port callback.
	MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS        MMALPortConnectionFlags = 0x0020 // Specify that the connection should not modify the port formats.
	MMAL_CONNECTION_FLAG_MIN                                              = MMAL_CONNECTION_FLAG_TUNNELLING
	MMAL_CONNECTION_FLAG_MAX                                              = MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t MMALFormatType) String() string {
	switch t {
	case MMAL_FORMAT_UNKNOWN:
		return "MMAL_FORMAT_UNKNOWN"
	case MMAL_FORMAT_CONTROL:
		return "MMAL_FORMAT_CONTROL"
	case MMAL_FORMAT_AUDIO:
		return "MMAL_FORMAT_AUDIO"
	case MMAL_FORMAT_VIDEO:
		return "MMAL_FORMAT_VIDEO"
	case MMAL_FORMAT_SUBPICTURE:
		return "MMAL_FORMAT_SUBPICTURE"
	default:
		return "[?? Unknown MMALFormatType value]"
	}
}

func (e MMALEncodingType) String() string {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(e))
	return "'" + string(buf) + "'"
}

func (t MMALDisplayTransform) String() string {
	switch t {
	case MMAL_DISPLAY_TRANSFORM_NONE:
		return "MMAL_DISPLAY_TRANSFORM_NONE"
	case MMAL_DISPLAY_TRANSFORM_MIRROR:
		return "MMAL_DISPLAY_TRANSFORM_MIRROR"
	case MMAL_DISPLAY_TRANSFORM_ROT180_MIRROR:
		return "MMAL_DISPLAY_TRANSFORM_ROT180_MIRROR"
	case MMAL_DISPLAY_TRANSFORM_ROT180:
		return "MMAL_DISPLAY_TRANSFORM_ROT180"
	case MMAL_DISPLAY_TRANSFORM_ROT90_MIRROR:
		return "MMAL_DISPLAY_TRANSFORM_ROT90_MIRROR"
	case MMAL_DISPLAY_TRANSFORM_ROT270:
		return "MMAL_DISPLAY_TRANSFORM_ROT270"
	case MMAL_DISPLAY_TRANSFORM_ROT90:
		return "MMAL_DISPLAY_TRANSFORM_ROT90"
	case MMAL_DISPLAY_TRANSFORM_ROT270_MIRROR:
		return "MMAL_DISPLAY_TRANSFORM_ROT270_MIRROR"
	default:
		return "[?? Unknown MMALDisplayTransform value]"
	}
}

func (m MMALDisplayMode) String() string {
	switch m {
	case MMAL_DISPLAY_MODE_FILL:
		return "MMAL_DISPLAY_MODE_FILL"
	case MMAL_DISPLAY_MODE_LETTERBOX:
		return "MMAL_DISPLAY_MODE_LETTERBOX"
	case MMAL_DISPLAY_MODE_STEREO_LEFT_TO_LEFT:
		return "MMAL_DISPLAY_MODE_STEREO_LEFT_TO_LEFT"
	case MMAL_DISPLAY_MODE_STEREO_TOP_TO_TOP:
		return "MMAL_DISPLAY_MODE_STEREO_TOP_TO_TOP"
	case MMAL_DISPLAY_MODE_STEREO_LEFT_TO_TOP:
		return "MMAL_DISPLAY_MODE_STEREO_LEFT_TO_TOP"
	case MMAL_DISPLAY_MODE_STEREO_TOP_TO_LEFT:
		return "MMAL_DISPLAY_MODE_STEREO_TOP_TO_LEFT"
	default:
		return "[?? Unknown MMALDisplayMode value]"
	}
}

func (c MMALPortConnectionFlags) String() string {
	parts := ""
	for flag := MMAL_CONNECTION_FLAG_MIN; flag <= MMAL_CONNECTION_FLAG_MAX; flag <<= 1 {
		if c&flag == 0 {
			continue
		}
		switch flag {
		case MMAL_CONNECTION_FLAG_TUNNELLING:
			parts += "|" + "MMAL_CONNECTION_FLAG_TUNNELLING"
		case MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT:
			parts += "|" + "MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT"
		case MMAL_CONNECTION_FLAG_ALLOCATION_ON_OUTPUT:
			parts += "|" + "MMAL_CONNECTION_FLAG_ALLOCATION_ON_OUTPUT"
		case MMAL_CONNECTION_FLAG_KEEP_BUFFER_REQUIREMENTS:
			parts += "|" + "MMAL_CONNECTION_FLAG_KEEP_BUFFER_REQUIREMENTS"
		case MMAL_CONNECTION_FLAG_DIRECT:
			parts += "|" + "MMAL_CONNECTION_FLAG_DIRECT"
		case MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS:
			parts += "|" + "MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS"
		default:
			parts += "|" + "[?? Invalid MMALPortConnectionFlags value]"
		}
	}
	return strings.Trim(parts, "|")
}
