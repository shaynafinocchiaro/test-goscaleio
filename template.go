// Copyright © 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goscaleio

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// GetTemplateByID gets the node details based on ID
func (gc *GatewayClient) GetTemplateByID(id string) (*types.TemplateDetails, error) {
	defer TimeSpent("GetTemplateByID", time.Now())

	path := fmt.Sprintf("/Api/V1/template/%v", id)

	var template types.TemplateDetails
	req, httpError := http.NewRequest("GET", gc.host+path, nil)
	if httpError != nil {
		return nil, httpError
	}

	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return nil, httpRespError
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Template not found")
	}

	responseString, _ := extractString(httpResp)
	err := json.Unmarshal([]byte(responseString), &template)
	if err != nil {
		return nil, fmt.Errorf("Error parsing response data for template: %s", err)
	}
	return &template, nil
}

// GetAllTemplates gets all the Template details
func (gc *GatewayClient) GetAllTemplates() ([]types.TemplateDetails, error) {
	defer TimeSpent("GetAllTemplates", time.Now())

	path := "/Api/V1/template"

	var templates types.TemplateDetailsFilter
	req, httpError := http.NewRequest("GET", gc.host+path, nil)
	if httpError != nil {
		return nil, httpError
	}

	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return nil, httpRespError
	}

	if httpResp.StatusCode == 200 {
		responseString, _ := extractString(httpResp)
		parseError := json.Unmarshal([]byte(responseString), &templates)

		if parseError != nil {
			return nil, fmt.Errorf("Error While Parsing Response Data For Template: %s", parseError)
		}
	}

	return templates.TemplateDetails, nil
}

// GetTemplateByFilters gets the Template details based on the provided filter
func (gc *GatewayClient) GetTemplateByFilters(key string, value string) ([]types.TemplateDetails, error) {
	defer TimeSpent("GetTemplateByFilters", time.Now())

	encodedValue := url.QueryEscape(value)

	path := `/Api/V1/template?filter=` + key + `%20eq%20%22` + encodedValue + `%22`

	var templates types.TemplateDetailsFilter
	req, httpError := http.NewRequest("GET", gc.host+path, nil)
	if httpError != nil {
		return nil, httpError
	}

	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}

	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return nil, httpRespError
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Template not found")
	}

	responseString, _ := extractString(httpResp)
	parseError := json.Unmarshal([]byte(responseString), &templates)
	if parseError != nil {
		return nil, fmt.Errorf("Error While Parsing Response Data For Template: %s", parseError)
	}

	if len(templates.TemplateDetails) == 0 {
		return nil, fmt.Errorf("Template not found")
	}

	return templates.TemplateDetails, nil
}
