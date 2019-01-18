package qb

import (
	"mooncamp.com/dgraphtools/gql"

	"github.com/go-kit/kit/endpoint"
)

type EndpointSet struct {
	Template endpoint.Endpoint
	Parse    endpoint.Endpoint
}

type TemplateRequest struct {
	Queries   []gql.GraphQuery
	Alias     string
	Variables map[string]string
}

type TemplateResponse struct {
	Query string
	Error error
}

type ParseRequest struct {
	Query     string
	Variables map[string]string
}

type ParseResponse struct {
	Error   error
	Queries []gql.GraphQuery
}
