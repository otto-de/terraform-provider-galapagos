package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TestServer struct {
	URL          string
	accounts     []Account
	applications ApplicationController
	topics       TopicController
	accRand      *rand.Rand
	appRand      *rand.Rand
	ts           *httptest.Server
}

func NewServer(ctx context.Context) *TestServer {
	mux := http.NewServeMux()
	s := &TestServer{
		accounts:     []Account{},
		applications: ApplicationController{},
		topics:       TopicController{},
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
		s.applications.Create(ctx, w, r)
	}
	topicsHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			panic(r.Method)
		}
		s.topics.Create(ctx, w, r)
	}

	// /api is hardcoded in Galapagos so mirror that here
	mux.HandleFunc("/api/accounts", accountsHandler)
	mux.HandleFunc("/api/applications", applicationsHandler)
	mux.HandleFunc("/api/topics", topicsHandler)
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
			s.applications.Delete(ctx, w, r)
		case http.MethodGet:
			s.applications.Get(ctx, w, r)
		case http.MethodHead:
			s.applications.Head(ctx, w, r)
		default:
			panic(r.Method)
		}
	}
	topicHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			s.topics.Delete(ctx, w, r)
		case http.MethodGet:
			s.topics.GetConfig(ctx, w, r)
		default:
			panic(r.Method)
		}
	}
	mux.Handle("/api/account/", http.StripPrefix("/api/account/", http.HandlerFunc(accountHandler)))
	mux.Handle("/api/application/", http.StripPrefix("/api/application/", http.HandlerFunc(applicationHandler)))
	mux.Handle("/api/topic/", http.StripPrefix("/api/topic", http.HandlerFunc(topicHandler)))
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

func (s *TestServer) Accounts() []Account {
	return s.accounts
}

func (s *TestServer) Applications() []Application {
	return s.applications.Resources
}

func (s *TestServer) Topics() []Topic {
	return s.topics.Resources
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
	s.accounts = append(s.accounts, Account{
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
