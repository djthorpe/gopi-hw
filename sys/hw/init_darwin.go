/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register hardware
	gopi.RegisterModule(gopi.Module{
		Name: "hw/darwin",
		Type: gopi.MODULE_TYPE_HARDWARE,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Hardware{}, app.Logger)
		},
	})
}

/*
type Hardware interface {
	Driver

	// Return name of the hardware platform
	Name() string

	// Return unique serial number of this hardware
	SerialNumber() string

	// Return the number of possible displays for this hardware
	NumberOfDisplays() uint
}
*/
// +build darwin
