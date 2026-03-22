package service

import (
	"database/sql"
	"errors"
	"strings"
	"fmt"

	"github.com/ChaitanyaSai-Meka/devledger/models"
	"github.com/ChaitanyaSai-Meka/devledger/repository"
)

func AddExpense(db *sql.DB, expense models.Expense) error {

}

func ListExpensesByGroup(db *sql.DB, groupname string) ([]models.Expense,error){
	groupname =strings.TrimSpace(groupname)
	if groupname == "" {
		return nil, errors.New("group name cannot be empty")
	}
	group,err := repository.GetGroupByName(db,groupname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("group '%s' not found", groupname)
		}
		return nil,err
	}
	expenses, err := repository.GetExpensesByGroupID(db, group.GroupID)
	if err !=nil{
		return nil,err
	}
	return expenses,nil
}

func ListExpensesByUser(db *sql.DB, username string) ([]models.Expense,error) {
	username =strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	user,err := repository.GetUserByName(db,username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user '%s' not found", username)
		}
		return nil,err
	}
	expenses, err := repository.GetExpensesByUserID(db, user.UserID)
	if err !=nil{
		return nil,err
	}
	return expenses,nil
}

func DeleteExpense(db *sql.DB, expenseID int) error {
	err:=repository.DeleteExpenseByID(db, expenseID)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("expense with ID '%d' not found", expenseID)
		}
		return err
	}
	return nil
}

func SettleExpense(db *sql.DB,expenseID int,username string) error {
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

func ListUnsettledSplits(db *sql.DB, username string) ([]models.Split, error){
	username =strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	user,err := repository.GetUserByName(db,username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user '%s' not found", username)
		}
		return nil,err
	}
	splits, err := repository.GetUnsettledSplitsByUserID(db, user.UserID)
	if err !=nil{
		return nil,err
	}
	return splits,nil
}

func GetExpenseInDetail(db *sql.DB, expenseID int) (models.ExpenseDetail, error){

}