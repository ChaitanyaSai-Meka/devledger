package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ChaitanyaSai-Meka/devledger/models"
	"github.com/ChaitanyaSai-Meka/devledger/repository"
)

func isDuplicateUserError(err error) bool {
	return err != nil && err.Error() == "username already exists"
}

func CreateUser(db *sql.DB, username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("%w: username cannot be empty", ErrInvalidInput)
	}
	err := repository.CreateUser(db, username)
	if err != nil {
		if isDuplicateUserError(err) {
			return fmt.Errorf("%w: user '%s' already exists", ErrConflict, username)
		}
		return err
	}
	return nil
}

func DeleteUser(db *sql.DB, username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("%w: username cannot be empty", ErrInvalidInput)
	}
	user, err := repository.GetUserByName(db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user '%s' not found", ErrNotFound, username)
		}
		return err
	}

	splits, err := repository.GetUnsettledSplitsByUserID(db, user.UserID)
	if err != nil {
		return err
	}
	if len(splits) > 0 {
		return fmt.Errorf("%w: cannot delete user '%s' with unsettled debts", ErrConflict, username)
	}

	creditorSplits, err := repository.GetUnsettledSplitsForExpensesPaidByUserID(db, user.UserID)
	if err != nil {
		return err
	}
	if len(creditorSplits) > 0 {
		return fmt.Errorf("%w: cannot delete user '%s' because other users still owe them money", ErrConflict, username)
	}

	err = repository.DeleteUserByID(db, user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user '%s' not found", ErrNotFound, username)
		}
		return err
	}
	return nil
}

func GetAllUsers(db *sql.DB) ([]models.User, error) {
	users, err := repository.GetAllUsers(db)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserGroups(db *sql.DB, username string) ([]models.Group, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("%w: username cannot be empty", ErrInvalidInput)
	}
	user, err := repository.GetUserByName(db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: user '%s' not found", ErrNotFound, username)
		}
		return nil, err
	}
	groups, err := repository.GetGroupsByUserID(db, user.UserID)
	if err != nil {
		return nil, err
	}
	return groups, nil
}
