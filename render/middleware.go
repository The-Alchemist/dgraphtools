package render

import (
	"context"
	"fmt"

	"mooncamp.com/dgraphtools/gql"

	dgraphgql "github.com/dgraph-io/dgraph/gql"

	"github.com/go-kit/kit/endpoint"
	"github.com/stretchr/testify/assert"
)

func nullDefault(gq gql.GraphQuery) gql.GraphQuery {
	gq.Default = nil

	if gq.Children == nil {
		return gq
	}

	children := make([]gql.GraphQuery, 0, len(gq.Children))
	for _, e := range gq.Children {
		children = append(children, nullDefault(e))
	}
	gq.Children = children
	return gq
}

func TemplateErrorMiddleware(queryReader func(request interface{}) Query, errFormatter func(err error) interface{}) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			query := queryReader(request)
			q, err := Render(query)
			if err != nil {
				return errFormatter(err), nil
			}

			gqlVariables := make(map[string]string, len(query.Variables))
			for k := range query.Variables {
				gqlVariables[k] = k
			}

			parseResult, err := dgraphgql.Parse(dgraphgql.Request{Str: q, Variables: gqlVariables})
			if err != nil {
				return errFormatter(err), nil
			}

			actual := make([]gql.GraphQuery, len(query.Queries))
			copy(actual, query.Queries)
			for i := range actual {
				actual[i] = nullDefault(actual[i])
			}

			expected := gql.DecodeGraphQueries(parseResult.Query)

			if !assert.ObjectsAreEqual(expected, actual) {
				d := diff(expected, actual)
				return errFormatter(fmt.Errorf("parsing difference: %s", d)), nil
			}

			return next(ctx, request)
		}
	}
}
