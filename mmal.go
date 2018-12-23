/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	"github.com/djthorpe/gopi"
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

	// Connect and Disconnect, Flush and Commit Format
	Connect(other MMALPort) error
	Disconnect() error
	Flush() error

	// Formats
	Format() MMALFormat
	CommitFormatChange() error
}

type MMALFormat interface {
}
