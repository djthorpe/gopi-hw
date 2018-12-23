// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package mmal

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *format) String() string {
	return fmt.Sprintf("<sys.hw.mmal.format>{}")
}
