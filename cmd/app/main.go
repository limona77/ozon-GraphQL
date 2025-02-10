package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"ozon-GraphQL/graph"
	graph2 "ozon-GraphQL/graph"
	"ozon-GraphQL/internal/database"
	"ozon-GraphQL/internal/database/storage"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	storageType := os.Getenv("STORAGE_TYPE")

	var repo database.Repository

	if storageType == "" {
		log.Fatalf("Error: STORAGE_TYPE environment variable not set")
	}

	if storageType == "postgres" {
		err = waitForDatabase(dbUser, dbPassword, dbName, dbHost, dbPort)
		if err != nil {
			log.Fatalf("Error waiting for database: %v", err)
		}

		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			dbUser, dbPassword, dbName, dbHost, dbPort)

		db, err := pgx.Connect(context.Background(), connStr)
		if err != nil {
			log.Fatalf("failed to connect to the database: %v", err)
		}
		defer db.Close(context.Background())

		repo = storage.NewPostgresSQLRepository(db)

	} else if storageType == "in_memory" {
		repo = storage.NewInMemoryRepository()
	}
	resolver := graph.NewResolver(repo)
	srv := handler.New(graph2.NewExecutableSchema(graph2.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func waitForDatabase(user, password, dbName, host, port string) error {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbName, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for database to become available: %v", ctx.Err())
		default:
			db, err := pgx.Connect(ctx, connStr)
			if err == nil {
				db.Close(ctx)
				return nil
			}
			log.Printf("waiting for database to become available on %s:%s - %v", host, port, err)
			time.Sleep(1 * time.Second)
		}
	}
}
