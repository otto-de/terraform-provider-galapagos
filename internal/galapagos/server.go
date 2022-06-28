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

type TestAccount struct {
	Id   string
	Name string
}

type TestApplication struct {
	Id      string
	Name    string
	Bcap    string
	Aliases []string
}

type TestServer struct {
	URL          string
	accounts     []TestAccount
	applications []TestApplication
	accRand      *rand.Rand
	appRand      *rand.Rand
	ts           *httptest.Server
}

func NewTestServer(ctx context.Context) *TestServer {
	mux := http.NewServeMux()
	s := &TestServer{
		accounts:     []TestAccount{},
		applications: []TestApplication{},
	}
	accountsHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			panic(r.Method)
		}
		s.accountCreate(ctx, w, r)
	}
	applicationsHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			panic(r.Method)
		}
		s.applicationCreate(ctx, w, r)
	}

	mux.HandleFunc("/accounts", accountsHandler)
	mux.HandleFunc("/applications", applicationsHandler)
	accountHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			s.accountDelete(ctx, w, r)
		case http.MethodGet:
			s.accountGet(ctx, w, r)
		case http.MethodHead:
			s.accountHead(ctx, w, r)
		default:
			panic(r.Method)
		}
	}
	applicationHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			s.applicationDelete(ctx, w, r)
		case http.MethodGet:
			s.applicationGet(ctx, w, r)
		case http.MethodHead:
			s.applicationHead(ctx, w, r)
		default:
			panic(r.Method)
		}
	}
	mux.Handle("/account/", http.StripPrefix("/account/", http.HandlerFunc(accountHandler)))
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

func (s *TestServer) Accounts() []TestAccount {
	return s.accounts
}

func (s *TestServer) Applications() []TestApplication {
	return s.applications
}

func (s *TestServer) accountCreate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	reqJson := struct {
		Name string `json:"name"`
	}{}
	jd := json.NewDecoder(r.Body)
	err := jd.Decode(&reqJson)
	if err != nil {
		tflog.Error(ctx, "Decode failed", map[string]interface{}{
			"error": err,
		})
	}
	id := fmt.Sprintf("acc%d", rand.Int())
	s.accounts = append(s.accounts, TestAccount{
		Id:   id,
		Name: reqJson.Name,
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

func (s *TestServer) accountDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deleteI := -1
	for i, app := range s.accounts {
		if app.Id == r.URL.Path {
			deleteI = i
			break
		}
	}
	s.accounts = append(s.accounts[0:deleteI], s.accounts[deleteI+1:]...)
}

func (s *TestServer) accountGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	foundI := -1
	for i, acc := range s.accounts {
		if acc.Id == r.URL.Path {
			foundI = i
			break
		}
	}
	if foundI == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	acc := s.accounts[foundI]
	je := json.NewEncoder(w)
	err := je.Encode(&acc)
	if err != nil {
		panic(err)
	}
}

func (s *TestServer) accountHead(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	foundI := -1
	for i, acc := range s.accounts {
		if acc.Id == r.URL.Path {
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

func (s *TestServer) applicationCreate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func (s *TestServer) applicationDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deleteI := -1
	for i, app := range s.applications {
		if app.Id == r.URL.Path {
			deleteI = i
			break
		}
	}
	s.applications = append(s.applications[0:deleteI], s.applications[deleteI+1:]...)
}

func (s *TestServer) applicationGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func (s *TestServer) applicationHead(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
