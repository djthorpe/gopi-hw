/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Outputs a table of displays - works on RPi at the moment
package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/olekukonko/tablewriter"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/gpio"
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/i2c"
	_ "github.com/djthorpe/gopi-hw/sys/lirc"
	_ "github.com/djthorpe/gopi-hw/sys/spi"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func moduleName(name string) string {
	if module := gopi.ModuleByName(name); module == nil {
		return "-"
	} else {
		return module.Name
	}
}

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	if app.Hardware == nil {
		return errors.New("No hardware detected")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Hardware
	table.Append([]string{"name", fmt.Sprint(app.Hardware.Name())})
	table.Append([]string{"serial_number", fmt.Sprint(app.Hardware.SerialNumber())})
	table.Append([]string{"number_of_displays", fmt.Sprint(app.Hardware.NumberOfDisplays())})
	table.Append([]string{"uptime_host", fmt.Sprint(app.Hardware.UptimeHost().Truncate(time.Second))})
	table.Append([]string{"load_average", fmt.Sprint(app.Hardware.LoadAverage())})

	// Module names
	table.Append([]string{"hw", moduleName("hw")})
	table.Append([]string{"gpio", moduleName("gpio")})
	table.Append([]string{"i2c", moduleName("i2c")})
	table.Append([]string{"spi", moduleName("spi")})
	table.Append([]string{"lirc", moduleName("lirc")})

	table.Render()

	return nil
}
