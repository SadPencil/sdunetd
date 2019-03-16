/*
Copyright © 2018-2019 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func version() {
	fmt.Println(NAME, VERSION)
	fmt.Println(DESCRIPTION)

}

func main() {
	var FlagShowHelp bool
	flag.BoolVar(&FlagShowHelp, "h", false, "standalone: show the help.")
	var FlagShowVersion bool
	flag.BoolVar(&FlagShowVersion, "v", false, "standalone: show the version")

	var FlagConfigFile string
	flag.StringVar(&FlagConfigFile, "c", "", "the config.json file. Leave it blank to use the interact mode.")

	var FlagNoAttribute bool
	flag.BoolVar(&FlagNoAttribute, "m", false, "option: output log without attributes. Turn it on when running as a systemd service.")

	var FlagOneshoot bool
	flag.BoolVar(&FlagOneshoot, "f", false, "standalone: login to the network, regardless of whether Internet is available.")

	var FlagIPDetect bool
	flag.BoolVar(&FlagIPDetect, "a", false, "standalone: detect the IP address from the authenticate server. Useful when behind a NAT router.")

	flag.Parse()

	version()
	if FlagShowVersion {
		return
	}
	if FlagShowHelp {
		flag.Usage()
		return
	}

	fileExist, err := PathExists(FlagConfigFile)
	if err != nil {
		panic(err)
	}

	if !fileExist {
		cartman()
		return
	}

	Settings, err := LoadSettings(FlagConfigFile)
	if err != nil {
		panic(err)
	}
	//required arguments
	checkInterval(&Settings)
	err = checkAuthServer(&Settings)
	err = checkPassword(&Settings)
	err = checkScheme(&Settings)
	err = checkUsername(&Settings)
	if err != nil {
		panic(err)
	}

	if Settings.Network.CustomIP == "" && Settings.Network.Interface == "" {
		panic("You need to specify either the local IP or the interface name.")
	}

	if FlagIPDetect {
		ret, err := getIPFromChallenge(Settings.Account.Scheme,
			Settings.Account.AuthServer,
			Settings.Account.Username,
			Settings.Network.CustomIP,
			Settings.Network.Interface,
			Settings.Control.StrictMode)
		if err != nil {
			panic(err)
		}
		fmt.Println(ret)
		return
	}

	if FlagOneshoot {
		log.Println("Log in via web portal...")
		err := login(Settings.Account.Scheme,
			Settings.Account.AuthServer,
			Settings.Account.Username,
			Settings.Account.Password,
			Settings.Network.CustomIP,
			Settings.Network.Interface,
			Settings.Control.StrictMode,
			Settings.Network.DetectIP)
		if err != nil {
			log.Println("Login failed.", err)
		}
		return
	}

	//write log to stdout
	log.SetOutput(os.Stdout)

	//openfile
	if Settings.Log.Filename != "" {
		logFile, err := os.OpenFile(Settings.Log.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		defer logFile.Close()
		if err != nil {
			log.Panicln(err)
			return
		}
		log.SetOutput(logFile)
	}

	log.SetFlags(0)

	//loop
	for {
		if !detectNetwork() {
			log.Println("Network is down. Log in via web portal...")
			err := login(Settings.Account.Scheme,
				Settings.Account.AuthServer,
				Settings.Account.Username,
				Settings.Account.Password,
				Settings.Network.CustomIP,
				Settings.Network.Interface,
				Settings.Control.StrictMode,
				Settings.Network.DetectIP)
			if err != nil {
				log.Println("Login failed.", err)
			}
			time.Sleep(time.Duration(5) * time.Second)
			if !detectNetwork() {
				log.Println("Network is still down. Retry after", Settings.Control.Interval, " seconds.")
			} else {
				log.Println("Network is up.")
			}
		} else {
			log.Println("Network is up. Nothing to do.")
		}
		time.Sleep(time.Duration(Settings.Control.Interval) * time.Second)
	}

}
