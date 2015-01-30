// Don't build this file on Raspberry Pi (i.e. not arm)
// +build !linux,!arm

package lightwavego

import "fmt"

func init() {
    fmt.Println("NOGPIO - lightwavego:interface::init")
}

func Shutdown() {
    fmt.Println("NOGPIO - lightwavego:interface::Shutdown")
}

func (lw *LwTx) sendBuffer(message LwBuffer) error {
    fmt.Println("NOGPIO - lightwavego:interface::sendBuffer:", message)
    return nil
}
