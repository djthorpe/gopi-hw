# Go parameters
GOCMD=go
GOFLAGS=-tags "rpi"
GOINSTALL=$(GOCMD) install $(GOFLAGS)
GOTEST=$(GOCMD) test $(GOFLAGS) 
GOCLEAN=$(GOCMD) clean

# Freetype parameters
FT_CFLAGS=-I/usr/include/freetype2
FT_LDFLAGS=-lfreetype

all: test_rpi test_freetype install

install:
	$(GOINSTALL) ./cmd/gpio_ctrl
	$(GOINSTALL) ./cmd/hw_list
	$(GOINSTALL) ./cmd/hw_service
	$(GOINSTALL) ./cmd/i2c_detect
	$(GOINSTALL) ./cmd/lirc_receive
	$(GOINSTALL) ./cmd/mmal_camera_preview
	$(GOINSTALL) ./cmd/pwm_ctrl
	$(GOINSTALL) ./cmd/spi_ctrl

test_rpi:
	$(GOTEST) -v ./rpi

test_freetype:
	CGO_CFLAGS="${FT_CFLAGS}" CGO_LDFLAGS="${FT_LDFLAGS}" $(GOTEST) -v ./freetype

clean: 
	$(GOCLEAN)
