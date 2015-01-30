package main

import (
    "fmt"
    "time"
    "github.com/jimjibone/lightwavego"
)

func main() {
    fmt.Println("lightwavego blink!")
    defer lightwavego.Shutdown()

    // Set up the LightwaveRF TX instance.
    lwtx := lightwavego.NewLwTx()

    // Define some messages to test with.
    lightOn, _ := lightwavego.NewBuffer([]byte{0x9,0xf,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2}) // dim to max
    lightDim, _ := lightwavego.NewBuffer([]byte{0x5,0x4,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2}) // dim to 10
    lightOff, _ := lightwavego.NewBuffer([]byte{0x4,0x0,0x3,0x0,0x5,0x9,0x3,0x0,0x1,0x2}) // off
    lightLast, _ := lightwavego.NewBuffer([]byte{0x0,0x0,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2}) // on (to last level)

    // Define some commands to test them also.
    commandOn := lightwavego.LwCommand{
        Command: lightwavego.Dim,
        Value: 31,
        Device: 3,
        Address: []byte{5, 9, 3, 0, 1},
        Room: 2,
    }

    commandDim := lightwavego.LwCommand{
        Command: lightwavego.Dim,
        Value: 20,
        Device: 3,
        Address: []byte{5, 9, 3, 0, 1},
        Room: 2,
    }

    commandOff := lightwavego.LwCommand{
        Command: lightwavego.Off,
        Value: 0,
        Device: 3,
        Address: []byte{5, 9, 3, 0, 1},
        Room: 2,
    }

    commandLast := lightwavego.LwCommand{
        Command: lightwavego.On,
        Value: 0,
        Device: 3,
        Address: []byte{5, 9, 3, 0, 1},
        Room: 2,
    }

    // Send a couple of these messages.
    for i := 0; i < 1; i++ {
        fmt.Println("lightOn:", lightOn)
        lwtx.SendBuffer(lightOn)
        time.Sleep(4 * time.Second)

        fmt.Println("lightDim:", lightDim)
        lwtx.SendBuffer(lightDim)
        time.Sleep(4 * time.Second)

        fmt.Println("lightOff:", lightOff)
        lwtx.SendBuffer(lightOff)
        time.Sleep(4 * time.Second)

        fmt.Println("lightLast:", lightLast)
        lwtx.SendBuffer(lightLast)
        time.Sleep(4 * time.Second)

        fmt.Println("commandOn:", commandOn)
        lwtx.SendCommand(commandOn)
        time.Sleep(4 * time.Second)

        fmt.Println("commandDim:", commandDim)
        lwtx.SendCommand(commandDim)
        time.Sleep(4 * time.Second)

        fmt.Println("commandOff:", commandOff)
        lwtx.SendCommand(commandOff)
        time.Sleep(4 * time.Second)

        fmt.Println("commandLast:", commandLast)
        lwtx.SendCommand(commandLast)
        time.Sleep(4 * time.Second)
    }

    // Turn your lights back on in case you're plunged into darkness like me.
    fmt.Println("commandOn:", commandOn)
    lwtx.SendCommand(commandOn)

    fmt.Println("All done")
}
