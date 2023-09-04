package api

import (
	"context"
	"encoding/json"
	"net/http"

	db "github.com/andriiklymiuk/go_server_user_data/v2/src/db/sqlc"
	"github.com/andriiklymiuk/go_server_user_data/v2/src/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(dbPool *pgxpool.Pool) *chi.Mux {
	router := chi.NewRouter()
	router.Use(utils.VersionMiddleware)
	router.Use(middleware.Logger)
	router.Get("/status", successHandler)

	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {

		productList, err := db.New(dbPool).GetUsers(context.Background())
		if err != nil {
			utils.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.NewEncoder(w).Encode(productList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	})

	return router
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Successful"))
}
