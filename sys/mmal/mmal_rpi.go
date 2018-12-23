// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package mmal

import (
	"fmt"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
	"github.com/djthorpe/gopi/util/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type MMAL struct {
	Hardware gopi.Hardware
}

type mmal struct {
	log        gopi.Logger
	hardware   gopi.Hardware
	components map[string]*component
}

type component struct {
	log     gopi.Logger
	handle  rpi.MMAL_ComponentHandle
	control *port
	input   []*port
	output  []*port
	clock   []*port
}

type port struct {
	log    gopi.Logger
	handle rpi.MMAL_PortHandle
}

type format struct {
	handle rpi.MMAL_StreamFormat
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config MMAL) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.mmal>Open{ hw=%v }", config.Hardware)
	this := new(mmal)
	this.log = log
	this.hardware = config.Hardware
	this.components = make(map[string]*component, 0)

	return this, nil
}

func (this *mmal) Close() error {
	this.log.Debug("<sys.hw.mmal>Close{ hw=%v }", this.hardware)

	// Close components
	err := new(errors.CompoundError)
	for _, component := range this.components {
		if err_ := component.Close(); err != nil {
			err.Add(err_)
		}
	}

	// Release resources
	this.hardware = nil
	this.components = nil

	return err.ErrorOrSelf()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *mmal) String() string {
	parts := ""
	for k, v := range this.components {
		parts += fmt.Sprintf("%v=%v ", k, v)
	}
	return fmt.Sprintf("<sys.hw.mmal>{ %v }", strings.TrimSpace(parts))
}

////////////////////////////////////////////////////////////////////////////////
// CREATE COMPONENT

func (this *mmal) ComponentWithName(name string) (hw.MMALComponent, error) {
	var handle rpi.MMAL_ComponentHandle

	if c, exists := this.components[name]; exists {
		return c, nil
	} else if err := rpi.MMALComponentCreate(name, &handle); err != nil {
		return nil, err
	} else {
		// Create the component
		c = &component{
			handle: handle,
			log:    this.log,
		}
		// Set control port
		c.control = this.NewPort(rpi.MMALComponentControlPort(handle))
		// Input ports
		c.input = make([]*port, rpi.MMALComponentInputPortNum(handle))
		for i := range c.input {
			c.input[i] = this.NewPort(rpi.MMALComponentInputPortAtIndex(handle, uint(i)))
		}
		// Output ports
		c.output = make([]*port, rpi.MMALComponentOutputPortNum(handle))
		for i := range c.output {
			c.output[i] = this.NewPort(rpi.MMALComponentOutputPortAtIndex(handle, uint(i)))
		}
		// Clock ports
		c.clock = make([]*port, rpi.MMALComponentClockPortNum(handle))
		for i := range c.clock {
			c.clock[i] = this.NewPort(rpi.MMALComponentClockPortAtIndex(handle, uint(i)))
		}
		// Add to the map
		this.components[name] = c
		return c, nil
	}
}

func (this *mmal) NewPort(handle rpi.MMAL_PortHandle) *port {
	return &port{
		handle: handle,
		log:    this.log,
	}
}
