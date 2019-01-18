package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"mooncamp.com/dgraphtools/gql"
	"mooncamp.com/dgraphtools/render"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type handler struct {
	dg *dgo.Dgraph
}

func (h *handler) Query(ctx context.Context, q string, vars map[string]string) (*api.Response, error) {
	return h.dg.NewTxn().QueryWithVars(ctx, q, vars)
}

type person struct {
	Name    string   `json:"name,omitempty"`
	Friends []person `json:"friends,omitempty"`
}

func Test_defaults_are_set(t *testing.T) {
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("grpc dial: %v", err)
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	defer dg.Alter(context.Background(), &api.Operation{DropAll: true})

	filled := person{
		Name:    "harry",
		Friends: []person{{Name: "peter"}},
	}

	empty := person{
		Name: "anna",
	}

	filledJSON, _ := json.Marshal(filled)
	emptyJSON, _ := json.Marshal(empty)

	res, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: filledJSON, CommitNow: true})
	if err != nil {
		t.Fatalf("insert filled: %v", err)
	}
	uidFilled, _ := strconv.ParseInt(res.GetUids()["blank-0"], 0, 64)

	res, err = dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: emptyJSON, CommitNow: true})
	if err != nil {
		t.Fatalf("insert empty: %v", err)
	}
	uidEmpty, _ := strconv.ParseInt(res.GetUids()["blank-0"], 0, 64)

	graphQuery := gql.GraphQuery{
		Alias: "persons",
		UID:   []uint64{uint64(uidFilled), uint64(uidEmpty)},
		Func:  &gql.Function{Name: "uid"},
		Children: []gql.GraphQuery{
			{Attr: "name"},
			{Attr: "friends", Default: []person{}, Children: []gql.GraphQuery{{Attr: "name"}}},
		},
	}

	q, err := render.Render(render.Query{Queries: []gql.GraphQuery{graphQuery}})
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	resp, err := (&handler{dg: dg}).Query(context.Background(), q, nil)
	if err != nil {
		t.Fatalf("query: %v", err)
	}

	queryResp := make(map[string]interface{})
	_ = json.Unmarshal(resp.GetJson(), &queryResp)

	defaulted := ApplyDefaults([]gql.GraphQuery{graphQuery}, queryResp)
	for _, e := range defaulted.(map[string]interface{})[graphQuery.Alias].([]interface{}) {
		m := e.(map[string]interface{})
		_, ok := m["friends"]
		require.True(t, ok, fmt.Sprintf("%s's friends undefined", m["name"]))
	}
}
