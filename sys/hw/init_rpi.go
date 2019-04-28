// +build rpi

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
		Name: "hw/rpi",
		Type: gopi.MODULE_TYPE_HARDWARE,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Hardware{}, app.Logger)
		},
	})
}
