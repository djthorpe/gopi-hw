/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package pwm

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name:     "hw/pwm/pi-blaster",
		Type:     gopi.MODULE_TYPE_PWM,
		Requires: []string{"gpio"},
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("pi-blaster.fifo", "/dev/pi-blaster", "Path to FIFO")
			config.AppFlags.FlagString("pi-blaster.exec", "/usr/sbin/pi-blaster", "Path to executable")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			fifo, _ := app.AppFlags.GetString("pi-blaster.fifo")
			exec, _ := app.AppFlags.GetString("pi-blaster.exec")
			return gopi.Open(PiBlaster{
				FIFO: fifo,
				Exec: exec,
			}, app.Logger)
		},
	})
}
