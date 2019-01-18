package dgraphtools

import (
	"context"

	"github.com/dgraph-io/dgo/protos/api"
	"mooncamp.com/dgraphtools/gql"
)

type QueryVerifier interface {
	QueryAllowed(ctx context.Context, queries []gql.GraphQuery, userID int, proofs map[int]gql.GraphQuery) (bool, error)
}

type QueryHandler interface {
	Query(ctx context.Context, q string, vars map[string]string) (*api.Response, error)
}

type QueryRequest struct {
	Identity int

	Queries   []gql.GraphQuery
	Alias     string
	Variables map[string]string
	Proof     map[int]gql.GraphQuery
}

type QueryResponse struct {
	Response []byte
	Error    error
}

type Unauthorized struct{}

func (Unauthorized) Error() string {
	return "unauthorized action"
}
