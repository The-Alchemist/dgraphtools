package render

import (
	"fmt"
	"strings"

	"mooncamp.com/dgraphtools/gql"
)

func isOperand(mt gql.MathTree) bool {
	return mt.Fn != ""
}

func mathValue(mt gql.MathTree) string {
	if (mt.Const == gql.Val{}) {
		return mt.Var
	}

	return fmt.Sprintf("%v", mt.Const.Value)
}

func brace(fn string) string {
	if fn[0] == '(' && fn[len(fn)-1] == ')' {
		return fn
	}

	return fmt.Sprintf("(%s)", fn)
}

func function(function string) (string, bool) {
	functions := map[string]string{
		"u-":   "-%s",
		"ln":   "ln(%s)",
		"exp":  "exp(%s)",
		"max":  "max(%s)",
		"sqrt": "sqrt(%s)",
		"cond": "cond(%s)",
	}

	f, ok := functions[function]
	return f, ok
}

func encodeChilds(childs []gql.MathTree) []string {
	res := []string{}
	for _, e := range childs {
		res = append(res, preToInf(e))
	}

	return res
}

func preToInf(mt gql.MathTree) string {
	if !isOperand(mt) {
		return mathValue(mt)
	}

	if f, ok := function(mt.Fn); ok {
		return fmt.Sprintf(f, strings.Join(encodeChilds(mt.Child), ","))
	}

	return fmt.Sprintf("(%s)", strings.Join(encodeChilds(mt.Child), mt.Fn))
}
