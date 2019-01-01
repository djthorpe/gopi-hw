package freetype_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	// Frameworks
	ft "github.com/djthorpe/gopi-hw/freetype"
)

////////////////////////////////////////////////////////////////////////////////
// TEST ENUMS

func TestStatus_000(t *testing.T) {
	for status := ft.FT_ERROR_MIN; status <= ft.FT_ERROR_MAX; status++ {
		status_error := fmt.Sprint(status.Error())
		if strings.HasPrefix(status_error, "FT_ERROR_") {
			t.Logf("%v => %s", int(status), status_error)
		} else {
			t.Logf("No status error for value: %v", status)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST LIBRARY

func TestLibrary_000(t *testing.T) {
	if library, err := ft.FT_Init(); err != nil {
		t.Error(err)
	} else if err := ft.FT_Destroy(library); err != nil {
		t.Error(err)
	}
}
func TestLibrary_001(t *testing.T) {
	if library, err := ft.FT_Init(); err != nil {
		t.Error(err)
	} else if major, minor, patch := ft.FT_Library_Version(library); major == 0 {
		t.Error("unexpected major version, ", major)
	} else {
		t.Logf("version={%v,%v,%v}", major, minor, patch)
		if err := ft.FT_Destroy(library); err != nil {
			t.Error(err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST FACE
func TestFace_001(t *testing.T) {
	if library, err := ft.FT_Init(); err != nil {
		t.Error(err)
	} else if face, err := ft.FT_NewFace(library, "", 0); err != ft.FT_ERROR_Cannot_Open_Resource {
		t.Error(err)
	} else if face != nil {
		t.Error("Expected face to be nil, is", face)
	} else if err := ft.FT_Destroy(library); err != nil {
		t.Error(err)
	}
}
func TestFace_002(t *testing.T) {
	if wd, err := os.Getwd(); err != nil {
		t.Error(err)
	} else if filenames, err := GetFontFilenames(filepath.Join(wd, "../etc/fonts")); err != nil {
		t.Error(err)
	} else if library, err := ft.FT_Init(); err != nil {
		t.Error(err)
	} else {
		for _, path := range filenames {
			if face, err := ft.FT_NewFace(library, path, 0); err != nil {
				t.Errorf("%v: %v", path, err)
			} else if name := ft.FT_FaceFamily(face); name == "" {
				t.Errorf("%v: Empty face family", path)
			} else if style := ft.FT_FaceStyle(face); style == "" {
				t.Errorf("%v: Empty face style", path)
			} else if num_faces := ft.FT_FaceNumFaces(face); num_faces == 0 {
				t.Errorf("%v: Unexpected num_faces", path)
			} else if num_glyphs := ft.FT_FaceNumGlyphs(face); num_glyphs == 0 {
				t.Errorf("%v: Unexpected num_glyphs", path)
			} else {
				style_flags := ft.FT_FaceStyleFlags(face)
				t.Log("File:", path)
				t.Log("     Family:", name)
				t.Log("      Style:", style)
				t.Log("  Num Faces:", num_faces)
				t.Log(" Num Glyphs:", num_glyphs)
				t.Log("      Flags:", style_flags)
				if err := ft.FT_DoneFace(face); err != nil {
					t.Errorf("%v: %v", path, err)
				}
			}

		}
		if err := ft.FT_Destroy(library); err != nil {
			t.Error(err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET FONT FILENAMES

func GetFontFilenames(path string) ([]string, error) {
	filenames := make([]string, 0)
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%v: %v", path, err)
	} else if stat.IsDir() != true {
		return nil, fmt.Errorf("%v: Expected dir", path)
	} else if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.ToLower(filepath.Ext(path)) != ".ttf" {
			return nil
		}
		// Accept TTF files
		filenames = append(filenames, path)
		return nil
	}); err != nil {
		return nil, err
	} else {
		return filenames, nil
	}
}
