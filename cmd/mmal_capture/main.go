/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Example to capture from the camera using the MMAL library
package main

import (
	"errors"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/mmal"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		return errors.New("Missing MMAL module")
	} else if capture, err := NewApp(mmal, app.Logger); err != nil {
		return err
	} else if err := capture.Setup(app.AppFlags); err != nil {
		return err
	} else if err := capture.CameraInfo(); err != nil {
		return err
	} else if err := capture.Capture(); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the MMAL instance
	config := gopi.NewAppConfig("hw/mmal")

	// Command-line parameters
	config.AppFlags.FlagUint("rotate", 0, "Camera rotation 0-360")
	config.AppFlags.FlagUint("camera", 0, "Camera number")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main))
}
