# Keylogger

> **Statement**: In consideration of all possible negative effects this software will bring, this software is developed
> only for academic usage. Any malicious usage of this software is forbidden and unrelated to the author.

Keylogger is consisted with an application and a server. The application listens events on keyboard and changes of
clipboard, then send them to the target server.

## Usage

Go to the `keylogger-application` and `keylogger-server` directory and execute `go mod tidy`to download libraries, then
run `go run keylogger.go` to start the application and server, respectively.

## Acknowledgements

- [Robotgo](https://github.com/go-vgo/robotgo): Go Native cross-platform GUI automation.
- [clipboard](https://github.com/atotto/clipboard): Clipboard for Golang.
- [GoLand](https://www.jetbrains.com/go/): A Clever IDE to Go by JetBrains.
