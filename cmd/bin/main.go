package main

import (
	"log"
	"net/http"
	"os"
	"ozon-test/internal/gql"
	"ozon-test/internal/inmemory"
	"ozon-test/internal/models"
	"ozon-test/internal/postgres"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	storageType := os.Getenv("STORAGE_TYPE")
	var storage models.Storage

	if storageType == "postgres" {
		db, err := sqlx.Connect("postgres", "host="+os.Getenv("DB_HOST")+" port="+os.Getenv("DB_PORT")+" user="+os.Getenv("DB_USER")+" password="+os.Getenv("DB_PASSWORD")+" dbname="+os.Getenv("DB_NAME")+" sslmode=disable")
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		storage = postgres.NewPostgresStorage(db)
	} else {
		storage = inmemory.NewInMemoryStorage()
	}

	srv := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{Storage: storage}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
