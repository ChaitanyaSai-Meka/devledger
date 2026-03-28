package repository

import (
	"database/sql"
	"github.com/ChaitanyaSai-Meka/devledger/models"
	"strings"
)

func CreateUser(db *sql.DB, username string) error {
	_, err := db.Exec("INSERT INTO Users (UserName) VALUES (?)", username)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func GetUserByID(db *sql.DB, userID int) (models.User, error) {
	var user models.User

	err := db.QueryRow(
		"SELECT UserID, UserName FROM Users WHERE UserID = ? AND DeletedAt IS NULL",
		userID,
	).Scan(&user.UserID, &user.UserName)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserByName(db DBTX, username string) (models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT UserID, UserName FROM Users WHERE UserName = ? AND DeletedAt IS NULL", username).Scan(&user.UserID, &user.UserName)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetAllUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query("SELECT UserID, UserName FROM Users WHERE DeletedAt IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.UserName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUserByID(db *sql.DB, userID int) error {
	result, err := db.Exec("UPDATE Users SET DeletedAt = CURRENT_TIMESTAMP WHERE UserID = ?",
		userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func GetUserByIDIncludingDeleted(db *sql.DB, userID int) (models.User, error) {
	var user models.User

	err := db.QueryRow(
		"SELECT UserID, UserName FROM Users WHERE UserID = ? ",
		userID,
	).Scan(&user.UserID, &user.UserName)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
