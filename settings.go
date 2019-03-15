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
	Interface string `json:"interface"`
	CustomIP  string `json:"custom_ip"`
}

type Control struct {
	StrictMode bool  `json:"strict"`
	Interval   int32 `json:"interval"`
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
		Control: Control{Interval: DEFAULT_INTERVAL},
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

	err = json.Unmarshal(file, &settings)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurs while unmarshalling config file.")
		return NewSettings(), err
	}

	return settings, nil
}
