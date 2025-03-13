package opensearch

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog/log"
)

func resourceDefaultIndexPattern() *schema.Resource {
	return &schema.Resource{
		ReadContext:   defaultIndexPatternRead,
		CreateContext: defaultIndexPatternWrite,
		UpdateContext: defaultIndexPatternWrite,
		DeleteContext: defaultIndexPatternDelete,
		Schema: map[string]*schema.Schema{
			"index_pattern_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func defaultIndexPatternDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	hc, isAssertedType := m.(*OpensearchDashboardsClient)

	if !isAssertedType {
		return diag.Errorf("unexpected type provided as client: %T", m)
	}

	diagnostics := hc.DefaultIndexPattern.SetDefaultIndexPattern(ctx, nil)
	if diagnostics != nil {
		log.Error().Msgf("could not remove default index pattern. Terraform diagnostics: %v", diagnostics)

		return diagnostics
	}

	d.SetId("default-pattern")

	return nil
}

func defaultIndexPatternRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	hc, isAssertedType := m.(*OpensearchDashboardsClient)
	if !isAssertedType {
		return diag.Errorf("unexpected type provided as client: %T", m)
	}

	resp, diagnostics := hc.DefaultIndexPattern.GetDefaultIndexPattern(ctx)
	if diagnostics != nil {
		return diagnostics
	}
	err := d.Set("index_pattern_id", resp.IndexPatternId)
	if err != nil {
		return diag.Errorf("could not read index_pattern_id after fetching from api: %v+", err)
	}

	d.SetId("default-pattern")

	return nil
}

func defaultIndexPatternWrite(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	hc, isAssertedType := m.(*OpensearchDashboardsClient)

	if !isAssertedType {
		return diag.Errorf("unexpected type provided as client: %T", m)
	}

	patternId := d.Get("index_pattern_id").(string)
	diagnostics := hc.DefaultIndexPattern.SetDefaultIndexPattern(ctx, &patternId)
	if diagnostics != nil {
		log.Error().Msgf("could not set default index pattern. Terraform diagnostics: %v", diagnostics)

		return diagnostics
	}

	d.SetId("default-pattern")

	return nil
}
