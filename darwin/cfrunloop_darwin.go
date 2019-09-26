// +build darwin

package darwin

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"

import (
	"runtime"

	// Frameworks
	"github.com/djthorpe/gopi"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	CFRunLoop C.CFRunLoopRef
)

///////////////////////////////////////////////////////////////////////////////
// METHODS

func CurrentRunLoop() CFRunLoop {
	return CFRunLoop(C.CFRunLoopGetCurrent())
}

func RunLoopInCurrentThread(stop <-chan struct{}) error {
	runtime.LockOSThread()
	stopped := make(chan struct{})
	if runloop := C.CFRunLoopGetCurrent(); runloop != 0 {
		go func() {
			<-stop
			C.CFRunLoopStop(C.CFRunLoopRef(runloop))
			close(stopped)
		}()
		C.CFRetain(C.CFTypeRef(runloop))
		C.CFRunLoopRun()
		<-stopped
		C.CFRelease(C.CFTypeRef(runloop))
		return nil
	} else {
		return gopi.ErrAppError
	}
}
