package render

import (
	"fmt"
	"strings"

	"mooncamp.com/dgraphtools/gql"
)

func renderCheckPwd(query gql.GraphQuery) string {
	return fmt.Sprintf(`checkpwd(%s, "%s")`, formatAttribute(query.Func.Attr), query.Func.Args[0].Value)
}

func renderFunc(query gql.GraphQuery) string {
	if isInternal(query) {
		return ""
	}

	fn := fmt.Sprintf(
		"%s (func: %s",
		query.Alias,
		renderFuncBody(query),
	)

	if len(query.Order) == 0 {
		return fmt.Sprintf("%s%s)", fn, renderFuncArgs(query))
	}

	return fmt.Sprintf("%s, %s%s)", fn, renderOrders(query), renderFuncArgs(query))
}

func renderFuncArgs(query gql.GraphQuery) string {
	if len(query.Args) == 0 {
		return ""
	}

	res := make([]string, 0, len(query.Args))
	for k, v := range query.Args {
		res = append(res, fmt.Sprintf("%s: %s", k, v))
	}

	return fmt.Sprintf(", %s", strings.Join(res, ","))
}

func isInternal(query gql.GraphQuery) bool {
	if query.IsInternal {
		return query.IsInternal
	}

	if len(query.Children) == 0 {
		return true
	}

	return false
}

func renderInternalFunc(query gql.GraphQuery) string {
	if !isInternal(query) {
		return ""
	}

	if query.Func.Name == "checkpwd" {
		return renderCheckPwd(query)
	}

	if query.Alias != "" {
		return fmt.Sprintf("%s: %s", query.Alias, renderFuncBody(query))
	}

	return renderFuncBody(query)
}

func encodeUID(query gql.GraphQuery) (string, bool) {
	if len(query.UID) == 0 {
		return "", false
	}

	res := make([]string, 0, len(query.UID))
	for _, e := range query.UID {
		res = append(res, fmt.Sprintf("0x%02x", e))
	}

	return strings.Join(res, ", "), true
}

func quoteMeta(s string) string {
	return strings.Replace(s, "/", "\\/", -1)
}

func encodeArgs(f *gql.Function) (string, bool) {
	if len(f.Args) == 0 {
		return "", false
	}

	res := make([]string, 0, len(f.Args))
	for _, e := range f.Args {
		if (e == gql.Arg{}) {
			continue
		}

		if f.Name == "regexp" {
			if len(res) > 0 {
				res = append(res, e.Value)
				continue
			}

			res = append(res, fmt.Sprintf("/%s/", quoteMeta(e.Value)))
			continue
		}

		if e.IsValueVar {
			res = append(res, fmt.Sprintf("val(%s)", e.Value))
			continue
		}

		if e.IsGraphQLVar || f.IsValueVar {
			res = append(res, e.Value)
			continue
		}

		res = append(res, fmt.Sprintf(`"%s"`, e.Value))
	}

	if f.Name == "regexp" {
		return strings.Join(res, ""), true
	}

	return strings.Join(res, ", "), true
}

func renderArguments(query gql.GraphQuery) string {
	if uids, ok := encodeUID(query); ok {
		return uids
	}

	if args, ok := encodeArgs(query.Func); ok {
		return args
	}

	if len(query.Func.NeedsVar) != 0 {
		return strings.Join(encodeNeedVar(query.Func.NeedsVar, false), ",")
	}

	return ""
}

func renderFuncBody(query gql.GraphQuery) string {
	arguments := renderArguments(query)

	if query.Attr != "" && arguments == "" {
		return fmt.Sprintf("%s(%s%s)", query.Func.Name, query.Attr, renderLangs(query.Langs))
	}

	if query.Attr != "" {
		return fmt.Sprintf("%s(%s(%s))", query.Func.Name, query.Attr, arguments)
	}

	if query.Func.Attr != "" {
		langs := []string{}
		if query.Func.Lang != "" {
			langs = []string{query.Func.Lang}
		}

		if query.Func.IsCount {
			return fmt.Sprintf("%s(count(%s), %s)", query.Func.Name, query.Func.Attr, arguments)
		}

		if arguments == "" {
			return fmt.Sprintf("%s(%s%s)", query.Func.Name, formatAttribute(query.Func.Attr), renderLangs(langs))
		}

		return fmt.Sprintf("%s(%s%s, %s)", query.Func.Name, formatAttribute(query.Func.Attr), renderLangs(langs), arguments)
	}

	return fmt.Sprintf("%s(%s)", query.Func.Name, arguments)
}

func orderDirection(order gql.Order) string {
	if order.Desc {
		return "orderdesc"
	}

	return "orderasc"
}

func renderOrder(needsVar []gql.VarContext, order gql.Order) string {
	m := mapifyVarContext(needsVar)
	if _, ok := m[order.Attr]; ok {
		return fmt.Sprintf(
			"%s: val(%s)",
			orderDirection(order),
			encodeVarContext(order.Attr, m[order.Attr]),
		)
	}

	return fmt.Sprintf("%s: %s%s", orderDirection(order), order.Attr, renderLangs(order.Langs))
}

func renderOrders(gq gql.GraphQuery) string {
	res := []string{}
	for _, e := range gq.Order {
		res = append(res, renderOrder(gq.NeedsVar, e))
	}

	return strings.Join(res, ",")

}

func renderGroupBy(gq gql.GraphQuery) string {
	if len(gq.GroupbyAttrs) == 0 {
		return ""
	}

	args := []string{}
	for _, e := range gq.GroupbyAttrs {
		alias := ""
		if e.Alias != "" {
			alias = fmt.Sprintf("%s: ", e.Alias)
		}

		args = append(args, fmt.Sprintf("%s%s%s", alias, e.Attr, renderLangs(e.Langs)))
	}

	return fmt.Sprintf("@groupby(%s)", strings.Join(args, ", "))
}
