package endpoint

import (
	"context"
	"encoding/json"

	"mooncamp.com/dgraphtools"
	"mooncamp.com/dgraphtools/gql/extension"
	"mooncamp.com/dgraphtools/render"

	"github.com/go-kit/kit/endpoint"
)

func Query(qh dgraphtools.QueryHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(dgraphtools.QueryRequest)

		q := render.Query{
			Queries:   req.Queries,
			Alias:     req.Alias,
			Variables: req.Variables,
		}
		renderedQuery, err := render.Render(q)
		if err != nil {
			return dgraphtools.QueryResponse{Error: err}, nil
		}

		resp, err := qh.Query(ctx, renderedQuery, req.Variables)
		if err != nil {
			return dgraphtools.QueryResponse{Error: err}, nil
		}

		var data map[string]interface{}
		if err := json.Unmarshal(resp.GetJson(), &data); err != nil {
			return dgraphtools.QueryResponse{Error: err}, nil
		}

		defaultedData := extension.ApplyDefaults(req.Queries, data)
		defaultedJSON, err := json.Marshal(defaultedData)
		if err != nil {
			return dgraphtools.QueryResponse{Error: err}, nil
		}

		return dgraphtools.QueryResponse{Response: defaultedJSON, Error: err}, nil
	}
}
