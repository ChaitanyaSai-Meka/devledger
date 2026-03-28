package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ChaitanyaSai-Meka/devledger/api/respond"
	"github.com/ChaitanyaSai-Meka/devledger/service"
	"github.com/go-chi/chi/v5"
)

func AddExpenseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupname := chi.URLParam(r, "groupname")
		if groupname == "" {
			respond.WriteError(w, http.StatusBadRequest, "group name is required")
			return
		}
		var req struct {
			Description string `json:"description"`
			Amount      string `json:"amount"`
			PaidBy      string `json:"paid_by"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respond.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		err := service.AddExpense(db, groupname, req.Description, req.PaidBy, req.Amount)
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
		respond.WriteCreated(w, map[string]any{"message": "Expense added successfully"})

	}
}

func ListExpensesByGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupname := chi.URLParam(r, "groupname")
		if groupname == "" {
			respond.WriteError(w, http.StatusBadRequest, "group name is required")
			return
		}
		expenses, err := service.ListExpensesByGroup(db, groupname)
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
		respond.WriteOK(w, expenses)
	}
}

func ListExpensesByUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		if username == "" {
			respond.WriteError(w, http.StatusBadRequest, "username is required")
			return
		}
		expenses, err := service.ListExpensesByUser(db, username)
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
		respond.WriteOK(w, expenses)
	}
}

func DeleteExpenseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expenseIDStr := chi.URLParam(r, "expenseID")
		if expenseIDStr == "" {
			respond.WriteError(w, http.StatusBadRequest, "expense ID is required")
			return
		}
		expenseID, err := strconv.ParseInt(expenseIDStr, 10, 64)
		if err != nil {
			respond.WriteError(w, http.StatusBadRequest, "invalid expense ID")
			return
		}
		err = service.DeleteExpense(db, expenseID)
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
		respond.WriteOK(w, map[string]any{"message": "Expense deleted successfully"})
	}
}

func SettleExpenseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expenseIDStr := chi.URLParam(r, "expenseID")
		if expenseIDStr == "" {
			respond.WriteError(w, http.StatusBadRequest, "expense ID is required")
			return
		}
		expenseID, err := strconv.ParseInt(expenseIDStr, 10, 64)
		if err != nil {
			respond.WriteError(w, http.StatusBadRequest, "invalid expense ID")
			return
		}
		username := chi.URLParam(r, "username")
		if username == "" {
			respond.WriteError(w, http.StatusBadRequest, "username is required")
			return
		}
		err = service.SettleExpense(db, expenseID, username)
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
		respond.WriteOK(w, map[string]any{"message": "Expense settled successfully"})
	}
}

func ListUnsettledSplitsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		if username == "" {
			respond.WriteError(w, http.StatusBadRequest, "username is required")
			return
		}
		splits, err := service.ListUnsettledSplits(db, username)
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
		respond.WriteOK(w, splits)
	}
}

func ExpenseInDetailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expenseIDStr := chi.URLParam(r, "expenseID")
		if expenseIDStr == "" {
			respond.WriteError(w, http.StatusBadRequest, "expense ID is required")
			return
		}
		expenseID, err := strconv.ParseInt(expenseIDStr, 10, 64)
		if err != nil {
			respond.WriteError(w, http.StatusBadRequest, "invalid expense ID")
			return
		}
		expense, err := service.GetExpenseInDetail(db, expenseID)
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
		respond.WriteOK(w, expense)
	}
}
