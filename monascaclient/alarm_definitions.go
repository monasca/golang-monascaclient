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

func GetAlarmDefinitions(alarmDefinitionQuery *models.AlarmDefinitionQuery) (*models.AlarmDefinitionsResponse, error) {
	return monClient.GetAlarmDefinitions(alarmDefinitionQuery)

}

func GetAlarmDefinition(alarmDefinitionID string) (*models.AlarmDefinitionElement, error) {
	return monClient.GetAlarmDefinition(alarmDefinitionID)
}

func CreateAlarmDefinition(alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody)(*models.AlarmDefinitionElement, error) {
	return monClient.CreateAlarmDefinition(alarmDefinitionRequestBody)
}

func UpdateAlarmDefinition(alarmDefinitionID string, alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody)(*models.AlarmDefinitionElement, error) {
	return monClient.UpdateAlarmDefinition(alarmDefinitionID, alarmDefinitionRequestBody)
}

func PatchAlarmDefinition(alarmDefinitionID string, alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody)(*models.AlarmDefinitionElement, error) {
	return monClient.PatchAlarmDefinition(alarmDefinitionID, alarmDefinitionRequestBody)
}

func DeleteAlarmDefinition(alarmDefinitionID string) (error) {
	return monClient.DeleteAlarmDefinition(alarmDefinitionID)
}

func (c *Client) GetAlarmDefinitions(alarmDefinitionQuery *models.AlarmDefinitionQuery) (*models.AlarmDefinitionsResponse, error) {
	urlValues := convertStructToQueryParameters(alarmDefinitionQuery)

	monascaURL, URLerr := c.createMonascaAPIURL("v2.0/alarm-definitions", urlValues)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := c.callMonascaReturnBody(monascaURL, "GET", nil)
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

func (c *Client) GetAlarmDefinition(alarmDefinitionID string) (*models.AlarmDefinitionElement, error) {
	path := "v2.0/alarm-definitions" + alarmDefinitionID

	monascaURL, URLerr := c.createMonascaAPIURL(path, nil)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := c.callMonascaReturnBody(monascaURL,"GET", nil)
	if monascaErr != nil {
		return nil, monascaErr
	}

	alarmDefinitionElement := models.AlarmDefinitionElement{}
	err := json.Unmarshal(body, &alarmDefinitionElement)
	if err != nil {
		return nil, err
	}

	return &alarmDefinitionElement, nil
}

func (c *Client) CreateAlarmDefinition(alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody)(*models.AlarmDefinitionElement, error) {
	path := "v2.0/alarm-definitions"

	monascaURL, URLerr := c.createMonascaAPIURL(path, nil)
	if URLerr != nil {
		return nil, URLerr
	}

	byteInput, marshalErr  := json.Marshal(*alarmDefinitionRequestBody)
	if marshalErr != nil{
		return nil, marshalErr
	}
	body, monascaErr := c.callMonascaReturnBody(monascaURL,"POST", &byteInput)
	if monascaErr != nil {
		return nil, monascaErr
	}

	alarmDefinitionsElement := models.AlarmDefinitionElement{}
	err := json.Unmarshal(body, &alarmDefinitionsElement)
	if err != nil {
		return nil, err
	}

	return &alarmDefinitionsElement, nil
}

func (c *Client) UpdateAlarmDefinition(alarmDefinitionID string, alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody)(*models.AlarmDefinitionElement, error) {
	path := "v2.0/alarm-definitions/" + alarmDefinitionID

	monascaURL, URLerr := c.createMonascaAPIURL(path, nil)
	if URLerr != nil {
		return nil, URLerr
	}

	byteInput, marshalErr  := json.Marshal(*alarmDefinitionRequestBody)
	if marshalErr != nil{
		return nil, marshalErr
	}
	body, monascaErr := c.callMonascaReturnBody(monascaURL,"PUT", &byteInput)
	if monascaErr != nil {
		return nil, monascaErr
	}

	alarmDefinitionsElement := models.AlarmDefinitionElement{}
	err := json.Unmarshal(body, &alarmDefinitionsElement)
	if err != nil {
		return nil, err
	}

	return &alarmDefinitionsElement, nil
}

func (c *Client) PatchAlarmDefinition(alarmDefinitionID string, alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody)(*models.AlarmDefinitionElement, error) {
	path := "v2.0/alarm-definitions/" + alarmDefinitionID

	monascaURL, URLerr := c.createMonascaAPIURL(path, nil)
	if URLerr != nil {
		return nil, URLerr
	}

	byteInput, marshalErr  := json.Marshal(*alarmDefinitionRequestBody)
	if marshalErr != nil{
		return nil, marshalErr
	}
	body, monascaErr := c.callMonascaReturnBody(monascaURL,"PATCH", &byteInput)
	if monascaErr != nil {
		return nil, monascaErr
	}

	alarmDefinitionsElement := models.AlarmDefinitionElement{}
	err := json.Unmarshal(body, &alarmDefinitionsElement)
	if err != nil {
		return nil, err
	}

	return &alarmDefinitionsElement, nil
}

func (c *Client) DeleteAlarmDefinition(alarmDefinitionID string) (error) {
	path := "v2.0/alarm-definitions/" + alarmDefinitionID

	monascaURL, URLerr := c.createMonascaAPIURL(path, nil)
	if URLerr != nil {
		return URLerr
	}

	return c.callMonascaNoContent(monascaURL,"DELETE", nil)
}