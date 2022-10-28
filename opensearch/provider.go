package opensearch

import (
	"context"
	"net/http"

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
		},
		ResourcesMap: map[string]*schema.Resource{
			"opensearch_saved_object": resourceSavedObjects(),
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
	SavedObjects *saved_objects.SavedObjectsProvider
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
			Summary:  "Failed to get AWS credentials via the default chain",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	cfg.RoundTripper = signer

	// init providers
	savedObjectsProvider := saved_objects.NewSavedObjectsProvider(cfg.BaseUrl, &http.Client{Transport: cfg.RoundTripper})

	// pass providers to the client
	client := &OpensearchDashboardsClient{
		SavedObjects: savedObjectsProvider,
	}

	return client, nil
}
