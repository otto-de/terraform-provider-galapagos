package resources

import "github.com/hashicorp/terraform-plugin-framework/types"

type accountState struct {
	id     types.String `tfsdk:"id"`
	name   types.String `tfsdk:"name"`
	status types.String `tfsdk:"status"`
}
