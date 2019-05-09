// +build curl

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
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func getRawUserInfo(scheme, server string, interfaceWtf string) (info map[string]interface{}, err error) {
	req, err := http.NewRequest("GET", scheme+"://"+server+"/cgi-bin/rad_user_info", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("callback", "jQuery")
	req.URL.RawQuery = q.Encode()

	//fmt.Println(req.URL.String())
	var respBody []byte

	if len(interfaceWtf) == 0 {
		respBody, err = commandCurl("-H", "Accept:application/json", req.URL.String())
	} else {
		respBody, err = commandCurl("--interface", "if!"+interfaceWtf, "-H", "Accept:application/json", req.URL.String())
	}

	if err != nil {
		return nil, err
	}

	if bytes.Equal(respBody[0:7], []byte("jQuery(")) {
		respBody = respBody[7:]
		respBody = respBody[:len(respBody)-1]
	}

	err = json.Unmarshal(respBody, &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func getRawChallenge(scheme, server, rawUsername, sduIPv4 string, interfaceWtf string) (challenge map[string]interface{}, err error) {
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

	var respBody []byte

	if len(interfaceWtf) == 0 {
		respBody, err = commandCurl("-H", "Accept:application/json", req.URL.String())
	} else {
		respBody, err = commandCurl("--interface", "if!"+interfaceWtf, "-H", "Accept:application/json", req.URL.String())
	}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &challenge)
	if err != nil {
		return nil, err
	}
	return challenge, nil
}

func getSduUserInfo(scheme, server string, interfaceWtf string) (logined bool, sduIPv4 string, err error) {
	output, err := getRawUserInfo(scheme, server, interfaceWtf)
	if err != nil {
		return false, "", err
	}
	errorStr := output["error"].(string)
	sduIPv4 = output["online_ip"].(string)
	//fmt.Println(sduIPv4)
	return errorStr == "ok", sduIPv4, nil
}

func getChallengeID(scheme, server, rawUsername, sduIPv4 string, interfaceWtf string) (challenge string, err error) {
	output, err := getRawChallenge(scheme, server, rawUsername, sduIPv4, interfaceWtf)
	if err != nil {
		return "", err
	}
	challenge = output["challenge"].(string)
	//fmt.Println(challenge)
	return challenge, nil
}

func commandCurl(args ...string) (stdout []byte, err error) {
	curlName := os.Getenv("SDUNETD_CURL_NAME")
	if len(curlName) == 0 {
		curlName = "curl"
	}
	if runtime.GOOS == `windows` {
		return nil, errors.New("curl doesn't support specified interface name on Windows. Please use the normal build of sdunetd")
	}

	return command(curlName, args...)
}

func command(name string, args ...string) (stdout []byte, err error) {
	var cmd *exec.Cmd
	cmd = exec.Command(name, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Println("[ERROR] Failed to execute command", name, "with arguments", strings.Join(args, ", "), "\n", err, "\n", stderr.String())
		return out.Bytes(), err
	}
	return out.Bytes(), nil
}

func loginDigest(scheme, server, rawUsername, rawPassword, sduIPv4 string, interfaceWtf string) (err error) {
	challenge, err := getChallengeID(scheme, server, rawUsername, sduIPv4, interfaceWtf)
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

	var respBody []byte

	if len(interfaceWtf) == 0 {
		respBody, err = commandCurl("-H", "Accept:application/json", req.URL.String())
	} else {
		respBody, err = commandCurl("--interface", "if!"+interfaceWtf, "-H", "Accept:application/json", req.URL.String())
	}

	if err != nil {
		return err
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
