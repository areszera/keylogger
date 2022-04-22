package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	fmt.Println("Start removing malware")
	result, err := exec.Command("tasklist").CombinedOutput()
	if err != nil {
		fmt.Println("Failed to get task list:", err.Error())
		return
	}
	temp := regexp.MustCompile(`sync\.exe[\s\S]*?Console`).FindString(string(result))
	if temp != "" {
		pid := regexp.MustCompile(`\d+`).FindString(temp)
		err = exec.Command("taskkill", "/pid", pid, "-t", "-f").Run()
		if err != nil {
			fmt.Println("Failed to kill task:", err.Error())
			return
		}
	}
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("Failed to open key:", err.Error())
		return
	}
	defer key.Close()
	path, _, err := key.GetStringValue("sync")
	if err != nil {
		fmt.Println("Malware has been moved")
		return
	}
	filename := regexp.MustCompile(`(^")|(".*?$)`).ReplaceAllString(path, "")
	_ = key.DeleteValue("sync")
	_ = os.Remove(filename)
	fmt.Println("Malware has been moved successfully")
}
