package service

import (
	"fmt"
	"testing"

	"github.com/ChaitanyaSai-Meka/devledger/models"
)

func TestSimplifyDebts(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		transactions := SimplifyDebts(nil)
		assertTransactions(t, transactions, nil)
	})

	t.Run("everyone balanced", func(t *testing.T) {
		balances := []models.UserBalance{
			newBalance("Chaitanya", 0),
			newBalance("Ravi", 0),
			newBalance("Arjun", 0),
		}

		transactions := SimplifyDebts(balances)
		assertTransactions(t, transactions, nil)
	})

	t.Run("simple A owes B", func(t *testing.T) {
		balances := []models.UserBalance{
			newBalance("A", -500),
			newBalance("B", 500),
		}

		transactions := SimplifyDebts(balances)
		assertTransactions(t, transactions, []transactionExpectation{
			{from: "A", to: "B", amount: 500},
		})
	})

	t.Run("single debtor single creditor unequal", func(t *testing.T) {
		balances := []models.UserBalance{
			newBalance("A", 1000),
			newBalance("B", -600),
			newBalance("C", -400),
		}

		transactions := SimplifyDebts(balances)
		if len(transactions) != 2 {
			t.Fatalf("expected 2 transactions, got %d", len(transactions))
		}

		var total int64
		for _, tx := range transactions {
			total += tx.Amount
		}
		if total != 1000 {
			t.Fatalf("expected total 1000, got %d", total)
		}
	})

	t.Run("simulation", func(t *testing.T) {
		balances := []models.UserBalance{
			newBalance("Chaitanya", 2000),
			newBalance("Ravi", -1000),
			newBalance("Arjun", -1000),
		}

		transactions := SimplifyDebts(balances)
		assertTransactions(t, transactions, []transactionExpectation{
			{from: "Ravi", to: "Chaitanya", amount: 1000},
			{from: "Arjun", to: "Chaitanya", amount: 1000},
		})
	})

	t.Run("complex case", func(t *testing.T) {
		balances := []models.UserBalance{
			newBalance("A", 700),
			newBalance("B", 300),
			newBalance("C", -400),
			newBalance("D", -600),
		}

		transactions := SimplifyDebts(balances)
		assertTransactions(t, transactions, []transactionExpectation{
			{from: "D", to: "A", amount: 600},
			{from: "C", to: "A", amount: 100},
			{from: "C", to: "B", amount: 300},
		})
	})
}

type transactionExpectation struct {
	from   string
	to     string
	amount int64
}

func newBalance(username string, netBalance int64) models.UserBalance {
	return models.UserBalance{
		User: models.User{
			UserName: username,
		},
		NetBalance: netBalance,
	}
}

func assertTransactions(t *testing.T, got []models.Transaction, want []transactionExpectation) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %d transactions, got %d: %#v", len(want), len(got), got)
	}

	gotCounts := make(map[string]int)
	var gotTotal int64
	for _, tx := range got {
		gotCounts[transactionKey(tx.From.UserName, tx.To.UserName, tx.Amount)]++
		gotTotal += tx.Amount
	}

	wantCounts := make(map[string]int)
	var wantTotal int64
	for _, tx := range want {
		wantCounts[transactionKey(tx.from, tx.to, tx.amount)]++
		wantTotal += tx.amount
	}

	if gotTotal != wantTotal {
		t.Fatalf("expected total transferred amount %d, got %d", wantTotal, gotTotal)
	}

	if len(gotCounts) != len(wantCounts) {
		t.Fatalf("expected transaction set %v, got %v", wantCounts, gotCounts)
	}

	for key, count := range wantCounts {
		if gotCounts[key] != count {
			t.Fatalf("expected transaction set %v, got %v", wantCounts, gotCounts)
		}
	}
}

func transactionKey(from string, to string, amount int64) string {
	return fmt.Sprintf("%s->%s:%d", from, to, amount)
}
