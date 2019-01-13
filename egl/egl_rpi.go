/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package egl

import (
	"strings"
	"unsafe"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: egl
#include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EGL_Display         C.EGLDisplay
	EGL_Config          C.EGLConfig
	EGL_ConfigAttrib    C.EGLint
	EGL_Context         C.EGLContext
	EGL_Surface         C.EGLSurface
	EGL_Error           C.EGLint
	EGL_Query           C.EGLint
	EGL_RenderableFlag  C.EGLint
	EGL_SurfaceTypeFlag C.EGLint
	EGL_API             C.EGLint
	EGL_NativeWindow    uintptr
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EGL_FALSE = C.EGLBoolean(0)
	EGL_TRUE  = C.EGLBoolean(1)
)

const (
	// EGL_Error
	EGL_SUCCESS             EGL_Error = 0x3000
	EGL_NOT_INITIALIZED     EGL_Error = 0x3001
	EGL_BAD_ACCESS          EGL_Error = 0x3002
	EGL_BAD_ALLOC           EGL_Error = 0x3003
	EGL_BAD_ATTRIBUTE       EGL_Error = 0x3004
	EGL_BAD_CONFIG          EGL_Error = 0x3005
	EGL_BAD_CONTEXT         EGL_Error = 0x3006
	EGL_BAD_CURRENT_SURFACE EGL_Error = 0x3007
	EGL_BAD_DISPLAY         EGL_Error = 0x3008
	EGL_BAD_MATCH           EGL_Error = 0x3009
	EGL_BAD_NATIVE_PIXMAP   EGL_Error = 0x300A
	EGL_BAD_NATIVE_WINDOW   EGL_Error = 0x300B
	EGL_BAD_PARAMETER       EGL_Error = 0x300C
	EGL_BAD_SURFACE         EGL_Error = 0x300D
	EGL_CONTEXT_LOST        EGL_Error = 0x300E // EGL 1.1 - IMG_power_management
)

const (
	// EGL_Query
	EGL_QUERY_VENDOR      EGL_Query = 0x3053
	EGL_QUERY_VERSION     EGL_Query = 0x3054
	EGL_QUERY_EXTENSIONS  EGL_Query = 0x3055
	EGL_QUERY_CLIENT_APIS EGL_Query = 0x308D
)

const (
	// EGL_ConfigAttrib
	EGL_BUFFER_SIZE             EGL_ConfigAttrib = 0x3020
	EGL_ALPHA_SIZE              EGL_ConfigAttrib = 0x3021
	EGL_BLUE_SIZE               EGL_ConfigAttrib = 0x3022
	EGL_GREEN_SIZE              EGL_ConfigAttrib = 0x3023
	EGL_RED_SIZE                EGL_ConfigAttrib = 0x3024
	EGL_DEPTH_SIZE              EGL_ConfigAttrib = 0x3025
	EGL_STENCIL_SIZE            EGL_ConfigAttrib = 0x3026
	EGL_CONFIG_CAVEAT           EGL_ConfigAttrib = 0x3027
	EGL_CONFIG_ID               EGL_ConfigAttrib = 0x3028
	EGL_LEVEL                   EGL_ConfigAttrib = 0x3029
	EGL_MAX_PBUFFER_HEIGHT      EGL_ConfigAttrib = 0x302A
	EGL_MAX_PBUFFER_PIXELS      EGL_ConfigAttrib = 0x302B
	EGL_MAX_PBUFFER_WIDTH       EGL_ConfigAttrib = 0x302C
	EGL_NATIVE_RENDERABLE       EGL_ConfigAttrib = 0x302D
	EGL_NATIVE_VISUAL_ID        EGL_ConfigAttrib = 0x302E
	EGL_NATIVE_VISUAL_TYPE      EGL_ConfigAttrib = 0x302F
	EGL_SAMPLES                 EGL_ConfigAttrib = 0x3031
	EGL_SAMPLE_BUFFERS          EGL_ConfigAttrib = 0x3032
	EGL_SURFACE_TYPE            EGL_ConfigAttrib = 0x3033
	EGL_TRANSPARENT_TYPE        EGL_ConfigAttrib = 0x3034
	EGL_TRANSPARENT_BLUE_VALUE  EGL_ConfigAttrib = 0x3035
	EGL_TRANSPARENT_GREEN_VALUE EGL_ConfigAttrib = 0x3036
	EGL_TRANSPARENT_RED_VALUE   EGL_ConfigAttrib = 0x3037
	EGL_NONE                    EGL_ConfigAttrib = 0x3038 // Attrib list terminator
	EGL_BIND_TO_TEXTURE_RGB     EGL_ConfigAttrib = 0x3039
	EGL_BIND_TO_TEXTURE_RGBA    EGL_ConfigAttrib = 0x303A
	EGL_MIN_SWAP_INTERVAL       EGL_ConfigAttrib = 0x303B
	EGL_MAX_SWAP_INTERVAL       EGL_ConfigAttrib = 0x303C
	EGL_LUMINANCE_SIZE          EGL_ConfigAttrib = 0x303D
	EGL_ALPHA_MASK_SIZE         EGL_ConfigAttrib = 0x303E
	EGL_COLOR_BUFFER_TYPE       EGL_ConfigAttrib = 0x303F
	EGL_RENDERABLE_TYPE         EGL_ConfigAttrib = 0x3040
	EGL_MATCH_NATIVE_PIXMAP     EGL_ConfigAttrib = 0x3041 // Pseudo-attribute (not queryable)
	EGL_CONFORMANT              EGL_ConfigAttrib = 0x3042
	EGL_COMFIG_ATTRIB_MIN                        = EGL_BUFFER_SIZE
	EGL_COMFIG_ATTRIB_MAX                        = EGL_CONFORMANT
)

const (
	EGL_RENDERABLE_FLAG_OPENGL_ES  EGL_RenderableFlag = 0x0001
	EGL_RENDERABLE_FLAG_OPENVG     EGL_RenderableFlag = 0x0002
	EGL_RENDERABLE_FLAG_OPENGL_ES2 EGL_RenderableFlag = 0x0004
	EGL_RENDERABLE_FLAG_OPENGL     EGL_RenderableFlag = 0x0008
	EGL_RENDERABLE_FLAG_MIN                           = EGL_RENDERABLE_FLAG_OPENGL_ES
	EGL_RENDERABLE_FLAG_MAX                           = EGL_RENDERABLE_FLAG_OPENGL
)

const (
	EGL_SURFACETYPE_FLAG_PBUFFER                 EGL_SurfaceTypeFlag = 0x0001 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_PIXMAP                  EGL_SurfaceTypeFlag = 0x0002 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_WINDOW                  EGL_SurfaceTypeFlag = 0x0004 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_VG_COLORSPACE_LINEAR    EGL_SurfaceTypeFlag = 0x0020 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_VG_ALPHA_FORMAT_PRE     EGL_SurfaceTypeFlag = 0x0040 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_MULTISAMPLE_RESOLVE_BOX EGL_SurfaceTypeFlag = 0x0200 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_SWAP_BEHAVIOR_PRESERVED EGL_SurfaceTypeFlag = 0x0400 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_MIN                                         = EGL_SURFACETYPE_FLAG_PBUFFER
	EGL_SURFACETYPE_FLAG_MAX                                         = EGL_SURFACETYPE_FLAG_SWAP_BEHAVIOR_PRESERVED
)

const (
	EGL_API_NONE      EGL_API = 0
	EGL_API_OPENGL_ES EGL_API = 0x30A0
	EGL_API_OPENVG    EGL_API = 0x30A1
	EGL_API_OPENGL    EGL_API = 0x30A2
	EGL_API_MIN               = EGL_API_OPENGL_ES
	EGL_API_MAX               = EGL_API_OPENGL
)

var (
	EGL_SurfaceTypeMap = map[string]gopi.SurfaceFlags{
		"OpenGL":     gopi.SURFACE_FLAG_OPENGL,
		"OpenGL_ES":  gopi.SURFACE_FLAG_OPENGL_ES,
		"OpenGL_ES2": gopi.SURFACE_FLAG_OPENGL_ES2,
		"OpenVG":     gopi.SURFACE_FLAG_OPENVG,
	}
	EGL_APIMap = map[gopi.SurfaceFlags]EGL_API{
		gopi.SURFACE_FLAG_OPENGL_ES: EGL_API_OPENGL_ES,
		gopi.SURFACE_FLAG_OPENVG:    EGL_API_OPENVG,
		gopi.SURFACE_FLAG_OPENGL:    EGL_API_OPENGL,
	}
	EGL_RenderableMap = map[gopi.SurfaceFlags]EGL_RenderableFlag{
		gopi.SURFACE_FLAG_OPENGL:    EGL_RENDERABLE_FLAG_OPENGL,
		gopi.SURFACE_FLAG_OPENGL_ES: EGL_RENDERABLE_FLAG_OPENGL_ES,
		gopi.SURFACE_FLAG_OPENVG:    EGL_RENDERABLE_FLAG_OPENVG,
	}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func EGL_GetError() error {
	if err := EGL_Error(C.eglGetError()); err == EGL_SUCCESS {
		return nil
	} else {
		return err
	}
}

func EGL_Initialize(display EGL_Display) (int, int, error) {
	var major, minor C.EGLint
	if C.eglInitialize(C.EGLDisplay(display), (*C.EGLint)(unsafe.Pointer(&major)), (*C.EGLint)(unsafe.Pointer(&minor))) != EGL_TRUE {
		return 0, 0, EGL_GetError()
	} else {
		return int(major), int(minor), nil
	}
}

func EGL_Terminate(display EGL_Display) error {
	if C.eglTerminate(C.EGLDisplay(display)) != EGL_TRUE {
		return EGL_GetError()
	} else {
		return nil
	}
}

func EGL_GetDisplay(display uint) EGL_Display {
	return EGL_Display(C.eglGetDisplay(C.EGLNativeDisplayType(uintptr(display))))
}

func EGL_QueryString(display EGL_Display, value EGL_Query) string {
	return C.GoString(C.eglQueryString(C.EGLDisplay(display), C.EGLint(value)))
}

////////////////////////////////////////////////////////////////////////////////
// SURFACE CONFIGS

func EGL_GetConfigs(display EGL_Display) ([]EGL_Config, error) {
	var num_config C.EGLint
	if C.eglGetConfigs(C.EGLDisplay(display), (*C.EGLConfig)(nil), C.EGLint(0), &num_config) != EGL_TRUE {
		return nil, EGL_GetError()
	}
	if num_config == C.EGLint(0) {
		return nil, EGL_BAD_CONFIG
	}
	// configs is a slice so we need to pass the slice pointer
	configs := make([]EGL_Config, num_config)
	if C.eglGetConfigs(C.EGLDisplay(display), (*C.EGLConfig)(unsafe.Pointer(&configs[0])), num_config, &num_config) != EGL_TRUE {
		return nil, EGL_GetError()
	} else {
		return configs, nil
	}
}

func EGL_GetConfigAttrib(display EGL_Display, config EGL_Config, attrib EGL_ConfigAttrib) (int, error) {
	var value C.EGLint
	if C.eglGetConfigAttrib(C.EGLDisplay(display), C.EGLConfig(config), C.EGLint(attrib), &value) != EGL_TRUE {
		return 0, EGL_GetError()
	} else {
		return int(value), nil
	}
}

func EGL_GetConfigAttribs(display EGL_Display, config EGL_Config) (map[EGL_ConfigAttrib]int, error) {
	attribs := make(map[EGL_ConfigAttrib]int, 0)
	for k := EGL_COMFIG_ATTRIB_MIN; k <= EGL_COMFIG_ATTRIB_MAX; k++ {
		if v, err := EGL_GetConfigAttrib(display, config, k); err == EGL_BAD_ATTRIBUTE {
			continue
		} else if err != nil {
			return nil, err
		} else {
			attribs[k] = v
		}
	}
	return attribs, nil
}

func EGL_ChooseConfig_(display EGL_Display, attributes map[EGL_ConfigAttrib]int) ([]EGL_Config, error) {
	var num_config C.EGLint

	// Make list of attributes as eglInt values
	attribute_list := make([]C.EGLint, len(attributes)*2+1)
	i := 0
	for k, v := range attributes {
		attribute_list[i] = C.EGLint(k)
		attribute_list[i+1] = C.EGLint(v)
		i = i + 2
	}
	attribute_list[i] = C.EGLint(EGL_NONE)

	// Get number of configurations this matches
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(nil), C.EGLint(0), &num_config) != EGL_TRUE {
		return nil, EGL_GetError()
	}
	// Return EGL_BAD_ATTRIBUTE if the attribute set doesn't match
	if num_config == 0 {
		return nil, EGL_BAD_ATTRIBUTE
	}
	// Allocate an array
	configs := make([]EGL_Config, num_config)
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(nil), C.EGLint(0), &num_config) != EGL_TRUE {
		return nil, EGL_GetError()
	}
	// Return the configurations
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(unsafe.Pointer(&configs[0])), num_config, &num_config) != EGL_TRUE {
		return nil, EGL_GetError()
	} else {
		return configs, nil
	}
}

func EGL_ChooseConfig(display EGL_Display, r_bits, g_bits, b_bits, a_bits uint, surface_type EGL_SurfaceTypeFlag, renderable_type EGL_RenderableFlag) (EGL_Config, error) {
	if configs, err := EGL_ChooseConfig_(display, map[EGL_ConfigAttrib]int{
		EGL_RED_SIZE:        int(r_bits),
		EGL_GREEN_SIZE:      int(g_bits),
		EGL_BLUE_SIZE:       int(b_bits),
		EGL_ALPHA_SIZE:      int(a_bits),
		EGL_SURFACE_TYPE:    int(surface_type),
		EGL_RENDERABLE_TYPE: int(renderable_type),
	}); err != nil {
		return nil, err
	} else if len(configs) == 0 {
		return nil, EGL_BAD_CONFIG
	} else {
		return configs[0], nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// API

func EGL_BindAPI(api EGL_API) error {
	if success := C.eglBindAPI(C.EGLenum(api)); success != EGL_TRUE {
		return EGL_GetError()
	} else {
		return nil
	}
}

func EGL_QueryAPI() (EGL_API, error) {
	if api := EGL_API(C.eglQueryAPI()); api == 0 {
		return EGL_API_NONE, EGL_GetError()
	} else {
		return api, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// CONTEXT

func EGL_CreateContext(display EGL_Display, config EGL_Config, share_context EGL_Context) (EGL_Context, error) {
	if context := EGL_Context(C.eglCreateContext(C.EGLDisplay(display), C.EGLConfig(config), C.EGLContext(share_context), nil)); context == nil {
		return nil, EGL_GetError()
	} else {
		return context, nil
	}
}

func EGL_DestroyContext(display EGL_Display, context EGL_Context) error {
	if C.eglDestroyContext(C.EGLDisplay(display), C.EGLContext(context)) != EGL_TRUE {
		return EGL_GetError()
	} else {
		return nil
	}
}

func EGL_MakeCurrent(display EGL_Display, draw, read EGL_Surface, context EGL_Context) error {
	if C.eglMakeCurrent(C.EGLDisplay(display), C.EGLSurface(draw), C.EGLSurface(read), C.EGLContext(context)) != EGL_TRUE {
		return EGL_GetError()
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// SURFACE

func EGL_CreateSurface(display EGL_Display, config EGL_Config, window EGL_NativeWindow) (EGL_Surface, error) {
	if surface := EGL_Surface(C.eglCreateWindowSurface(C.EGLDisplay(display), C.EGLConfig(config), C.EGLNativeWindowType(window), nil)); surface == nil {
		return nil, EGL_GetError()
	} else {
		return surface, nil
	}
}

func EGL_DestroySurface(display EGL_Display, surface EGL_Surface) error {
	if C.eglDestroySurface(C.EGLDisplay(display), C.EGLSurface(surface)) != EGL_TRUE {
		return EGL_GetError()
	} else {
		return nil
	}
}

func EGL_SwapBuffers(display EGL_Display, surface EGL_Surface) error {
	if C.eglSwapBuffers(C.EGLDisplay(display), C.EGLSurface(surface)) != EGL_TRUE {
		return EGL_GetError()
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e EGL_Error) Error() string {
	switch e {
	case EGL_SUCCESS:
		return "EGL_SUCCESS"
	case EGL_NOT_INITIALIZED:
		return "EGL_NOT_INITIALIZED"
	case EGL_BAD_ACCESS:
		return "EGL_BAD_ACCESS"
	case EGL_BAD_ALLOC:
		return "EGL_BAD_ALLOC"
	case EGL_BAD_ATTRIBUTE:
		return "EGL_BAD_ATTRIBUTE"
	case EGL_BAD_CONFIG:
		return "EGL_BAD_CONFIG"
	case EGL_BAD_CONTEXT:
		return "EGL_BAD_CONTEXT"
	case EGL_BAD_CURRENT_SURFACE:
		return "EGL_BAD_CURRENT_SURFACE"
	case EGL_BAD_DISPLAY:
		return "EGL_BAD_DISPLAY"
	case EGL_BAD_MATCH:
		return "EGL_BAD_MATCH"
	case EGL_BAD_NATIVE_PIXMAP:
		return "EGL_BAD_NATIVE_PIXMAP"
	case EGL_BAD_NATIVE_WINDOW:
		return "EGL_BAD_NATIVE_WINDOW"
	case EGL_BAD_PARAMETER:
		return "EGL_BAD_PARAMETER"
	case EGL_BAD_SURFACE:
		return "EGL_BAD_SURFACE"
	case EGL_CONTEXT_LOST:
		return "EGL_CONTEXT_LOST"
	default:
		return "[?? Unknown EGL_Error value]"
	}
}

func (a EGL_ConfigAttrib) Error() string {
	switch a {
	case EGL_BUFFER_SIZE:
		return "EGL_BUFFER_SIZE"
	case EGL_ALPHA_SIZE:
		return "EGL_ALPHA_SIZE"
	case EGL_BLUE_SIZE:
		return "EGL_BLUE_SIZE"
	case EGL_GREEN_SIZE:
		return "EGL_GREEN_SIZE"
	case EGL_RED_SIZE:
		return "EGL_RED_SIZE"
	case EGL_DEPTH_SIZE:
		return "EGL_DEPTH_SIZE"
	case EGL_STENCIL_SIZE:
		return "EGL_STENCIL_SIZE"
	case EGL_CONFIG_CAVEAT:
		return "EGL_CONFIG_CAVEAT"
	case EGL_CONFIG_ID:
		return "EGL_CONFIG_ID"
	case EGL_LEVEL:
		return "EGL_LEVEL"
	case EGL_MAX_PBUFFER_HEIGHT:
		return "EGL_MAX_PBUFFER_HEIGHT"
	case EGL_MAX_PBUFFER_PIXELS:
		return "EGL_MAX_PBUFFER_PIXELS"
	case EGL_MAX_PBUFFER_WIDTH:
		return "EGL_MAX_PBUFFER_WIDTH"
	case EGL_NATIVE_RENDERABLE:
		return "EGL_NATIVE_RENDERABLE"
	case EGL_NATIVE_VISUAL_ID:
		return "EGL_NATIVE_VISUAL_ID"
	case EGL_NATIVE_VISUAL_TYPE:
		return "EGL_NATIVE_VISUAL_TYPE"
	case EGL_SAMPLES:
		return "EGL_SAMPLES"
	case EGL_SAMPLE_BUFFERS:
		return "EGL_SAMPLE_BUFFERS"
	case EGL_SURFACE_TYPE:
		return "EGL_SURFACE_TYPE"
	case EGL_TRANSPARENT_TYPE:
		return "EGL_TRANSPARENT_TYPE"
	case EGL_TRANSPARENT_BLUE_VALUE:
		return "EGL_TRANSPARENT_BLUE_VALUE"
	case EGL_TRANSPARENT_GREEN_VALUE:
		return "EGL_TRANSPARENT_GREEN_VALUE"
	case EGL_TRANSPARENT_RED_VALUE:
		return "EGL_TRANSPARENT_RED_VALUE"
	case EGL_NONE:
		return "EGL_NONE"
	case EGL_BIND_TO_TEXTURE_RGB:
		return "EGL_BIND_TO_TEXTURE_RGB"
	case EGL_BIND_TO_TEXTURE_RGBA:
		return "EGL_BIND_TO_TEXTURE_RGBA"
	case EGL_MIN_SWAP_INTERVAL:
		return "EGL_MIN_SWAP_INTERVAL"
	case EGL_MAX_SWAP_INTERVAL:
		return "EGL_MAX_SWAP_INTERVAL"
	case EGL_LUMINANCE_SIZE:
		return "EGL_LUMINANCE_SIZE"
	case EGL_ALPHA_MASK_SIZE:
		return "EGL_ALPHA_MASK_SIZE"
	case EGL_COLOR_BUFFER_TYPE:
		return "EGL_COLOR_BUFFER_TYPE"
	case EGL_RENDERABLE_TYPE:
		return "EGL_RENDERABLE_TYPE"
	case EGL_MATCH_NATIVE_PIXMAP:
		return "EGL_MATCH_NATIVE_PIXMAP"
	case EGL_CONFORMANT:
		return "EGL_CONFORMANT"
	default:
		return "[?? Invalid EGL_ConfigAttrib value]"
	}
}

func (a EGL_API) String() string {
	switch a {
	case EGL_API_OPENGL_ES:
		return "EGL_API_OPENGL_ES"
	case EGL_API_OPENGL:
		return "EGL_API_OPENGL"
	case EGL_API_OPENVG:
		return "EGL_API_OPENVG"
	default:
		return "[?? Invalid EGL_API value]"
	}
}

func (f EGL_RenderableFlag) String() string {
	parts := ""
	for flag := EGL_RENDERABLE_FLAG_MIN; flag <= EGL_RENDERABLE_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case EGL_RENDERABLE_FLAG_OPENGL_ES:
			parts += "|" + "EGL_RENDERABLE_FLAG_OPENGL_ES"
		case EGL_RENDERABLE_FLAG_OPENVG:
			parts += "|" + "EGL_RENDERABLE_FLAG_OPENVG"
		case EGL_RENDERABLE_FLAG_OPENGL_ES2:
			parts += "|" + "EGL_RENDERABLE_FLAG_OPENGL_ES2"
		case EGL_RENDERABLE_FLAG_OPENGL:
			parts += "|" + "EGL_RENDERABLE_FLAG_OPENGL"
		default:
			parts += "|" + "[?? Invalid EGL_RenderableFlag value]"
		}
	}
	return strings.Trim(parts, "|")
}

func (f EGL_SurfaceTypeFlag) String() string {
	parts := ""
	for flag := EGL_SURFACETYPE_FLAG_MIN; flag <= EGL_SURFACETYPE_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case EGL_SURFACETYPE_FLAG_PBUFFER:
			parts += "|" + "EGL_SURFACETYPE_FLAG_PBUFFER"
		case EGL_SURFACETYPE_FLAG_PIXMAP:
			parts += "|" + "EGL_SURFACETYPE_FLAG_PIXMAP"
		case EGL_SURFACETYPE_FLAG_WINDOW:
			parts += "|" + "EGL_SURFACETYPE_FLAG_WINDOW"
		case EGL_SURFACETYPE_FLAG_VG_COLORSPACE_LINEAR:
			parts += "|" + "EGL_SURFACETYPE_FLAG_VG_COLORSPACE_LINEAR"
		case EGL_SURFACETYPE_FLAG_VG_ALPHA_FORMAT_PRE:
			parts += "|" + "EGL_SURFACETYPE_FLAG_VG_ALPHA_FORMAT_PRE"
		case EGL_SURFACETYPE_FLAG_MULTISAMPLE_RESOLVE_BOX:
			parts += "|" + "EGL_SURFACETYPE_FLAG_MULTISAMPLE_RESOLVE_BOX"
		case EGL_SURFACETYPE_FLAG_SWAP_BEHAVIOR_PRESERVED:
			parts += "|" + "EGL_SURFACETYPE_FLAG_SWAP_BEHAVIOR_PRESERVED"
		default:
			parts += "|" + "[?? Invalid EGL_SurfaceTypeFlag value]"
		}
	}
	return strings.Trim(parts, "|")
}

/*
////////////////////////////////////////////////////////////////////////////////
// TYPES

// Native window structure
type EGLNativeWindowType struct {
	// TODO	element DXElement
	width  int
	height int
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func eglCreateWindowSurface(display EGLDisplay, config eglConfig, native EGLNativeWindowType) (EGLSurface, EGLError) {
	return nil, EGL_BAD_SURFACE
}

func eglCreatePbufferSurface(display EGLDisplay, config eglConfig, native EGLNativeWindowType) (EGLSurface, EGLError) {
	return nil, EGL_BAD_SURFACE
}

func eglCreatePixmapSurface(display EGLDisplay, config eglConfig, native EGLNativeWindowType) (EGLSurface, EGLError) {
	return nil, EGL_BAD_SURFACE
}

func EGLDestroySurface(display EGLDisplay, surface EGLSurface) EGLError {
	if C.eglDestroySurface(C.EGLDisplay(display), C.EGLSurface(surface)) != C.EGLBoolean(EGL_TRUE) {
		return eglGetError()
	} else {
		return EGL_SUCCESS
	}
}

EGLAPI EGLSurface EGLAPIENTRY eglCreateWindowSurface(EGLDisplay dpy, EGLConfig config,
	EGLNativeWindowType win,
	const EGLint *attrib_list);
EGLAPI EGLSurface EGLAPIENTRY eglCreatePbufferSurface(EGLDisplay dpy, EGLConfig config,
	 const EGLint *attrib_list);
EGLAPI EGLSurface EGLAPIENTRY eglCreatePixmapSurface(EGLDisplay dpy, EGLConfig config,
	EGLNativePixmapType pixmap,
	const EGLint *attrib_list);
EGLAPI EGLBoolean EGLAPIENTRY eglDestroySurface(EGLDisplay dpy, EGLSurface surface);
EGLAPI EGLBoolean EGLAPIENTRY eglQuerySurface(EGLDisplay dpy, EGLSurface surface,
EGLint attribute, EGLint *value);

*/
