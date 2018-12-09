// +build rpi

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
	"strconv"
	"strings"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type PiBlaster struct {
	GPIO gopi.GPIO
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
	pins       map[gopi.GPIOPin]float32
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new PiBlaster object, returns error if not possible
func (config PiBlaster) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.PWM.PiBlaster>Open{ fifo='%v' exec='%v' }", config.FIFO, config.Exec)

	// create new PWM driver
	this := new(piblaster)

	// Set logging & device
	this.log = log
	this.fifo = config.FIFO

	// Check for fifo
	if this.fifo == "" {
		return nil, gopi.ErrBadParameter
	}

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

	// Check parameters
	if this.period == 0 {
		return nil, fmt.Errorf("Unable to determine PWM frequency")
	}

	// Set up pins
	this.pins = make(map[gopi.GPIOPin]float32)

	// success
	return this, nil
}

// Close PiBlaster object
func (this *piblaster) Close() error {
	this.log.Debug("<sys.hw.PWM.PiBlaster>Close{ fifo='%v' }", this.fifo)

	// Release pins
	for _, pin := range this.all_pins() {
		if err := this.Release(pin); err != nil {
			return err
		}
	}

	// Release resources
	if err := this.fh.Close(); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *piblaster) String() string {
	return fmt.Sprintf("<sys.hw.PWM.PiBlaster>{ fifo='%v' pins=%v period=%v min_period=%v max_period=%v samples=%v }", this.fifo, this.all_pins(), this.period, this.min_period, this.max_period, this.samples)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *piblaster) Pins() []gopi.GPIOPin {
	return this.all_pins()
}

func (this *piblaster) Period(pin gopi.GPIOPin) (time.Duration, error) {
	this.log.Debug2("<sys.hw.PWM.PiBlaster>Period{ pin=%v }", pin)

	// If pin is not in list of pins, then set duty cycle of 0
	if _, exists := this.pins[pin]; exists {
		return this.period, nil
	} else if err := this.SetDutyCycle(0, pin); err != nil {
		return 0, err
	} else {
		return this.period, nil
	}
}

func (this *piblaster) SetPeriod(period time.Duration, pins ...gopi.GPIOPin) error {
	this.log.Debug2("<sys.hw.PWM.PiBlaster>SetPeriod{ period=%v pins=%v }", period, pins)

	if period != this.period {
		return fmt.Errorf("<sys.hw.PWM.PiBlaster>SetPeriod: Unable to set period other than %v on pins %v", this.period, pins)
	} else if len(pins) == 0 {
		return gopi.ErrBadParameter
	} else if err := this.SetDutyCycle(0, pins...); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *piblaster) DutyCycle(pin gopi.GPIOPin) (float32, error) {
	this.log.Debug2("<sys.hw.PWM.PiBlaster>DutyCycle{ pin=%v }", pin)

	if duty_cycle, exists := this.pins[pin]; exists {
		return duty_cycle, nil
	} else if err := this.SetDutyCycle(0, pin); err != nil {
		return 0, err
	} else {
		return 0, nil
	}
}

func (this *piblaster) SetDutyCycle(duty_cycle float32, pins ...gopi.GPIOPin) error {
	this.log.Debug2("<sys.hw.PWM.PiBlaster>SetDutyCycle{ duty_cycle=%v pins=%v }", duty_cycle, pins)

	if duty_cycle < 0.0 || duty_cycle > 1.0 {
		return gopi.ErrBadParameter
	} else if len(pins) == 0 {
		return gopi.ErrBadParameter
	} else {
		// Clamp duty cycle between minimum and maxiumum unless it's 0 or 1
		min_duty_cycle := float32(this.min_period.Nanoseconds()) / float32(this.period.Nanoseconds())
		max_duty_cycle := float32(this.max_period.Nanoseconds()) / float32(this.period.Nanoseconds())
		if duty_cycle != 0.0 && duty_cycle < min_duty_cycle {
			duty_cycle = min_duty_cycle
		}
		if duty_cycle != 1.0 && duty_cycle > max_duty_cycle {
			duty_cycle = max_duty_cycle
		}
		params := make([]string, len(pins))
		for i, pin := range pins {
			params[i] = fmt.Sprintf("%v=%v", uint(pin), duty_cycle)
		}
		this.log.Debug2("<sys.hw.PWM.PiBlaster>SetDutyCycle{ write=>\"%v\" }", strings.Join(params, " "))
		if _, err := fmt.Fprintf(this.fh, "%v\n", strings.Join(params, " ")); err != nil {
			return err
		} else {
			for _, pin := range pins {
				this.pins[pin] = duty_cycle
			}
			return nil
		}
	}
}

func (this *piblaster) Release(pin gopi.GPIOPin) error {
	this.log.Debug2("<sys.hw.PWM.PiBlaster>Release{ pin=%v }", pin)

	if _, exists := this.pins[pin]; exists == false {
		return gopi.ErrNotFound
	} else if _, err := fmt.Fprintf(this.fh, "release %v\n", uint(pin)); err != nil {
		return err
	} else {
		delete(this.pins, pin)
	}
	return nil
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
		if samples, err := strconv.ParseUint(value, 10, 64); err != nil {
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

// all_pins returns the enabled pins
func (this *piblaster) all_pins() []gopi.GPIOPin {
	pins := make([]gopi.GPIOPin, 0, len(this.pins))
	for pin := range this.pins {
		pins = append(pins, pin)
	}
	return pins
}
