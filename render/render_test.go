package render

import (
	"testing"

	"mooncamp.com/dgraphtools/gql"

	dgraphgql "github.com/dgraph-io/dgraph/gql"
	"github.com/stretchr/testify/require"
)

func TestListJoin(t *testing.T) {
	res, err := listjoin(
		".Foo",
		",",
		[]map[string]string{
			{"Foo": "Bar"},
			{"Foo": "Baz"},
		},
	)
	if err != nil {
		t.Fatalf("listjoin: %v", err)
	}

	require.Equal(t, "Bar,Baz", res)
}

func TestRender(t *testing.T) {
	table := []struct {
		name      string
		query     string
		variables map[string]string
		alias     string
	}{
		{
			name: "filter within count",
			query: `
			       query {
				 userCount (func: uid(0x186b4)) {
				   user.company {
				     count(~user.company @filter(regexp(user.fullName, /joschka/i)))
				   }
				 }
			       }
			       `,
		},
		{
			name: "dgraph docs 1",
			query: `
			       {
				 bladerunner(func: uid(0x107b2c)) {
				   uid
				   name@en
				   initial_release_date
				   netflix_id
				 }
			       }
			       `,
		},
		{
			name: "order with lang",
			query: `
			       {
				       me(func: uid(0x1), orderasc: name@en:fr:., orderdesc: lastname@ci, orderasc: salary) {
					       name
				       }
			       }
		       `,
		},
		{
			name: "order by var and pred",
			query: `{
				var(func: uid(0x0a)) {
					friends (orderasc: name, orderdesc: genre) {
						name
					}
				}

			}`,
		},
		{
			name: "lang with dash",
			query: `{
				q(func: uid(1)) {
					text@en-us
				}
			}`,
		},
		{
			name: "eq arg with dollar",
			query: `
			       {
				       ab(func: eq(name@en, "$pringfield (or, How)")) {
					       uid
				       }
			       }
			       `,
		},
		{
			name: "order 2",
			query: `
			       {
				       me(func: uid(0x01)) {
					       friend(orderasc: alias, orderdesc: name) @filter(lt(alias, "Pat")) {
						       alias
					       }
				       }
			       }
		       `,
		},
		{
			name: "order 1",
			query: `
			       {
				       me(func: uid(1), orderdesc: name, orderasc: age) {
					       name
				       }
			       }
		       `,
		},
		{
			name: "agg root 1",
			query: `
			       {
				       var(func: anyofterms(name, "Rick Michonne Andrea")) {
					       a as age
				       }

				       me() {
					       sum(val(a))
					       avg(val(a))
				       }
			       }
		       `,
		},
		{
			name: "filter uid",
			query: `
			       {
				       me(func: uid(1, 3 , 5, 7)) @filter(uid(3, 7)) {
					       name
				       }
			       }
			       `,
		},
		{
			name: "eq arg 2",
			query: `
			       {
				       me(func: eq(age, [1, 20])) @filter(eq(name, ["Andrea", "Bob"])) {
					name
				       }
			       }
		       `,
		},
		{
			name: "eq arg",
			query: `
			{
				me(func: uid(1, 20)) @filter(eq(name, ["Andrea", "Bob"])) {
				 name
				}
			}
		`,
		},
		{
			name: "multiple equal",
			query: `{
				me(func: eq(name,["Steven Spielberg", "Tom Hanks"])) {
					name
				}
			}`,
		},
		{
			name: "has filter at child",
			query: `{
				me(func: anyofterms(name, "Steven Tom")) {
					name
					director.film @filter(has(genre)) {
					}
				}
			}`,
		},
		{
			name: "has filter at root",
			query: `{
				me(func: allofterms(name, "Steven Tom")) @filter(has(director.film)) {
					name
				}
			}`,
		},
		{
			name: "has func at root",
			query: `{
			       me(func: has(name@en)) {
				       name
			       }
		       }`,
		},
		{
			name: "count at root",
			query: `{
				me(func: uid( 1)) {
					count(uid)
					count(enemy)
				}
			}`,
		},
		{
			name: "regexp 3",
			query: `
			       {
				 me(func:allofterms(name, "barack")) @filter(regexp(secret, /whitehouse[0-9]{1,4}/fLaGs)) {
				   name
				 }
			   }
		       `,
		},
		{
			name: "regexp 2",
			query: `
			       {
				 me(func:regexp(name@en, /another\/compilicated ("") regexp('')/)) {
				   name
				 }
			   }
		       `,
		},
		{
			name: "regexp 1",
			query: `
			       {
				 me(func: uid(0x1)) {
				   name
				       friend @filter(regexp(name@en, /case .+*?()|[]{}^$" INSENSITIVE regexp with \/ escaped value/i)) {
				     name@en
				   }
				 }
			   }
		       `,
		},
		{
			name: "with attr lang 2",
			query: `
			       {
				 me(func:regexp(name, /^[a-zA-z]*[^Kk ]?[Nn]ight/), orderasc: name@en, first:5) {
				       name@en
				       name@de
				       name@it
				 }
			       }
		       `,
		},
		{
			name: "with attr lang",
			query: `
			       {
				       me(func: uid(0x1)) {
					       name
					       friend(first:5, orderasc: name@en:fr) {
						       name@en
					       }
				       }
			       }
		       `,
		},
		{
			name: "facets at value",
			query: `
			       {
				       me(func: uid(0x1)) {
					       friend	{
						       name @facets(eq(some.facet, true))
					       }
				       }
			       }
		       `,
		},
		{
			name: "facets filter all",
			query: `
			       {
				me(func: uid(0x1)) {
					name
					friend @facets(eq(close, true) or eq(family, true)) @facets(close, family, since) {
						name @facets
						gender
					}
				}
			}
		`,
		},
		{
			name: "facets filter simple",
			query: `
			       {
				me(func: uid(0x1)) {
					name
					friend @facets(eq(close, true)) {
						name
						gender
					}
				}
			}
		`,
		},
		{
			name: "facets empty",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @facets() {
					       }
					       hometown
					       age
				       }
			       }
		       `,
		},
		{
			name: "facets multiple repeat",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @facets {
						       name @facets(key1, key2, key3, key1)
					       }
					       hometown
					       age
				       }
			       }
		       `,
		},
		{
			name: "facets multiple var",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @facets {
						       name @facets(a as key1, key2, b as key3)
					       }
					       hometown
					       age
				       }
				       h(func: uid(a, b)) {
					       uid
				       }
			       }
		       `,
		},
		{
			name: "facets alias",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @facets {
						       name @facets(a1: key1, a2: key2, a3: key3)
					       }
				       }
			       }
		       `,
		},
		{
			name: "facets multiple",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @facets {
						       name @facets(key1, key2, key3)
					       }
					       hometown
					       age
				       }
			       }
		       `,
		},
		{
			name: "order by facet",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @facets {
						       name @facets(facet1)
					       }
					       hometown
					       age
				       }
			       }
		       `,
		},
		{
			name: "facets",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @facets(orderdesc: closeness) {
						       name
					       }
				       }
			       }
		       `,
		},
		{
			name: "facets order var",
			query: `
			       query {
				       me1(func: uid(0x1)) {
					       friends @facets(orderdesc: a as b) {
						       name
					       }
				       }
				       me2(func: uid(a)) {
					 foo
				       }
			       }
		       `,
		},
		{
			name: "facets order with alias",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @facets(orderdesc: closeness, b as some, order: abc, key, key1: val, abcd) {
						       val(b)
					       }
				       }
			       }
		       `,
		},
		{
			name: "group by with alias for key",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @groupby(Name: name, SchooL: school) {
						       count(uid)
					       }
					       hometown
					       age
				       }
			       }
		       `,
		},
		{
			name: "group by with alias",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @groupby(name) {
						       GroupCount: count(uid)
					       }
					       hometown
					       age
				       }
			       }
		       `,
		},
		{
			name: "group by",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @groupby(name@en) {
						       count(uid)
					       }
					       hometown
					       age
				       }
			       }
		       `,
		},
		{
			name: "group by with max var",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @groupby(friends) {
						       a as max(first-name@en:ta)
					       }
					       hometown
					       age
				       }

				       groups(func: uid(a)) {
					       uid
					       val(a)
				       }
			       }
		       `,
		},
		{
			name: "group by with count var",
			query: `
			       query {
				       me(func: uid(0x1)) {
					       friends @groupby(friends) {
						       a as count(uid)
					       }
					       hometown
					       age
				       }

				       groups(func: uid(a)) {
					       uid
					       val(a)
				       }
			       }
		       `,
		},
		{
			name: "group by root",
			query: `
			       query {
				       me(func: uid(1, 2, 3)) @groupby(friends) {
						       a as count(uid)
				       }

				       groups(func: uid(a)) {
					       uid
					       val(a)
				       }
			       }
		       `,
		},
		{
			name: "normalize",
			query: `
			       query {
				       me(func: uid( 0x3)) @normalize {
					       friends {
						       name
					       }
					       gender
					       hometown
				       }
		       }
		       `,
		},
		{
			name: "langs function",
			query: `
			       query {
				       me(func:alloftext(descr@en, "something")) {
					       friends {
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "langs filter",
			query: `
			       query {
				       me(func: uid(0x0a)) {
					       friends @filter(alloftext(descr@en, "something")) {
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "langs",
			query: `
			       query {
				       me(func: uid(1)) {
					       name@en,name@en:ru:hu
				       }
			       }
			       `,
		},
		{
			name: "IRIRef 2",
			query: `{
				me(func:anyofterms(<http://helloworld.com/how/are/you>, "good better bad")) {
					<http://verygood.com/what/about/you>
					friends @filter(allofterms(<http://verygood.com/what/about/you>,
						"good better bad")){
						name
					}
				}
			}`,
		},
		{
			name: "IRIREF",
			query: `{
			       me(func: uid( 0x1)) {
				       <http://verygood.com/what/about/you>
				       friends @filter(allofterms(<http://verygood.com/what/about/you>,
					       "good better bad")){
					       name
				       }
				       gender,age
				       hometown
			       }
		       }`,
		},
		{
			name: "generator",
			query: `{
			       me(func:allofterms(name, "barack")) {
				       friends {
					       name
				       }
				       gender,age
				       hometown
			       }
		       }
	       `,
		},
		{
			name: "check pwd",
			query: `{
				me(func: uid(1)) {
					checkpwd(password, "123456")
					hometown
				}
			}
		`,
		},
		{
			name: "parse count as func",
			query: `{
				me(func: uid(1)) {
					count(friends)
					gender,age
					hometown
				}
			}
		`,
		},
		{
			name: "count as func multiple",
			query: `{
			       me(func: uid(1)) {
				       count(friends), count(relatives)
				       count(classmates)
				       gender,age
				       hometown
			       }
		       }
	       `,
		},
		{
			name: "filter geo 2",
			query: `
			       query {
				       me(func: uid(0x0a)) {
					       friends @filter(within(loc, [[11.2 , -2.234 ], [ -31.23, 4.3214] , [5.312, 6.53]] )) {
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "filter geo 1",
			query: `
			       query {
				       me(func: uid(0x0a)) {
					       friends @filter(near(loc, [-1.12 , 2.0123 ], 100.123 )) {
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "filter brac",
			query: `
			       query {
				       me(func: uid(0x0a)) {
					       friends @filter(  a(name, "hello") or b(name, "world", "is") and (c(aa, "aaa") or (d(dd, "haha") or e(ee, "aaa"))) and f(ff, "aaa")){
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "filter op 2",
			query: `
			       query {
				       me(func: uid(0x0a)) {
					       friends @filter((a(aa, "aaa") Or b(bb, "bbb"))
						and c(cc, "ccc")) {
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "filter op not 2",
			query: `
			       query {
				       me(func: uid(0x0a)) {
					       friends @filter(not(a(aa, "aaa") or (b(bb, "bbb"))) and c(cc, "ccc")) {
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "filter op not 1",
			query: `
			       query {
				       me(func: uid(0x0a)) {
					       friends @filter(not a(aa, "aaa")) {
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "filter op",
			query: `
				query {
					me(func: uid(0x0a)) {
						friends @filter(a(aa, "aaa") or b(bb, "bbb")
						and c(cc, "ccc")) {
							name
						}
						gender,age
						hometown
					}
				}
			`,
		},
		{
			name: "filter simplest",
			query: `
			       query {
				       me(func: uid(0x0a)) {
					       friends @filter() {
						       name @filter(namefilter(name, "a"))
					       }
					       gender @filter(eq(g, "a")),age @filter(neq(a, "b"))
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "filter root 2",
			query: `
			       query {
				       me(func:anyofterms(abc, "Abc")) @filter(gt(count(friends), 10)) {
					       friends @filter() {
						       name
					       }
					       hometown
				       }
			       }
			`,
		},
		{
			name: "func nested 2",
			query: `
			       query {
				       var(func:uid(1)) {
					       a as name
				       }
				       me(func: eq(name, val(a))) {
					       friends @filter() {
						       name
					       }
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "func nested",
			query: `
			       query {
				       me(func: gt(count(friend), 10)) {
					       friends @filter() {
						       name
					       }
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "filter root",
			query: `
			query {
			       me(func:anyofterms(abc, "Abc")) @filter(allofterms(name, "alice")) {
				       friends @filter() {
					       name @filter(namefilter(name, "a"))
				       }
				       gender @filter(eq(g, "a")),age @filter(neq(a, "b"))
				       hometown
			       }
		       }
			       `,
		},
		{
			name: "string var in filter",
			query: `
			       query me($version: string)
			       {
				       versions(func:eq(type, "version"))
				       {
					       versions @filter(eq(version_number, $version))
					       {
						       version_number
					       }
				       }
			       }
		       `,
			variables: map[string]string{"$version": "string"},
			alias:     "me",
		},
		{
			name: "fragment no nest 2",
			query: `
		       query {
			       user(func: uid(0x0a)) {
				       friends {
					       ...fragmenta
				       }
			       }
		       }
		       fragment fragmenta {
			       name
			       ...fragmentb
		       }
		       fragment fragmentb {
			       nickname
		       }
	       `,
		},
		{
			name: "fragment no nest 1",
			query: `
		       query {
			       user(func: uid(0x0a)) {
				       ...fragmenta
				       friends {
					       name
				       }
			       }
		       }

		       fragment fragmenta {
			       id
			       ...fragmentb
		       }

		       fragment fragmentb {
			       hobbies
		       }
	       `,
		},
		{
			name: "fragmenrt no nesting",
			query: `
		       query {
			       user(func: uid(0x0a)) {
				       ...fragmenta,...fragmentb
				       friends {
					       name
				       }
				       ...fragmentc
				       hobbies
				       ...fragmentd
			       }
		       }

		       fragment fragmenta {
			       name
		       }

		       fragment fragmentb {
			       id
		       }

		       fragment fragmentc {
			       name
		       }

		       fragment fragmentd {
			       id
		       }
	       `,
		},
		{
			name: "fragment multi query",
			query: `
		       {
			       user(func: uid(0x0a)) {
				       ...fragmenta,...fragmentb
				       friends {
					       name
				       }
				       ...fragmentc
				       hobbies
				       ...fragmentd
			       }

			       me(func: uid(0x01)) {
				       ...fragmenta
				       ...fragmentb
			       }
		       }

		       fragment fragmenta {
			       name
		       }

		       fragment fragmentb {
			       id
		       }

		       fragment fragmentc {
			       name
		       }

		       fragment fragmentd {
			       id
		       }
	       `,
		},

		{
			name: "block",
			query: `
			       {
				       root(func: uid( 0x0a)) {
					       type.object.name.es.419
				       }
			       }
		       `,
		},
		{
			name: "alias 1",
			query: `
			       {
				       me(func: uid(0x0a)) {
					       name: type.object.name.en
					       bestFriend: friends(first: 10) {
						       name: type.object.name.hi
					       }
				       }
			       }
		       `,
		},
		{
			name: "parse alias",
			query: `
			       {
				       me(func: uid(0x0a)) {
					       name,
					       bestFriend: friends(first: 10) {
						       name
					       }
				       }
			       }
		       `,
		},
		{
			name: "alias max",
			query: `
			       {
				       me(func: uid(0x0a)) {
					       name,
					       bestFriend: friends(first: 10) {
						       x as count(friends)
					       }
					       maxfriendcount: max(val(x))
				       }
			       }
		       `,
		},
		{
			name: "alias var",
			query: `
			       {
				       me(func: uid(0x0a)) {
					       name,
					       f as bestFriend: friends(first: 10) {
						       c as count(friend)
					       }
				       }

				       friend(func: uid(f)) {
					       name
					       fcount: val(c)
				       }
			       }
		       `,
		},
		{
			name: "alias count",
			query: `
			       {
				       me(func: uid(0x0a)) {
					       name,
					       bestFriend: friends(first: 10) {
						       nameCount: count(name)
					       }
				       }
			       }
		       `,
		},
		{
			name: "first",
			query: `
			       query {
				       user(func: uid( 0x1)) {
					       type.object.name
					       friends (first: 10) {
					       }
				       }
			       }`,
		},
		{
			name: "id list 1",
			query: `
			       query {
				       user(func: uid(0x1, 0x34)) {
					       type.object.name
				       }
			       }`,
		},
		{
			name: "root args 2",
			query: `
			       query {
				       me(func: uid(0x0a), first: 1, offset:0) {
					       friends {
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "root args 1",
			query: `
			       query {
				       me(func: uid(0x0a), first: -4, offset: +1) {
					       friends {
						       name
					       }
					       gender,age
					       hometown
				       }
			       }
		       `,
		},
		{
			name: "multiple queries",
			query: `
			       {
				       you(func: uid(0x0a)) {
					       name
				       }

				       me(func: uid(0x0b)) {
					friends
				       }
			       }
		       `,
		},
		{
			name: "shortest path",
			query: `
			       {
				       shortest(from:0x0a, to:0x0b, numpaths: 3) {
					       friends
					       name
				       }
			       }
		       `,
		},
		{
			name: "with multiple var",
			query: `
			       {
				       var(func: uid(0x0a)) {
					       L AS friends {
						       B AS relatives
					       }
				       }

				       me(func: uid(L)) {
					name
				       }

				       relatives(func: uid(B)) {
					       name
				       }
			       }
		       `,
		},
		{
			name: "with var in ineq",
			query: `
			       {
				       var(func: uid(0x0a)) {
					       fr as friends {
						       a as age
					       }
				       }

				       me(func: uid(fr)) @filter(gt(val(a), 10)) {
					name
				       }
			       }
		       `,
		},
		{
			name: "with var at root",
			query: `
			       {
				       K AS var(func: uid(0x0a)) {
					       fr as friends
				       }
				       me(func: uid(fr)) @filter(uid(K)) {
					name	@filter(uid(fr))
				       }
			       }
		       `,
		},
		{
			name: "with var at root filter id",
			query: `
			       {
				       K as var(func: uid(0x0a)) {
					       L AS friends
				       }
				       me(func: uid(K)) @filter(uid(L)) {
					name
				       }
			       }
		       `,
		},
		{
			name: "with var",
			query: `
			       {
				       me(func: uid( L, J, K)) {name}
				       var(func: uid(0x0a)) {L AS friends}
				       var(func: uid(0x0a)) {J AS friends}
				       var(func: uid(0x0a)) {K AS friends}
			       }
		       `,
		},
		{
			name: "with var multi root",
			query: `
			       {
				       me(func: uid( L, J, K)) {name}
				       var(func: uid(0x0a)) {L AS friends}
				       var(func: uid(0x0a)) {J AS friends}
				       var(func: uid(0x0a)) {K AS friends}
			       }
		       `,
		},
		{
			name: "with var val",
			query: `
			       {
				       me(func: uid(L), orderasc: val(n) ) {
					       name
				       }

				       var(func: uid(0x0a)) {
					       L AS friends {
						       n AS name
					       }
				       }
			       }
		       `,
		},
		{
			name: "with var val count",
			query: `
			       {
				       me(func: uid(L), orderasc: val(n) ) {
					       name
				       }

				       var(func: uid(0x0a)) {
					       L AS friends {
						       na as name
					       }
					       n as min(val(na))
				       }
			       }
		       `,
		},
		{
			name: "with var val agg",
			query: `
			       {
				       me(func: uid(L), orderasc: val(n) ) {
					       name
				       }

				       var(func: uid(0x0a)) {
					       L AS friends {
						       na as name
					       }
					       n as min(val(na))
				       }
			       }
		       `,
		},
		{
			name: "with var val agg combination",
			query: `
			       {
				       me(func: uid(L), orderasc: val(c) ) {
					       name
					       val(c)
				       }

				       var(func: uid(0x0a)) {
					       L as friends {
						       x as age
					       }
					       a as min(val(x))
					       b as max(val(x))
					       c as math(a + b)
				       }
			       }
		       `,
		},
		{
			name: "with level agg",
			query: `
			       {
				       var(func: uid(0x0a)) {
					       friends {
						       a as count(age)
					       }
					       s as sum(val(a))
				       }

				       sumage(func: uid( 0x0a)) {
					       val(s)
				       }
			       }
		       `,
		},
		{
			name: "with var val agg nested 3",
			query: `
			       {
				       me(func: uid(L), orderasc: val(d) ) {
					       name
			       }

				       var(func: uid(0x0a)) {
					       L as friends {
						       a as age
						       b as count(friends)
						       c as count(relatives)
						       d as math(a + b * c / a + exp(a + b + 1) - ln(c))
					       }
				       }
			       }
		       `,
		},
		{
			name: "with var val agg nested conditional",
			query: `
			       {
				       me(func: uid(L), orderasc: val(d) ) {
					       name
					       val(f)
				       }

				       var(func: uid(0x0a)) {
					       L as friends {
						       a as age
						       b as count(friends)
						       c as count(relatives)
						       d as math(cond(a <= 10, exp(a + b + 1), ln(c)) + 10*a)
						       e as math(cond(a!=10, exp(a + b + 1), ln(d)))
						       f as math(cond(a==10, exp(a + b + 1), ln(e)))
					       }
				       }
			       }
		       `,
		},
		{
			name: "with var val agg log sqrt",
			query: `
			       {
				       me(func: uid(L), orderasc: val(d) ) {
					       name
					       val(e)
				       }

				       var(func: uid(0x0a)) {
					       L as friends {
						       a as age
						       d as math(ln(sqrt(a)))
						       e as math(sqrt(ln(a)))
					       }
				       }
			       }
		       `,
		},
		{
			name: "with var val agg nested 4",
			query: `
			       {
				       me(func: uid(L), orderasc: val(d) ) {
					       name
				       }

				       var(func: uid(0x0a)) {
					       L as friends {
						       a as age
						       b as count(friends)
						       c as count(relatives)
						       d as math(exp(a + b + 1) - max(c,ln(c)) + sqrt(a%b))
					       }
				       }
			       }
		       `,
		},
		{
			name: "with var val agg nested2",
			query: `
			       {
				       me(func: uid(L), orderasc: val(d)) {
					       name
					       val(q)
				       }

				       var(func: uid(0x0a)) {
					       L as friends {
						       a as age
						       b as count(friends)
						       c as count(relatives)
						       d as math(exp(a + b + 1) - ln(c))
						       q as math(c*-1+-b+(-b*c))
					       }
				       }
			       }
		       `,
		},
		{
			name: "with var val agg nested",
			query: `
			       {
				       me(func: uid(L), orderasc: val(d)) {
					       name
				       }

				       var(func: uid(0x0a)) {
					       L as friends {
						       a as age
						       b as count(friends)
						       c as count(relatives)
						       d as math(a + b*c)
					       }
				       }
			       }
		       `,
		},
		{
			name: "with no var val error",
			query: `
			       {
				       me(func: uid(), orderasc: val(n)) {
					       name
				       }

				       var(func: uid(0x0a)) {
					       friends {
						       n AS name
					       }
				       }
			       }
		       `,
		},
		{
			name: "list pred 2",
			query: `
			       {
				       var(func: uid(0x0a)) {
					       f as friends
				       }

				       var(func: uid(f)) {
					       l as _predicate_
				       }

				       var(func: uid( 0x0a)) {
					       friends {
						       expand(val(l))
					       }
				       }
			       }
		       `,
		},
		{
			name: "count list pred",
			query: `
			       {
				       me(func: uid(0x0a)) {
					       count(_predicate_)
				       }
			       }
		       `,
		},
		{
			name: "query alias list pred",
			query: `
			       {
				       me(func: uid(0x0a)) {
					       pred: _predicate_
				       }
			       }
		       `,
		},
		{
			name: "query expand reverse",
			query: `
			       {
				       var(func: uid( 0x0a)) {
					       friends {
						       expand(_reverse_)
					       }
				       }
			       }
		       `,
		},
		{
			name: "query list pred",
			query: `
				{
					var(func: uid( 0x0a)) {
						friends {
							expand(_all_)
						}
					}
				}
			`,
		},
	}

	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			gqlVariables := make(map[string]string, len(e.variables))
			for k := range e.variables {
				gqlVariables[k] = k
			}

			gqlExpected, err := dgraphgql.Parse(dgraphgql.Request{Str: e.query, Variables: gqlVariables})
			if err != nil {
				t.Fatalf("parse: %v", err)
			}

			expected := gql.DecodeGraphQueries(gqlExpected.Query)
			renderedQuery, err := Render(Query{Queries: expected, Alias: e.alias, Variables: e.variables})
			if err != nil {
				t.Fatalf("render: %v", err)
			}

			gqlActual, err := dgraphgql.Parse(dgraphgql.Request{Str: renderedQuery, Variables: gqlVariables})
			if err != nil {
				t.Fatalf("parse: %v", err)
			}

			actual := gql.DecodeGraphQueries(gqlActual.Query)
			require.Equal(t, expected, actual)
		})
	}
}
