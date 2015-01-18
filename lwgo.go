/*
Package lwgo is a LightwaveRF library for the Raspberry Pi.

Basic usage:

    import (
        "github.com/jimjibone/lwgo"
        "fmt"
    )

    func main() {
        lwtx := lwgo.NewLwTx()
        lightOn, err := lwgo.NewMessage([]byte{0x9,0xf,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2}, 2, time.Millisecond * 500)
        if err != nil {
            fmt.Error(err)
        }
        lwtx.Send(lightOn)
    }
 */
package lwgo

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lpigpio -lpthread -lrt
#include "lwgo.h"
*/
import "C"

import (
    "fmt"
    "unsafe"
)

func init() {
    if int(C.init_gpio()) == 0 {
        C.gpioTerminate()
        fmt.Println("lwgo::init: failed")
    }
}

func Shutdown() {
    C.gpioTerminate()
}

/*
LwTx contains the configuration of your LightwaveRF setup.
The best way to create this struct, with all appropriate defaults, is to do the
following e.g:
    lwtx := lwgo.NewLwTx()
 */
type LwTx struct {
    Pin int         // 22
    Repeats int     // 10
    Invert bool     // false
    Translate bool  // true
    Period int      // 140 (us)
}

// NewLwTx returns a pointer to a LwTx struct initialised with the recommended
// defaults.
func NewLwTx() *LwTx {
    // Apply defaults, allowing the user to change them afterwards if needed.
    return &LwTx{
        Pin: 22,
        Repeats: 10,
        Invert: false,
        Translate: true,
        Period: 140,
    }
}

// Send a constructed LwMessage via the 433 MHz module.
func (lw *LwTx) Send(message LwMessage) { // TODO: return `error`
    // Convert the boolean types to ints for passing to C.
    var translate, invert int
    if lw.Translate {
        translate = 1
    }
    if lw.Invert {
        invert = 1
    }

    // Create a C byte buffer from the Go byte slice.
    var buffer = unsafe.Pointer(C.calloc(C.size_t(len(message)), 1))
    var bufferptr = uintptr(buffer)
    for i := 0; i < len(message); i ++ {
        *(*C.byte)(unsafe.Pointer(bufferptr)) = C.byte(message[i])
        bufferptr++
    }
    defer C.free(buffer)

    // Send the message.
    result := int(C.send_bytes(C.int(lw.Pin),
                               C.int(lw.Period),
                               C.int(lw.Repeats),
                               C.int(translate),
                               C.int(invert),
                               (*C.byte)(buffer),
                               C.int(len(message))))

    if result == 0 {
        fmt.Println("lwgo::Send: send FAIL!")
    }
}

// Send a constructed LwCommand via the 433 MHz module.
func (lw *LwTx) SendCommand(command LwCommand) { // TODO: return `error`
    lw.Send(command.Message())
}
