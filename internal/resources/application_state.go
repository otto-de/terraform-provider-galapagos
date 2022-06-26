package resources

import "github.com/hashicorp/terraform-plugin-framework/types"

type applicationState struct {
	name    types.String   `tfsdk:"name"`
	id      types.String   `tfsdk:"id"`
	bcap    types.String   `tfsdk:"bcap"`
	aliases []types.String `tfsdk:"aliases"`
}
