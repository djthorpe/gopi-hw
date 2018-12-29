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
	"sync"

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
	log      gopi.Logger
	handle   rpi.MMAL_ComponentHandle
	control  *port
	input    []*port
	output   []*port
	clock    []*port
	port_map map[rpi.MMAL_PortHandle]uint
}

type port struct {
	log    gopi.Logger
	handle rpi.MMAL_PortHandle
	pool   rpi.MMAL_Pool
	queue  rpi.MMAL_Queue

	sync.WaitGroup
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

type buffer struct {
	log    gopi.Logger
	handle rpi.MMAL_Buffer
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
		c.control = this.NewPort(c, rpi.MMALComponentControlPort(handle))
		// Input ports
		c.input = make([]*port, rpi.MMALComponentInputPortNum(handle))
		for i := range c.input {
			c.input[i] = this.NewPort(c, rpi.MMALComponentInputPortAtIndex(handle, uint(i)))
		}
		// Output ports
		c.output = make([]*port, rpi.MMALComponentOutputPortNum(handle))
		for i := range c.output {
			c.output[i] = this.NewPort(c, rpi.MMALComponentOutputPortAtIndex(handle, uint(i)))
		}
		// Clock ports
		c.clock = make([]*port, rpi.MMALComponentClockPortNum(handle))
		for i := range c.clock {
			c.clock[i] = this.NewPort(c, rpi.MMALComponentClockPortAtIndex(handle, uint(i)))
		}

		// Map port handles to port index
		if port_num := rpi.MMALComponentPortNum(handle); port_num == 0 {
			// There should be at least one port
			return nil, gopi.ErrAppError
		} else {
			c.port_map = make(map[rpi.MMAL_PortHandle]uint, int(port_num))
			for i := uint(0); i < port_num; i++ {
				if port_handle := rpi.MMALComponentPortAtIndex(handle, i); port_handle == nil {
					// Port handle should not be nil
					return nil, gopi.ErrAppError
				} else if _, exists := c.port_map[port_handle]; exists {
					// Each port should only exist once
					return nil, gopi.ErrAppError
				} else {
					c.port_map[port_handle] = rpi.MMALPortIndex(port_handle)
				}
			}
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
		// Make a connection object
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

func (this *mmal) NewPort(c *component, handle rpi.MMAL_PortHandle) *port {
	var pool rpi.MMAL_Pool
	var queue rpi.MMAL_Queue
	var err error

	// Create the pool & queue if it's an input or output port
	if rpi.MMALPortType(handle) == rpi.MMAL_PORT_TYPE_INPUT || rpi.MMALPortType(handle) == rpi.MMAL_PORT_TYPE_OUTPUT {
		if pool, err = rpi.MMALPortPoolCreate(handle, 0, 0); err != nil {
			this.log.Error("<sys.hw.mmal>NewPort: %v", err)
			return nil
		} else if queue = rpi.MMALQueueCreate(); queue == nil {
			this.log.Error("<sys.hw.mmal>NewPort: queue == nil")
			rpi.MMALPortPoolDestroy(handle, pool)
			return nil
		}
	}

	// Register the port callback
	rpi.MMALPortRegisterCallback(handle, func(port rpi.MMAL_PortHandle, buffer rpi.MMAL_Buffer) {
		if rpi.MMALPortType(port) == rpi.MMAL_PORT_TYPE_CONTROL {
			// Callback from a control port. Error events will be received there
			fmt.Printf("CALLBACK CONTROL PORT BUFFER: port=%v, buffer=%v\n", rpi.MMALPortName(port), rpi.MMALBufferString(buffer))
		} else if rpi.MMALPortType(port) == rpi.MMAL_PORT_TYPE_INPUT {
			// Callback from an input port. Buffer is released
			if err := rpi.MMALBufferRelease(buffer); err != nil {
				fmt.Printf("CALLBACK INPUT PORT BUFFER: port=%v, buffer=%v: %v\n", rpi.MMALPortName(port), rpi.MMALBufferString(buffer), err)
			}
		} else if rpi.MMALPortType(port) == rpi.MMAL_PORT_TYPE_OUTPUT {
			// Callback from an output port. Buffer is queued for the next component
			fmt.Printf("CALLBACK OUTPUT PORT BUFFER: port=%v, buffer=%v\n", rpi.MMALPortName(port), rpi.MMALBufferString(buffer))
		} else {
			fmt.Printf("CALLBACK OTHER PORT CALLBACK: component=%v port=%v, buffer=%v\n", c.Name(), rpi.MMALPortName(port), rpi.MMALBufferString(buffer))
		}
	})

	// Return port
	return &port{
		handle: handle,
		pool:   pool,
		queue:  queue,
		log:    this.log,
	}
}
