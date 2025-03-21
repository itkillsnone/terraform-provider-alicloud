package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudBrainIndustrialPidOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBrainIndustrialPidOrganizationCreate,
		Read:   resourceAlicloudBrainIndustrialPidOrganizationRead,
		Update: resourceAlicloudBrainIndustrialPidOrganizationUpdate,
		Delete: resourceAlicloudBrainIndustrialPidOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"parent_pid_organization_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pid_organization_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudBrainIndustrialPidOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePidOrganization"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("parent_pid_organization_id"); ok {
		request["ParentOrganizationId"] = v
	}

	request["OrganizationName"] = d.Get("pid_organization_name")
	request["ClientToken"] = buildClientToken("CreatePidOrganization")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_brain_industrial_pid_organization", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["OrganizationId"]))

	return resourceAlicloudBrainIndustrialPidOrganizationRead(d, meta)
}
func resourceAlicloudBrainIndustrialPidOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	brain_industrialService := Brain_industrialService{client}
	object, err := brain_industrialService.DescribeBrainIndustrialPidOrganization(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_brain_industrial_pid_organization brain_industrialService.DescribeBrainIndustrialPidOrganization Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("pid_organization_name", object["OrganizationName"])
	return nil
}
func resourceAlicloudBrainIndustrialPidOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	if d.HasChange("pid_organization_name") {
		request := map[string]interface{}{
			"OrganizationId": d.Id(),
		}
		request["OrganizationName"] = d.Get("pid_organization_name")
		action := "UpdatePidOrganization"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudBrainIndustrialPidOrganizationRead(d, meta)
}
func resourceAlicloudBrainIndustrialPidOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePidOrganization"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"OrganizationId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
