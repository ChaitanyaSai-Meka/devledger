package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ChaitanyaSai-Meka/devledger/models"
	"github.com/ChaitanyaSai-Meka/devledger/money"
	"github.com/ChaitanyaSai-Meka/devledger/repository"
)

func AddExpense(db *sql.DB, groupname string, description string, paidbyusername string, amount string) error {
	groupname = strings.TrimSpace(groupname)
	description = strings.TrimSpace(description)
	paidbyusername = strings.TrimSpace(paidbyusername)
	convertedamount, err := money.ParseToMinorUnit(amount)
	if err != nil {
		return fmt.Errorf("%w: invalid amount: %v", ErrInvalidInput, err)
	}
	if groupname == "" {
		return fmt.Errorf("%w: group name cannot be empty", ErrInvalidInput)
	}
	if description == "" {
		return fmt.Errorf("%w: description cannot be empty", ErrInvalidInput)
	}
	if paidbyusername == "" {
		return fmt.Errorf("%w: paid by username cannot be empty", ErrInvalidInput)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	group, err := repository.GetGroupByName(tx, groupname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: group '%s' not found", ErrNotFound, groupname)
		}
		return err
	}
	user, err := repository.GetUserByName(tx, paidbyusername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user '%s' not found", ErrNotFound, paidbyusername)
		}
		return err
	}
	members, err := repository.GetGroupMembers(tx, group.GroupID)
	if err != nil {
		return err
	}
	if len(members) == 0 {
		return fmt.Errorf("%w: group '%s' has no members", ErrConflict, groupname)
	}
	if convertedamount < int64(len(members)) {
		return fmt.Errorf("%w: amount is too small to split among %d members", ErrInvalidInput, len(members))
	}
	isMember := false
	for _, member := range members {
		if member.UserID == user.UserID {
			isMember = true
			break
		}
	}
	if !isMember {
		return fmt.Errorf("%w: user '%s' is not a member of group '%s'", ErrConflict, paidbyusername, groupname)
	}
	expense := models.Expense{
		Description:  description,
		Amount:       convertedamount,
		PaidByUserID: user.UserID,
		GroupID:      group.GroupID,
	}
	splitAmount := convertedamount / int64(len(members))
	remainder := convertedamount % int64(len(members))

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

func ListExpensesByGroup(db *sql.DB, groupname string) ([]models.ExpenseView, error) {
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
	expenses, err := repository.GetExpensesByGroupID(db, group.GroupID)
	if err != nil {
		return nil, err
	}
	return buildExpenseViews(expenses), nil
}

func ListExpensesByUser(db *sql.DB, username string) ([]models.ExpenseView, error) {
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
	expenses, err := repository.GetExpensesByUserID(db, user.UserID)
	if err != nil {
		return nil, err
	}
	return buildExpenseViews(expenses), nil
}

func DeleteExpense(db *sql.DB, expenseID int64) error {
	err := repository.DeleteExpenseByID(db, expenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: expense with ID '%d' not found", ErrNotFound, expenseID)
		}
		return err
	}
	return nil
}

func SettleExpense(db *sql.DB, expenseID int64, username string) error {
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
	err = repository.SettleSplit(db, expenseID, user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user '%s' has no unsettled split for expense ID %d", ErrConflict, username, expenseID)
		}
		return err
	}
	return nil
}

func ListUnsettledSplits(db *sql.DB, username string) ([]models.SplitView, error) {
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
	splits, err := repository.GetUnsettledSplitsByUserID(db, user.UserID)
	if err != nil {
		return nil, err
	}
	return buildSplitViews(splits), nil
}

func GetExpenseInDetail(db *sql.DB, expenseID int64) (models.ExpenseDetailView, error) {
	expense, err := repository.GetExpenseByID(db, expenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ExpenseDetailView{}, fmt.Errorf("%w: expense with ID %d not found", ErrNotFound, expenseID)
		}
		return models.ExpenseDetailView{}, err
	}
	user, err := repository.GetUserByIDIncludingDeleted(db, expense.PaidByUserID)
	if err != nil {
		return models.ExpenseDetailView{}, err
	}
	group, err := repository.GetGroupByID(db, expense.GroupID)
	if err != nil {
		return models.ExpenseDetailView{}, err
	}

	splitsDetails, err := repository.GetSplitsWithUsersByExpenseID(db, expenseID)
	if err != nil {
		return models.ExpenseDetailView{}, err
	}
	details := models.ExpenseDetailView{
		Description: expense.Description,
		Amount:      money.FormatFromMinorUnit(expense.Amount),
		PaidBy:      user.UserName,
		GroupName:   group.GroupName,
		CreatedAt:   expense.CreatedAt,
		Splits:      buildSplitDetailViews(splitsDetails),
	}
	return details, nil
}

func buildExpenseViews(expenses []models.Expense) []models.ExpenseView {
	views := make([]models.ExpenseView, 0, len(expenses))
	for _, expense := range expenses {
		views = append(views, models.ExpenseView{
			ExpenseID:    expense.ExpenseID,
			Amount:       money.FormatFromMinorUnit(expense.Amount),
			Description:  expense.Description,
			PaidByUserID: expense.PaidByUserID,
			GroupID:      expense.GroupID,
			CreatedAt:    expense.CreatedAt,
		})
	}
	return views
}

func buildSplitViews(splits []models.Split) []models.SplitView {
	views := make([]models.SplitView, 0, len(splits))
	for _, split := range splits {
		views = append(views, models.SplitView{
			ExpenseID: split.ExpenseID,
			UserID:    split.UserID,
			Amount:    money.FormatFromMinorUnit(split.Amount),
			Settled:   split.Settled,
		})
	}
	return views
}

func buildSplitDetailViews(splits []models.SplitDetail) []models.SplitDetailView {
	views := make([]models.SplitDetailView, 0, len(splits))
	for _, split := range splits {
		views = append(views, models.SplitDetailView{
			UserName: split.UserName,
			Amount:   money.FormatFromMinorUnit(split.Amount),
			Settled:  split.Settled,
		})
	}
	return views
}
