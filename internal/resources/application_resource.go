package resources

import (
	"context"
	"errors"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/otto-de/terraform-provider-galapagos/internal/credentials"
	"github.com/otto-de/terraform-provider-galapagos/internal/rest"
	"github.com/otto-de/terraform-provider-galapagos/internal/typeconvert"
)

type applicationCreateRequest struct {
	Name    string   `json:"name"`
	Bcap    string   `json:"bcap,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
}

type applicationCreateResponse struct {
	Id string `json:"id"`
}

type applicationDeleteRequest struct {
	id string
}

type applicationDescribeRequest struct {
	id string
}

type applicationDescribeResponse struct {
	Name    string   `json:"name"`
	Bcap    string   `json:"bcap"`
	Aliases []string `json:"aliases"`
}

type applicationResource struct {
	lateClient  credentials.ClientWithCredentials
	restDetails *rest.RESTConfig
}

func (r *applicationDeleteRequest) String() string {
	return r.id
}

func (r *applicationDescribeRequest) String() string {
	return r.id
}

func (r *applicationResource) WithClient(client *http.Client) *applicationResource {
	r.lateClient = credentials.WithClient(client)
	return r
}

func (r *applicationResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var d applicationState
	diags := req.Plan.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.sendCreateToREST(ctx, &d)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, d)
	resp.Diagnostics.Append(diags...)
}

func (r *applicationResource) sendCreateToREST(ctx context.Context, d *applicationState) (diags diag.Diagnostics) {

	cd := applicationCreateRequest{
		Name:    d.name.Value,
		Bcap:    d.bcap.Value,
		Aliases: typeconvert.ToStringSlice(d.aliases),
	}
	ans := applicationCreateResponse{}
	client := r.lateClient.Client(ctx)
	err := rest.NewClient(r.restDetails, client).SendCreate(ctx, rest.CreateOptions{}, &cd, &ans)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic(err.Error(), ""))
		return
	}

	d.id = types.String{
		Value: ans.Id,
	}
	return
}

func (r *applicationResource) sendDeleteToREST(ctx context.Context, id string) (diags diag.Diagnostics) {

	deleteReq := applicationDeleteRequest{
		id: id,
	}
	deleteResp := struct{}{}
	client := r.lateClient.Client(ctx)
	err := rest.NewClient(r.restDetails, client).SendDelete(ctx, &deleteReq, &deleteResp)
	if err != nil {
		return append(diags, diag.NewErrorDiagnostic(err.Error(), ""))
	}
	return
}

func (r *applicationResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var d applicationState
	diags := req.State.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.sendDeleteToREST(ctx, d.id.Value)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func (r *applicationResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var s applicationState
	diags := req.State.Get(ctx, &s)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.getStateFromREST(ctx, s.id.Value)
	if errors.Is(err, rest.ErrNotExist) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	s.name = state.name
	s.bcap = state.bcap
	s.aliases = state.aliases

	diags = resp.State.Set(ctx, &s)
	resp.Diagnostics.Append(diags...)
}

func (r *applicationResource) getStateFromREST(ctx context.Context, id string) (*applicationState, error) {

	client := r.lateClient.Client(ctx)
	req := applicationDescribeRequest{
		id: id,
	}
	ans := applicationDescribeResponse{}
	err := rest.NewClient(r.restDetails, client).SendDescribe(ctx, &req, &ans)
	if err != nil {
		return nil, err
	}

	s := &applicationState{
		id: types.String{
			Value: id,
		},
		name: types.String{
			Value: ans.Name,
		},
		bcap: types.String{
			Value: ans.Bcap,
		},
		aliases: typeconvert.ToTypesStrings(ans.Aliases),
	}
	return s, nil
}

func (r *applicationResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	panic("Update not implemented")
}

func (r *applicationResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {

	// TODO: Clarify whether names are unique
	state, err := r.getStateFromREST(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	var d applicationState
	d.name = state.name
	d.bcap = state.bcap
	d.aliases = state.aliases

	diags := resp.State.Set(ctx, &d)
	resp.Diagnostics.Append(diags...)
	return
}
