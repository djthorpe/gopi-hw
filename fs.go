/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2018-2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	// Frameworks
	"strings"
	"time"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type FSFlag uint64

type FSNotify interface {
	gopi.Driver
	gopi.Publisher

	Watch(string) error
	Unwatch(string) error
}

type FSEvent interface {
	gopi.Event

	Root() string
	Path() string
	RelPath() string
	Flags() FSFlag
	Timestamp() time.Time
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FS_FLAG_RENAMED FSFlag = (1 << iota)
	FS_FLAG_CREATED
	FS_FLAG_DELETED
	FS_FLAG_MODIFIED
	FS_FLAG_CHMOD
	FS_FLAG_ISFILE
	FS_FLAG_ISFOLDER
	FS_FLAG_ISSYMLINK
	FS_FLAG_NONE FSFlag = 0
	FS_FLAG_MIN         = FS_FLAG_RENAMED
	FS_FLAG_MAX         = FS_FLAG_ISSYMLINK
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f FSFlag) String() string {
	if f == FS_FLAG_NONE {
		return "FS_FLAG_NONE"
	}
	parts := ""
	for flag := FS_FLAG_MIN; flag <= FS_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case FS_FLAG_RENAMED:
			parts += "|" + "FS_FLAG_RENAMED"
		case FS_FLAG_CREATED:
			parts += "|" + "FS_FLAG_CREATED"
		case FS_FLAG_DELETED:
			parts += "|" + "FS_FLAG_DELETED"
		case FS_FLAG_MODIFIED:
			parts += "|" + "FS_FLAG_MODIFIED"
		case FS_FLAG_CHMOD:
			parts += "|" + "FS_FLAG_CHMOD"
		case FS_FLAG_ISFILE:
			parts += "|" + "FS_FLAG_ISFILE"
		case FS_FLAG_ISFOLDER:
			parts += "|" + "FS_FLAG_ISFOLDER"
		case FS_FLAG_ISSYMLINK:
			parts += "|" + "FS_FLAG_ISSYMLINK"
		default:
			parts += "|" + "[?? Invalid FSFlag value]"
		}
	}
	return strings.Trim(parts, "|")
}
