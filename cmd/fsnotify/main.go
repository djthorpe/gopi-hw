/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/fsnotify"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func EventLoop(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {
	fsnotify := app.ModuleInstance("hw/fsnotify").(hw.FSNotify)
	messages := fsnotify.Subscribe()
	fmt.Printf("%-20s %-20s %s\n", "Root", "Path", "Flags")
	fmt.Printf("%20s %20s\n", "--------------------", "--------------------")

	start <- gopi.DONE

FOR_LOOP:
	for {
		select {
		case evt := <-messages:
			if event, ok := evt.(hw.FSEvent); ok {
				fmt.Printf("%-20s %-20s %s\n", event.Root(), event.RelPath(), fmt.Sprint(event.Flags()))
			} else {
				fmt.Println(evt)
			}
		case <-stop:
			break FOR_LOOP
		}
	}

	// End of routine
	fsnotify.Unsubscribe(messages)
	return nil
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	fsnotify := app.ModuleInstance("hw/fsnotify").(hw.FSNotify)

	if args := app.AppFlags.Args(); len(args) == 0 {
		return gopi.ErrHelp
	} else {
		for _, path := range args {
			if err := fsnotify.Watch(path); err != nil {
				return err
			}
		}
	}

	// Wait for signal
	app.Logger.Info("Press CTRL+C to terminate")
	app.WaitForSignal()

	// Finished
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	config := gopi.NewAppConfig("hw/fsnotify")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main, EventLoop))
}
