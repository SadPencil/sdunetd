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
	"github.com/SadPencil/sdunetd/setting"
	"github.com/SadPencil/sdunetd/utils"
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

func detectNetwork(settings *setting.Settings, manager *sdunet.Manager) (bool, error) {
	if settings.Control.OnlineDetectionMethod == setting.ONLINE_DETECTION_METHOD_AUTH {
		return detectNetworkWithAuthServer(manager)
	} else if settings.Control.OnlineDetectionMethod == setting.ONLINE_DETECTION_METHOD_MS {
		return detectNetworkWithMicrosoft(manager)
	} else {
		return detectNetworkWithAuthServer(manager)
	}
}

func detectNetworkWithMicrosoft(manager *sdunet.Manager) (bool, error) {
	client, err := manager.GetHttpClient()
	if err != nil {
		return false, err
	}

	resp, err := client.Get("http://www.msftconnecttest.com/connecttest.txt")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return bytes.Compare(body, []byte("Microsoft Connect Test")) == 0, nil
}

func detectNetworkWithAuthServer(manager *sdunet.Manager) (bool, error) {
	info, err := manager.GetUserInfo()
	if err != nil {
		return false, err
	} else {
		return info.LoggedIn, nil
	}
}

type Action func() error

func retryWithSettings(settings *setting.Settings, action Action, exitNow *bool) error {
	return retry(settings.Control.MaxRetryCount, settings.Control.RetryIntervalSec, action, exitNow)
}
func retry(retryTimes int32, retryIntervalSec int32, action Action, exitNow *bool) error {
	errorList := make([]error, 0)
	aggregateErrors := func(errorList []error) error {
		errorMsg := ""
		for j := 0; j < len(errorList); j++ {
			errorMsg += "[" + fmt.Sprint(j) + "] " + errorList[j].Error()
		}
		return errors.New(errorMsg)
	}
	for i := int32(0); i < retryTimes; i++ {
		err := action()
		if err == nil {
			return nil
		}
		errorList = append(errorList, err)
		if i == retryTimes-1 {
			return aggregateErrors(errorList)
		}

		if exitNow != nil && *exitNow {
			return aggregateErrors(errorList)
		}

		for i := int32(0); i < retryIntervalSec; i++ {
			time.Sleep(time.Second)

			if exitNow != nil && *exitNow {
				return aggregateErrors(errorList)
			}
		}
	}
	// impossible path
	// make compiler happy
	return errors.New("assert failed")
}

var _manager *sdunet.Manager

func getManager(settings *setting.Settings) (*sdunet.Manager, error) {
	if _manager == nil {
		networkInterface := ""
		if settings.Network.StrictMode {
			networkInterface = settings.Network.Interface
		}

		manager, err := sdunet.GetManager(
			settings.Account.Scheme,
			settings.Account.AuthServer,
			settings.Account.Username,
			networkInterface,
			time.Duration(settings.Network.Timeout)*time.Second,
			int(settings.Network.MaxRetryCount),
			time.Duration(settings.Network.RetryIntervalSec)*time.Second,
		)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		_manager = &manager
	}
	return _manager, nil
}

func logout(settings *setting.Settings, exitNow *bool) error {
	return retryWithSettings(settings, func() error {
		log.Println("Logout via web portal...")
		manager, err := getManager(settings)
		if err != nil {
			return err
		}
		err = manager.Logout()
		if err != nil {
			return err
		}
		log.Println("Logged out.")
		return nil
	}, exitNow)
}

func login(settings *setting.Settings, exitNow *bool) error {
	return retryWithSettings(settings, func() error {
		manager, err := getManager(settings)
		if err != nil {
			return err
		}
		return manager.Login(settings.Account.Password)
	}, exitNow)
}

func loginIfNotOnline(settings *setting.Settings, exitNow *bool) error {
	isOnline := false

	err := retryWithSettings(settings, func() error {
		manager, err := getManager(settings)
		if err != nil {
			return err
		}
		isOnline, err = detectNetwork(settings, manager)
		return err
	}, exitNow)

	if err == nil && isOnline {
		log.Println("Network is up. Nothing to do.")
		return nil
	} else {
		// not online
		if err != nil {
			log.Println(err)
		}

		log.Println("Network is down. Log in via web portal...")
		err = login(settings, exitNow)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Logged in.")
		}
		return err
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

	fileExist, err := utils.PathExists(FlagConfigFile)
	if err != nil {
		panic(err)
	}

	if !fileExist {
		version()
		cartman()
		return
	}

	settings, err := setting.LoadSettings(FlagConfigFile)
	if err != nil {
		panic(err)
	}
	//required arguments
	err = checkInterval(settings)
	if err != nil {
		panic(err)
	}
	err = checkAuthServer(settings)
	if err != nil {
		panic(err)
	}
	err = checkPassword(settings)
	if err != nil {
		panic(err)
	}
	err = checkScheme(settings)
	if err != nil {
		panic(err)
	}
	err = checkUsername(settings)
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
		err := retryWithSettings(settings, func() error {
			manager, err := getManager(settings)
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
		err := login(settings, nil)
		if err != nil {
			panic(err)
		}

		return
	}

	if FlagTryOneshoot {
		version()
		err := loginIfNotOnline(settings, nil)
		if err != nil {
			panic(err)
		}
		return
	}

	if FlagLogout {
		version()
		err := logout(settings, nil)
		if err != nil {
			panic(err)
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
					_ = retryWithSettings(settings, func() error {
						log.Println("Logging out...")
						manager, err := getManager(settings)
						if err != nil {
							return err
						}
						err = manager.Logout()
						if err != nil {
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
		_ = loginIfNotOnline(settings, &exitNow)

		for i := int32(0); i < settings.Control.LoopIntervalSec; i++ {
			if exitNow {
				break
			}
			time.Sleep(time.Second)
		}

	}

}
