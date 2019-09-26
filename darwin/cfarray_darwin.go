// +build darwin

package darwin

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"
import "unsafe"

///////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	CFArray C.CFMutableArrayRef
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewCFArray(capacity uint) CFArray {
	this := C.CFArrayCreateMutable(0, C.CFIndex(capacity), nil)
	return CFArray(this)
}

func (this CFArray) Free() {
	if this != 0 {
		C.CFRelease(C.CFTypeRef(this))
	}
}

func (this CFArray) Len() uint {
	if this == 0 {
		return 0
	}
	return uint(C.CFArrayGetCount(C.CFArrayRef(this)))
}

func (this CFArray) Append(element CFType) {
	if this != 0 {
		C.CFArrayAppendValue(C.CFMutableArrayRef(this), unsafe.Pointer(element))
	}
}

func (this CFArray) AtIndex(i uint) CFType {
	if this == 0 {
		return 0
	}
	return CFType(C.CFArrayGetValueAtIndex(C.CFArrayRef(this), C.CFIndex(i)))
}
