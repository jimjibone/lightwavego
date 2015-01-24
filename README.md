# lwgo

[![GoDoc](https://godoc.org/github.com/jimjibone/lwgo?status.svg)](https://godoc.org/github.com/jimjibone/lwgo) [![Build Status](https://travis-ci.org/jimjibone/lwgo.svg?branch=master)](https://travis-ci.org/jimjibone/lwgo)

LightwaveRF library for the Raspberry Pi written in Go.

See [my C library](https://github.com/jimjibone/LightwaveRF) for the Raspberry Pi from which this is based.


## Connecting your modules

### TX

- Connect the data pin (to the default) Broadcom pin 22. This is GPIO pin 15 on the Pi.
- Connect the VCC to the Pi's 3.3V and GND to the Pi's GND.
- For best performance add an antenna to your TX module. This can just be a quarter-wavelength wire soldered to the module. 433 MHz quarter wavelength is 173 mm but I found best performance with my module using a length of 165 mm.


## Installation

- [pigpio](http://abyz.co.uk/rpi/pigpio/download.html): This is required to control the GPIO pins on the Raspberry Pi
- `go get github.com/jimjibone/lwgo`
- `go install github.com/jimjibone/lwgo`


## Blink example

- `git clone https://github.com/jimjibone/lwgo.git`
- `cd lwgo`
- `go get github.com/jimjibone/lwgo` (if not done already)
- `go install github.com/jimjibone/lwgo` (if not done already)
- `go build examples/blink.co`
- `sudo ./blink` (sudo required for GPIO access)


## Development

For development there are some other requirements. Don't worry about these if you are just building to run.

### `go generate`

To automatically create a nice human friendly representation of the various constants within `types.go` we will use `stringer` along with `go generate`. When we call `go generate` on the command line, `generate` will automatically find the stringer command within the types.go file and execute it. It will create a new go file that will help translate the constants into strings, easy.

- `go get golang.org/x/tools/cmd/stringer`
- `go generate`
- Build away.
