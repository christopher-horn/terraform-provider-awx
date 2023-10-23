/*

Example Usage

```hcl
data "awx_inventory" "default" {
  name            = "private_services"
  organization_id = data.awx_organization.default.id
}

data "awx_project" "baseconfig" {
  name = "base_config"
}

resource "awx_job_template" "baseconfig" {
  name                     = "baseconfig"
  job_type                 = "run"
  inventory_id             = data.awx_inventory.default.id
  project_id               = data.awx_project.baseconfig.id
  playbook                 = "master-configure-system.yml"
  become_enabled           = true
  ask_credential_on_launch = true
  ask_limit_on_launch      = true
}

resource "awx_job_template_launch" "postinstall" {
  job_template_id        = awx_job_template.baseconfig.id
  limit                  = "sample-hostname"
  credential_ids         = [3]
  monitor_for_completion = true
}
```

*/

package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/mrcrilly/goawx/client"
)

func resourceJobTemplateLaunch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobTemplateLaunchCreate,
		ReadContext:   resourceJobRead,
		DeleteContext: resourceJobDelete,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Job template ID",
				ForceNew:    true,
			},
			"monitor_for_completion": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If true monitor job for successful completion",
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Job type, one of run, check or scan",
			},
			"playbook": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Playbook file name",
			},
			"forks": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Forks",
			},
			"limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Limit",
			},
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "One of 0,1,2,3,4,5",
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Extra variables",
			},
			"job_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Job tags",
			},
			"skip_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Skip tags",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Timeout",
			},
			"scm_revision": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SCM revision",
			},
			"diff_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Diff mode",
			},
			"credential_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Computed:    true,
				Description: "A list of credential IDs",
			},
			"execution_environment_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Execution environment ID",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
		},
	}
}

func resourceJobTemplateLaunchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	jobTemplateID := d.Get("job_template_id").(int)
	_, err := awxService.GetJobTemplateByID(jobTemplateID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("job template", jobTemplateID, err)
	}

	params := make(map[string]interface{})

	if p, ok := d.GetOk("job_type"); ok {
		params["job_type"] = p.(string)
	}
	if p, ok := d.GetOk("playbook"); ok {
		params["playbook"] = p.(string)
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
	if p, ok := d.GetOk("skip_tags"); ok {
		params["skip_tags"] = p.(string)
	}
	if p, ok := d.GetOk("timeout"); ok {
		params["timeout"] = p.(int)
	}
	if p, ok := d.GetOk("scm_revision"); ok {
		params["scm_revision"] = p.(string)
	}
	if p, ok := d.GetOk("diff_mode"); ok {
		params["diff_mode"] = p.(bool)
	}
	if p, ok := d.GetOk("credential_ids"); ok {
		params["credentials"] = expandIntList(p.(*schema.Set).List())
	}
	if p, ok := d.GetOk("execution_environment_id"); ok {
		params["execution_environment"] = p.(int)
	}

	res, err := awxService.Launch(jobTemplateID, params, map[string]string{})
	if err != nil {
		log.Printf("Failed to launch job template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to launch job from job template",
			Detail:   fmt.Sprintf("JobTemplate with name %s in the project id %d, failed to create %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	var jobID = strconv.Itoa(res.ID)
	d.SetId(jobID)

	if d.Get("monitor_for_completion").(bool) {
		_, err = isWaitForJobComplete(ctx, client, jobID)
		if err != nil {
			log.Printf("Job failed %v", err)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Job %s", err.Error()),
				Detail:   fmt.Sprintf("Job %s status: %s", jobID, err.Error()),
			})
			return diags
		}
	}

	return resourceJobRead(ctx, d, m)
}

func resourceJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobService
	jobID, diags := convertStateIDToNummeric("Read Job", d)
	job, err := awxService.GetJob(jobID, map[string]string{})
	if err != nil {
		time.Sleep(3 * time.Second)
		retryJob, retryErr := awxService.GetJob(jobID, map[string]string{})
		if retryErr != nil {
			if check := strings.Contains(retryErr.Error(), "404 Not Found"); check {
				d.SetId("")
				return diags
			}
			return buildDiagNotFoundFail("job", jobID, retryErr)
		}
		job = retryJob
	}
	d = setJobResourceData(d, job)
	return diags
}

func resourceJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobService
	jobID, diags := convertStateIDToNummeric("Delete Job", d)
	_, err := awxService.GetJob(jobID, map[string]string{})
	if err != nil {
		log.Printf("Job already deleted %d", jobID)
		d.SetId("")
		return diags
	}
	// Should provide user option and if desired actually delete it
	d.SetId("")
	return diags
}

func setJobResourceData(d *schema.ResourceData, r *awx.Job) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("status", r.Status)
	d.Set("failed", r.Failed)

	d.SetId(strconv.Itoa(r.ID))
	return d
}

func isWaitForJobComplete(ctx context.Context, client *awx.AWX, id string) (interface{}, error) {
	log.Printf("Waiting for job %s to complete ", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{awx.JobStatusNew, awx.JobStatusPending, awx.JobStatusWaiting, awx.JobStatusRunning},
		Target:     []string{awx.JobStatusSuccessful, awx.JobStatusCanceled, awx.JobStatusError, awx.JobStatusFailed},
		Refresh:    isJobRefreshFunc(client, id),
		Delay:      30 * time.Second,
		MinTimeout: 10 * time.Second,
		Timeout:    120 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isJobRefreshFunc(client *awx.AWX, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var jobID int
		jobID, _ = strconv.Atoi(id)
		job, err := client.JobService.GetJob(jobID, map[string]string{})
		if err != nil {
			time.Sleep(5 * time.Second)
			retryJob, retryErr := client.JobService.GetJob(jobID, map[string]string{})
			if retryErr != nil {
				return nil, "", retryErr
			}
			job = retryJob
		}
		if job.Status == awx.JobStatusSuccessful {
			return job, job.Status, nil
		}
		if job.Status == awx.JobStatusError {
			return job, job.Status, fmt.Errorf("error")
		}
		if job.Status == awx.JobStatusFailed {
			return job, job.Status, fmt.Errorf("failed")
		}
		if job.Status == awx.JobStatusCanceled {
			return job, job.Status, fmt.Errorf("canceled")
		}
		return job, job.Status, nil
	}
}
