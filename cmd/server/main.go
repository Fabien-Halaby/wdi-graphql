package main

import (
	"log"
	"net/http"

	"wdi/internal/interface/graph"
	"wdi/internal/interface/graph/resolvers"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &resolvers.Resolver{},
			},
		),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Println("GraphQL server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
