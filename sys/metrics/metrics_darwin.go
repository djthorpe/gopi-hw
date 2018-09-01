/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package metrics

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"time"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
	#include <stdlib.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// RETURN UPTIME

func (this *metrics) UptimeHost() time.Duration {
	tv := syscall.Timeval32{}

	if err := sysctlbyname("kern.boottime", &tv); err != nil {
		this.log.Error("<metrics.darwin>UptimeHost: %v", err)
		return 0
	} else {
		return time.Since(time.Unix(int64(tv.Sec), int64(tv.Usec)*1000))
	}
}

////////////////////////////////////////////////////////////////////////////////
// LOAD AVERAGES

func (this *metrics) LoadAverage() (float64, float64, float64) {
	avg := []C.double{0, 0, 0}
	if C.getloadavg(&avg[0], C.int(len(avg))) == C.int(-1) {
		this.log.Error("<metrics.darwin>LoadAverage: Unavailable")
		return 0, 0, 0
	} else {
		return float64(avg[0]), float64(avg[1]), float64(avg[2])
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET SYSTEM INFO

// Generic Sysctl buffer unmarshalling
// from https://github.com/cloudfoundry/gosigar/blob/master/sigar_darwin.go
func sysctlbyname(name string, data interface{}) error {
	if val, err := syscall.Sysctl(name); err != nil {
		return err
	} else {
		buf := []byte(val)
		switch v := data.(type) {
		case *uint64:
			*v = *(*uint64)(unsafe.Pointer(&buf[0]))
			return nil
		case []byte:
			for i := 0; i < len(val) && i < len(v); i++ {
				v[i] = val[i]
			}
			return nil
		}
		bbuf := bytes.NewBuffer([]byte(val))
		return binary.Read(bbuf, binary.LittleEndian, data)
	}
}
