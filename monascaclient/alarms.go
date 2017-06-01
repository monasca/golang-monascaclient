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
	"time"
)

func GetAlarms(alarmDefinitionID *string, metricName *string, state *string, severity *string, lifecycleState *string,
	link *string, stateUpdatedStartTime *time.Time, sortBy *string,
	metricDimensions map[string]string) (*models.AlarmsResponse, error) {
	return monClient.GetAlarms(alarmDefinitionID, metricName, state, severity, lifecycleState, link,
		stateUpdatedStartTime, sortBy, metricDimensions)
}
func (p *Client) GetAlarms(alarmDefinitionID *string, metricName *string, state *string, severity *string,
	lifecycleState *string, link *string, stateUpdatedStartTime *time.Time, sortBy *string,
	metricDimensions map[string]string) (*models.AlarmsResponse, error) {
	queryParameters := map[string]string{}
	if alarmDefinitionID != nil {
		queryParameters["alarm_definition_id"] = *alarmDefinitionID
	}
	if metricName != nil {
		queryParameters["metric_name"] = *metricName
	}
	if state != nil {
		queryParameters["state"] = *state
	}
	if severity != nil {
		queryParameters["severity"] = *severity
	}
	if lifecycleState != nil {
		queryParameters["lifecycle_state"] = *lifecycleState
	}
	if state != nil {
		queryParameters["state_updated_start_time"] = *state
	}
	if sortBy != nil {
		queryParameters["sort_by"] = *sortBy
	}
	if link != nil {
		queryParameters["link"] = *link
	}
	if stateUpdatedStartTime != nil {
		queryParameters["state_updated_start_time"] = (*stateUpdatedStartTime).UTC().Format(timeFormat)
	}

	monascaURL, URLerr := p.createMonascaAPIURL("v2.0/alarms", queryParameters, metricDimensions)

	if URLerr != nil {
		return nil, URLerr
	}

	body, monascaErr := p.callMonasca(monascaURL)
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
