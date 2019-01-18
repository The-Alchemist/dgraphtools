package render

import (
	"fmt"
	"strings"

	"mooncamp.com/dgraphtools/gql"
)

func encodeFacetVar(query gql.GraphQuery) []string {
	if query.Facets == nil {
		return []string{}
	}

	res := []string{}
	for k, v := range query.FacetVar {
		if k == query.FacetOrder {
			continue
		}
		res = append(res, fmt.Sprintf("%s as %s", v, k))
	}

	return res
}

func encodeFacetOrder(query gql.GraphQuery) []string {
	if query.Facets == nil {
		return []string{}
	}

	if query.FacetOrder == "" {
		return []string{}
	}

	dir := "orderasc"
	if query.FacetDesc {
		dir = "orderdesc"
	}

	orderVar := query.FacetOrder
	if v, ok := query.FacetVar[query.FacetOrder]; ok {
		orderVar = fmt.Sprintf("%s as %s", v, query.FacetOrder)
	}

	return []string{fmt.Sprintf("%s: %s", dir, orderVar)}
}

func encodeFacetParams(query gql.GraphQuery) []string {
	if query.Facets == nil {
		return []string{}
	}

	res := []string{}
	for _, e := range query.Facets.Param {
		if _, ok := query.FacetVar[e.Key]; ok {
			continue
		}

		if e.Key == query.FacetOrder {
			continue
		}

		alias := ""
		if e.Alias != "" {
			alias = fmt.Sprintf("%s: ", e.Alias)
		}

		res = append(res, fmt.Sprintf("%s%s", alias, e.Key))
	}

	return res
}

func renderFacetsFilter(query gql.GraphQuery) string {
	if query.FacetsFilter == nil {
		return ""
	}

	return fmt.Sprintf("@facets(%s)", renderFilterTree(*query.FacetsFilter))
}

func renderFacets(query gql.GraphQuery) string {
	if query.Facets == nil {
		return ""
	}

	arguments := []string{}
	arguments = append(arguments, encodeFacetOrder(query)...)
	arguments = append(arguments, encodeFacetVar(query)...)
	arguments = append(arguments, encodeFacetParams(query)...)

	if query.Facets != nil && query.Facets.AllKeys {
		return "@facets"
	}

	if len(arguments) == 0 {
		return "@facets() { }"
	}

	return fmt.Sprintf("@facets(%s)", strings.Join(arguments, ", "))
}
