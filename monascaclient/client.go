// Copyright 2017 Hewlett Packard Enterprise Development LP
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package monascaclient

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"bytes"
)

var (
	defaultURL      = "http://localhost:8070"
	defaultTimeout  = 60
	defaultInsecure = false
)

var (
	monClient = &Client{
		baseURL:        defaultURL,
		requestTimeout: defaultTimeout,
		allowInsecure:  defaultInsecure,
	}
)

func SetBaseURL(url string) {
	monClient.SetBaseURL(url)
}

func SetDefaultBaseURL(url string) {
	defaultURL = url
}

func SetInsecure(insecure bool) {
	monClient.SetInsecure(insecure)
}

func SetDefaultInsecure(insecure bool) {
	defaultInsecure = insecure
}

func SetTimeout(timeout int) {
	monClient.SetTimeout(timeout)
}

func SetDefaultTimeout(timeout int) {
	defaultTimeout = timeout
}

func SetHeaders(headers http.Header) {
	monClient.SetHeaders(headers)
}

type Client struct {
	baseURL        string
	requestTimeout int
	allowInsecure  bool
	headers        http.Header
}

func New() *Client {
	return &Client{
		baseURL:        defaultURL,
		requestTimeout: defaultTimeout,
		allowInsecure:  defaultInsecure,
	}
}


func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// Value of true should only be used for testing!!!
func (c *Client) SetInsecure(insecure bool) {
	c.allowInsecure = insecure
}

func (c *Client) SetTimeout(timeout int) {
	c.requestTimeout = timeout
}

func (c *Client) SetHeaders(headers http.Header) {
	c.headers = headers
}

func (c *Client) callMonasca(monascaURL string, method string, requestBody *[]byte) ([]byte, error) {
	var req *http.Request
	var err error

	if requestBody == nil {
		req, err = http.NewRequest(method, monascaURL, nil)
	} else {
		req, err = http.NewRequest(method, monascaURL, bytes.NewBuffer(*requestBody))
	}


	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	for header, values := range c.headers {
		for index := range values {
			value := values[index]
			if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
				value = value[1 : len(value)-1]
			}
			req.Header.Add(header, value)
		}
	}

	timeout := time.Duration(c.requestTimeout) * time.Second
	var client *http.Client
	if !c.allowInsecure {
		client = &http.Client{Timeout: timeout}
	} else {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}

		client = &http.Client{Timeout: timeout, Transport: transCfg}
	}
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 201  {
		return nil, fmt.Errorf("Error: %d %s", resp.StatusCode, string([]byte(body)))
	}

	return body, nil
}

func (c *Client) createMonascaAPIURL(path string, urlValues url.Values) (string, error) {

	monascaURL, parseErr := url.Parse(c.baseURL)
	if parseErr != nil {
		return "", parseErr
	}
	monascaURL.Path = path

	if urlValues != nil {
		monascaURL.RawQuery = urlValues.Encode()
	}

	return monascaURL.String(), nil
}
