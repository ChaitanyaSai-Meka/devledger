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

	r.Post("/groups", handler.CreateGroupsHandler(db))
	r.Get("/groups", handler.ListGroupsHandler(db))
	r.Get("/groups/{groupname}/members", handler.ListGroupMembersHandler(db))
	r.Delete("/groups/{groupname}/members/{username}", handler.RemoveMemberFromGroupHandler(db))
	r.Delete("/groups/{groupname}",handler.DeleteGroupHandler(db))
	r.Post("/groups/{groupname}/members", handler.AddMemberToGroupHandler(db))

	return r
}
