// Example `curl` commmand to use to turn on a light to maximum brightness:
// `curl -i -d '{"Buffer":"090f0301050903000102"}' http://localhost:8080/send`

package main

import (
    "fmt"
    "log"
    "sync"
    "time"
    "net/http"
    "encoding/hex"
    "github.com/jimjibone/lightwavego"
    "github.com/ant0ine/go-json-rest/rest"
)

var lwtx *lightwavego.LwTx

func main() {
    fmt.Println("lightwavego server starting on port 8080")
    lwtx = lightwavego.NewLwTx()

    handler := rest.ResourceHandler{
        EnableRelaxedContentType: true,
    }
    err := handler.SetRoutes(
        &rest.Route{"POST", "/send", PostSend},
        &rest.Route{"GET", "/send/history", GetSendHistory},
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Fatal(http.ListenAndServe(":8080", &handler))
}

type Command struct {
    Sequence int
    Time     time.Time
    Buffer   string
    Bytes    []byte
    Command  lightwavego.LwCommand
}

var lock = sync.RWMutex{}
var sent = []Command{}
var lastSequence = 0

func PostSend(w rest.ResponseWriter, r *rest.Request) {
    lastSequence++
    command := Command{Time: time.Now(), Sequence: lastSequence}

    err := r.DecodeJsonPayload(&command)
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if command.Buffer == "" {
        rest.Error(w, "command buffer required", 400)
        return
    }

    command.Bytes, err = hex.DecodeString(command.Buffer)
    if err != nil {
        errstring := fmt.Sprint("command message could not be converted to bytes:", err)
        rest.Error(w, errstring, 400)
        return
    }

    // Convert the raw bytes to LwBuffer format.
    // TODO: handle error
    buffer, converr := lightwavego.NewBuffer(command.Bytes)
    if converr != nil {
        errstring := fmt.Sprint("could not convert input buffer to LwBuffer:", converr)
        rest.Error(w, errstring, 400)
        return
    } else {
        command.Command, _ = buffer.Command()
    }

    log.Print("Sending: seq: ", command.Sequence, ", time: ", command.Time,
              ", buffer: ", buffer, ", command: ", command.Command)

    lock.Lock()
    sent = append(sent, command)
    lwtx.SendBuffer(buffer)
    lock.Unlock()
    w.WriteJson(&command)
}

func GetSendHistory(w rest.ResponseWriter, r *rest.Request) {
    lock.RLock()
    commands := make([]Command, len(sent))
    i := 0
    for _, command := range sent {
        commands[i] = command
        i++
    }
    lock.RUnlock()
    w.WriteJson(&commands)
}
