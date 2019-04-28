// +build linux,!rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Hardware struct{}

type hardware struct {
	log     gopi.Logger
	serial  net.HardwareAddr
	sysinfo syscall.Utsname
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config Hardware) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("<hw.linux>Open{}")

	// Create hardware object
	this := new(hardware)
	this.log = logger

	// For Linux we set the serial number from the MAC address
	if ifaces, err := net.Interfaces(); err != nil {
		return nil, err
	} else if len(ifaces) == 0 {
		logger.Error("hw.linux: No network interfaces")
		return nil, gopi.ErrAppError
	} else {
		for _, iface := range ifaces {
			if iface.HardwareAddr != nil {
				this.serial = iface.HardwareAddr
			}
		}
	}

	// Grab the machine details
	if err := syscall.Uname(&this.sysinfo); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}

// Close
func (this *hardware) Close() error {
	this.log.Debug("<hw.linux>Close{ serial=%v }", this.serial)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GetName returns the name of the hardware
func (this *hardware) Name() string {
	return fmt.Sprintf("%v %v (%v)", string(this.sysinfo.Sysname[:]), string(this.sysinfo.Machine[:]), string(this.sysinfo.Release[:]))
}

// SerialNumber returns the serial number of the hardware. For linux
// this is the MAC address with the colons removed
func (this *hardware) SerialNumber() string {
	return strings.Replace(strings.ToUpper(this.serial.String()), ":", "", -1)
}

// Return the number of displays which can be opened
func (this *hardware) NumberOfDisplays() uint {
	return 0
}

// Return Host Uptime
func (this *hardware) UptimeHost() time.Duration {
	if info := this.sysinfo_(); info != nil {
		return time.Second * time.Duration(info.Uptime)
	} else {
		return 0
	}
}

// Return load averages
func (this *hardware) LoadAverage() (float64, float64, float64) {
	if info := this.sysinfo_(); info != nil {
		return float64(info.Loads[0]) / float64(1<<16), float64(info.Loads[1]) / float64(1<<16), float64(info.Loads[2]) / float64(1<<16)
	} else {
		return 0, 0, 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET SYSTEM INFO STRUCTURE

func (this *hardware) sysinfo_() *syscall.Sysinfo_t {
	info := syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(&info); err != nil {
		this.log.Error("<hw.linux>sysinfo: %v", err)
		return nil
	} else {
		return &info
	}
}
