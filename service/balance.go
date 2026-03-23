package service

import (
	"database/sql"
	"errors"
	"github.com/ChaitanyaSai-Meka/devledger/models"
	"github.com/ChaitanyaSai-Meka/devledger/repository"
	"sort"
)

type balanceEntry struct {
	user   models.User
	amount int64
}

func CalculateBalances(db *sql.DB, groupName string) ([]models.UserBalance, error) {
	group, err := repository.GetGroupByName(db, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("group not found")
		}
		return nil, err
	}
	expenses, err := repository.GetExpensesByGroupID(db, group.GroupID)
	if err != nil {
		return nil, err
	}
	balance := make(map[int]int64)
	for _, expense := range expenses {
		balance[expense.PaidByUserID] += expense.Amount
		splits, err := repository.GetSplitsByExpenseID(db, expense.ExpenseID)
		if err != nil {
			return nil, err
		}
		for _, split := range splits {
			if !split.Settled {
				balance[split.UserID] -= split.Amount
			}
		}
	}
	members, err := repository.GetGroupMembers(db, group.GroupID)
	if err != nil {
		return nil, err
	}
	var result []models.UserBalance
	for _, member := range members {
		result = append(result, models.UserBalance{
			User:       member,
			NetBalance: balance[member.UserID],
		})
	}
	return result, nil
}

func SimplifyDebts(balances []models.UserBalance) []models.Transaction {
	creditors := []balanceEntry{}
	debtors := []balanceEntry{}
	for _, balance := range balances {
		if balance.NetBalance > 0 {
			creditors = append(creditors, balanceEntry{
				user:   balance.User,
				amount: balance.NetBalance,
			})
		} else {
			debtors = append(debtors, balanceEntry{
				user:   balance.User,
				amount: -balance.NetBalance,
			})
		}
	}

	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].amount > creditors[j].amount
	})
	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].amount > debtors[j].amount
	})
	var transactions []models.Transaction
	i, j := 0, 0
	for i < len(creditors) && j < len(debtors) {
		amount := min(creditors[i].amount, debtors[j].amount)
		transactions = append(transactions, models.Transaction{
			From:   debtors[j].user,
			To:     creditors[i].user,
			Amount: amount,
		})
		creditors[i].amount -= amount
		debtors[j].amount -= amount
		if creditors[i].amount == 0 {
			i++
		}
		if debtors[j].amount == 0 {
			j++
		}
	}
	return transactions
}
