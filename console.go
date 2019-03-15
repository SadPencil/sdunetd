package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
)

func cartman() {
	Settings := NewSettings()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Looks like the config file doesn't exist.")
	fmt.Println("That's okay if you just want to do a Portal authentication for once, or you want me to help you generate a config file.")
	fmt.Println("A few questions need to be answered. Leave it blank if you want the default answer in the bracket.")

	var err error

	for {
		fmt.Println("Question 1. What's your username? []")

		Settings.Account.Username, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		err = checkUsername(&Settings)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}
	for {
		fmt.Println("Question 2. What's your password? []")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println() // it's necessary to add a new line after user's input
		if err != nil {
			panic(err)
		} else {
			Settings.Account.Password = string(bytePassword)
		}
		err = checkPassword(&Settings)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Great. Your password contains", fmt.Sprint(len(Settings.Account.Password)), "characters.")
			break
		}
	}

	for {
		fmt.Println("Question 3. What's your school's authenticate server's ip address? [" + DEFAULT_AUTH_SERVER + "]")
		fmt.Println("Hint: You can also write down the server's FQDN if necessary.")
		Settings.Account.AuthServer, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		err = checkAuthServer(&Settings)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}
	for {
		fmt.Println("Question 4. Is the authenticate server use HTTP protocol, or HTTPS? [" + DEFAULT_AUTH_SCHEME + "]")
		Settings.Account.Scheme, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		err = checkScheme(&Settings)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}

	for {
		fmt.Println("Question 5.")
		interfaces, err := net.Interfaces()
		if err != nil {
			panic(err)
		}
		var interfaceName string
		for _, interfaceWtf := range interfaces {
			name, err := GetIPFromInterface(interfaceWtf.Name)
			if err == nil {
				_, err = fmt.Println(interfaceWtf.Name + "\t" + name)
				interfaceName = interfaceWtf.Name
			}
		}
		if len(interfaces) == 0 {
			fmt.Println("There is not even a network interface with a valid IPv4 address.")
			fmt.Println("Screw you guys, I'm going home.")
			return
		}
		fmt.Println("Network interfaces are listed above. Which one is connected to the Portal network? [" + interfaceName + "]")
		Settings.Network.Interface, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		Settings.Network.Interface = strings.TrimSpace(Settings.Network.Interface)
		if Settings.Network.Interface == "" {
			Settings.Network.Interface = interfaceName
		}

		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}
	var saveFile bool
	for {
		fmt.Println("That's it. That's all the information needed. Would you like to save it to a config file? [y/n]")
		yesOrNoStr, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		yesOrNoStr = strings.ToLower(strings.TrimSpace(yesOrNoStr))

		if yesOrNoStr == "y" {
			saveFile = true
			break
		} else if yesOrNoStr == "n" {
			saveFile = false
			break
		} else {
			fmt.Println("All you need to do is to answer me yes or no. Don't be a pussy.")
		}
	}

	if saveFile {
		fmt.Println("Where to save the file? [" + DEFAULT_CONFIG_FILENAME + "]")
		filename, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		filename = strings.TrimSpace(filename)
		if filename == "" {
			filename = DEFAULT_CONFIG_FILENAME
		}
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
		} else {

			jsonBytes, err := json.Marshal(Settings)
			if err != nil {
				fmt.Println(err)
				saveFile = false
			} else {
				_, err := f.WriteString(string(jsonBytes))
				if err != nil {
					fmt.Println(err)
					saveFile = false
				} else {
					fmt.Println(`File saved. You may re-run the program with the "-c config" flag.`)
				}
			}
		}
	}
	if !saveFile {
		var yesOrNo bool
		for {
			fmt.Println("File not saved. Login to the network right away? [Y/n]")
			if err != nil {
				panic(err)
			}
			yesOrNoStr, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			yesOrNoStr = strings.ToLower(strings.TrimSpace(yesOrNoStr))

			if yesOrNoStr == "y" || yesOrNoStr == "" {
				yesOrNo = true
				break
			} else if yesOrNoStr == "n" {
				yesOrNo = false
				break
			} else {
				fmt.Println("All you need to do is to answer me yes or no. Don't be a pussy.")
			}
		}
		if yesOrNo {
			log.Println("Log in via web portal...")
			err := login(Settings.Account.Scheme,
				Settings.Account.AuthServer,
				Settings.Account.Username,
				Settings.Account.Password,
				Settings.Network.CustomIP,
				Settings.Network.Interface,
				false)
			if err != nil {
				log.Println("Login failed.", err)
			}
		}
	}
}
