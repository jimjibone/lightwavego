// Don't build this file on Raspberry Pi (not linux && not arm)
// +build !linux,!arm

package lwgo

import "fmt"

func init() {
    fmt.Println("NOGPIO - lwgo:interface::init")
}

func Shutdown() {
    fmt.Println("NOGPIO - lwgo:interface::Shutdown")
}

func (lw *LwTx) sendBuffer(message LwBuffer) error {
    fmt.Println("NOGPIO - lwgo:interface::sendBuffer:", message)
    return nil
}
