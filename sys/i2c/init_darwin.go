// +build darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package i2c

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register I2C
	gopi.RegisterModule(gopi.Module{
		Name: "hw/i2c/darwin",
		Type: gopi.MODULE_TYPE_I2C,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return nil, gopi.ErrNotImplemented
		},
	})
}
