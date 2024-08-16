/*
Use this data source to query Credential Type by ID or name.

# Example Usage

```hcl

	data "awx_credential_type" "project" {
	  name = "Project Credentials"
	}

```
*/
package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/mrcrilly/goawx/client"
)

func dataSourceCredentialType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inputs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"injectors": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCredentialTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	params := make(map[string]string)

	if name, ok := d.GetOk("name"); ok {
		params["name"] = name.(string)
	}

	if id, ok := d.GetOk("id"); ok {
		params["id"] = strconv.Itoa(id.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or id)",
		)
	}
	ctypes, _, err := client.CredentialTypeService.ListCredentialTypes(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch credential types",
			"Fail to find the credential type got: %s",
			err.Error(),
		)
	}
	if len(ctypes) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one credential type, %d",
			len(ctypes),
		)
	}
	credType := ctypes[0]

	d.Set("name", credType.Name)
	d.Set("description", credType.Description)
	d.Set("kind", credType.Kind)
	d.Set("inputs", credType.Inputs)
	d.Set("injectors", credType.Injectors)
	d.SetId(strconv.Itoa(credType.ID))
	return diags
}
