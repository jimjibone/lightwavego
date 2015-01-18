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
    //lightOn, _ := lwgo.NewMessage([]byte{0x0,0x0,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2}) // switch on (to last level)
    lightOff, _ := lwgo.NewMessage([]byte{0x4,0x0,0x3,0x0,0x5,0x9,0x3,0x0,0x1,0x2}) // off

    // Send a couple of these messages.
    for i := 0; i < 2; i++ {
        fmt.Println("lightOn:", lightOn)
        lwtx.Send(lightOn)
        time.Sleep(4 * time.Second)

        fmt.Println("lightOff:", lightOff)
        lwtx.Send(lightOff)
        time.Sleep(4 * time.Second)
    }

    fmt.Println("all done")
}