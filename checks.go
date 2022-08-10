/*
Copyright © 2018-2022 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"errors"
	"github.com/SadPencil/sdunetd/setting"
	"strings"
)

func checkAuthServer(settings *setting.Settings) (err error) {
	settings.Account.AuthServer = strings.TrimSpace(settings.Account.AuthServer)

	if settings.Account.AuthServer == "" {
		settings.Account.AuthServer = setting.DEFAULT_AUTH_SERVER
	}

	if len(settings.Account.AuthServer) >= 5 {
		if settings.Account.AuthServer[0:5] == "http:" {
			return errors.New(`I'm asking you about the server's FQDN, not the URI. Please remove "http://".`)
		}
	}
	if len(settings.Account.AuthServer) >= 6 {
		if settings.Account.AuthServer[0:6] == "https:" {
			return errors.New(`I'm asking you about the server's FQDN, not the URI. Please remove "https:".`)
		}
	}

	return nil
}
func checkUsername(settings *setting.Settings) (err error) {
	settings.Account.Username = strings.TrimSpace(settings.Account.Username)

	if settings.Account.Username == "" {
		return errors.New("Dude. I can't login without a username.")
	}

	return nil
}
func checkPassword(settings *setting.Settings) (err error) {
	settings.Account.Password = strings.TrimSpace(settings.Account.Password)

	if settings.Account.Password == "" {
		return errors.New("Give me your password, bitch. 'Cause I need it to hack into your... Ah, I mean, login the network.")
	}

	return nil
}
func checkScheme(settings *setting.Settings) (err error) {
	settings.Account.Scheme = strings.ToLower(strings.TrimSpace(settings.Account.Scheme))

	if settings.Account.Scheme == "" {
		settings.Account.Scheme = setting.DEFAULT_AUTH_SCHEME
	} else if strings.Contains(settings.Account.Scheme, "fuck") {
		return errors.New("Fuck you! You are a fucking asshole.")
	} else if !(settings.Account.Scheme == "http" || settings.Account.Scheme == "https") {
		return errors.New("HTTP or HTTPS. Douchebag.")
	}
	return nil
}
func checkInterval(settings *setting.Settings) error {
	if settings.Control.LoopIntervalSec == 0 {
		return errors.New("interval should be more than 0 seconds")
	}
	return nil
}
