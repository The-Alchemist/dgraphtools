package extension

import "mooncamp.com/dgraphtools/gql"

func ApplyDefaults(gqs []gql.GraphQuery, resp map[string]interface{}) interface{} {
	res := make(map[string]interface{})
	for _, e := range gqs {
		res[e.Alias] = setDefault(e, resp[e.Alias])
	}

	return res

}

func getNodeName(gq gql.GraphQuery) string {
	if gq.Alias != "" {
		return gq.Alias
	}

	return gq.Attr
}

func setDefault(gq gql.GraphQuery, node interface{}) interface{} {
	// 1. if at []interface, apply gq to all entities
	// 2. if at map[string]interface{}, check for all nodes if defaults are set if empty

	switch t := node.(type) {
	case []interface{}:
		res := make([]interface{}, 0, len(t))
		for _, e := range t {
			res = append(res, setDefault(gq, e))
		}

		return res

	case map[string]interface{}:
		res := make(map[string]interface{})

		for _, e := range gq.Children {
			_, ok := t[getNodeName(e)]
			if ok {
				res[getNodeName(e)] = setDefault(e, t[getNodeName(e)])
				continue
			}

			if e.Default != nil {
				res[getNodeName(e)] = e.Default

			}
		}

		return res
	default:
		return node
	}
}
