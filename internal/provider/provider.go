package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/otto-de/terraform-provider-galapagos/internal/resources"
)

type provider struct {
}

type providerConfig struct {
}

func New() func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{}
	}
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerConfig
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"galapagos_application": resources.ApplicationType,
		//"galapagos_client_access_key": resource.ClientAccessKeyType,
		//"galapagos_topic": resource.TopicType,
		//"galapagos_schema": resource.SchemaType,

		/* galapagos_user -> OIDC
			// create -> send request / approval
			//
			 * galapagos_client_access_key
			// - secret_id
			// - key

			 * galapagos_topic
			   - name
			   - partitions
			   - criticality
			   - application
			   - env
			   - schema list(string)/list(ref)
			   - cleanup
			   - retention

			 * galapagos_schema
			   - content = string
			   - type
			   - name computed

		         * galapagos_application_subscription
			   - topic
			   - application
		*/
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		/* galapagos_application_ownership
		 */
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{}, nil
	/* a file with
	- access token
	- expiry
	- refresh_token
	- refresh_expiry
	- id_token
	*/
}
