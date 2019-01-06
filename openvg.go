/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Color struct {
	R, G, B, A float32
}

type (
	VGStrokeCapStyle  uint16
	VGStrokeJoinStyle uint16
	VGFillRule        uint16
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type VG interface {
	// Inherit general driver interface
	gopi.Driver

	/*
		// Atomic drawing operation
		Do(surface gopi.Surface, callback func() error) error

		// Create a path
		CreatePath() (VGPath, error)

		// Destroy a path
		DestroyPath(VGPath) error

		// Create a paintbrush
		CreatePaint(color Color) (VGPaint, error)

		// Destroy a paintbrush
		DestroyPaint(VGPaint) error

		// Clear surface to color
		Clear(surface gopi.Surface, color Color) error

		// Translate co-ordinate system
		Translate(offset gopi.Point) error

		// Scale co-ordinate system
		Scale(x, y float32) error

		// Shear co-ordinate system
		Shear(x, y float32) error

		// Rotate co-ordinate system
		Rotate(r float32) error

		// Load Identity Matrix (reset co-ordinate system)
		LoadIdentity() error
	*/
}

// Drawing Path
type VGPath interface {
	// Draw the path with both stroke and fill
	Draw(stroke, fill VGPaint) error

	// Stroke the path
	Stroke(stroke VGPaint) error

	// Fill the path
	Fill(fill VGPaint) error

	// Reset to empty path
	Clear() error

	// Close Path
	Close() error

	// Move To
	MoveTo(gopi.Point) error

	// Line To
	LineTo(...gopi.Point) error

	// Quad To
	QuadTo(p1, p2 gopi.Point) error

	// Cubic To
	CubicTo(p1, p2, p3 gopi.Point) error

	// Append a line to the path
	Line(start, end gopi.Point) error

	// Append a rectangle to the path
	Rect(origin gopi.Point, size gopi.Size) error

	// Append an ellipse to the path
	Ellipse(origin gopi.VGPoint, diameter gopi.Size) error

	// Append a circle to the path
	Circle(origin gopi.VGPoint, diameter float32) error
}

// Paintbrush for Fill and Stroke
type VGPaint interface {
	// Set color
	SetColor(color Color) error

	// Set fill rule
	SetFillRule(style VGFillRule) error

	// Set stroke width
	SetStrokeWidth(width float32) error

	// Set miter limit
	SetMiterLimit(value float32) error

	// Set stroke path endpoint styles (for joins and cap). Use
	// NONE as a value for no change
	SetStrokeStyle(VGStrokeJoinStyle, VGStrokeCapStyle) error

	// Set stroke dash pattern, call with no arguments to reset
	SetStrokeDash(...float32) error
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_STYLE_CAP_NONE   VGStrokeCapStyle = 0x0000
	VG_STYLE_CAP_BUTT   VGStrokeCapStyle = 0x1700 // Default
	VG_STYLE_CAP_ROUND  VGStrokeCapStyle = 0x1701
	VG_STYLE_CAP_SQUARE VGStrokeCapStyle = 0x1702
)

const (
	VG_STYLE_JOIN_NONE  VGStrokeJoinStyle = 0x0000
	VG_STYLE_JOIN_MITER VGStrokeJoinStyle = 0x1800 // Default
	VG_STYLE_JOIN_ROUND VGStrokeJoinStyle = 0x1801
	VG_STYLE_JOIN_BEVEL VGStrokeJoinStyle = 0x1802
)

const (
	VG_STYLE_FILL_NONE    VGFillRule = 0x0000
	VG_STYLE_FILL_NONZERO VGFillRule = 0x1900
	VG_STYLE_FILL_EVENODD VGFillRule = 0x1901 // Default
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

// Standard Colors
var (
	ColorRed       = Color{1.0, 0.0, 0.0, 1.0}
	ColorGreen     = Color{0.0, 1.0, 0.0, 1.0}
	ColorBlue      = Color{0.0, 0.0, 1.0, 1.0}
	ColorWhite     = Color{1.0, 1.0, 1.0, 1.0}
	ColorBlack     = Color{0.0, 0.0, 0.0, 1.0}
	ColorPurple    = Color{1.0, 0.0, 1.0, 1.0}
	ColorCyan      = Color{0.0, 1.0, 1.0, 1.0}
	ColorYellow    = Color{1.0, 1.0, 0.0, 1.0}
	ColorDarkGrey  = Color{0.25, 0.25, 0.25, 1.0}
	ColorLightGrey = Color{0.75, 0.75, 0.75, 1.0}
	ColorMidGrey   = Color{0.5, 0.5, 0.5, 1.0}
)
