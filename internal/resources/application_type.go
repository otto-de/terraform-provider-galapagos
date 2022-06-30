package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/otto-de/terraform-provider-galapagos/internal/credentials"
)

var (
	ApplicationType = &applicationType{}
)

type applicationType struct {
}

func (t *applicationType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "FIXME",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Type:        types.StringType,
				Description: "Name of Application",
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					// TODO: Implement Update
					tfsdk.RequiresReplace(),
				},
			},
			"id": {
				Type:        types.StringType,
				Description: "FIXME",
				Computed:    true,
			},
			"bcap": {
				Type:        types.StringType,
				Description: "FIXME",
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					// TODO: Implement Update
					tfsdk.RequiresReplace(),
				},
			},
			"aliases": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Description: "FIXME",
				Optional:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					// TODO: Implement Update
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (t *applicationType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	cfg := p.(credentials.ClientWithCredentials)
	return &applicationResource{
		lateClient: cfg,
	}, nil
}
