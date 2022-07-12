package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TestApplicationController mimics the ApplicationController found in Galapagos Java.
type ApplicationController struct {
	Resources []Application
}

type Application struct {
	Id      string
	Name    string
	Bcap    string
	Aliases []string
}

func (c *ApplicationController) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	reqJson := struct {
		Name    string   `json:"name"`
		Bcap    string   `json:"bcap"`
		Aliases []string `json:"aliases"`
	}{}
	jd := json.NewDecoder(r.Body)
	err := jd.Decode(&reqJson)
	if err != nil {
		tflog.Error(ctx, "Decode failed", map[string]interface{}{
			"error": err,
		})
	}
	id := fmt.Sprintf("app%d", rand.Int())
	c.Resources = append(c.Resources, Application{
		Id:      id,
		Name:    reqJson.Name,
		Bcap:    reqJson.Bcap,
		Aliases: reqJson.Aliases,
	})

	respJson := struct {
		Id string `json:"id"`
	}{
		Id: id,
	}
	je := json.NewEncoder(w)
	err = je.Encode(&respJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tflog.Error(ctx, "Encode failed", map[string]interface{}{
			"error": err,
		})
		return
	}
}

func (c *ApplicationController) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deleteI := -1
	for i, app := range c.Resources {
		if app.Id == r.URL.Path {
			deleteI = i
			break
		}
	}
	c.Resources = append(c.Resources[0:deleteI], c.Resources[deleteI+1:]...)
}

func (c *ApplicationController) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	foundI := -1
	for i, app := range c.Resources {
		if app.Id == r.URL.Path {
			foundI = i
			break
		}
	}
	if foundI == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	app := c.Resources[foundI]
	je := json.NewEncoder(w)
	err := je.Encode(&app)
	if err != nil {
		panic(err)
	}
}

func (c *ApplicationController) Head(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	foundI := -1
	for i, app := range c.Resources {
		if app.Id == r.URL.Path {
			foundI = i
			break
		}
	}
	if foundI == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
POST /api/application addApplication(
PUT /api/me/requests submitApplicationOwnerRequest(
DELETE /api/me/requests/{id} cancelApplicationOwnerRequest
POST /api/user/createApplicationRequests submitApplicationCreationRequest(
*/
