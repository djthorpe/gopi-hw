/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// MMAL example to encode an image
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unsafe"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/metrics"
	_ "github.com/djthorpe/gopi-hw/sys/mmal"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func ByteArray(ptr uintptr, len int) []byte {
	var array []byte
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&array)))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = ptr
	return array
}

// CreateRGBImage returns uncompressed RGBA image with red, green and blue stripes
// every 100 pixels
func CreateRGBImage(width, height uint32) []byte {
	// Create the image
	image := make([]uint32, width*height)
	i := 0
	for y := uint32(0); y < height; y++ {
		for x := uint32(0); x < width; x++ {
			// ALPHA,BLUE,GREEN,RED
			switch {
			case x%300 < 100:
				image[i] = 0xFF00FF00
			case x%300 < 200:
				image[i] = 0xFFFF0000
			default:
				image[i] = 0xFF0000FF
			}
			i++
		}
	}
	return ByteArray(uintptr(unsafe.Pointer(&image[0])), len(image)*4)
}

func MMALEncodeTest(encoder hw.MMALComponent, format hw.MMALEncodingType, width, height uint32) error {
	port_in := encoder.Inputs()[0]
	port_out := encoder.Outputs()[0]
	if port_in.Enabled() {
		if err := port_in.SetEnabled(false); err != nil {
			return err
		}
	}
	if port_out.Enabled() {
		if err := port_out.SetEnabled(false); err != nil {
			return err
		}
	}

	// Set input port to uncompressed RGBA
	port_in.VideoFormat().SetEncoding(hw.MMAL_ENCODING_RGBA)
	port_in.VideoFormat().SetWidthHeight(width, height)
	port_in.VideoFormat().SetCrop(hw.MMALRect{0, 0, width, height})
	if err := port_in.CommitFormatChange(); err != nil {
		return err
	}

	// Set output port to input parameter and set encoding
	if err := port_out.CopyFormat(port_in.VideoFormat()); err != nil {
		return err
	}
	port_out.VideoFormat().SetEncoding(format)
	if err := port_out.CommitFormatChange(); err != nil {
		return err
	}

	// Set JPEG Quality factor
	if err := port_out.SetJPEGQFactor(5); err != nil {
		return err
	}

	// Enable the ports
	if err := port_in.SetEnabled(true); err != nil {
		return err
	}
	if err := port_out.SetEnabled(true); err != nil {
		return err
	}

	// Get filename
	ext := strings.ToLower(strings.TrimSpace(strings.Trim(fmt.Sprint(format), "'")))
	reader := bytes.NewReader(CreateRGBImage(width, height))
	if writer, err := os.Create("encoded_image." + ext); err != nil {
		return err
	} else {
		defer writer.Close()
		eof := false
		eoe := false

		// Report start of loop
		fmt.Println("Encoding to:", writer.Name())

		// Feed input port and accept output
		for {
			// Get an empty buffer on from output pool, block until we get one, then send it
			// to the port so that it can be used for filling the result of the encode
			if buffer, err := encoder.GetEmptyBufferOnPort(port_out, true); err != nil {
				return err
			} else if err := port_out.Send(buffer); err != nil {
				return err
			}

			// Get an empty buffer on input port, block until we get one, then fill it
			// with uncompressed image data and send it
			if eof {
				// Do nothing when all bytes have been sent
			} else if buffer, err := encoder.GetEmptyBufferOnPort(port_in, true); err != nil {
				return err
			} else if _, err := buffer.Fill(reader); err != nil && err != io.EOF {
				return err
			} else if err := port_in.Send(buffer); err != nil {
				return err
			} else if buffer.Flags()&hw.MMAL_BUFFER_FLAG_EOS != 0 {
				eof = true
			}

			// Get a full buffer on the output port, block until we get one,
			// and write out to file
			if buffer, err := encoder.GetFullBufferOnPort(port_out, true); err != nil {
				return err
			} else if buffer != nil {
				if _, err := writer.Write(buffer.Data()); err != nil {
					return err
				}
				eoe = buffer.Flags()&hw.MMAL_BUFFER_FLAG_EOS != 0
				if err := port_out.Release(buffer); err != nil {
					return err
				}
			}

			// Check for end of input and output, break out of loop
			// when both input and outputs are finished
			if eof && eoe {
				break
			}
		}
	}

	// Flush output port
	if err := port_out.Flush(); err != nil {
		return err
	}

	// Disable the ports
	if err := port_in.SetEnabled(false); err != nil {
		return err
	}
	if err := port_out.SetEnabled(false); err != nil {
		return err
	}

	// Success
	return nil
}

// Create JPEG, PNG and BMP encoded images
func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	if mmal := app.ModuleInstance("hw/mmal").(hw.MMAL); mmal == nil {
		return fmt.Errorf("Missing MMAL module")
	} else if image_encoder, err := mmal.ImageEncoderComponent(); err != nil {
		return err
	} else if err := MMALEncodeTest(image_encoder, hw.MMAL_ENCODING_JPEG, 1920, 1080); err != nil {
		return err
	} else if err := MMALEncodeTest(image_encoder, hw.MMAL_ENCODING_PNG, 1920, 1080); err != nil {
		return err
	} else if err := MMALEncodeTest(image_encoder, hw.MMAL_ENCODING_BMP, 1920, 1080); err != nil {
		return err
	}

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the MMAL instance
	config := gopi.NewAppConfig("hw/mmal")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main))
}
