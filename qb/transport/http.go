package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"mooncamp.com/dgraphtools/gql"
	"mooncamp.com/dgraphtools/qb"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(eps qb.EndpointSet, pathPrefix string) http.Handler {
	r := mux.NewRouter().PathPrefix(pathPrefix).Subrouter()
	r.Handle("/template", httptransport.NewServer(
		eps.Template,
		decodeTemplateRequest,
		encodeTemplateResponse,
	)).Methods(http.MethodPost)
	r.Handle("/parse", httptransport.NewServer(
		eps.Parse,
		decodeParseRequest,
		encodeParseResponse,
	)).Methods(http.MethodPost)

	return r
}

func decodeTemplateRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	req := struct {
		Queries   []gql.GraphQuery  `json:"queries"`
		Alias     string            `json:"alias"`
		Variables map[string]string `json:"variables"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return qb.TemplateRequest{
		Queries:   req.Queries,
		Alias:     req.Alias,
		Variables: req.Variables,
	}, nil
}

func encodeTemplateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(qb.TemplateResponse)
	if resp.Error != nil {
		http.Error(w, fmt.Sprintf("%v", resp.Error), http.StatusInternalServerError)
		return nil
	}

	_, err := w.Write([]byte(resp.Query))
	return err
}

func decodeParseRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	req := struct {
		Query     string            `json:"query"`
		Variables map[string]string `json:"variables"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return qb.ParseRequest{
		Query:     req.Query,
		Variables: req.Variables,
	}, nil
}

func encodeParseResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(qb.ParseResponse)
	if resp.Error != nil {
		http.Error(w, fmt.Sprintf("%v", resp.Error), http.StatusInternalServerError)
		return nil
	}

	json, err := json.Marshal(resp.Queries)
	if err != nil {
		return err
	}

	_, _ = w.Write(json)
	return nil
}
