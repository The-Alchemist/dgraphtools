package render

import (
	"fmt"
	"strings"

	"mooncamp.com/dgraphtools/gql"
)

func renderAttribute(query gql.GraphQuery) string {
	if query.Attr == "" && query.Alias != "" {
		return fmt.Sprintf("%s()", query.Alias)
	}

	if query.Alias != "" {
		return fmt.Sprintf("%s: %s%s%s%s", query.Alias, formatAttribute(query.Attr), renderLangs(query.Langs), renderAttributeArgs(query), renderExpand(query))
	}

	return fmt.Sprintf("%s%s%s%s", formatAttribute(query.Attr), renderLangs(query.Langs), renderAttributeArgs(query), renderExpand(query))
}

func renderLangs(langs []string) string {
	if len(langs) == 0 {
		return ""
	}

	return fmt.Sprintf("@%s", strings.Join(langs, ":"))
}

func renderAttributeArgs(query gql.GraphQuery) string {
	if len(query.Args) == 0 && len(query.Order) == 0 {
		return ""
	}

	res := make([]string, 0, len(query.Args))
	for k, v := range query.Args {
		res = append(res, fmt.Sprintf("%s: %s", k, v))
	}

	for _, e := range query.Order {
		dir := "orderasc"
		if e.Desc {
			dir = "orderdesc"
		}

		res = append(res, fmt.Sprintf("%s: %s%s", dir, e.Attr, renderLangs(e.Langs)))
	}

	return fmt.Sprintf("(%s)", strings.Join(res, ","))
}

func renderExpand(query gql.GraphQuery) string {
	if query.Expand != "" && len(query.NeedsVar) > 0 {
		vars := encodeNeedVar(query.NeedsVar, true)
		return fmt.Sprintf("(%s)", strings.Join(vars, ","))
	}

	if query.Expand != "" {
		return fmt.Sprintf("(%s)", query.Expand)
	}

	return ""
}
