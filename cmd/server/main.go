package main

import (
	"log"
	"net/http"

	"wdi/internal/infrastructure/db"
	"wdi/internal/interface/graph"
	"wdi/internal/interface/graph/resolvers"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	//! Config and connect to the database
	cfg := db.LoadConfigFromEnv()
	db, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//! Set up the GraphQL server
	resolvers := &resolvers.Resolver{
		DB: db,
	}
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: resolvers,
			},
		),
	)
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Println("GraphQL server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
