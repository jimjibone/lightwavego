// gcc -o c_test  -lpigpio -lpthread -lrt
#include <pigpio.h>
#include <stdio.h>

#define byte unsigned char
#define bool int
#define true 1
#define false 0
const byte nibbles[] = {0xF6, 0xEE, 0xED, 0xEB, 0xDE, 0xDD, 0xDB, 0xBE,
                        0xBD, 0xBB, 0xB7, 0x7E, 0x7D, 0x7B, 0x77, 0x6F};

const max_pulses = 1600;
gpioPulse_t pulses[1600];
int pulse_count = 0;
int wave_duration = 0;

const int buffer_size = 10;
byte input_buffer[10] = {0};
byte output_buffer[10] = {0};

bool initGpio()
{
    bool ok = true;

    int errorCode = gpioInitialise();
    if (errorCode < 0) {
        ok = false;
        printf("gpioInitialise failed with code %d\n", errorCode);
    } else {
        printf("gpioInitialise success\n");
    }

    return ok;
}

void reset_pulses()
{
    wave_duration = 0;
    pulse_count = 0;
    pulses[0].gpioOn  = 0;
    pulses[0].gpioOff = 0;
    pulses[0].usDelay = 0;
    pulse_count++;
}

bool add_pulse(int on_mask, int off_mask, int us_delay, bool invert)
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

void add_delay(int us_delay)
{
    // This is easy. Just add delay to the current pulse.
    pulses[pulse_count-1].usDelay += us_delay;
}

bool send_pulses(int pin)
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

bool sendBytes(int pin, int period, int repeats, bool translate, bool invert,
               byte input_buffer[], int input_length)
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

        // Should we translate the input bytes to nibbles? Probably yes.
        if (translate == true) {
            int i = 0;
            for (i = 0; i < buffer_size; i++) {
                output_buffer[i] = nibbles[input_buffer[i] & 0x0F];
            }
        } else {
            int i = 0;
            for (i = 0; i < buffer_size; i++) {
                output_buffer[i] = input_buffer[i];
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

int main ()
{
    bool ok = initGpio(22);
    if (ok == false) {
        gpioTerminate();
        printf("terminating early");
        return 1;
    }

    byte on_cmd[]  = {0x9,0xf,0x3,0x1,0x5,0x9,0x3,0x0,0x1,0x2};
    byte off_cmd[] = {0x4,0x0,0x3,0x0,0x5,0x9,0x3,0x0,0x1,0x2};

    int i = 0;
    for (i = 0; i < 2; i++) {

        printf("%2d ON\n", i);
        //bool sendBytes(int pin, int period, int repeats, bool translate, bool invert,
        //               byte input_buffer[], int input_length)
        ok = sendBytes(22, 140, 10, true, false, on_cmd, 10);
        sleep(2);

        printf("%2d OFF\n", i);
        ok = sendBytes(22, 140, 10, true, false, off_cmd, 10);
        sleep(2);
    }

    gpioTerminate();
    printf("finished");
    return !ok;
}
