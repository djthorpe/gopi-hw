# gopi-hw

This repository contains hardware implementations and hardware components for
[gopi](http://github.com/djthorpe/gopi) and some example programs in the `cmd` folder.
The repository depends on golang version 1.11 and above (in order to support modules). For
Linux support, Raspbian GNU/Linux v9 has been tested and for Mac support.

## Components

The gopi components provided by this repository are:

| Component Path | Plaform/Tag      | Description                            | Conforms to   |
| -------------- | ---------------- | -------------------------------------- |-------------- |
| sys/filepoll   | linux            | Watch for changes to files and folders | 
| sys/gpio       | linux,rpi        | General Purpose Hardware Input/Output  | gopi.GPIO     |
| sys/hw         | linux,rpi,darwin | Hardware information, capabilities     | gopi.Hardware | 
| sys/i2c        | linux            | I2C interface                          | gopi.I2C      |
| sys/lirc       | linux            | Linux IR control (LIRC) interface      | gopi.LIRC     |
| sys/mmal       | rpi              | Multimedia Abstraction Layer           | hw.MMAL       |
| sys/pwm        | rpi              | Pulse Wide Modulation (PWM) interface  | gopi.PWM      |
| sys/spi        | linux            | SPI interface                          | gopi.SPI      |

## Bindings

In addition to these general interfaces, the repository contains golang bindings for underlying
libraries:

| Package  | Description                                         |
| -------- | --------------------------------------------------- |
| rpi      | Raspberry Pi bindings for DispmanX, MMAL, Videocore |
| egl      | Bindings for Khronos Group OpenEGL                  |
| freetype | Bindings for Freetype2                              |

## Building and installing examples

In order to build, `pkg-config` is required. For the Raspberry Pi, the configuration files are located
under `/opt/vc/lib/pkg-config` so set the environment variable as follows:

```
bash% PKG_CONFIG_PATH="/opt/vc/lib/pkgconfig"
bash% go install -tags rpi ./cmd/...
```

There is a makefile which can be used for testing and installing bindings and examples, on a per-platform
basis:

```
bash% make linux   # makes for generic linux
bash% make darwin  # makes for MacOS
bash% make rpi     # makes for Raspberry Pi
```

The resulting binaries are as follows. Use the `-help` flag to see the different options for each:

  * `hw_list` Provide information on hardware capabilities
  * `gpio_ctrl` Control the GPIO interface
  * `i2c_detect` Detect I2C devices
  * `lirc_receive` Display IR pulses from an IR device
  * `pwm_ctrl` Control PWM signals on the GPIO interface
  * `spi_ctrl` Control SPI communication
  * `mmal_camera_preview` Preview the camera output on the screen
  * `mmal_encode_image` Demonstrates image decoding and encoding using the GPU
  * `mmal_video_preview` Demonstrates playback of a H264 video on the screen using the GPU

## Construction Ahead

More information about the individual components implemented in this repository will be
forthcoming later!

  
