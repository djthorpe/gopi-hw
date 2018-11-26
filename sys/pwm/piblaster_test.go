package pwm_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/gpio"
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/metrics"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE MODULE LISTS

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
