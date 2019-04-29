// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	"syscall"
	"time"
)

// Return Host Uptime
func (this *hardware) UptimeHost() time.Duration {
	if info := this.sysinfo_(); info != nil {
		return time.Second * time.Duration(info.Uptime)
	} else {
		return 0
	}
}

// Return load averages
func (this *hardware) LoadAverage() (float64, float64, float64) {
	if info := this.sysinfo_(); info != nil {
		return float64(info.Loads[0]) / float64(1<<16), float64(info.Loads[1]) / float64(1<<16), float64(info.Loads[2]) / float64(1<<16)
	} else {
		return 0, 0, 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET SYSTEM INFO STRUCTURE

func (this *hardware) sysinfo_() *syscall.Sysinfo_t {
	info := syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(&info); err != nil {
		this.log.Error("<hw.linux>sysinfo: %v", err)
		return nil
	} else {
		return &info
	}
}
