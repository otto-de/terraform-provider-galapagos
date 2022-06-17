package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/otto-de/terraform-provider-galapagos/internal/resources"
	"golang.org/x/oauth2/clientcredentials"
)

type provider struct {
}

type providerConfig struct {
	Config types.String
}

type authConfig struct {
	Issuer       string    `json:"issuer"`
	ClientId     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	AccessToken  string    `json:"access_token"`
	Expiry       time.Time `json:"expiry"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
}

func New() func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{}
	}
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var pc providerConfig
	diags := req.Config.Get(ctx, &pc)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	configureFromConfigAttribute(ctx, pc.Config, resp)
}

func configureFromConfigAttribute(ctx context.Context, config types.String, resp *tfsdk.ConfigureProviderResponse) {

	if config.Null || config.Value == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		configureFromConfigPath(ctx, fmt.Sprintf("%s/galapagos/config.json", dir), resp)
		return
	}

	configureFromConfigPath(ctx, config.Value, resp)
}

func configureFromConfigPath(ctx context.Context, configFilePath string, resp *tfsdk.ConfigureProviderResponse) {

	auth := authConfig{}
	func() {
		f, err := os.Open(configFilePath)
		if err != nil {
			ap := tftypes.NewAttributePath().WithAttributeName("config")
			resp.Diagnostics.AddAttributeError(ap, fmt.Sprintf("Could not find file %s", configFilePath), err.Error())
			return
		}
		defer f.Close()
		d := json.NewDecoder(f)
		err = d.Decode(&auth)
		if err != nil {
			ap := tftypes.NewAttributePath().WithAttributeName("config")
			resp.Diagnostics.AddAttributeError(ap, fmt.Sprintf("Could not decode config %s", configFilePath), err.Error())
		}
	}()

	if resp.Diagnostics.HasError() {
		return
	}

	p, err := oidc.NewProvider(ctx, auth.Issuer)
	if err != nil {
		resp.Diagnostics.AddError("Creating OIDC provider failed", err.Error())
		return
	}

	cfg := clientcredentials.Config{
		ClientID:     auth.ClientId,
		ClientSecret: auth.ClientSecret,
		TokenURL:     p.Endpoint().TokenURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	tkn, err := cfg.Token(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Fetching token failed", err.Error())
		return
	}

	auth.AccessToken = tkn.AccessToken
	auth.Expiry = tkn.Expiry
	auth.RefreshToken = tkn.RefreshToken
	auth.TokenType = tkn.TokenType
	f, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		resp.Diagnostics.AddWarning("Writing back to config failed", err.Error())
		return
	}
	defer f.Close()
	je := json.NewEncoder(f)
	je.SetIndent("", "  ")
	err = je.Encode(&auth)
	if err != nil {
		resp.Diagnostics.AddWarning("Encoding to config failed", err.Error())
	}
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"galapagos_account":     resources.AccountType,
		"galapagos_account_key": resources.AccountKeyType,
		"galapagos_application": resources.ApplicationType,
		//"galapagos_client_access_key": resource.ClientAccessKeyType,
		//"galapagos_topic": resource.TopicType,
		//"galapagos_schema": resource.SchemaType,

		/*
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
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"config": {
				Description: `Configuration file for client access to Galapagos. Defaults to ~/.config/galapagos/config.json.`,
				Type:        types.StringType,
				Optional:    true,
			},
		},
	}, nil
}
