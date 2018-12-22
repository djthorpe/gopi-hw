package rpi_test

import (
	"fmt"
	"strings"
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TEST ENUMS

func TestStatus_000(t *testing.T) {
	for status := rpi.MMAL_SUCCESS; status <= rpi.MMAL_MAX; status++ {
		status_error := fmt.Sprint(status.Error())
		status_string := fmt.Sprint(status.String())
		if strings.HasPrefix(status_error, "MMAL_") == false {
			t.Error("Invalid status error ", status_error)
		} else {
			t.Logf("%v => %s, %s", int(status), status_error, status_string)
		}
	}
}
func TestPortType_000(t *testing.T) {
	for port_type := rpi.MMAL_PORT_TYPE_NONE; port_type <= rpi.MMAL_PORT_TYPE_MAX; port_type++ {
		port_type_string := fmt.Sprint(port_type)
		if strings.HasPrefix(port_type_string, "MMAL_PORT_TYPE_") == false {
			t.Error("Invalid port type string ", port_type_string)
		} else {
			t.Logf("%v => %s", int(port_type), port_type_string)
		}
	}
}
func TestPortCapability_000(t *testing.T) {
	for port_cap := rpi.MMAL_PORT_CAPABILITY_MIN; port_cap <= rpi.MMAL_PORT_CAPABILITY_MAX; port_cap++ {
		port_cap_string := fmt.Sprint(port_cap)
		if strings.HasPrefix(port_cap_string, "MMAL_PORT_CAPABILITY_") == false {
			t.Error("Invalid port capability string ", port_cap_string)
		} else {
			t.Logf("%v => %s", int(port_cap), port_cap_string)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST COMPONENTS

func TestComponent_000(t *testing.T) {
	var handle rpi.MMAL_ComponentHandle
	if err := rpi.MMALComponentCreate(rpi.MMAL_COMPONENT_DEFAULT_CLOCK, &handle); err != nil {
		t.Error("Component create error:", err)
	} else if err := rpi.MMALComponentDestroy(handle); err != nil {
		t.Error("Component destroy error:", err)
	}
}

func TestComponent_001(t *testing.T) {
	var handle rpi.MMAL_ComponentHandle
	if err := rpi.MMALComponentCreate(rpi.MMAL_COMPONENT_DEFAULT_CLOCK, &handle); err != nil {
		t.Error("Component create error:", err)
	} else if name := rpi.MMALComponentName(handle); name != rpi.MMAL_COMPONENT_DEFAULT_CLOCK {
		t.Error("Unepxected component name:", name)
	} else if err := rpi.MMALComponentDestroy(handle); err != nil {
		t.Error("Component destroy error:", err)
	}
}

func TestComponent_002(t *testing.T) {
	var handle rpi.MMAL_ComponentHandle
	if err := rpi.MMALComponentCreate(rpi.MMAL_COMPONENT_DEFAULT_CLOCK, &handle); err != nil {
		t.Error("Component create error:", err)
	} else if component_id := rpi.MMALComponentId(handle); component_id == 0 {
		t.Error("Unepxected component_id:", component_id)
	} else if err := rpi.MMALComponentDestroy(handle); err != nil {
		t.Error("Component destroy error:", err)
	} else {
		t.Logf("component_id=%08X", component_id)
	}
}
func TestComponent_003(t *testing.T) {
	var handle rpi.MMAL_ComponentHandle
	if err := rpi.MMALComponentCreate(rpi.MMAL_COMPONENT_DEFAULT_CLOCK, &handle); err != nil {
		t.Error("Component create error:", err)
	} else if control := rpi.MMALComponentControlPort(handle); control == nil {
		t.Error("Unepxected control port")
	} else if err := rpi.MMALComponentDestroy(handle); err != nil {
		t.Error("Component destroy error:", err)
	} else if control_type := rpi.MMALPortType(control); control_type != rpi.MMAL_PORT_TYPE_CONTROL {
		t.Logf("Unexpected control port type=%v", control_type)
	}
}
