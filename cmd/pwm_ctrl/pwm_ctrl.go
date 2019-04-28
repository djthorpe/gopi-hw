/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Controls PWM on GPIO pins
package main

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/gpio"
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/pwm"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if app.PWM == nil {
		return app.Logger.Error("Missing PWM module instance")
	}

	if err := app.PWM.SetDutyCycle(0, 4); err != nil {
		return err
	}

	fmt.Println(app.PWM)

	// Finished
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the spi instance
	config := gopi.NewAppConfig("pwm")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
