package openvg_test

import (
	"fmt"
	"strings"
	"testing"

	// Frameworks
	vg "github.com/djthorpe/gopi-hw/openvg"
)

////////////////////////////////////////////////////////////////////////////////
// TEST ENUMS

func TestStatus_000(t *testing.T) {
	for status := vg.VG_ERROR_MIN; status <= vg.VG_ERROR_MAX; status++ {
		status_error := fmt.Sprint(status.Error())
		if strings.HasPrefix(status_error, "VG_ERROR_") {
			t.Logf("%v => %s", int(status), status_error)
		} else {
			t.Errorf("No status error for value: %v", status)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST FLUSH

func TestFlush_000(t *testing.T) {
	if err := vg.VG_Flush(); err != nil {
		t.Error(err)
	}
}

func TestFinish_000(t *testing.T) {
	if err := vg.VG_Finish(); err != nil {
		t.Error(err)
	}
}
