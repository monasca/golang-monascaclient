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

package models

import "time"

type AlarmsResponse struct {
	Links    []Link           `json:"links"`
	Elements []Alarm 	  `json:"elements"`
}

type Alarm struct {
	Metrics 	      []Metric  `json:"metrics"`
	LifecycleState 	      string    `json:"lifecycle_state"`
	State 	    	      string    `json:"state"`
	Link		      string    `json:"link"`
	ID	 	      string    `json:"id"`
	UpdatedTimestamp      time.Time `json:"updated_timestamp"`
	CreatedTimestamp      time.Time `json:"created_timestamp"`
	StateUpdatedTimestamp time.Time `json:"state_updated_timestamp"`
	Links    	      []Link    		`json:"links"`
	//AlarmDefinition	      []AlarmDefinitionSummary  `json:"links"`
}
