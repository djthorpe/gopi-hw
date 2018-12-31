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
	"fmt"
	"strings"

	// Frameworks
	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi-hw/rpi"
)

type displayregion struct {
	handle rpi.MMAL_DisplayRegion
}

func (this *displayregion) String() string {
	parts := ""
	parts += fmt.Sprintf("display=%v ", this.Display())
	parts += fmt.Sprintf("fullscreen=%v ", this.FullScreen())
	parts += fmt.Sprintf("layer=%v ", this.Layer())
	parts += fmt.Sprintf("alpha=0x%02X ", this.Alpha())
	parts += fmt.Sprintf("src_rect=%v ", rpi.MMALDisplayRegionGetSrcRect(this.handle))
	parts += fmt.Sprintf("dest_rect=%v ", rpi.MMALDisplayRegionGetDestRect(this.handle))
	parts += fmt.Sprintf("transform=%v ", this.Transform())
	parts += fmt.Sprintf("mode=%v ", this.Mode())
	parts += fmt.Sprintf("noaspect=%v ", this.NoAspect())
	parts += fmt.Sprintf("copyprotect=%v ", this.CopyProtect())
	return fmt.Sprintf("<sys.hw.mmal.displayregion>{ %v }", strings.TrimSpace(parts))
}

func (this *displayregion) Display() uint16 {
	return uint16(rpi.MMALDisplayRegionGetDisplayNum(this.handle))
}

func (this *displayregion) FullScreen() bool {
	return rpi.MMALDisplayRegionGetFullScreen(this.handle)
}

func (this *displayregion) SetFullScreen(value bool) {
	rpi.MMALDisplayRegionSetFullScreen(this.handle, value)
}

func (this *displayregion) Layer() int16 {
	return int16(rpi.MMALDisplayRegionGetLayer(this.handle))
}

func (this *displayregion) Alpha() uint8 {
	return uint8(rpi.MMALDisplayRegionGetAlpha(this.handle) & 0xFF)
}

func (this *displayregion) SetLayer(value int16) {
	rpi.MMALDisplayRegionSetLayer(this.handle, int32(value))
}

func (this *displayregion) SetAlpha(value uint8) {
	/** Bits 7-0: Level of opacity of the layer, where zero is fully transparent and 255 is fully opaque.
	 * Bits 31-8: Flags from MMAL_DISPLAYALPHAFLAGS_T for alpha mode selection */
	alpha := (rpi.MMALDisplayRegionGetAlpha(this.handle) & 0xFFFFFF00) | uint32(value)
	rpi.MMALDisplayRegionSetAlpha(this.handle, alpha)
}

func (this *displayregion) Transform() hw.MMALDisplayTransform {
	return hw.MMALDisplayTransform(rpi.MMALDisplayRegionGetTransform(this.handle))
}

func (this *displayregion) SetTransform(value hw.MMALDisplayTransform) {
	rpi.MMALDisplayRegionSetTransform(this.handle, value)
}

func (this *displayregion) Mode() hw.MMALDisplayMode {
	return hw.MMALDisplayMode(rpi.MMALDisplayRegionGetMode(this.handle))
}

func (this *displayregion) SetMode(value hw.MMALDisplayMode) {
	rpi.MMALDisplayRegionSetMode(this.handle, value)
}

func (this *displayregion) NoAspect() bool {
	return rpi.MMALDisplayRegionGetNoAspect(this.handle)
}

func (this *displayregion) CopyProtect() bool {
	return rpi.MMALDisplayRegionGetCopyProtect(this.handle)
}

func (this *displayregion) SetNoAspect(value bool) {
	rpi.MMALDisplayRegionSetNoAspect(this.handle, value)
}

func (this *displayregion) SetCopyProtect(value bool) {
	rpi.MMALDisplayRegionSetCopyProtect(this.handle, value)
}

func (this *displayregion) DestRect() hw.MMALRect {
	return rpi.MMALDisplayRegionGetDestRect(this.handle)
}

func (this *displayregion) SrcRect() hw.MMALRect {
	return rpi.MMALDisplayRegionGetSrcRect(this.handle)
}

func (this *displayregion) SetDestRect(value hw.MMALRect) {
	rpi.MMALDisplayRegionSetDestRect(this.handle, value)
}

func (this *displayregion) SetSrcRect(value hw.MMALRect) {
	rpi.MMALDisplayRegionSetSrcRect(this.handle, value)
}
