# Go parameters
GOCMD=go
GOFLAGS=-tags "rpi"
GOINSTALL=$(GOCMD) install $(GOFLAGS)
GOTEST=$(GOCMD) test $(GOFLAGS) 
GOCLEAN=$(GOCMD) clean

all: test_rpi test_egl test_freetype test_openvg install

pkg-config:
	PKG_CONFIG_PATH="/opt/vc/lib/pkgconfig"

install: pkg-config
	$(GOINSTALL) ./cmd/gpio_ctrl
	$(GOINSTALL) ./cmd/hw_list
	$(GOINSTALL) ./cmd/hw_service
	$(GOINSTALL) ./cmd/i2c_detect
	$(GOINSTALL) ./cmd/lirc_receive
	$(GOINSTALL) ./cmd/mmal_camera_preview
	$(GOINSTALL) ./cmd/mmal_encode_image
	$(GOINSTALL) ./cmd/pwm_ctrl
	$(GOINSTALL) ./cmd/spi_ctrl

test_rpi: pkg-config
	$(GOTEST) -v ./rpi

test_egl: pkg-config
	$(GOTEST) -v ./egl

test_dx: pkg-config
	$(GOTEST) -v ./rpi/dispmanx_test.go

test_freetype: pkg-config
	$(GOTEST) -v ./freetype

test_openvg: pkg-config
	$(GOTEST) -v ./openvg

clean: 
	$(GOCLEAN)
