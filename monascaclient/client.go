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
	"encoding/json"
	"fmt"
	log "github.hpe.com/kronos/kelog"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const timeFormat = "2006-01-02T15:04:05Z"

var (
	baseURL        = "http://monasca.monasca:8070"
	requestTimeout = 10
	allowInsecure  = false
)

func SetBaseURL(url string) {
	log.Infof("Using Monasca API URL of '%v'", url)
	baseURL = url
}

// Value of true should only be used for testing!!!
func SetInsecure(insecure bool) {
	log.Infof("Setting allowInsecure '%v'", insecure)
	allowInsecure = insecure
}

func SetTimeout(timeout int) {
	log.Infof("Using Monasca API Timeout of '%v'", timeout)
	requestTimeout = timeout
}

type Client struct {
}

func (p *Client) GetStatistics(metricName string, startTime time.Time, endTime time.Time, period int64, dimensions map[string]string, headers http.Header) (*StatisticsResponse, error) {
	// TODO: Review this entire function
	// TODO: Handle errors. How does Go work?
	queryParameters := map[string]string{
		"name":       metricName,
		"statistics": "avg",
		"start_time": startTime.UTC().Format(timeFormat),
		"end_time":   endTime.UTC().Format(timeFormat),
		"period":     fmt.Sprintf("%d", period),
	}

	monascaURL := createMonascaAPIURL("v2.0/metrics/statistics", queryParameters, dimensions)

	body, monascaErr := callMonasca(monascaURL, headers)
	if monascaErr != nil {
		return nil, monascaErr
	}
	statisticsResponse := StatisticsResponse{}
	err := json.Unmarshal(body, &statisticsResponse)

	if err != nil {
		panic(err.Error())
	}

	return &statisticsResponse, nil
}

func callMonasca(monascaURL string, headers http.Header) ([]byte, error) {

	log.Infof("Calling Monasca with '%v'", monascaURL)

	req, err := http.NewRequest("GET", monascaURL, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	for header, values := range headers {
		for index := range values {
			value := values[index]
			if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
				log.Infof("Fixing the value '%s' for header '%s'", value, header)
				value = value[1 : len(value)-1]
			}
			req.Header.Add(header, value)
		}
	}

	timeout := time.Duration(requestTimeout) * time.Second
	var client *http.Client
	if !allowInsecure {
		client = &http.Client{Timeout: timeout}
	} else {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}

		client = &http.Client{Timeout: timeout, Transport: transCfg}
	}
	resp, err := client.Do(req)

	if err != nil || resp == nil {
		log.Errorf("Monasca request failed '%v'", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode != 200 {
		log.Errorf("Monasca request failed %d; '%v'", resp.StatusCode, string([]byte(body)))
		panic(fmt.Errorf("Error: %d %s", resp.StatusCode, string([]byte(body))))
	}

	return body, nil
}

func createMonascaAPIURL(path string, queryParameters map[string]string, dimensions map[string]string) string {

	monascaURL, parseErr := url.Parse(baseURL)
	if parseErr != nil {
		panic(parseErr.Error())
	}
	monascaURL.Path = path

	q := url.Values{}
	for key := range queryParameters {
		q.Add(key, queryParameters[key])
	}
	if len(dimensions) > 0 {
		dimensionsSlice := make([]string, 0, len(dimensions))
		for key := range dimensions {
			dimensionsSlice = append(dimensionsSlice, key+":"+dimensions[key])
		}
		// Make sure dimensions are always in correct order to ensure tests pass
		sort.Strings(dimensionsSlice)
		q.Add("dimensions", strings.Join(dimensionsSlice, ","))
	}
	monascaURL.RawQuery = q.Encode()

	return monascaURL.String()
}

type MetricsResponse struct {
	Links    []Link   `json:"links"`
	Elements []Metric `json:"elements"`
}

type Metric struct {
	Name       string            `json:"name"`
	Dimensions map[string]string `json:"dimensions"`
}

func (p *Client) GetMetrics(metricName string, dimensions map[string]string, headers http.Header) ([]Metric, error) {
	queryParameters := map[string]string{
		"name": metricName,
	}

	monascaURL := createMonascaAPIURL("v2.0/metrics", queryParameters, dimensions)

	body, monascaErr := callMonasca(monascaURL, headers)
	if monascaErr != nil {
		return nil, monascaErr
	}

	metricsResponse := MetricsResponse{}
	err := json.Unmarshal(body, &metricsResponse)

	if err != nil {
		panic(err.Error())
	}
	return metricsResponse.Elements, nil
}
