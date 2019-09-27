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
	"path/filepath"
	"strconv"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
	errors "github.com/djthorpe/gopi/util/errors"
	event "github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type FSNotify struct{}

type fsnotify struct {
	log     gopi.Logger
	watches map[string]interface{}

	fsnotify_impl
	event.Publisher
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config FSNotify) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<fsnotify.Open>{ %+v }", config)

	this := new(fsnotify)
	this.log = log
	this.watches = make(map[string]interface{})

	// Implementation init
	if err := this.init(this.emit, this.log); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}

func (this *fsnotify) Close() error {
	this.log.Debug("<fsnotify.Close>{}")
	var errs errors.CompoundError

	// Unwatch all paths being watched
	for k := range this.watches {
		errs.Add(this.Unwatch(k))
	}

	// Implementation close
	errs.Add(this.close())

	// Release resources
	this.watches = nil

	return errs.ErrorOrSelf()
}

////////////////////////////////////////////////////////////////////////////////
// WATCH AND UNWATCH

func (this *fsnotify) Watch(path string) error {
	this.log.Debug2("<fsnotify.Watch>{ path=%v }", strconv.Quote(path))

	this.Lock()
	defer this.Unlock()

	if path = filepath.Clean(path); path == "" {
		return gopi.ErrBadParameter
	} else if _, exists := this.watches[path]; exists {
		return fmt.Errorf("Already watched: %v", strconv.Quote(path))
	} else if watch, err := this.watch(path); err != nil {
		return err
	} else {
		this.watches[path] = watch
	}

	// Success
	return nil
}

func (this *fsnotify) Unwatch(path string) error {
	this.log.Debug2("<fsnotify.Unwatch>{ path=%v }", strconv.Quote(path))

	this.Lock()
	defer this.Unlock()

	if path = filepath.Clean(path); path == "" {
		return gopi.ErrBadParameter
	} else if watch, exists := this.watches[path]; exists == false {
		return gopi.ErrBadParameter
	} else {
		delete(this.watches, path)
		if err := this.unwatch(watch); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *fsnotify) String() string {
	paths := make([]string, len(this.watches))
	for k := range this.watches {
		paths = append(paths, strconv.Quote(k))
	}
	return fmt.Sprintf("<fsnotify>{ watches=%v }", paths)
}

////////////////////////////////////////////////////////////////////////////////
// EMIT EVENTS

func (this *fsnotify) emit(evt hw.FSEvent) {
	root := evt.Root()
	if root == "" || evt.Flags() == hw.FS_FLAG_NONE {
		// Ignore events with no flags or root
	} else if _, exists := this.watches[root]; exists {
		this.Emit(evt)
	}
}
