package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/otto-de/terraform-provider-galapagos/internal/credentials"
	"github.com/otto-de/terraform-provider-galapagos/internal/typeconvert"
)

var (
	ErrNotExist = errors.New("Application does not exist")
)

type applicationCreateRequest struct {
	Name    string   `json:"name"`
	Bcap    string   `json:"bcap,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
}

type applicationCreateResponse struct {
	Id string `json:"id"`
}

type applicationDescribeResponse struct {
	Name    string   `json:"name"`
	Bcap    string   `json:"bcap"`
	Aliases []string `json:"aliases"`
}

type applicationResource struct {
	lateClient credentials.ClientWithCredentials
	restUrl    string
}

func (r *applicationResource) WithClient(client *http.Client) *applicationResource {
	r.lateClient = credentials.WithClient(client)
	return r
}

func (r *applicationResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var d applicationResourceData
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

func (r *applicationResource) sendCreateToREST(ctx context.Context, d *applicationResourceData) (diags diag.Diagnostics) {

	jsonReader, jsonWriter := io.Pipe()
	defer jsonReader.Close()

	go func() {

		defer jsonWriter.Close()
		cd := applicationCreateRequest{
			Name:    d.name.Value,
			Bcap:    d.bcap.Value,
			Aliases: typeconvert.ToStringSlice(d.aliases),
		}
		je := json.NewEncoder(jsonWriter)
		err := je.Encode(&cd)
		if err != nil {
			tflog.Error(ctx, "Async JSON encoding failed", map[string]interface{}{
				"err": err,
			})
		}
	}()

	c := r.lateClient.Client(ctx)
	postResp, err := c.Post(fmt.Sprintf("%s/applications", r.restUrl), "application/json", jsonReader)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic(fmt.Sprintf("Post call to %s failed", r.restUrl), err.Error()))
		return
	}

	ans := applicationCreateResponse{}
	jd := json.NewDecoder(postResp.Body)
	err = jd.Decode(&ans)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic("JSON decoding failed", err.Error()))
		return
	}

	d.id = types.String{
		Value: ans.Id,
	}
	return
}

func (r *applicationResource) sendDeleteToREST(ctx context.Context, id string) (diags diag.Diagnostics) {

	c := r.lateClient.Client(ctx)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/application/%s", r.restUrl, id), nil)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic(fmt.Sprintf("Building request to %s failed", r.restUrl), err.Error()))
		return
	}
	_, err = c.Do(req)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic(fmt.Sprintf("Delete call to %s failed", r.restUrl), err.Error()))
		return
	}

	return
}

func (r *applicationResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var d applicationResourceData
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
	var d applicationResourceData
	diags := req.State.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.getStateFromREST(ctx, d.id.Value)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	if err == ErrNotExist {
		resp.State.RemoveResource(ctx)
		return
	}

	d.name = state.name
	d.bcap = state.bcap
	d.aliases = state.aliases

	diags = resp.State.Set(ctx, &d)
	resp.Diagnostics.Append(diags...)
}

func (r *applicationResource) getStateFromREST(ctx context.Context, id string) (d *applicationResourceData, err error) {

	c := r.lateClient.Client(ctx)
	resp, err := c.Get(fmt.Sprintf("%s/application/%s", r.restUrl, id))
	if err != nil {
		err = fmt.Errorf("Delete call to %s failed: %w", r.restUrl, err)
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		err = ErrNotExist
		return
	case http.StatusOK:
		break
	default:
		panic(fmt.Sprintf("Status code %d unexpected", resp.StatusCode))
	}

	ans := applicationDescribeResponse{}
	jd := json.NewDecoder(resp.Body)
	err = jd.Decode(&ans)
	if err != nil {
		err = fmt.Errorf("JSON decoding failed: %w", err)
		return
	}

	d = &applicationResourceData{
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
	return
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

	var d applicationResourceData
	d.name = state.name
	d.bcap = state.bcap
	d.aliases = state.aliases

	diags := resp.State.Set(ctx, &d)
	resp.Diagnostics.Append(diags...)
	return
}
