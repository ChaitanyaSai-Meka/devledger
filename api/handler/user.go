package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ChaitanyaSai-Meka/devledger/api/respond"
	"github.com/ChaitanyaSai-Meka/devledger/service"
	"github.com/go-chi/chi/v5"
)

func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username string `json:"username"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			respond.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		err := service.CreateUser(db, input.Username)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidInput):
				respond.WriteError(w, http.StatusBadRequest, err.Error())
			case errors.Is(err, service.ErrNotFound):
				respond.WriteError(w, http.StatusNotFound, err.Error())
			case errors.Is(err, service.ErrConflict):
				respond.WriteError(w, http.StatusConflict, err.Error())
			default:
				respond.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		respond.WriteCreated(w, map[string]any{"message": "User created successfully"})
	}
}

func ListUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := service.GetAllUsers(db)
		if err != nil {
			respond.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		respond.WriteOK(w, users)
	}
}

func GetUserGroupsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		if username == "" {
			respond.WriteError(w, http.StatusBadRequest, "username is required")
			return
		}

		groups, err := service.GetUserGroups(db, username)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidInput):
				respond.WriteError(w, http.StatusBadRequest, err.Error())
			case errors.Is(err, service.ErrNotFound):
				respond.WriteError(w, http.StatusNotFound, err.Error())
			case errors.Is(err, service.ErrConflict):
				respond.WriteError(w, http.StatusConflict, err.Error())
			default:
				respond.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		respond.WriteOK(w, groups)
	}
}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		if username == "" {
			respond.WriteError(w, http.StatusBadRequest, "username is required")
			return
		}

		err := service.DeleteUser(db, username)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidInput):
				respond.WriteError(w, http.StatusBadRequest, err.Error())
			case errors.Is(err, service.ErrNotFound):
				respond.WriteError(w, http.StatusNotFound, err.Error())
			case errors.Is(err, service.ErrConflict):
				respond.WriteError(w, http.StatusConflict, err.Error())
			default:
				respond.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}
		respond.WriteOK(w, map[string]any{"message": "User deleted successfully"})
	}
}
