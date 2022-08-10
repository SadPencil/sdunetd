/*
Copyright © 2018-2022 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/SadPencil/sdunetd/sdunet"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func version() {
	fmt.Println(NAME, VERSION)
	fmt.Println(DESCRIPTION)
}

func getNetworkInterface(settings *Settings) (networkInterface string) {
	if settings.Network.StrictMode {
		return settings.Network.Interface
	}
	return ""
}

func detectNetwork(settings *Settings, manager *sdunet.Manager) bool {
	if settings.Control.OnlineDetectionMethod == ONLINE_DETECTION_METHOD_AUTH {
		return detectNetworkWithAuthServer(manager)
	} else if settings.Control.OnlineDetectionMethod == ONLINE_DETECTION_METHOD_MS {
		return detectNetworkWithMicrosoft(manager)
	} else {
		return detectNetworkWithAuthServer(manager)
	}
}

func detectNetworkWithMicrosoft(manager *sdunet.Manager) bool {
	client, err := manager.GetNewHttpClient()
	resp, err := client.Get("http://www.msftconnecttest.com/connecttest.txt")
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

func detectNetworkWithAuthServer(manager *sdunet.Manager) bool {
	info, err := manager.GetUserInfo()
	if err != nil {
		log.Println(err)
		return false
	} else {
		log.Println("Logged in to sdunet. IP Address:", info.ClientIP)
		return true
	}
}

type Action func() error

func retry(retryTimes int32, retryInterval int32, action Action, exitNow *bool) error {
	for i := int32(0); i < retryTimes; i++ {
		err := action()
		if err == nil {
			return nil
		}
		if i == retryTimes-1 {
			return err
		}

		if exitNow != nil && *exitNow {
			return errors.New("exit triggered")
		}

		for i := int32(0); i < retryInterval; i++ {
			time.Sleep(time.Second)

			if exitNow != nil && *exitNow {
				return errors.New("exit triggered")
			}
		}
	}
	// impossible path
	// make compiler happy
	return errors.New("assert failed")
}

var _manager *sdunet.Manager

func getManager(settings *Settings) (*sdunet.Manager, error) {
	if _manager == nil {
		manager, err := sdunet.GetManager(
			settings.Account.Scheme,
			settings.Account.AuthServer,
			settings.Account.Username,
			getNetworkInterface(settings),
		)
		if err != nil {
			return nil, err
		}
		_manager = &manager
	}
	return _manager, nil
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

	settings, err := LoadSettings(FlagConfigFile)
	if err != nil {
		panic(err)
	}
	//required arguments
	checkInterval(&settings)
	if err != nil {
		panic(err)
	}
	err = checkAuthServer(&settings)
	if err != nil {
		panic(err)
	}
	err = checkPassword(&settings)
	if err != nil {
		panic(err)
	}
	err = checkScheme(&settings)
	if err != nil {
		panic(err)
	}
	err = checkUsername(&settings)
	if err != nil {
		panic(err)
	}

	//write log to stdout
	log.SetOutput(os.Stdout)

	//open the log file for writing
	if settings.Log.Filename != "" {
		logFile, err := os.OpenFile(settings.Log.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
		err := retry(settings.Control.MaxRetryCount, RETRY_INTERVAL, func() error {
			manager, err := getManager(&settings)
			if err != nil {
				log.Println(err)
				return err
			}
			info, err := manager.GetUserInfo()
			if err != nil {
				log.Println(err)
				return err
			}
			fmt.Println(info.ClientIP)
			return nil
		}, nil)
		if err != nil {
			log.Panicln(err)
		}
		return
	}

	if FlagOneshoot {
		version()
		log.Println("Log in via web portal...")
		manager, err := getManager(&settings)
		if err != nil {
			log.Println("Login failed.", err)
			os.Exit(1)
		}
		err = manager.Login(settings.Account.Password)
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
		manager, err := getManager(&settings)
		if err != nil {
			log.Println("Login failed.", err)
			os.Exit(1)
		}
		if !detectNetwork(&settings, manager) {
			log.Println("Network is down. Log in via web portal...")
			err := manager.Login(settings.Account.Password)
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
		err := retry(settings.Control.MaxRetryCount, RETRY_INTERVAL, func() error {
			log.Println("Logout via web portal...")
			manager, err := getManager(&settings)
			if err != nil {
				log.Println("Login failed.", err)
				return err
			}
			err = manager.Logout()
			if err != nil {
				log.Println("Logout failed.", err)
				return err
			}
			log.Println("Succeed.")
			return nil

		}, nil)
		if err != nil {
			os.Exit(1)
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
				log.Println("Exiting...")
				if settings.Control.LogoutWhenExit {
					pauseNow = true
					_ = retry(settings.Control.MaxRetryCount, RETRY_INTERVAL, func() error {
						log.Println("Logging out...")
						manager, err := getManager(&settings)
						if err != nil {
							log.Println("Logout failed.", err)
							return err
						}
						err = manager.Logout()
						if err != nil {
							log.Println("Logout failed.", err)
							return err
						}
						return nil
					}, nil)
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
		retry(settings.Control.MaxRetryCount, RETRY_INTERVAL, func() error {
			manager, err := getManager(&settings)
			if err != nil {
				log.Println(err)
				return err
			}
			if detectNetwork(&settings, manager) {
				log.Println("Network is up. Nothing to do.")
				return nil
			}
			log.Println("Network is down. Log in via web portal...")
			err = manager.Login(settings.Account.Password)
			if err != nil {
				log.Println("Login failed.", err)
			}
			return err
		}, &exitNow)
		for i := int32(0); i < settings.Control.Interval; i++ {
			if exitNow {
				break
			}
			time.Sleep(time.Second)
		}

	}

}
