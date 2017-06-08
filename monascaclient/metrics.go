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
	"encoding/json"
	"github.com/monasca/golang-monascaclient/monascaclient/models"
	"net/url"
)

const timeFormat = "2006-01-02T15:04:05Z"

func GetMetrics(metricQuery *models.MetricQuery) ([]models.Metric, error) {
	return monClient.GetMetrics(metricQuery)
}

func GetDimensionValues(dimensionQuery *models.DimensionValueQuery) ([]string, error) {
	return monClient.GetDimensionValues(dimensionQuery)
}

func GetDimensionNames(dimensionQuery *models.DimensionNameQuery) ([]string, error) {
	return monClient.GetDimensionNames(dimensionQuery)
}

func GetStatistics(statisticsQuery *models.StatisticQuery) (*models.StatisticsResponse, error) {
	return monClient.GetStatistics(statisticsQuery)
}

func GetMeasurements(measurementQuery *models.MeasurementQuery) (*models.MeasurementsResponse, error) {
	return monClient.GetMeasurements(measurementQuery)
}

func CreateMetric(tenantID *string, metricRequestBody *models.MetricRequestBody) (error) {
	return monClient.CreateMetric(tenantID, metricRequestBody)
}


func (c *Client) CreateMetric(tenantID *string, metricRequestBody *models.MetricRequestBody) (error) {
	urlValues := url.Values{}
	if tenantID != nil {
		urlValues.Add("tenant_id", *tenantID)
	}
	monascaURL, URLerr := c.createMonascaAPIURL("v2.0/metrics", urlValues)
	if URLerr != nil {
		return URLerr
	}
	byteInput, marshalErr  := json.Marshal(*metricRequestBody)
	if marshalErr != nil{
		return marshalErr
	}
	return c.callMonascaNoContent(monascaURL, "POST", &byteInput)
}

func (c *Client) GetMetrics(metricQuery *models.MetricQuery) ([]models.Metric, error) {
	urlValues := convertStructToQueryParameters(metricQuery)

	monascaURL, URLerr := c.createMonascaAPIURL("v2.0/metrics", urlValues)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := c.callMonascaReturnBody(monascaURL, "GET", nil)
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

func (c *Client) GetDimensionValues(dimensionQuery *models.DimensionValueQuery) ([]string, error) {
	urlValues := convertStructToQueryParameters(dimensionQuery)

	monascaURL, err := c.createMonascaAPIURL("/v2.0/metrics/dimensions/names/values", urlValues)
	if err != nil {
		return nil, err
	}

	body, monascaErr := c.callMonascaReturnBody(monascaURL, "GET", nil)
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

func (c *Client) GetDimensionNames(dimensionQuery *models.DimensionNameQuery) ([]string, error) {
	urlValues := convertStructToQueryParameters(dimensionQuery)

	monascaURL, err := c.createMonascaAPIURL("/v2.0/metrics/dimensions/names", urlValues)
	if err != nil {
		return nil, err
	}

	body, monascaErr := c.callMonascaReturnBody(monascaURL, "GET", nil)
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

func (c *Client) GetStatistics(statisticsQuery *models.StatisticQuery) (*models.StatisticsResponse, error) {
	urlValues := convertStructToQueryParameters(statisticsQuery)

	monascaURL, URLerr := c.createMonascaAPIURL("v2.0/metrics/statistics", urlValues)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := c.callMonascaReturnBody(monascaURL, "GET", nil)
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

func (c *Client) GetMeasurements(measurementQuery *models.MeasurementQuery) (*models.MeasurementsResponse, error) {
	urlValues := convertStructToQueryParameters(measurementQuery)

	monascaURL, URLerr := c.createMonascaAPIURL("v2.0/metrics/measurements", urlValues)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := c.callMonascaReturnBody(monascaURL, "GET", nil)
	if monascaErr != nil {
		return nil, monascaErr
	}
	measurementResponse := models.MeasurementsResponse{}
	err := json.Unmarshal(body, &measurementResponse)
	if err != nil {
		return nil, err
	}

	return &measurementResponse, nil
}
