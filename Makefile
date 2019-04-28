# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
PKG_CONFIG_PATH="/opt/vc/lib/pkgconfig"

# App parameters
GOPI=github.com/djthorpe/gopi
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 

linux: install-linux

darwin: install-darwin
	
rpi: test-rpi test-dx test-freetype install-rpi install-mmal

install-darwin:
	$(GOINSTALL) -tags "darwin" ./cmd/hw_list/...

install-linux:
	$(GOINSTALL) -tags "linux" ./cmd/hw_list/...

install-rpi:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi" $(GOFLAGS) ./cmd/hw_list
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi" $(GOFLAGS) ./cmd/gpio_ctrl
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi" $(GOFLAGS) ./cmd/i2c_detect
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi" $(GOFLAGS) ./cmd/lirc_receive
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi" $(GOFLAGS) ./cmd/pwm_ctrl
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi" $(GOFLAGS) ./cmd/spi_ctrl

install-mmal:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi" $(GOFLAGS) ./cmd/mmal_camera_preview
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi" $(GOFLAGS) ./cmd/mmal_encode_image
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi" $(GOFLAGS) ./cmd/mmal_video_preview

test-rpi:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -tags "rpi"  -v ./rpi

test-egl:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -tags "rpi" -v ./egl

test-dx:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -tags "rpi" -v ./rpi/dispmanx_test.go

test-freetype:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -tags "rpi" -v ./freetype

test-openvg:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOTEST) -tags "rpi" -v ./openvg

clean: 
	$(GOCLEAN)
