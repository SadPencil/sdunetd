// +build !curl

/*
Copyright © 2018-2019 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func getRawUserInfo(scheme, server string, interfaceWtf string) (info map[string]interface{}, err error) {
	client, err := getHttpClientNormal(interfaceWtf)
	if err != nil {
		return nil, err
	}
	return getRawUserInfoNormal(scheme, server, client)
}

func getRawUserInfoNormal(scheme, server string, client *http.Client) (info map[string]interface{}, err error) {
	req, err := http.NewRequest("GET", scheme+"://"+server+"/cgi-bin/rad_user_info", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("callback", "jQuery")
	req.URL.RawQuery = q.Encode()

	//fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if bytes.Equal(respBody[0:7], []byte("jQuery(")) {
		respBody = respBody[7:]
		respBody = respBody[:len(respBody)-1]
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	err = json.Unmarshal(respBody, &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

//func getRawChallenge(scheme, server, rawUsername, sduIPv4 string, interfaceWtf string) (challenge map[string]interface{}, err error) {
//	client, err := getHttpClientNormal(interfaceWtf)
//	if err != nil {
//		return nil, err
//	}
//	return getRawChallengeNormal(scheme, server, rawUsername, sduIPv4, client)
//}

func getRawChallengeNormal(scheme, server, rawUsername, sduIPv4 string, client *http.Client) (challenge map[string]interface{}, err error) {
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", scheme+"://"+server+"/cgi-bin/get_challenge", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("username", rawUsername)
	q.Add("ip", sduIPv4)

	req.URL.RawQuery = q.Encode()

	//fmt.Println(req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	err = json.Unmarshal(respBody, &challenge)
	if err != nil {
		return nil, err
	}
	return challenge, nil
}
func getSduUserInfo(scheme, server string, interfaceWtf string) (logined bool, sduIPv4 string, err error) {
	client, err := getHttpClientNormal(interfaceWtf)
	if err != nil {
		return false, "", err
	}
	return getSduUserInfoNormal(scheme, server, client)
}
func getSduUserInfoNormal(scheme, server string, client *http.Client) (logined bool, sduIPv4 string, err error) {
	output, err := getRawUserInfoNormal(scheme, server, client)
	if err != nil {
		return false, "", err
	}
	errorStr := output["error"].(string)
	sduIPv4 = output["online_ip"].(string)
	//fmt.Println(sduIPv4)
	return errorStr == "ok", sduIPv4, nil
}

//func getChallengeID(scheme, server, rawUsername, sduIPv4 string, interfaceWtf string) (challenge string, err error) {
//	client, err := getHttpClientNormal(interfaceWtf)
//	if err != nil {
//		return "", nil
//	}
//	return getChallengeIDNormal(scheme, server, rawUsername, sduIPv4, client)
//}

func getChallengeIDNormal(scheme, server, rawUsername, sduIPv4 string, client *http.Client) (challenge string, err error) {
	output, err := getRawChallengeNormal(scheme, server, rawUsername, sduIPv4, client)
	if err != nil {
		return "", err
	}
	challenge = output["challenge"].(string)
	//fmt.Println(challenge)
	return challenge, nil
}

func getHttpClientNormal(interfaceWtf string) (client *http.Client, err error) {
	var c http.Client
	if len(interfaceWtf) == 0 {
		c, err = getHttpClient(false, "", "")
	} else {
		c, err = getHttpClient(true, "", interfaceWtf)
	}
	if err != nil {
		return nil, err
	} else {
		return &c, nil
	}
}

func getHttpClient(strict bool, localIPv4, interfaceWtf string) (client http.Client, err error) {
	client = http.Client{}
	if strict {
		var ipv4 string
		if interfaceWtf == "" {
			ipv4 = localIPv4
		} else if localIPv4 == "" {
			ipv4, err = GetIPv4FromInterface(interfaceWtf)
			if err != nil {
				return client, err
			}
		} else {
			ipv4 = localIPv4
		}

		//fmt.Println("STRICT MODE: ",ipv4+":0")
		localAddress, err := net.ResolveTCPAddr("tcp", ipv4+":0")
		if err != nil {
			return client, err
		}

		//https://stackoverflow.com/questions/30552447/how-to-set-which-ip-to-use-for-a-http-request
		// Create a transport like http.DefaultTransport, but with a specified localAddr

		client.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				LocalAddr: localAddress,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

	}
	return client, nil
}
func loginDigest(scheme, server, rawUsername, rawPassword, sduIPv4 string, interfaceWtf string) (err error) {
	client, err := getHttpClientNormal(interfaceWtf)
	if err != nil {
		return err
	}
	return loginDigestNormal(scheme, server, rawUsername, rawPassword, sduIPv4, client)
}

func loginDigestNormal(scheme, server, rawUsername, rawPassword, sduIPv4 string, client *http.Client) (err error) {
	challenge, err := getChallengeIDNormal(scheme, server, rawUsername, sduIPv4, client)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", scheme+"://"+server+"/cgi-bin/srun_portal", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")

	var dataInfoStr, dataPasswordMd5Str, dataChecksumStr string

	err = sdunetChallenge(rawUsername, rawPassword, sduIPv4, challenge, &dataInfoStr, &dataPasswordMd5Str, &dataChecksumStr)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("action", "login")
	q.Add("username", rawUsername)
	q.Add("password", dataPasswordMd5Str)
	q.Add("ac_id", "1")
	q.Add("ip", sduIPv4)
	q.Add("info", dataInfoStr)
	q.Add("chksum", dataChecksumStr)
	q.Add("n", "200")
	q.Add("type", "1")

	req.URL.RawQuery = q.Encode()

	//fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	var output map[string]interface{}
	err = json.Unmarshal(respBody, &output)
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
