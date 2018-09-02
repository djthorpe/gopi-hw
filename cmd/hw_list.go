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

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/olekukonko/tablewriter"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/metrics"
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

	// Module names
	table.Append([]string{"hw", moduleName("hw")})
	table.Append([]string{"metrics", moduleName("metrics")})
	table.Append([]string{"gpio", moduleName("gpio")})
	table.Append([]string{"i2c", moduleName("i2c")})
	table.Append([]string{"spi", moduleName("spi")})
	table.Append([]string{"lirc", moduleName("lirc")})
	table.Append([]string{"display", moduleName("display")})

	// Metrics
	if metrics := app.ModuleInstance("metrics").(gopi.Metrics); metrics != nil {
		// Uptime
		table.Append([]string{"uptime_host", fmt.Sprint(metrics.UptimeHost())})
		table.Append([]string{"uptime_app", fmt.Sprint(metrics.UptimeApp())})

		// Load Averages
		loadav_1m, loadav_5m, loadav_15m := metrics.LoadAverage()
		table.Append([]string{"load_average_1m", fmt.Sprintf("%.2f", loadav_1m)})
		table.Append([]string{"load_average_5m", fmt.Sprintf("%.2f", loadav_5m)})
		table.Append([]string{"load_average_15m", fmt.Sprintf("%.2f", loadav_15m)})

		// Other metrics
		for _, metric := range metrics.Metrics(gopi.METRIC_TYPE_NONE) {
			table.Append([]string{metric.Name(), fmt.Sprintf("%v%v", metric.FloatValue(), metric.Unit())})
		}
	}

	table.Render()

	app.Logger.Info("Waiting for CTRL+C")
	app.WaitForSignal()

	return nil
}

func main() {
	// Create the configuration, load the gpio instance
	config := gopi.NewAppConfig("hw", "metrics")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
