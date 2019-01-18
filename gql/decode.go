package gql

import (
	"github.com/dgraph-io/dgraph/gql"
	"github.com/dgraph-io/dgraph/protos/pb"
	"github.com/dgraph-io/dgraph/types"
)

func DecodeGraphQuery(source *gql.GraphQuery) GraphQuery {
	return GraphQuery{
		UID:        source.UID,
		Attr:       source.Attr,
		Langs:      source.Langs,
		Alias:      source.Alias,
		IsCount:    source.IsCount,
		IsInternal: source.IsInternal,
		IsGroupby:  source.IsGroupby,
		Var:        source.Var,
		NeedsVar:   decodeVarContexts(source.NeedsVar),
		Func:       decodeFunc(source.Func),
		Expand:     source.Expand,

		Args:         decodeArgMap(source.Args),
		Order:        decodeOrders(source.Order),
		Children:     DecodeGraphQueries(source.Children),
		Filter:       decodeFilterTree(source.Filter),
		MathExp:      decodeMathTree(source.MathExp),
		Normalize:    source.Normalize,
		Recurse:      source.Recurse,
		Cascade:      source.Cascade,
		IgnoreReflex: source.IgnoreReflex,
		Facets:       decodeFacetParams(source.Facets),
		FacetsFilter: decodeFilterTree(source.FacetsFilter),
		GroupbyAttrs: decodeGroupByAttrs(source.GroupbyAttrs),
		FacetVar:     source.FacetVar,
		FacetOrder:   source.FacetOrder,
		FacetDesc:    source.FacetDesc,
	}
}

func decodeArgMap(source map[string]string) map[string]string {
	if len(source) == 0 {
		return nil
	}

	return source
}

func decodeGroupByAttrs(source []gql.GroupByAttr) []GroupByAttr {
	if len(source) == 0 {
		return nil
	}

	attrs := make([]GroupByAttr, 0, len(source))
	for _, e := range source {
		attrs = append(attrs, decodeGroupByAttr(e))
	}
	return attrs
}

func decodeGroupByAttr(source gql.GroupByAttr) GroupByAttr {
	return GroupByAttr{
		Attr:  source.Attr,
		Alias: source.Alias,
		Langs: source.Langs,
	}
}

func decodeFacetParams(source *pb.FacetParams) *FacetParams {
	if source == nil {
		return nil
	}

	return &FacetParams{
		AllKeys: source.AllKeys,
		Param:   decodeFacetParamSlice(source.Param),
	}
}

func decodeFacetParamSlice(source []*pb.FacetParam) []FacetParam {
	if len(source) == 0 {
		return nil
	}

	params := make([]FacetParam, 0, len(source))
	for _, e := range source {
		params = append(
			params,
			FacetParam{
				Key:   e.Key,
				Alias: e.Alias,
			},
		)
	}
	return params
}

func decodeMathTree(source *gql.MathTree) *MathTree {
	if source == nil {
		return nil
	}

	return &MathTree{
		Fn:    source.Fn,
		Var:   source.Var,
		Const: decodeVal(source.Const),
		Val:   decodeVals(source.Val),
		Child: decodeMathTrees(source.Child),
	}
}

func decodeVal(source types.Val) Val {
	return Val{
		Tid:   TypeID(source.Tid),
		Value: source.Value,
	}
}

func decodeVals(source map[uint64]types.Val) map[uint64]Val {
	if len(source) == 0 {
		return nil
	}

	vals := make(map[uint64]Val, len(source))
	for k, v := range source {
		vals[k] = decodeVal(v)
	}

	return vals
}

func decodeMathTrees(source []*gql.MathTree) []MathTree {
	if len(source) == 0 {
		return nil
	}

	trees := make([]MathTree, 0, len(source))
	for _, e := range source {
		mathTree := decodeMathTree(e)
		trees = append(trees, *mathTree)
	}
	return trees
}

func decodeFilterTree(source *gql.FilterTree) *FilterTree {
	if source == nil {
		return nil
	}

	return &FilterTree{
		Op:    source.Op,
		Child: decodeFilterTrees(source.Child),
		Func:  decodeFunc(source.Func),
	}
}

func decodeFilterTrees(source []*gql.FilterTree) []FilterTree {
	if len(source) == 0 {
		return nil
	}

	filterTrees := make([]FilterTree, 0, len(source))
	for _, e := range source {
		filterTree := decodeFilterTree(e)
		filterTrees = append(filterTrees, *filterTree)
	}
	return filterTrees
}

func DecodeGraphQueries(source []*gql.GraphQuery) []GraphQuery {
	if len(source) == 0 {
		return nil
	}

	graphQueries := make([]GraphQuery, 0, len(source))
	for _, e := range source {
		graphQueries = append(graphQueries, DecodeGraphQuery(e))
	}

	return graphQueries
}

func decodeOrders(source []*pb.Order) []Order {
	if len(source) == 0 {
		return nil
	}

	orders := make([]Order, 0, len(source))
	for _, e := range source {
		orders = append(orders, decodeOrder(e))
	}

	return orders
}

func decodeOrder(source *pb.Order) Order {
	return Order{
		Attr:  source.Attr,
		Desc:  source.Desc,
		Langs: source.Langs,
	}
}

func decodeVarContexts(source []gql.VarContext) []VarContext {
	if len(source) == 0 {
		return nil
	}
	varContexts := make([]VarContext, 0, len(source))
	for _, e := range source {
		varContexts = append(varContexts, decodeVarContext(e))
	}

	return varContexts
}

func decodeVarContext(source gql.VarContext) VarContext {
	return VarContext{
		Name: source.Name,
		Typ:  source.Typ,
	}
}

func decodeArg(source gql.Arg) Arg {
	return Arg{
		Value:        source.Value,
		IsValueVar:   source.IsValueVar,
		IsGraphQLVar: source.IsGraphQLVar,
	}
}

func decodeArgs(source []gql.Arg) []Arg {
	if len(source) == 0 {
		return nil
	}

	args := make([]Arg, 0, len(source))
	for _, e := range source {
		args = append(args, decodeArg(e))
	}

	return args
}

func decodeFunc(source *gql.Function) *Function {
	if source == nil {
		return nil
	}

	return &Function{
		Attr:       source.Attr,
		Lang:       source.Lang,
		Name:       source.Name,
		Args:       decodeArgs(source.Args),
		UID:        source.UID,
		NeedsVar:   decodeVarContexts(source.NeedsVar),
		IsCount:    source.IsCount,
		IsValueVar: source.IsValueVar,
	}
}
