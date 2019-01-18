// Copyright © 2019 mooncamp.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"net/http"

	"mooncamp.com/dgraphtools/qb"
	"mooncamp.com/dgraphtools/qb/endpoint"
	"mooncamp.com/dgraphtools/qb/transport"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var querybuilderCmd = &cobra.Command{
	Use:   "querybuilder",
	Short: "allows translating between the data representation of a GraphQL+- query",
	Long: `QueryBuilder provides a simple interface that helps creating the data
representation of a graphql+- query. GraphQL+- queries can be inserted
and the respective data representation will be generated. In the same
manner, a data representation can be inspected by rendering it back to
the actual GraphQL+- query.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiRoute := "/api/v1"
		apiHandler := transport.NewHTTPHandler(endpoint.NewEndpointSet(), apiRoute)

		handler := mux.NewRouter()
		handler.PathPrefix(apiRoute).Handler(apiHandler)
		handler.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(qb.MustAsset("static/index.html")))
		})

		fmt.Printf("access the querybuilder through http://localhost:%s\n", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
			log.Fatalf("listen: %v", err)
		}
	},
}

var port string

func init() {
	rootCmd.AddCommand(querybuilderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// querybuilderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// querybuilderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	querybuilderCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to listen on")
}
