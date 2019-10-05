/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Example to capture from the camera using the MMAL library
package main

import (
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
	tablewriter "github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

type CaptureApp struct {
	log      gopi.Logger
	mmal     hw.MMAL
	camera   hw.MMALComponent
	info     hw.MMALComponent
	renderer hw.MMALComponent
	encoder  hw.MMALComponent
}

////////////////////////////////////////////////////////////////////////////////

const (
	CAMERA_PREVIEW_OUTPORT = 0
	CAMERA_CAPTURE_OUTPORT = 2
	PREVIEW_INPORT         = 0
	ENCODER_INPORT         = 0
	ENCODER_OUTPORT        = 0
)

////////////////////////////////////////////////////////////////////////////////

func NewApp(mmal hw.MMAL, log gopi.Logger) (*CaptureApp, error) {
	this := new(CaptureApp)
	this.log = log
	this.mmal = mmal

	// Obtain components
	if info, err := this.mmal.CameraInfoComponent(); err != nil {
		return nil, err
	} else {
		this.info = info
	}
	if camera, err := this.mmal.CameraComponent(); err != nil {
		return nil, err
	} else {
		this.camera = camera
	}
	if renderer, err := this.mmal.VideoRendererComponent(); err != nil {
		return nil, err
	} else {
		this.renderer = renderer
	}
	if encoder, err := this.mmal.ImageEncoderComponent(); err != nil {
		return nil, err
	} else {
		this.encoder = encoder
	}

	return this, nil
}

////////////////////////////////////////////////////////////////////////////////

func (this *CaptureApp) disablePort(port hw.MMALPort) error {
	if port.Enabled() {
		if err := port.SetEnabled(false); err != nil {
			return err
		}
	}
	// Success
	return nil
}

func (this *CaptureApp) disableComponent(component hw.MMALComponent) error {
	if component.Enabled() {
		if err := component.SetEnabled(false); err != nil {
			return err
		}
	}
	// Success
	return nil
}

func (this *CaptureApp) Setup(flags *gopi.Flags) error {
	// Disable components and ports to start
	if err := this.disablePort(this.renderer.Inputs()[PREVIEW_INPORT]); err != nil {
		return err
	}
	if err := this.disablePort(this.encoder.Inputs()[ENCODER_INPORT]); err != nil {
		return err
	}
	if err := this.disablePort(this.camera.Outputs()[CAMERA_PREVIEW_OUTPORT]); err != nil {
		return err
	}
	if err := this.disablePort(this.camera.Outputs()[CAMERA_CAPTURE_OUTPORT]); err != nil {
		return err
	}
	if err := this.disablePort(this.encoder.Outputs()[ENCODER_OUTPORT]); err != nil {
		return err
	}

	// Connect ports
	if conn, err := this.mmal.Connect(this.renderer.Inputs()[PREVIEW_INPORT], this.camera.Outputs()[CAMERA_PREVIEW_OUTPORT], hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
		return err
	} else if err := conn.SetEnabled(true); err != nil {
		return err
	}
	if conn, err := this.mmal.Connect(this.encoder.Inputs()[ENCODER_INPORT], this.camera.Outputs()[CAMERA_CAPTURE_OUTPORT], hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
		return err
	} else if err := conn.SetEnabled(true); err != nil {
		return err
	}

	// Set output port to input parameter and set encoding
	this.encoder.Outputs()[ENCODER_OUTPORT].VideoFormat().SetEncoding(hw.MMAL_ENCODING_JPEG)
	this.encoder.Outputs()[ENCODER_OUTPORT].SetJPEGQFactor(5)
	if err := this.encoder.Outputs()[ENCODER_OUTPORT].CommitFormatChange(); err != nil {
		return err
	}

	// Enable the ports
	if err := this.encoder.Outputs()[ENCODER_OUTPORT].SetEnabled(true); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *CaptureApp) Save() error {
	return gopi.ErrNotImplemented
}

func (this *CaptureApp) CameraInfo() error {
	table := tablewriter.NewWriter(os.Stdout)
	if info, err := this.info.Control().CameraInfo(); err != nil {
		return err
	} else {
		table.SetHeader([]string{"Device", "Name", "Size", "Lens Present"})
		for _, camera := range info.Cameras() {
			id := fmt.Sprintf("Camera [%v]", camera.Id())
			x, y := camera.Size()
			table.Append([]string{id, camera.Name(), fmt.Sprintf("%vx%v", x, y), fmt.Sprint(camera.LensPresent())})
		}
	}
	table.Render()
	return nil
}

func (this *CaptureApp) Capture() error {
	if err := this.camera.Outputs()[CAMERA_CAPTURE_OUTPORT].SetCapture(true); err != nil {
		return err
	}

	// Get filename
	if writer, err := os.Create("encoded_image.jpg"); err != nil {
		return err
	} else {
		defer writer.Close()
		eoe := false

		// Report start of loop
		fmt.Println("Encoding to:", writer.Name())
		port_out := this.encoder.Outputs()[ENCODER_OUTPORT]

		// Feed input port and accept output
		for {
			// Get an empty buffer on from output pool, block until we get one, then send it
			// to the port so that it can be used for filling the result of the encode
			this.log.Debug("GetEmptyBufferOnPort")
			if buffer, err := this.encoder.GetEmptyBufferOnPort(port_out, true); err != nil {
				return err
			} else if err := port_out.Send(buffer); err != nil {
				return err
			}

			// Get a full buffer on the output port, block until we get one,
			// and write out to file
			this.log.Debug("GetFullBufferOnPort")
			if buffer, err := this.encoder.GetFullBufferOnPort(port_out, true); err != nil {
				return err
			} else if buffer != nil {
				if _, err := writer.Write(buffer.Data()); err != nil {
					return err
				}
				eoe = buffer.Flags()&hw.MMAL_BUFFER_FLAG_EOS != 0
				if err := port_out.Release(buffer); err != nil {
					return err
				}
				fmt.Println(eoe, buffer)
			}
		}

		// Flush output port
		if err := port_out.Flush(); err != nil {
			return err
		}

		// Disable the ports
		if err := port_out.SetEnabled(false); err != nil {
			return err
		}

		// Return success
		return nil
	}

}
