package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/otto-de/terraform-provider-galapagos/internal/credentials"
)

var (
	AccountType = &accountType{}
)

type accountType struct {
}

func (t *accountType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "FIXME",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Type:        types.StringType,
				Description: "Name of Account",
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
		},
	}, nil
}

func (t *accountType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	cfg := p.(credentials.ClientWithCredentials)
	return &accountResource{
		client: cfg,
	}, nil
}
