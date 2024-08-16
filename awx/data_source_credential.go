/*
Use this data source to query for a Credential by Name or ID.

# Example Usage

```hcl

	data "awx_credential" "provisioning_credentials" {
	  name = "Provisioning Credentials"
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

func dataSourceCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialRead,
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
			"tower_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	creds, _, err := client.CredentialsService.ListCredentials(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch credentials",
			"Fail to find the credential got: %s",
			err.Error(),
		)
	}
	if len(creds) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one credential, %d",
			len(creds),
		)
	}
	cred := creds[0]

	d.Set("name", cred.Name)
	d.Set("username", cred.Inputs["username"])
	d.Set("kind", cred.Kind)
	d.Set("tower_id", cred.ID)
	d.SetId(strconv.Itoa(cred.ID))
	return diags
}
