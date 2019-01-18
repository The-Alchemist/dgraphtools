package render

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"mooncamp.com/dgraphtools/gql"

	"github.com/Masterminds/sprig"
)

//go:generate go-bindata -pkg render tmpl/

func formatAttribute(attr string) string {
	u, err := url.Parse(attr)
	if err == nil && u.Scheme != "" {
		return fmt.Sprintf("<%s>", attr)
	}

	return attr
}

func listjoin(property, sep string, data interface{}) (string, error) {
	tmplString := fmt.Sprintf("{{- range . -}}{{ %s }}%s{{- end -}}", property, sep)
	t, err := template.New("").Parse(tmplString)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return strings.TrimSuffix(buf.String(), sep), nil
}

func include(t *template.Template) interface{} {
	return func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := t.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}

func encodeVarContext(val string, typ int) string {
	switch typ {
	case 1: //UID
		i, _ := strconv.Atoi(val)
		return fmt.Sprintf("0x%02x", i)
	default:
		return val
	}
}

func mapifyVarContext(vc []gql.VarContext) map[string]int {
	m := make(map[string]int, len(vc))
	for _, e := range vc {
		m[e.Name] = e.Typ
	}
	return m
}

func math(mt gql.MathTree) string {
	return fmt.Sprintf("math%s", brace(preToInf(mt)))
}

type Query struct {
	Queries   []gql.GraphQuery  `yaml:"queries,omitempty" json:"queries,omitempty"`
	Alias     string            `yaml:"alias,omitempty" json:"alias,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"`
}

func renderGraphqlVariables(variables map[string]string) string {
	if variables == nil {
		return ""
	}

	res := make([]string, 0, len(variables))
	for k, e := range variables {
		res = append(res, fmt.Sprintf("%s: %s", k, e))
	}

	return strings.Join(res, ", ")
}

func splitIndent(s string) (string, string) {
	indent := ""
	for i, c := range s {
		if c == ' ' {
			indent = fmt.Sprintf("%s%c", indent, c)
			continue
		}

		if c == '\t' {
			indent = fmt.Sprintf("%s%c", indent, c)
			continue
		}

		return indent, s[i:]
	}

	return "", s
}

func cleanWhitespace(query string) string {
	space := regexp.MustCompile(`\s+`)

	res := []string{}
	lines := strings.Split(query, "\n")
	for _, l := range lines {
		if strings.Trim(l, " \t") == "" {
			continue
		}

		indent, q := splitIndent(l)
		q = space.ReplaceAllString(q, " ")

		res = append(res, strings.Join([]string{indent, q}, ""))
	}

	return strings.Join(res, "\n")
}

func Render(query Query) (string, error) {
	templ, err := Asset("tmpl/query.tmpl")
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("query"), nil)
	funcMap := template.FuncMap{
		"include":          include(t),
		"listjoin":         listjoin,
		"math":             math,
		"filter":           renderFilter,
		"fn":               renderFunc,
		"internalFn":       renderInternalFunc,
		"shortest":         renderShortest,
		"attribute":        renderAttribute,
		"graphqlVariables": renderGraphqlVariables,
		"groupBy":          renderGroupBy,
		"facets":           renderFacets,
		"facetsFilter":     renderFacetsFilter,
	}

	t.Funcs(funcMap).Funcs(sprig.TxtFuncMap()).Parse(string(templ))
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, query)
	if err != nil {
		return "", err
	}

	return cleanWhitespace(buf.String()), nil
}
