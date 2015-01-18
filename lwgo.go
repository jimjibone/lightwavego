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
#include <pigpio.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#define byte unsigned char
#define bool int
#define true 1
#define false 0

// Our possible nibbles when translating.
static const byte nibbles[] = {0xF6, 0xEE, 0xED, 0xEB, 0xDE, 0xDD, 0xDB, 0xBE,
                               0xBD, 0xBB, 0xB7, 0x7E, 0x7D, 0x7B, 0x77, 0x6F};

// Wave generation variables.
static const int max_pulses = 1600;
static gpioPulse_t pulses[1600];
static int pulse_count = 0;
static int wave_duration = 0;

// Byte buffers with message contents.
static const int buffer_size = 10;
static byte input_buffer[10] = {0};
static byte output_buffer[10] = {0};

///
/// Initialise the pigpio library.
/// @return True if initialisation was a success.
///
static bool init_gpio()
{
    bool ok = true;

    int errorCode = gpioInitialise();
    if (errorCode < 0) {
        ok = false;
        printf("C::init_gpio: gpioInitialise failed with code %d\n", errorCode);
    }

    return ok;
}

///
/// Reset the waveform pulses ready for generation of a new wave.
///
static void reset_pulses()
{
    wave_duration = 0;
    pulse_count = 0;
    pulses[0].gpioOn  = 0;
    pulses[0].gpioOff = 0;
    pulses[0].usDelay = 0;
    pulse_count++;
}

///
/// Add a new pulse to the waveform. This will add a new pulse to the wave if it
/// differs from the previous pulse. If the on and off masks are identical to
/// the previous pulse then the us_delay will be added to the previous pulse
/// only. This helps to reduce the number of pulses required in the array.
/// @return True if the new pulse was added (or not) successfully.
///
static bool add_pulse(int on_mask, int off_mask, int us_delay, bool invert)
{
    bool ok = true;

    // Set the next pulse with the new masks.
    if (invert == true)
    {
        pulses[pulse_count].gpioOn  = off_mask;
        pulses[pulse_count].gpioOff = on_mask;
    }
    else
    {
        pulses[pulse_count].gpioOn  = on_mask;
        pulses[pulse_count].gpioOff = off_mask;
    }


    // Check if the next pulse differs from the current.
    if (pulses[pulse_count-1].gpioOn  != pulses[pulse_count].gpioOn ||
        pulses[pulse_count-1].gpioOff != pulses[pulse_count].gpioOff)
    {
        // The next pulse is different to the current.
        // Set the next pulses delay and increment the pulse count.
        pulses[pulse_count].usDelay = us_delay;
        pulse_count++;
    }
    else
    {
        // The next pulse is the same as the current so we just bump up the
        // delay of the current pulse.
        pulses[pulse_count-1].usDelay += us_delay;
    }

    // Bump up the total waveform duration.
    wave_duration += us_delay;

    if (pulse_count >= max_pulses) {
        ok = false;
        printf("C::add_pulse: ERROR: Pulse count of %d is above max of %d\n", pulse_count, max_pulses);
    }

    return ok;
}

///
/// Just add a delay to the current pulse.
///
static void add_delay(int us_delay)
{
    // This is easy. Just add delay to the current pulse.
    pulses[pulse_count-1].usDelay += us_delay;
}

///
/// Send all the pulses in the waveform to the specified pin.
/// @param  pin The Broadcom number of the GPIO pin to transmit the pulses over.
/// @return     True if the sending of pulses was a success.
///
static bool send_pulses(int pin)
{
    bool ok = true;
    int lastMode = gpioGetMode(pin);
    int errorCode = 0;
    int wave_id = 0;

    if (ok)
    {
        errorCode = gpioSetMode(pin, PI_OUTPUT);
        if (errorCode != 0) {
            ok = false;
            printf("C::send_pulses: ERROR: gpioSetMode failed with code %d\n", errorCode);
        }
    }

    if (ok)
    {
        errorCode = gpioWaveClear();
        if (errorCode != 0) {
            ok = false;
            printf("C::send_pulses: ERROR: gpioWaveClear failed with code %d\n", errorCode);
        }
    }

    if (ok)
    {
        errorCode = gpioWaveAddGeneric(pulse_count, pulses);
        if (errorCode != pulse_count) {
            ok = false;
            printf("C::send_pulses: ERROR: gpioWaveAddGeneric failed with code %d and %d pulses\n", errorCode, pulse_count);
        }
    }

    if (ok)
    {
        wave_id = gpioWaveCreate();
        if (wave_id < 0)
        {
            // Wave create failed.
            ok = false;
            printf("C::send_pulses: ERROR: gpioWaveCreate failed with code %d, cause: ", wave_id);
            switch (wave_id)
            {
                case PI_EMPTY_WAVEFORM:
                    printf("PI_EMPTY_WAVEFORM\n");
                    break;
                case PI_NO_WAVEFORM_ID:
                    printf("PI_NO_WAVEFORM_ID\n");
                    break;
                case PI_TOO_MANY_CBS:
                    printf("PI_TOO_MANY_CBS\n");
                    break;
                case PI_TOO_MANY_OOL:
                    printf("PI_TOO_MANY_OOL\n");
                    break;
                default:
                    printf("NO CAUSE FOUND\n");
                    break;
            }
        }
        else
        {
            // Send the waveform!
            //errorCode = gpioWaveTxSend(wave_id, PI_WAVE_MODE_REPEAT);
            errorCode = gpioWaveTxSend(wave_id, PI_WAVE_MODE_ONE_SHOT);

            if (errorCode == PI_BAD_WAVE_ID ||
                errorCode == PI_BAD_WAVE_MODE)
            {
                ok = false;
                printf("C::send_pulses: ERROR: gpioWaveTxSend failed with code %d, cause: ", errorCode);
                switch (errorCode)
                {
                    case PI_BAD_WAVE_ID:
                        printf("PI_BAD_WAVE_ID\n");
                        break;
                    case PI_BAD_WAVE_MODE:
                        printf("PI_BAD_WAVE_MODE\n");
                        break;
                    default:
                        printf("NO CAUSE FOUND\n");
                        break;
                }
            }
            else
            {
                // Wait for the waveform to finish transmitting.
                int wait = 0;
                int duration = wave_duration;
                if (duration <= 0)
                {
                    duration = 100; // set the minimum duration to 100 us.
                    printf("C::send_pulses: Sending waveform with duration %d us, was %d.\n", duration, wave_duration);
                }

                while (gpioWaveTxBusy())
                {
                    gpioDelay(duration);
                    wait++;
                    if (wait > 10) {
                        errorCode = gpioWaveTxStop();
                        if (errorCode != 0)
                        {
                            printf("\nC::send_pulses: forced gpioWaveTxStop done with code %d\n", errorCode);
                        }
                    }
                }

                errorCode = gpioWaveTxStop();
                if (errorCode != 0)
                {
                    ok = false;
                    printf("C::send_pulses: ERROR: gpioWaveTxStop failed with code %d\n", errorCode);
                }

                errorCode = gpioWaveDelete(wave_id);
                if (errorCode != 0)
                {
                    ok = false;
                    printf("C::send_pulses: ERROR: gpioWaveTxDelete failed with code %d\n", errorCode);
                }
            } // end if (send ok)
        }  // end if (wave ok)
    } // end if (ok)

    if (ok)
    {
        // Reset the pin mode.
        errorCode = gpioSetMode(pin, lastMode);
        if (errorCode != 0) {
            ok = false;
            printf("C::send_pulses: ERROR: gpioSetMode failed with code %d\n", errorCode);
        }
    }

    return ok;
}

///
/// Create the waveform from the input buffer, then send it over the pin
/// specified to the 433 MHz transmitter.
/// @param  pin          The Broadcom number of the pin to transmit over.
/// @param  period       The inter-bit period to use for transmission.
/// @param  repeats      The number of times to repeat transmission.
/// @param  translate    Translate the input buffer to nibbles?
/// @param  invert       Invert the bits?
/// @param  input_buffer The buffer to transmit.
/// @param  input_length The length of the buffer.
/// @return              True if wave creation and transmission was a success.
///
static bool send_bytes(int pin, int period, int repeats, bool translate,
                       bool invert, byte input_buffer[], int input_length)
{
    bool ok = true;
    bool building = true;

    // Check that the input buffer is big enough.
    if (input_length != buffer_size)
    {
        ok = false;
        printf("C::send_pulses: ERROR: The input buffer size (%d) is not equal to the output buffer size (%d).\n", input_length, buffer_size);
    }

    if (ok)
    {
        // Reset the pulse array for our new waveform.
        reset_pulses();

        // Copy the input_buffer into output buffer.
        memcpy(output_buffer, input_buffer, buffer_size);

        // Should we translate the input bytes to nibbles? Probably yes.
        if (translate == true) {
            for (int i = 0; i < buffer_size; i++) {
                output_buffer[i] = nibbles[input_buffer[i] & 0x0F];
            }
        }

        // Prepare state variables for building the waveform.
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
        int tx_bit_mask = 0; // bit mask in current byte
        int tx_num_bytes = 0; // number of bytes sent

        while (building == true && ok == true)
        {
            //Set low after toggle count interrupts
            tx_toggle_count--;
            if (tx_toggle_count == tx_trail_count) {
                // Add an OFF pulse.
                ok = add_pulse(0, (1<<pin), period, invert);
            } else if (tx_toggle_count == 0) {
                tx_toggle_count = tx_high_count; //default high pulse duration
                switch (tx_state) {
                    case tx_state_idle:
                        if(building) {
                            tx_repeat = 0;
                            tx_state = tx_state_msgstart;
                        }
                        break;
                    case tx_state_msgstart:
                        // Add an ON pulse.
                        ok = add_pulse((1<<pin), 0, period, invert);

                        tx_num_bytes = 0;
                        tx_state = tx_state_bytestart;
                        break;
                    case tx_state_bytestart:
                        // Add an ON pulse.
                        ok = add_pulse((1<<pin), 0, period, invert);

                        tx_bit_mask = 0x80;
                        tx_state = tx_state_sendbyte;
                        break;
                    case tx_state_sendbyte:
                        if(output_buffer[tx_num_bytes] & tx_bit_mask) {
                            // Add an ON pulse.
                            ok = add_pulse((1<<pin), 0, period, invert);
                        } else {
                            // toggle count for the 0 pulse
                            tx_toggle_count = tx_low_count;
                        }
                        tx_bit_mask >>=1;
                        if(tx_bit_mask == 0) {
                            tx_num_bytes++;
                            if(tx_num_bytes >= buffer_size) {
                                tx_state = tx_state_msgend;
                            } else {
                                tx_state = tx_state_bytestart;
                            }
                        }
                        break;
                    case tx_state_msgend:
                        // Add an ON pulse.
                        ok = add_pulse((1<<pin), 0, period, invert);

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
                            // We have finished adding repeats so stop building
                            // the waveform now.
                            building = false;
                            tx_state = tx_state_idle;
                        } else {
                            tx_state = tx_state_msgstart;
                        }
                        break;
                } // end switch
            } // end if (toggle count states)
            else
            {
                // Just add delay.
                add_delay(period);
            }
        } // end while (building == true && ok == true)
    } // end if (ok)

    if (ok)
    {
        // Send the waveform.
        ok = send_pulses(pin);
    }

    return ok;
}
*/
import "C"

import (
    "fmt"
    "time"
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
type LwMessage struct {
    Buffer [10]byte
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
        Pin: 22,
        Repeats: 10,
        Invert: false,
        Translate: true,
        Period: 140,
    }
}

func byteSliceToCArray(byteSlice [10]byte) unsafe.Pointer {
    var array = unsafe.Pointer(C.calloc(C.size_t(len(byteSlice)), 1))
    var arrayptr = uintptr(array)

    for i := 0; i < len(byteSlice); i ++ {
        *(*C.byte)(unsafe.Pointer(arrayptr)) = C.byte(byteSlice[i])
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
    buffer := byteSliceToCArray(message.Buffer)
    defer C.free(buffer)

    // Send the message.
    result := int(C.send_bytes(C.int(lw.Pin),
                               C.int(lw.Period),
                               C.int(lw.Repeats),
                               C.int(translate),
                               C.int(invert),
                               (*C.byte)(buffer),
                               C.int(len(message.Buffer))))

    if result == 0 {
        fmt.Println("lwgo::Send: send FAIL!")
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

        message := LwMessage{}
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
    return fmt.Sprint(message.command().String())
}

// Raw returns the raw byte buffer stored within the LwMessage.
func (message LwMessage) Raw() []byte {
    out := make([]byte, len(message.Buffer))
    for i, val := range message.Buffer {
        out[i] = val
    }
    return out
}
