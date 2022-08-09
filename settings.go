/*
Copyright © 2018-2019 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Account struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	AuthServer string `json:"server"`
	Scheme     string `json:"scheme"`
}

type Log struct {
	Filename string `json:"filename"`
}

type Network struct {
	Interface  string `json:"interface"`
	CustomIP   string `json:"custom_ip"`
	StrictMode bool   `json:"strict"`
}

type Control struct {
	Interval              int32  `json:"interval"`
	LogoutWhenExit        bool   `json:"logout_when_exit"`
	OnlineDetectionMethod string `json:"online_detection_method"`
}

type Settings struct {
	Account Account `json:"account"`
	Log     Log     `json:"log"`
	Network Network `json:"network"`
	Control Control `json:"control"`
}

func NewSettings() Settings {
	return Settings{
		Account: Account{Scheme: DEFAULT_AUTH_SCHEME, AuthServer: DEFAULT_AUTH_SERVER},
		Control: Control{
			Interval:              DEFAULT_INTERVAL,
			OnlineDetectionMethod: DEFAULT_ONLINE_DETECTION_METHOD,
			LogoutWhenExit:        false,
		},
	}
}

// LoadSettings -- Load settings from config file
func LoadSettings(configPath string) (settings Settings, err error) {
	// LoadSettings from config file
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurs while reading config file.")
		return NewSettings(), err
	}

	settings = NewSettings()
	err = json.Unmarshal(file, &settings)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurs while unmarshalling config file.")
		return NewSettings(), err
	}

	return settings, nil
}
