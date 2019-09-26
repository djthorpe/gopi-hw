// +build darwin

package darwin

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"

///////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	CFType C.CFTypeRef
)
