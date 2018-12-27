/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// MMAL example to encode an image
package main

import (
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

func MMALEncodeTest(app *gopi.AppInstance, mmal hw.MMAL, encoder hw.MMALComponent, format hw.MMALEncodingType, width, height uint32) error {
	port_in := encoder.Inputs()[0]
	port_out := encoder.Outputs()[0]
	if port_in.Enabled() {
		if err := port_in.SetEnabled(false); err != nil {
			return err
		}
	}
	if port_out.Enabled() {
		if err := port_out.SetEnabled(false); err != nil {
			return err
		}
	}

	// Set input port to uncompressed RGBA
	port_in.VideoFormat().SetEncoding(hw.MMAL_ENCODING_RGBA)
	port_in.VideoFormat().SetWidthHeight(width, height)
	port_in.VideoFormat().SetCrop(hw.MMALRect{0, 0, width, height})
	if err := port_in.CommitFormatChange(); err != nil {
		return err
	}

	// Set output port to input parameter and set encoding
	if err := port_out.CopyFormat(port_in.VideoFormat()); err != nil {
		return err
	}
	port_out.VideoFormat().SetEncoding(format)
	if err := port_out.CommitFormatChange(); err != nil {
		return err
	}

	// Set JPEG Quality factor
	if err := port_out.SetJPEGQFactor(100); err != nil {
		return err
	}

	// Enable the ports
	if err := port_in.SetEnabled(true); err != nil {
		return err
	}
	if err := port_out.SetEnabled(true); err != nil {
		return err
	}

FOR_LOOP:
	for {
		// Wait for an empty buffer to be available on a port, block until one becomes available
		if buffer, err := encoder.GetEmptyBufferOnPort(port_out, true); err != nil {
			return err
		} else {
			fmt.Println("buffer=", buffer)
			break FOR_LOOP
		}
	}

	// Flush output port
	if err := port_out.Flush(); err != nil {
		return err
	}

	// Disable the ports
	if err := port_in.SetEnabled(false); err != nil {
		return err
	}
	if err := port_out.SetEnabled(false); err != nil {
		return err
	}

	// Success
	return nil
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		return fmt.Errorf("Missing MMAL module")
	} else if image_encoder, err := mmal.ImageEncoderComponent(); err != nil {
		return err
	} else if err := MMALEncodeTest(app, mmal, image_encoder, hw.MMAL_ENCODING_JPEG, 1920, 1080); err != nil {
		return err
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
