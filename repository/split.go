package repository

import (
	"database/sql"
	"github.com/ChaitanyaSai-Meka/devledger/models"
)

func CreateSplit(tx *sql.Tx, split models.Split) error {
	_, err := tx.Exec("INSERT INTO Splits (ExpenseID, UserID, Amount, Settled) VALUES (?,?,?,?)", split.ExpenseID, split.UserID, split.Amount, split.Settled)
	if err != nil {
		return err
	}
	return nil
}

func GetSplitsByExpenseID(db *sql.DB, expenseID int64) ([]models.Split, error) {
	rows, err := db.Query("SELECT ExpenseID, UserID, Amount, Settled FROM Splits WHERE ExpenseID = ?", expenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	splits := []models.Split{}
	for rows.Next() {
		var split models.Split
		err := rows.Scan(&split.ExpenseID, &split.UserID, &split.Amount, &split.Settled)
		if err != nil {
			return nil, err
		}
		splits = append(splits, split)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return splits, nil
}

func GetUnsettledSplitsByUserID(db *sql.DB, userID int) ([]models.Split, error) {
	rows, err := db.Query("SELECT ExpenseID, UserID, Amount, Settled FROM Splits WHERE UserID = ? AND Settled = 0", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	splits := []models.Split{}
	for rows.Next() {
		var split models.Split
		err := rows.Scan(&split.ExpenseID, &split.UserID, &split.Amount, &split.Settled)
		if err != nil {
			return nil, err
		}
		splits = append(splits, split)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return splits, nil
}

func SettleSplit(db *sql.DB, expenseID int64, userID int) error {
	result, err := db.Exec("UPDATE Splits SET Settled = 1 WHERE ExpenseID = ? AND UserID = ? AND Settled = 0", expenseID, userID)
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

func GetSplitsWithUsersByExpenseID(db *sql.DB, expenseID int64) ([]models.SplitDetail, error) {
	rows, err := db.Query(`
        SELECT u.UserName, s.Amount, s.Settled
        FROM Splits s
        JOIN Users u ON s.UserID = u.UserID
        WHERE s.ExpenseID = ?
    `, expenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splits []models.SplitDetail
	for rows.Next() {
		var split models.SplitDetail
		if err := rows.Scan(&split.UserName, &split.Amount, &split.Settled); err != nil {
			return nil, err
		}
		splits = append(splits, split)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return splits, nil
}

func GetSplitsWithExpensesByGroupID(db *sql.DB, groupID int) ([]models.SplitWithExpense, error) {
	rows, err := db.Query(`
        SELECT e.PaidByUserID, s.UserID, s.Amount, s.Settled
        FROM Splits s
        JOIN Expenses e ON s.ExpenseID = e.ExpenseID
        WHERE e.GroupID = ?
    `, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SplitWithExpense
	for rows.Next() {
		var s models.SplitWithExpense
		if err := rows.Scan(&s.PaidByUserID, &s.UserID, &s.Amount, &s.Settled); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
