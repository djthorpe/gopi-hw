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
	log         gopi.Logger
	hardware    gopi.Hardware
	components  map[string]*component
	connections []*connection
}

type component struct {
	log     gopi.Logger
	handle  rpi.MMAL_ComponentHandle
	control *port
	input   []*port
	output  []*port
	clock   []*port
	sema    rpi.VCSemaphore
}

type port struct {
	log    gopi.Logger
	handle rpi.MMAL_PortHandle
	pool   rpi.MMAL_Pool
}

type format struct {
	log    gopi.Logger
	handle rpi.MMAL_StreamFormat
}

type connection struct {
	log           gopi.Logger
	handle        rpi.MMAL_PortConnection
	input, output hw.MMALPort
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config MMAL) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.mmal>Open{ hw=%v }", config.Hardware)
	this := new(mmal)
	this.log = log
	this.hardware = config.Hardware
	this.components = make(map[string]*component, 0)
	this.connections = make([]*connection, 0)
	return this, nil
}

func (this *mmal) Close() error {
	this.log.Debug("<sys.hw.mmal>Close{ components=%v connections=%v }", this.components, this.connections)
	err := new(errors.CompoundError)

	// Disconnect connections
	for _, connection := range this.connections {
		if err_ := connection.Close(); err_ != nil {
			err.Add(err_)
		}
	}

	// Close components
	for _, component := range this.components {
		if err_ := component.Close(); err_ != nil {
			err.Add(err_)
		}
	}

	// Release resources
	this.hardware = nil
	this.components = nil
	this.connections = nil

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
// COMPONENTS

func (this *mmal) ComponentWithName(name string) (hw.MMALComponent, error) {
	this.log.Debug2("<sys.hw.mmal>ComponentWithName{ name='%v' }", name)

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

		// Create Semaphore
		if sema, err := rpi.VCSemaphoreCreate(name, 0); err != nil {
			return nil, err
		} else {
			c.sema = sema
		}

		// Enable control port
		if err := rpi.MMALPortEnable(rpi.MMALComponentControlPort(handle)); err != nil {
			return nil, err
		}

		// Add to the map
		this.components[name] = c
		return c, nil
	}
}

func (this *mmal) CameraComponent() (hw.MMALCameraComponent, error) {
	return this.ComponentWithName(rpi.MMAL_COMPONENT_DEFAULT_CAMERA)
}

func (this *mmal) CameraInfoComponent() (hw.MMALComponent, error) {
	return this.ComponentWithName(rpi.MMAL_COMPONENT_DEFAULT_CAMERA_INFO)
}

func (this *mmal) VideoDecoderComponent() (hw.MMALComponent, error) {
	return this.ComponentWithName(rpi.MMAL_COMPONENT_DEFAULT_VIDEO_DECODER)
}

func (this *mmal) VideoEncoderComponent() (hw.MMALComponent, error) {
	return this.ComponentWithName(rpi.MMAL_COMPONENT_DEFAULT_VIDEO_ENCODER)
}

func (this *mmal) VideoRendererComponent() (hw.MMALComponent, error) {
	return this.ComponentWithName(rpi.MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER)
}

func (this *mmal) ImageEncoderComponent() (hw.MMALComponent, error) {
	return this.ComponentWithName(rpi.MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER)
}

func (this *mmal) ImageDecoderComponent() (hw.MMALComponent, error) {
	return this.ComponentWithName(rpi.MMAL_COMPONENT_DEFAULT_IMAGE_DECODER)
}

func (this *mmal) ReaderComponent() (hw.MMALComponent, error) {
	return this.ComponentWithName(rpi.MMAL_COMPONENT_DEFAULT_CONTAINER_READER)
}

func (this *mmal) WriterComponent() (hw.MMALComponent, error) {
	return this.ComponentWithName(rpi.MMAL_COMPONENT_DEFAULT_CONTAINER_WRITER)
}

////////////////////////////////////////////////////////////////////////////////
// CONNECTIONS

func (this *mmal) Connect(input, output hw.MMALPort, flags hw.MMALPortConnectionFlags) (hw.MMALPortConnection, error) {
	this.log.Debug2("<sys.hw.mmal>Connect{ input=%v output=%v flags=%v }", input, output, flags)

	var conn_ rpi.MMAL_PortConnection

	if input_, ok := input.(*port); ok == false {
		return nil, gopi.ErrBadParameter
	} else if output_, ok := output.(*port); ok == false {
		return nil, gopi.ErrBadParameter
	} else if err := rpi.MMALPortConnectionCreate(&conn_, input_.handle, output_.handle, flags); err != nil {
		return nil, err
	} else {
		// Append connection
		conn := &connection{
			handle: conn_,
			log:    this.log,
			input:  input,
			output: output,
		}
		this.connections = append(this.connections, conn)
		return conn, nil
	}
}

func (this *mmal) Disconnect(conn hw.MMALPortConnection) error {
	this.log.Debug2("<sys.hw.mmal>Disconnect{ conn=%v }", conn)
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *mmal) NewPort(handle rpi.MMAL_PortHandle) *port {
	var pool rpi.MMAL_Pool
	var err error

	// Create the pool & queue if it's an input or output port
	if rpi.MMALPortType(handle) == rpi.MMAL_PORT_TYPE_INPUT || rpi.MMALPortType(handle) == rpi.MMAL_PORT_TYPE_OUTPUT {
		pool, err = rpi.MMALPortPoolCreate(handle, 0, 0)
		if err != nil {
			this.log.Error("<sys.hw.mmal>NewPort: %v", err)
			return nil
		}
	}

	// Return port
	return &port{
		handle: handle,
		pool:   pool,
		log:    this.log,
	}
}
