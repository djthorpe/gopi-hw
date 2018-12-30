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
	if renderer, err := mmal.ComponentWithName("vc.ril.video_render"); err != nil {
		return nil, err
	} else if port := renderer.Inputs()[0]; port == nil {
		return nil, gopi.ErrBadParameter
	} else if display_region, err := port.DisplayRegion(); err != nil {
		return nil, err
	} else {
		display_region.SetFullScreen(true)
		if err := port.SetDisplayRegion(display_region); err != nil {
			return nil, err
		} else {
			return port, nil
		}
	}
}

func ReaderOutputPort(mmal hw.MMAL, uri string) (hw.MMALPort, error) {
	if reader, err := mmal.ComponentWithName("container_reader"); err != nil {
		return nil, err
	} else if output_port := reader.Outputs()[0]; output_port == nil {
		return nil, gopi.ErrBadParameter
	} else if err := reader.Control().SetUri(uri); err != nil {
		return nil, err
	} else {
		return output_port, nil
	}
}

func DecoderInputOutputPorts(mmal hw.MMAL) (hw.MMALPort, hw.MMALPort, error) {
	if decoder, err := mmal.ComponentWithName("vc.ril.video_decode"); err != nil {
		return nil, nil, err
	} else if input_port := decoder.Inputs()[0]; input_port == nil {
		return nil, nil, gopi.ErrBadParameter
	} else if output_port := decoder.Outputs()[0]; output_port == nil {
		return nil, nil, gopi.ErrBadParameter
	} else {
		return input_port, output_port, nil
	}
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	args := app.AppFlags.Args()
	if len(args) == 0 {
		return fmt.Errorf("Missing filename")
	}

	if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		return errors.New("Missing MMAL module")
	} else if reader_port, err := ReaderOutputPort(mmal, args[0]); err != nil {
		return err
	} else if decoder_in_port, decoder_out_port, err := DecoderInputOutputPorts(mmal); err != nil {
		return err
	} else if renderer_port, err := RendererInputPort(mmal); err != nil {
		return err
	} else if c1, err := mmal.Connect(reader_port, decoder_in_port, hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
		return err
	} else if c2, err := mmal.Connect(decoder_out_port, renderer_port, hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
		return err
	} else if err := c1.SetEnabled(true); err != nil {
		return err
	} else if err := c2.SetEnabled(true); err != nil {
		return err
	} else {
		// Display video until interrupted
		fmt.Println("Press CTRL+C to exit")
		app.WaitForSignal()

		// Disconnect
		if err := mmal.Disconnect(c1); err != nil {
			return err
		} else if err := mmal.Disconnect(c2); err != nil {
			return err
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
