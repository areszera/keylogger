# Keylogger

> **Statement**: In consideration of all possible negative effects this software will bring, this software is developed
> only for academic usage. Any malicious usage of this software is forbidden and unrelated to the author.

Keylogger is consisted with an application and a server. The application listens events on keyboard and changes of
clipboard, then send them to the target server.

## Application

The keylogger application is programmed to record keyboard events, clipboard, and the current running process title.
When there is no operation in the recent 10 seconds, it will serialise to JSON and send to the keylogger server via TCP.
Every time the application runs, it checks, copies and registers autostart.

### Compile and Run

The application is suggested to be compiled to executable files then run. Make sure Go and GCC has been installed in
your device.

#### On Windows

Execute `go build -ldflags "-H windowsgui" -o keylogger.exe`. The `ldflags` of `-H windowsgui` hides the command line
window when the application is running. Then double-click the generated `keylogger.exe` file to run it. In the view of
users, it looks like nothing has happened, but in the Desktop Window Manager, a process named `keylogger.exe` will be in
the list.

#### On Linux

To compile the keylogger application, other support libraries are required to be installed.

- GCC for compiling.
- `xcb`, `xkb`, `libxkbcommon` for listening keyboard events.
- `xsel`, `xclip` for listening clipboards.

To install them using `apt`:

> ```bash
> # Update
> sudo apt update
> 
> # gcc
> sudo apt install gcc libc6-dev
> 
> # x11
> sudo apt install libx11-dev xorg-dev libxtst-dev
>
> # Hook
> sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
>
> # Clipboard
> sudo apt install xsel xcli
> ```

When all the necessary libraries are installed, execute `go build -o keylogger` to compile and `sudo ./keylogger` to run
the application. `sudo` is necessary here because the directories of `/usr/local/bin/` and `/usr/lib/systemd/system/`
requires root privileges to modify.

So far, the application **cannot** run silently in the background or as a daemon process. And the service **cannot** start (
result by core-dump).

### Recover

After running the Keylogger application, it will copy the file to the current user directory and set up autostart. Thus,
even though the compiled application has been deleted, it still exists in other places.

#### On Windows

1. Go to the current user directory, e.g., `C:\Users\CurrentUserName\`, show all the hidden files. Delete the hidden
   file named `keylogger.exe`.
2. Press `Windows` and `R` keys at the same time, then type `regedit` to open the Registry Editor.
3. Go to the `HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Run` directory.
4. Delete the key-value pair named `keylogger`.

#### On Linux

1. Delete `/usr/local/bin/.keylogger` by executing `sudo rm /usr/local/bin/.keylogger`. It deletes the binary copy in
   the system.
2. Delete `/usr/lib/systemd/system/keylogger.service` by executing `sudo rm /usr/lib/systemd/system/keylogger.service`.
   This step removes keylogger service from the system.

## Server

The server is programmed to receive key event logs via TCP. It initialises a TCP listener which binds port `8722`, then
continuously wait for the logs. Every time the server received logs, it will unmarshal the JSON data and try to update
log files. The server uses `.log` text file or SQLite database to record logs. In the current stage, the server will not
write the `KeyLog.Title` field data to texts.

### Compile and Run

Execute `go build` to compile the server. Using `go run server.go` to run is also acceptable. The following arguments
are available:

- `-db`: Use SQLite database to store logs.
- `-both`: Use both SQLite database and `.log` text file to store logs.
- No extra argument or other arguments will use `.log` text file to store logs.

## Platforms

- [x] Windows (passed test on Windows 10 and 11)
- [x] Linux (passed partial functions on Ubuntu 20.4)
- [ ] Darwin (have not tested)

## Acknowledgements

- [Robotgo](https://github.com/go-vgo/robotgo): Go Native cross-platform GUI automation.
- [clipboard](https://github.com/atotto/clipboard): Clipboard for Golang.
- [GoLand](https://www.jetbrains.com/go/): A Clever IDE to Go by JetBrains.
