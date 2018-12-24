# Go parameters
GOCMD=go
GOFLAGS=-tags "rpi"
GOINSTALL=$(GOCMD) install $(GOFLAGS)
GOTEST=$(GOCMD) test $(GOFLAGS) 
GOCLEAN=$(GOCMD) clean
    
all: test install

install:
	$(GOINSTALL) ./cmd/gpio_ctrl
	$(GOINSTALL) ./cmd/hw_list
	$(GOINSTALL) ./cmd/hw_service
	$(GOINSTALL) ./cmd/i2c_detect
	$(GOINSTALL) ./cmd/lirc_receive
	$(GOINSTALL) ./cmd/mmal_camera_preview
	$(GOINSTALL) ./cmd/pwm_ctrl
	$(GOINSTALL) ./cmd/spi_ctrl

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
