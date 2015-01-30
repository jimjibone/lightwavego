# lightwavego

[![GoDoc](https://godoc.org/github.com/jimjibone/lightwavego?status.svg)](https://godoc.org/github.com/jimjibone/lightwavego) [![Build Status](https://travis-ci.org/jimjibone/lightwavego.svg?branch=master)](https://travis-ci.org/jimjibone/lightwavego)

LightwaveRF library for the Raspberry Pi written in Go.

See [my C library](https://github.com/jimjibone/LightwaveRF) for the Raspberry Pi from which this is based.


## Connecting your modules

### TX

- Connect the data pin (to the default) Broadcom pin 22. This is GPIO pin 15 on the Pi.
- Connect the VCC to the Pi's 3.3V and GND to the Pi's GND.
- For best performance add an antenna to your TX module. This can just be a quarter-wavelength wire soldered to the module. 433 MHz quarter wavelength is 173 mm but I found best performance with my module using a length of 165 mm.


## Installation

- [pigpio](http://abyz.co.uk/rpi/pigpio/download.html): This is required to control the GPIO pins on the Raspberry Pi
- `go get github.com/jimjibone/lightwavego`
- `go get github.com/ant0ine/go-json-rest/rest` (required for `lwserver` only)
- `go install github.com/jimjibone/lightwavego/examples/...`
- `lwblink` and `lwserver` will now both be in your PATH


## Examples

### Blink

This is basically the blink example we all know from microelectronics except that this time it's using your houses lights and it is able to dim them.

- `git clone https://github.com/jimjibone/lightwavego.git`
- `cd lightwavego`
- `go get github.com/jimjibone/lightwavego` (if not done already)
- `go install github.com/jimjibone/lightwavego` (if not done already)
- `go build examples/lwblink/lwblink.go`
- `sudo ./lwblink` (sudo required for GPIO access)


### Server

This sets up a simple RESTful JSON server that is able to receive byte buffers of pre-compiled commands and use the library to broadcast them over your transmitter.

- Follow the instructions above for the Blink example, but now do:
- `go build examples/lwserver/lwserver.go`
- `sudo ./lwserver` (sudo also required here for the GPIO access)
- `curl -i -d '{"Buffer":"090f0301050903000102"}' http://localhost:8080/send` will turn a light on to max brightness


## Development

For development there are some other requirements. Don't worry about these if you are just building to run.

### `go generate`

To automatically create a nice human friendly representation of the various constants within `types.go` we will use `stringer` along with `go generate`. When we call `go generate` on the command line, `generate` will automatically find the stringer command within the types.go file and execute it. It will create a new go file that will help translate the constants into strings, easy.

- `go get golang.org/x/tools/cmd/stringer`
- `go generate`
- Build away.
