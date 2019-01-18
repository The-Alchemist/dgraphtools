package proof

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"mooncamp.com/dgraphtools/gql"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	dgraphHost string
	dgraphPort string
)

func init() {
	dgraphHost = "localhost"
	dgraphPort = "9080"
}

type User struct {
	UID  string `json:"uid"`
	Name string `json:"name"`

	Company Company `json:"user.company"`
}

type Company struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type handler struct {
	dg *dgo.Dgraph
}

func (h *handler) Query(ctx context.Context, q string, vars map[string]string) (*api.Response, error) {
	return h.dg.NewTxn().QueryWithVars(ctx, q, vars)
}

func data(t *testing.T, dg *dgo.Dgraph) []User {
	companies := []Company{
		{Name: "CompanyA"},
		{Name: "CompanyB"},
	}

	for i, e := range companies {
		js, err := json.Marshal(e)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}

		op := &api.Mutation{
			SetJson:   js,
			CommitNow: true,
		}

		assigned, err := dg.NewTxn().Mutate(context.Background(), op)
		if err != nil {
			t.Fatalf("write to db: %v", err)
		}

		companies[i].UID = assigned.Uids["blank-0"]
	}

	users := []User{
		{Name: "UserA", Company: companies[0]},
		{Name: "UserB", Company: companies[1]},
	}

	for i, e := range users {
		js, err := json.Marshal(e)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}

		op := &api.Mutation{
			SetJson:   js,
			CommitNow: true,
		}

		assigned, err := dg.NewTxn().Mutate(context.Background(), op)
		if err != nil {
			t.Fatalf("write to db: %v", err)
		}

		users[i].UID = assigned.Uids["blank-0"]
	}

	return users
}

func clear(t *testing.T, dg *dgo.Dgraph) {
	op := api.Operation{
		DropAll: true,
	}
	if err := dg.Alter(context.Background(), &op); err != nil {
		t.Fatalf("clear database: %v", err)
	}
}

func parseID(t *testing.T, id string) uint64 {
	n, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		t.Fatalf("parse int: %v", err)
	}

	return uint64(n)
}

func Test_query_allowed(t *testing.T) {
	conn, err := grpc.Dial(dgraphHost+":"+dgraphPort, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("grpc dial: %v", err)
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)
	defer clear(t, dg)

	users := data(t, dg)

	cases := []struct {
		name     string
		queries  []gql.GraphQuery
		identity int
		proofs   map[int]gql.GraphQuery
		allowed  bool
	}{
		{
			name:     "disallow non uid functions as root",
			queries:  []gql.GraphQuery{{Func: &gql.Function{Name: "eq"}}},
			identity: 0,
			proofs:   nil,
			allowed:  false,
		},
		{
			name:     "disallow uid inputs without proof query",
			queries:  []gql.GraphQuery{{UID: []uint64{0}, Func: &gql.Function{Name: "uid"}}},
			identity: 1,
			proofs:   nil,
			allowed:  false,
		},
		{
			name:     "allow uid inputs with proof query",
			queries:  []gql.GraphQuery{{UID: []uint64{parseID(t, users[0].Company.UID)}, Func: &gql.Function{Name: "uid"}}},
			identity: int(parseID(t, users[0].UID)),
			proofs: map[int]gql.GraphQuery{
				int(parseID(t, users[0].Company.UID)): gql.GraphQuery{
					Alias: "proof",
					UID:   []uint64{parseID(t, users[0].UID)},
					Func:  &gql.Function{Name: "uid"},
					Children: []gql.GraphQuery{
						{
							Attr:     "user.company",
							Alias:    "proof",
							Children: []gql.GraphQuery{{Attr: "uid", Alias: "proof"}},
							Filter: &gql.FilterTree{
								Func: &gql.Function{Name: "uid", UID: []uint64{parseID(t, users[0].Company.UID)}},
							},
						},
					},
				},
			},
			allowed: true,
		},
		{
			name:     "disallow uid inputs with misleading proof query",
			queries:  []gql.GraphQuery{{UID: []uint64{parseID(t, users[1].Company.UID)}, Func: &gql.Function{Name: "uid"}}},
			identity: int(parseID(t, users[0].UID)),
			proofs: map[int]gql.GraphQuery{
				int(parseID(t, users[0].Company.UID)): gql.GraphQuery{
					Alias: "proof",
					UID:   []uint64{parseID(t, users[0].UID)},
					Func:  &gql.Function{Name: "uid"},
					Children: []gql.GraphQuery{
						{
							Attr:     "user.company",
							Alias:    "proof",
							Children: []gql.GraphQuery{{Attr: "uid", Alias: "proof"}},
							Filter: &gql.FilterTree{
								Func: &gql.Function{Name: "uid", UID: []uint64{parseID(t, users[1].Company.UID)}},
							},
						},
					},
				},
			},
			allowed: false,
		},
	}

	proof := &Proof{
		QueryHandler: &handler{dg: dg},
	}

	for _, e := range cases {
		allowed, queries, identity, proofs := e.allowed, e.queries, e.identity, e.proofs
		t.Run(e.name, func(t *testing.T) {
			ok, err := proof.QueryAllowed(context.Background(), queries, identity, proofs)
			if err != nil {
				t.Fatalf("check query: %v", err)
			}

			assert.Equal(t, allowed, ok)
		})
	}
}
