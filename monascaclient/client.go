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
	"github.com/monasca/golang-monascaclient/monascaclient/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const timeFormat = "2006-01-02T15:04:05Z"

var (
	defaultURL      = "http://localhost:8070"
	defaultTimeout  = 10
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

func GetMetrics(metricName string, dimensions map[string]string) ([]models.Metric, error) {
	return monClient.GetMetrics(metricName, dimensions)
}

func GetStatistics(metricName string, startTime time.Time, endTime time.Time, period int64, dimensions map[string]string) (*models.StatisticsResponse, error) {
	return monClient.GetStatistics(metricName, startTime, endTime, period, dimensions)
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

func (p *Client) GetMetrics(metricName string, dimensions map[string]string) ([]models.Metric, error) {
	queryParameters := map[string]string{
		"name": metricName,
	}

	monascaURL, URLerr := p.createMonascaAPIURL("v2.0/metrics", queryParameters, dimensions)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := p.callMonasca(monascaURL)
	if monascaErr != nil {
		return nil, monascaErr
	}

	metricsResponse := models.MetricsResponse{}
	err := json.Unmarshal(body, &metricsResponse)
	if err != nil {
		return nil, err
	}

	return metricsResponse.Elements, nil
}

func (p *Client) GetDimensionValues(metricName, dimensionName string) ([]string, error) {
	queryParameters := map[string]string{
		"dimension_name": dimensionName,
	}
	if metricName != "" {
		queryParameters["metric_name"] = metricName
	}

	monascaURL, err := p.createMonascaAPIURL("/v2.0/metrics/dimensions/names/values", queryParameters, map[string]string{})
	if err != nil {
		return nil, err
	}

	body, monascaErr := p.callMonasca(monascaURL)
	if monascaErr != nil {
		return nil, monascaErr
	}

	var response models.DimensionValueResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err.Error())
	}

	results := []string{}
	for _, value := range response.Elements {
		results = append(results, value.Value)
	}

	return results, nil
}

func (p *Client) GetStatistics(metricName string, startTime time.Time, endTime time.Time, period int64, dimensions map[string]string) (*models.StatisticsResponse, error) {
	// TODO: Review this entire function
	// TODO: Handle errors. How does Go work?
	queryParameters := map[string]string{
		"name":       metricName,
		"statistics": "avg",
		"start_time": startTime.UTC().Format(timeFormat),
		"end_time":   endTime.UTC().Format(timeFormat),
		"period":     fmt.Sprintf("%d", period),
	}

	monascaURL, URLerr := p.createMonascaAPIURL("v2.0/metrics/statistics", queryParameters, dimensions)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := p.callMonasca(monascaURL)
	if monascaErr != nil {
		return nil, monascaErr
	}
	statisticsResponse := models.StatisticsResponse{}
	err := json.Unmarshal(body, &statisticsResponse)
	if err != nil {
		return nil, err
	}

	return &statisticsResponse, nil
}

func (c *Client) callMonasca(monascaURL string) ([]byte, error) {

	req, err := http.NewRequest("GET", monascaURL, nil)

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
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: %d %s", resp.StatusCode, string([]byte(body)))
	}

	return body, nil
}

func (c *Client) createMonascaAPIURL(path string, queryParameters map[string]string, dimensions map[string]string) (string, error) {

	monascaURL, parseErr := url.Parse(c.baseURL)
	if parseErr != nil {
		return "", parseErr
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

	return monascaURL.String(), nil
}
