// +build linux

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
		Name: "hw/i2c",
		Type: gopi.MODULE_TYPE_I2C,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("i2c.bus", 1, "I2C Bus")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			bus, _ := app.AppFlags.GetUint("i2c.bus")
			return gopi.Open(I2C{
				Bus: bus,
			}, app.Logger)
		},
	})
}
