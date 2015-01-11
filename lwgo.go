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
    "time"
)

func wiringPiSetup() error {
    if -1 == int(C.wiringPiSetup()) {
        return fmt.Errorf("lwgo::init: wiringPiSetup() failed to call")
    }
    err := C.piHiPri(C.int(99));
    if err < 0 {
        return fmt.Errorf("lwgo::init: piHiPri() failed to set thread priority")
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
LwTx contains the configuration of your LightwaveRF setup.
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

// LwMessage contains the command you wish to send along with repeat configuration.
type LwMessage struct {
    Buffer [10]byte
    Repeats int // repeat transmission more than once?
    Period time.Duration // for repeats
}

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
        Pin: 3,
        Repeats: 10,
        Onval: 1,
        Offval: 0,
        Translate: 1,
        Period: 140,
    }
}

// Send a constructed LwMessage via the 433 MHz module.
func (lw *LwTx) Send(message LwMessage) {
    // Check that the transmitter is setup.
    if lw.setup == false {
        lw.setupPins()
    }

    //fmt.Println("LwTx::Run: send:", message)

    for count := message.Repeats; count >= 0; count-- {
        // Send the message.
        C.sendBytes(C.int(lw.Pin), C.int(lw.Onval),
                    C.int(lw.Offval), C.int(lw.Period),
                    C.int(lw.Repeats), C.int(lw.Translate),
                    C.byte(message.Buffer[0]), C.byte(message.Buffer[1]),
                    C.byte(message.Buffer[2]), C.byte(message.Buffer[3]),
                    C.byte(message.Buffer[4]), C.byte(message.Buffer[5]),
                    C.byte(message.Buffer[6]), C.byte(message.Buffer[7]),
                    C.byte(message.Buffer[8]), C.byte(message.Buffer[9]))

        // Wait if there are remaining repeats.
        if count > 0 {
            time.Sleep(message.Period)
        }
    }
}

func NewMessage(buffer []byte, repeats int, period time.Duration) (LwMessage, error) {
    if len(buffer) <= 0 {
        return LwMessage{}, fmt.Errorf("input buffer size is too small: %d", len(buffer))
    } else if len(buffer) > 10 {
        return LwMessage{}, fmt.Errorf("input buffer size is too big: %d", len(buffer))
    } else {
        if repeats < 0 {
            repeats = 0
        }
        if period < 100 * time.Millisecond {
            period = 100 * time.Millisecond
        }

        message := LwMessage{Repeats: repeats, Period: period}
        for i, val := range buffer {
            message.Buffer[i] = val;
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
        device: int(message.Buffer[2]),
        address: message.Buffer[4:8],
        room: int(message.Buffer[9]),
    }

    command := int(message.Buffer[3])
    param := int(message.Buffer[1])
    param += int(message.Buffer[0] << 4)

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
    return fmt.Sprint(message.command().String(),
                      ", Repeats: ", message.Repeats,
                      ", Period: ", message.Period)
}

// Raw returns the raw byte buffer stored within the LwMessage.
func (message LwMessage) Raw() []byte {
    out := make([]byte, len(message.Buffer))
    for i, val := range message.Buffer {
        out[i] = val
    }
    return out
}
