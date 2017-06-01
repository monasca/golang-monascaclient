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
	"github.com/monasca/golang-monascaclient/monascaclient/models"
	"encoding/json"
)

func GetAlarmDefinitions(alarmName *string, severity *string,
	dimensions map[string]string) (*models.AlarmDefinitionsResponse, error) {
	return monClient.GetAlarmDefinitions(alarmName, severity, dimensions)

}

func (p *Client) GetAlarmDefinitions(alarmName *string, severity *string,
	dimensions map[string]string) (*models.AlarmDefinitionsResponse, error) {
	queryParameters := map[string]string{}
	if alarmName != nil {
		queryParameters["name"] = *alarmName
	}
	if severity != nil {
		queryParameters["severity"] = *severity
	}
	monascaURL, URLerr := p.createMonascaAPIURL("v2.0/alarm-definitions", queryParameters, dimensions)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := p.callMonasca(monascaURL)
	if monascaErr != nil {
		return nil, monascaErr
	}

	alarmDefinitionsResponse := models.AlarmDefinitionsResponse{}
	err := json.Unmarshal(body, &alarmDefinitionsResponse)
	if err != nil {
		return nil, err
	}

	return &alarmDefinitionsResponse, nil
}
//name (string(255), optional) - Name of alarm to filter by.
//dimensions (string, optional) - Dimensions of metrics to filter by specified as a comma separated array of (key, value) pairs as key1:value1,key1:value1, ..., leaving the value empty key1,key2:value2 will return all values for that key, multiple values for a key may be specified as key1:value1|value2|...,key2:value4,...
//severity (string, optional) - One or more severities to filter by, separated with |, ex. severity=LOW|MEDIUM.
//offset (integer, optional)
//limit (integer, optional)
//sort_by (string, optional) - Comma separated list of fields to sort by, defaults to 'id', 'created_at'. Fields may be followed by 'asc' or 'desc' to set the direction, ex 'severity desc' Allowed fields for sort_by are: 'id', 'name', 'severity', 'updated_at', 'created_at'
//alarm-definition-create  Create an alarm definition.
//alarm-definition-delete  Delete the alarm definition.
//alarm-definition-list    List alarm definitions for this tenant.
//alarm-definition-patch   Patch the alarm definition.
//alarm-definition-show    Describe the alarm definition.
//alarm-definition-update  Update the alarm definition.