/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package spi

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name: "hw/spi",
		Type: gopi.MODULE_TYPE_SPI,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("spi.bus", 0, "SPI Bus")
			config.AppFlags.FlagUint("spi.slave", 0, "SPI Slave")
			config.AppFlags.FlagUint("spi.delay", 0, "SPI Transfer delay in microseconds")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			bus, _ := app.AppFlags.GetUint("spi.bus")
			slave, _ := app.AppFlags.GetUint("spi.slave")
			delay, _ := app.AppFlags.GetUint16("spi.delay")
			return gopi.Open(SPI{
				Bus:   bus,
				Slave: slave,
				Delay: delay,
			}, app.Logger)
		},
	})
}
