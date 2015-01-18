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

// LwMessage contains the command you wish to send along with repeat configuration.
type LwMessage [10]byte

// lwCommand is a helper struct to pull out the meaning of a LwMessage, useful
// for logging.
type lwCommand struct {
    parameter string
    device int
    command string
    address []byte
    room int
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

// Convert the LwMessage contents to a C array that we can pass to our C
// counterpart.
func messageToCArray(message LwMessage) unsafe.Pointer {
    var array = unsafe.Pointer(C.calloc(C.size_t(len(message)), 1))
    var arrayptr = uintptr(array)

    for i := 0; i < len(message); i ++ {
        *(*C.byte)(unsafe.Pointer(arrayptr)) = C.byte(message[i])
        arrayptr++
    }

    return array
}

// Send a constructed LwMessage via the 433 MHz module.
func (lw *LwTx) Send(message LwMessage) {
    // Convert the boolean types to ints for passing to C.
    var translate, invert int
    if lw.Translate {
        translate = 1
    }
    if lw.Invert {
        invert = 1
    }

    // Create a C byte array from the Go byte slice.
    buffer := messageToCArray(message)
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

func NewMessage(buffer []byte) (LwMessage, error) {
    if len(buffer) <= 0 {
        return LwMessage{}, fmt.Errorf("input buffer size is too small: %d", len(buffer))
    } else if len(buffer) > 10 {
        return LwMessage{}, fmt.Errorf("input buffer size is too big: %d", len(buffer))
    } else {
        message := LwMessage{}
        for i, val := range buffer {
            message[i] = val;
        }
        return message, nil
    }
}

// Convert the LwMessage to a lwCommand.
func (message LwMessage) command() lwCommand {
    // parameter (2 [0,1])
    // device    (1 [2])
    // command   (1 [3])
    // address   (5 [4-8])
    // room      (1 [9])
    cmd := lwCommand{
        device: int(message[2]),
        address: message[4:8],
        room: int(message[9]),
    }

    command := int(message[3])
    param := int(message[1])
    param += int(message[0] << 4)

    // Get the parameter
    switch {
        // Command off
        case command == 0 && param >= 0 && param <= 127: {
            cmd.parameter = "off"
        }
        case command == 0 && param >= 128 && param <= 159: {
            cmd.parameter = fmt.Sprint("off with level:", param-128)
        }
        case command == 0 && param >= 160 && param <= 191: {
            cmd.parameter = "decrease brightness"
        }
        case command == 0 && param >= 192 && param <= 255: {
            cmd.parameter = "all off"
        }

        // Command on
        case command == 1 && param >= 0 && param <= 31: {
            cmd.parameter = "on to last level"
        }
        case command == 1 && param >= 32 && param <= 63: {
            cmd.parameter = fmt.Sprint("on with level:", param-32)
        }
        case command == 1 && param >= 64 && param <= 95: {
            cmd.parameter = fmt.Sprint("on with level:", param-64)
        }
        case command == 1 && param >= 96 && param <= 127: {
            cmd.parameter = fmt.Sprint("on with level:", param-96)
        }
        case command == 1 && param >= 128 && param <= 159: {
            cmd.parameter = fmt.Sprint("on with level:", param-128)
        }
        case command == 1 && param >= 160 && param <= 191: {
            cmd.parameter = "increase brightness"
        }
        case command == 1 && param >= 192 && param <= 223: {
            cmd.parameter = fmt.Sprint("set all to level:", param-192)
        }
        case command == 1 && param >= 224 && param <= 255: {
            cmd.parameter = fmt.Sprint("set all to level:", param-224)
        }

        // Command mood
        case command == 2 && param >= 130 && param <= 255: {
            cmd.parameter = fmt.Sprint("start mood:", param-129)
        }
        case command == 2 && param >= 2 && param <= 129: {
            cmd.parameter = fmt.Sprint("define mood:", param-1)
        }

        default: cmd.parameter = "unknown"
    }

    // Get the command
    switch command {
        case 0: cmd.command = "off"
        case 1: cmd.command = "on"
        case 2: cmd.command = "mood"
        default: cmd.command = "unknown"
    }

    return cmd
}

// Get a string version of the lwCommand.
func (cmd lwCommand) String() string {
    return fmt.Sprint("Parameter: ", cmd.parameter,
                      ", Device: ", cmd.device,
                      ", Command: ", cmd.command,
                      ", Address: ", cmd.address,
                      ", Room: ", cmd.room)
}

// Get a nicely formatted string version of the LwMessage.
func (message LwMessage) String() string {
    return fmt.Sprint(message.command().String())
}

// Raw returns the raw byte buffer stored within the LwMessage.
func (message LwMessage) Raw() []byte {
    out := make([]byte, len(message))
    for i, val := range message {
        out[i] = val
    }
    return out
}
