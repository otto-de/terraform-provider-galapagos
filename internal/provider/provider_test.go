package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestProvider(t *testing.T) {
	resp := &tfsdk.ConfigureProviderResponse{}
	configureFromConfigAttribute(context.TODO(), types.String{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp)
	}
}
