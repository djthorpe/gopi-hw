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
// MMAL_PARAMETER_SUPPORTED_ENCODINGS

func (this *port) SupportedEncodings() ([]hw.MMALEncodingType, error) {
	if param, err := rpi.MMALPortParameterAllocGet(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_SUPPORTED_ENCODINGS, rpi.MMAL_PARAMETER_SUPPORTED_ENCODINGS_ARRAY_SIZE); err != nil {
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
// MMAL_PARAMETER_URI

func (this *port) SetUri(value string) error {
	return rpi.MMALPortSetURI(this.handle, value)
}

////////////////////////////////////////////////////////////////////////////////
// MMAL_PARAMETER_ZERO_COPY

func (this *port) ZeroCopy() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_ZERO_COPY)
}

func (this *port) SetZeroCopy(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_ZERO_COPY, value)
}

////////////////////////////////////////////////////////////////////////////////
// MMAL_PARAMETER_NO_IMAGE_PADDING

func (this *port) NoImagePadding() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_NO_IMAGE_PADDING)
}

func (this *port) SetNoImagePadding(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_NO_IMAGE_PADDING, value)
}

////////////////////////////////////////////////////////////////////////////////
// MMAL_PARAMETER_LOCKSTEP_ENABLE

func (this *port) LockstepEnable() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_LOCKSTEP_ENABLE)
}

func (this *port) SetLockstepEnable(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_LOCKSTEP_ENABLE, value)
}

////////////////////////////////////////////////////////////////////////////////
// MMAL_PARAMETER_POWERMON_ENABLE

func (this *port) PowermonEnable() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_POWERMON_ENABLE)
}

func (this *port) SetPowermonEnable(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_POWERMON_ENABLE, value)
}

////////////////////////////////////////////////////////////////////////////////
// MMAL_PARAMETER_BUFFER_FLAG_FILTER

func (this *port) BufferFlagFilter() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_BUFFER_FLAG_FILTER)
}

func (this *port) SetBufferFlagFilter(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_BUFFER_FLAG_FILTER, value)
}

////////////////////////////////////////////////////////////////////////////////
// MMAL_PARAMETER_SYSTEM_TIME

func (this *port) SystemTime() (uint64, error) {
	return rpi.MMALPortParameterGetUint64(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_SYSTEM_TIME)
}

////////////////////////////////////////////////////////////////////////////////
// MMAL_PARAMETER_SEEK

func (this *port) SetSeek(offset int64, precise, forward bool) error {
	flags := uint32(0)
	if precise {
		flags |= rpi.MMAL_PARAMETER_SEEK_FLAG_PRECISE
	}
	if forward {
		flags |= rpi.MMAL_PARAMETER_SEEK_FLAG_FORWARD
	}
	value := rpi.MMAL_ParameterSeek{}
	value.SetOffset(offset)
	value.SetFlags(flags)
	return rpi.MMALPortParameterSetSeek(this.handle, rpi.MMAL_PARAMETER_GROUP_COMMON|rpi.MMAL_PARAMETER_SEEK, value)
}

/*
TODO:
MMAL_PARAMETER_CHANGE_EVENT_REQUEST                                          // Takes a MMAL_PARAMETER_CHANGE_EVENT_REQUEST_T
MMAL_PARAMETER_BUFFER_REQUIREMENTS                                           // Takes a MMAL_PARAMETER_BUFFER_REQUIREMENTS_T
MMAL_PARAMETER_STATISTICS                                                    // Takes a MMAL_PARAMETER_STATISTICS_T
MMAL_PARAMETER_CORE_STATISTICS                                               // Takes a MMAL_PARAMETER_CORE_STATISTICS_T
MMAL_PARAMETER_MEM_USAGE                                                     // Takes a MMAL_PARAMETER_MEM_USAGE_T
MMAL_PARAMETER_LOGGING                                                       // Takes a MMAL_PARAMETER_LOGGING_T
*/
