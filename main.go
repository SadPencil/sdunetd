/*
Copyright © 2018-2019 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func version() {
	fmt.Println(NAME, VERSION)
	fmt.Println(DESCRIPTION)
}

func _getNetworkInterface(settings *Settings) (interfaceWtf string) {
	if settings.Network.StrictMode {
		interfaceWtf = settings.Network.Interface
	}
	return interfaceWtf
}

func _login(settings *Settings) (err error) {
	interfaceWtf := _getNetworkInterface(settings)

	_, sduIPv4, err := getSduUserInfo(settings.Account.Scheme, settings.Account.AuthServer, interfaceWtf)
	if err != nil {
		return err
	}

	return loginDigest(settings.Account.Scheme, settings.Account.AuthServer, settings.Account.Username, settings.Account.Password, sduIPv4, interfaceWtf)
}

func _logout(settings *Settings) (err error) {
	interfaceWtf := _getNetworkInterface(settings)
	return logout(settings.Account.Scheme, settings.Account.AuthServer, settings.Account.Username, interfaceWtf)
}

func detectNetwork(settings *Settings) bool {
	if settings.Control.OnlineDetectionMethod == ONLINE_DETECTION_METHOD_AUTH {
		return detectNetworkWithAuthServer(settings)
	} else if settings.Control.OnlineDetectionMethod == ONLINE_DETECTION_METHOD_MS {
		return detectNetworkWithMicrosoft()
	} else {
		return detectNetworkWithAuthServer(settings)
	}
}

func detectNetworkWithMicrosoft() bool {
	resp, err := http.Get("http://www.msftconnecttest.com/connecttest.txt")
	if err != nil {
		log.Println("Error: " + err.Error())
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error: " + err.Error())
		return false
	}

	return bytes.Compare(body, []byte("Microsoft Connect Test")) == 0
}

func detectNetworkWithAuthServer(settings *Settings) bool {
	//client, err := getHttpClient(settings.Network.StrictMode, settings.Network.CustomIP, settings.Network.Interface)
	//if err != nil {
	//	log.Println("[ERROR]", err)
	//	return false
	//}
	interfaceWtf := _getNetworkInterface(settings)
	ret, ipv4, err := getSduUserInfo(settings.Account.Scheme, settings.Account.AuthServer, interfaceWtf)
	if err != nil {
		log.Println("[ERROR]", err)
		return false
	} else {
		log.Println("Logined to sdunet. IP Address:", ipv4)
		return ret
	}
}

func main() {
	var FlagShowHelp bool
	flag.BoolVar(&FlagShowHelp, "h", false, "standalone: show the help.")

	var FlagShowVersion bool
	flag.BoolVar(&FlagShowVersion, "v", false, "standalone: show the version.")

	var FlagConfigFile string
	flag.StringVar(&FlagConfigFile, "c", "", "the config.json file. Leave it blank to use the interact mode.")

	var FlagNoAttribute bool
	flag.BoolVar(&FlagNoAttribute, "m", false, "option: output log without attributes. Turn it on when running as a systemd service.")

	var FlagIPDetect bool
	flag.BoolVar(&FlagIPDetect, "a", false, "standalone: detect the IP address from the authenticate server. Useful when behind a NAT router.")

	var FlagOneshoot bool
	flag.BoolVar(&FlagOneshoot, "f", false, "standalone: login to the network for once, regardless of whether Internet is available.")

	var FlagTryOneshoot bool
	flag.BoolVar(&FlagTryOneshoot, "t", false, "standalone: login to the network for once, only if Internet isn't available.")

	var FlagLogout bool
	flag.BoolVar(&FlagLogout, "l", false, "standalone: logout from the network for once.")

	flag.Parse()

	if FlagShowVersion {
		version()
		return
	}
	if FlagShowHelp {
		version()
		flag.Usage()
		return
	}

	fileExist, err := PathExists(FlagConfigFile)
	if err != nil {
		panic(err)
	}

	if !fileExist {
		version()
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
	if FlagNoAttribute {
		log.SetFlags(0)
	}

	if FlagIPDetect {
		//client, err := getHttpClient(Settings.Network.StrictMode, Settings.Network.CustomIP, Settings.Network.Interface)
		interfaceWtf := _getNetworkInterface(&Settings)
		_, ret, err := getSduUserInfo(Settings.Account.Scheme,
			Settings.Account.AuthServer, interfaceWtf)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(ret)
		return
	}

	if FlagOneshoot {
		version()
		log.Println("Log in via web portal...")
		err := _login(&Settings)
		if err != nil {
			log.Println("Login failed.", err)
			os.Exit(1)
		} else {
			log.Println("Succeed.")
		}
		return
	}

	if FlagTryOneshoot {
		version()
		if !detectNetwork(&Settings) {
			log.Println("Network is down. Log in via web portal...")
			err := _login(&Settings)
			if err != nil {
				log.Println("Login failed.", err)
				os.Exit(1)
			} else {
				log.Println("Succeed.")
			}
		} else {
			log.Println("Network is up. Nothing to do.")
		}
		return
	}

	if FlagLogout {
		version()
		log.Println("Logout via web portal...")
		err := _logout(&Settings)
		if err != nil {
			log.Println("Logout failed.", err)
			os.Exit(1)
		} else {
			log.Println("Succeed.")
		}
		return
	}

	version()

	// main loop

	// set up handler for SIGINT and SIGTERM
	pauseNow := false
	exitNow := false
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Exiting...")
				if Settings.Control.LogoutWhenExit {
					pauseNow = true
					for retry := 0; retry < 3; retry++ {
						fmt.Println("Logging out...")
						err := _logout(&Settings)
						if err != nil {
							log.Println("Logout failed.", err)
						} else {
							break
						}
					}
					pauseNow = false
				}
				exitNow = true
			default:
			}
		}
	}()
	for {
		for pauseNow && !exitNow {
			time.Sleep(time.Second)
		}
		if exitNow {
			break
		}
		if !detectNetwork(&Settings) {
			log.Println("Network is down. Log in via web portal...")
			err := _login(&Settings)
			if err != nil {
				log.Println("Login failed.", err)
			}
			time.Sleep(time.Duration(5) * time.Second)
			if !detectNetwork(&Settings) {
				log.Println("Network is still down. Retry after", Settings.Control.Interval, " seconds.")
			} else {
				log.Println("Network is up.")
			}
		} else {
			log.Println("Network is up. Nothing to do.")
		}
		for i := int32(0); i < Settings.Control.Interval; i++ {
			if exitNow {
				break
			}
			time.Sleep(time.Second)
		}

	}

}
