# go-freakble

A simple tool to send messages into [FreakWAN](https://github.com/antirez/sx1276-micropython-driver/)
over Bluetooth low energy.

**This project is in pre-alpha please use the [python implementation](https://pypi.org/project/freakble/)
instead.**

## Installation

### From source

To build the latest version of `go-freakble` run:

```
cd /tmp; GOPATH=/tmp/go go install noa.mornie.org/eriol/go-freakble/cmd/freakble@latest
```

You will find the binary at `/tmp/go/bin/freakble`. You don't need to set
`GOPATH` if you already have it set and you are fine having freakble installed
there.

## Usage

```console
A simple tool to send messages into FreakWAN.

Usage:
  freakble [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  scan        Scan to find BLE devices.

Flags:
  -h, --help   help for freakble

Use "freakble [command] --help" for more information about a command.
```
