We have been using Dgraph for over a year at Mooncamp and while loving
the product we found some issues that we needed to write additional
tools for.

## Query Rendering

DgraphTools allows you to render GraphQL+- queries based on its data
representation. The types for the data representation are mostly
copied from the `github.com/dgraph-io/dgraph/gql` package.

Start testing the translation between the data and string
representation of the query by starting the `querybuilder` tool.

```bash
$ go get -u mooncamp.com/dgraphtools/cmd/dgraphtools
$ dgraphtools querybuilder
```

We find that using the data representation over the string form gives
us several advantages. When defining queries, we realized that
composing queries can speed up development and reduce
maintenance. However, in the string representation queries are hard to
compose, as string concatenations are messy and error prone. When the
query is in the data form compositions are substantially easier. Using
a typed language like Go brings additional safety, as the query
structure can be verified using the language's type system.

### GraphQL+- Query
```graphql
query {
  bladerunner (func: eq(name@en, "Blade Runner")) {
	name@en
	initial_release_date
   }
}
```

### Data Representation in Go
```Go
gql.GraphQuery{
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
```

## Moving Query Ownership to the Frontend

While it is standard practice for GraphQL clients to take ownership of
the queries, with Dgraph this is usually not an option as there is no
safe way for the backend to implement access restriction based on the
string representation of the query. However, moving the queries to the
client offers great flexibility and should therefore be the preferable
option. Using the data representation the backend now has access to
every aspect of the query.

This project comes with an example implementation of how minimal
access restriction could be implemented. The approach assumes that any
identity accessing the graph is stored in the graph itself and that
there are no edges connecting the network of the identity itself and
the network the identity doesn't have access to.

![alt text](/doc/graph_example.png "example network")

In the example the approach assumes that there is no connection
between the networks {0,1,2,3,4} and {5,6,7,8,9,10,12}. This is often
the case when you implement multi tenancy with a single database
instance. For example a user John Doe has access to the data of his
own company Goggles Inc, but no access to the other company Amazing
Corp. At the same time there is no connection between the data of
Goggles Inc and Amazing Corp.

The approach works by only allowing `uid` as the root function, where
the corresponding uid is the users identity. Additionally, if the uid
is not the users identity a separate query can be provided that proofs
the connection between the input uid and the users identity.

## Dgraph Extensions

When taking control over the query language we have the ability to
extend the language itself. For example we sometimes want to set a
default value for a non found relation in dgraph. Therefore the
`GraphQuery` type was enhanced with a `Default interface{}`
property. Extensions are applied on the response of dgraph. Another
extension is planned, which will allow setting the cardinality between
two nodes. This solves the issue of representing one-to-one
relationships within dgraph.

## Don't trust us

Although the rendering code is pretty well tested using the actual
dgraph tests as an inspiration there could still be bugs which could
potentially lead to security issues. Luckily we can use the dgraph
query parser as a source of truth, meaning any bug in the rendering
code will be detected right away. We are actually testing any query
coming from the client for potential issues and return 500 if the
query couldn't be verified.

## Example Application

Checkout `example/main.go` for an example usage of all components
described above in a web application setting. The example makes use of
the gokit architecture, so if you are confused about what is going on
checkout [gokit.io](https://gokit.io).
