// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for educational usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

package main

import (
	"fmt"
	"net"
)

func main() {
	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8722")
	if err != nil {
		fmt.Printf("Failed to resolve UDP address: %s\n", err.Error())
		return
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Printf("Failed to listen UDP: %s\n", err.Error())
		return
	}
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, raddr, e := conn.ReadFromUDP(buffer)
		if e != nil {
			fmt.Printf("Failed to read from UDP: %s\n", e.Error())
			return
		}
		fmt.Printf("Sent from %s: %s\n", raddr, buffer[:n])
	}
}
