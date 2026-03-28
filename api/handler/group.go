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

func CreateGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var input struct {
			Groupname string `json:"groupname"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			respond.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		err := service.CreateGroup(db,input.Groupname)
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
		respond.WriteCreated(w, map[string]string{"message": "group created successfully"})
	}
}

func DeleteGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupname := chi.URLParam(r, "groupname")
		if groupname == "" {
			respond.WriteError(w, http.StatusBadRequest, "groupname is required")
			return
		}
		err := service.DeleteGroup(db, groupname)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidInput):
				respond.WriteError(w, http.StatusBadRequest, err.Error())
			case errors.Is(err, service.ErrNotFound):
				respond.WriteError(w, http.StatusNotFound, err.Error())
			default:
				respond.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}
		respond.WriteOK(w, map[string]string{"message": "group deleted successfully"})
	}
}

func AddMemberToGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupname := chi.URLParam(r, "groupname")
		if groupname == "" {
			respond.WriteError(w, http.StatusBadRequest, "groupname is required")
			return
		}
		var input struct {
			Username string `json:"username"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			respond.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		err := service.AddMemberToGroup(db, groupname, input.Username)
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
		respond.WriteCreated(w, map[string]string{"message": "member added to group successfully"})
	}
}

func RemoveMemberFromGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupname := chi.URLParam(r, "groupname")
		username := chi.URLParam(r, "username")
		if groupname == "" {
			respond.WriteError(w, http.StatusBadRequest, "groupname is required")
			return
		}
		if username == "" {
			respond.WriteError(w, http.StatusBadRequest, "username is required")
			return
		}

		err := service.RemoveMemberFromGroup(db, groupname, username)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidInput):
				respond.WriteError(w, http.StatusBadRequest, err.Error())
			case errors.Is(err, service.ErrNotFound):
				respond.WriteError(w, http.StatusNotFound, err.Error())
			default:
				respond.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}
		respond.WriteOK(w, map[string]string{"message": "member removed from group successfully"})
	}
}

func ListGroupMembersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupname := chi.URLParam(r, "groupname")
		if groupname == "" {
			respond.WriteError(w, http.StatusBadRequest, "groupname is required")
			return
		}
		members, err := service.ListGroupMembers(db, groupname)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidInput):
				respond.WriteError(w, http.StatusBadRequest, err.Error())
			case errors.Is(err, service.ErrNotFound):
				respond.WriteError(w, http.StatusNotFound, err.Error())
			default:
				respond.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}
		respond.WriteOK(w, members)
	}
}

func ListGroupsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := service.ListGroups(db)
		if err != nil {
			respond.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		respond.WriteOK(w, groups)
	}
}
