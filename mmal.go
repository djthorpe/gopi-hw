/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2018-2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALFormatType          uint
	MMALEncodingType        uint32
	MMALDisplayTransform    uint
	MMALDisplayMode         uint
	MMALPortConnectionFlags uint
	MMALVideoEncProfile     uint
	MMALVideoEncLevel       uint
	MMALCameraFlashType     uint
	MMALCameraMeteringMode  uint
	MMALCameraExposureMode  uint
	MMALTextJustify         uint
	MMALBufferFlag          uint
)

type MMALVideoProfile struct {
	Profile MMALVideoEncProfile
	Level   MMALVideoEncLevel
}

type MMALRationalNum struct {
	Num, Den int32
}

type MMALRect struct {
	X, Y int32
	W, H uint32
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type MMAL interface {
	gopi.Driver

	// Return components
	ComponentWithName(name string) (MMALComponent, error)

	// Return specific components
	CameraComponent() (MMALCameraComponent, error)
	CameraInfoComponent() (MMALComponent, error)
	VideoDecoderComponent() (MMALComponent, error)
	VideoEncoderComponent() (MMALComponent, error)
	VideoRendererComponent() (MMALComponent, error)
	ImageEncoderComponent() (MMALComponent, error)
	ImageDecoderComponent() (MMALComponent, error)
	ReaderComponent() (MMALComponent, error)
	WriterComponent() (MMALComponent, error)

	// Connect and disconnect component ports
	Connect(input, output MMALPort, flags MMALPortConnectionFlags) (MMALPortConnection, error)
	Disconnect(MMALPortConnection) error
}

type MMALComponent interface {
	Name() string
	Id() uint32

	// Enable and disable
	Enabled() bool
	SetEnabled(value bool) error

	// Return ports
	Control() MMALPort
	Clocks() []MMALPort
	Inputs() []MMALPort
	Outputs() []MMALPort

	// Get buffer from port/connection and optionally block
	GetEmptyBufferOnPort(MMALPort, bool) (MMALBuffer, error)
	GetFullBufferOnPort(MMALPort, bool) (MMALBuffer, error)
}

type MMALCameraComponent interface {
	MMALComponent
}

type MMALPort interface {
	Name() string
	CapabilityPassthrough() bool
	CapabilityAllocation() bool
	CapabilitySupportsEventFormatChange() bool

	// Enable and disable
	Enabled() bool
	SetEnabled(value bool) error

	// Connect, disconnect, flush and return any errors
	Connect(other MMALPort) error
	Disconnect() error
	Flush() error
	Error() error

	// Format
	Format() MMALFormat
	CopyFormat(MMALFormat) error
	VideoFormat() MMALVideoFormat
	AudioFormat() MMALAudioFormat
	SubpictureFormat() MMALSubpictureFormat
	CommitFormatChange() error

	// Send buffer to port, release buffer
	Send(MMALBuffer) error
	Release(MMALBuffer) error

	// Port Parameters
	MMALCommonParameters
	MMALVideoParameters
	MMALCameraParameters
}

type MMALBuffer interface {
	// Fill buffer
	Fill(io.Reader) (uint32, error)

	// Return buffer data
	Data() []byte

	// Buffer flags
	Flags() MMALBufferFlag
}

type MMALPortConnection interface {
	// Input and Output ports
	Input() MMALPort
	Output() MMALPort

	// Enable and disable
	Enabled() bool
	SetEnabled(value bool) error

	// Acquire and release
	Acquire() error
	Release() error
}

type MMALCommonParameters interface {
	// Get Parameters
	SupportedEncodings() ([]MMALEncodingType, error)
	ZeroCopy() (bool, error)
	NoImagePadding() (bool, error)
	LockstepEnable() (bool, error)
	PowermonEnable() (bool, error)
	BufferFlagFilter() (uint32, error)
	SystemTime() (uint64, error)

	// Set Parameters
	SetUri(value string) error
	SetZeroCopy(value bool) error
	SetNoImagePadding(value bool) error
	SetLockstepEnable(value bool) error
	SetPowermonEnable(value bool) error
	SetBufferFlagFilter(value uint32) error
}

type MMALVideoParameters interface {
	// Get Parameters
	DisplayRegion() (MMALDisplayRegion, error)
	SupportedVideoProfiles() ([]MMALVideoProfile, error)
	VideoProfile() (MMALVideoProfile, error)
	IntraPeriod() (uint32, error)
	MBRowsPerSlice() (uint32, error)
	Bitrate() (uint32, error)
	EncodeMinQuant() (uint32, error)
	EncodeMaxQuant() (uint32, error)
	ExtraBuffers() (uint32, error)
	AlignHoriz() (uint32, error)
	AlignVert() (uint32, error)
	EncodeInitialQuant() (uint32, error)
	EncodeQPP() (uint32, error)
	EncodeRCSliceDQuant() (uint32, error)
	EncodeFrameLimitBits() (uint32, error)
	EncodePeakRate() (uint32, error)
	EncodeH264DeblockIDC() (uint32, error)
	MaxNumCallbacks() (uint32, error)
	DroppablePFrameLength() (uint32, error)
	MinimiseFragmentation() (bool, error)
	RequestIFrame() (bool, error)
	ImmutableInput() (bool, error)
	DroppablePFrames() (bool, error)
	EncodeH264DisableCABAC() (bool, error)
	EncodeH264AUDelimiters() (bool, error)
	EncodeHeaderOnOpen() (bool, error)
	EncodePrecodeForQP() (bool, error)
	TimestampFIFO() (bool, error)
	DecodeErrorConcealment() (bool, error)
	Encode264VCLHRDParameters() (bool, error)
	Encode264LowDelayHRDFlag() (bool, error)
	Encode264EncodeInlineHeader() (bool, error)
	EncodeSEIEnable() (bool, error)
	EncodeInlineVectors() (bool, error)
	InterpolateTimestamps() (bool, error)
	EncodeSPSTiming() (bool, error)
	EncodeSeparateNALBufs() (bool, error)
	EncodeH264LowLatency() (bool, error)

	// Set Parameters
	SetDisplayRegion(MMALDisplayRegion) error
	SetVideoProfile(MMALVideoProfile) error
	SetMBRowsPerSlice(uint32) error
	SetBitrate(uint32) error
	SetEncodeMinQuant(uint32) error
	SetEncodeMaxQuant(uint32) error
	SetExtraBuffers(uint32) error
	SetAlignHoriz(uint32) error
	SetAlignVert(uint32) error
	SetEncodeInitialQuant(uint32) error
	SetEncodeQPP(uint32) error
	SetEncodeRCSliceDQuant(uint32) error
	SetEncodeFrameLimitBits(uint32) error
	SetEncodePeakRate(uint32) error
	SetEncodeH264DeblockIDC(uint32) error
	SetMaxNumCallbacks(uint32) error
	SetDroppablePFrameLength(uint32) error
	SetMinimiseFragmentation(bool) error
	SetRequestIFrame(bool) error
	SetImmutableInput(bool) error
	SetDroppablePFrames(bool) error
	SetEncodeH264DisableCABAC(bool) error
	SetEncodeH264AUDelimiters(bool) error
	SetEncodeHeaderOnOpen(bool) error
	SetEncodePrecodeForQP(bool) error
	SetTimestampFIFO(bool) error
	SetDecodeErrorConcealment(bool) error
	SetEncode264VCLHRDParameters(bool) error
	SetEncode264LowDelayHRDFlag(bool) error
	SetEncode264EncodeInlineHeader(bool) error
	SetEncodeSEIEnable(bool) error
	SetEncodeInlineVectors(bool) error
	SetInterpolateTimestamps(bool) error
	SetEncodeSPSTiming(bool) error
	SetEncodeSeparateNALBufs(bool) error
	SetEncodeH264LowLatency(bool) error
}

type MMALCameraParameters interface {
	// Get Parameters
	CameraInfo() (MMALCameraInfo, error)
	CaptureExposureComp() (int32, error)
	OutputShift() (int32, error)
	CCMShift() (int32, error)
	Rotation() (int32, error)
	CameraNum() (uint32, error)
	JPEGQFactor() (uint32, error)
	ISO() (uint32, error)
	MinISO() (uint32, error)
	CustomSensorConfig() (uint32, error)
	ShutterSpeed() (uint32, error)
	DPFConfig() (uint32, error)
	JPEGRestartInterval() (uint32, error)
	CameraISPBlockOverride() (uint32, error)
	BlackLevel() (uint32, error)
	EXIFDisable() (bool, error)
	Capture() (bool, error)
	DrawBoxFacesAndFocus() (bool, error)
	VideoStabilisation() (bool, error)
	EnableRAWCapture() (bool, error)
	EnableDPFFile() (bool, error)
	DPFFailIsFatal() (bool, error)
	HighDynamicRange() (bool, error)
	Antishake() (bool, error)
	BurstCapture() (bool, error)
	CaptureStatsPass() (bool, error)
	EnableRegisterFile() (bool, error)
	RegisterFailIsFatal() (bool, error)
	JPEGAttachLog() (bool, error)
	SWSharpenDisable() (bool, error)
	FlashRequired() (bool, error)
	SWSaturationDisable() (bool, error)
	VideoDenoise() (bool, error)
	StillsDenoise() (bool, error)
	ExposureComp() (MMALRationalNum, error)
	Sharpness() (MMALRationalNum, error)
	Contrast() (MMALRationalNum, error)
	Brightness() (MMALRationalNum, error)
	Saturation() (MMALRationalNum, error)
	AnalogGain() (MMALRationalNum, error)
	DigitalGain() (MMALRationalNum, error)
	MeteringMode() (MMALCameraMeteringMode, error)
	ExposureMode() (MMALCameraExposureMode, error)
	Annotation() (MMALCameraAnnotation, error)

	// Set Parameters
	SetCaptureExposureComp(value int32) error
	SetOutputShift(value int32) error
	SetCCMShift(value int32) error
	SetRotation(value int32) error
	SetCameraNum(value uint32) error
	SetJPEGQFactor(value uint32) error
	SetISO(value uint32) error
	SetMinISO(value uint32) error
	SetCustomSensorConfig(value uint32) error
	SetShutterSpeed(value uint32) error
	SetDPFConfig(value uint32) error
	SetJPEGRestartInterval(value uint32) error
	SetCameraISPBlockOverride(value uint32) error
	SetBlackLevel(value uint32) error
	SetEXIFDisable(value bool) error
	SetCapture(value bool) error
	SetDrawBoxFacesAndFocus(value bool) error
	SetVideoStabilisation(value bool) error
	SetEnableRAWCapture(value bool) error
	SetEnableDPFFile(value bool) error
	SetDPFFailIsFatal(value bool) error
	SetHighDynamicRange(value bool) error
	SetAntishake(value bool) error
	SetBurstCapture(value bool) error
	SetCaptureStatsPass(value bool) error
	SetEnableRegisterFile(value bool) error
	SetRegisterFailIsFatal(value bool) error
	SetJPEGAttachLog(value bool) error
	SetSWSharpenDisable(value bool) error
	SetFlashRequired(value bool) error
	SetSWSaturationDisable(value bool) error
	SetVideoDenoise(value bool) error
	SetStillsDenoise(value bool) error
	SetExposureComp(value MMALRationalNum) error
	SetSharpness(value MMALRationalNum) error
	SetContrast(value MMALRationalNum) error
	SetBrightness(value MMALRationalNum) error
	SetSaturation(value MMALRationalNum) error
	SetAnalogGain(value MMALRationalNum) error
	SetDigitalGain(value MMALRationalNum) error
	SetMeteringMode(MMALCameraMeteringMode) error
	SetExposureMode(MMALCameraExposureMode) error
	SetAnnotation(MMALCameraAnnotation) error
}

type MMALFormat interface {
	Type() MMALFormatType
	Bitrate() uint32
	SetBitrate(uint32)
	Encoding() (MMALEncodingType, MMALEncodingType)
	SetEncoding(MMALEncodingType)
	SetEncodingVariant(MMALEncodingType, MMALEncodingType)
}

type MMALVideoFormat interface {
	MMALFormat

	// Get and set video format parameters
	WidthHeight() (uint32, uint32)
	SetWidthHeight(uint32, uint32)
	Crop() MMALRect
	SetCrop(MMALRect)
	FrameRate() MMALRationalNum
	SetFrameRate(MMALRationalNum)
	PixelAspectRatio() MMALRationalNum
	SetPixelAspectRatio(MMALRationalNum)
	ColorSpace() MMALEncodingType
	SetColorSpace(MMALEncodingType)
}

type MMALAudioFormat interface {
	MMALFormat

	// Get and set audio format parameters
	/* TODO
	Channels() uint32
	SetChannels(uint32)
	SampleRate() uint32
	SetSampleRate(uint32)
	BitsPerSample() uint32
	SetBitsPerSample(uint32)
	BlockAlign() uint32
	SetBlockAlign(uint32)*/
}

type MMALSubpictureFormat interface {
	MMALFormat
	/* TODO
	XYOffset() (uint32, uint32)
	SetXYOffset(uint32, uint32)
	*/
}

type MMALDisplayRegion interface {
	// Get properties
	Display() uint16
	FullScreen() bool
	Layer() int16
	Alpha() uint8
	Transform() MMALDisplayTransform
	NoAspect() bool
	Mode() MMALDisplayMode
	CopyProtect() bool
	DestRect() MMALRect
	SrcRect() MMALRect

	// Set properties
	SetFullScreen(bool)
	SetLayer(int16)
	SetAlpha(uint8)
	SetTransform(MMALDisplayTransform)
	SetNoAspect(bool)
	SetMode(MMALDisplayMode)
	SetCopyProtect(bool)
	SetDestRect(MMALRect)
	SetSrcRect(MMALRect)
}

type MMALCameraInfo interface {
	Cameras() []MMALCamera
	Flashes() []MMALCameraFlashType
}

type MMALCamera interface {
	Id() uint32
	Name() string
	Size() (uint32, uint32)
	LensPresent() bool
}

type MMALCameraAnnotation interface {
	// Get parameters
	ShowShutter() bool
	ShowAnalogGain() bool
	ShowLens() bool
	ShowCAF() bool
	ShowMotion() bool
	ShowFrameNum() bool
	BackgroundColor() (uint8, uint8, uint8)
	TextColor() (uint8, uint8, uint8)
	TextSize() uint8
	Text() string
	TextJustify() MMALTextJustify
	TextOffset() (uint32, uint32)

	// Set parameters
	SetShowShutter(bool)
	SetShowAnalogGain(bool)
	SetShowLens(bool)
	SetShowCAF(bool)
	SetShowMotion(bool)
	SetShowFrameNum(bool)
	SetBackgroundColor(y, u, v uint8)
	SetTextColor(y, u, v uint8)
	SetText(string)
	SetTextSize(uint8)
	SetTextJustify(MMALTextJustify)
	SetTextOffset(x, y uint32)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_FORMAT_UNKNOWN MMALFormatType = iota
	MMAL_FORMAT_CONTROL
	MMAL_FORMAT_AUDIO
	MMAL_FORMAT_VIDEO
	MMAL_FORMAT_SUBPICTURE
	MMAL_FORMAT_MAX = MMAL_FORMAT_SUBPICTURE
)

const (
	MMAL_DISPLAY_TRANSFORM_NONE MMALDisplayTransform = iota
	MMAL_DISPLAY_TRANSFORM_MIRROR
	MMAL_DISPLAY_TRANSFORM_ROT180_MIRROR
	MMAL_DISPLAY_TRANSFORM_ROT180
	MMAL_DISPLAY_TRANSFORM_ROT90_MIRROR
	MMAL_DISPLAY_TRANSFORM_ROT270
	MMAL_DISPLAY_TRANSFORM_ROT90
	MMAL_DISPLAY_TRANSFORM_ROT270_MIRROR
	MMAL_DISPLAY_TRANSFORM_MAX = MMAL_DISPLAY_TRANSFORM_ROT270_MIRROR
)

const (
	MMAL_DISPLAY_MODE_FILL MMALDisplayMode = iota
	MMAL_DISPLAY_MODE_LETTERBOX
	MMAL_DISPLAY_MODE_STEREO_LEFT_TO_LEFT
	MMAL_DISPLAY_MODE_STEREO_TOP_TO_TOP
	MMAL_DISPLAY_MODE_STEREO_LEFT_TO_TOP
	MMAL_DISPLAY_MODE_STEREO_TOP_TO_LEFT
	MMAL_DISPLAY_MODE_MAX = MMAL_DISPLAY_MODE_STEREO_TOP_TO_LEFT
)

const (
	// MMALPortConnectionFlags
	MMAL_CONNECTION_FLAG_TUNNELLING               MMALPortConnectionFlags = 0x0001 // The connection is tunnelled. Buffer headers do not transit via the client but directly from the output port to the input port.
	MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT      MMALPortConnectionFlags = 0x0002 // Force the pool of buffer headers used by the connection to be allocated on the input port.
	MMAL_CONNECTION_FLAG_ALLOCATION_ON_OUTPUT     MMALPortConnectionFlags = 0x0004 // Force the pool of buffer headers used by the connection to be allocated on the output port.
	MMAL_CONNECTION_FLAG_KEEP_BUFFER_REQUIREMENTS MMALPortConnectionFlags = 0x0008 // Specify that the connection should not modify the buffer requirements.
	MMAL_CONNECTION_FLAG_DIRECT                   MMALPortConnectionFlags = 0x0010 // The connection is flagged as direct. This doesn't change the behaviour of the connection itself but is used by the the graph utility to specify that the buffer should be sent to the input port from with the port callback.
	MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS        MMALPortConnectionFlags = 0x0020 // Specify that the connection should not modify the port formats.
	MMAL_CONNECTION_FLAG_MIN                                              = MMAL_CONNECTION_FLAG_TUNNELLING
	MMAL_CONNECTION_FLAG_MAX                                              = MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS
)

const (
	// MMALVideoEncProfile
	MMAL_VIDEO_PROFILE_H263_BASELINE MMALVideoEncProfile = iota
	MMAL_VIDEO_PROFILE_H263_H320CODING
	MMAL_VIDEO_PROFILE_H263_BACKWARDCOMPATIBLE
	MMAL_VIDEO_PROFILE_H263_ISWV2
	MMAL_VIDEO_PROFILE_H263_ISWV3
	MMAL_VIDEO_PROFILE_H263_HIGHCOMPRESSION
	MMAL_VIDEO_PROFILE_H263_INTERNET
	MMAL_VIDEO_PROFILE_H263_INTERLACE
	MMAL_VIDEO_PROFILE_H263_HIGHLATENCY
	MMAL_VIDEO_PROFILE_MP4V_SIMPLE
	MMAL_VIDEO_PROFILE_MP4V_SIMPLESCALABLE
	MMAL_VIDEO_PROFILE_MP4V_CORE
	MMAL_VIDEO_PROFILE_MP4V_MAIN
	MMAL_VIDEO_PROFILE_MP4V_NBIT
	MMAL_VIDEO_PROFILE_MP4V_SCALABLETEXTURE
	MMAL_VIDEO_PROFILE_MP4V_SIMPLEFACE
	MMAL_VIDEO_PROFILE_MP4V_SIMPLEFBA
	MMAL_VIDEO_PROFILE_MP4V_BASICANIMATED
	MMAL_VIDEO_PROFILE_MP4V_HYBRID
	MMAL_VIDEO_PROFILE_MP4V_ADVANCEDREALTIME
	MMAL_VIDEO_PROFILE_MP4V_CORESCALABLE
	MMAL_VIDEO_PROFILE_MP4V_ADVANCEDCODING
	MMAL_VIDEO_PROFILE_MP4V_ADVANCEDCORE
	MMAL_VIDEO_PROFILE_MP4V_ADVANCEDSCALABLE
	MMAL_VIDEO_PROFILE_MP4V_ADVANCEDSIMPLE
	MMAL_VIDEO_PROFILE_H264_BASELINE
	MMAL_VIDEO_PROFILE_H264_MAIN
	MMAL_VIDEO_PROFILE_H264_EXTENDED
	MMAL_VIDEO_PROFILE_H264_HIGH
	MMAL_VIDEO_PROFILE_H264_HIGH10
	MMAL_VIDEO_PROFILE_H264_HIGH422
	MMAL_VIDEO_PROFILE_H264_HIGH444
	MMAL_VIDEO_PROFILE_H264_CONSTRAINED_BASELINE
	MMAL_VIDEO_PROFILE_MIN = MMAL_VIDEO_PROFILE_H263_BASELINE
	MMAL_VIDEO_PROFILE_MAX = MMAL_VIDEO_PROFILE_H264_CONSTRAINED_BASELINE
)

const (
	// MMALVideoEncLevel
	MMAL_VIDEO_LEVEL_H263_10 MMALVideoEncLevel = iota
	MMAL_VIDEO_LEVEL_H263_20
	MMAL_VIDEO_LEVEL_H263_30
	MMAL_VIDEO_LEVEL_H263_40
	MMAL_VIDEO_LEVEL_H263_45
	MMAL_VIDEO_LEVEL_H263_50
	MMAL_VIDEO_LEVEL_H263_60
	MMAL_VIDEO_LEVEL_H263_70
	MMAL_VIDEO_LEVEL_MP4V_0
	MMAL_VIDEO_LEVEL_MP4V_0b
	MMAL_VIDEO_LEVEL_MP4V_1
	MMAL_VIDEO_LEVEL_MP4V_2
	MMAL_VIDEO_LEVEL_MP4V_3
	MMAL_VIDEO_LEVEL_MP4V_4
	MMAL_VIDEO_LEVEL_MP4V_4a
	MMAL_VIDEO_LEVEL_MP4V_5
	MMAL_VIDEO_LEVEL_MP4V_6
	MMAL_VIDEO_LEVEL_H264_1
	MMAL_VIDEO_LEVEL_H264_1b
	MMAL_VIDEO_LEVEL_H264_11
	MMAL_VIDEO_LEVEL_H264_12
	MMAL_VIDEO_LEVEL_H264_13
	MMAL_VIDEO_LEVEL_H264_2
	MMAL_VIDEO_LEVEL_H264_21
	MMAL_VIDEO_LEVEL_H264_22
	MMAL_VIDEO_LEVEL_H264_3
	MMAL_VIDEO_LEVEL_H264_31
	MMAL_VIDEO_LEVEL_H264_32
	MMAL_VIDEO_LEVEL_H264_4
	MMAL_VIDEO_LEVEL_H264_41
	MMAL_VIDEO_LEVEL_H264_42
	MMAL_VIDEO_LEVEL_H264_5
	MMAL_VIDEO_LEVEL_H264_51
	MMAL_VIDEO_LEVEL_MIN = MMAL_VIDEO_LEVEL_H263_10
	MMAL_VIDEO_LEVEL_MAX = MMAL_VIDEO_LEVEL_H264_51
)

const (
	MMAL_CAMERA_FLASH_TYPE_XENON MMALCameraFlashType = iota
	MMAL_CAMERA_FLASH_TYPE_LED
	MMAL_CAMERA_FLASH_TYPE_OTHER
)

const (
	MMAL_CAMERA_METERINGMODE_AVERAGE MMALCameraMeteringMode = iota
	MMAL_CAMERA_METERINGMODE_SPOT
	MMAL_CAMERA_METERINGMODE_BACKLIT
	MMAL_CAMERA_METERINGMODE_MATRIX
	MMAL_CAMERA_METERINGMODE_MAX = MMAL_CAMERA_METERINGMODE_MATRIX
)

const (
	MMAL_CAMERA_EXPOSUREMODE_OFF MMALCameraExposureMode = iota
	MMAL_CAMERA_EXPOSUREMODE_AUTO
	MMAL_CAMERA_EXPOSUREMODE_NIGHT
	MMAL_CAMERA_EXPOSUREMODE_NIGHTPREVIEW
	MMAL_CAMERA_EXPOSUREMODE_BACKLIGHT
	MMAL_CAMERA_EXPOSUREMODE_SPOTLIGHT
	MMAL_CAMERA_EXPOSUREMODE_SPORTS
	MMAL_CAMERA_EXPOSUREMODE_SNOW
	MMAL_CAMERA_EXPOSUREMODE_BEACH
	MMAL_CAMERA_EXPOSUREMODE_VERYLONG
	MMAL_CAMERA_EXPOSUREMODE_FIXEDFPS
	MMAL_CAMERA_EXPOSUREMODE_ANTISHAKE
	MMAL_CAMERA_EXPOSUREMODE_FIREWORKS
	MMAL_CAMERA_EXPOSUREMODE_MAX = MMAL_CAMERA_EXPOSUREMODE_FIREWORKS
)

const (
	MMAL_TEXT_JUSTIFY_CENTER MMALTextJustify = iota
	MMAL_TEXT_JUSTIFY_LEFT
	MMAL_TEXT_JUSTIFY_RIGHT
	MMAL_TEXT_JUSTIFY_CENTRE = MMAL_TEXT_JUSTIFY_CENTER
)

const (
	MMAL_BUFFER_FLAG_EOS                 MMALBufferFlag = (1 << iota)
	MMAL_BUFFER_FLAG_FRAME_START         MMALBufferFlag = (1 << iota)                                                 // Signals that the start of the current payload starts a frame
	MMAL_BUFFER_FLAG_FRAME_END           MMALBufferFlag = (1 << iota)                                                 // Signals that the end of the current payload ends a frame
	MMAL_BUFFER_FLAG_KEYFRAME            MMALBufferFlag = (1 << iota)                                                 // Signals that the current payload is a keyframe (i.e. self decodable)
	MMAL_BUFFER_FLAG_DISCONTINUITY       MMALBufferFlag = (1 << iota)                                                 // Signals a discontinuity in the stream of data (e.g. after a seek). Can be used for instance by a decoder to reset its state
	MMAL_BUFFER_FLAG_CONFIG              MMALBufferFlag = (1 << iota)                                                 // Signals a buffer containing some kind of config data for the component (e.g. codec config data)
	MMAL_BUFFER_FLAG_ENCRYPTED           MMALBufferFlag = (1 << iota)                                                 // Signals an encrypted payload
	MMAL_BUFFER_FLAG_CODECSIDEINFO       MMALBufferFlag = (1 << iota)                                                 // Signals a buffer containing side information
	MMAL_BUFFER_FLAGS_SNAPSHOT           MMALBufferFlag = (1 << iota)                                                 // Signals a buffer which is the snapshot/postview image from a stills capture
	MMAL_BUFFER_FLAG_CORRUPTED           MMALBufferFlag = (1 << iota)                                                 // Signals a buffer which contains data known to be corrupted
	MMAL_BUFFER_FLAG_TRANSMISSION_FAILED MMALBufferFlag = (1 << iota)                                                 // Signals that a buffer failed to be transmitted
	MMAL_BUFFER_FLAG_DECODEONLY          MMALBufferFlag = (1 << iota)                                                 // Signals the output buffer won't be used, just update reference frames
	MMAL_BUFFER_FLAG_NAL_END             MMALBufferFlag = (1 << iota)                                                 // Signals that the end of the current payload ends a NAL
	MMAL_BUFFER_FLAG_FRAME               MMALBufferFlag = (MMAL_BUFFER_FLAG_FRAME_START | MMAL_BUFFER_FLAG_FRAME_END) // Signals that the current payload contains only complete frames (1 or more)
	MMAL_BUFFER_FLAG_MIN                                = MMAL_BUFFER_FLAG_EOS
	MMAL_BUFFER_FLAG_MAX                                = MMAL_BUFFER_FLAG_NAL_END
)

////////////////////////////////////////////////////////////////////////////////
// VIDEO ENCODINGS

var (
	MMAL_ENCODING_H264   = MMAL_FOURCC('H', '2', '6', '4')
	MMAL_ENCODING_MVC    = MMAL_FOURCC('M', 'V', 'C', ' ')
	MMAL_ENCODING_H263   = MMAL_FOURCC('H', '2', '6', '3')
	MMAL_ENCODING_MP4V   = MMAL_FOURCC('M', 'P', '4', 'V')
	MMAL_ENCODING_MP2V   = MMAL_FOURCC('M', 'P', '2', 'V')
	MMAL_ENCODING_MP1V   = MMAL_FOURCC('M', 'P', '1', 'V')
	MMAL_ENCODING_WMV3   = MMAL_FOURCC('W', 'M', 'V', '3')
	MMAL_ENCODING_WMV2   = MMAL_FOURCC('W', 'M', 'V', '2')
	MMAL_ENCODING_WMV1   = MMAL_FOURCC('W', 'M', 'V', '1')
	MMAL_ENCODING_WVC1   = MMAL_FOURCC('W', 'V', 'C', '1')
	MMAL_ENCODING_VP8    = MMAL_FOURCC('V', 'P', '8', ' ')
	MMAL_ENCODING_VP7    = MMAL_FOURCC('V', 'P', '7', ' ')
	MMAL_ENCODING_VP6    = MMAL_FOURCC('V', 'P', '6', ' ')
	MMAL_ENCODING_THEORA = MMAL_FOURCC('T', 'H', 'E', 'O')
	MMAL_ENCODING_SPARK  = MMAL_FOURCC('S', 'P', 'R', 'K')
	MMAL_ENCODING_MJPEG  = MMAL_FOURCC('M', 'J', 'P', 'G')
)

////////////////////////////////////////////////////////////////////////////////
// IMAGE ENCODINGS

var (
	MMAL_ENCODING_JPEG = MMAL_FOURCC('J', 'P', 'E', 'G')
	MMAL_ENCODING_GIF  = MMAL_FOURCC('G', 'I', 'F', ' ')
	MMAL_ENCODING_PNG  = MMAL_FOURCC('P', 'N', 'G', ' ')
	MMAL_ENCODING_PPM  = MMAL_FOURCC('P', 'P', 'M', ' ')
	MMAL_ENCODING_TGA  = MMAL_FOURCC('T', 'G', 'A', ' ')
	MMAL_ENCODING_BMP  = MMAL_FOURCC('B', 'M', 'P', ' ')
)

////////////////////////////////////////////////////////////////////////////////
// UNCOMPRESSED ENCODINGS

var (
	MMAL_ENCODING_I420        = MMAL_FOURCC('I', '4', '2', '0')
	MMAL_ENCODING_I420_SLICE  = MMAL_FOURCC('S', '4', '2', '0')
	MMAL_ENCODING_YV12        = MMAL_FOURCC('Y', 'V', '1', '2')
	MMAL_ENCODING_I422        = MMAL_FOURCC('I', '4', '2', '2')
	MMAL_ENCODING_I422_SLICE  = MMAL_FOURCC('S', '4', '2', '2')
	MMAL_ENCODING_YUYV        = MMAL_FOURCC('Y', 'U', 'Y', 'V')
	MMAL_ENCODING_YVYU        = MMAL_FOURCC('Y', 'V', 'Y', 'U')
	MMAL_ENCODING_UYVY        = MMAL_FOURCC('U', 'Y', 'V', 'Y')
	MMAL_ENCODING_VYUY        = MMAL_FOURCC('V', 'Y', 'U', 'Y')
	MMAL_ENCODING_NV12        = MMAL_FOURCC('N', 'V', '1', '2')
	MMAL_ENCODING_NV21        = MMAL_FOURCC('N', 'V', '2', '1')
	MMAL_ENCODING_ARGB        = MMAL_FOURCC('A', 'R', 'G', 'B')
	MMAL_ENCODING_ARGB_SLICE  = MMAL_FOURCC('a', 'r', 'g', 'b')
	MMAL_ENCODING_RGBA        = MMAL_FOURCC('R', 'G', 'B', 'A')
	MMAL_ENCODING_RGBA_SLICE  = MMAL_FOURCC('r', 'g', 'b', 'a')
	MMAL_ENCODING_ABGR        = MMAL_FOURCC('A', 'B', 'G', 'R')
	MMAL_ENCODING_ABGR_SLICE  = MMAL_FOURCC('a', 'b', 'g', 'r')
	MMAL_ENCODING_BGRA        = MMAL_FOURCC('B', 'G', 'R', 'A')
	MMAL_ENCODING_BGRA_SLICE  = MMAL_FOURCC('b', 'g', 'r', 'a')
	MMAL_ENCODING_RGB16       = MMAL_FOURCC('R', 'G', 'B', '2')
	MMAL_ENCODING_RGB16_SLICE = MMAL_FOURCC('r', 'g', 'b', '2')
	MMAL_ENCODING_RGB24       = MMAL_FOURCC('R', 'G', 'B', '3')
	MMAL_ENCODING_RGB24_SLICE = MMAL_FOURCC('r', 'g', 'b', '3')
	MMAL_ENCODING_RGB32       = MMAL_FOURCC('R', 'G', 'B', '4')
	MMAL_ENCODING_RGB32_SLICE = MMAL_FOURCC('r', 'g', 'b', '4')
	MMAL_ENCODING_BGR16       = MMAL_FOURCC('B', 'G', 'R', '2')
	MMAL_ENCODING_BGR16_SLICE = MMAL_FOURCC('b', 'g', 'r', '2')
	MMAL_ENCODING_BGR24       = MMAL_FOURCC('B', 'G', 'R', '3')
	MMAL_ENCODING_BGR24_SLICE = MMAL_FOURCC('b', 'g', 'r', '3')
	MMAL_ENCODING_BGR32       = MMAL_FOURCC('B', 'G', 'R', '4')
	MMAL_ENCODING_BGR32_SLICE = MMAL_FOURCC('b', 'g', 'r', '4')
)

////////////////////////////////////////////////////////////////////////////////
// CONTROL EVENTS

var (
	MMAL_EVENT_ERROR             = MMAL_FOURCC('E', 'R', 'R', 'O')
	MMAL_EVENT_EOS               = MMAL_FOURCC('E', 'E', 'O', 'S') // End-of-stream event. Data contains a MMAL_EVENT_END_OF_STREAM_T
	MMAL_EVENT_FORMAT_CHANGED    = MMAL_FOURCC('E', 'F', 'C', 'H') // Format changed event. Data contains a MMAL_EVENT_FORMAT_CHANGED_T
	MMAL_EVENT_PARAMETER_CHANGED = MMAL_FOURCC('E', 'P', 'C', 'H') // Parameter changed event. Data contains a MMAL_EVENT_PARAMETER_CHANGED_T
)

func MMAL_FOURCC(a, b, c, d uint8) MMALEncodingType {
	return MMALEncodingType(uint32(a) | uint32(b)<<8 | uint32(c)<<16 | uint32(d)<<24)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (j MMALTextJustify) String() string {
	switch j {
	case MMAL_TEXT_JUSTIFY_CENTER:
		return "MMAL_TEXT_JUSTIFY_CENTER"
	case MMAL_TEXT_JUSTIFY_LEFT:
		return "MMAL_TEXT_JUSTIFY_LEFT"
	case MMAL_TEXT_JUSTIFY_RIGHT:
		return "MMAL_TEXT_JUSTIFY_RIGHT"
	default:
		return "[?? Unknown MMALTextJustify value]"
	}
}

func (n MMALRationalNum) String() string {
	return fmt.Sprintf("(%v/%v)", n.Num, n.Den)
}

func (t MMALFormatType) String() string {
	switch t {
	case MMAL_FORMAT_UNKNOWN:
		return "MMAL_FORMAT_UNKNOWN"
	case MMAL_FORMAT_CONTROL:
		return "MMAL_FORMAT_CONTROL"
	case MMAL_FORMAT_AUDIO:
		return "MMAL_FORMAT_AUDIO"
	case MMAL_FORMAT_VIDEO:
		return "MMAL_FORMAT_VIDEO"
	case MMAL_FORMAT_SUBPICTURE:
		return "MMAL_FORMAT_SUBPICTURE"
	default:
		return "[?? Unknown MMALFormatType value]"
	}
}

func (e MMALEncodingType) String() string {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(e))
	return "'" + string(buf) + "'"
}

func (t MMALDisplayTransform) String() string {
	switch t {
	case MMAL_DISPLAY_TRANSFORM_NONE:
		return "MMAL_DISPLAY_TRANSFORM_NONE"
	case MMAL_DISPLAY_TRANSFORM_MIRROR:
		return "MMAL_DISPLAY_TRANSFORM_MIRROR"
	case MMAL_DISPLAY_TRANSFORM_ROT180_MIRROR:
		return "MMAL_DISPLAY_TRANSFORM_ROT180_MIRROR"
	case MMAL_DISPLAY_TRANSFORM_ROT180:
		return "MMAL_DISPLAY_TRANSFORM_ROT180"
	case MMAL_DISPLAY_TRANSFORM_ROT90_MIRROR:
		return "MMAL_DISPLAY_TRANSFORM_ROT90_MIRROR"
	case MMAL_DISPLAY_TRANSFORM_ROT270:
		return "MMAL_DISPLAY_TRANSFORM_ROT270"
	case MMAL_DISPLAY_TRANSFORM_ROT90:
		return "MMAL_DISPLAY_TRANSFORM_ROT90"
	case MMAL_DISPLAY_TRANSFORM_ROT270_MIRROR:
		return "MMAL_DISPLAY_TRANSFORM_ROT270_MIRROR"
	default:
		return "[?? Unknown MMALDisplayTransform value]"
	}
}

func (m MMALDisplayMode) String() string {
	switch m {
	case MMAL_DISPLAY_MODE_FILL:
		return "MMAL_DISPLAY_MODE_FILL"
	case MMAL_DISPLAY_MODE_LETTERBOX:
		return "MMAL_DISPLAY_MODE_LETTERBOX"
	case MMAL_DISPLAY_MODE_STEREO_LEFT_TO_LEFT:
		return "MMAL_DISPLAY_MODE_STEREO_LEFT_TO_LEFT"
	case MMAL_DISPLAY_MODE_STEREO_TOP_TO_TOP:
		return "MMAL_DISPLAY_MODE_STEREO_TOP_TO_TOP"
	case MMAL_DISPLAY_MODE_STEREO_LEFT_TO_TOP:
		return "MMAL_DISPLAY_MODE_STEREO_LEFT_TO_TOP"
	case MMAL_DISPLAY_MODE_STEREO_TOP_TO_LEFT:
		return "MMAL_DISPLAY_MODE_STEREO_TOP_TO_LEFT"
	default:
		return "[?? Unknown MMALDisplayMode value]"
	}
}

func (c MMALPortConnectionFlags) String() string {
	parts := ""
	for flag := MMAL_CONNECTION_FLAG_MIN; flag <= MMAL_CONNECTION_FLAG_MAX; flag <<= 1 {
		if c&flag == 0 {
			continue
		}
		switch flag {
		case MMAL_CONNECTION_FLAG_TUNNELLING:
			parts += "|" + "MMAL_CONNECTION_FLAG_TUNNELLING"
		case MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT:
			parts += "|" + "MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT"
		case MMAL_CONNECTION_FLAG_ALLOCATION_ON_OUTPUT:
			parts += "|" + "MMAL_CONNECTION_FLAG_ALLOCATION_ON_OUTPUT"
		case MMAL_CONNECTION_FLAG_KEEP_BUFFER_REQUIREMENTS:
			parts += "|" + "MMAL_CONNECTION_FLAG_KEEP_BUFFER_REQUIREMENTS"
		case MMAL_CONNECTION_FLAG_DIRECT:
			parts += "|" + "MMAL_CONNECTION_FLAG_DIRECT"
		case MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS:
			parts += "|" + "MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS"
		default:
			parts += "|" + "[?? Invalid MMALPortConnectionFlags value]"
		}
	}
	return strings.Trim(parts, "|")
}

func (p MMALVideoEncProfile) String() string {
	switch p {
	case MMAL_VIDEO_PROFILE_H263_BASELINE:
		return "MMAL_VIDEO_PROFILE_H263_BASELINE"
	case MMAL_VIDEO_PROFILE_H263_H320CODING:
		return "MMAL_VIDEO_PROFILE_H263_H320CODING"
	case MMAL_VIDEO_PROFILE_H263_BACKWARDCOMPATIBLE:
		return "MMAL_VIDEO_PROFILE_H263_BACKWARDCOMPATIBLE"
	case MMAL_VIDEO_PROFILE_H263_ISWV2:
		return "MMAL_VIDEO_PROFILE_H263_ISWV2"
	case MMAL_VIDEO_PROFILE_H263_ISWV3:
		return "MMAL_VIDEO_PROFILE_H263_ISWV3"
	case MMAL_VIDEO_PROFILE_H263_HIGHCOMPRESSION:
		return "MMAL_VIDEO_PROFILE_H263_HIGHCOMPRESSION"
	case MMAL_VIDEO_PROFILE_H263_INTERNET:
		return "MMAL_VIDEO_PROFILE_H263_INTERNET"
	case MMAL_VIDEO_PROFILE_H263_INTERLACE:
		return "MMAL_VIDEO_PROFILE_H263_INTERLACE"
	case MMAL_VIDEO_PROFILE_H263_HIGHLATENCY:
		return "MMAL_VIDEO_PROFILE_H263_HIGHLATENCY"
	case MMAL_VIDEO_PROFILE_MP4V_SIMPLE:
		return "MMAL_VIDEO_PROFILE_MP4V_SIMPLE"
	case MMAL_VIDEO_PROFILE_MP4V_SIMPLESCALABLE:
		return "MMAL_VIDEO_PROFILE_MP4V_SIMPLESCALABLE"
	case MMAL_VIDEO_PROFILE_MP4V_CORE:
		return "MMAL_VIDEO_PROFILE_MP4V_CORE"
	case MMAL_VIDEO_PROFILE_MP4V_MAIN:
		return "MMAL_VIDEO_PROFILE_MP4V_MAIN"
	case MMAL_VIDEO_PROFILE_MP4V_NBIT:
		return "MMAL_VIDEO_PROFILE_MP4V_NBIT"
	case MMAL_VIDEO_PROFILE_MP4V_SCALABLETEXTURE:
		return "MMAL_VIDEO_PROFILE_MP4V_SCALABLETEXTURE"
	case MMAL_VIDEO_PROFILE_MP4V_SIMPLEFACE:
		return "MMAL_VIDEO_PROFILE_MP4V_SIMPLEFACE"
	case MMAL_VIDEO_PROFILE_MP4V_SIMPLEFBA:
		return "MMAL_VIDEO_PROFILE_MP4V_SIMPLEFBA"
	case MMAL_VIDEO_PROFILE_MP4V_BASICANIMATED:
		return "MMAL_VIDEO_PROFILE_MP4V_BASICANIMATED"
	case MMAL_VIDEO_PROFILE_MP4V_HYBRID:
		return "MMAL_VIDEO_PROFILE_MP4V_HYBRID"
	case MMAL_VIDEO_PROFILE_MP4V_ADVANCEDREALTIME:
		return "MMAL_VIDEO_PROFILE_MP4V_ADVANCEDREALTIME"
	case MMAL_VIDEO_PROFILE_MP4V_CORESCALABLE:
		return "MMAL_VIDEO_PROFILE_MP4V_CORESCALABLE"
	case MMAL_VIDEO_PROFILE_MP4V_ADVANCEDCODING:
		return "MMAL_VIDEO_PROFILE_MP4V_ADVANCEDCODING"
	case MMAL_VIDEO_PROFILE_MP4V_ADVANCEDCORE:
		return "MMAL_VIDEO_PROFILE_MP4V_ADVANCEDCORE"
	case MMAL_VIDEO_PROFILE_MP4V_ADVANCEDSCALABLE:
		return "MMAL_VIDEO_PROFILE_MP4V_ADVANCEDSCALABLE"
	case MMAL_VIDEO_PROFILE_MP4V_ADVANCEDSIMPLE:
		return "MMAL_VIDEO_PROFILE_MP4V_ADVANCEDSIMPLE"
	case MMAL_VIDEO_PROFILE_H264_BASELINE:
		return "MMAL_VIDEO_PROFILE_H264_BASELINE"
	case MMAL_VIDEO_PROFILE_H264_MAIN:
		return "MMAL_VIDEO_PROFILE_H264_MAIN"
	case MMAL_VIDEO_PROFILE_H264_EXTENDED:
		return "MMAL_VIDEO_PROFILE_H264_EXTENDED"
	case MMAL_VIDEO_PROFILE_H264_HIGH:
		return "MMAL_VIDEO_PROFILE_H264_HIGH"
	case MMAL_VIDEO_PROFILE_H264_HIGH10:
		return "MMAL_VIDEO_PROFILE_H264_HIGH10"
	case MMAL_VIDEO_PROFILE_H264_HIGH422:
		return "MMAL_VIDEO_PROFILE_H264_HIGH422"
	case MMAL_VIDEO_PROFILE_H264_HIGH444:
		return "MMAL_VIDEO_PROFILE_H264_HIGH444"
	case MMAL_VIDEO_PROFILE_H264_CONSTRAINED_BASELINE:
		return "MMAL_VIDEO_PROFILE_H264_CONSTRAINED_BASELINE"
	default:
		return "[?? Invalid MMALVideoEncProfile value]"
	}
}

func (l MMALVideoEncLevel) String() string {
	switch l {
	case MMAL_VIDEO_LEVEL_H263_10:
		return "MMAL_VIDEO_LEVEL_H263_10"
	case MMAL_VIDEO_LEVEL_H263_20:
		return "MMAL_VIDEO_LEVEL_H263_20"
	case MMAL_VIDEO_LEVEL_H263_30:
		return "MMAL_VIDEO_LEVEL_H263_30"
	case MMAL_VIDEO_LEVEL_H263_40:
		return "MMAL_VIDEO_LEVEL_H263_40"
	case MMAL_VIDEO_LEVEL_H263_45:
		return "MMAL_VIDEO_LEVEL_H263_45"
	case MMAL_VIDEO_LEVEL_H263_50:
		return "MMAL_VIDEO_LEVEL_H263_50"
	case MMAL_VIDEO_LEVEL_H263_60:
		return "MMAL_VIDEO_LEVEL_H263_60"
	case MMAL_VIDEO_LEVEL_H263_70:
		return "MMAL_VIDEO_LEVEL_H263_70"
	case MMAL_VIDEO_LEVEL_MP4V_0:
		return "MMAL_VIDEO_LEVEL_MP4V_0"
	case MMAL_VIDEO_LEVEL_MP4V_0b:
		return "MMAL_VIDEO_LEVEL_MP4V_0b"
	case MMAL_VIDEO_LEVEL_MP4V_1:
		return "MMAL_VIDEO_LEVEL_MP4V_1"
	case MMAL_VIDEO_LEVEL_MP4V_2:
		return "MMAL_VIDEO_LEVEL_MP4V_2"
	case MMAL_VIDEO_LEVEL_MP4V_3:
		return "MMAL_VIDEO_LEVEL_MP4V_3"
	case MMAL_VIDEO_LEVEL_MP4V_4:
		return "MMAL_VIDEO_LEVEL_MP4V_4"
	case MMAL_VIDEO_LEVEL_MP4V_4a:
		return "MMAL_VIDEO_LEVEL_MP4V_4a"
	case MMAL_VIDEO_LEVEL_MP4V_5:
		return "MMAL_VIDEO_LEVEL_MP4V_5"
	case MMAL_VIDEO_LEVEL_MP4V_6:
		return "MMAL_VIDEO_LEVEL_MP4V_6"
	case MMAL_VIDEO_LEVEL_H264_1:
		return "MMAL_VIDEO_LEVEL_H264_1"
	case MMAL_VIDEO_LEVEL_H264_1b:
		return "MMAL_VIDEO_LEVEL_H264_1b"
	case MMAL_VIDEO_LEVEL_H264_11:
		return "MMAL_VIDEO_LEVEL_H264_11"
	case MMAL_VIDEO_LEVEL_H264_12:
		return "MMAL_VIDEO_LEVEL_H264_12"
	case MMAL_VIDEO_LEVEL_H264_13:
		return "MMAL_VIDEO_LEVEL_H264_13"
	case MMAL_VIDEO_LEVEL_H264_2:
		return "MMAL_VIDEO_LEVEL_H264_2"
	case MMAL_VIDEO_LEVEL_H264_21:
		return "MMAL_VIDEO_LEVEL_H264_21"
	case MMAL_VIDEO_LEVEL_H264_22:
		return "MMAL_VIDEO_LEVEL_H264_22"
	case MMAL_VIDEO_LEVEL_H264_3:
		return "MMAL_VIDEO_LEVEL_H264_3"
	case MMAL_VIDEO_LEVEL_H264_31:
		return "MMAL_VIDEO_LEVEL_H264_31"
	case MMAL_VIDEO_LEVEL_H264_32:
		return "MMAL_VIDEO_LEVEL_H264_32"
	case MMAL_VIDEO_LEVEL_H264_4:
		return "MMAL_VIDEO_LEVEL_H264_4"
	case MMAL_VIDEO_LEVEL_H264_41:
		return "MMAL_VIDEO_LEVEL_H264_41"
	case MMAL_VIDEO_LEVEL_H264_42:
		return "MMAL_VIDEO_LEVEL_H264_42"
	case MMAL_VIDEO_LEVEL_H264_5:
		return "MMAL_VIDEO_LEVEL_H264_5"
	default:
		return "[?? Invalid MMALVideoEncLevel value]"
	}
}

func (f MMALCameraFlashType) String() string {
	switch f {
	case MMAL_CAMERA_FLASH_TYPE_XENON:
		return "MMAL_CAMERA_FLASH_TYPE_XENON"
	case MMAL_CAMERA_FLASH_TYPE_LED:
		return "MMAL_CAMERA_FLASH_TYPE_LED"
	case MMAL_CAMERA_FLASH_TYPE_OTHER:
		return "MMAL_CAMERA_FLASH_TYPE_OTHER"
	default:
		return "[?? Invalid MMALCameraFlashType value]"
	}
}

func (m MMALCameraMeteringMode) String() string {
	switch m {
	case MMAL_CAMERA_METERINGMODE_AVERAGE:
		return "MMAL_CAMERA_METERINGMODE_AVERAGE"
	case MMAL_CAMERA_METERINGMODE_SPOT:
		return "MMAL_CAMERA_METERINGMODE_SPOT"
	case MMAL_CAMERA_METERINGMODE_BACKLIT:
		return "MMAL_CAMERA_METERINGMODE_BACKLIT"
	case MMAL_CAMERA_METERINGMODE_MATRIX:
		return "MMAL_CAMERA_METERINGMODE_MATRIX"
	default:
		return "[?? Invalid MMALCameraMeteringMode value]"
	}
}

func (m MMALCameraExposureMode) String() string {
	switch m {
	case MMAL_CAMERA_EXPOSUREMODE_OFF:
		return "MMAL_CAMERA_EXPOSUREMODE_OFF"
	case MMAL_CAMERA_EXPOSUREMODE_AUTO:
		return "MMAL_CAMERA_EXPOSUREMODE_AUTO"
	case MMAL_CAMERA_EXPOSUREMODE_NIGHT:
		return "MMAL_CAMERA_EXPOSUREMODE_NIGHT"
	case MMAL_CAMERA_EXPOSUREMODE_NIGHTPREVIEW:
		return "MMAL_CAMERA_EXPOSUREMODE_NIGHTPREVIEW"
	case MMAL_CAMERA_EXPOSUREMODE_BACKLIGHT:
		return "MMAL_CAMERA_EXPOSUREMODE_BACKLIGHT"
	case MMAL_CAMERA_EXPOSUREMODE_SPOTLIGHT:
		return "MMAL_CAMERA_EXPOSUREMODE_SPOTLIGHT"
	case MMAL_CAMERA_EXPOSUREMODE_SPORTS:
		return "MMAL_CAMERA_EXPOSUREMODE_SPORTS"
	case MMAL_CAMERA_EXPOSUREMODE_SNOW:
		return "MMAL_CAMERA_EXPOSUREMODE_SNOW"
	case MMAL_CAMERA_EXPOSUREMODE_BEACH:
		return "MMAL_CAMERA_EXPOSUREMODE_BEACH"
	case MMAL_CAMERA_EXPOSUREMODE_VERYLONG:
		return "MMAL_CAMERA_EXPOSUREMODE_VERYLONG"
	case MMAL_CAMERA_EXPOSUREMODE_FIXEDFPS:
		return "MMAL_CAMERA_EXPOSUREMODE_FIXEDFPS"
	case MMAL_CAMERA_EXPOSUREMODE_ANTISHAKE:
		return "MMAL_CAMERA_EXPOSUREMODE_ANTISHAKE"
	case MMAL_CAMERA_EXPOSUREMODE_FIREWORKS:
		return "MMAL_CAMERA_EXPOSUREMODE_FIREWORKS"
	default:
		return "[?? Invalid MMALCameraExposureMode value]"
	}
}

func (f MMALBufferFlag) String() string {
	parts := ""
	for flag := MMAL_BUFFER_FLAG_MIN; flag <= MMAL_BUFFER_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case MMAL_BUFFER_FLAG_EOS:
			parts += "|" + "MMAL_BUFFER_FLAG_EOS"
		case MMAL_BUFFER_FLAG_FRAME_START:
			parts += "|" + "MMAL_BUFFER_FLAG_FRAME_START"
		case MMAL_BUFFER_FLAG_FRAME_END:
			parts += "|" + "MMAL_BUFFER_FLAG_FRAME_END"
		case MMAL_BUFFER_FLAG_KEYFRAME:
			parts += "|" + "MMAL_BUFFER_FLAG_KEYFRAME"
		case MMAL_BUFFER_FLAG_DISCONTINUITY:
			parts += "|" + "MMAL_BUFFER_FLAG_DISCONTINUITY"
		case MMAL_BUFFER_FLAG_CONFIG:
			parts += "|" + "MMAL_BUFFER_FLAG_CONFIG"
		case MMAL_BUFFER_FLAG_ENCRYPTED:
			parts += "|" + "MMAL_BUFFER_FLAG_ENCRYPTED"
		case MMAL_BUFFER_FLAG_CODECSIDEINFO:
			parts += "|" + "MMAL_BUFFER_FLAG_CODECSIDEINFO"
		case MMAL_BUFFER_FLAGS_SNAPSHOT:
			parts += "|" + "MMAL_BUFFER_FLAGS_SNAPSHOT"
		case MMAL_BUFFER_FLAG_CORRUPTED:
			parts += "|" + "MMAL_BUFFER_FLAG_CORRUPTED"
		case MMAL_BUFFER_FLAG_TRANSMISSION_FAILED:
			parts += "|" + "MMAL_BUFFER_FLAG_TRANSMISSION_FAILED"
		case MMAL_BUFFER_FLAG_DECODEONLY:
			parts += "|" + "MMAL_BUFFER_FLAG_DECODEONLY"
		case MMAL_BUFFER_FLAG_NAL_END:
			parts += "|" + "MMAL_BUFFER_FLAG_NAL_END"
		default:
			parts += "|" + "[?? Invalid MMALBufferFlag value]"
		}
	}
	return strings.Trim(parts, "|")
}
