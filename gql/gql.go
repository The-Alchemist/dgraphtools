package gql

// GraphQuery stores the parsed Query in a tree format. This gets converted to
// pb.y used query.SubGraph before processing the query.
type GraphQuery struct {
	UID        []uint64     `yaml:"uid,omitempty" json:"uid,omitempty"`
	Attr       string       `yaml:"attr,omitempty" json:"attr,omitempty"`
	Langs      []string     `yaml:"langs,omitempty" json:"langs,omitempty"`
	Alias      string       `yaml:"alias,omitempty" json:"alias,omitempty"`
	Default    interface{}  `yaml:"default,omitempty" json:"default,omitempty"`
	IsCount    bool         `yaml:"isCount,omitempty" json:"isCount,omitempty"`
	IsInternal bool         `yaml:"isInternal,omitempty" json:"isInternal,omitempty"`
	IsGroupby  bool         `yaml:"isGroupby,omitempty" json:"isGroupby,omitempty"`
	Var        string       `yaml:"var,omitempty" json:"var,omitempty"`
	NeedsVar   []VarContext `yaml:"needsVar,omitempty" json:"needsVar,omitempty"`
	Func       *Function    `yaml:"func,omitempty" json:"func,omitempty"`
	Expand     string       `yaml:"expand,omitempty" json:"expand,omitempty"` // Which variable to expand with.

	Args map[string]string `yaml:"args,omitempty" json:"args,omitempty"`
	// Query can have multiple sort parameters.
	Order        []Order           `yaml:"order,omitempty" json:"order,omitempty"`
	Children     []GraphQuery      `yaml:"children,omitempty" json:"children,omitempty"`
	Filter       *FilterTree       `yaml:"filter,omitempty" json:"filter,omitempty"`
	MathExp      *MathTree         `yaml:"mathExp,omitempty" json:"mathExp,omitempty"`
	Normalize    bool              `yaml:"normalize,omitempty" json:"normalize,omitempty"`
	Recurse      bool              `yaml:"recurse,omitempty" json:"recurse,omitempty"`
	Cascade      bool              `yaml:"cascade,omitempty" json:"cascade,omitempty"`
	IgnoreReflex bool              `yaml:"ignoreReflex,omitempty" json:"ignoreReflex,omitempty"`
	Facets       *FacetParams      `yaml:"facets,omitempty" json:"facets,omitempty"`
	FacetsFilter *FilterTree       `yaml:"facetsFilter,omitempty" json:"facetsFilter,omitempty"`
	GroupbyAttrs []GroupByAttr     `yaml:"groupbyAttrs,omitempty" json:"groupbyAttrs,omitempty"`
	FacetVar     map[string]string `yaml:"facetVar,omitempty" json:"facetVar,omitempty"`
	FacetOrder   string            `yaml:"facetOrder,omitempty" json:"facetOrder,omitempty"`
	FacetDesc    bool              `yaml:"facetDesc,omitempty" json:"facetDesc,omitempty"`
}

type FacetParams struct {
	AllKeys bool         `yaml:"allKeys,omitempty" json:"allKeys,omitempty"`
	Param   []FacetParam `yaml:"param,omitempty" json:"param,omitempty"`
}

type FacetParam struct {
	Key   string `yaml:"key,omitempty" json:"key,omitempty"`
	Alias string `yaml:"alias,omitempty" json:"alias,omitempty"`
}

type Order struct {
	Attr  string   `yaml:"attr,omitempty" json:"attr,omitempty"`
	Desc  bool     `yaml:"desc,omitempty" json:"desc,omitempty"`
	Langs []string `yaml:"langs,omitempty" json:"langs,omitempty"`
}

type VarContext struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	Typ  int    `yaml:"typ,omitempty" json:"typ,omitempty"` //  1 for UID vars, 2 for value vars
}

type Function struct {
	Attr       string       `yaml:"attr,omitempty" json:"attr,omitempty"`
	Lang       string       `yaml:"lang,omitempty" json:"lang,omitempty"` // language of the attribute value
	Name       string       `yaml:"name,omitempty" json:"name,omitempty"` // Specifies the name of the function.
	Args       []Arg        `yaml:"args,omitempty" json:"args,omitempty"` // Contains the arguments of the function.
	UID        []uint64     `yaml:"uid,omitempty" json:"uid,omitempty"`
	NeedsVar   []VarContext `yaml:"needsVar,omitempty" json:"needsVar,omitempty"`     // If the function requires some variable
	IsCount    bool         `yaml:"isCount,omitempty" json:"isCount,omitempty"`       // gt(count(friends),0)
	IsValueVar bool         `yaml:"isValueVar,omitempty" json:"isValueVar,omitempty"` // eq(val(s), 5)
}

type Arg struct {
	Value        string `yaml:"value,omitempty" json:"value,omitempty"`
	IsValueVar   bool   `yaml:"isValueVar,omitempty" json:"isValueVar,omitempty"` // If argument is val(a)
	IsGraphQLVar bool   `yaml:"isGraphQLVar,omitempty" json:"isGraphQLVar,omitempty"`
}

type FilterTree struct {
	Op    string       `yaml:"op,omitempty" json:"op,omitempty"`
	Child []FilterTree `yaml:"child,omitempty" json:"child,omitempty"`
	Func  *Function    `yaml:"func,omitempty" json:"func,omitempty"`
}

type MathTree struct {
	Fn    string         `yaml:"fn,omitempty" json:"fn,omitempty"`
	Var   string         `yaml:"var,omitempty" json:"var,omitempty"`
	Const Val            `yaml:"const,omitempty" json:"const,omitempty"` // This will always be parsed as a float value
	Val   map[uint64]Val `yaml:"val,omitempty" json:"val,omitempty"`
	Child []MathTree     `yaml:"child,omitempty" json:"child,omitempty"`
}

type TypeID int32

type Val struct {
	Tid   TypeID      `yaml:"tid,omitempty" json:"tid,omitempty"`
	Value interface{} `yaml:"value,omitempty" json:"value,omitempty"`
}

type GroupByAttr struct {
	Attr  string   `yaml:"attr,omitempty" json:"attr,omitempty"`
	Alias string   `yaml:"alias,omitempty" json:"alias,omitempty"`
	Langs []string `yaml:"langs,omitempty" json:"langs,omitempty"`
}
