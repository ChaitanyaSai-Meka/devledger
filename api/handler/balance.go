package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ChaitanyaSai-Meka/devledger/api/respond"
	"github.com/ChaitanyaSai-Meka/devledger/service"
	"github.com/go-chi/chi/v5"
)

func CalculateBalanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupname := chi.URLParam(r, "groupname")
		if groupname == "" {
			respond.WriteError(w, http.StatusBadRequest, "groupname is required")
			return
		}

		balance, err := service.CalculateBalances(db, groupname)
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
		respond.WriteOK(w, balance)
	}
}

func SimplifyDebtHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupname := chi.URLParam(r, "groupname")
		if groupname == "" {
			respond.WriteError(w, http.StatusBadRequest, "groupname is required")
			return
		}
		balance, err := service.CalculateBalances(db, groupname)
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
		simplifiedDebts := service.SimplifyDebts(balance)
		respond.WriteOK(w, simplifiedDebts)
	}
}
