package gql

import (
	"github.com/dgraph-io/dgraph/gql"
	"github.com/dgraph-io/dgraph/protos/pb"
	"github.com/dgraph-io/dgraph/types"
)

func EncodeGraphQuery(source GraphQuery) *gql.GraphQuery {
	return &gql.GraphQuery{
		UID:        source.UID,
		Attr:       source.Attr,
		Langs:      source.Langs,
		Alias:      source.Alias,
		IsCount:    source.IsCount,
		IsInternal: source.IsInternal,
		IsGroupby:  source.IsGroupby,
		Var:        source.Var,
		NeedsVar:   encodeVarContexts(source.NeedsVar),
		Func:       encodeFunc(source.Func),
		Expand:     source.Expand,

		Args:         source.Args,
		Order:        encodeOrders(source.Order),
		Children:     EncodeGraphQueries(source.Children),
		Filter:       encodeFilterTree(source.Filter),
		MathExp:      encodeMathTree(source.MathExp),
		Normalize:    source.Normalize,
		Recurse:      source.Recurse,
		Cascade:      source.Cascade,
		IgnoreReflex: source.IgnoreReflex,
		Facets:       encodeFacetParams(source.Facets),
		FacetsFilter: encodeFilterTree(source.FacetsFilter),
		GroupbyAttrs: encodeGroupByAttrs(source.GroupbyAttrs),
		FacetVar:     source.FacetVar,
		FacetOrder:   source.FacetOrder,
		FacetDesc:    source.FacetDesc,
	}
}

func encodeGroupByAttrs(source []GroupByAttr) []gql.GroupByAttr {
	attrs := make([]gql.GroupByAttr, 0, len(source))
	for _, e := range source {
		attrs = append(attrs, encodeGroupByAttr(e))
	}
	return attrs
}

func encodeGroupByAttr(source GroupByAttr) gql.GroupByAttr {
	return gql.GroupByAttr{
		Attr:  source.Attr,
		Alias: source.Alias,
		Langs: source.Langs,
	}
}

func encodeFacetParams(source *FacetParams) *pb.FacetParams {
	return &pb.FacetParams{
		AllKeys: source.AllKeys,
		Param:   encodeFacetParamSlice(source.Param),
	}
}

func encodeFacetParamSlice(source []FacetParam) []*pb.FacetParam {
	params := make([]*pb.FacetParam, 0, len(source))
	for _, e := range source {
		params = append(
			params,
			&pb.FacetParam{
				Key:   e.Key,
				Alias: e.Alias,
			},
		)
	}
	return params
}

func encodeMathTree(source *MathTree) *gql.MathTree {
	return &gql.MathTree{
		Fn:    source.Fn,
		Var:   source.Var,
		Const: encodeVal(source.Const),
		Val:   encodeVals(source.Val),
		Child: encodeMathTrees(source.Child),
	}
}

func encodeVal(source Val) types.Val {
	return types.Val{
		Tid:   types.TypeID(source.Tid),
		Value: source.Value,
	}
}

func encodeVals(source map[uint64]Val) map[uint64]types.Val {
	vals := make(map[uint64]types.Val, len(source))
	for k, v := range source {
		vals[k] = encodeVal(v)
	}

	return vals
}

func encodeMathTrees(source []MathTree) []*gql.MathTree {
	trees := make([]*gql.MathTree, 0, len(source))
	for _, e := range source {
		trees = append(trees, encodeMathTree(&e))
	}
	return trees
}

func encodeFilterTree(source *FilterTree) *gql.FilterTree {
	return &gql.FilterTree{
		Op:    source.Op,
		Child: encodeFilterTrees(source.Child),
		Func:  encodeFunc(source.Func),
	}
}

func encodeFilterTrees(source []FilterTree) []*gql.FilterTree {
	filterTrees := make([]*gql.FilterTree, 0, len(source))
	for _, e := range source {
		filterTrees = append(filterTrees, encodeFilterTree(&e))
	}
	return filterTrees
}

func EncodeGraphQueries(source []GraphQuery) []*gql.GraphQuery {
	graphQueries := make([]*gql.GraphQuery, 0, len(source))
	for _, e := range source {
		graphQueries = append(graphQueries, EncodeGraphQuery(e))
	}

	return graphQueries
}

func encodeOrders(source []Order) []*pb.Order {
	orders := make([]*pb.Order, 0, len(source))
	for _, e := range source {
		orders = append(orders, encodeOrder(e))
	}

	return orders
}

func encodeOrder(source Order) *pb.Order {
	return &pb.Order{
		Attr:  source.Attr,
		Desc:  source.Desc,
		Langs: source.Langs,
	}
}

func encodeVarContexts(source []VarContext) []gql.VarContext {
	varContexts := make([]gql.VarContext, 0, len(source))
	for _, e := range source {
		varContexts = append(varContexts, encodeVarContext(e))
	}

	return varContexts
}

func encodeVarContext(source VarContext) gql.VarContext {
	return gql.VarContext{
		Name: source.Name,
		Typ:  source.Typ,
	}
}

func encodeArg(source Arg) gql.Arg {
	return gql.Arg{
		Value:        source.Value,
		IsValueVar:   source.IsValueVar,
		IsGraphQLVar: source.IsGraphQLVar,
	}
}

func encodeArgs(source []Arg) []gql.Arg {
	args := make([]gql.Arg, 0, len(source))
	for _, e := range source {
		args = append(args, encodeArg(e))
	}

	return args
}

func encodeFunc(source *Function) *gql.Function {
	return &gql.Function{
		Attr:       source.Attr,
		Lang:       source.Lang,
		Name:       source.Name,
		Args:       encodeArgs(source.Args),
		UID:        source.UID,
		NeedsVar:   encodeVarContexts(source.NeedsVar),
		IsCount:    source.IsCount,
		IsValueVar: source.IsValueVar,
	}
}
