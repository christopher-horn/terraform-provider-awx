/*
*TBD*

# Example Usage

```hcl

	data "awx_inventory" "default" {
	  name            = "private_services"
	  organization_id = data.awx_organization.default.id
	}

	resource "awx_job_template" "baseconfig" {
	  name           = "baseconfig"
	  job_type       = "run"
	  inventory_id   = data.awx_inventory.default.id
	  project_id     = awx_project.base_service_config.id
	  playbook       = "master-configure-system.yml"
	  become_enabled = true
	}

```
*/
package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/mrcrilly/goawx/client"
)

func resourceJobTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobTemplateCreate,
		ReadContext:   resourceJobTemplateRead,
		UpdateContext: resourceJobTemplateUpdate,
		DeleteContext: resourceJobTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Run, Check, Scan
			"job_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of: run, check, scan",
			},
			"inventory_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"playbook": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scm_branch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forks": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "One of 0,1,2,3,4,5",
			},
			"extra_vars": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force_handlers": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"skip_tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_at_task": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"use_fact_cache": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"execution_environment_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"host_config_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ask_scm_branch_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_diff_mode_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_variables_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_limit_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_skip_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_job_type_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_verbosity_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_inventory_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_credential_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_execution_environment_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_labels_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_forks_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_job_slice_count_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_timeout_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ask_instance_groups_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"survey_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"become_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"diff_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"allow_simultaneous": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"custom_virtualenv": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_slice_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"webhook_service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"webhook_credential_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"prevent_instance_group_fallback": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"credential_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
		//
		//Timeouts: &schema.ResourceTimeout{
		//	Create: schema.DefaultTimeout(1 * time.Minute),
		//	Update: schema.DefaultTimeout(1 * time.Minute),
		//	Delete: schema.DefaultTimeout(1 * time.Minute),
		//},
	}
}

func resourceJobTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService

	params := map[string]interface{}{
		"name":      d.Get("name").(string),
		"job_type":  d.Get("job_type").(string),
		"inventory": AtoipOr(d.Get("inventory_id").(string), nil),
		"project":   d.Get("project_id").(int),
		"playbook":  d.Get("playbook").(string),
	}
	if p, ok := d.GetOk("description"); ok {
		params["description"] = p.(string)
	}
	if p, ok := d.GetOk("scm_branch"); ok {
		params["scm_branch"] = p.(string)
	}
	if p, ok := d.GetOk("forks"); ok {
		params["forks"] = p.(int)
	}
	if p, ok := d.GetOk("limit"); ok {
		params["limit"] = p.(string)
	}
	if p, ok := d.GetOk("verbosity"); ok {
		params["verbosity"] = p.(int)
	}
	if p, ok := d.GetOk("extra_vars"); ok {
		params["extra_vars"] = p.(string)
	}
	if p, ok := d.GetOk("job_tags"); ok {
		params["job_tags"] = p.(string)
	}
	if p, ok := d.GetOk("force_handlers"); ok {
		params["force_handlers"] = p.(bool)
	}
	if p, ok := d.GetOk("skip_tags"); ok {
		params["skip_tags"] = p.(string)
	}
	if p, ok := d.GetOk("start_at_task"); ok {
		params["start_at_task"] = p.(string)
	}
	if p, ok := d.GetOk("timeout"); ok {
		params["timeout"] = p.(int)
	}
	if p, ok := d.GetOk("use_fact_cache"); ok {
		params["use_fact_cache"] = p.(bool)
	}
	if p, ok := d.GetOk("organization_id"); ok {
		params["organization_id"] = p.(int)
	}
	if p, ok := d.GetOk("execution_environment_id"); ok {
		params["execution_environment_id"] = p.(int)
	}
	if p, ok := d.GetOk("host_config_key"); ok {
		params["host_config_key"] = p.(string)
	}
	if p, ok := d.GetOk("ask_scm_branch_on_launch"); ok {
		params["ask_scm_branch_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_diff_mode_on_launch"); ok {
		params["ask_diff_mode_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_variables_on_launch"); ok {
		params["ask_variables_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_limit_on_launch"); ok {
		params["ask_limit_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_tags_on_launch"); ok {
		params["ask_tags_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_skip_tags_on_launch"); ok {
		params["ask_skip_tags_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_job_type_on_launch"); ok {
		params["ask_job_type_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_verbosity_on_launch"); ok {
		params["ask_verbosity_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_inventory_on_launch"); ok {
		params["ask_inventory_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_credential_on_launch"); ok {
		params["ask_credential_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_execution_environment_on_launch"); ok {
		params["ask_execution_environment_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_labels_on_launch"); ok {
		params["ask_labels_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_forks_on_launch"); ok {
		params["ask_forks_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_job_slice_on_launch"); ok {
		params["ask_job_slice_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_timeout_on_launch"); ok {
		params["ask_timeout_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_instance_groups_on_launch"); ok {
		params["ask_instance_groups_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("survey_enabled"); ok {
		params["survey_enabled"] = p.(bool)
	}
	if p, ok := d.GetOk("become_enabled"); ok {
		params["become_enabled"] = p.(bool)
	}
	if p, ok := d.GetOk("diff_mode"); ok {
		params["diff_mode"] = p.(bool)
	}
	if p, ok := d.GetOk("allow_simultaneous"); ok {
		params["allow_simultaneous"] = p.(bool)
	}
	if p, ok := d.GetOk("job_slice_count"); ok {
		params["job_slice_count"] = p.(int)
	}
	if p, ok := d.GetOk("webhook_service"); ok {
		params["webhook_service"] = p.(string)
	}
	if p, ok := d.GetOk("prevent_instance_group_fallback"); ok {
		params["prevent_instance_group_fallback"] = p.(bool)
	}
	if p, ok := d.GetOk("custom_virtualenv"); ok {
		params["custom_virtualenv"] = AtoipOr(p.(string), nil)
	}
	if p, ok := d.GetOk("credential_id"); ok {
		params["credential"] = p.(int)
	}
	if p, ok := d.GetOk("webhook_credential_id"); ok {
		params["webhook_credential"] = p.(int)
	}

	result, err := awxService.CreateJobTemplate(params, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create JobTemplate",
			Detail:   fmt.Sprintf("JobTemplate %s: %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceJobTemplateRead(ctx, d, m)
}

func resourceJobTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	id, diags := convertStateIDToNummeric("Update JobTemplate", d)
	if diags.HasError() {
		return diags
	}

	checkParams := make(map[string]string)
	_, err := awxService.GetJobTemplateByID(id, checkParams)
	if err != nil {
		return buildDiagNotFoundFail("job template", id, err)
	}

	params := map[string]interface{}{
		"name":      d.Get("name").(string),
		"job_type":  d.Get("job_type").(string),
		"inventory": AtoipOr(d.Get("inventory_id").(string), nil),
		"project":   d.Get("project_id").(int),
		"playbook":  d.Get("playbook").(string),
	}
	if p, ok := d.GetOk("description"); ok {
		params["description"] = p.(string)
	}
	if p, ok := d.GetOk("scm_branch"); ok {
		params["scm_branch"] = p.(string)
	}
	if p, ok := d.GetOk("forks"); ok {
		params["forks"] = p.(int)
	}
	if p, ok := d.GetOk("limit"); ok {
		params["limit"] = p.(string)
	}
	if p, ok := d.GetOk("verbosity"); ok {
		params["verbosity"] = p.(int)
	}
	if p, ok := d.GetOk("extra_vars"); ok {
		params["extra_vars"] = p.(string)
	}
	if p, ok := d.GetOk("job_tags"); ok {
		params["job_tags"] = p.(string)
	}
	if p, ok := d.GetOk("force_handlers"); ok {
		params["force_handlers"] = p.(bool)
	}
	if p, ok := d.GetOk("skip_tags"); ok {
		params["skip_tags"] = p.(string)
	}
	if p, ok := d.GetOk("start_at_task"); ok {
		params["start_at_task"] = p.(string)
	}
	if p, ok := d.GetOk("timeout"); ok {
		params["timeout"] = p.(int)
	}
	if p, ok := d.GetOk("use_fact_cache"); ok {
		params["use_fact_cache"] = p.(bool)
	}
	if p, ok := d.GetOk("organization_id"); ok {
		params["organization_id"] = p.(int)
	}
	if p, ok := d.GetOk("execution_environment_id"); ok {
		params["execution_environment_id"] = p.(int)
	}
	if p, ok := d.GetOk("host_config_key"); ok {
		params["host_config_key"] = p.(string)
	}
	if p, ok := d.GetOk("ask_scm_branch_on_launch"); ok {
		params["ask_scm_branch_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_diff_mode_on_launch"); ok {
		params["ask_diff_mode_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_variables_on_launch"); ok {
		params["ask_variables_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_limit_on_launch"); ok {
		params["ask_limit_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_tags_on_launch"); ok {
		params["ask_tags_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_skip_tags_on_launch"); ok {
		params["ask_skip_tags_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_job_type_on_launch"); ok {
		params["ask_job_type_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_verbosity_on_launch"); ok {
		params["ask_verbosity_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_inventory_on_launch"); ok {
		params["ask_inventory_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_credential_on_launch"); ok {
		params["ask_credential_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_execution_environment_on_launch"); ok {
		params["ask_execution_environment_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_labels_on_launch"); ok {
		params["ask_labels_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_forks_on_launch"); ok {
		params["ask_forks_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_job_slice_on_launch"); ok {
		params["ask_job_slice_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_timeout_on_launch"); ok {
		params["ask_timeout_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("ask_instance_groups_on_launch"); ok {
		params["ask_instance_groups_on_launch"] = p.(bool)
	}
	if p, ok := d.GetOk("survey_enabled"); ok {
		params["survey_enabled"] = p.(bool)
	}
	if p, ok := d.GetOk("become_enabled"); ok {
		params["become_enabled"] = p.(bool)
	}
	if p, ok := d.GetOk("diff_mode"); ok {
		params["diff_mode"] = p.(bool)
	}
	if p, ok := d.GetOk("allow_simultaneous"); ok {
		params["allow_simultaneous"] = p.(bool)
	}
	if p, ok := d.GetOk("job_slice_count"); ok {
		params["job_slice_count"] = p.(int)
	}
	if p, ok := d.GetOk("webhook_service"); ok {
		params["webhook_service"] = p.(string)
	}
	if p, ok := d.GetOk("prevent_instance_group_fallback"); ok {
		params["prevent_instance_group_fallback"] = p.(bool)
	}
	if p, ok := d.GetOk("custom_virtualenv"); ok {
		params["custom_virtualenv"] = AtoipOr(p.(string), nil)
	}
	if p, ok := d.GetOk("credential_id"); ok {
		params["credential_id"] = p.(int)
	}
	if p, ok := d.GetOk("webhook_credential_id"); ok {
		params["webhook_credential"] = p.(int)
	}

	_, err = awxService.UpdateJobTemplate(id, params, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update JobTemplate",
			Detail:   fmt.Sprintf("JobTemplate with name %s in the project id %d faild to update %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	return resourceJobTemplateRead(ctx, d, m)
}

func resourceJobTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	id, diags := convertStateIDToNummeric("Read JobTemplate", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetJobTemplateByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("job template", id, err)

	}
	d = setJobTemplateResourceData(d, res)
	return nil
}

func setJobTemplateResourceData(d *schema.ResourceData, r *awx.JobTemplate) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("job_type", r.JobType)
	d.Set("inventory_id", r.Inventory)
	d.Set("project_id", r.Project)
	d.Set("playbook", r.Playbook)
	d.Set("scm_branch", r.SCMBranch)
	d.Set("forks", r.Forks)
	d.Set("limit", r.Limit)
	d.Set("verbosity", r.Verbosity)
	d.Set("extra_vars", normalizeJsonYaml(r.ExtraVars))
	d.Set("job_tags", r.JobTags)
	d.Set("force_handlers", r.ForceHandlers)
	d.Set("skip_tags", r.SkipTags)
	d.Set("start_at_task", r.StartAtTask)
	d.Set("timeout", r.Timeout)
	d.Set("use_fact_cache", r.UseFactCache)
	d.Set("organization_id", r.Organization)
	d.Set("execution_environment_id", r.ExecutionEnvironment)
	d.Set("host_config_key", r.HostConfigKey)
	d.Set("ask_scm_branch_on_launch", r.AskSCMBranchOnLaunch)
	d.Set("ask_diff_mode_on_launch", r.AskDiffModeOnLaunch)
	d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch)
	d.Set("ask_limit_on_launch", r.AskLimitOnLaunch)
	d.Set("ask_tags_on_launch", r.AskTagsOnLaunch)
	d.Set("ask_skip_tags_on_launch", r.AskSkipTagsOnLaunch)
	d.Set("ask_job_type_on_launch", r.AskJobTypeOnLaunch)
	d.Set("ask_verbosity_on_launch", r.AskVerbosityOnLaunch)
	d.Set("ask_inventory_on_launch", r.AskInventoryOnLaunch)
	d.Set("ask_credential_on_launch", r.AskCredentialOnLaunch)
	d.Set("ask_execution_environment_on_launch", r.AskExecutionEnvironmentOnLaunch)
	d.Set("ask_labels_on_launch", r.AskLabelsOnLaunch)
	d.Set("ask_forks_on_launch", r.AskForksOnLaunch)
	d.Set("ask_job_slice_count_on_launch", r.AskJobSliceCountOnLaunch)
	d.Set("ask_timeout_on_launch", r.AskTimeoutOnLaunch)
	d.Set("ask_inventory_groups_on_launch", r.AskInstanceGroupsOnLaunch)
	d.Set("survey_enabled", r.SurveyEnabled)
	d.Set("become_enabled", r.BecomeEnabled)
	d.Set("diff_mode", r.DiffMode)
	d.Set("allow_simultaneous", r.AllowSimultaneous)
	d.Set("custom_virtualenv", r.CustomVirtualenv)
	d.Set("job_slice_count", r.JobSliceCount)
	d.Set("webhook_service", r.WebHookService)
	d.Set("webhook_credential_id", r.WebHookCredential)
	d.Set("prevent_instance_group_fallback", r.PreventInstanceGroupFallback)
	d.Set("credential_id", r.Credential)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
