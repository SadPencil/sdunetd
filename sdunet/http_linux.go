//go:build linux
// +build linux

/*
Copyright © 2018-2022 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package sdunet

import (
	"golang.org/x/sys/unix"
	"net"
	"net/http"
	"syscall"
	"time"
)

func getHttpClient(forceNetworkInterface string, timeout time.Duration) (client *http.Client, err error) {
	client = &http.Client{Timeout: timeout}
	if forceNetworkInterface != "" {
		// https://iximiuz.com/en/posts/go-net-http-setsockopt-example/
		// https://linux.die.net/man/7/socket
		dialer := &net.Dialer{
			Control: func(network, address string, conn syscall.RawConn) error {
				var operr error
				if err := conn.Control(func(fd uintptr) {
					operr = unix.BindToDevice(int(fd), forceNetworkInterface)
				}); err != nil {
					return err
				}
				return operr
			},
		}

		client.Transport = &http.Transport{
			DialContext: dialer.DialContext,
		}
	}
	return client, nil
}
