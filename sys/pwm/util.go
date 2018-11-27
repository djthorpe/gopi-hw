/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package pwm

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	regexpFrequency = regexp.MustCompile("^([0-9]*\\.?[0-9]+)\\s*(hz|khz|mhz)$")
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// ParseFrequency converts a string into a period. It will return an
// error if the string could not be parsed.
func ParseFrequency(value string) (time.Duration, error) {
	if parts := regexpFrequency.FindStringSubmatch(strings.ToLower(value)); len(parts) != 3 {
		return 0, gopi.ErrBadParameter
	} else if num, err := strconv.ParseFloat(parts[1], 64); err != nil {
		return 0, err
	} else {
		switch parts[2] {
		case "hz":
			return time.Nanosecond * time.Duration((1E9 / num)), nil
		case "khz":
			return time.Nanosecond * time.Duration((1E6 / num)), nil
		case "mhz":
			return time.Nanosecond * time.Duration((1E3 / num)), nil
		default:
			return 0, gopi.ErrBadParameter
		}
	}
}
