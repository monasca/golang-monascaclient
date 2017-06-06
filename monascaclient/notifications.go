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

func GetNotificationMethods(notificationQuery *models.NotificationQuery) (*models.NotificationResponse, error){
	return monClient.GetNotificationMethods(notificationQuery)
}

func GetNotificationMethod(notificationMethodID string, notificationQuery *models.NotificationQuery) (*models.NotificationElement, error) {
	return monClient.GetNotificationMethod(notificationMethodID, notificationQuery)
}

func (c *Client) GetNotificationMethods(notificationQuery *models.NotificationQuery) (*models.NotificationResponse, error) {
	urlValues := convertStructToQueryParameters(notificationQuery)

	monascaURL, URLerr := c.createMonascaAPIURL("v2.0/notification-methods", urlValues)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := c.callMonasca(monascaURL,"GET", nil)
	if monascaErr != nil {
		return nil, monascaErr
	}

	notificationMethodsResponse := models.NotificationResponse{}
	err := json.Unmarshal(body, &notificationMethodsResponse)
	if err != nil {
		return nil, err
	}

	return &notificationMethodsResponse, nil
}

func (c *Client) GetNotificationMethod(notificationMethodID string, notificationQuery *models.NotificationQuery) (*models.NotificationElement, error) {
	urlValues := convertStructToQueryParameters(notificationQuery)

	path := "v2.0/notification-methods/" + notificationMethodID

	monascaURL, URLerr := c.createMonascaAPIURL(path, urlValues)
	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := c.callMonasca(monascaURL,"GET", nil)
	if monascaErr != nil {
		return nil, monascaErr
	}

	notificationMethodElement := models.NotificationElement{}
	err := json.Unmarshal(body, &notificationMethodElement)
	if err != nil {
		return nil, err
	}

	return &notificationMethodElement, nil
}

func (c *Client) CreateNotificationMethod(name *string, notificationType *string, address *string,
	period *int) (*models.NotificationElement, error) {
	notificationMethod := models.Notification{}
	if name != nil {
		notificationMethod.Name = *name
	}
	if notificationType != nil {
		notificationMethod.Type = *notificationType
	}
	if address != nil {
		notificationMethod.Address = *address
	}
	if period != nil {
		notificationMethod.Period = *period
	}

	path := "v2.0/notification-methods"

	monascaURL, URLerr := c.createMonascaAPIURL(path, nil)
	if URLerr != nil {
		return nil, URLerr
	}

	byteInput, marshalErr  := json.Marshal(notificationMethod)
	if marshalErr != nil{
		return nil, marshalErr
	}
	body, monascaErr := c.callMonasca(monascaURL,"POST", &byteInput)
	if monascaErr != nil {
		return nil, monascaErr
	}

	notificationMethodElement := models.NotificationElement{}
	err := json.Unmarshal(body, &notificationMethodElement)
	if err != nil {
		return nil, err
	}

	return &notificationMethodElement, nil
}

