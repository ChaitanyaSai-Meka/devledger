package api

import (
	"database/sql"
	"net/http"

	"github.com/ChaitanyaSai-Meka/devledger/api/handler"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Post("/users", handler.CreateUserHandler(db))
	r.Get("/users", handler.ListUsersHandler(db))
	r.Get("/users/{username}/groups", handler.GetUserGroupsHandler(db))
	r.Delete("/users/{username}", handler.DeleteUserHandler(db))

	return r
}
