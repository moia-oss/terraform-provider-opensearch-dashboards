package saved_objects

/*
Copyright 2022 MOIA GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type SavedObjectsProvider struct {
	BaseUrl    string
	httpClient *http.Client
}

func NewSavedObjectsProvider(baseUrl string, client *http.Client) *SavedObjectsProvider {
	return &SavedObjectsProvider{
		BaseUrl:    baseUrl,
		httpClient: client,
	}
}

func (p *SavedObjectsProvider) GetObject(ctx context.Context, obj *SavedObjectOSD) (*SavedObjectTF, diag.Diagnostics) {
	// build request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		p.URL(fmt.Sprintf("/%s/%s", obj.Type, obj.ID)),
		nil,
	)
	req.Header.Set("osd-xsrf", "true")
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, diag.FromErr(err)
	}

	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("request failed, err %w ", err))
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if res.StatusCode != http.StatusOK {
		return nil, diag.FromErr(fmt.Errorf("request failed, statuscode %d", err))
	}
	// parse result
	obj = &SavedObjectOSD{}
	// we need to stringify the attributes field here
	err = json.NewDecoder(res.Body).Decode(obj)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("request failed, cannot decode response body, err %w ", err))
	}

	stringifiedAttributes, err := json.Marshal(obj.Attributes)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("request failed, cannot stringify attibutes from response body, err %w ", err))
	}

	result := &SavedObjectTF{
		Type:       obj.Type,
		ID:         obj.ID,
		Attributes: string(stringifiedAttributes),
		References: obj.References,
	}

	return result, nil
}

func (p *SavedObjectsProvider) SaveObject(ctx context.Context, obj *SavedObjectOSD) diag.Diagnostics {
	url := p.URL(fmt.Sprintf("/%s/%s%s", obj.Type, obj.ID, "?overwrite=true"))

	jsonBytes, err := json.Marshal(obj.SavedObjectPostPayload)
	if err != nil {
		return diag.Errorf("failed to encode saved object as JSON: %v", err)
	}
	payload := bytes.NewBuffer(jsonBytes)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		payload,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("osd-xsrf", "true")
	req.Header.Set("Content-Type", "application/json")

	res, err := p.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("request failed, err %s, statuscode %d, url %s\npayload: %s", err, res.StatusCode, req.URL.String(), string(jsonBytes)))
	}
	return nil
}

func (p *SavedObjectsProvider) DeleteObject(ctx context.Context, obj *SavedObjectOSD) diag.Diagnostics {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		p.URL(fmt.Sprintf("/%s/%s", obj.Type, obj.ID)),
		nil,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("osd-xsrf", "true")
	req.Header.Set("Content-Type", "application/json")

	res, err := p.httpClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("request failed, err %s, statuscode %d", err, res.StatusCode))
	}

	return nil
}

func (provider *SavedObjectsProvider) URL(path string) string {
	base := fmt.Sprintf("%s/_dashboards/api/saved_objects%s", provider.BaseUrl, path)
	return base
}
