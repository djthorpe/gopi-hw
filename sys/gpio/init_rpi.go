// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package gpio

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name:     "hw/gpio/rpi",
		Type:     gopi.MODULE_TYPE_GPIO,
		Requires: []string{"hw/rpi"},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(GPIO{Hardware: app.Hardware}, app.Logger)
		},
	})
}
