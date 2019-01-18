package render

import (
	"fmt"

	"mooncamp.com/dgraphtools/gql"
)

func renderShortest(query gql.GraphQuery) string {
	return fmt.Sprintf(
		"shortest(from: %s, to:%s, numpaths: %s)",
		query.Args["from"],
		query.Args["to"],
		query.Args["numpaths"],
	)
}
