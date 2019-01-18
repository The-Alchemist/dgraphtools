package doc

import "mooncamp.com/dgraphtools/gql"

var _ = gql.GraphQuery{
	Alias: "bladerunner",
	Func: &gql.Function{
		Attr: "name",
		Lang: "en",
		Name: "eq",
		Args: []gql.Arg{
			{Value: "Blade Runner"},
		},
	},
	Children: []gql.GraphQuery{
		{Attr: "Name", Langs: []string{"en"}},
		{Attr: "initial_release_date"},
	},
}
