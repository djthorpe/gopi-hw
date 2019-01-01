/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2019
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package freetype

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
  #cgo CFLAGS:   -I/usr/include/freetype2
  #cgo LDFLAGS:  -lfreetype
  #include <ft2build.h>
  #include FT_FREETYPE_H
*/
import "C"
import "unsafe"

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
