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

func GetAlarms(alarmQuery *models.AlarmQuery) (*models.AlarmsResponse, error) {
	return monClient.GetAlarms(alarmQuery)
}

func GetAlarm(alarmID string) (*models.Alarm, error) {
	return monClient.GetAlarm(alarmID)
}

func UpdateAlarm(alarmID string, alarmRequestBody *models.AlarmRequestBody) (*models.Alarm, error) {
	return monClient.UpdateAlarm(alarmID, alarmRequestBody)
}

func PatchAlarm(alarmID string, alarmRequestBody *models.AlarmRequestBody) (*models.Alarm, error) {
	return monClient.PatchAlarm(alarmID, alarmRequestBody)
}

func (c *Client) GetAlarms(alarmQuery *models.AlarmQuery) (*models.AlarmsResponse, error) {
	urlValues := convertStructToQueryParameters(alarmQuery)


	monascaURL, URLerr := c.createMonascaAPIURL("v2.0/alarms", urlValues)

	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := c.callMonasca(monascaURL, "GET", nil)
	if monascaErr != nil {
		return nil, monascaErr
	}

	alarmsResponse := models.AlarmsResponse{}
	err := json.Unmarshal(body, &alarmsResponse)
	if err != nil {
		return nil, err
	}

	return &alarmsResponse, nil
}

func (c *Client) GetAlarm(alarmID string) (*models.Alarm, error) {
	path := "v2.0/alarms/" + alarmID

	monascaURL, URLerr := c.createMonascaAPIURL(path, nil)

	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := c.callMonasca(monascaURL, "GET", nil)
	if monascaErr != nil {
		return nil, monascaErr
	}

	alarm := models.Alarm{}
	err := json.Unmarshal(body, &alarm)
	if err != nil {
		return nil, err
	}

	return &alarm, nil
}

func (c *Client) UpdateAlarm(alarmID string, alarmRequestBody *models.AlarmRequestBody) (*models.Alarm, error) {
	path := "v2.0/alarms/" + alarmID

	monascaURL, URLerr := c.createMonascaAPIURL(path, nil)

	if URLerr != nil {
		return nil, URLerr
	}

	byteInput, marshalErr  := json.Marshal(*alarmRequestBody)
	if marshalErr != nil{
		return nil, marshalErr
	}

	body, monascaErr := c.callMonasca(monascaURL, "PUT", &byteInput)
	if monascaErr != nil {
		return nil, monascaErr
	}

	alarm := models.Alarm{}
	err := json.Unmarshal(body, &alarm)
	if err != nil {
		return nil, err
	}

	return &alarm, nil
}

func (c *Client) PatchAlarm(alarmID string, alarmRequestBody *models.AlarmRequestBody) (*models.Alarm, error) {
	path := "v2.0/alarms/" + alarmID

	monascaURL, URLerr := c.createMonascaAPIURL(path, nil)

	if URLerr != nil {
		return nil, URLerr
	}

	byteInput, marshalErr  := json.Marshal(*alarmRequestBody)
	if marshalErr != nil{
		return nil, marshalErr
	}

	body, monascaErr := c.callMonasca(monascaURL, "PATCH", &byteInput)
	if monascaErr != nil {
		return nil, monascaErr
	}

	alarm := models.Alarm{}
	err := json.Unmarshal(body, &alarm)
	if err != nil {
		return nil, err
	}

	return &alarm, nil
}