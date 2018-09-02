// +build linux,!rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package gpio

import (
	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi-hw/sys/filepoll"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name:     "hw/gpio/linux",
		Requires: []string{"hw/filepoll"},
		Type:     gopi.MODULE_TYPE_GPIO,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagBool("gpio.unexport", true, "Unexport exported pins on exit")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			unexport, _ := app.AppFlags.GetBool("gpio.unexport")
			return gopi.Open(GPIO{
				UnexportOnClose: unexport,
				FilePoll:        app.ModuleInstance("hw/filepoll").(filepoll.FilePollInterface),
			}, app.Logger)
		},
	})
}
