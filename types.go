package lwgo

import "fmt"

type Command int
//go:generate stringer -type=Command

const (
    On Command = iota
    Off
    Dim
    Increase
    Decrease
    Mood
    AllOn
    AllOff
)

// LwBuffer contains the command you wish to send along with repeat
// configuration.
type LwBuffer [10]byte

// LwCommand is a helper struct to pull out the meaning of a LwBuffer, useful
// for logging.
type LwCommand struct {
    Command Command
    Value   int      // 0-31 dim levels (TODO: Moods)
    Device  int
    Address []byte
    Room    int
}

func NewBuffer(bytebuffer []byte) (LwBuffer, error) {
    if len(bytebuffer) <= 0 {
        return LwBuffer{}, fmt.Errorf("input buffer size is too small: %d", len(bytebuffer))
    } else if len(bytebuffer) > 10 {
        return LwBuffer{}, fmt.Errorf("input buffer size is too big: %d", len(bytebuffer))
    } else {
        buffer := LwBuffer{}
        for i, val := range bytebuffer {
            buffer[i] = val;
        }
        return buffer, nil
    }
}

// Convert the LwCommand to a LwBuffer.
func (cmd LwCommand) Buffer() LwBuffer {
    // parameter (2 [0,1])
    // device    (1 [2])
    // command   (1 [3])
    // address   (5 [4-8])
    // room      (1 [9])
    parameter, command := 0, 0
    device := cmd.Device
    address := cmd.Address
    room := cmd.Room

    // Get the parameter
    switch cmd.Command {
        // Command off
        case Off: {
            // Off
            parameter = 64; // 0-127, usually 64
            command = 0;
        }
        // case Dim: {
        //     // Off with level
        //     parameter = cmd.Value + 128; // 128-159
        //     command = 0;
        // }
        case Decrease: {
            // Decrease brightness
            parameter = 160; // 160-191
            command = 0;
        }
        case AllOff: {
            // All off
            parameter = 192; // 192-255, usually 192
            command = 0;
        }

        // Command on
        case On: {
            // On (to last level)
            parameter = 0; // 0-31
            command = 1;
        }
        case Dim: {
            // On with level
            parameter = cmd.Value + 64; // 32-63, 64-95, 96-127, 128-159
            command = 1;
        }
        case Increase: {
            // Increase brightness
            parameter = cmd.Value + 160; // 160-191
            command = 1;
        }
        case AllOn: {
            // All on with level
            parameter = cmd.Value + 192; // 192-223, 224-255
            command = 1;
        }

        // Command mood
        case Mood: {
            // Mood
            parameter = cmd.Value + 2; // 2-129, 130-255
            command = 2;
        }

        // Unknown case!
        default: {
            fmt.Println("lwgo::Message: cannot determine case for Command: ", cmd)
        }
    }

    // Build the message
    buffer := LwBuffer{
        ((byte(parameter) >> 4) & 0x0F), // 0 parameter[1/2]
        (byte(parameter) & 0x0F), // 1 parameter[2/2]
        byte(device),             // 2 device
        byte(command),            // 3 command
        byte(address[0]),         // 4 address[1/5]
        byte(address[1]),         // 5 address[2/5]
        byte(address[2]),         // 6 address[3/5]
        byte(address[3]),         // 7 address[4/5]
        byte(address[4]),         // 8 address[5/5]
        byte(room),               // 9 room
    }

    return buffer
}

// Convert the LwBuffer to a LwCommand.
func (buffer LwBuffer) Command() (LwCommand, error) {
    // parameter (2 [0,1])
    // device    (1 [2])
    // command   (1 [3])
    // address   (5 [4-8])
    // room      (1 [9])
    command := LwCommand{
        Device: int(buffer[2]),
        Address: buffer[4:9],
        Room: int(buffer[9]),
    }
    var err error = nil

    cmd := int(buffer[3])
    param := int(buffer[1])
    param += int(buffer[0] << 4)

    // Get the parameter
    switch {
        // Command off
        case cmd == 0 && param >= 0 && param <= 127: {
            // Off
            command.Command = Off
            command.Value = 0
        }
        case cmd == 0 && param >= 128 && param <= 159: {
            // Off with level
            command.Command = Dim
            command.Value = param - 128
        }
        case cmd == 0 && param >= 160 && param <= 191: {
            // Decrease brightness
            command.Command = Decrease
            command.Value = 160
        }
        case cmd == 0 && param >= 192 && param <= 255: {
            // All off
            command.Command = AllOff
            command.Value = 192
        }

        // Command on
        case cmd == 1 && param >= 0 && param <= 31: {
            // On (to last level)
            command.Command = On
            command.Value = 0
        }
        case cmd == 1 && param >= 32 && param <= 63: {
            // On with level
            command.Command = Dim
            command.Value = param - 32
        }
        case cmd == 1 && param >= 64 && param <= 95: {
            // On with level
            command.Command = Dim
            command.Value = param - 64
        }
        case cmd == 1 && param >= 96 && param <= 127: {
            // On with level
            command.Command = Dim
            command.Value = param - 96
        }
        case cmd == 1 && param >= 128 && param <= 159: {
            // On with level
            command.Command = Dim
            command.Value = param - 128
        }
        case cmd == 1 && param >= 160 && param <= 191: {
            // Increase brightness
            command.Command = Increase
            command.Value = 160
        }
        case cmd == 1 && param >= 192 && param <= 223: {
            // All on with level
            command.Command = AllOn
            command.Value = param - 192
        }
        case cmd == 1 && param >= 224 && param <= 255: {
            // All on with level
            command.Command = AllOn
            command.Value = param - 224
        }

        // Command mood
        case cmd == 2 && param >= 130 && param <= 255: {
            // Mood
            command.Command = Mood
            command.Value = param - 192
        }
        case cmd == 2 && param >= 2 && param <= 129: {
            // Mood
            command.Command = Mood
            command.Value = param - 1
        }

        default: {
            err = fmt.Errorf("could not convert the buffers command and parameter values to a valid Command state")
        }
    }

    return command, err
}

// String gives you a human friendly version of the LwBuffer.
func (buffer LwBuffer) String() (string, error) {
    cmd, err := buffer.Command()
    return cmd.String(), err
}

// String gives you a human friendly version of the LwCommand.
func (command LwCommand) String() string {
    return fmt.Sprint("Command: ", command.Command,
                      ", Value: ", command.Value,
                      ", Device: ", command.Device,
                      ", Address: ", command.Address,
                      ", Room: ", command.Room)
}

// Raw returns the raw byte buffer stored within the LwBuffer.
func (buffer LwBuffer) Raw() []byte {
    out := make([]byte, len(buffer))
    for i, val := range buffer {
        out[i] = val
    }
    return out
}
