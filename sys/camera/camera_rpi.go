// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package camera

import (
	"fmt"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Camera struct {
	MMAL   hw.MMAL
	Camera uint32
}

type camera struct {
	log         gopi.Logger
	mmal        hw.MMAL
	camera_id   uint32
	camera_info hw.MMALCamera

	// Components
	camera   hw.MMALComponent
	renderer hw.MMALComponent
	encoder  hw.MMALComponent
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Camera) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<media.camera>Open{ camera_id=%v mmal=%v }", config.Camera, config.MMAL)

	this := new(camera)
	this.log = log
	this.mmal = config.MMAL
	this.camera_id = config.Camera

	// Check incoming parameters
	if this.mmal == nil {
		return nil, gopi.ErrBadParameter
	}

	// Obtain information about attached cameras
	if info, err := this.mmal.CameraInfoComponent(); err != nil {
		return nil, err
	} else if camera_info, err := info.Control().CameraInfo(); err != nil {
		return nil, err
	} else {
		for _, camera := range camera_info.Cameras() {
			if camera.Id() == this.camera_id {
				this.camera_info = camera
			}
		}
	}

	// Get camera, render and encoder components
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

	// If no camera info set, then return error
	if this.camera_info == nil {
		return nil, hw.ErrInvalidCameraId
	} else if err := this.CameraControlPort().SetCameraNum(this.camera_id); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *camera) Close() error {
	this.log.Debug("<media.camera>Close{ camera_id=%v }", this.camera_id)

	// Release resources
	this.mmal = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *camera) String() string {
	return fmt.Sprintf("<media.camera>{ camera_id=%v info=%v }", this.camera_id, this.camera_info)
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *camera) CameraId() uint32 {
	return this.camera_id
}

func (this *camera) CameraModel() string {
	return this.camera_info.Name()
}

func (this *camera) CameraFrameSize() gopi.Size {
	w, h := this.camera_info.Size()
	return gopi.Size{float32(w), float32(h)}
}

////////////////////////////////////////////////////////////////////////////////
// CAPTURE CONFIGURATION

func (this *camera) CameraConfig() (hw.CameraConfig, error) {
	var config hw.CameraConfig

	// Rotation values
	if rotation, err := this.CameraImageOutPort().Rotation(); err != nil {
		return hw.CameraConfig{}, err
	} else {
		config.ImageRotation = rotation
	}
	if rotation, err := this.CameraPreviewOutPort().Rotation(); err != nil {
		return hw.CameraConfig{}, err
	} else {
		config.PreviewRotation = rotation
	}
	if rotation, err := this.CameraVideoOutPort().Rotation(); err != nil {
		return hw.CameraConfig{}, err
	} else {
		config.VideoRotation = rotation
	}

	// Encoder Image Encoding format
	config.ImageFormatEncoding, config.ImageFormatEncodingVariant = this.EncoderOutPort().VideoFormat().Encoding()

	// Camera Image Frame size
	if format := this.CameraImageOutPort().VideoFormat(); format != nil {
		w, h := format.WidthHeight()
		config.ImageFrameSize = gopi.Size{float32(w), float32(h)}
	}

	// Preview Frame size
	if format := this.PreviewInPort().VideoFormat(); format != nil {
		w, h := format.WidthHeight()
		config.PreviewFrameSize = gopi.Size{float32(w), float32(h)}
	}

	// JPEG Quality
	if jpegquality, err := this.EncoderOutPort().JPEGQFactor(); err != nil {
		return hw.CameraConfig{}, err
	} else {
		config.ImageJPEGQuality = jpegquality
	}

	// Return success
	return config, nil
}

func (this *camera) SetCameraConfig(config hw.CameraConfig) error {
	if config.HasFlags(hw.FLAG_IMAGE_ROTATION) {
		if err := this.CameraImageOutPort().SetRotation(config.ImageRotation); err != nil {
			return err
		}
	}
	if config.HasFlags(hw.FLAG_PREVIEW_ROTATION) {
		if err := this.CameraPreviewOutPort().SetRotation(config.PreviewRotation); err != nil {
			return err
		}
	}
	if config.HasFlags(hw.FLAG_VIDEO_ROTATION) {
		if err := this.CameraVideoOutPort().SetRotation(config.VideoRotation); err != nil {
			return err
		}
	}
	if config.HasFlags(hw.FLAG_IMAGE_ENCODING_FORMAT) {
		this.CameraImageOutPort().VideoFormat().SetEncodingVariant(config.ImageFormatEncoding, config.ImageFormatEncodingVariant)
		if err := this.CameraImageOutPort().CommitFormatChange(); err != nil {
			return err
		}
	}
	if config.HasFlags(hw.FLAG_IMAGE_FRAMESIZE) {
		w, h := uint32(config.ImageFrameSize.W), uint32(config.ImageFrameSize.H)
		format := this.CameraImageOutPort().VideoFormat()
		format.SetWidthHeight(w, h)
		format.SetCrop(hw.MMALRect{0, 0, w, h})
		if err := this.CameraImageOutPort().CommitFormatChange(); err != nil {
			return err
		}
		this.EncoderInPort().CopyFormat(format)
	}
	if config.HasFlags(hw.FLAG_PREVIEW_FRAMESIZE) {
		w, h := uint32(config.PreviewFrameSize.W), uint32(config.PreviewFrameSize.H)
		format := this.PreviewInPort().VideoFormat()
		format.SetWidthHeight(w, h)
		format.SetCrop(hw.MMALRect{0, 0, w, h})
		if err := this.PreviewInPort().CommitFormatChange(); err != nil {
			return err
		}
	}
	if config.HasFlags(hw.FLAG_IMAGE_ENCODING_JPEGQUALITY) {
		if err := this.CameraImageOutPort().SetJPEGQFactor(config.ImageJPEGQuality); err != nil {
			return err
		}
	}
	// Success
	return nil
}

func (this *camera) Preview() error {
	this.log.Debug2("<media.camera.Preview>{}")

	// Connect camera to video renderer port
	if conn, err := this.mmal.Connect(this.PreviewInPort(), this.CameraPreviewOutPort(), hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
		return err
	} else if err := conn.SetEnabled(true); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *camera) ImageCapture() error {
	this.log.Debug2("<media.camera.ImageCapture>{}")

	// Connect camera to the image capture port
	if conn, err := this.mmal.Connect(this.EncoderInPort(), this.CameraImageOutPort(), hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
		return err
	} else if err := conn.SetEnabled(true); err != nil {
		return err
	}

	// Enable the ports
	if err := this.EncoderOutPort().SetEnabled(true); err != nil {
		return err
	}
	defer this.EncoderOutPort().SetEnabled(false)

	// Start Capture
	if err := this.CameraImageOutPort().SetCapture(true); err != nil {
		return err
	}
	defer this.CameraImageOutPort().SetCapture(false)

	port_out := this.EncoderOutPort()
FOR_LOOP:
	for {
		// Get an empty buffer on from output pool, block until we get one, then send it
		// to the port so that it can be used for filling the result of the encode
		this.log.Debug("media.camera.ImageCapture: GetEmptyBufferOnPort")
		if buffer, err := this.encoder.GetEmptyBufferOnPort(port_out, true); err != nil {
			return err
		} else if err := port_out.Send(buffer); err != nil {
			return err
		}

		// Get a full buffer on the output port, block until we get one,
		// and write out to file
		this.log.Debug("media.camera.ImageCapture: GetFullBufferOnPort")
		if buffer, err := this.encoder.GetFullBufferOnPort(port_out, true); err != nil {
			return err
		} else if buffer != nil {
			eos := buffer.Flags() & hw.MMAL_BUFFER_FLAG_EOS
			fmt.Printf("TODO: Copy %v bytes out\n", len(buffer.Data()))
			this.log.Debug("<media.camera.ImageCapture>{ buffer=%v eos=%v }", buffer, eos)
			if err := port_out.Release(buffer); err != nil {
				return err
			}
			if eos == hw.MMAL_BUFFER_FLAG_EOS {
				break FOR_LOOP
			}
		}
	}

	// Success
	return nil
}
