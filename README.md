# Go OpenVPN Management Protocol Client

This repository contains a Go package that implements a client for
[the OpenVPN management interface](https://openvpn.net/index.php/open-source/documentation/miscellaneous/79-management-interface.html).

This can be used to monitor and control an OpenVPN process running with
its management port enabled.

Currently only a subset of the protocol is supported, primarily focused on
monitoring status changes and cleanly shutting down or restarting connections.
Patches that add additional functionality in a manner consistent with the
existing features will be gratefully accepted.

## Installation

```
go get github.com/apparentlymart/go-openvpn-mgmt/openvpn
```

## Usage

First, we can import the package:

```go
import (
    "github.com/apparentlymart/go-openvpn-mgmt/openvpn"
)
```

The default imported package name is `openvpn`.

OpenVPN's management interface can be used in two different modes: either
OpenVPN opens a TCP listen port and the management client connects to it,
or the management client opens a TCP listen port and *OpenVPN* connects
to it.

The default configuration is for the OpenVPN process to act as the TCP
server. Launch OpenVPN with the argument ``--management <listen-addr> <port>``
to activate this mode, and then connect to it by passing the same address
and port to `openvpn.Dial`. For example:

```go
eventCh := make(chan openvpn.Event, 10)
client, err := openvpn.Dial("127.0.0.1:6061", eventCh);
if err != nil {
    panic(err)
}

// "client" is now connected to the OpenVPN process and can send
// it commands. "eventCh" can be read to get asynchronous events from
// OpenVPN.
```

The alternative configuration, with OpenVPN acting as the TCP client, can be
useful for processes that manage a number of separate OpenVPN instances running
as child processes, since they can all be commanded to connect to the same
listen port on the parent management process. Launch OpenVPN with the arguments
``--management-client --management <remote-addr> <port>`` to activate this
mode, after creating a listen server using `openvpn.ListenAndServe`.
For example:

```go
func main() {
    log.Fatal(openvpn.ListenAndServe(
        "127.0.0.1:6061",
        openvpn.IncomingConnHandlerFunc(newConnection),
    ))
}

func newConnection(conn openvpn.IncomingConn) {
    eventCh := make(chan openvpn.Event, 10)
    client := conn.Open(eventCh)

    // "client" is now connected to the OpenVPN process and can send
    // it commands. "eventCh" can be read to get asynchronous events from
    // OpenVPN.
}
```

For more usage information on both modes, see
[the reference documentation](https://godoc.org/github.com/apparentlymart/go-openvpn-mgmt/openvpn).

# License

Copyright (c) 2016 Martin Atkins

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
