/*
LightwaveRF library for the Raspberry Pi written in Go.

Basic usage:

    import "github.com/jimjibone/lwgo"

    func main() {
        lwtx := lwgo.NewLwTx()
        lightOn  := lwgo.LwBuffer{0x9,0xf,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2}
        lwtx.Send(lightOn)
    }
 */
package lwgo

/*
#cgo LDFLAGS: -lwiringPi
#include <wiringPi.h>

#define byte unsigned char
#define bool int
#define true 1
#define false 0
static const byte nibbles[] = {0xF6, 0xEE, 0xED, 0xEB, 0xDE, 0xDD, 0xDB, 0xBE,
                               0xBD, 0xBB, 0xB7, 0x7E, 0x7D, 0x7B, 0x77, 0x6F};

static void sendBytes(int pin, int onval, int offval, int period,
                      int repeats, int translate,
                      byte b1, byte b2, byte b3, byte b4, byte b5,
                      byte b6, byte b7, byte b8, byte b9, byte b10)
{
    bool sending = true;
    const int buflen = 10;
    byte in_buf[] = {b1,b2,b3,b4,b5,b6,b7,b8,b9,b10};
    byte out_buf[10] = {0};

    // Should we translate the input bytes to nibbles? Probably yes.
    if (translate > 0) {
        byte i = 0;
        for (i = 0; i < buflen; i++) {
            out_buf[i] = nibbles[in_buf[i] & 0x0F];
        }
    } else {
        byte i = 0;
        for (i = 0; i < buflen; i++) {
            out_buf[i] = in_buf[i];
        }
    }

    int tx_low_count = 7;   // total number of ticks in a low (980 uSec)
    int tx_high_count = 4;  // total number of ticks in a high (560 uSec)
    int tx_trail_count = 2; //tick count to set line low (280 uSec)

    int tx_gap_count = 72; // Inter-message gap count (10.8 msec)
    //Gap multiplier byte is used to multiply gap if longer periods are needed for experimentation
    //If gap is 255 (35msec) then this to give a max of 9 seconds
    //Used with low repeat counts to find if device times out
    int tx_gap_multiplier = 0; //Gap extension byte

    int tx_repeat = 0; //counter for repeats
    int tx_toggle_count = 3;
    int tx_gap_repeat = 0;  //unsigned int

    typedef enum TxState_ {
        tx_state_idle = 0,
        tx_state_msgstart,
        tx_state_bytestart,
        tx_state_sendbyte,
        tx_state_msgend,
        tx_state_gapstart,
        tx_state_gapend
    } TxState;
    TxState tx_state = tx_state_idle;

    int tx_bit_mask = 0; // bit mask in current byte
    int tx_num_bytes = 0; // number of bytes sent

    while (sending == true)
    {
        //Set low after toggle count interrupts
        tx_toggle_count--;
        if (tx_toggle_count == tx_trail_count) {
            digitalWrite(pin, offval);
        } else if (tx_toggle_count == 0) {
            tx_toggle_count = tx_high_count; //default high pulse duration
            switch (tx_state) {
                case tx_state_idle:
                    if(sending) {
                        tx_repeat = 0;
                        tx_state = tx_state_msgstart;
                    }
                    break;
                case tx_state_msgstart:
                    digitalWrite(pin, onval);
                    tx_num_bytes = 0;
                    tx_state = tx_state_bytestart;
                    break;
                case tx_state_bytestart:
                    digitalWrite(pin, onval);
                    tx_bit_mask = 0x80;
                    tx_state = tx_state_sendbyte;
                    break;
                case tx_state_sendbyte:
                    if(out_buf[tx_num_bytes] & tx_bit_mask) {
                        digitalWrite(pin, onval);
                    } else {
                        // toggle count for the 0 pulse
                        tx_toggle_count = tx_low_count;
                    }
                    tx_bit_mask >>=1;
                    if(tx_bit_mask == 0) {
                        tx_num_bytes++;
                        if(tx_num_bytes >= buflen) {
                            tx_state = tx_state_msgend;
                        } else {
                            tx_state = tx_state_bytestart;
                        }
                    }
                    break;
                case tx_state_msgend:
                    digitalWrite(pin, onval);
                    tx_state = tx_state_gapstart;
                    tx_gap_repeat = tx_gap_multiplier;
                    break;
                case tx_state_gapstart:
                    tx_toggle_count = tx_gap_count;
                    if (tx_gap_repeat == 0) {
                        tx_state = tx_state_gapend;
                    } else {
                        tx_gap_repeat--;
                    }
                    break;
                case tx_state_gapend:
                    tx_repeat++;
                    if(tx_repeat >= repeats) {
                        //disable timer interrupt
                        //lw_timer_Stop();
                        sending = false;
                        tx_state = tx_state_idle;
                    } else {
                        tx_state = tx_state_msgstart;
                    }
                    break;
            } // end switch
        } // end if

        // Sleep for period (default 140 us)
        delayMicroseconds(period);
    } // end while (sending == true)
} // end func
*/
import "C"

import (
    "fmt"
    "errors"
)

func wiringPiSetup() error {
    if -1 == int(C.wiringPiSetup()) {
        return errors.New("lwgo::init: wiringPiSetup() failed to call")
    }
    err := C.piHiPri(C.int(99));
    if err < 0 {
        return errors.New("lwgo::init: piHiPri() failed to set thread priority")
    }
    return nil
}

func init() {
    err := wiringPiSetup()
    if err != nil {
        fmt.Println("lwgo::init: failed")
    }
}

func (lw *LwTx) setupPins() {
    C.pinMode(C.int(lw.Pin), C.OUTPUT)
    C.digitalWrite(C.int(lw.Pin), C.LOW)
    lw.setup = true
}

/*
Send LightwaveRF commands using this struct and its functions.
The best way to create this struct, with all appropriate defaults, is to do the
following e.g:
    lwtx := lwgo.NewLwTx()
 */
type LwTx struct {
    setup bool
    Pin int
    Repeats int
    Onval, Offval int
    Translate int
    Period int
}

// A 10-byte buffer containing the command you wish to send.
type LwBuffer [10]byte

// A helper struct to pull out the meaning of a LwBuffer, useful for logging.
type LwCommand struct {
    parameter string
    device int
    command string
    address []byte
    room int
}

// Get a pointer to a LwTx initialised with recommended defaults.
func NewLwTx() *LwTx {
    // Apply defaults, allowing the user to change them afterwards if needed.
    return &LwTx{
        Pin: 3,
        Repeats: 10,
        Onval: 1,
        Offval: 0,
        Translate: 1,
        Period: 140,
    }
}

// Send a constructed LwBuffer via the 433 MHz module.
func (lw *LwTx) Send(buffer LwBuffer) {
    // Check that the transmitter is setup.
    if lw.setup == false {
        lw.SetupPins()
    }

    //fmt.Println("LwTx::Run: send:", buffer)

    // Send the message.
    C.sendBytes(C.int(lw.Pin), C.int(lw.Onval),
                C.int(lw.Offval), C.int(lw.Period),
                C.int(lw.Repeats), C.int(lw.Translate),
                C.byte(buffer[0]), C.byte(buffer[1]),
                C.byte(buffer[2]), C.byte(buffer[3]),
                C.byte(buffer[4]), C.byte(buffer[5]),
                C.byte(buffer[6]), C.byte(buffer[7]),
                C.byte(buffer[8]), C.byte(buffer[9]))
}

// Convert the LwBuffer to a LwCommand.
func (buf LwBuffer) Command() LwCommand {
    // parameter (2 [0,1])
    // device    (1 [2])
    // command   (1 [3])
    // address   (5 [4-8])
    // room      (1 [9])
    cmd := LwCommand{
        device: int(buf[2]),
        address: buf[4:8],
        room: int(buf[9]),
    }

    command := int(buf[3])
    param := int(buf[1])
    param += int(buf[0] << 4)

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

// Get a string version of the LwCommand.
func (cmd LwCommand) String() string {
    return fmt.Sprint("Parameter: ", cmd.parameter,
                      ", Device: ", cmd.device,
                      ", Command: ", cmd.command,
                      ", Address: ", cmd.address,
                      ", Room: ", cmd.room)
}

// Get a nicely formatted string version of the LwBuffer.
func (buf LwBuffer) String() string {
    return fmt.Sprint(buf.Command().String())
}

// Get the raw byte buffer within the LwBuffer.
func (buf LwBuffer) Raw() []byte {
    out := make([]byte, len(buf))
    for i, val := range buf {
        out[i] = val
    }
    return out
}
