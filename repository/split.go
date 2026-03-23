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

func GetSplitsByExpenseID(db *sql.DB, expenseID int) ([]models.Split, error) {
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

func SettleSplit(db *sql.DB, expenseID int, userID int) error {
	result, err := db.Exec("UPDATE Splits SET Settled = 1 WHERE ExpenseID = ? AND UserID = ?", expenseID, userID)
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
