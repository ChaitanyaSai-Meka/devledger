package service

import (
	"database/sql"
	"errors"
	"strings"
	"fmt"

	"github.com/ChaitanyaSai-Meka/devledger/models"
	"github.com/ChaitanyaSai-Meka/devledger/repository"
)

func CreateUser(db *sql.DB, username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username cannot be empty")
	}
	err := repository.CreateUser(db, username)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *sql.DB, username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username cannot be empty")
	}
	user, err := repository.GetUserByName(db, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user '%s' not found", username)
		}
		return err
	}
	err = repository.DeleteUserByID(db, user.UserID)
	if err != nil {
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
