/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package lirc

import (
	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi-hw/sys/filepoll"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register LIRC
	gopi.RegisterModule(gopi.Module{
		Name:     "hw/lirc",
		Type:     gopi.MODULE_TYPE_LIRC,
		Requires: []string{"hw/filepoll"},
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("lirc.device", "", "LIRC device")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			device, _ := app.AppFlags.GetString("lirc.device")
			return gopi.Open(LIRC{
				Device:   device,
				FilePoll: app.ModuleInstance("hw/filepoll").(filepoll.FilePollInterface),
			}, app.Logger)
		},
	})
}
