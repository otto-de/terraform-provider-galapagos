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
	ts                     = test.NewServer(context.Background())
	restApplicationDetails = &rest.RESTConfig{
		BaseUrl: ts.URL,
		Type:    galapagos.APPLICATION_REST,
	}
)

func TestApplicationCreate(t *testing.T) {

	p := (&applicationResource{
		restDetails: restApplicationDetails,
	}).WithClient(ts.Client())
	d := applicationState{
		name: types.String{Value: "foo"},
		bcap: types.String{Value: "capyb"},
	}
	diags := p.sendCreateToREST(context.Background(), &d)
	if diags.HasError() {
		t.Fatal("sendCreateToREST failed", diags)
	}

	applications := ts.Applications()
	if !reflect.DeepEqual(applications, []test.Application{
		{
			Id:   "app6129484611666145821",
			Name: "foo",
			Bcap: "capyb",
		},
	}) {
		t.Fatalf("Unexpected applications %#v", applications)
	}
}

func TestApplicationDelete(t *testing.T) {

	applications := ts.Applications()
	if len(applications) != 1 {
		panic("Test setup invalid")
	}
	p := (&applicationResource{
		restDetails: restApplicationDetails,
	}).WithClient(ts.Client())
	diags := p.sendDeleteToREST(context.Background(), applications[0].Id)
	if diags.HasError() {
		t.Fatal("sendDeleteToREST failed", diags)
	}

	applications = ts.Applications()
	if !reflect.DeepEqual(applications, []test.Application{}) {
		t.Fatalf("Expected empty applications. Got: %#v", applications)
	}
}

func TestApplicationGet(t *testing.T) {
	applications := ts.Applications()
	if !reflect.DeepEqual(applications, []test.Application{}) {
		panic("Test setup invalid")
	}

	p := (&applicationResource{
		restDetails: restApplicationDetails,
	}).WithClient(ts.Client())
	_, err := p.getStateFromREST(context.Background(), "notapp")
	if err == nil {
		t.Fatal("Error not returned from getState")
	}
	if err != rest.ErrNotExist {
		t.Fatal("Wrong error returned from getState:", err)
	}

	d := &applicationState{
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
