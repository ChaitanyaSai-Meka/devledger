package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ChaitanyaSai-Meka/devledger/models"
	"github.com/ChaitanyaSai-Meka/devledger/repository"
)

func AddExpense(db *sql.DB, groupname string, description string, paidbyusername string, amount int64) error {
	groupname = strings.TrimSpace(groupname)
	description = strings.TrimSpace(description)
	paidbyusername = strings.TrimSpace(paidbyusername)
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if groupname == "" {
		return errors.New("group name cannot be empty")
	}
	if description == "" {
		return errors.New("description cannot be empty")
	}
	if paidbyusername == "" {
		return errors.New("paid by username cannot be empty")
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	group, err := repository.GetGroupByName(tx, groupname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("group '%s' not found", groupname)
		}
		return err
	}
	user, err := repository.GetUserByName(tx, paidbyusername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user '%s' not found", paidbyusername)
		}
		return err
	}
	members, err := repository.GetGroupMembers(tx, group.GroupID)
	if err != nil {
		return err
	}
	if len(members) == 0 {
		return fmt.Errorf("group '%s' has no members", groupname)
	}
	isMember := false
	for _, member := range members {
		if member.UserID == user.UserID {
			isMember = true
			break
		}
	}
	if !isMember {
		return fmt.Errorf("user '%s' is not a member of group '%s'", paidbyusername, groupname)
	}
	expense := models.Expense{
		Description:  description,
		Amount:       amount,
		PaidByUserID: user.UserID,
		GroupID:      group.GroupID,
	}
	splitAmount := amount / int64(len(members))
	remainder := amount % int64(len(members))

	expenseID, err := repository.CreateExpense(tx, expense)
	if err != nil {
		return err
	}

	for _, member := range members {
		memberAmount := splitAmount
		if member.UserID == user.UserID {
			memberAmount += remainder
		}
		split := models.Split{
			ExpenseID: expenseID,
			UserID:    member.UserID,
			Amount:    memberAmount,
			Settled:   member.UserID == user.UserID,
		}
		err = repository.CreateSplit(tx, split)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func ListExpensesByGroup(db *sql.DB, groupname string) ([]models.Expense, error) {
	groupname = strings.TrimSpace(groupname)
	if groupname == "" {
		return nil, errors.New("group name cannot be empty")
	}
	group, err := repository.GetGroupByName(db, groupname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("group '%s' not found", groupname)
		}
		return nil, err
	}
	expenses, err := repository.GetExpensesByGroupID(db, group.GroupID)
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func ListExpensesByUser(db *sql.DB, username string) ([]models.Expense, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	user, err := repository.GetUserByName(db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user '%s' not found", username)
		}
		return nil, err
	}
	expenses, err := repository.GetExpensesByUserID(db, user.UserID)
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func DeleteExpense(db *sql.DB, expenseID int64) error {
	err := repository.DeleteExpenseByID(db, expenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("expense with ID '%d' not found", expenseID)
		}
		return err
	}
	return nil
}

func SettleExpense(db *sql.DB, expenseID int64, username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username cannot be empty")
	}
	user, err := repository.GetUserByName(db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user '%s' not found", username)
		}
		return err
	}
	err = repository.SettleSplit(db, expenseID, user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user '%s' has no unsettled split for expense ID %d", username, expenseID)
		}
		return err
	}
	return nil
}

func ListUnsettledSplits(db *sql.DB, username string) ([]models.Split, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	user, err := repository.GetUserByName(db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user '%s' not found", username)
		}
		return nil, err
	}
	splits, err := repository.GetUnsettledSplitsByUserID(db, user.UserID)
	if err != nil {
		return nil, err
	}
	return splits, nil
}

func GetExpenseInDetail(db *sql.DB, expenseID int64) (models.ExpenseDetail, error) {
	expense, err := repository.GetExpenseByID(db, expenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ExpenseDetail{}, fmt.Errorf("expense with ID %d not found", expenseID)
		}
		return models.ExpenseDetail{}, err
	}
	user, err := repository.GetUserByID(db, expense.PaidByUserID)
	if err != nil {
		return models.ExpenseDetail{}, err
	}
	group, err := repository.GetGroupByID(db, expense.GroupID)
	if err != nil {
		return models.ExpenseDetail{}, err
	}

	splitsDetails, err := repository.GetSplitsWithUsersByExpenseID(db, expenseID)
	if err != nil {
		return models.ExpenseDetail{}, err
	}
	details := models.ExpenseDetail{
		Description: expense.Description,
		Amount:      expense.Amount,
		PaidBy:      user.UserName,
		GroupName:   group.GroupName,
		CreatedAt:   expense.CreatedAt,
		Splits:      splitsDetails,
	}
	return details, nil
}
