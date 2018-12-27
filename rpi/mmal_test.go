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
	for status := rpi.MMAL_MIN; status <= rpi.MMAL_MAX; status++ {
		status_error := fmt.Sprint(status.Error())
		if strings.HasPrefix(status_error, "MMAL_") == false {
			t.Error("Invalid status error ", status_error)
		} else {
			t.Logf("%v => %s", int(status), status_error)
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

func TestStreamType_000(t *testing.T) {
	for stream_type := rpi.MMAL_STREAM_TYPE_MIN; stream_type <= rpi.MMAL_STREAM_TYPE_MAX; stream_type++ {
		stream_type_string := fmt.Sprint(stream_type)
		if strings.HasPrefix(stream_type_string, "MMAL_STREAM_TYPE_") == false {
			t.Error("Invalid stream type string ", stream_type_string)
		} else {
			t.Logf("%v => %s", int(stream_type), stream_type_string)
		}
	}
}

func TestStreamCompareFlags_000(t *testing.T) {
	for stream_flag := rpi.MMAL_STREAM_COMPARE_FLAG_MIN; stream_flag <= rpi.MMAL_STREAM_COMPARE_FLAG_MAX; stream_flag <<= 1 {
		stream_flag_string := fmt.Sprint(stream_flag)
		if strings.HasPrefix(stream_flag_string, "MMAL_STREAM_COMPARE_FLAG_") == false {
			t.Logf("%08X => [not used]", int(stream_flag))
		} else {
			t.Logf("%08X => %s", int(stream_flag), stream_flag_string)
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

func TestComponent_004(t *testing.T) {
	var handle rpi.MMAL_ComponentHandle
	if err := rpi.MMALComponentCreate(rpi.MMAL_COMPONENT_DEFAULT_CLOCK, &handle); err != nil {
		t.Error("Component create error:", err)
	} else if err := rpi.MMALComponentAcquire(handle); err != nil {
		t.Error("Unepxected acquire error:", err)
	} else if err := rpi.MMALComponentRelease(handle); err != nil {
		t.Error("Unepxected release error:", err)
	} else if err := rpi.MMALComponentDestroy(handle); err != nil {
		t.Error("Component destroy error:", err)
	}
}

func TestComponent_005(t *testing.T) {
	var handle rpi.MMAL_ComponentHandle
	if err := rpi.MMALComponentCreate(rpi.MMAL_COMPONENT_DEFAULT_CLOCK, &handle); err != nil {
		t.Error("Component create error:", err)
	} else if err := rpi.MMALComponentEnable(handle); err != nil {
		t.Error("Unepxected enable error:", err)
	} else if enabled := rpi.MMALComponentIsEnabled(handle); enabled == false {
		t.Error("Unepxected enabled value:", enabled)
	} else if err := rpi.MMALComponentDisable(handle); err != nil {
		t.Error("Unepxected disable error:", err)
	} else if enabled := rpi.MMALComponentIsEnabled(handle); enabled == true {
		t.Error("Unepxected enabled value:", enabled)
	} else if err := rpi.MMALComponentDestroy(handle); err != nil {
		t.Error("Component destroy error:", err)
	}
}

func TestComponent_006(t *testing.T) {
	var handle rpi.MMAL_ComponentHandle
	if err := rpi.MMALComponentCreate(rpi.MMAL_COMPONENT_DEFAULT_CLOCK, &handle); err != nil {
		t.Error("Component create error:", err)
	} else if num_ports := rpi.MMALComponentPortNum(handle); num_ports == 0 {
		t.Error("Unexpected number of ports:", num_ports)
	} else if err := rpi.MMALComponentDestroy(handle); err != nil {
		t.Error("Component destroy error:", err)
	}
}

func TestComponent_007(t *testing.T) {
	var handle rpi.MMAL_ComponentHandle
	if err := rpi.MMALComponentCreate(rpi.MMAL_COMPONENT_DEFAULT_CLOCK, &handle); err != nil {
		t.Error("Component create error:", err)
	} else if num_ports := rpi.MMALComponentPortNum(handle); num_ports == 0 {
		t.Error("Unexpected number of ports:", num_ports)
	} else {
		for port_index := uint(0); port_index < uint(num_ports); port_index++ {
			port := rpi.MMALComponentPortAtIndex(handle, port_index)
			port_type := rpi.MMALPortType(port)
			port_name := rpi.MMALPortName(port)
			t.Log(port_index, port_type, port_name)
		}
		if err := rpi.MMALComponentDestroy(handle); err != nil {
			t.Error("Component destroy error:", err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST FORMATS

func TestFormats_001(t *testing.T) {
	var handle rpi.MMAL_ComponentHandle
	if err := rpi.MMALComponentCreate(rpi.MMAL_COMPONENT_DEFAULT_VIDEO_DECODER, &handle); err != nil {
		t.Error("Component create error:", err)
	} else if num_ports := rpi.MMALComponentPortNum(handle); num_ports == 0 {
		t.Error("Unexpected number of ports:", num_ports)
	} else {
		for port_index := uint(0); port_index < uint(num_ports); port_index++ {
			port := rpi.MMALComponentPortAtIndex(handle, port_index)
			port_name := rpi.MMALPortName(port)
			format := rpi.MMALPortFormat(port)
			format_type := rpi.MMALStreamFormatType(format)
			format_encoding, format_variant := rpi.MMALStreamFormatEncoding(format)
			format_bitrate := rpi.MMALStreamFormatBitrate(format)
			t.Log("PORT", port_name)
			t.Log("  FORMAT TYPE", format_type)
			if format_encoding != 0 {
				t.Log("  FORMAT ENCODING", format_encoding)
			}
			if format_variant != 0 {
				t.Log("  FORMAT ENCODING VARIANT", format_variant)
			}
			if format_bitrate != 0 {
				t.Log("  FORMAT BITRATE", format_bitrate)
			}
			if format_type == rpi.MMAL_STREAM_TYPE_VIDEO {
				w, h := rpi.MMALStreamFormatVideoWidthHeight(format)
				t.Log("  FORMAT VIDEO FRAME SIZE {", w, ",", h, "}")
				t.Log("  FORMAT VIDEO CROP ", rpi.MMALStreamFormatVideoCrop(format))
				t.Log("  FORMAT VIDEO FRAME RATE ", rpi.MMALStreamFormatVideoFrameRate(format))
				t.Log("  FORMAT VIDEO PIXEL ASPECT RATIO ", rpi.MMALStreamFormatVideoPixelAspectRatio(format))
				t.Log("  FORMAT VIDEO COLOR SPACE ", rpi.MMALStreamFormatColorSpace(format))
			}
		}
		if err := rpi.MMALComponentDestroy(handle); err != nil {
			t.Error("Component destroy error:", err)
		}
	}
}
