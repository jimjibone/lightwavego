// This file should only build on the Raspberry Pi (linux && arm)

package lwgo

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lpigpio -lpthread -lrt
#include "interface.h"
*/
import "C"

import (
    "fmt"
    "unsafe"
)

func init() {
    if int(C.init_gpio()) == 0 {
        C.gpioTerminate()
        panic(fmt.Sprint("failed to initialise pigpio library"))
    }
}

func Shutdown() {
    C.gpioTerminate()
}

func (lw *LwTx) sendBuffer(buffer LwBuffer) error {
    // Convert the boolean types to ints for passing to C.
    var translate, invert int
    if lw.Translate {
        translate = 1
    }
    if lw.Invert {
        invert = 1
    }

    // Create a C byte buffer from the Go byte slice.
    var cbuffer = unsafe.Pointer(C.calloc(C.size_t(len(buffer)), 1))
    var cbufferptr = uintptr(cbuffer)
    for i := 0; i < len(buffer); i ++ {
        *(*C.byte)(unsafe.Pointer(cbufferptr)) = C.byte(buffer[i])
        cbufferptr++
    }
    defer C.free(cbuffer)

    // Send the buffer.
    result := int(C.send_bytes(C.int(lw.Pin),
                               C.int(lw.Period),
                               C.int(lw.Repeats),
                               C.int(translate),
                               C.int(invert),
                               (*C.byte)(cbuffer),
                               C.int(len(buffer))))

    if result == 0 {
        return fmt.Errorf("failed to send buffer")
    } else {
        return nil
    }
}
