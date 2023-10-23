/*
*TBD*

# Example Usage

```hcl

	data "awx_execution_environment" "default" {
	  name = "Default"
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

func dataSourceExecutionEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExecutionEnvironmentsRead,
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
		},
	}
}

func dataSourceExecutionEnvironmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or group_id)",
		)
	}
	ExecutionEnvironments, _, err := client.ExecutionEnvironmentService.ListExecutionEnvironments(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Execution Environments",
			"Fail to find the execution environment got: %s",
			err.Error(),
		)
	}
	if len(ExecutionEnvironments) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one execution environment, %d",
			len(ExecutionEnvironments),
		)
	}

	ExecutionEnvironment := ExecutionEnvironments[0]
	d = setExecutionEnvironmentResourceData(d, ExecutionEnvironment)
	return diags
}

func setExecutionEnvironmentResourceData(d *schema.ResourceData, r *awx.ExecutionEnvironment) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("organization", r.Organization)
	d.Set("image", r.Image)
	d.Set("managed", r.Managed)
	d.Set("credential", r.Credential)
	d.Set("pull", r.Pull)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
