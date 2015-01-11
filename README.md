# lwgo

LightwaveRF library for the Raspberry Pi written in Go. [![GoDoc](https://godoc.org/github.com/jimjibone/lwgo?status.svg)](https://godoc.org/github.com/jimjibone/lwgo)

See [my C library](https://github.com/jimjibone/LightwaveRF) for the Raspberry Pi from which this is based.


## Connecting your modules

### TX

- Connect the data pin (to the default) WiringPi pin 3. This is GPIO pin 15 on the Pi.
- Connect the VCC to the Pi's 3.3V and GND to the Pi's GND.
- For best performance add an antenna to your TX module. This can just be a quarter-wavelength wire soldered to the module. 433 MHz quarter wavelength is 173 mm but I found best performance with my module using a length of 165 mm.


## Installation

- [WiringPi](http://wiringpi.com/download-and-install/): This is required to control the GPIO pins on the Raspberry Pi
- `go get github.com/jimjibone/lwgo`
- `go install github.com/jimjibone/lwgo`


## Blink example

- `git clone https://github.com/jimjibone/lwgo.git`
- `cd lwgo`
- `go get github.com/jimjibone/lwgo` (if not done already)
- `go install github.com/jimjibone/lwgo` (if not done already)
- `go build examples/blink.co`
- `sudo ./blink` (sudo required for WiringPi GPIO access)
