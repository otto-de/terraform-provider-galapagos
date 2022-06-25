package resources

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/otto-de/terraform-provider-galapagos/internal/galapagos"
)

var (
	ts = galapagos.NewTestServer(context.Background())
)

func TestCreate(t *testing.T) {

	p := (&applicationResource{
		restUrl: ts.URL,
	}).WithClient(ts.Client())
	d := applicationResourceData{
		name: types.String{Value: "foo"},
		bcap: types.String{Value: "capyb"},
	}
	diags := p.sendCreateToREST(context.Background(), &d)
	if diags.HasError() {
		t.Fatal("sendCreateToREST failed", diags)
	}

	applications := ts.Applications()
	if !reflect.DeepEqual(applications, []galapagos.TestApplication{
		{
			Id:   "app5577006791947779410",
			Name: "foo",
			Bcap: "capyb",
		},
	}) {
		t.Fatalf("Unexpected applications %#v", applications)
	}
}

func TestDelete(t *testing.T) {

	applications := ts.Applications()
	if len(applications) != 1 {
		panic("Test setup invalid")
	}
	p := (&applicationResource{
		restUrl: ts.URL,
	}).WithClient(ts.Client())
	diags := p.sendDeleteToREST(context.Background(), applications[0].Id)
	if diags.HasError() {
		t.Fatal("sendDeleteToREST failed", diags)
	}

	applications = ts.Applications()
	if !reflect.DeepEqual(applications, []galapagos.TestApplication{}) {
		t.Fatalf("Expected empty applications. Got: %#v", applications)
	}
}

func TestGet(t *testing.T) {
	applications := ts.Applications()
	if !reflect.DeepEqual(applications, []galapagos.TestApplication{}) {
		panic("Test setup invalid")
	}

	p := (&applicationResource{
		restUrl: ts.URL,
	}).WithClient(ts.Client())
	_, err := p.getStateFromREST(context.Background(), "notapp")
	if err == nil {
		t.Fatal("Error not returned from getState")
	}
	if err != ErrNotExist {
		t.Fatal("Wrong error returned from getState:", err)
	}

	d := &applicationResourceData{
		id:   types.String{Value: "bar"},
		name: types.String{Value: "foobar"},
		bcap: types.String{Value: "capbp"},
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
