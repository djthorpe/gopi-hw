// +build rpi

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
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rpi "github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Hardware struct {
	Metrics gopi.Metrics
}

type hardware struct {
	log     gopi.Logger
	service int
	serial  uint64
	product uint32
	done    chan struct{}
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config Hardware) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("hw.rpi.Open{ metrics=%v  }", config.Metrics)

	// Create hardware object
	this := new(hardware)
	this.log = logger

	// Initialise
	if err := rpi.BCMHostInit(); err != nil {
		return nil, err
	}
	if service, err := rpi.VCGencmdInit(); err != nil {
		return nil, err
	} else {
		this.service = service
	}

	// Set serial and revision
	if serial, product, err := rpi.VCGetSerialRevision(); err != nil {
		return nil, err
	} else {
		this.serial = serial
		this.product = product
	}

	// Get channel for updating core CPU temperature
	if core_temp_chan, err := config.Metrics.NewMetricUint(gopi.METRIC_TYPE_CELCIUS, gopi.METRIC_RATE_HOUR, "core_temp"); err != nil {
		return nil, err
	} else {
		this.done = make(chan struct{})
		// record the temperature every minute
		go this.recordTemperature(core_temp_chan, time.Minute)
	}

	// Success
	return this, nil
}

// Close
func (this *hardware) Close() error {
	this.log.Debug("hw.rpi.Close{ }")

	// Stop recording temperature
	if this.done != nil {
		this.done <- gopi.DONE
		<-this.done
		this.done = nil
	}

	// vcgencmd interface
	if this.service != rpi.GENCMD_SERVICE_NONE {
		if err := rpi.VCGencmdTerminate(); err != nil {
			rpi.BCMHostTerminate()
			return err
		}
		this.service = rpi.GENCMD_SERVICE_NONE
	}

	// host terminate
	if err := rpi.BCMHostTerminate(); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GetName returns the name of the hardware
func (this *hardware) Name() string {
	if product_info := rpi.GetProductInfo(this.product); product_info == nil {
		this.log.Warn("rpi.ProductInfo: Invalid product")
		return ""
	} else {
		return fmt.Sprintf("%v (revision %v)", product_info.Model, product_info.Revision)
	}
}

// SerialNumber returns the serial number of the hardware, if available
func (this *hardware) SerialNumber() string {
	return fmt.Sprintf("%08X", this.serial)
}

// Return the number of displays which can be opened
func (this *hardware) NumberOfDisplays() uint {
	return uint(rpi.DX_ID_MAX) + 1
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *hardware) String() string {
	if product_info := rpi.GetProductInfo(this.product); product_info == nil {
		return fmt.Sprintf("hw.rpi{ INVALID PRODUCT }")
	} else {
		params := []string{
			fmt.Sprintf("name=%v", this.Name()),
			fmt.Sprintf("serial=0x%X", this.serial),
			fmt.Sprintf("product=%v", product_info),
			fmt.Sprintf("displays=%v", this.NumberOfDisplays()),
			fmt.Sprintf("peripheral_addr=0x%08X", rpi.BCMHostGetPeripheralAddress()),
			fmt.Sprintf("peripheral_size=0x%08X", rpi.BCMHostGetPeripheralSize()),
		}
		return fmt.Sprintf("hw.rpi{ %v }", strings.Join(params, " "))
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *hardware) recordTemperature(core_temp_chan chan<- uint, delta time.Duration) {
	this.log.Debug2("recordTemperature started with delta=%v", delta)
	interval := time.NewTimer(0)
FOR_LOOP:
	for {
		select {
		case <-interval.C:
			if cpu_temp, err := rpi.VCGetCoreTemperatureCelcius(); err != nil {
				this.log.Error("recordTemperature: %v", err)
			} else if core_temp_chan != nil {
				core_temp_chan <- uint(cpu_temp)
			}
			interval.Reset(delta)
		case <-this.done:
			interval.Stop()
			close(this.done)
			break FOR_LOOP
		}
	}
	this.log.Debug2("recordTemperature ended")
}
