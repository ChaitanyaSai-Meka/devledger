package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ChaitanyaSai-Meka/devledger/models"
	"github.com/ChaitanyaSai-Meka/devledger/repository"
)

func isDuplicateGroupError(err error) bool {
	return err != nil && err.Error() == "group name already exists"
}

func isDuplicateGroupMemberError(err error) bool {
	return err != nil && err.Error() == "user is already a member of this group"
}

func CreateGroup(db *sql.DB, groupname string) error {
	groupname = strings.TrimSpace(groupname)
	if groupname == "" {
		return fmt.Errorf("%w: group name cannot be empty", ErrInvalidInput)
	}
	err := repository.CreateGroup(db, groupname)
	if err != nil {
		if isDuplicateGroupError(err) {
			return fmt.Errorf("%w: group '%s' already exists", ErrConflict, groupname)
		}
		return err
	}
	return nil
}

func DeleteGroup(db *sql.DB, groupname string) error {
	groupname = strings.TrimSpace(groupname)
	if groupname == "" {
		return fmt.Errorf("%w: group name cannot be empty", ErrInvalidInput)
	}
	group, err := repository.GetGroupByName(db, groupname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: group '%s' not found", ErrNotFound, groupname)
		}
		return err
	}
	err = repository.DeleteGroupByID(db, group.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: group '%s' not found", ErrNotFound, groupname)
		}
		return err
	}
	return nil
}

func AddMemberToGroup(db *sql.DB, groupname string, username string) error {
	groupname = strings.TrimSpace(groupname)
	username = strings.TrimSpace(username)

	if groupname == "" {
		return fmt.Errorf("%w: group name cannot be empty", ErrInvalidInput)
	}
	if username == "" {
		return fmt.Errorf("%w: username cannot be empty", ErrInvalidInput)
	}
	group, err := repository.GetGroupByName(db, groupname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: group '%s' not found", ErrNotFound, groupname)
		}
		return err
	}
	user, err := repository.GetUserByName(db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user '%s' not found", ErrNotFound, username)
		}
		return err
	}
	err = repository.AddMember(db, group.GroupID, user.UserID)
	if err != nil {
		if isDuplicateGroupMemberError(err) {
			return fmt.Errorf("%w: user '%s' is already a member of group '%s'", ErrConflict, username, groupname)
		}
		return err
	}
	return nil
}

func RemoveMemberFromGroup(db *sql.DB, groupname string, username string) error {
	groupname = strings.TrimSpace(groupname)
	username = strings.TrimSpace(username)

	if groupname == "" {
		return fmt.Errorf("%w: group name cannot be empty", ErrInvalidInput)
	}
	if username == "" {
		return fmt.Errorf("%w: username cannot be empty", ErrInvalidInput)
	}
	group, err := repository.GetGroupByName(db, groupname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: group '%s' not found", ErrNotFound, groupname)
		}
		return err
	}
	user, err := repository.GetUserByName(db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user '%s' not found", ErrNotFound, username)
		}
		return err
	}
	err = repository.RemoveMember(db, group.GroupID, user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user '%s' is not a member of group '%s'", ErrNotFound, username, groupname)
		}
		return err
	}
	return nil
}

func ListGroupMembers(db *sql.DB, groupname string) ([]models.User, error) {
	groupname = strings.TrimSpace(groupname)
	if groupname == "" {
		return nil, fmt.Errorf("%w: group name cannot be empty", ErrInvalidInput)
	}
	group, err := repository.GetGroupByName(db, groupname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: group '%s' not found", ErrNotFound, groupname)
		}
		return nil, err
	}
	members, err := repository.GetGroupMembers(db, group.GroupID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func ListGroups(db *sql.DB) ([]models.Group, error) {
	groups, err := repository.GetAllGroups(db)
	if err != nil {
		return nil, err
	}
	return groups, nil
}
