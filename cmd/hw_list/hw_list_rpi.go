// +build rpi

/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
)

func main() {
	// Create the configuration, load the gpio instance
	config := gopi.NewAppConfig("hw")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
