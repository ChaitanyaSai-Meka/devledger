package repository

import (
	"database/sql"
	"github.com/ChaitanyaSai-Meka/devledger/models"
)

func CreateExpense(tx *sql.Tx, expense models.Expense) (int64,error) {
	result, err := tx.Exec("INSERT INTO Expenses (Amount, Description, PaidByUserID, GroupID) VALUES (?,?,?,?)", expense.Amount, expense.Description, expense.PaidByUserID, expense.GroupID)
	if err != nil {
		return 0,err
	}

	id,err := result.LastInsertId()
	if err!= nil {
		return 0,err
	}
	return id,nil
}

func GetExpenseByID(db *sql.DB, expenseID int) (models.Expense, error) {
	var expense models.Expense
	err := db.QueryRow(
		"SELECT ExpenseID, Amount, Description, PaidByUserID, GroupID, CreatedAt FROM Expenses WHERE ExpenseID = ?",
		expenseID,
	).Scan(&expense.ExpenseID, &expense.Amount, &expense.Description, &expense.PaidByUserID, &expense.GroupID, &expense.CreatedAt)
	if err != nil {
		return models.Expense{}, err
	}
	return expense, nil
}

func GetExpensesByGroupID(db *sql.DB, groupID int) ([]models.Expense, error) {
	rows, err := db.Query("SELECT ExpenseID, Amount, Description, PaidByUserID, GroupID, CreatedAt FROM Expenses WHERE GroupID = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []models.Expense{}
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ExpenseID, &expense.Amount, &expense.Description, &expense.PaidByUserID, &expense.GroupID, &expense.CreatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}

func GetExpensesByUserID(db *sql.DB, userID int) ([]models.Expense, error) {
	rows, err := db.Query("SELECT ExpenseID, Amount, Description, PaidByUserID, GroupID, CreatedAt FROM Expenses WHERE PaidByUserID = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []models.Expense{}
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ExpenseID, &expense.Amount, &expense.Description, &expense.PaidByUserID, &expense.GroupID, &expense.CreatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}

func DeleteExpenseByID(db *sql.DB, expenseID int) error {
	result, err := db.Exec("DELETE FROM Expenses WHERE ExpenseID = ?", expenseID)
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
