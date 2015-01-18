package main

import (
    "fmt"
    "time"
    "github.com/jimjibone/lwgo"
)

func main() {
    fmt.Println("lwgo blink!")
    defer lwgo.Shutdown()

    // Set up the LightwaveRF TX instance.
    lwtx := lwgo.NewLwTx()

    // Define some messages to test with.
    lightOn, _ := lwgo.NewMessage([]byte{0x9,0xf,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2}) // dim to max
    lightDim, _ := lwgo.NewMessage([]byte{0x5,0x4,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2}) // dim to 10
    lightOff, _ := lwgo.NewMessage([]byte{0x4,0x0,0x3,0x0,0x5,0x9,0x3,0x0,0x1,0x2}) // off
    lightLast, _ := lwgo.NewMessage([]byte{0x0,0x0,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2}) // on (to last level)

    // Define some commands to test them also.
    commandOn := lwgo.LwCommand{
        Command: "Dim",
        Value: 31,
        Device: 3,
        Address: []byte{5, 9, 3, 0, 1},
        Room: 2,
    }

    commandDim := lwgo.LwCommand{
        Command: "Dim",
        Value: 20,
        Device: 3,
        Address: []byte{5, 9, 3, 0, 1},
        Room: 2,
    }

    commandOff := lwgo.LwCommand{
        Command: "Off",
        Value: 0,
        Device: 3,
        Address: []byte{5, 9, 3, 0, 1},
        Room: 2,
    }

    commandLast := lwgo.LwCommand{
        Command: "On",
        Value: 0,
        Device: 3,
        Address: []byte{5, 9, 3, 0, 1},
        Room: 2,
    }

    // Send a couple of these messages.
    for i := 0; i < 1; i++ {
        fmt.Println("lightOn:", lightOn)
        lwtx.Send(lightOn)
        time.Sleep(4 * time.Second)

        fmt.Println("lightDim:", lightDim)
        lwtx.Send(lightDim)
        time.Sleep(4 * time.Second)

        fmt.Println("lightOff:", lightOff)
        lwtx.Send(lightOff)
        time.Sleep(4 * time.Second)

        fmt.Println("lightLast:", lightLast)
        lwtx.Send(lightLast)
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
