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

	// Frameworks
	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// SUPPORTED ENCODINGS

func (this *port) SupportedEncodings() ([]hw.MMALEncodingType, error) {
	if param, err := rpi.MMALPortParameterAllocGet(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_SUPPORTED_ENCODINGS, 0); err != nil {
		return nil, err
	} else {
		defer rpi.MMALPortParameterAllocFree(param)
		encodings := make([]hw.MMALEncodingType, 0)
		for _, encoding := range rpi.MMALParamGetArrayUint32(param) {
			encodings = append(encodings, hw.MMALEncodingType(encoding))
		}
		return encodings, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// URI

func (this *port) Uri() (string, error) {
	if param, err := rpi.MMALPortParameterAllocGet(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_URI, 0); err != nil {
		return "", err
	} else {
		defer rpi.MMALPortParameterAllocFree(param)
		return "<TODO>", nil
	}
}

func (this *port) SetUri(value string) error {
	return rpi.MMALPortSetURI(this.handle, value)
}

////////////////////////////////////////////////////////////////////////////////
// ZEROCOPY

func (this *port) ZeroCopy() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_ZERO_COPY)
}

func (this *port) SetZeroCopy(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_ZERO_COPY, value)
}

/*
MMAL_PARAMETER_URI                                                           // Takes a MMAL_PARAMETER_URI_T
MMAL_PARAMETER_CHANGE_EVENT_REQUEST                                          // Takes a MMAL_PARAMETER_CHANGE_EVENT_REQUEST_T
MMAL_PARAMETER_ZERO_COPY                                                     // Takes a MMAL_PARAMETER_BOOLEAN_T
MMAL_PARAMETER_BUFFER_REQUIREMENTS                                           // Takes a MMAL_PARAMETER_BUFFER_REQUIREMENTS_T
MMAL_PARAMETER_STATISTICS                                                    // Takes a MMAL_PARAMETER_STATISTICS_T
MMAL_PARAMETER_CORE_STATISTICS                                               // Takes a MMAL_PARAMETER_CORE_STATISTICS_T
MMAL_PARAMETER_MEM_USAGE                                                     // Takes a MMAL_PARAMETER_MEM_USAGE_T
MMAL_PARAMETER_BUFFER_FLAG_FILTER                                            // Takes a MMAL_PARAMETER_UINT32_T
MMAL_PARAMETER_SEEK                                                          // Takes a MMAL_PARAMETER_SEEK_T
MMAL_PARAMETER_POWERMON_ENABLE                                               // Takes a MMAL_PARAMETER_BOOLEAN_T
MMAL_PARAMETER_LOGGING                                                       // Takes a MMAL_PARAMETER_LOGGING_T
MMAL_PARAMETER_SYSTEM_TIME                                                   // Takes a MMAL_PARAMETER_UINT64_T
MMAL_PARAMETER_NO_IMAGE_PADDING                                              // Takes a MMAL_PARAMETER_BOOLEAN_T
MMAL_PARAMETER_LOCKSTEP_ENABLE                                               // Takes a MMAL_PARAMETER_BOOLEAN_T
*/
