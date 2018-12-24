package pwm_test

import (
	"testing"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	pwm "github.com/djthorpe/gopi-hw/sys/pwm"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
)

var (
	frequency_map = map[string]time.Duration{
		"1hz":     time.Second,
		"1 hz":    time.Second,
		"1 Hz":    time.Second,
		"1   HZ":  time.Second,
		"2hz":     time.Millisecond * 500,
		"1kHz":    time.Microsecond * 1000,
		"100kHz":  time.Microsecond * 10,
		"100MHz":  time.Nanosecond * 10,
		".1 HZ":   time.Second * 10,
		"0.5 kHz": time.Millisecond * 2,
	}
)

////////////////////////////////////////////////////////////////////////////////
// CREATE MODULES / APPS

func TestConfig_000(t *testing.T) {
	// Create config file
	config := gopi.NewAppConfig("pwm")
	t.Log(config)
}

func TestApp_000(t *testing.T) {
	// Create app
	config := gopi.NewAppConfig("pwm")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else {
		t.Log(app)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PARSE FREQUENCY

func TestParseFrequency(t *testing.T) {
	for k, expected_result := range frequency_map {
		if result, err := pwm.ParseFrequency(k); err != nil {
			t.Error(err)
		} else if result != expected_result {
			t.Errorf("TestParseFrequency: '%v' expected %v but got %v", k, expected_result, result)
		} else {
			t.Logf("TestParseFrequency: '%v' => %v", k, result)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET PINS AND SET DUTY CYCLE
func TestDutyCycle(t *testing.T) {
	pin := gopi.GPIOPin(4)
	config := gopi.NewAppConfig("pwm")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if pwm := app.PWM; pwm == nil {
		t.Fatal("Unable to create PWM module")
	} else if err := app.PWM.SetDutyCycle(0, pin); err != nil {
		t.Error("Unable to SetDutyCycle for pin", pin)
	} else if duty_cycle, err := app.PWM.DutyCycle(pin); err != nil {
		t.Error(err)
	} else if duty_cycle != 0 {
		t.Error("Unexpected DutyCycle for pin", pin, " = ", duty_cycle)
	} else if err := app.PWM.SetDutyCycle(1, pin); err != nil {
		t.Error("Unable to SetDutyCycle for pin", pin)
	} else if duty_cycle, err := app.PWM.DutyCycle(pin); err != nil {
		t.Error(err)
	} else if duty_cycle != 1 {
		t.Error("Unexpected DutyCycle for pin", pin, " = ", duty_cycle)
	} else if err := app.PWM.SetDutyCycle(0.5, pin); err != nil {
		t.Error("Unable to SetDutyCycle for pin", pin)
	} else if duty_cycle, err := app.PWM.DutyCycle(pin); err != nil {
		t.Error(err)
	} else if duty_cycle != 0.5 {
		t.Error("Unexpected DutyCycle for pin", pin, " = ", duty_cycle)
	} else if pins := pwm.Pins(); len(pins) != 1 {
		t.Log("pins() failed, pins=", pins)
	} else if pins[0] != pin {
		t.Log("pins() failed, expected pin=", pin, " but got ", pins[0])
	}
}
