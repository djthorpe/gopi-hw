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
	_ "github.com/djthorpe/gopi-hw/sys/mmal"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		return errors.New("Missing MMAL module")
	} else if component, err := mmal.ComponentWithName("vc.ril.video_decode"); err != nil {
		return err
	} else {
		fmt.Println(component)
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