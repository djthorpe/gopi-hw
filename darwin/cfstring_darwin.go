// +build darwin

package darwin

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"

import (
	"unsafe"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	CFString C.CFStringRef
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewCFString(value string) CFString {
	s := C.CString(value)
	defer C.free(unsafe.Pointer(s))
	this := C.CFStringCreateWithCString(C.kCFAllocatorDefault, s, C.kCFStringEncodingUTF8)
	return CFString(this)
}

func (this CFString) Free() {
	if this != 0 {
		C.CFRelease(C.CFTypeRef(this))
	}
}

func (this CFString) String() string {
	if this == 0 {
		return ""
	}
	// Quick version without allocation
	ptr := C.CFStringGetCStringPtr(C.CFStringRef(this), C.kCFStringEncodingUTF8)
	if ptr != nil {
		return C.GoString(ptr)
	}
	// Next find the buffer size necessary
	var length C.CFIndex
	length_ := C.CFStringGetLength(C.CFStringRef(this))
	range_ := C.CFRange{0, length_}
	if C.CFStringGetBytes(C.CFStringRef(this), range_, C.kCFStringEncodingUTF8, 0, C.false, nil, 0, &length) == 0 {
		return ""
	}
	bytes := make([]byte, length)
	if C.CFStringGetBytes(C.CFStringRef(this), range_, C.kCFStringEncodingUTF8, 0, C.false, (*C.uchar)(unsafe.Pointer(&bytes[0])), length, nil) == 0 {
		return ""
	}
	return string(bytes)
}
