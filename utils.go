/*
Copyright © 2018-2019 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"errors"
	"log"
	"net"
	"os"
)

// PathExists returns whether the given file or directory exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//GetIPv4FromInterface gets an IPv4 address from the specific interface
func GetIPv4FromInterface(interfaceWtf string) (string, error) {
	ifaces, err := net.InterfaceByName(interfaceWtf)
	if err != nil {
		log.Println("Can't get network device "+interfaceWtf+".", err)
		return "", err
	}

	addrs, err := ifaces.Addrs()
	if err != nil {
		log.Println("Can't get ip address from "+interfaceWtf+".", err)
		return "", err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil {
			continue
		}

		if !(ip.IsGlobalUnicast() && !(ip.IsUnspecified() || ip.IsMulticast() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast())) {
			continue
		}

		//the code is not ready for updating an AAAA record
		/*
			if (isIPv4(ip.String())){
				if (configuration.IPType=="IPv6"){
					continue;
				}
			}else{
				if (configuration.IPType!="IPv6"){
					continue;
				}
			} */
		if ip.To4() == nil {
			continue
		}

		return ip.String(), nil

	}
	return "", errors.New("can't get a vaild address from " + interfaceWtf)
}
