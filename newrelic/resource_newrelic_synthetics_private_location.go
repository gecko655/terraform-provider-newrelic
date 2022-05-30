package newrelic

import (
	"context"
	"log"

	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func resourceNewRelicSyntheticsPrivateLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsPrivateLocationCreate,
		ReadContext:   resourceNewRelicSyntheticsPrivateLocationRead,
		UpdateContext: resourceNewRelicSyntheticsPrivateLocationUpdate,
		DeleteContext: resourceNewRelicSyntheticsPrivateLocationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the account in New Relic.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The private location description.",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the private location.",
				ForceNew:    true,
				Required:    true,
			},
			"verified_script_execution": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "The private location requires a password to edit if value is true.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The private location globally unique identifier.",
			},
			"guid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The guid of the entity to tag.",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The private locations key.",
			},
			"location_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An alternate identifier based on name.",
			},
		},
	}
}

func resourceNewRelicSyntheticsPrivateLocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	var diags diag.Diagnostics

	description := d.Get("description").(string)
	name := d.Get("name").(string)
	verifiedScriptExecution := d.Get("verified_script_execution").(bool)
	res, err := client.Synthetics.SyntheticsCreatePrivateLocationWithContext(ctx, accountID, description, name, verifiedScriptExecution)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Errors) > 0 {
		for _, err := range res.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Description,
			})
		}
	}

	if len(diags) > 0 {
		return diags
	}
	d.SetId(string(res.GUID))
	_ = d.Set("domain_id", res.DomainId)
	_ = d.Set("key", res.Key)
	_ = d.Set("location_id", res.LocationId)
	_ = d.Set("guid", res.GUID)

	return nil
}

func resourceNewRelicSyntheticsPrivateLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Reading New Relic Synthetics Private Location %s", d.Id())

	guid := common.EntityGUID(d.Id())
	resp, err := client.Entities.GetEntity(guid)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	setCommonSyntheticsPrivateLocationAttributes(resp, d)
	return nil
}

func setCommonSyntheticsPrivateLocationAttributes(v *entities.EntityInterface, d *schema.ResourceData) {
	switch e := (*v).(type) {
	case *entities.GenericEntityOutline:
		_ = d.Set("guid", e.GUID)
		_ = d.Set("name", e.Name)

	}
}

func resourceNewRelicSyntheticsPrivateLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	var diags diag.Diagnostics
	description := d.Get("description").(string)
	guid := synthetics.EntityGUID(d.Id())
	verifiedScriptExecution := d.Get("verified_script_execution").(bool)
	res, err := client.Synthetics.SyntheticsUpdatePrivateLocation(description, guid, verifiedScriptExecution)

	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Errors) > 0 {
		for _, err := range res.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Description,
			})
		}
	}

	if len(diags) > 0 {
		return diags
	}

	return nil
}

func resourceNewRelicSyntheticsPrivateLocationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	var diags diag.Diagnostics
	guid := synthetics.EntityGUID(d.Id())
	res, err := client.Synthetics.SyntheticsDeletePrivateLocationWithContext(ctx, guid)

	if err != nil {
		return diag.FromErr(err)
	}
	if res != nil {
		for _, err := range res.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Description,
			})
		}
	}

	if len(diags) > 0 {
		return diags
	}

	d.SetId("")
	return nil
}
