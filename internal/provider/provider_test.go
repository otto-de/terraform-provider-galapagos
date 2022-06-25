package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/otto-de/terraform-provider-galapagos/internal/credentials"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	Provider = New()()
)

func TestProvider(t *testing.T) {
	resp := &tfsdk.ConfigureProviderResponse{}
	cfg := clientcredentials.Config{}
	configureFromConfigAttribute(context.TODO(), types.String{}, &cfg, resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp)
	}
}

func TestProviderClientWithCredentials(t *testing.T) {
	p := &provider{}
	_ = credentials.ClientWithCredentials(p)
}
