# Keylogger

> **Statement**: In consideration of all possible negative effects this software will bring, this software is developed
> only for academic usage. Any malicious usage of this software is forbidden and unrelated to the author.

Keylogger is consisted with an application and a server. The application listens events on keyboard and changes of
clipboard, then send them to the target server.

## Application

The keylogger application is programmed to record keyboard events and clipboard. When there is no operation in the
recent 10 seconds, it will serialise to JSON and send to the keylogger server via TCP. Every time the application runs,
it checks and registers autostart.

### Compile and Run

The application is suggested to be compiled to executable files then run. Make sure Go has been installed in your
device, to build for Windows, execute `go build -ldflags "-H windowsgui"`. The `ldflags` of `-H windowsgui` hides the
command line window when the application is running (Windows only).

### Recover

After running the Keylogger application, it will set up autostart. To revoke autostart on Windows systems:

1. Press `Windows` and `R` keys at the same time, then type `regedit` to open the Registry Editor.
2. Go to the `HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Run` directory.
3. Delete the key-value pair whose name is `Keylogger`.

## Server

The server is programmed to receive key event logs via TCP. It initialises a TCP listener which binds port `8722`, then
continuously wait for the logs. Every time the server received logs, it will unmarshal the JSON data and try to append
to the `keylogger.log` file.

### Compile and Run

Execute `go build` to compile the server. Using `go run keylogger.go` to run is also acceptable.

## Platforms

- [x] Windows (passed test on Windows 10 and 11)
- [ ] Linux (have not tested)
- [ ] Darwin (have not tested)

## Acknowledgements

- [Robotgo](https://github.com/go-vgo/robotgo): Go Native cross-platform GUI automation.
- [clipboard](https://github.com/atotto/clipboard): Clipboard for Golang.
- [GoLand](https://www.jetbrains.com/go/): A Clever IDE to Go by JetBrains.
