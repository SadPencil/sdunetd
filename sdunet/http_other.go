//go:build !linux
// +build !linux

/*
Copyright © 2018-2022 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package sdunet

import (
	"errors"
	retryableHttp "github.com/hashicorp/go-retryablehttp"
	"log"
	"net/http"
	"time"
)

func getHttpClient(forceNetworkInterface string, timeout time.Duration, retryCount int, retryWait time.Duration, logger *log.Logger) (*http.Client, error) {
	client := retryableHttp.NewClient()
	client.HTTPClient.Timeout = timeout
	client.RetryMax = retryCount
	client.RetryWaitMin = retryWait
	client.RetryWaitMax = retryWait
	client.Logger = logger

	if forceNetworkInterface != "" {
		return nil, errors.New("the strict mode is only available in Linux")
	}
	return client.StandardClient(), nil
}
