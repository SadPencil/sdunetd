/*
Copyright © 2018-2022 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/SadPencil/sdunetd/sdunet"
	"github.com/SadPencil/sdunetd/setting"
	"github.com/SadPencil/sdunetd/utils"
	"github.com/flowchartsman/retry"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger *log.Logger = log.New(os.Stderr, "", log.LstdFlags)
var verboseLogger *log.Logger = log.New(ioutil.Discard, "", 0)

func version() {
	fmt.Println(NAME, VERSION)
	fmt.Println(DESCRIPTION)
}

func detectNetwork(ctx context.Context, settings *setting.Settings, manager *sdunet.Manager) (bool, error) {
	if settings.Control.OnlineDetectionMethod == setting.ONLINE_DETECTION_METHOD_AUTH {
		return detectNetworkWithAuthServer(ctx, manager)
	} else if settings.Control.OnlineDetectionMethod == setting.ONLINE_DETECTION_METHOD_MS {
		return detectNetworkWithMicrosoft(ctx, manager)
	} else {
		return detectNetworkWithAuthServer(ctx, manager)
	}
}

func detectNetworkWithMicrosoft(ctx context.Context, manager *sdunet.Manager) (bool, error) {
	client, err := manager.GetHttpClient()
	if err != nil {
		return false, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://www.msftconnecttest.com/connecttest.txt", nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
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

func detectNetworkWithAuthServer(ctx context.Context, manager *sdunet.Manager) (bool, error) {
	info, err := manager.GetUserInfo(ctx)
	if err != nil {
		return false, err
	} else {
		return info.LoggedIn, nil
	}
}

func retryWithSettings(ctx context.Context, settings *setting.Settings, action func() error) error {
	return _retry(ctx, int(settings.Control.MaxRetryCount), int(settings.Control.RetryIntervalSec), action)
}

func _retry(ctx context.Context, retryTimes int, retryIntervalSec int, action func() error) error {
	retrier := retry.NewRetrier(retryTimes, time.Duration(retryIntervalSec)*time.Second, time.Duration(retryIntervalSec)*time.Second)
	return retrier.RunContext(ctx, func(ctx context.Context) error {
		return action()
	})
}

func logout(ctx context.Context, settings *setting.Settings) error {
	logger.Println("Logout via web portal...")
	return retryWithSettings(ctx, settings, func() error {
		manager, err := getManager(ctx, settings)
		if err != nil {
			return err
		}
		err = manager.Logout(ctx)
		if err != nil {
			return err
		}
		logger.Println("Logged out.")
		return nil
	})
}

func login(ctx context.Context, settings *setting.Settings) error {
	logger.Println("Log in via web portal...")
	return retryWithSettings(ctx, settings, func() error {
		manager, err := getManager(ctx, settings)
		if err != nil {
			return err
		}
		err = manager.Login(ctx, settings.Account.Password)
		if err != nil {
			return err
		}
		logger.Println("Logged in.")
		return nil
	})
}

func loginIfNotOnline(ctx context.Context, settings *setting.Settings) error {
	isOnline := false

	err := retryWithSettings(ctx, settings, func() error {
		manager, err := getManager(ctx, settings)
		if err != nil {
			return err
		}
		isOnline, err = detectNetwork(ctx, settings, manager)
		return err
	})

	if err == nil && isOnline {
		logger.Println("Network is up. Nothing to do.")
		return nil
	} else {
		// not online
		if err != nil {
			logger.Println(err)
		}

		logger.Println("Network is down.")
		err = login(ctx, settings)
		if err != nil {
			logger.Println(err)
		}
		return err
	}
}

func onExit(action func()) {
	// set up handler for SIGINT and SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				action()
				return
			default:
			}
		}
	}()
}

func main() {
	var FlagShowHelp bool
	flag.BoolVar(&FlagShowHelp, "h", false, "standalone: show this help.")

	var FlagShowVersion bool
	flag.BoolVar(&FlagShowVersion, "V", false, "standalone: show the version.")

	var FlagConfigFile string
	flag.StringVar(&FlagConfigFile, "c", "", "the path to the config.json file. Leave it blank to generate a new configuration file interactively.")

	var FlagLogOutput string
	flag.StringVar(&FlagLogOutput, "o", "", "the path to the output file of log message. Empty means stderr, and - means stdout.")

	var FlagNoAttribute bool
	flag.BoolVar(&FlagNoAttribute, "m", false, "option: output log without the timestamp prefix. Turn it on when running as a systemd service.")

	var FlagIPDetect bool
	flag.BoolVar(&FlagIPDetect, "a", false, "standalone: detect the IP address from the authenticate server. Useful when behind a NAT router.")

	var FlagOneshoot bool
	flag.BoolVar(&FlagOneshoot, "f", false, "standalone: login to the network for once, regardless of whether the network is offline.")

	var FlagTryOneshoot bool
	flag.BoolVar(&FlagTryOneshoot, "t", false, "standalone: login to the network for once, only if the network is offline.")

	var FlagLogout bool
	flag.BoolVar(&FlagLogout, "l", false, "standalone: logout from the network for once.")

	var FlagVerbose bool
	flag.BoolVar(&FlagVerbose, "v", false, "option: output verbose log")

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

	//open the log file for writing
	if FlagLogOutput == "" {
		log.New(os.Stderr, "", log.LstdFlags)
	} else if FlagLogOutput == "-" {
		log.New(os.Stdout, "", log.LstdFlags)
	} else {
		logFile, err := os.OpenFile(FlagLogOutput, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		defer logFile.Close()
		if err != nil {
			logger.Panicln(err)
			return
		}
		logger = log.New(logFile, "", log.LstdFlags)
	}
	if FlagNoAttribute {
		logger.SetFlags(0)
	}

	if FlagVerbose {
		verboseLogger = logger
	}

	if FlagIPDetect {
		err := retryWithSettings(context.Background(), settings, func() error {
			manager, err := getManager(context.Background(), settings)
			if err != nil {
				return err
			}
			info, err := manager.GetUserInfo(context.Background())
			if err != nil {
				return err
			}
			fmt.Println(info.ClientIP)
			return nil
		})
		if err != nil {
			logger.Panicln(err)
		}
		return
	}

	if FlagOneshoot {
		version()
		err := login(context.Background(), settings)
		if err != nil {
			logger.Panicln(err)
		}

		return
	}

	if FlagTryOneshoot {
		version()
		err := loginIfNotOnline(context.Background(), settings)
		if err != nil {
			logger.Panicln(err)
		}
		return
	}

	if FlagLogout {
		version()
		err := logout(context.Background(), settings)
		if err != nil {
			logger.Panicln(err)
		}
		return
	}

	version()

	// main loop
	ctx, cancelFunc := context.WithCancel(context.Background())
	onExit(func() {
		logger.Println("Exiting...")
		cancelFunc()
	})

	_ = loginIfNotOnline(ctx, settings)

	for {
		canceled := false

		select {
		case <-time.After(time.Duration(settings.Control.LoopIntervalSec) * time.Second):
			_ = loginIfNotOnline(ctx, settings)
		case <-ctx.Done():
			canceled = true
		}

		if canceled {
			break
		}
	}

	// Cleanup
	if settings.Control.LogoutWhenExit {
		ctx, cancelFunc := context.WithCancel(context.Background())
		onExit(func() {
			logger.Println("Force exiting. Abort logging out action...")
			cancelFunc()
		})
		_ = logout(ctx, settings)
	}
}
