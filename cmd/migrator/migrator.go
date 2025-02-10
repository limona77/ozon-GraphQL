package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	var migrationsPath, migrationsTable, action string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.StringVar(&action, "action", "up", "migration action: up or down")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("Error loading .env file in migratgor"))
	}

	postgresUser := os.Getenv("DB_USER")
	postgresPassword := os.Getenv("DB_PASSWORD")
	postgresHost := "localhost"
	postgresPort := "5434"
	postgresDB := os.Getenv("DB_NAME")
	fmt.Println(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB))
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB),
	)
	if err != nil {
		panic(err)
	}

	switch action {
	case "up":
		// Применяем миграции
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")
				return
			}
			panic(err)
		}
		fmt.Println("migrations applied")
	case "down":
		// Откатываем последнюю миграцию
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to revert")
				return
			}
			panic(err)
		}
		fmt.Println("last migration reverted")
	default:
		fmt.Println("Invalid action. Use 'up' or 'down'.")
	}
	fmt.Println("migrations applied")
}
