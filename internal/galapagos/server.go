package galapagos

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TestApplication struct {
	Id      string
	Name    string
	Bcap    string
	Aliases []string
}

type TestServer struct {
	URL          string
	applications []TestApplication
	ts           *httptest.Server
}

func NewTestServer(ctx context.Context) *TestServer {
	mux := http.NewServeMux()
	s := &TestServer{
		applications: []TestApplication{},
	}
	applicationsHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			panic(r.Method)
		}
		s.create(ctx, w, r)
	}

	mux.HandleFunc("/applications", applicationsHandler)
	applicationHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			s.delete(ctx, w, r)
		case http.MethodGet:
			s.get(ctx, w, r)
		case http.MethodHead:
			s.head(ctx, w, r)
		default:
			panic(r.Method)
		}
	}
	mux.Handle("/application/", http.StripPrefix("/application/", http.HandlerFunc(applicationHandler)))
	s.ts = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		mux.ServeHTTP(w, r)
	}))
	s.URL = s.ts.URL
	return s
}

func (s *TestServer) Close() error {
	s.ts.Close()
	return nil
}

func (s *TestServer) Client() *http.Client {
	return s.ts.Client()
}

func (s *TestServer) Applications() []TestApplication {
	return s.applications
}

func (s *TestServer) create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	s.applications = append(s.applications, TestApplication{
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

func (s *TestServer) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deleteI := -1
	for i, app := range s.applications {
		if app.Id == r.URL.Path {
			deleteI = i
			break
		}
	}
	s.applications = append(s.applications[0:deleteI], s.applications[deleteI+1:]...)
}

func (s *TestServer) get(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	foundI := -1
	for i, app := range s.applications {
		if app.Id == r.URL.Path {
			foundI = i
			break
		}
	}
	if foundI == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	app := s.applications[foundI]
	je := json.NewEncoder(w)
	err := je.Encode(&app)
	if err != nil {
		panic(err)
	}
}

func (s *TestServer) head(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	foundI := -1
	for i, app := range s.applications {
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
