package resources

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/otto-de/terraform-provider-galapagos/internal/galapagos"
	"github.com/otto-de/terraform-provider-galapagos/internal/galapagos/test"
	"github.com/otto-de/terraform-provider-galapagos/internal/rest"
)

var (
	restAccountDetails = &rest.RESTConfig{
		BaseUrl: ts.URL,
		Type:    galapagos.ACCOUNT_REST,
	}
)

func TestAccountCreate(t *testing.T) {

	p := (&accountResource{
		restDetails: restAccountDetails,
	}).WithClient(ts.Client())
	d := accountState{
		name: types.String{Value: "foo"},
	}
	diags := p.sendCreateToREST(context.Background(), &d)
	if diags.HasError() {
		t.Fatal("sendCreateToREST failed", diags)
	}

	accounts := ts.Accounts()
	if !reflect.DeepEqual(accounts, []test.Account{
		{
			Id:   "acc5577006791947779410",
			Name: "foo",
		},
	}) {
		t.Fatalf("Unexpected account %#v", accounts)
	}
}

func TestAccountDelete(t *testing.T) {

	accounts := ts.Accounts()
	if len(accounts) != 1 {
		panic("Test setup invalid")
	}
	p := (&accountResource{
		restDetails: restAccountDetails,
	}).WithClient(ts.Client())
	diags := p.sendDeleteToREST(context.Background(), accounts[0].Id)
	if diags.HasError() {
		t.Fatal("sendDeleteToREST failed", diags)
	}

	accounts = ts.Accounts()
	if !reflect.DeepEqual(accounts, []test.Account{}) {
		t.Fatalf("Expected empty accounts. Got: %#v", accounts)
	}
}

func TestAccountGet(t *testing.T) {
	accounts := ts.Accounts()
	if !reflect.DeepEqual(accounts, []test.Account{}) {
		panic("Test setup invalid")
	}

	p := (&accountResource{
		restDetails: restAccountDetails,
	}).WithClient(ts.Client())
	_, err := p.getStateFromREST(context.Background(), "notapp")
	if err == nil {
		t.Fatal("Error not returned from getState")
	}
	if err != rest.ErrNotExist {
		t.Fatal("Wrong error returned from getState:", err)
	}

	d := &accountState{
		id:   types.String{Value: "bar"},
		name: types.String{Value: "foobar"},
	}
	diags := p.sendCreateToREST(context.Background(), d)
	if diags.HasError() {
		t.Fatal("sendCreateToREST failed", diags)
	}

	state, err := p.getStateFromREST(context.Background(), d.id.Value)
	if err != nil {
		t.Fatal("getStateFromREST failed:", err)
	}
	if !reflect.DeepEqual(d, state) {
		t.Fatalf("Remote State mismatch:\n%#v\n%#v\n", state, d)
	}
}
