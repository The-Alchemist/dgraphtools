package endpoint

import (
	"context"

	"mooncamp.com/dgraphtools/gql"
	"mooncamp.com/dgraphtools/qb"
	"mooncamp.com/dgraphtools/render"

	dgraphgql "github.com/dgraph-io/dgraph/gql"
	"github.com/go-kit/kit/endpoint"
)

func errFormatter(err error) interface{} {
	return qb.TemplateResponse{Error: err}
}

func queryReader(request interface{}) render.Query {
	req := request.(qb.TemplateRequest)
	return render.Query{
		Queries:   req.Queries,
		Alias:     req.Alias,
		Variables: req.Variables,
	}
}

func NewEndpointSet() qb.EndpointSet {
	var templateEndpoint endpoint.Endpoint
	{
		templateEndpoint = render.TemplateErrorMiddleware(queryReader, errFormatter)(MakeTemplateEndpoint())
	}

	var parseEndpoint endpoint.Endpoint
	{
		parseEndpoint = MakeParseEndpoint()
	}

	return qb.EndpointSet{
		Template: templateEndpoint,
		Parse:    parseEndpoint,
	}
}

func MakeTemplateEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(qb.TemplateRequest)
		q, err := render.Render(render.Query{
			Queries:   req.Queries,
			Alias:     req.Alias,
			Variables: req.Variables,
		})

		return qb.TemplateResponse{
			Query: q,
			Error: err,
		}, nil
	}
}

func MakeParseEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(qb.ParseRequest)
		result, err := dgraphgql.Parse(dgraphgql.Request{Str: req.Query, Variables: req.Variables})
		if err != nil {
			return qb.ParseResponse{Error: err}, nil
		}

		queries := gql.DecodeGraphQueries(result.Query)
		return qb.ParseResponse{
			Error:   nil,
			Queries: queries,
		}, nil
	}
}
