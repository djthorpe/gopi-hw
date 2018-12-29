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
	"github.com/djthorpe/gopi-hw/rpi"
)

// Fill the buffer and return number of bytes read and error
func (this *buffer) Fill(handle io.Reader) (uint32, error) {
	n, err := handle.Read(rpi.MMALBufferBytes(this.handle))
	this.log.Debug2("<sys.hw.mmal.buffer>Fill{ n=%v err=%v }", n, err)
	rpi.MMALBufferSetLength(this.handle, uint32(n))
	return uint32(n), err
}
