package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type accountResource struct {
}

func (r *accountResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// TODO: Implement me
	// create -> send request / approval
}

func (r *accountResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// TODO: Implement me
}

func (r *accountResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// TODO: Implement me
}

func (r *accountResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	panic("Update not implemented")
}

func (r *accountResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// TODO: Implement me
}
