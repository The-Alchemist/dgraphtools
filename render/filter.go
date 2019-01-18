package render

import (
	"fmt"
	"strings"

	"mooncamp.com/dgraphtools/gql"
)

func renderFilter(filter *gql.FilterTree) string {
	if filter == nil {
		return ""
	}

	if filter.Op != "" {
		return fmt.Sprintf("@filter(%s)", renderFilterTree(*filter))
	}

	if filter.Func == nil {
		return ""
	}

	if len(filter.Func.UID) != 0 && len(filter.Func.Args) == 0 {
		uids := make([]string, 0, len(filter.Func.UID))
		for _, e := range filter.Func.UID {
			uids = append(uids, fmt.Sprintf("%d", e))
		}

		return fmt.Sprintf("@filter(%s(%s))", filter.Func.Name, strings.Join(uids, ", "))
	}

	if filter.Func.Attr != "" && len(filter.Func.Args) == 0 {
		return fmt.Sprintf("@filter(%s(%s))", filter.Func.Name, filter.Func.Attr)
	}

	if len(filter.Func.Args) == 0 {
		return fmt.Sprintf("@filter(%s(%s))", filter.Func.Name, strings.Join(encodeNeedVar(filter.Func.NeedsVar, true), ","))
	}

	if filter.Func.Attr != "" {
		arguments, _ := encodeArgs(filter.Func)
		attr := formatAttribute(filter.Func.Attr)
		if filter.Func.IsValueVar {
			attr = fmt.Sprintf("val(%s)", formatAttribute(filter.Func.Attr))
		}

		if filter.Func.IsCount {
			return fmt.Sprintf("@filter(%s(count(%s), %s))", filter.Func.Name, attr, arguments)
		}

		langs := []string{}
		if filter.Func.Lang != "" {
			langs = []string{filter.Func.Lang}
		}

		return fmt.Sprintf("@filter(%s(%s%s, %s))", filter.Func.Name, attr, renderLangs(langs), arguments)
	}

	return fmt.Sprintf(
		"@filter(%s(%s, %s))",
		filter.Func.Name,
		strings.Join(encodeNeedVar(filter.Func.NeedsVar, true), ","),
		strings.Join(encodeFilterArgs(filter.Func.Args), ","),
	)

}

func encodeFilterArgs(args []gql.Arg) []string {
	res := make([]string, 0, len(args))
	for _, e := range args {
		res = append(res, encodeFilterArg(e))
	}
	return res
}

func encodeFilterArg(arg gql.Arg) string {
	return arg.Value
}

func encodeNeedVar(varContexts []gql.VarContext, annotate bool) []string {
	s := make([]string, 0, len(varContexts))
	for _, e := range varContexts {
		if e.Typ > 1 && annotate {
			s = append(s, fmt.Sprintf("val(%s)", e.Name))
			continue
		}
		s = append(s, e.Name)
	}

	return s
}

func renderFilterfunc(fn *gql.Function) string {
	args := []string{fn.Attr}
	for _, e := range fn.Args {
		args = append(args, fmt.Sprintf(`"%s"`, e.Value))
	}

	return fmt.Sprintf(`%s(%s)`, fn.Name, strings.Join(args, ", "))
}

func renderSingleArgFunction(tree gql.FilterTree) (string, bool) {
	for _, e := range []string{"not", "eq"} {
		args := []string{}
		for _, e := range tree.Child {
			args = append(args, renderFilterTree(e))
		}

		if e == tree.Op {
			return fmt.Sprintf("%s (%s)", tree.Op, strings.Join(args, ", ")), true
		}
	}

	return "", false
}

func encodeFilterTreeChilds(childs []gql.FilterTree) []string {
	res := []string{}
	for _, e := range childs {
		res = append(res, renderFilterTree(e))
	}
	return res
}

func renderFilterTree(tree gql.FilterTree) string {
	if tree.Func != nil {
		return renderFilterfunc(tree.Func)
	}

	if n, ok := renderSingleArgFunction(tree); ok {
		return n
	}

	return fmt.Sprintf("(%s)", strings.Join(encodeFilterTreeChilds(tree.Child), fmt.Sprintf(" %s ", tree.Op)))
}
