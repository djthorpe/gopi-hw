/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

// This module provides the ability to do pulse-width-modulation (PWM)
// on GPIO pins, affecting both the period of the pulses and the "duty
// cycle" which determines how long a pulse is "high" for, compared to
// the "low" state.
//
// RASPBERRY PI IMPLEMENTATION
//
// The Raspberry Pi implementation uses "Pi Blaster" which can be cloned
// from https://github.com/sarfata/pi-blaster and installed using the
// following set of commands:
//
//    git clone https://github.com/sarfata/pi-blaster
//    cd pi-blaster/
//    sudo apt-get install autoconf
//    ./autogen.sh && ./configure && make
//    sudo make install
//
// This will run the system daemon and you will end up with a named pipe
// called /dev/pi-blaster
//
// You can then run the included pwm_ctrl command in order to control the
// PWM function. For example:
//
//    go run -tags rpi ./cmd/pwm_ctrl/...
//
package pwm
