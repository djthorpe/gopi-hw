/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package pwm

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	fifo       string
	log        gopi.Logger
	fh         *os.File
	period     time.Duration
	samples    uint
	min_period time.Duration
	max_period time.Duration
}

var (
	regexpFrequency = regexp.MustCompile("^([0-9]+\\.?[0-9]*)\\s*(hz|khz|mhz)$")
)

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
	} else if stat.Mode().IsDir() {
		return nil, fmt.Errorf("Not a file: %v", this.fifo)
	} else if fh, err := os.OpenFile(this.fifo, os.O_WRONLY|os.O_APPEND, 0755); err != nil {
		return nil, err
	} else {
		this.fh = fh
	}

	// Read parameters
	if _, err := os.Stat(config.Exec); os.IsNotExist(err) == true {
		return nil, err
	} else if output, _ := exec.Command(config.Exec, "-D").CombinedOutput(); len(output) == 0 {
		return nil, fmt.Errorf("No output: %v", config.Exec)
	} else if err := this.set_values_from_output(output); err != nil {
		return nil, err
	}

	// success
	return this, nil
}

// Close PiBlaster object
func (this *piblaster) Close() error {
	this.log.Debug("<sys.hw.PWM.PiBlaster>Close{ fifo='%v' }", this.fifo)

	// Release resources
	if err := this.fh.Close(); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *piblaster) String() string {
	return fmt.Sprintf("<sys.hw.PWM.PiBlaster>{ fifo='%v' period=%v min_period=%v max_period=%v samples=%v }", this.fifo, this.period,this.min_period,this.max_period,this.samples)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *piblaster) Frequency(gopi.GPIOPin) (float32, error) {
	return 0, gopi.ErrNotImplemented
}

func (this *piblaster) SetFrequency(float32, gopi.GPIOPin) error {
	return gopi.ErrNotImplemented
}

func (this *piblaster) DutyCycle(gopi.GPIOPin) (float32, error) {
	return 0, gopi.ErrNotImplemented
}

func (this *piblaster) SetDutyCycle(float32, gopi.GPIOPin) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *piblaster) set_values_from_output(output []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		if tuple := strings.Split(scanner.Text(), ":"); len(tuple) != 2 {
			continue
		} else if err := this.set_value(strings.ToLower(strings.TrimSpace(tuple[0])), strings.TrimSpace(tuple[1])); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (this *piblaster) set_value(name, value string) error {
	this.log.Debug2("<sys.hw.PWM.PiBlaster>set_value{ name='%v' value='%v' }", name, value)
	switch {
	case strings.HasSuffix(name, "frequency"):
		if period, err := ParseFrequency(value); err != nil {
			return err
		} else {
			this.period = period
		}
	case strings.HasSuffix(name, "steps"):
		if samples, err := strconv.ParseUint(value,10, 64); err != nil {
			return err
		} else {
			this.samples = uint(samples)
		}
	case strings.HasPrefix(name, "maximum period"):
		if max, err := time.ParseDuration(value); err != nil {
			return err
		} else {
			this.max_period = max
		}
	case strings.HasPrefix(name, "minimum period"):
		if min, err := time.ParseDuration(value); err != nil {
			return err
		} else {
			this.min_period = min
		}
	}
	return nil
}

func ParseFrequency(value string) (time.Duration, error) {
	if parts := regexpFrequency.FindStringSubmatch(strings.ToLower(value)); len(parts) != 3 {
		return 0, gopi.ErrBadParameter
	} else if num, err := strconv.ParseFloat(parts[1], 64); err != nil {
		return 0, err
	} else {
		switch parts[2] {
		case "hz":
			return time.Nanosecond * time.Duration((1E9 / num)), nil
		default:
			return 0, gopi.ErrBadParameter
		}
	}
}
