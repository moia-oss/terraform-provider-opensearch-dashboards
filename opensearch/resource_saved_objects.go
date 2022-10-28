package opensearch

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/moia-oss/terraform-provider-opensearch-dashboards/pkg/saved_objects"
	"github.com/rs/zerolog/log"
)

func resourceSavedObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSavedObjectRead,
		CreateContext: resourceSavedObjectWrite,
		UpdateContext: resourceSavedObjectWrite,
		DeleteContext: resourceSavedObjectsDelete,
		Schema: map[string]*schema.Schema{
			"obj_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attributes": {
				// This will contain stringified JSON. We'll just send the content to the OpenSearch Dashboards API
				Type:     schema.TypeString,
				Required: true,
			},
			"references": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSavedObjectsToRequest(resource *schema.ResourceData) (*saved_objects.SavedObjectOSD, diag.Diagnostics) {
	objId, ok := resource.GetOk("obj_id")
	if !ok {
		return nil, diag.Errorf("failed to convert schema to request struct, obj_id field missing: %v", resource)
	}

	objType, ok := resource.GetOk("type")
	if !ok {
		return nil, diag.Errorf("failed to convert schema to request struct, type field missing: %v", resource)
	}

	objAttr, ok := resource.GetOk("attributes")
	if !ok {
		return nil, diag.Errorf("failed to convert schema to request struct, attribute field missing: %v", resource)
	}

	result := &saved_objects.SavedObjectOSD{}

	resource.SetId(objId.(string))
	result.ID = objId.(string)
	result.Type = objType.(string)

	attrMap := make(map[string]any)
	attrStr := objAttr.(string)
	err := json.Unmarshal([]byte(attrStr), &attrMap)
	if err != nil {
		return nil, diag.Errorf("attributes is not valid json: %v", err)
	}
	result.Attributes = attrMap

	if refsAny, ok := resource.GetOk("references"); ok {
		var refs []saved_objects.Reference

		for _, rAny := range refsAny.([]any) {
			rMap := rAny.(map[string]any)

			refs = append(refs, saved_objects.Reference{
				ID:   rMap["id"].(string),
				Type: rMap["type"].(string),
			})
		}
		result.References = refs
	}

	return result, nil
}

func resourceSavedObjectsDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	hc, isAssertedType := m.(*OpensearchDashboardsClient)
	if !isAssertedType {
		return diag.Errorf("unexpected type provided as client: %T", m)
	}

	req, diag := resourceSavedObjectsToRequest(d)
	if diag != nil {
		return diag
	}

	diagnostics := hc.SavedObjects.DeleteObject(ctx, req)
	if diagnostics != nil {
		log.Error().Msgf("could not get saved object. Terraform diagnostics: %v", diagnostics)

		return diagnostics
	}

	return nil
}

func resourceSavedObjectRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	hc, isAssertedType := m.(*OpensearchDashboardsClient)
	if !isAssertedType {
		return diag.Errorf("unexpected type provided as client: %T", m)
	}

	obj, diagnostics := resourceSavedObjectsToRequest(d)
	if diagnostics != nil {
		return diagnostics
	}

	resp, diagnostics := hc.SavedObjects.GetObject(ctx, obj)
	if diagnostics != nil {
		return diagnostics
	}

	if resp == nil {
		// signals the resource must be (re)created
		d.SetId("")
		return diagnostics
	}

	err := d.Set("type", resp.Type)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("attributes", resp.Attributes)
	if err != nil {
		return diag.FromErr(err)
	}

	refs := make([]any, len(resp.References))
	for i := range resp.References {
		a := make(map[string]any)
		a["id"] = resp.References[i].ID
		a["type"] = resp.References[i].Type
		refs[i] = a
	}

	err = d.Set("references", refs)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return nil
}

func resourceSavedObjectWrite(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	hc, isAssertedType := m.(*OpensearchDashboardsClient)
	if !isAssertedType {
		return diag.Errorf("unexpected type provided as client: %T", m)
	}

	req, diagnostics := resourceSavedObjectsToRequest(d)
	if diagnostics != nil {
		return diagnostics
	}

	diagnostics = hc.SavedObjects.SaveObject(ctx, req)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}
