//+build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	// Frameworks
	"fmt"
	"unsafe"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#include <vc_dispmanx.h>
#include <bcm_host.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	DX_DisplayId     uint16
	DX_DisplayHandle C.DISPMANX_DISPLAY_HANDLE_T
	DX_InputFormat   C.uint32_t
	DX_Transform     C.int
	DX_Rect          (*C.VC_RECT_T)
	DX_Resource      C.DISPMANX_RESOURCE_HANDLE_T
	DX_ImageType     C.VC_IMAGE_TYPE_T
	DX_Element       C.DISPMANX_ELEMENT_HANDLE_T
	DX_Update        C.DISPMANX_UPDATE_HANDLE_T
	DX_Protection    C.uint32_t
	DX_AlphaFlags    C.uint32_t
	DX_ClampMode     C.int
)

type DX_DisplayModeInfo struct {
	Handle      DX_DisplayHandle
	Size        DX_Size
	Transform   DX_Transform
	InputFormat DX_InputFormat
}

type DX_Point struct {
	X int32
	Y int32
}

type DX_Size struct {
	W uint32
	H uint32
}

type DX_Clamp struct {
	Mode    DX_ClampMode
	Flags   int
	Opacity uint32
	Mask    DX_Resource
}

type DX_Alpha struct {
	Flags   DX_AlphaFlags
	Opacity uint32
	Mask    DX_Resource
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DX_DISPLAY_NONE DX_DisplayHandle = 0
)

const (
	/* Success and failure conditions */
	DX_SUCCESS   = 0
	DX_INVALID   = -1
	DX_NO_HANDLE = 0
)

const (
	// DX_Transform values
	DX_TRANSFORM_NONE DX_Transform = iota
	DX_TRANSFORM_ROTATE_90
	DX_TRANSFORM_ROTATE_180
	DX_TRANSFORM_ROTATE_270
	DX_TRANSFORM_MAX = DX_TRANSFORM_ROTATE_270
)

const (
	// DX_InputFormat values
	DX_INPUT_FORMAT_INVALID DX_InputFormat = iota
	DX_INPUT_FORMAT_RGB888
	DX_INPUT_FORMAT_RGB565
	DX_INPUT_FORMAT_MAX = DX_INPUT_FORMAT_RGB565
)

const (
	// DX_DisplayId values
	DX_DISPLAYID_MAIN_LCD DX_DisplayId = iota
	DX_DISPLAYID_AUX_LCD
	DX_DISPLAYID_HDMI
	DX_DISPLAYID_SDTV
	DX_DISPLAYID_FORCE_LCD
	DX_DISPLAYID_FORCE_TV
	DX_DISPLAYID_FORCE_OTHER
	DX_DISPLAYID_MAX = DX_DISPLAYID_FORCE_OTHER
)

const (
	DX_IMAGE_TYPE_NONE DX_ImageType = iota
	DX_IMAGE_TYPE_RGB565
	DX_IMAGE_TYPE_1BPP
	DX_IMAGE_TYPE_YUV420
	DX_IMAGE_TYPE_48BPP
	DX_IMAGE_TYPE_RGB888
	DX_IMAGE_TYPE_8BPP
	DX_IMAGE_TYPE_4BPP          // 4bpp palettised image
	DX_IMAGE_TYPE_3D32          // A separated format of 16 colour/light shorts followed by 16 z values
	DX_IMAGE_TYPE_3D32B         // 16 colours followed by 16 z values
	DX_IMAGE_TYPE_3D32MAT       // A separated format of 16 material/colour/light shorts followed by 16 z values
	DX_IMAGE_TYPE_RGB2X9        // 32 bit format containing 18 bits of 6.6.6 RGB 9 bits per short
	DX_IMAGE_TYPE_RGB666        // 32-bit format holding 18 bits of 6.6.6 RGB
	DX_IMAGE_TYPE_PAL4_OBSOLETE // 4bpp palettised image with embedded palette
	DX_IMAGE_TYPE_PAL8_OBSOLETE // 8bpp palettised image with embedded palette
	DX_IMAGE_TYPE_RGBA32        // RGB888 with an alpha byte after each pixel */ /* xxx: isn't it BEFORE each pixel?
	DX_IMAGE_TYPE_YUV422        // a line of Y (32-byte padded) a line of U (16-byte padded) and a line of V (16-byte padded)
	DX_IMAGE_TYPE_RGBA565       // RGB565 with a transparent patch
	DX_IMAGE_TYPE_RGBA16        // Compressed (4444) version of RGBA32
	DX_IMAGE_TYPE_YUV_UV        // VCIII codec format
	DX_IMAGE_TYPE_TF_RGBA32     // VCIII T-format RGBA8888
	DX_IMAGE_TYPE_TF_RGBX32     // VCIII T-format RGBx8888
	DX_IMAGE_TYPE_TF_FLOAT      // VCIII T-format float
	DX_IMAGE_TYPE_TF_RGBA16     // VCIII T-format RGBA4444
	DX_IMAGE_TYPE_TF_RGBA5551   // VCIII T-format RGB5551
	DX_IMAGE_TYPE_TF_RGB565     // VCIII T-format RGB565
	DX_IMAGE_TYPE_TF_YA88       // VCIII T-format 8-bit luma and 8-bit alpha
	DX_IMAGE_TYPE_TF_BYTE       // VCIII T-format 8 bit generic sample
	DX_IMAGE_TYPE_TF_PAL8       // VCIII T-format 8-bit palette
	DX_IMAGE_TYPE_TF_PAL4       // VCIII T-format 4-bit palette
	DX_IMAGE_TYPE_TF_ETC1       // VCIII T-format Ericsson Texture Compressed
	DX_IMAGE_TYPE_BGR888        // RGB888 with R & B swapped
	DX_IMAGE_TYPE_BGR888_NP     // RGB888 with R & B swapped but with no pitch i.e. no padding after each row of pixels
	DX_IMAGE_TYPE_BAYER         // Bayer image extra defines which variant is being used
	DX_IMAGE_TYPE_CODEC         // General wrapper for codec images e.g. JPEG from camera
	DX_IMAGE_TYPE_YUV_UV32      // VCIII codec format
	DX_IMAGE_TYPE_TF_Y8         // VCIII T-format 8-bit luma
	DX_IMAGE_TYPE_TF_A8         // VCIII T-format 8-bit alpha
	DX_IMAGE_TYPE_TF_SHORT      // VCIII T-format 16-bit generic sample
	DX_IMAGE_TYPE_TF_1BPP       // VCIII T-format 1bpp black/white
	DX_IMAGE_TYPE_OPENGL
	DX_IMAGE_TYPE_YUV444I      // VCIII-B0 HVS YUV 4:4:4 interleaved samples
	DX_IMAGE_TYPE_YUV422PLANAR // Y U & V planes separately (DX_IMAGE_TYPE_YUV422 has them interleaved on a per line basis)
	DX_IMAGE_TYPE_ARGB8888     // 32bpp with 8bit alpha at MS byte with R G B (LS byte)
	DX_IMAGE_TYPE_XRGB8888     // 32bpp with 8bit unused at MS byte with R G B (LS byte)
	DX_IMAGE_TYPE_YUV422YUYV   // interleaved 8 bit samples of Y U Y V
	DX_IMAGE_TYPE_YUV422YVYU   // interleaved 8 bit samples of Y V Y U
	DX_IMAGE_TYPE_YUV422UYVY   // interleaved 8 bit samples of U Y V Y
	DX_IMAGE_TYPE_YUV422VYUY   // interleaved 8 bit samples of V Y U Y
	DX_IMAGE_TYPE_RGBX32       // 32bpp like RGBA32 but with unused alpha
	DX_IMAGE_TYPE_RGBX8888     // 32bpp corresponding to RGBA with unused alpha
	DX_IMAGE_TYPE_BGRX8888     // 32bpp corresponding to BGRA with unused alpha
	DX_IMAGE_TYPE_YUV420SP     // Y as a plane then UV byte interleaved in plane with with same pitch half height
	DX_IMAGE_TYPE_YUV444PLANAR // Y U & V planes separately 4:4:4
	DX_IMAGE_TYPE_TF_U8        // T-format 8-bit U - same as TF_Y8 buf from U plane
	DX_IMAGE_TYPE_TF_V8        // T-format 8-bit U - same as TF_Y8 buf from V plane
	DX_IMAGE_TYPE_YUV420_16    // YUV4:2:0 planar 16bit values
	DX_IMAGE_TYPE_YUV_UV_16    // YUV4:2:0 codec format 16bit values
	DX_IMAGE_TYPE_YUV420_S     // YUV4:2:0 with UV in side-by-side format
	DX_IMAGE_TYPE_MIN          = DX_IMAGE_TYPE_RGB565
	DX_IMAGE_TYPE_MAX          = DX_IMAGE_TYPE_YUV420_S
)

const (
	/* Protection values */
	DX_PROTECTION_NONE DX_Protection = 0
	DX_PROTECTION_HDCP DX_Protection = 11
)

const (
	/* Clamp flags */
	DX_CLAMP_MODE_NONE DX_ClampMode = iota
	DX_CLAMP_MODE_LUMA_TRANSPARENT
	DX_CLAMP_MODE_CHROMA_TRANSPARENT
	DX_CLAMP_MODE_REPLACE
)

const (
	/* Alpha flags */
	DX_ALPHA_FLAG_FROM_SOURCE DX_AlphaFlags = iota // Bottom 2 bits sets the alpha mode
	DX_ALPHA_FLAG_FIXED_ALL_PIXELS
	DX_ALPHA_FLAG_FIXED_NON_ZERO
	DX_ALPHA_FLAG_FIXED_EXCEED_0X07
	DX_ALPHA_FLAG_PREMULT               DX_AlphaFlags = 1 << 16
	DX_ALPHA_FLAG_MIX                   DX_AlphaFlags = 1 << 17
	DX_ALPHA_FLAG__DISCARD_LOWER_LAYERS DX_AlphaFlags = 1 << 18
)

////////////////////////////////////////////////////////////////////////////////
// DISPLAY

func DX_Init() {
	C.bcm_host_init()
}

func DX_Stop() {
	C.vc_dispmanx_stop()
}

func DX_DisplayOpen(display DX_DisplayId) (DX_DisplayHandle, error) {

	if handle := DX_DisplayHandle(C.vc_dispmanx_display_open(C.uint32_t(display))); handle != DX_DISPLAY_NONE {
		return handle, nil
	} else {
		return DX_DISPLAY_NONE, gopi.ErrBadParameter
	}
}

func DX_DisplayClose(display DX_DisplayHandle) error {
	if C.vc_dispmanx_display_close(C.DISPMANX_DISPLAY_HANDLE_T(display)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrUnexpectedResponse
	}
}

func DX_DisplayGetInfo(display DX_DisplayHandle) (*DX_DisplayModeInfo, error) {
	info := &DX_DisplayModeInfo{}
	if C.vc_dispmanx_display_get_info(C.DISPMANX_DISPLAY_HANDLE_T(display), (*C.DISPMANX_MODEINFO_T)(unsafe.Pointer(info))) == DX_SUCCESS {
		return info, nil
	} else {
		return nil, gopi.ErrUnexpectedResponse
	}
}

func DX_DisplaySnapshot(display DX_DisplayHandle, resource DX_Resource, transform DX_Transform) error {
	if C.vc_dispmanx_snapshot(C.DISPMANX_DISPLAY_HANDLE_T(display), C.DISPMANX_RESOURCE_HANDLE_T(resource), C.DISPMANX_TRANSFORM_T(transform)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

////////////////////////////////////////////////////////////////////////////////
// RESOURCES

func DX_ResourceCreate(image_type DX_ImageType, size DX_Size) (DX_Resource, error) {
	var dummy C.uint32_t
	if handle := DX_Resource(C.vc_dispmanx_resource_create(C.VC_IMAGE_TYPE_T(image_type), C.uint32_t(size.W), C.uint32_t(size.H), (*C.uint32_t)(unsafe.Pointer(&dummy)))); handle == DX_NO_HANDLE {
		return DX_NO_HANDLE, gopi.ErrBadParameter
	} else {
		return handle, nil
	}
}

func DX_ResourceDelete(handle DX_Resource) error {
	if C.vc_dispmanx_resource_delete(C.DISPMANX_RESOURCE_HANDLE_T(handle)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

func DX_ResourceWriteData(handle DX_Resource, image_type DX_ImageType, src_pitch uint32, src uintptr, dest DX_Rect) error {
	if C.vc_dispmanx_resource_write_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), C.VC_IMAGE_TYPE_T(image_type), C.int(src_pitch), unsafe.Pointer(src), dest) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

func DX_ResourceReadData(handle DX_Resource, src DX_Rect, dest uintptr, dest_pitch uint32) error {
	if C.vc_dispmanx_resource_read_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), src, unsafe.Pointer(dest), C.uint32_t(dest_pitch)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

/*
func dxAlignUp(value, alignment uint32) uint32 {
	return ((value - 1) & ^(alignment - 1)) + alignment
}
*/

////////////////////////////////////////////////////////////////////////////////
// UPDATES

func DX_UpdateStart(priority int32) (DX_Update, error) {
	if handle := C.vc_dispmanx_update_start(C.int32_t(priority)); handle != DX_NO_HANDLE {
		return DX_Update(handle), nil
	} else {
		return DX_NO_HANDLE, gopi.ErrBadParameter
	}
}

func DX_UpdateSubmitSync(handle DX_Update) error {
	if C.vc_dispmanx_update_submit_sync(C.DISPMANX_UPDATE_HANDLE_T(handle)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

////////////////////////////////////////////////////////////////////////////////
// ELEMENTS

func DX_ElementAdd(update DX_Update, display DX_DisplayHandle, layer uint16, dest_rect DX_Rect, src_resource DX_Resource, src_size DX_Size, protection DX_Protection, alpha DX_Alpha, clamp DX_Clamp, transform DX_Transform) (DX_Element, error) {
	src_rect := DX_NewRect(0, 0, uint32(src_size.W)<<16, uint32(src_size.H)<<16)
	if handle := C.vc_dispmanx_element_add(
		C.DISPMANX_UPDATE_HANDLE_T(update),
		C.DISPMANX_DISPLAY_HANDLE_T(display),
		C.int32_t(layer),
		dest_rect,
		C.DISPMANX_RESOURCE_HANDLE_T(src_resource),
		src_rect,
		C.DISPMANX_PROTECTION_T(protection),
		(*C.VC_DISPMANX_ALPHA_T)(unsafe.Pointer(&alpha)),
		(*C.DISPMANX_CLAMP_T)(unsafe.Pointer(&clamp)),
		C.DISPMANX_TRANSFORM_T(transform)); handle != DX_NO_HANDLE {
		return DX_Element(handle), nil
	} else {
		return DX_NO_HANDLE, gopi.ErrBadParameter
	}
}

func DX_ElementRemove(update DX_Update, element DX_Element) error {
	if C.vc_dispmanx_element_remove(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_ELEMENT_HANDLE_T(element)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

func DX_ElementModified(update DX_Update, element DX_Element, rect DX_Rect) error {
	if C.vc_dispmanx_element_modified(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_ELEMENT_HANDLE_T(element), rect) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

/*
func dxElementChangeDestinationFrame(update dxUpdateHandle, element dxElementHandle, frame *DXFrame) bool {
	return C.vc_dispmanx_element_change_attributes(
		C.DISPMANX_UPDATE_HANDLE_T(update),
		C.DISPMANX_ELEMENT_HANDLE_T(element),
		C.uint32_t(DX_ELEMENT_CHANGE_DEST_RECT),
		C.int32_t(0),                          // layer
		C.uint8_t(0),                          // opacity
		(*C.VC_RECT_T)(unsafe.Pointer(frame)), // dest_rect
		(*C.VC_RECT_T)(unsafe.Pointer(nil)),   // src_rect
		C.DISPMANX_RESOURCE_HANDLE_T(0),       // mask
		C.DISPMANX_TRANSFORM_T(0),             // transform
	) == DX_ELEMENT_SUCCESS
}
*/

////////////////////////////////////////////////////////////////////////////////
// RECT

func DX_NewRect(x, y int32, w, h uint32) DX_Rect {
	return DX_Rect(&C.VC_RECT_T{C.int32_t(x), C.int32_t(y), C.int32_t(w), C.int32_t(h)})
}

func DX_RectSet(rect DX_Rect, x, y int32, w, h uint32) error {
	if C.vc_dispmanx_rect_set(rect, C.uint32_t(x), C.uint32_t(y), C.uint32_t(w), C.uint32_t(h)) != DX_SUCCESS {
		return gopi.ErrBadParameter
	} else {
		return nil
	}
}

func DX_RectSize(rect DX_Rect) DX_Size {
	return DX_Size{uint32(rect.width), uint32(rect.height)}
}

func DX_RectOrigin(rect DX_Rect) DX_Point {
	return DX_Point{int32(rect.x), int32(rect.y)}
}

func DX_RectString(r DX_Rect) string {
	return fmt.Sprintf("DX_Rect{ origin={%v,%v} size={%v,%v} }", r.x, r.y, r.width, r.height)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *DX_DisplayModeInfo) String() string {
	return fmt.Sprintf("DX_DisplayModeInfo{ size=%v transform=%v input_format=%v }", this.Size, this.Transform, this.InputFormat)
}

func (size DX_Size) String() string {
	return fmt.Sprintf("DX_Size{%v,%v}", size.W, size.H)
}

func (t DX_Transform) String() string {
	switch t {
	case DX_TRANSFORM_NONE:
		return "DX_TRANSFORM_NONE"
	case DX_TRANSFORM_ROTATE_90:
		return "DX_TRANSFORM_ROTATE_90"
	case DX_TRANSFORM_ROTATE_180:
		return "DX_TRANSFORM_ROTATE_180"
	case DX_TRANSFORM_ROTATE_270:
		return "DX_TRANSFORM_ROTATE_270"
	default:
		return "[?? Invalid DX_Transform value]"
	}
}

func (f DX_InputFormat) String() string {
	switch f {
	case DX_INPUT_FORMAT_RGB888:
		return "DX_INPUT_FORMAT_RGB888"
	case DX_INPUT_FORMAT_RGB565:
		return "DX_INPUT_FORMAT_RGB565"
	default:
		return "DX_INPUT_FORMAT_INVALID"
	}
}

func (d DX_DisplayId) String() string {
	switch d {
	case DX_DISPLAYID_MAIN_LCD:
		return "DX_DISPLAYID_MAIN_LCD"
	case DX_DISPLAYID_AUX_LCD:
		return "DX_DISPLAYID_AUX_LCD"
	case DX_DISPLAYID_HDMI:
		return "DX_DISPLAYID_HDMI"
	case DX_DISPLAYID_SDTV:
		return "DX_DISPLAYID_SDTV"
	case DX_DISPLAYID_FORCE_LCD:
		return "DX_DISPLAYID_FORCE_LCD"
	case DX_DISPLAYID_FORCE_TV:
		return "DX_DISPLAYID_FORCE_TV"
	case DX_DISPLAYID_FORCE_OTHER:
		return "DX_DISPLAYID_FORCE_OTHER"
	default:
		return "[?? Invalid DX_DisplayId value]"
	}
}

func (t DX_ImageType) String() string {
	switch t {
	case DX_IMAGE_TYPE_NONE:
		return "DX_IMAGE_TYPE_NONE"
	case DX_IMAGE_TYPE_RGB565:
		return "DX_IMAGE_TYPE_RGB565"
	case DX_IMAGE_TYPE_1BPP:
		return "DX_IMAGE_TYPE_1BPP"
	case DX_IMAGE_TYPE_YUV420:
		return "DX_IMAGE_TYPE_YUV420"
	case DX_IMAGE_TYPE_48BPP:
		return "DX_IMAGE_TYPE_48BPP"
	case DX_IMAGE_TYPE_RGB888:
		return "DX_IMAGE_TYPE_RGB888"
	case DX_IMAGE_TYPE_8BPP:
		return "DX_IMAGE_TYPE_8BPP"
	case DX_IMAGE_TYPE_4BPP:
		return "DX_IMAGE_TYPE_4BPP"
	case DX_IMAGE_TYPE_3D32:
		return "DX_IMAGE_TYPE_3D32"
	case DX_IMAGE_TYPE_3D32B:
		return "DX_IMAGE_TYPE_3D32B"
	case DX_IMAGE_TYPE_3D32MAT:
		return "DX_IMAGE_TYPE_3D32MAT"
	case DX_IMAGE_TYPE_RGB2X9:
		return "DX_IMAGE_TYPE_RGB2X9"
	case DX_IMAGE_TYPE_RGB666:
		return "DX_IMAGE_TYPE_RGB666"
	case DX_IMAGE_TYPE_PAL4_OBSOLETE:
		return "DX_IMAGE_TYPE_PAL4_OBSOLETE"
	case DX_IMAGE_TYPE_PAL8_OBSOLETE:
		return "DX_IMAGE_TYPE_PAL8_OBSOLETE"
	case DX_IMAGE_TYPE_RGBA32:
		return "DX_IMAGE_TYPE_RGBA32"
	case DX_IMAGE_TYPE_YUV422:
		return "DX_IMAGE_TYPE_YUV422"
	case DX_IMAGE_TYPE_RGBA565:
		return "DX_IMAGE_TYPE_RGBA565"
	case DX_IMAGE_TYPE_RGBA16:
		return "DX_IMAGE_TYPE_RGBA16"
	case DX_IMAGE_TYPE_YUV_UV:
		return "DX_IMAGE_TYPE_YUV_UV"
	case DX_IMAGE_TYPE_TF_RGBA32:
		return "DX_IMAGE_TYPE_TF_RGBA32"
	case DX_IMAGE_TYPE_TF_RGBX32:
		return "DX_IMAGE_TYPE_TF_RGBX32"
	case DX_IMAGE_TYPE_TF_FLOAT:
		return "DX_IMAGE_TYPE_TF_FLOAT"
	case DX_IMAGE_TYPE_TF_RGBA16:
		return "DX_IMAGE_TYPE_TF_RGBA16"
	case DX_IMAGE_TYPE_TF_RGBA5551:
		return "DX_IMAGE_TYPE_TF_RGBA5551"
	case DX_IMAGE_TYPE_TF_RGB565:
		return "DX_IMAGE_TYPE_TF_RGB565"
	case DX_IMAGE_TYPE_TF_YA88:
		return "DX_IMAGE_TYPE_TF_YA88"
	case DX_IMAGE_TYPE_TF_BYTE:
		return "DX_IMAGE_TYPE_TF_BYTE"
	case DX_IMAGE_TYPE_TF_PAL8:
		return "DX_IMAGE_TYPE_TF_PAL8"
	case DX_IMAGE_TYPE_TF_PAL4:
		return "DX_IMAGE_TYPE_TF_PAL4"
	case DX_IMAGE_TYPE_TF_ETC1:
		return "DX_IMAGE_TYPE_TF_ETC1"
	case DX_IMAGE_TYPE_BGR888:
		return "DX_IMAGE_TYPE_BGR888"
	case DX_IMAGE_TYPE_BGR888_NP:
		return "DX_IMAGE_TYPE_BGR888_NP"
	case DX_IMAGE_TYPE_BAYER:
		return "DX_IMAGE_TYPE_BAYER"
	case DX_IMAGE_TYPE_CODEC:
		return "DX_IMAGE_TYPE_CODEC"
	case DX_IMAGE_TYPE_YUV_UV32:
		return "DX_IMAGE_TYPE_YUV_UV32"
	case DX_IMAGE_TYPE_TF_Y8:
		return "DX_IMAGE_TYPE_TF_Y8"
	case DX_IMAGE_TYPE_TF_A8:
		return "DX_IMAGE_TYPE_TF_A8"
	case DX_IMAGE_TYPE_TF_SHORT:
		return "DX_IMAGE_TYPE_TF_SHORT"
	case DX_IMAGE_TYPE_TF_1BPP:
		return "DX_IMAGE_TYPE_TF_1BPP"
	case DX_IMAGE_TYPE_OPENGL:
		return "DX_IMAGE_TYPE_OPENGL"
	case DX_IMAGE_TYPE_YUV444I:
		return "DX_IMAGE_TYPE_YUV444I"
	case DX_IMAGE_TYPE_YUV422PLANAR:
		return "DX_IMAGE_TYPE_YUV422PLANAR"
	case DX_IMAGE_TYPE_ARGB8888:
		return "DX_IMAGE_TYPE_ARGB8888"
	case DX_IMAGE_TYPE_XRGB8888:
		return "DX_IMAGE_TYPE_XRGB8888"
	case DX_IMAGE_TYPE_YUV422YUYV:
		return "DX_IMAGE_TYPE_YUV422YUYV"
	case DX_IMAGE_TYPE_YUV422YVYU:
		return "DX_IMAGE_TYPE_YUV422YVYU"
	case DX_IMAGE_TYPE_YUV422UYVY:
		return "DX_IMAGE_TYPE_YUV422UYVY"
	case DX_IMAGE_TYPE_YUV422VYUY:
		return "DX_IMAGE_TYPE_YUV422VYUY"
	case DX_IMAGE_TYPE_RGBX32:
		return "DX_IMAGE_TYPE_RGBX32"
	case DX_IMAGE_TYPE_RGBX8888:
		return "DX_IMAGE_TYPE_RGBX8888"
	case DX_IMAGE_TYPE_BGRX8888:
		return "DX_IMAGE_TYPE_BGRX8888"
	case DX_IMAGE_TYPE_YUV420SP:
		return "DX_IMAGE_TYPE_YUV420SP"
	case DX_IMAGE_TYPE_YUV444PLANAR:
		return "DX_IMAGE_TYPE_YUV444PLANAR"
	case DX_IMAGE_TYPE_TF_U8:
		return "DX_IMAGE_TYPE_TF_U8"
	case DX_IMAGE_TYPE_TF_V8:
		return "DX_IMAGE_TYPE_TF_V8"
	case DX_IMAGE_TYPE_YUV420_16:
		return "DX_IMAGE_TYPE_YUV420_16"
	case DX_IMAGE_TYPE_YUV_UV_16:
		return "DX_IMAGE_TYPE_YUV_UV_16"
	case DX_IMAGE_TYPE_YUV420_S:
		return "DX_IMAGE_TYPE_YUV420_S"
	default:
		return "[?? Invalid DX_ImageType value]"
	}
}
