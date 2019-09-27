/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package fsnotify

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
	mac "github.com/djthorpe/gopi-hw/darwin"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type fsnotify_impl struct {
	log      gopi.Logger
	runloop  mac.CFRunLoop
	stop     chan struct{}
	userInfo map[uintptr]string
	callback func(hw.FSEvent)
}

type fsevent_impl struct {
	ts      time.Time
	root    string
	fsevent *mac.FSEvent
}

////////////////////////////////////////////////////////////////////////////////
// WATCH AND UNWATCH CONCRETE IMPLEMENTATION

func (this *fsnotify_impl) init(callback func(hw.FSEvent), log gopi.Logger) error {
	this.log = log
	this.stop = make(chan struct{})
	this.userInfo = make(map[uintptr]string)
	this.callback = callback

	var wait sync.WaitGroup
	wait.Add(1)

	// Call runloop in background
	go func() {
		this.runloop = mac.CurrentRunLoop()
		wait.Done()

		// Watch a temporary path - need at least one slot in the runloop
		if path, err := ioutil.TempDir(os.TempDir(), "fsnotify"); err != nil {
			this.log.Error("fsnotify: %v", err)
		} else if stream, err := this.watch(path); err != nil {
			this.log.Error("fsnotify: %v", err)
			os.Remove(path)
		} else {
			defer this.unwatch(stream)
			if err := mac.RunLoopInCurrentThread(this.stop); err != nil {
				this.log.Error("fsnotify: %v", err)
			}
			if err := os.Remove(path); err != nil {
				this.log.Error("fsnotify: %v", err)
			}
		}
		// Signal end of runloop
		close(this.stop)
	}()

	// Wait until the runloop member is set before returning
	wait.Wait()

	return nil
}

func (this *fsnotify_impl) close() error {
	// Signal end of runloop
	this.stop <- gopi.DONE
	// Wait for close()
	<-this.stop
	return nil
}

func (this *fsnotify_impl) watch(path string) (interface{}, error) {
	latest := mac.LatestEventID()
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return nil, gopi.ErrNotFound
	} else if err != nil {
		return nil, err
	} else if stat.IsDir() == false {
		return nil, gopi.ErrBadParameter
	} else if userInfo := this.userinfo_for_path(path); userInfo == 0 {
		return nil, gopi.ErrBadParameter
	} else if stream, err := mac.NewEventStream([]string{path}, userInfo, latest, 100*time.Millisecond, mac.FS_STREAM_CREATE_FLAG_FILEEVENTS|mac.FS_STREAM_CREATE_FLAG_WATCHROOT); err != nil {
		return nil, err
	} else {
		mac.SetEventCallback(userInfo, this.emit)
		mac.StartEventStreamInRunloop(stream, this.runloop)
		return stream, nil
	}
}

func (this *fsnotify_impl) unwatch(stream interface{}) error {
	if stream == nil {
		return gopi.ErrBadParameter
	} else if stream_, ok := stream.(mac.FSEventStream); ok == false {
		return gopi.ErrAppError
	} else {
		mac.StopEventStreamInRunLoop(stream_, this.runloop)
		// Delete callback for event stream and release event stream
		mac.SetEventCallback(0, nil)
		mac.ReleaseEventStream(stream_)
		return nil
	}
}

func (this *fsnotify_impl) userinfo_for_path(path string) uintptr {
	h := fnv.New64()
	if _, err := h.Write([]byte(path)); err != nil {
		return 0
	} else if userInfo := uintptr(h.Sum64()); userInfo == 0 {
		return 0
	} else if _, exists := this.userInfo[userInfo]; exists {
		return 0
	} else {
		this.userInfo[userInfo] = path
		return userInfo
	}
}

func (this *fsnotify_impl) path_for_userinfo(userInfo uintptr) string {
	if path, exists := this.userInfo[userInfo]; exists {
		return path
	} else {
		return ""
	}
}

func (this *fsnotify_impl) emit(evt *mac.FSEvent) {
	if rootPath := this.path_for_userinfo(evt.UserInfo); rootPath == "" {
		// Ignore when we don't know the rootPath
		return
	} else {
		this.callback(&fsevent_impl{time.Now(), rootPath, evt})
	}
}

////////////////////////////////////////////////////////////////////////////////
// FSEvent implementation

func (*fsevent_impl) Name() string {
	return "FSEvent"
}

func (*fsevent_impl) Source() gopi.Driver {
	return nil
}

func (this *fsevent_impl) Flags() hw.FSFlag {
	f := hw.FS_FLAG_NONE
	if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_REMOVED != 0 {
		f |= hw.FS_FLAG_DELETED
	} else if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_RENAMED != 0 {
		f |= hw.FS_FLAG_RENAMED
	} else if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_INODEMETAMOD != 0 {
		f |= hw.FS_FLAG_CHMOD
	} else if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_XATTRMOD != 0 {
		f |= hw.FS_FLAG_CHMOD
	} else if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_CHANGEOWNER != 0 {
		f |= hw.FS_FLAG_CHMOD
	} else if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_FINDERINFOMOD != 0 {
		f |= hw.FS_FLAG_CHMOD
	} else if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_MODIFIED != 0 {
		f |= hw.FS_FLAG_MODIFIED
	} else if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_CREATED != 0 {
		f |= hw.FS_FLAG_CREATED
	}

	if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_ISDIR != 0 {
		f |= hw.FS_FLAG_ISFOLDER
	} else if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_ISFILE != 0 {
		f |= hw.FS_FLAG_ISFILE
	} else if this.fsevent.Flags&mac.FS_STREAM_FLAG_ITEM_ISSYMLINK != 0 {
		f |= hw.FS_FLAG_ISSYMLINK
	}
	return f
}

func (this *fsevent_impl) Path() string {
	return this.fsevent.Path
}

func (this *fsevent_impl) RelPath() string {
	if rel, err := filepath.Rel(this.root, this.fsevent.Path); err == nil {
		return rel
	} else {
		return ""
	}
}

func (this *fsevent_impl) Timestamp() time.Time {
	return this.ts
}

func (this *fsevent_impl) Root() string {
	return this.root
}

func (this *fsevent_impl) String() string {
	if rel, err := filepath.Rel(this.root, this.fsevent.Path); err == nil {
		return fmt.Sprintf("<fsevent>{ root=%v path=%v flags=%v ts=%v }", strconv.Quote(this.root), strconv.Quote(rel), this.Flags(), this.ts.Format(time.Kitchen))
	} else {
		return fmt.Sprintf("<fsevent>{ path=%v flags=%v ts=%v }", strconv.Quote(this.fsevent.Path), this.Flags(), this.ts.Format(time.Kitchen))
	}
}
