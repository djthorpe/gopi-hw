# Go parameters
GOCMD=go
GOFLAGS=-tags "rpi"
GOINSTALL=$(GOCMD) install $(GOFLAGS)
GOTEST=$(GOCMD) test $(GOFLAGS) 
GOCLEAN=$(GOCMD) clean
PKG_CONFIG_PATH="/opt/vc/lib/pkgconfig"

	
all: test install

test: test_rpi test_egl test_freetype test_openvg

install:
	$(GOINSTALL) ./cmd/gpio_ctrl
	$(GOINSTALL) ./cmd/hw_list
	$(GOINSTALL) ./cmd/hw_service
	$(GOINSTALL) ./cmd/i2c_detect
	$(GOINSTALL) ./cmd/lirc_receive
	$(GOINSTALL) ./cmd/pwm_ctrl
	$(GOINSTALL) ./cmd/spi_ctrl
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) ./cmd/mmal_camera_preview
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) ./cmd/mmal_encode_image

test_rpi:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -v ./rpi

test_egl:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -v ./egl

test_dx:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -v ./rpi/dispmanx_test.go

test_freetype:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -v ./freetype

test_openvg:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -v ./openvg

clean: 
	$(GOCLEAN)
