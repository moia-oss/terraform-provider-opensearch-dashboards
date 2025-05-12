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
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/moia-oss/terraform-provider-opensearch-dashboards/pkg/saved_objects"
	"github.com/rs/zerolog/log"
)

func resourceSavedObjects() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages saved objects in OpenSearch Dashboards.",
		ReadContext:   resourceSavedObjectRead,
		CreateContext: resourceSavedObjectWrite,
		UpdateContext: resourceSavedObjectWrite,
		DeleteContext: resourceSavedObjectsDelete,
		Schema: map[string]*schema.Schema{
			"obj_id": {
				Description: "ID of the saved object.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description: "Type of the saved object.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"attributes": {
				// This will contain stringified JSON. We'll just send the content to the OpenSearch Dashboards API
				Description: "Attributes of the saved object.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"references": {
				Description: "References of the saved object.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
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

		for _, rAny := range refsAny.(*schema.Set).List() {
			rMap := rAny.(map[string]any)

			refs = append(refs, saved_objects.Reference{
				ID:   rMap["id"].(string),
				Name: rMap["name"].(string),
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
		a["name"] = resp.References[i].Name
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
