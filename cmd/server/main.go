package main

import (
	"log"
	"net/http"

	"wdi/internal/infrastructure/db"
	"wdi/internal/infrastructure/repository"
	"wdi/internal/interface/graph"
	"wdi/internal/interface/graph/resolvers"
	"wdi/internal/usecase"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

func main() {
	//! Config and connect to the database
	cfg := db.LoadConfigFromEnv()
	db, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//!
	indicatorRepo := repository.NewIndicatorRepository(db)
	indicatorUC := usecase.NewIndicatorUsecase(indicatorRepo)

	//! Set up the GraphQL server
	resolvers := &resolvers.Resolver{
		IndicatorUC: indicatorUC,
	}
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: resolvers,
			},
		),
	)

	//! Set up HTTP handlers
	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	//! Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	//! Start the server
	log.Println("GraphQL server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
