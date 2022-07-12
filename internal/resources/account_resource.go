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
)

type accountCreateRequest struct {
	Name string `json:"name"`
}

type accountCreateResponse struct {
	Id string `json:"id"`
}

type accountDeleteRequest struct {
	id string
}

type accountDescribeRequest struct {
	id string
}

type accountDescribeResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type accountResource struct {
	lateClient  credentials.ClientWithCredentials
	restDetails *rest.RESTConfig
}

func (r *accountDescribeRequest) String() string {
	return r.id
}

func (r *accountDeleteRequest) String() string {
	return r.id
}

func (r *accountResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var s accountState
	diags := req.Plan.Get(ctx, &s)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.sendCreateToREST(ctx, &s)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, s)
	resp.Diagnostics.Append(diags...)
}

func (r *accountResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var s accountState
	diags := req.State.Get(ctx, &s)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.sendDeleteToREST(ctx, s.id.Value)
}

func (r *accountResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	s := accountState{}
	diags := req.State.Get(ctx, &s)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rs, err := r.getStateFromREST(ctx, s.id.Value)
	if errors.Is(err, rest.ErrNotExist) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	s.name = rs.name
	s.status = rs.status
	diags = resp.State.Set(ctx, &s)
	resp.Diagnostics.Append(diags...)
}

func (r *accountResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	panic("Update not implemented")
}

func (r *accountResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// TODO: Clarify whether names are unique
	rs, err := r.getStateFromREST(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	var s accountState
	s.id = rs.id
	s.name = rs.name
	s.status = rs.status

	diags := resp.State.Set(ctx, &s)
	resp.Diagnostics.Append(diags...)
	return
}

func (r *accountResource) WithClient(client *http.Client) *accountResource {
	r.lateClient = credentials.WithClient(client)
	return r
}

func (r *accountResource) getStateFromREST(ctx context.Context, id string) (*accountState, error) {

	client := r.lateClient.Client(ctx)
	req := accountDescribeRequest{
		id: id,
	}
	ans := accountDescribeResponse{}
	err := rest.NewClient(r.restDetails, client).SendDescribe(ctx, rest.DescribeOptions{}, &req, &ans)
	if err != nil {
		return nil, err
	}

	s := &accountState{
		id: types.String{
			Value: id,
		},
		name: types.String{
			Value: ans.Name,
		},
		status: types.String{
			Value: ans.Status,
		},
	}
	return s, nil
}

func (r *accountResource) sendCreateToREST(ctx context.Context, s *accountState) diag.Diagnostics {

	createReq := accountCreateRequest{
		Name: s.name.Value,
	}
	createResp := accountCreateResponse{}

	client := r.lateClient.Client(ctx)
	err := rest.NewClient(r.restDetails, client).SendCreate(ctx, rest.CreateOptions{}, &createReq, &createResp)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic(err.Error(), "")}
	}

	s.id = types.String{
		Value: createResp.Id,
	}
	return nil
}

func (r *accountResource) sendDeleteToREST(ctx context.Context, id string) diag.Diagnostics {
	deleteReq := accountDeleteRequest{
		id: id,
	}
	deleteResp := struct{}{}

	client := r.lateClient.Client(ctx)
	err := rest.NewClient(r.restDetails, client).SendDelete(ctx, &deleteReq, &deleteResp)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic(err.Error(), "")}
	}
	return nil
}
