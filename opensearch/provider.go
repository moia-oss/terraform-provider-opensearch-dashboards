package opensearch

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
	"context"
	"net/http"

	"github.com/moia-oss/terraform-provider-opensearch-dashboards/pkg/default_index_pattern"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/moia-oss/terraform-provider-opensearch-dashboards/pkg/saved_objects"
	"github.com/moia-oss/terraform-provider-opensearch-dashboards/pkg/sigv4"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_BASE_URL", nil),
				Description: "The Opensearch base url",
			},
			"sync_index_pattern_fields": {
				Type:      schema.TypeBool,
				Optional:  true,
				Sensitive: false,
				Default:   false,
				Description: "Usually in index-patterns the fields are automatically generated from the matched indices. " +
					"If you instead explicitly want to track index-pattern-fields with terraform, set this value to true.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"opensearch_saved_object": resourceSavedObjects(),
			"default_index_pattern":   resourceDefaultIndexPattern(),
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		return providerConfigure(ctx, d)
	}
	return p
}

type ProviderConfig struct {
	// public settings
	BaseUrl string

	// internal
	RoundTripper http.RoundTripper
}

type OpensearchDashboardsClient struct {
	SavedObjects        *saved_objects.SavedObjectsProvider
	DefaultIndexPattern *default_index_pattern.Provider
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	var diags diag.Diagnostics

	cfg := &ProviderConfig{
		BaseUrl:      d.Get("base_url").(string),
		RoundTripper: http.DefaultTransport,
	}

	sess, err := session.NewSession()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to get AWS credentials via the default chain",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	signer, err := sigv4.NewSigner(
		&sigv4.Config{
			Service: "es",
			Region:  *sess.Config.Region,
		},
		sess.Config.Credentials,
		http.DefaultTransport)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not create sigv4 HTTP request signer",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	cfg.RoundTripper = signer

	var syncIndexPatternFields bool
	if v, ok := d.GetOk("sync_index_pattern_fields"); ok {
		syncIndexPatternFields = v.(bool)
	}

	// init providers
	savedObjectsProvider := saved_objects.NewSavedObjectsProvider(cfg.BaseUrl, &http.Client{Transport: cfg.RoundTripper}, syncIndexPatternFields)
	defaultIndexPatternProvider := default_index_pattern.NewProvider(cfg.BaseUrl, &http.Client{Transport: cfg.RoundTripper})

	// pass providers to the client
	client := &OpensearchDashboardsClient{
		SavedObjects:        savedObjectsProvider,
		DefaultIndexPattern: defaultIndexPatternProvider,
	}

	return client, nil
}
