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
	"fmt"
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
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Usually in index-patterns the fields are automatically generated from the matched indices. If you instead explicitly want to track index-pattern-fields with terraform, set this value to true.",
			},
			"disable_authentication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "In all production environments, authentication is expected but with this flag it can be disabled for example for the purpose of local testing",
			},
			"path_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/_dashboards",
				Description: "prefix to be prepended to any path. The default is '/_dashboards' to prevent breaking change since this is needed for AWS Opensearch on which this provider was first used. You will want to set this to an empty string for development on a local Opensearch for example",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"opensearch_saved_object":          resourceSavedObjects(),
			"opensearch_default_index_pattern": resourceDefaultIndexPattern(),
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
		BaseUrl:      d.Get("base_url").(string) + d.Get("path_prefix").(string),
		RoundTripper: http.DefaultTransport,
	}

	var disableAuthentication bool
	if v, ok := d.GetOk("disable_authentication"); ok {
		disableAuthentication = v.(bool)
	}

	signer, err := getRoundTripper(disableAuthentication)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not create round-tripper for request-signing: " + err.Error(),
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

func getRoundTripper(disableAuthentication bool) (http.RoundTripper, error) {
	if disableAuthentication {
		return sigv4.RoundTripperFunc(http.DefaultTransport.RoundTrip), nil
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS-credentials: %w", err)
	}

	signer, err := sigv4.NewSigner(
		&sigv4.Config{
			Service: "es",
			Region:  *sess.Config.Region,
		},
		sess.Config.Credentials,
		http.DefaultTransport)
	if err != nil {
		return nil, fmt.Errorf("could not create sigv4 Http request signer: %w", err)
	}

	return signer, nil
}
