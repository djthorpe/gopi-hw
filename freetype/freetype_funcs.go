/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2019
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package freetype

import (
	"fmt"
	"unsafe"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: freetype2
#include <ft2build.h>
#include FT_FREETYPE_H
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// LIBRARY FUNCTIONS

func FT_Init() (FT_Library, error) {
	var handle C.FT_Library
	if err := FT_Error(C.FT_Init_FreeType((*C.FT_Library)(&handle))); err != FT_SUCCESS {
		return nil, err
	} else {
		return FT_Library(handle), nil
	}
}

func FT_Destroy(handle FT_Library) error {
	if err := FT_Error(C.FT_Done_FreeType(handle)); err != FT_SUCCESS {
		return err
	} else {
		return nil
	}
}

func FT_Library_Version(handle FT_Library) (int, int, int) {
	var major, minor, patch C.FT_Int
	C.FT_Library_Version(handle, (*C.FT_Int)(&major), (*C.FT_Int)(&minor), (*C.FT_Int)(&patch))
	return int(major), int(minor), int(patch)
}

////////////////////////////////////////////////////////////////////////////////
// FACE FUNCTIONS

func FT_NewFace(handle FT_Library, path string, index uint) (FT_Face, error) {
	var face C.FT_Face
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))
	if err := FT_Error(C.FT_New_Face(handle, cstr, C.FT_Long(index), &face)); err != FT_SUCCESS {
		return nil, err
	} else {
		return FT_Face(face), nil
	}
}

func FT_DoneFace(handle FT_Face) error {
	if err := FT_Error(C.FT_Done_Face(handle)); err != FT_SUCCESS {
		return err
	} else {
		return nil
	}
}

func FT_SelectCharmap(handle FT_Face, encoding FT_Encoding) error {
	if err := FT_Error(C.FT_Select_Charmap(handle, C.FT_Encoding(encoding))); err != FT_SUCCESS {
		return err
	} else {
		return nil
	}
}

func FT_FaceFamily(handle FT_Face) string {
	return C.GoString(handle.family_name)
}

func FT_FaceStyle(handle FT_Face) string {
	fmt.Println(handle.style_name)
	return C.GoString(handle.style_name)
}

func FT_FaceIndex(handle FT_Face) uint {
	return uint(handle.face_index)
}

func FT_FaceNumFaces(handle FT_Face) uint {
	return uint(handle.num_faces)
}

func FT_FaceNumGlyphs(handle FT_Face) uint {
	return uint(handle.num_glyphs)
}

func FT_FaceStyleFlags(handle FT_Face) gopi.FontFlags {
	return gopi.FontFlags(handle.style_flags)
}
