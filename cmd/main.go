package main

import (
	"context"
	"fmt"
	"github.com/aakosarev/transaction-system/internal/config"
	"github.com/aakosarev/transaction-system/internal/user"
	"github.com/aakosarev/transaction-system/pkg/client/postgresql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.GetConfig()

	router := httprouter.New()

	pgConfig := postgresql.NewPgConfig(
		cfg.PostgreSQL.Username, cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(ctx, 5, time.Second*5, pgConfig)
	if err != nil {
		log.Fatal(err)
	}

	userStorage := user.NewStorage(pgClient)
	userService := user.NewService(userStorage)
	userHandler := user.NewHandler(userService)
	worker := user.NewWorker(userService)

	userHandler.Register(router)

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", pgConfig.Username, pgConfig.Password, pgConfig.Host, pgConfig.Port, pgConfig.Database)
	runDBMigrations("file://migrations", dsn)

	worker.StartTransactionProcessing(ctx)
	start(router, cfg)
}

func start(router http.Handler, cfg *config.Config) {
	var server *http.Server

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.HTTP.IP, cfg.HTTP.Port))
	if err != nil {
		log.Fatal(err)
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func runDBMigrations(migrationsURL string, databaseURL string) {
	migrations, err := migrate.New(migrationsURL, databaseURL)
	if err != nil {
		log.Fatal("cannot create new migrate instance: ", err)
	}

	if err = migrations.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}

	log.Println("db migrated successfully")
}
