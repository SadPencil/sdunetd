/*
Copyright © 2018-2022 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package sdunet

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type MangerBase struct {
	client                *http.Client
	Scheme                string
	Server                string
	ForceNetworkInterface string
	Timeout               time.Duration
	MaxRetryCount         int
	RetryWait             time.Duration
}

type Manager struct {
	MangerBase
	Username string
	ClientIP string
}

type UserInfo struct {
	ClientIP string
	LoggedIn bool
}

func GetManager(scheme string, server string, username string, forceNetworkInterface string,
	timeout time.Duration, maxRetryCount int, retryWait time.Duration) (Manager, error) {
	base := MangerBase{
		Scheme:                scheme,
		Server:                server,
		ForceNetworkInterface: forceNetworkInterface,
		Timeout:               timeout,
		MaxRetryCount:         maxRetryCount,
		RetryWait:             retryWait,
	}
	info, err := base.GetUserInfo()
	if err != nil {
		return Manager{}, err
	}
	return Manager{
		MangerBase: base,
		Username:   username,
		ClientIP:   info.ClientIP,
	}, nil
}

func (m MangerBase) GetHttpClient() (*http.Client, error) {
	if m.client == nil {
		client, err := getHttpClient(m.ForceNetworkInterface, m.Timeout, m.MaxRetryCount, m.RetryWait)
		if err != nil {
			return nil, err
		}
		m.client = client
	}
	return m.client, nil
}

func (m Manager) getRawChallenge() (map[string]interface{}, error) {
	return m.httpJsonQuery(
		"/cgi-bin/get_challenge",
		map[string][]string{
			"username": {m.Username},
			"ip":       {m.ClientIP},
		},
		"jQuery",
	)
}

func (m MangerBase) getRawUserInfo() (map[string]interface{}, error) {
	return m.httpJsonQuery(
		"/cgi-bin/rad_user_info",
		map[string][]string{},
		"jQuery",
	)
}
func (m MangerBase) GetUserInfo() (UserInfo, error) {
	output, err := m.getRawUserInfo()
	if err != nil {
		return UserInfo{}, err
	}
	return UserInfo{
		ClientIP: output["online_ip"].(string),
		LoggedIn: output["error"].(string) == "ok",
	}, nil
}

func (m MangerBase) httpJsonQuery(relativeUrl string, getParams url.Values, jsonCallback string) (map[string]interface{}, error) {
	if relativeUrl[0] != '/' {
		return nil, errors.New("invalid relative url")
	}

	client, err := m.GetHttpClient()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", m.Scheme+"://"+m.Server+relativeUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	if jsonCallback != "" {
		getParams.Add("callback", jsonCallback)
	}
	req.URL.RawQuery = getParams.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(respBody[0:len(jsonCallback)+1], []byte(jsonCallback+"(")) {
		respBody = respBody[len(jsonCallback)+1:]
		respBody = respBody[:len(respBody)-1]
	}

	var output map[string]interface{}
	err = json.Unmarshal(respBody, &output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (m Manager) getChallengeID() (string, error) {
	output, err := m.getRawChallenge()
	if err != nil {
		return "", err
	}
	return output["challenge"].(string), nil
}

func (m Manager) Login(password string) error {
	challenge, err := m.getChallengeID()
	if err != nil {
		return err
	}

	var dataInfoStr, dataPasswordMd5Str, dataChecksumStr string
	err = sdunetChallenge(m.Username, password, m.ClientIP, challenge, &dataInfoStr, &dataPasswordMd5Str, &dataChecksumStr)
	if err != nil {
		return err
	}

	output, err := m.httpJsonQuery(
		"/cgi-bin/srun_portal",
		map[string][]string{
			"action":   {"login"},
			"username": {m.Username},
			"password": {dataPasswordMd5Str},
			"ac_id":    {"1"},
			"ip":       {m.ClientIP},
			"info":     {dataInfoStr},
			"chksum":   {dataChecksumStr},
			"n":        {"200"},
			"type":     {"1"},
		},
		"jQuery",
	)
	if err != nil {
		return err
	}

	errorStr := output["error"].(string)
	if errorStr == "ok" {
		return nil
	} else {
		return errors.New(errorStr)
	}
}

func (m Manager) Logout() error {
	output, err := m.httpJsonQuery(
		"/cgi-bin/srun_portal",
		map[string][]string{
			"ac_id":    {"1"},
			"action":   {"logout"},
			"username": {m.Username},
		},
		"jQuery",
	)
	if err != nil {
		return err
	}
	errorStr := output["error"].(string)
	if errorStr == "ok" {
		return nil
	} else {
		return errors.New(errorStr)
	}
}
