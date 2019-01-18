package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"mooncamp.com/dgraphtools"
	"mooncamp.com/dgraphtools/endpoint"
	"mooncamp.com/dgraphtools/gql"
	"mooncamp.com/dgraphtools/proof"
	"mooncamp.com/dgraphtools/render"

	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	gokitendpoint "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type queryHandler struct {
	dg *dgo.Dgraph
}

func (h *queryHandler) Query(ctx context.Context, q string, vars map[string]string) (*api.Response, error) {
	return h.dg.NewTxn().QueryWithVars(ctx, q, vars)
}

func queryReader(request interface{}) render.Query {
	req := request.(dgraphtools.QueryRequest)

	return render.Query{
		Queries:   req.Queries,
		Alias:     req.Alias,
		Variables: req.Variables,
	}
}

func errFormatter(err error) interface{} {
	return dgraphtools.QueryResponse{Error: err}
}

func main() {
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to dgraph: %v", err)
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	handler := mux.NewRouter()

	verifier := &proof.Proof{&queryHandler{dg: dg}}

	var queryEndpoint gokitendpoint.Endpoint
	{
		queryEndpoint = endpoint.Query(&queryHandler{dg: dg})
		queryEndpoint = proof.Middleware(verifier)(queryEndpoint)
		queryEndpoint = render.TemplateErrorMiddleware(queryReader, errFormatter)(queryEndpoint)
	}

	queryHandler := httptransport.NewServer(
		queryEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			userID, _ := r.Cookie("userid")

			req := struct {
				Queries   []gql.GraphQuery       `json:"queries"`
				Alias     string                 `json:"alias"`
				Variables map[string]string      `json:"variables"`
				Proof     map[int]gql.GraphQuery `json:"proof"`
			}{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return nil, err
			}

			id, err := strconv.ParseInt(userID.Value, 0, 64)
			if err != nil {
				return nil, err
			}

			return dgraphtools.QueryRequest{
				Queries:   req.Queries,
				Alias:     req.Alias,
				Variables: req.Variables,
				Proof:     req.Proof,
				Identity:  int(id),
			}, nil

		},
		func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
			resp := response.(dgraphtools.QueryResponse)

			if resp.Error != nil {
				http.Error(w, fmt.Sprintf("%v", resp.Error), http.StatusInternalServerError)
				return nil
			}

			_, err := w.Write(resp.Response)
			return err
		},
	)

	handler.Handle("/query", queryHandler)
}
