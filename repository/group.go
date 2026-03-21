package repository

import (
	"database/sql"
	"errors"
	"github.com/ChaitanyaSai-Meka/devledger/models"
	"strings"
)

func CreateGroup(db *sql.DB, groupName string) error {
	_, err := db.Exec("INSERT INTO Groups (GroupName) VALUES (?)", groupName)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("group name already exists")
		}
		return err
	}
	return nil
}

func GetAllGroups(db *sql.DB) ([]models.Group, error) {
	rows, err := db.Query("SELECT GroupID,GroupName FROM Groups")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.GroupID, &group.GroupName); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return groups, nil
}

func GetGroupByID(db *sql.DB, groupID int) (models.Group, error) {
	var group models.Group
	err := db.QueryRow(
		"SELECT GroupID, GroupName FROM Groups WHERE GroupID = ?",
		groupID,
	).Scan(&group.GroupID, &group.GroupName)

	if err != nil {
		return models.Group{}, err
	}
	return group, nil
}

func GetGroupByName(db *sql.DB, groupName string) (models.Group, error) {
	var group models.Group
	err := db.QueryRow(
		"SELECT GroupID, GroupName FROM Groups WHERE GroupName = ?",
		groupName,
	).Scan(&group.GroupID, &group.GroupName)

	if err != nil {
		return models.Group{}, err
	}
	return group, nil
}

func DeleteGroupByID(db *sql.DB, groupID int) error {
	result, err := db.Exec("DELETE FROM Groups WHERE GroupID = ?", groupID)
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

func AddMember(db *sql.DB, groupID int, userID int) error {
	_, err := db.Exec("INSERT INTO GroupMembers (GroupID, UserID) VALUES (?,?)", groupID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("user is already a member of this group")
		}
		return err
	}
	return nil
}

func RemoveMember(db *sql.DB, groupID int, userID int) error {
	result, err := db.Exec("DELETE FROM GroupMembers WHERE GroupID = ? AND UserID = ?", groupID, userID)
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

func GetGroupMembers(db *sql.DB, groupID int) ([]models.User, error) {
	rows, err := db.Query(" SELECT u.UserID, u.UserName FROM Users u JOIN GroupMembers gm ON u.UserID = gm.UserID WHERE gm.GroupID = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.UserName); err != nil {
			return nil, err
		}
		members = append(members, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func GetGroupsByUserID(db *sql.DB, userID int) ([]models.Group, error) {
	rows, err := db.Query("SELECT g.GroupID, g.GroupName FROM Groups g JOIN GroupMembers gm ON g.GroupID = gm.GroupID WHERE gm.UserID = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.GroupID, &group.GroupName); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return groups, nil
}
