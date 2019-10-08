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
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
	errors "github.com/djthorpe/gopi/util/errors"
	event "github.com/djthorpe/gopi/util/event"
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
	camera        hw.MMALComponent
	renderer      hw.MMALComponent
	image_encoder hw.MMALComponent
	video_encoder hw.MMALComponent

	// Connections
	preview_connection hw.MMALPortConnection
	image_connection   hw.MMALPortConnection

	// Embeds
	event.Publisher
	sync.Mutex
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
		this.image_encoder = encoder
	}
	if encoder, err := this.mmal.VideoEncoderComponent(); err != nil {
		return nil, err
	} else {
		this.video_encoder = encoder
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

	// Lock
	this.Lock()
	defer this.Unlock()

	// Close subscriber channels
	this.Publisher.Close()

	// Disconnect connections
	this.log.Debug("disable connections")
	if this.image_connection != nil {
		if err := this.image_connection.SetEnabled(false); err != nil {
			return err
		} else if err := this.mmal.Disconnect(this.image_connection); err != nil {
			return err
		}
	}
	if this.preview_connection != nil {
		if err := this.preview_connection.SetEnabled(false); err != nil {
			return err
		} else if err := this.mmal.Disconnect(this.preview_connection); err != nil {
			return err
		}
	}

	// Disable ports
	this.log.Debug("disable ports")
	errs := new(errors.CompoundError)
	errs.Add(this.enable_port(this.CameraPreviewOutPort(), false))
	errs.Add(this.enable_port(this.CameraImageOutPort(), false))
	errs.Add(this.enable_port(this.CameraVideoOutPort(), false))
	errs.Add(this.enable_port(this.PreviewInPort(), false))
	errs.Add(this.enable_port(this.ImageEncoderInPort(), false))
	errs.Add(this.enable_port(this.ImageEncoderOutPort(), false))
	errs.Add(this.enable_port(this.VideoEncoderInPort(), false))
	errs.Add(this.enable_port(this.VideoEncoderOutPort(), false))

	// Disable components
	this.log.Debug("disable components")
	errs.Add(this.camera.SetEnabled(false))
	errs.Add(this.renderer.SetEnabled(false))
	errs.Add(this.image_encoder.SetEnabled(false))
	errs.Add(this.video_encoder.SetEnabled(false))

	// Release resources
	this.camera = nil
	this.renderer = nil
	this.image_encoder = nil
	this.video_encoder = nil
	this.preview_connection = nil
	this.image_connection = nil
	this.mmal = nil

	// Return success
	return errs.ErrorOrSelf()
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

func (this *camera) SupportedImageEncodings() []hw.MMALEncodingType {
	if encodings, err := this.ImageEncoderOutPort().SupportedEncodings(); err != nil {
		this.log.Warn("SupportedImageEncodings: %v", err)
		return nil
	} else {
		return encodings
	}
}

func (this *camera) SupportedVideoEncodings() []hw.MMALEncodingType {
	if encodings, err := this.VideoEncoderOutPort().SupportedEncodings(); err != nil {
		this.log.Warn("SupportedVideoEncodings: %v", err)
		return nil
	} else {
		return encodings
	}
}

////////////////////////////////////////////////////////////////////////////////
// CAPTURE CONFIGURATION

func (this *camera) CameraConfig() (hw.CameraConfig, error) {
	this.log.Debug2("<media.camera.CameraConfig>{ }")
	this.Lock()
	defer this.Unlock()

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
	config.ImageFormatEncoding, config.ImageFormatEncodingVariant = this.ImageEncoderOutPort().VideoFormat().Encoding()

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
	if jpegquality, err := this.ImageEncoderOutPort().JPEGQFactor(); err != nil {
		return hw.CameraConfig{}, err
	} else {
		config.ImageJPEGQuality = jpegquality
	}

	// Return success
	return config, nil
}

func (this *camera) SetCameraConfig(config hw.CameraConfig) error {
	this.log.Debug2("<media.camera.SetCameraConfig>{ config=%v }", config)
	this.Lock()
	defer this.Unlock()

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
		this.ImageEncoderInPort().CopyFormat(format)
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

func (this *camera) Preview(enable bool) error {
	this.log.Debug2("<media.camera.Preview>{ enable=%v }", enable)
	this.Lock()
	defer this.Unlock()

	if enable {
		if this.preview_connection != nil {
			return gopi.ErrOutOfOrder
		} else if preview_connection, err := this.mmal.Connect(this.PreviewInPort(), this.CameraPreviewOutPort(), hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
			return err
		} else if err := preview_connection.SetEnabled(true); err != nil {
			return err
		} else {
			this.preview_connection = preview_connection
		}
	} else {
		if this.preview_connection == nil {
			return gopi.ErrOutOfOrder
		} else if err := this.preview_connection.SetEnabled(false); err != nil {
			return err
		} else if err := this.mmal.Disconnect(this.preview_connection); err != nil {
			return err
		} else {
			this.preview_connection = nil
		}
	}

	// Success
	return nil
}

func (this *camera) ImageCapture() error {
	this.log.Debug2("<media.camera.ImageCapture>{}")
	this.Lock()
	defer this.Unlock()

	// Connect camera to the image capture port
	if this.image_connection == nil {
		if image_connection, err := this.mmal.Connect(this.ImageEncoderInPort(), this.CameraImageOutPort(), hw.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT|hw.MMAL_CONNECTION_FLAG_TUNNELLING); err != nil {
			return err
		} else if err := image_connection.SetEnabled(true); err != nil {
			return err
		} else {
			this.image_connection = image_connection
		}
	}

	// Enable the ports
	this.log.Debug2("media.camera.ImageCapture: EncoderOutPort SetEnabled")
	if err := this.ImageEncoderOutPort().SetEnabled(true); err != nil {
		return err
	}
	defer this.ImageEncoderOutPort().SetEnabled(false)

	// Start Capture
	this.log.Debug2("media.camera.ImageCapture: CameraImageOutPort SetCapture")
	if err := this.CameraImageOutPort().SetCapture(true); err != nil {
		return err
	}
	defer this.CameraImageOutPort().SetCapture(false)

	port_out := this.ImageEncoderOutPort()
	start_of_stream := true
FOR_LOOP:
	for {
		// Get an empty buffer on from output pool, block until we get one, then send it
		// to the port so that it can be used for filling the result of the encode
		this.log.Debug2("media.camera.ImageCapture: GetEmptyBufferOnPort")
		if buffer, err := this.image_encoder.GetEmptyBufferOnPort(port_out, true); err != nil {
			return err
		} else if err := port_out.Send(buffer); err != nil {
			return err
		}

		// Get a full buffer on the output port, block until we get one,
		// and write out to file
		this.log.Debug2("media.camera.ImageCapture: GetFullBufferOnPort")
		if buffer, err := this.image_encoder.GetFullBufferOnPort(port_out, true); err != nil {
			return err
		} else if buffer != nil {
			end_of_stream := this.emit(start_of_stream, hw.FLAG_DATA_IMAGE, buffer)
			this.log.Debug2("<media.camera.ImageCapture>{ buffer=%v eos=%v }", buffer, end_of_stream)
			if err := port_out.Release(buffer); err != nil {
				return err
			}
			if end_of_stream {
				break FOR_LOOP
			} else {
				start_of_stream = false
			}
		}
	}

	// Flush the output port
	if err := port_out.Flush(); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *camera) emit(start_of_stream bool, flags hw.CameraDataFlag, buffer hw.MMALBuffer) bool {
	end_of_stream := (buffer.Flags() & hw.MMAL_BUFFER_FLAG_EOS) == hw.MMAL_BUFFER_FLAG_EOS
	if start_of_stream {
		flags |= hw.FLAG_DATA_STREAM_START
	}
	if end_of_stream {
		flags |= hw.FLAG_DATA_STREAM_END
	}
	this.Emit(NewEvent(this, buffer.Data(), flags))
	return end_of_stream
}

func (this *camera) enable_port(port hw.MMALPort, enable bool) error {
	this.log.Debug("<media.camera.enable_port>{ port=%v enable=%v }", port, enable)
	if enable != port.Enabled() {
		return port.SetEnabled(enable)
	} else {
		return nil
	}
}
