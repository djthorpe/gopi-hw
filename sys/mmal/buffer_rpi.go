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
	"io"

	// Frameworks
	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

// Fill the buffer and return number of bytes read and error
func (this *buffer) Fill(handle io.Reader) (uint32, error) {
	n, err := handle.Read(rpi.MMALBufferBytes(this.handle))
	rpi.MMALBufferSetLength(this.handle, uint32(n))
	// Try again to read to determine if EOF
	if err == nil {
		_, err = handle.Read([]byte{})
	}
	this.log.Debug2("<sys.hw.mmal.buffer>Fill{ buffer=%v n=%v err=%v }", rpi.MMALBufferString(this.handle), n, err)
	if err == io.EOF {
		// If EOF then set EOS flag
		rpi.MMALBufferSetFlags(this.handle, rpi.MMALBufferFlags(this.handle)|hw.MMAL_BUFFER_FLAG_EOS)
	}
	// check buffer is not empty
	return uint32(n), err
}

func (this *buffer) Flags() hw.MMALBufferFlag {
	return rpi.MMALBufferFlags(this.handle)
}

func (this *buffer) Data() []byte {
	return rpi.MMALBufferData(this.handle)
}

// Acquire buffer
func (this *buffer) Acquire() error {
	this.log.Debug2("<sys.hw.mmal.buffer>Acquire{ buffer=%v }", rpi.MMALBufferString(this.handle))
	return rpi.MMALBufferAcquire(this.handle)
}

// Release buffer
func (this *buffer) Release() error {
	this.log.Debug2("<sys.hw.mmal.buffer>Release{ buffer=%v }", rpi.MMALBufferString(this.handle))
	return rpi.MMALBufferRelease(this.handle)
}

// Reset buffer
func (this *buffer) Reset() error {
	this.log.Debug2("<sys.hw.mmal.buffer>Reset{ buffer=%v }", rpi.MMALBufferString(this.handle))
	return rpi.MMALBufferReset(this.handle)
}

// Stringify buffer
func (this *buffer) String() string {
	return rpi.MMALBufferString(this.handle)
}
