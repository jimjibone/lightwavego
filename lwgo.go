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

// Send a constructed LwBuffer via the 433 MHz module.
func (lw *LwTx) SendBuffer(buffer LwBuffer) error {
    return lw.sendBuffer(buffer)
}

// Send a constructed LwCommand via the 433 MHz module.
func (lw *LwTx) SendCommand(command LwCommand) error {
    return lw.SendBuffer(command.Buffer())
}
