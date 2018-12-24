/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// MMAL examples
package main

import (
	"errors"
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/metrics"
	_ "github.com/djthorpe/gopi-hw/sys/mmal"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		return errors.New("Missing MMAL module")
	} else if camera, err := mmal.ComponentWithName("vc.ril.camera"); err != nil {
		return err
	} else if renderer, err := mmal.ComponentWithName("vc.ril.video_render"); err != nil {
		return err
	} else if c, err := mmal.Connect(renderer.Input()[0], camera.Output()[0], hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
		return err
	} else if display_region, err := c.Input().GetDisplayRegion(); err != nil {
		return err
	} else {
		display_region.SetFullScreen(true)
		display_region.SetTransform(hw.MMAL_DISPLAY_TRANSFORM_ROT180_MIRROR)
		display_region.SetAlpha(0x50)
		if err := c.Input().SetDisplayRegion(display_region); err != nil {
			return err
		} else if err := c.SetEnabled(true); err != nil {
			return err
		} else {
			fmt.Println("Press CTRL+C to exit")
			app.WaitForSignal()
		}
	}

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the MMAL instance
	config := gopi.NewAppConfig("hw/mmal")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main))
}
