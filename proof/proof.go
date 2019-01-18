package proof

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"mooncamp.com/dgraphtools"
	"mooncamp.com/dgraphtools/gql"
	"mooncamp.com/dgraphtools/render"

	"github.com/go-kit/kit/endpoint"
)

type Proof struct {
	dgraphtools.QueryHandler
}

func (p *Proof) QueryAllowed(ctx context.Context, queries []gql.GraphQuery, identity int, proofs map[int]gql.GraphQuery) (bool, error) {
	for _, e := range queries {
		if e.Func == nil {
			return false, nil
		}

		if e.Func.Name != "uid" {
			return false, nil
		}

		for _, uid := range e.UID {
			if int(uid) == identity {
				continue
			}

			proofQuery, ok := proofs[int(uid)]
			if !ok {
				return false, nil
			}

			ok, err := p.hasPath(ctx, int(uid), identity, proofQuery)
			if err != nil {
				return false, err
			}

			if !ok {
				return false, nil
			}
		}
	}

	return true, nil
}

func (p *Proof) hasPath(ctx context.Context, uid, identity int, proofQuery gql.GraphQuery) (bool, error) {
	if proofQuery.Func == nil {
		return false, nil
	}

	if proofQuery.Func.Name != "uid" {
		return false, nil
	}

	if int(proofQuery.UID[0]) != identity {
		return false, nil
	}

	query, err := render.Render(render.Query{Queries: []gql.GraphQuery{proofQuery}})
	if err != nil {
		return false, err
	}

	resp, err := p.QueryHandler.Query(ctx, query, map[string]string{})
	if err != nil {
		return false, err
	}

	proof := make(map[string]interface{})
	if err := json.Unmarshal(resp.GetJson(), &proof); err != nil {
		return false, err
	}

	proofUID, err := followProof(proof)
	if err != nil {
		return false, err
	}

	puid, err := strconv.ParseInt(proofUID, 0, 64)
	if err != nil {
		return false, err
	}

	return int(puid) == uid, nil
}

func followProof(node interface{}) (string, error) {
	if next, ok := node.(string); ok {
		return next, nil
	}

	if next, ok := node.(map[string]interface{}); ok {
		return followProof(next["proof"])
	}

	if next, ok := node.([]interface{}); ok {
		if len(next) == 0 {
			return "", fmt.Errorf("incorrect proof")
		}

		return followProof(next[0])
	}

	return "", fmt.Errorf("incorrect proof")
}

func Middleware(verifier dgraphtools.QueryVerifier) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(dgraphtools.QueryRequest)

			ok, err := verifier.QueryAllowed(ctx, req.Queries, req.Identity, req.Proof)
			if err != nil {
				return dgraphtools.QueryResponse{Error: err}, nil
			}

			if !ok {
				return dgraphtools.QueryResponse{Error: dgraphtools.Unauthorized{}}, nil
			}

			return next(ctx, request)
		}
	}
}
