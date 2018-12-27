/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// MMAL example to connect camera to screen and render the camera
// image
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

func RendererInputPort(mmal hw.MMAL) (hw.MMALPort, error) {
	if renderer, err := mmal.VideoRendererComponent(); err != nil {
		return nil, err
	} else if port := renderer.Inputs()[0]; port == nil {
		return nil, gopi.ErrBadParameter
	} else if display_region, err := port.DisplayRegion(); err != nil {
		return nil, err
	} else {
		display_region.SetFullScreen(true)
		display_region.SetTransform(hw.MMAL_DISPLAY_TRANSFORM_ROT180_MIRROR)
		display_region.SetMode(hw.MMAL_DISPLAY_MODE_FILL)
		if err := port.SetDisplayRegion(display_region); err != nil {
			return nil, err
		} else {
			return port, nil
		}
	}
}

func CameraOutputPort(mmal hw.MMAL) (hw.MMALPort, error) {
	if camera, err := mmal.CameraComponent(); err != nil {
		return nil, err
	} else if port := camera.Outputs()[0]; port == nil {
		return nil, gopi.ErrBadParameter
	} else if annotation, err := camera.Control().Annotation(); err != nil {
		return nil, err
	} else {
		// Set camera framesize
		port.VideoFormat().SetWidthHeight(100, 100)
		if err := port.CommitFormatChange(); err != nil {
			return nil, err
		}

		annotation.SetText("Hello, world")
		annotation.SetTextSize(24)
		fmt.Println(annotation)
		if err := camera.Control().SetAnnotation(annotation); err != nil {
			return nil, err
		}
		return port, nil
	}
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		return errors.New("Missing MMAL module")
	} else if camera_port, err := CameraOutputPort(mmal); err != nil {
		return err
	} else if renderer_port, err := RendererInputPort(mmal); err != nil {
		return err
	} else if c, err := mmal.Connect(renderer_port, camera_port, hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
		return err
	} else if err := c.SetEnabled(true); err != nil {
		return err
	} else {
		// Display camera preview until interrupted
		fmt.Println("Press CTRL+C to exit")
		app.WaitForSignal()
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
