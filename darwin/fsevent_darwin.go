// +build darwin

package darwin

/*
#cgo LDFLAGS: -framework CoreServices
#include <CoreServices/CoreServices.h>

extern void fsevtCallback(FSEventStreamRef p0, uintptr_t info, size_t p1, char** p2, FSEventStreamEventFlags* p3, FSEventStreamEventId* p4);

static FSEventStreamRef FSEventStreamCreate_(FSEventStreamContext* context, uintptr_t info, CFArrayRef paths, FSEventStreamEventId since, CFTimeInterval latency,FSEventStreamCreateFlags flags) {	
	context->info = (void* )info;
	return FSEventStreamCreate(NULL,(FSEventStreamCallback)fsevtCallback, context, paths, since, latency, flags);
}

*/
import "C"

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	FSEventStream           C.FSEventStreamRef
	FSEventStreamContext    C.FSEventStreamContext
	FSEventStreamCreateFlag C.uint32
	FSEventStreamFlag       C.uint32
	FSEventID               C.FSEventStreamEventId
	FSCallback              func(*FSEvent)
)

type FSEvent struct {
	Stream   FSEventStream
	UserInfo uintptr
	Path     string
	Flags    FSEventStreamFlag
	Event    FSEventID
}

const (
	FS_STREAM_CREATE_FLAG_USECFTYPES FSEventStreamCreateFlag = 1 << iota
	FS_STREAM_CREATE_FLAG_NODEFER    FSEventStreamCreateFlag = 1 << iota
	FS_STREAM_CREATE_FLAG_WATCHROOT  FSEventStreamCreateFlag = 1 << iota
	FS_STREAM_CREATE_FLAG_IGNORESELF FSEventStreamCreateFlag = 1 << iota
	FS_STREAM_CREATE_FLAG_FILEEVENTS FSEventStreamCreateFlag = 1 << iota
	FS_STREAM_CREATE_FLAG_NONE       FSEventStreamCreateFlag = 0
)

const (
	FS_STREAM_FLAG_MUSTSCANSUBDIRS    FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_USERDROPPED        FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_KERNELDROPPED      FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_EVENTIDSWRAPPED    FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_HISTORYDONE        FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ROOTCHANGED        FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_MOUNT              FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_UNMOUNT            FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_CREATED       FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_REMOVED       FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_INODEMETAMOD  FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_RENAMED       FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_MODIFIED      FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_FINDERINFOMOD FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_CHANGEOWNER   FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_XATTRMOD      FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_ISFILE        FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_ISDIR         FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_ITEM_ISSYMLINK     FSEventStreamFlag = 1 << iota
	FS_STREAM_FLAG_MIN                                  = FS_STREAM_FLAG_MUSTSCANSUBDIRS
	FS_STREAM_FLAG_MAX                                  = FS_STREAM_FLAG_ITEM_ISSYMLINK
	FS_STREAM_FLAG_NONE               FSEventStreamFlag = 0
)

///////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	callback_map = make(map[uintptr]FSCallback, 1)
	callback_lock sync.Mutex
)

///////////////////////////////////////////////////////////////////////////////
// CALLBACKS

func SetEventCallback(userInfo uintptr, callback FSCallback) {
	callback_lock.Lock()
	defer callback_lock.Unlock()
	if callback == nil {
		delete(callback_map, userInfo)
	} else {
		callback_map[userInfo] = callback
	}
}

///////////////////////////////////////////////////////////////////////////////
// STREAMS

func NewEventStream(paths []string, userInfo uintptr, since FSEventID, latency time.Duration, flags FSEventStreamCreateFlag) (FSEventStream, error) {
	if len(paths) == 0 {
		return nil, gopi.ErrBadParameter
	}

	// Create an array for the paths
	cfarray := NewCFArray(uint(len(paths)))
	defer cfarray.Free()
	for _, path := range paths {
		if path_, err := filepath.Abs(path); err != nil {
			return nil, fmt.Errorf("%v: %v", err, path)
		} else {
			cfstr := NewCFString(path_)
			defer cfstr.Free()
			cfarray.Append(CFType(cfstr))
		}
	}

	// convert latency to CFTimeInterval
	latency_ := C.CFTimeInterval(float64(latency) / float64(time.Second))

	// Create the stream
	if stream := C.FSEventStreamCreate_(&C.FSEventStreamContext{}, C.uintptr_t(userInfo), C.CFArrayRef(cfarray), C.FSEventStreamEventId(since), latency_, C.FSEventStreamCreateFlags(flags)); stream == nil {
		return nil, gopi.ErrAppError
	} else {
		return FSEventStream(stream), nil
	}
}

func ReleaseEventStream(stream FSEventStream) {
	if stream != nil {
		C.FSEventStreamRelease(stream)
	}
}

func FlushEventStream(stream FSEventStream, sync bool) {
	if sync {
		C.FSEventStreamFlushSync(stream)
	} else {
		C.FSEventStreamFlushAsync(stream)
	}
}

func StartEventStreamInRunloop(stream FSEventStream, runloop CFRunLoop) bool {
	C.FSEventStreamScheduleWithRunLoop(stream, C.CFRunLoopRef(runloop), C.kCFRunLoopDefaultMode)
	if success := C.FSEventStreamStart(stream); success == C.Boolean(0) {
		return false
	} else {
		return true
	}
}

func StopEventStreamInRunLoop(stream FSEventStream, runloop CFRunLoop) {
	C.FSEventStreamStop(stream)
	C.FSEventStreamInvalidate(stream)
}

func LatestEventID() FSEventID {
	return FSEventID(C.FSEventsGetCurrentEventId())
}

///////////////////////////////////////////////////////////////////////////////
// CALLBACK

//export fsevtCallback
func fsevtCallback(stream C.FSEventStreamRef, info uintptr, num C.size_t, paths **C.char, flags *C.FSEventStreamEventFlags, event_ids *C.FSEventStreamEventId) {
	if cb, exists := callback_map[info]; exists {
		paths_ := fsevtPathSlice(int(num), paths)
		event_ids_ := fsevtEventIdSlice(int(num), event_ids)
		flags_ := fsevtFlagSlice(int(num), flags)
		for i := uint(0); i < uint(num); i++ {
			cb(&FSEvent{
				Stream:   FSEventStream(stream),
				UserInfo: info,
				Event:    event_ids_[i],
				Path:     C.GoString(paths_[i]),
				Flags:    flags_[i],
			})
		}
	}
}

func fsevtPathSlice(num int, paths **C.char) []*C.char {
	var paths_ []*C.char
	header := (*reflect.SliceHeader)(unsafe.Pointer(&paths_))
	header.Data = uintptr(unsafe.Pointer(paths))
	header.Len = num
	header.Cap = num
	return paths_
}

func fsevtEventIdSlice(num int, event_ids *C.FSEventStreamEventId) []FSEventID {
	var event_ids_ []FSEventID
	header := (*reflect.SliceHeader)(unsafe.Pointer(&event_ids_))
	header.Data = uintptr(unsafe.Pointer(event_ids))
	header.Len = num
	header.Cap = num
	return event_ids_
}

func fsevtFlagSlice(num int, flags *C.FSEventStreamEventFlags) []FSEventStreamFlag {
	var flags_ []FSEventStreamFlag
	header := (*reflect.SliceHeader)(unsafe.Pointer(&flags_))
	header.Data = uintptr(unsafe.Pointer(flags))
	header.Len = num
	header.Cap = num
	return flags_
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f FSEventStreamFlag) String() string {
	if f == FS_STREAM_FLAG_NONE {
		return "FS_STREAM_FLAG_NONE"
	}
	parts := ""
	for flag := FS_STREAM_FLAG_MIN; flag <= FS_STREAM_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case FS_STREAM_FLAG_MUSTSCANSUBDIRS:
			parts += "|" + "FS_STREAM_FLAG_MUSTSCANSUBDIRS"
		case FS_STREAM_FLAG_USERDROPPED:
			parts += "|" + "FS_STREAM_FLAG_USERDROPPED"
		case FS_STREAM_FLAG_KERNELDROPPED:
			parts += "|" + "FS_STREAM_FLAG_KERNELDROPPED"
		case FS_STREAM_FLAG_EVENTIDSWRAPPED:
			parts += "|" + "FS_STREAM_FLAG_EVENTIDSWRAPPED"
		case FS_STREAM_FLAG_HISTORYDONE:
			parts += "|" + "FS_STREAM_FLAG_HISTORYDONE"
		case FS_STREAM_FLAG_ROOTCHANGED:
			parts += "|" + "FS_STREAM_FLAG_ROOTCHANGED"
		case FS_STREAM_FLAG_MOUNT:
			parts += "|" + "FS_STREAM_FLAG_MOUNT"
		case FS_STREAM_FLAG_UNMOUNT:
			parts += "|" + "FS_STREAM_FLAG_UNMOUNT"
		case FS_STREAM_FLAG_ITEM_CREATED:
			parts += "|" + "FS_STREAM_FLAG_ITEM_CREATED"
		case FS_STREAM_FLAG_ITEM_REMOVED:
			parts += "|" + "FS_STREAM_FLAG_ITEM_REMOVED"
		case FS_STREAM_FLAG_ITEM_INODEMETAMOD:
			parts += "|" + "FS_STREAM_FLAG_ITEM_INODEMETAMOD"
		case FS_STREAM_FLAG_ITEM_RENAMED:
			parts += "|" + "FS_STREAM_FLAG_ITEM_RENAMED"
		case FS_STREAM_FLAG_ITEM_MODIFIED:
			parts += "|" + "FS_STREAM_FLAG_ITEM_MODIFIED"
		case FS_STREAM_FLAG_ITEM_FINDERINFOMOD:
			parts += "|" + "FS_STREAM_FLAG_ITEM_FINDERINFOMOD"
		case FS_STREAM_FLAG_ITEM_CHANGEOWNER:
			parts += "|" + "FS_STREAM_FLAG_ITEM_CHANGEOWNER"
		case FS_STREAM_FLAG_ITEM_XATTRMOD:
			parts += "|" + "FS_STREAM_FLAG_ITEM_XATTRMOD"
		case FS_STREAM_FLAG_ITEM_ISFILE:
			parts += "|" + "FS_STREAM_FLAG_ITEM_ISFILE"
		case FS_STREAM_FLAG_ITEM_ISDIR:
			parts += "|" + "FS_STREAM_FLAG_ITEM_ISDIR"
		case FS_STREAM_FLAG_ITEM_ISSYMLINK:
			parts += "|" + "FS_STREAM_FLAG_ITEM_ISSYMLINK"
		default:
			parts += "|" + "[?? Invalid FSEventStreamFlag value]"
		}
	}
	return strings.Trim(parts, "|")
}
