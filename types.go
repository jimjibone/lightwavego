package lwgo

import (
    "fmt"
)

// LwMessage contains the command you wish to send along with repeat configuration.
type LwMessage [10]byte

// LwCommand is a helper struct to pull out the meaning of a LwMessage, useful
// for logging.
type LwCommand struct {
    Command string // On, Off, Dim, Increase, Decrease, Mood, All On, All Off
    Value int      // 0-31 dim levels (TODO: Moods)
    Device int
    Address []byte
    Room int
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

// Convert the LwCommand to a LwMessage.
func (cmd LwCommand) Message() LwMessage {
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
    switch {
        // Command off
        case cmd.Command == "Off": {
            // Off
            parameter = 64; // 0-127, usually 64
            command = 0;
        }
        // case cmd.Command == "Dim": {
        //     // Off with level
        //     parameter = cmd.Value + 128; // 128-159
        //     command = 0;
        // }
        case cmd.Command == "Decrease": {
            // Decrease brightness
            parameter = 160; // 160-191
            command = 0;
        }
        case cmd.Command == "All Off": {
            // All off
            parameter = 192; // 192-255, usually 192
            command = 0;
        }

        // Command on
        case cmd.Command == "On": {
            // On (to last level)
            parameter = 0; // 0-31
            command = 1;
        }
        case cmd.Command == "Dim": {
            // On with level
            parameter = cmd.Value + 64; // 32-63, 64-95, 96-127, 128-159
            command = 1;
        }
        case cmd.Command == "Increase": {
            // Increase brightness
            parameter = cmd.Value + 160; // 160-191
            command = 1;
        }
        case cmd.Command == "All On": {
            // All on with level
            parameter = cmd.Value + 192; // 192-223, 224-255
            command = 1;
        }

        // Command mood
        case cmd.Command == "Mood": {
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
    message := LwMessage{
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

    return message
}

// Convert the LwMessage to a LwCommand.
func (message LwMessage) Command() LwCommand {
    // parameter (2 [0,1])
    // device    (1 [2])
    // command   (1 [3])
    // address   (5 [4-8])
    // room      (1 [9])
    command := LwCommand{
        Device: int(message[2]),
        Address: message[4:9],
        Room: int(message[9]),
    }

    cmd := int(message[3])
    param := int(message[1])
    param += int(message[0] << 4)

    // Get the parameter
    switch {
        // Command off
        case cmd == 0 && param >= 0 && param <= 127: {
            // Off
            command.Command = "Off"
            command.Value = 0
        }
        case cmd == 0 && param >= 128 && param <= 159: {
            // Off with level
            command.Command = "Dim"
            command.Value = param - 128
        }
        case cmd == 0 && param >= 160 && param <= 191: {
            // Decrease brightness
            command.Command = "Decrease"
            command.Value = 160
        }
        case cmd == 0 && param >= 192 && param <= 255: {
            // All off
            command.Command = "All Off"
            command.Value = 192
        }

        // Command on
        case cmd == 1 && param >= 0 && param <= 31: {
            // On (to last level)
            command.Command = "On"
            command.Value = 0
        }
        case cmd == 1 && param >= 32 && param <= 63: {
            // On with level
            command.Command = "Dim"
            command.Value = param - 32
        }
        case cmd == 1 && param >= 64 && param <= 95: {
            // On with level
            command.Command = "Dim"
            command.Value = param - 64
        }
        case cmd == 1 && param >= 96 && param <= 127: {
            // On with level
            command.Command = "Dim"
            command.Value = param - 96
        }
        case cmd == 1 && param >= 128 && param <= 159: {
            // On with level
            command.Command = "Dim"
            command.Value = param - 128
        }
        case cmd == 1 && param >= 160 && param <= 191: {
            // Increase brightness
            command.Command = "Increase"
            command.Value = 160
        }
        case cmd == 1 && param >= 192 && param <= 223: {
            // All on with level
            command.Command = "All On"
            command.Value = param - 192
        }
        case cmd == 1 && param >= 224 && param <= 255: {
            // All on with level
            command.Command = "All On"
            command.Value = param - 224
        }

        // Command mood
        case cmd == 2 && param >= 130 && param <= 255: {
            // Mood
            command.Command = "Mood"
            command.Value = param - 192
        }
        case cmd == 2 && param >= 2 && param <= 129: {
            // Mood
            command.Command = "Mood"
            command.Value = param - 1
        }

        default: command.Command = "unknown"
    }

    return command
}

// Convert the LwMessage to a String.
func (message LwMessage) String() string {
    // parameter (2 [0,1])
    // device    (1 [2])
    // command   (1 [3])
    // address   (5 [4-8])
    // room      (1 [9])
    var output string

    command := int(message[3])
    param := int(message[1])
    param += int(message[0] << 4)

    // Get the parameter
    output += "Parameter: "
    switch {
        // Command off
        case command == 0 && param >= 0 && param <= 127: {
            output += "off"
        }
        case command == 0 && param >= 128 && param <= 159: {
            output += fmt.Sprint("off with level ", param-128)
        }
        case command == 0 && param >= 160 && param <= 191: {
            output += "decrease brightness"
        }
        case command == 0 && param >= 192 && param <= 255: {
            output += "all off"
        }

        // Command on
        case command == 1 && param >= 0 && param <= 31: {
            output += "on to last level"
        }
        case command == 1 && param >= 32 && param <= 63: {
            output += fmt.Sprint("on with level ", param-32)
        }
        case command == 1 && param >= 64 && param <= 95: {
            output += fmt.Sprint("on with level ", param-64)
        }
        case command == 1 && param >= 96 && param <= 127: {
            output += fmt.Sprint("on with level ", param-96)
        }
        case command == 1 && param >= 128 && param <= 159: {
            output += fmt.Sprint("on with level ", param-128)
        }
        case command == 1 && param >= 160 && param <= 191: {
            output += "increase brightness"
        }
        case command == 1 && param >= 192 && param <= 223: {
            output += fmt.Sprint("set all to level ", param-192)
        }
        case command == 1 && param >= 224 && param <= 255: {
            output += fmt.Sprint("set all to level ", param-224)
        }

        // Command mood
        case command == 2 && param >= 130 && param <= 255: {
            output += fmt.Sprint("start mood ", param-129)
        }
        case command == 2 && param >= 2 && param <= 129: {
            output += fmt.Sprint("define mood ", param-1)
        }

        default: output += "unknown"
    }
    output += ", "

    // Get the command
    output += "Command: "
    switch command {
        case 0: output += "off"
        case 1: output += "on"
        case 2: output += "mood"
        default: output += "unknown"
    }
    output += ", "

    // Get the other values
    output += fmt.Sprint("Device: ", int(message[2]))
    output += fmt.Sprint(", Address: ", message[4:9])
    output += fmt.Sprint(", Room: ", int(message[9]))

    return output
}

// Get a string version of the LwCommand.
func (command LwCommand) String() string {
    return fmt.Sprint("Command: ", command.Command,
                      ", Value: ", command.Value,
                      ", Device: ", command.Device,
                      ", Address: ", command.Address,
                      ", Room: ", command.Room)
}

// Raw returns the raw byte buffer stored within the LwMessage.
func (message LwMessage) Raw() []byte {
    out := make([]byte, len(message))
    for i, val := range message {
        out[i] = val
    }
    return out
}