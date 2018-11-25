/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package pwm

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type PiBlaster struct {
	FIFO string
	Exec string
}

type piblaster struct {
	fifo string
	log  gopi.Logger
	fh   *os.File
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new PiBlaster object, returns error if not possible
func (config PiBlaster) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.PWM.PiBlaster>Open{ fifo='%v' exec='%v' }", config.FIFO, config.Exec)

	// create new GPIO driver
	this := new(piblaster)

	// Set logging & device
	this.log = log
	this.fifo = config.FIFO

	// Open FIFO
	if stat, err := os.Stat(this.fifo); os.IsNotExist(err) == true {
		return nil, err
	} else if stat.Mode() != os.ModeNamedPipe {
		return nil, fmt.Errorf("Not a named pipe: %v", this.fifo)
	} else if fh, err := os.OpenFile(this.fifo, os.O_WRONLY|os.O_APPEND, 0755); err != nil {
		return nil, err
	} else {
		this.fh = fh
	}

	// success
	return this, nil
}

// Close PiBlaster object
func (this *piblaster) Close() error {
	this.log.Debug("<sys.hw.PWM.PiBlaster>Close{ fifo='%v' }", this.fifo)

	// Release resources

	return nil
}
