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
	"fmt"
	"encoding/json"
	"time"
	"strconv"
	"github.com/monasca/golang-monascaclient/monascaclient/models"
)

const timeFormat = "2006-01-02T15:04:05Z"

func GetMetrics(metricName *string, dimensions map[string]string) ([]models.Metric, error) {
	return monClient.GetMetrics(metricName, dimensions)
}

func GetStatistics(metricName *string, statisticFunction *string, startTime *time.Time, endTime *time.Time,
	period *int64, dimensions map[string]string) (*models.StatisticsResponse, error) {
	return monClient.GetStatistics(metricName, statisticFunction, startTime, endTime, period, dimensions)
}

func GetMeasurements(metricName *string, startTime *time.Time, endTime *time.Time, mergeMetrics *bool, groupBy *string,
	dimensions map[string]string) (*models.MeasurementsResponse, error) {
	return monClient.GetMeasurements(metricName, startTime, endTime, mergeMetrics, groupBy, dimensions)
}

func (p *Client) GetMetrics(metricName *string, dimensions map[string]string) ([]models.Metric, error) {
	queryParameters := map[string]string{}
	if metricName != nil {
		queryParameters["name"] = *metricName
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

func (p *Client) GetDimensionValues(metricName, dimensionName *string) ([]string, error) {
	queryParameters := map[string]string{}
	if metricName != nil {
		queryParameters["metric_name"] = *metricName
	}
	if dimensionName != nil {
		queryParameters["dimension_name"] = *dimensionName
	}

	monascaURL, err := p.createMonascaAPIURL("/v2.0/metrics/dimensions/names/values", queryParameters, nil)
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

func (p *Client) GetStatistics(metricName *string, statisticFunction *string, startTime *time.Time, endTime *time.Time,
	period *int64, dimensions map[string]string) (*models.StatisticsResponse, error) {
	// TODO: Review this entire function
	// TODO: Handle errors. How does Go work?
	queryParameters := map[string]string{}
	if metricName != nil {
		queryParameters["name"] = *metricName
	}
	if statisticFunction != nil {
		queryParameters["statistics"] = *statisticFunction
	}
	if startTime != nil {
		queryParameters["start_time"] = (*startTime).UTC().Format(timeFormat)
	}
	if endTime != nil {
		queryParameters["end_time"] = (*endTime).UTC().Format(timeFormat)
	}
	if period != nil {
		queryParameters["period"] = fmt.Sprintf("%d", *period)
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

func (p *Client) GetMeasurements(metricName *string, startTime *time.Time, endTime *time.Time, mergeMetrics *bool,
	groupBy *string, dimensions map[string]string) (*models.MeasurementsResponse, error) {
	queryParameters := map[string]string{}
	if metricName != nil {
		queryParameters["name"] = *metricName
	}
	if startTime != nil {
		queryParameters["start_time"] = (*startTime).UTC().Format(timeFormat)
	}
	if endTime != nil {
		queryParameters["end_time"] = (*endTime).UTC().Format(timeFormat)
	}
	if mergeMetrics != nil {
		queryParameters["merge_metrics"] = strconv.FormatBool(*mergeMetrics)
	}
	if groupBy != nil {
		queryParameters["group_by"] = *groupBy
	}

	monascaURL, URLerr := p.createMonascaAPIURL("v2.0/metrics/measurements", queryParameters, dimensions)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := p.callMonasca(monascaURL)
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