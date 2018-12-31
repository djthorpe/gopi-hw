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

	"github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// MMAL_PARAMETER_DISPLAYREGION

func (this *port) DisplayRegion() (hw.MMALDisplayRegion, error) {
	if value, err := rpi.MMALPortParameterGetDisplayRegion(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_DISPLAYREGION); err != nil {
		return nil, err
	} else {
		return &displayregion{value}, nil
	}
}

func (this *port) SetDisplayRegion(value hw.MMALDisplayRegion) error {
	if value_, ok := value.(*displayregion); ok == false {
		return gopi.ErrBadParameter
	} else {
		return rpi.MMALPortParameterSetDisplayRegion(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_DISPLAYREGION, value_.handle)
	}
}

////////////////////////////////////////////////////////////////////////////////
// MMAL_PARAMETER_SUPPORTED_PROFILES, MMAL_PARAMETER_PROFILE

func (this *port) SupportedVideoProfiles() ([]hw.MMALVideoProfile, error) {
	size := uint32(0)
	if param, err := rpi.MMALPortParameterAllocGet(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_SUPPORTED_PROFILES, size); err != nil {
		return nil, err
	} else {
		defer rpi.MMALPortParameterAllocFree(param)
		profiles := make([]hw.MMALVideoProfile, 0)
		for _, profile := range rpi.MMALParamGetArrayVideoProfile(param) {
			profiles = append(profiles, profile)
		}
		return profiles, nil
	}
}

func (this *port) VideoProfile() (hw.MMALVideoProfile, error) {
	if value, err := rpi.MMALPortParameterGetVideoProfile(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_PROFILE); err != nil {
		return hw.MMALVideoProfile{}, err
	} else {
		return value, nil
	}
}

func (this *port) SetVideoProfile(value hw.MMALVideoProfile) error {
	return rpi.MMALPortParameterSetVideoProfile(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_PROFILE, value)
}

func (this *port) IntraPeriod() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_INTRAPERIOD)
}

func (this *port) SetIntraPeriod(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_INTRAPERIOD, value)
}

func (this *port) MBRowsPerSlice() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_MB_ROWS_PER_SLICE)
}

func (this *port) SetMBRowsPerSlice(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_MB_ROWS_PER_SLICE, value)
}

func (this *port) Bitrate() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_BIT_RATE)
}

func (this *port) SetBitrate(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_BIT_RATE, value)
}

func (this *port) EncodeMinQuant() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_MIN_QUANT)
}

func (this *port) SetEncodeMinQuant(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_MIN_QUANT, value)
}

func (this *port) EncodeMaxQuant() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_MAX_QUANT)
}

func (this *port) SetEncodeMaxQuant(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_MAX_QUANT, value)
}

func (this *port) ExtraBuffers() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_EXTRA_BUFFERS)
}

func (this *port) SetExtraBuffers(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_EXTRA_BUFFERS, value)
}

func (this *port) AlignHoriz() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ALIGN_HORIZ)
}
func (this *port) SetAlignHoriz(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ALIGN_HORIZ, value)
}

func (this *port) AlignVert() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ALIGN_VERT)
}

func (this *port) SetAlignVert(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ALIGN_VERT, value)
}

func (this *port) EncodeInitialQuant() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_INITIAL_QUANT)
}

func (this *port) SetEncodeInitialQuant(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_INITIAL_QUANT, value)
}

func (this *port) EncodeQPP() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_QP_P)
}

func (this *port) SetEncodeQPP(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_QP_P, value)
}

func (this *port) EncodeRCSliceDQuant() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_RC_SLICE_DQUANT)
}
func (this *port) SetEncodeRCSliceDQuant(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_RC_SLICE_DQUANT, value)
}

func (this *port) EncodeFrameLimitBits() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_FRAME_LIMIT_BITS)
}

func (this *port) SetEncodeFrameLimitBits(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_FRAME_LIMIT_BITS, value)
}

func (this *port) EncodePeakRate() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_PEAK_RATE)
}
func (this *port) SetEncodePeakRate(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_PEAK_RATE, value)
}

func (this *port) EncodeH264DeblockIDC() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_DEBLOCK_IDC)
}

func (this *port) SetEncodeH264DeblockIDC(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_DEBLOCK_IDC, value)
}

func (this *port) MaxNumCallbacks() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_MAX_NUM_CALLBACKS)
}

func (this *port) SetMaxNumCallbacks(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_MAX_NUM_CALLBACKS, value)
}

func (this *port) DroppablePFrameLength() (uint32, error) {
	return rpi.MMALPortParameterGetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_DROPPABLE_PFRAME_LENGTH)
}

func (this *port) SetDroppablePFrameLength(value uint32) error {
	return rpi.MMALPortParameterSetUint32(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_DROPPABLE_PFRAME_LENGTH, value)
}

func (this *port) MinimiseFragmentation() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_MINIMISE_FRAGMENTATION)
}
func (this *port) SetMinimiseFragmentation(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_MINIMISE_FRAGMENTATION, value)

}

func (this *port) RequestIFrame() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_REQUEST_I_FRAME)
}
func (this *port) SetRequestIFrame(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_REQUEST_I_FRAME, value)

}

func (this *port) ImmutableInput() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_IMMUTABLE_INPUT)
}

func (this *port) SetImmutableInput(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_IMMUTABLE_INPUT, value)

}

func (this *port) DroppablePFrames() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_DROPPABLE_PFRAMES)

}
func (this *port) SetDroppablePFrames(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_DROPPABLE_PFRAMES, value)

}

func (this *port) EncodeH264DisableCABAC() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_DISABLE_CABAC)

}
func (this *port) SetEncodeH264DisableCABAC(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_DISABLE_CABAC, value)
}

func (this *port) EncodeH264LowLatency() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_LOW_LATENCY)

}
func (this *port) SetEncodeH264LowLatency(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_LOW_LATENCY, value)
}

func (this *port) EncodeH264AUDelimiters() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_AU_DELIMITERS)

}
func (this *port) SetEncodeH264AUDelimiters(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_AU_DELIMITERS, value)
}

func (this *port) EncodeHeaderOnOpen() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_HEADER_ON_OPEN)

}
func (this *port) SetEncodeHeaderOnOpen(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_HEADER_ON_OPEN, value)

}

func (this *port) EncodePrecodeForQP() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_PRECODE_FOR_QP)

}
func (this *port) SetEncodePrecodeForQP(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_PRECODE_FOR_QP, value)

}

func (this *port) TimestampFIFO() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_TIMESTAMP_FIFO)

}
func (this *port) SetTimestampFIFO(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_TIMESTAMP_FIFO, value)

}

func (this *port) DecodeErrorConcealment() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_DECODE_ERROR_CONCEALMENT)

}
func (this *port) SetDecodeErrorConcealment(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_DECODE_ERROR_CONCEALMENT, value)

}

func (this *port) Encode264VCLHRDParameters() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_VCL_HRD_PARAMETERS)

}
func (this *port) SetEncode264VCLHRDParameters(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_VCL_HRD_PARAMETERS, value)

}

func (this *port) Encode264LowDelayHRDFlag() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_LOW_DELAY_HRD_FLAG)

}
func (this *port) SetEncode264LowDelayHRDFlag(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_H264_LOW_DELAY_HRD_FLAG, value)

}

func (this *port) Encode264EncodeInlineHeader() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_INLINE_HEADER)

}
func (this *port) SetEncode264EncodeInlineHeader(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_INLINE_HEADER, value)

}

func (this *port) EncodeSEIEnable() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_SEI_ENABLE)

}
func (this *port) SetEncodeSEIEnable(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_SEI_ENABLE, value)

}

func (this *port) EncodeInlineVectors() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_INLINE_VECTORS)

}
func (this *port) SetEncodeInlineVectors(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_INLINE_VECTORS, value)

}

func (this *port) InterpolateTimestamps() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_INTERPOLATE_TIMESTAMPS)

}
func (this *port) SetInterpolateTimestamps(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_INTERPOLATE_TIMESTAMPS, value)

}

func (this *port) EncodeSPSTiming() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_SPS_TIMING)

}
func (this *port) SetEncodeSPSTiming(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_SPS_TIMING, value)

}

func (this *port) EncodeSeparateNALBufs() (bool, error) {
	return rpi.MMALPortParameterGetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_SEPARATE_NAL_BUFS)

}
func (this *port) SetEncodeSeparateNALBufs(value bool) error {
	return rpi.MMALPortParameterSetBool(this.handle, rpi.MMAL_PARAMETER_GROUP_VIDEO|rpi.MMAL_PARAMETER_VIDEO_ENCODE_SEPARATE_NAL_BUFS, value)

}

/*
TODO:
MMAL_PARAMETER_RATECONTROL                                                    // Takes a MMAL_PARAMETER_VIDEO_RATECONTROL_T
MMAL_PARAMETER_NALUNITFORMAT                                                  // Takes a MMAL_PARAMETER_VIDEO_NALUNITFORMAT_T
MMAL_PARAMETER_VIDEO_LEVEL_EXTENSION                                          // Takes a MMAL_PARAMETER_VIDEO_LEVEL_EXTENSION_T
MMAL_PARAMETER_VIDEO_EEDE_ENABLE                                              // Takes a MMAL_PARAMETER_VIDEO_EEDE_ENABLE_T
MMAL_PARAMETER_VIDEO_EEDE_LOSSRATE                                            // Takes a MMAL_PARAMETER_VIDEO_EEDE_LOSSRATE_T
MMAL_PARAMETER_VIDEO_INTRA_REFRESH                                            // Takes a MMAL_PARAMETER_VIDEO_INTRA_REFRESH_T
MMAL_PARAMETER_VIDEO_FRAME_RATE                                               // Takes a MMAL_PARAMETER_FRAME_RATE_T
MMAL_PARAMETER_VIDEO_ENCODE_RC_MODEL                                          // Takes a MMAL_PARAMETER_VIDEO_ENCODE_RC_MODEL_T
MMAL_PARAMETER_VIDEO_ENCODE_H264_MB_INTRA_MODE                                // Takes a MMAL_PARAMETER_VIDEO_ENCODER_H264_MB_INTRA_MODES_T
MMAL_PARAMETER_VIDEO_DRM_INIT_INFO                                            // Takes a MMAL_PARAMETER_VIDEO_DRM_INIT_INFO_T
MMAL_PARAMETER_VIDEO_DRM_PROTECT_BUFFER                                       // Takes a MMAL_PARAMETER_VIDEO_DRM_PROTECT_BUFFER_T
MMAL_PARAMETER_VIDEO_DECODE_CONFIG_VD3                                        // Takes a MMAL_PARAMETER_BYTES_T
MMAL_PARAMETER_VIDEO_RENDER_STATS                                             // Takes a MMAL_PARAMETER_VIDEO_RENDER_STATS_T
MMAL_PARAMETER_VIDEO_INTERLACE_TYPE                                           // Takes a MMAL_PARAMETER_VIDEO_INTERLACE_TYPE_T
MMAL_PARAMETER_VIDEO_SOURCE_PATTERN                                           // Takes a MMAL_PARAMETER_SOURCE_PATTERN_T
*/
