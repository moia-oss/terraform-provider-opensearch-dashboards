package default_index_pattern

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type OpenSearchRequestBody struct {
	IndexPatternId *string `json:"index_pattern_id"`
}

type httpPayload struct {
	Changes httpPayloadChanges `json:"changes"`
}

type httpPayloadChanges struct {
	DefaultIndex *string `json:"defaultIndex"`
}

type Provider struct {
	Url                    string
	httpClient             *http.Client
	SyncIndexPatternFields bool
}

func NewProvider(baseUrl string, client *http.Client) *Provider {
	return &Provider{
		Url:        fmt.Sprintf("%s/_dashboards/api/opensearch-dashboards/settings", baseUrl),
		httpClient: client,
	}
}

func (p *Provider) GetDefaultIndexPattern(ctx context.Context) (*OpenSearchRequestBody, diag.Diagnostics) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.Url, nil)
	req.Header.Set("osd-xsrf", "true")
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("could not build request to GET %v %w", p.Url, err))
	}

	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("GET '%v' failed: %w", req.URL.String(), err))
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if res.StatusCode != http.StatusOK {
		response, bodyReadErr := io.ReadAll(res.Body)
		if bodyReadErr != nil {
			return nil, diag.FromErr(fmt.Errorf("GET '%v' failed with status %d", req.URL.String(), res.StatusCode))
		}
		return nil, diag.FromErr(fmt.Errorf("GET '%v' failed with status %d\nresponse_body: %v", req.URL.String(), res.StatusCode, string(response)))
	}

	result := &httpPayload{}
	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("request failed, cannot decode response body, err %w ", err))
	}

	return &OpenSearchRequestBody{IndexPatternId: result.Changes.DefaultIndex}, nil
}

func (p *Provider) SetDefaultIndexPattern(ctx context.Context, indexPatternId *string) diag.Diagnostics {
	requestBody := httpPayload{Changes: httpPayloadChanges{DefaultIndex: indexPatternId}}
	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		return diag.Errorf("failed to encode index pattern as JSON: %+v \n%v", requestBody, err)
	}
	payload := bytes.NewBuffer(jsonBytes)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.Url, payload)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("osd-xsrf", "true")
	req.Header.Set("Content-Type", "application/json")

	res, err := p.httpClient.Do(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("POST '%s' failed, err %s, statuscode %d, \nrequest_body: %s", req.URL.String(), err, res.StatusCode, string(jsonBytes)))
	}
	if res.StatusCode != http.StatusOK {
		rawBody, bodyReadErr := io.ReadAll(res.Body)
		if bodyReadErr != nil {
			return diag.FromErr(fmt.Errorf("POST '%s' failed with status %d\nrequest_body: %s", req.URL.String(), err, string(jsonBytes)))
		}
		return diag.FromErr(fmt.Errorf("POST '%s' failed with status %d\nrequest_body: %s\nresponse_body: %s", req.URL.String(), res.StatusCode, string(jsonBytes), string(rawBody)))
	}
	return nil
}
