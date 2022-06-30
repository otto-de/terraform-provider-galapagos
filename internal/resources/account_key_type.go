package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/otto-de/terraform-provider-galapagos/internal/credentials"
)

var (
	AccountKeyType = &accountKeyType{}
)

type accountKeyType struct {
}

func (t *accountKeyType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "FIXME",
		Attributes: map[string]tfsdk.Attribute{
			"account_id": {
				Type:        types.StringType,
				Description: "FIXME",
				Required:    true,
			},
			"key_algorithm": {
				Type:        types.StringType,
				Description: "The algorithm used to generate the key.",
				Required:    true,
			},
			"keepers": {
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Description: "Arbitrary map of values that, when changed, will trigger a new key to be generated.",
				Optional:    true,
			},
		},
	}, nil
}

func (t *accountKeyType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	cfg := p.(credentials.ClientWithCredentials)
	return &accountKeyResource{
		client: cfg,
	}, nil
}
