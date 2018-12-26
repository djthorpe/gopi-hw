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

type annotation struct {
	handle rpi.MMAL_CameraAnnotation
}

func (this *annotation) String() string {
	parts := ""
	parts += fmt.Sprintf("enabled=%v ", this.Enabled())

	if this.Enabled() {
		parts += fmt.Sprintf("show_shutter=%v ", this.ShowShutter())
		parts += fmt.Sprintf("show_analog_gain=%v ", this.ShowAnalogGain())
		parts += fmt.Sprintf("show_lens=%v ", this.ShowLens())
		parts += fmt.Sprintf("show_caf=%v ", this.ShowCAF())
		parts += fmt.Sprintf("show_motion=%v ", this.ShowMotion())
		parts += fmt.Sprintf("show_frame_num=%v ", this.ShowFrameNum())

		parts += fmt.Sprintf("text=\"%v\" ", this.Text())
		parts += fmt.Sprintf("text_background=%v ", this.TextBackground())
		parts += fmt.Sprintf("text_size=%v ", this.TextSize())
		parts += fmt.Sprintf("text_justify=%v ", this.TextJustify())
	}

	return fmt.Sprintf("<sys.hw.mmal.annotation>{ %v }", strings.TrimSpace(parts))
}

func (this *annotation) Enabled() bool {
	return rpi.MMALCameraAnnotationEnabled(this.handle)
}

func (this *annotation) ShowShutter() bool {
	return rpi.MMALCameraAnnotationShowShutter(this.handle)
}

func (this *annotation) SetShowShutter(value bool) {
	rpi.MMALCameraAnnotationSetShowShutter(this.handle, value)
	rpi.MMALCameraAnnotationSetEnabled(this.handle, true)
}

func (this *annotation) ShowAnalogGain() bool {
	return rpi.MMALCameraAnnotationShowAnalogGain(this.handle)
}
func (this *annotation) ShowLens() bool {
	return rpi.MMALCameraAnnotationShowLens(this.handle)
}

func (this *annotation) ShowCAF() bool {
	return rpi.MMALCameraAnnotationShowCAF(this.handle)
}

func (this *annotation) ShowMotion() bool {
	return rpi.MMALCameraAnnotationShowMotion(this.handle)
}

func (this *annotation) ShowFrameNum() bool {
	return rpi.MMALCameraAnnotationShowFrameNum(this.handle)
}

func (this *annotation) TextBackground() bool {
	return false
}
func (this *annotation) BackgroundColor() (uint8, uint8, uint8) {
	return 0, 0, 0
}
func (this *annotation) TextColor() (uint8, uint8, uint8) {
	return 0, 0, 0
}

func (this *annotation) TextSize() uint8 {
	return rpi.MMALCameraAnnotationTextSize(this.handle)
}
func (this *annotation) SetTextSize(value uint8) {
	rpi.MMALCameraAnnotationSetTextSize(this.handle, value)
	rpi.MMALCameraAnnotationSetEnabled(this.handle, true)
}

func (this *annotation) Text() string {
	return rpi.MMALCameraAnnotationText(this.handle)
}

func (this *annotation) SetText(value string) {
	rpi.MMALCameraAnnotationSetText(this.handle, value)
	rpi.MMALCameraAnnotationSetEnabled(this.handle, true)
}

func (this *annotation) TextJustify() hw.MMALTextJustify {
	return 0
}
func (this *annotation) TextOffset() (uint32, uint32) {
	return 0, 0
}

func (this *annotation) SetShowAnalogGain(bool)            {}
func (this *annotation) SetShowLens(bool)                  {}
func (this *annotation) SetShowCAF(bool)                   {}
func (this *annotation) SetShowMotion(bool)                {}
func (this *annotation) SetShowFrameNum(bool)              {}
func (this *annotation) SetTextBackground(bool)            {}
func (this *annotation) SetBackgroundColor(y, u, v uint8)  {}
func (this *annotation) SetTextColor(y, u, v uint8)        {}
func (this *annotation) SetTextJustify(hw.MMALTextJustify) {}
func (this *annotation) SetTextOffset(x, y uint32)         {}
