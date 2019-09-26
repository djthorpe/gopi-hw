package darwin_test

import (
	"fmt"
	"testing"

	// Frameworks
	mac "github.com/djthorpe/gopi-hw/darwin"
)

func Test_CFArray000(t *testing.T) {
	cfarray := mac.NewCFArray(0)
	defer cfarray.Free()

	if cfarray.Len() != 0 {
		t.Error("Unexpected array size,", cfarray.Len())
	}
}

func Test_CFArray001(t *testing.T) {
	cfarray := mac.NewCFArray(0)
	defer cfarray.Free()

	cfstr := mac.NewCFString("TEST")
	defer cfstr.Free()
	cfarray.Append(mac.CFType(cfstr))

	if cfarray.Len() != 1 {
		t.Error("Unexpected array size,", cfarray.Len())
	}

	if cfarray.AtIndex(0) != mac.CFType(cfstr) {
		t.Error("Unexpected array element, ", cfarray.AtIndex(0))
	}
}

func Test_CFArray002(t *testing.T) {
	cfarray := mac.NewCFArray(0)
	defer cfarray.Free()

	for i := uint(0); i < 100; i++ {
		gostr := fmt.Sprintf("%v TEST %v", i, i)
		cfstr := mac.NewCFString(gostr)
		defer cfstr.Free()

		cfarray.Append(mac.CFType(cfstr))

		if cfarray.Len() != i+1 {
			t.Error("Unexpected array size,", cfarray.Len())
		}

		for j := uint(0); j < cfarray.Len(); j++ {
			if mac.CFString(cfarray.AtIndex(j)).String() != fmt.Sprintf("%v TEST %v", j, j) {
				t.Error("Unexpected array element:", j)
			}
		}
	}
}
