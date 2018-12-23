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

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALFormatType   uint
	MMALEncodingType uint32
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type MMAL interface {
	gopi.Driver

	// Return component
	ComponentWithName(name string) (MMALComponent, error)
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

	// Implements common parameters
	MMALCommonParameters
}

type MMALCommonParameters interface {
	// Get Parameters
	SupportedEncodings() ([]MMALEncodingType, error)
	Uri() (string, error)
	ZeroCopy() (bool, error)

	// Set Parameters
	SetUri(value string) error
	SetZeroCopy(value bool) error
}

type MMALFormat interface {
	Type() MMALFormatType
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_FORMAT_UNKNOWN MMALFormatType = iota
	MMAL_FORMAT_CONTROL
	MMAL_FORMAT_AUDIO
	MMAL_FORMAT_VIDEO
	MMAL_FORMAT_SUBPICTURE
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
