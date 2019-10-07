// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package camera

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name:     "media/camera",
		Type:     gopi.MODULE_TYPE_OTHER,
		Requires: []string{"hw/mmal"},
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("camera.id", 0, "Camera selection")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			camera_id, _ := app.AppFlags.GetUint("camera.id")
			return gopi.Open(Camera{
				MMAL:   app.ModuleInstance("hw/mmal").(hw.MMAL),
				Camera: uint32(camera_id),
			}, app.Logger)
		},
	})
}
