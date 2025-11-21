package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"

	"wdi/graph"
	"wdi/internal/database"
	"wdi/internal/repository"
	"wdi/internal/usecase"
)

func main() {
	_ = godotenv.Load()

	db := database.New()

	CountryRepository := repository.NewCountryRepository(db)
	CountryUsecase := usecase.NewCountryUsecase(CountryRepository)

	resolver := &graph.Resolver{
		CountryUC: CountryUsecase,
	}

	//! Serveur GraphQL
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{Resolvers: resolver}),
	)

	//! Routes HTTP
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server ready at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
