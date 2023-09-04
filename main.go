package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andriiklymiuk/go_server_user_data/v2/src/api"
	"github.com/andriiklymiuk/go_server_user_data/v2/src/db/migrations"
	"github.com/andriiklymiuk/go_server_user_data/v2/src/utils"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	envConfig, err := utils.LoadConnectionConfig()
	if err != nil {
		color.Red("Couldn't load env variables: \n%v", err)
		os.Exit(1)
	}

	if migrations.CheckMigrationFlags(envConfig) {
		os.Exit(0)
	}

	migrations.MigrateUpDb(envConfig)

	dbpool, err := pgxpool.New(context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			envConfig.DbUser,
			envConfig.DbPassword,
			envConfig.DbHost,
			envConfig.DbPort,
			envConfig.DbName,
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	router := api.SetupRouter(dbpool)
	srv := setupServer(envConfig.ServerPort, router)

	defer dbpool.Close()

	testConnection(dbpool)

	// Shutdown signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			color.Red("Listen: %s\n", err)
			os.Exit(1)
		}
	}()

	color.Green("Server Started")

	sig := <-done
	color.Yellow("\nReceived signal: %s. Starting shutdown", sig)

	handleShutdownSignal(sig, srv)
}

func setupServer(port int, router *chi.Mux) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
}

func handleShutdownSignal(sig os.Signal, srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		color.Red("Server Shutdown Failed: %+v", err)
	} else {
		color.Green("Server Exited Properly")
	}
}

func testConnection(connection *pgxpool.Pool) {
	var greeting string
	err := connection.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)
}
