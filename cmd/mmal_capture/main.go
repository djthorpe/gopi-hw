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
	"fmt"
	"os"
	"reflect"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
	tablewriter "github.com/olekukonko/tablewriter"

	// Modules
	_ "github.com/djthorpe/gopi-graphics/sys/display"
	_ "github.com/djthorpe/gopi-hw/sys/camera"
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/mmal"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func ApplyConfig(camera hw.Camera, display gopi.Display, app *gopi.AppInstance) (hw.CameraConfig, error) {
	// Obtain the configuration
	config, err := camera.CameraConfig()
	if err != nil {
		return config, err
	}

	// Set rotation
	if rotation, exists := app.AppFlags.GetInt("rotate"); exists {
		config.PreviewRotation = int32(rotation)
		config.ImageRotation = int32(rotation)
		config.VideoRotation = int32(rotation)
		config.Flags |= hw.FLAG_ROTATION
	}

	// Set preview to full screen
	if preview, _ := app.AppFlags.GetBool("preview"); preview {
		w, h := display.Size()
		config.PreviewFrameSize = gopi.Size{float32(w), float32(h)}
		config.Flags |= hw.FLAG_PREVIEW_FRAMESIZE
	}

	// Set the image capture to the full frame
	config.ImageFrameSize = camera.CameraFrameSize()
	config.Flags |= hw.FLAG_IMAGE_FRAMESIZE

	// Update the configuration
	if err := camera.SetCameraConfig(config); err != nil {
		return config, err
	}
	return config, nil
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	camera := app.ModuleInstance("media/camera").(hw.Camera)
	display := app.Display
	if camera == nil || display == nil {
		return errors.New("Missing camera or display")
	} else if config, err := ApplyConfig(camera, display, app); err != nil {
		return err
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		config_ := reflect.ValueOf(config)
		table.SetHeader([]string{"Parameter", "Value"})
		for i := 0; i < config_.NumField(); i++ {
			name := config_.Type().Field(i).Name
			if name != "Flags" {
				value := config_.Field(i).Interface()
				table.Append([]string{name, fmt.Sprint(value)})
			}
		}
		table.Render()
	}

	// Start preview
	if preview, _ := app.AppFlags.GetBool("preview"); preview {
		if err := camera.Preview(); err != nil {
			return err
		}
	}

	// Perform image capture
	if err := camera.ImageCapture(); err != nil {
		return err
	}

	// Wait for CTRL+C
	fmt.Println("Press CTRL+C to exit")
	app.WaitForSignal()

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the Camera & Display modules
	config := gopi.NewAppConfig("media/camera", "display")

	// Main command-line parameters
	config.AppFlags.FlagInt("rotate", 0, "Rotation angle (0,90,180 or 270)")
	config.AppFlags.FlagBool("preview", false, "Display capture preview")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main))
}
