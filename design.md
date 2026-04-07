# DevLedger — Design Document

## Overview

DevLedger is a local-first CLI + REST API tool for tracking and splitting shared infrastructure costs among developer teams. This document explains the key architectural decisions made during development.

---

## Architecture
```text
CLI (Cobra)
    ↓ HTTP calls (localhost:38080)
API Layer (Chi router + handlers)
    ↓
Service Layer (business logic)
    ↓
Repository Layer (DB queries)
    ↓
SQLite (devledger.db)
```

Each layer has one responsibility:
- **CLI** — parses commands, makes HTTP calls, displays results
- **API** — routes requests, parses input, returns JSON
- **Service** — orchestrates repository calls, applies business logic
- **Repository** — raw DB operations, no logic

---

## Key Decisions

### 1. SQLite over PostgreSQL

SQLite was chosen deliberately for two reasons:

- DevLedger is a local-first tool — a separate DB process adds unnecessary setup complexity for small dev teams
- `modernc.org/sqlite` is a pure Go port with no CGo dependency, meaning `go build` works on any platform without a C compiler

The tradeoff is no concurrent writes at scale, which is acceptable for this use case.

---

### 2. Monetary Values Stored in Paise (int64)

All amounts are stored as `int64` paise internally rather than `float64` rupees. ₹30.00 is stored as `3000`.

This eliminates floating-point precision bugs common in financial applications:
```
0.1 + 0.2 = 0.30000000000000004  ← float64
100 + 200 = 300                  ← int64, exact
```

The `shopspring/decimal` library handles parsing user input ("30.00" → 3000) and formatting output (3000 → "30.00").

---

### 3. Debt Simplification Algorithm

The core algorithm simplifies debts and reduces the number of transactions needed to settle all debts in a group.

**Approach — Greedy net balance:**

1. Compute net balance per member:
```text
netBalance = totalPaid - totalOwed
positive → creditor (is owed money)
negative → debtor (owes money)
```

2. Sort creditors and debtors by amount (highest first)

3. Greedily match largest debtor with largest creditor:
```text
transfer = min(debtor.amount, creditor.amount)
reduce both by transfer amount
advance pointer of whichever reaches zero
repeat until all settled
```

**Complexity:** O(n log n) — dominated by the sorting step

**Why not exact minimum?**
The theoretically optimal minimum transaction problem is NP-hard. The greedy approach gives near-optimal results in O(n log n) which is more than sufficient for realistic team sizes.

---

### 4. Database Transactions for Atomic Expense Creation

When adding an expense, two tables are written to — `Expenses` and `Splits`. If split creation fails halfway, the expense would exist without complete splits, corrupting balance calculations.

Solved using a database transaction:
```text
BEGIN TRANSACTION
  INSERT INTO Expenses → get expenseID
  INSERT INTO Splits × N members
COMMIT (or ROLLBACK if anything fails)
```

Either all splits are created or none are — the DB stays consistent.

---

### 5. N+1 Query Pattern — Fixed

Two functions originally had N+1 query patterns:

**`GetExpenseInDetail`** — originally fetched username per split in a loop:
```text
1 query → get splits
N queries → get username per split
```

Fixed with a single JOIN:
```sql
SELECT u.UserName, s.Amount, s.Settled
FROM Splits s
JOIN Users u ON s.UserID = u.UserID
WHERE s.ExpenseID = ?
```

**`CalculateBalances`** — originally fetched splits per expense in a loop:
```text
1 query → get expenses
N queries → get splits per expense
```

Fixed with a single JOIN:
```sql
SELECT e.PaidByUserID, s.UserID, s.Amount, s.Settled
FROM Splits s
JOIN Expenses e ON s.ExpenseID = e.ExpenseID
WHERE e.GroupID = ?
```

Result: balance calculation reduced from O(n) queries to O(1).

---

### 6. Soft Delete for Users

Users are never hard deleted. Instead a `DeletedAt` timestamp is set:
```sql
UPDATE Users SET DeletedAt = CURRENT_TIMESTAMP WHERE UserID = ?
```

This preserves historical expense data. Without soft delete, deleting a user who paid expenses would leave orphaned records — "AWS Bill paid by ???".

Active user queries filter with `AND DeletedAt IS NULL`. Historical views use `GetUserByIDIncludingDeleted` to always show complete records.

---

### 7. Payer's Split Auto-Settled

When an expense is created, the payer's own split is automatically marked `Settled = true`. This prevents the payer from appearing in their own unsettled debts list and simplifies balance calculation — settled splits are skipped entirely.

---

### 8. DBTX Interface

Repository functions called inside transactions accept a `DBTX` interface instead of `*sql.DB`:
```go
type DBTX interface {
    Exec(query string, args ...any) (sql.Result, error)
    Query(query string, args ...any) (*sql.Rows, error)
    QueryRow(query string, args ...any) *sql.Row
}
```

Both `*sql.DB` and `*sql.Tx` implement this interface, allowing the same repository functions to work inside and outside transactions without duplication.

---

### 9. Defensive Balance Calculation

The balance calculation skips a split if either:
- `split.Settled = true` — already paid
- `split.UserID == split.PaidByUserID` — payer's own share

The second condition is redundant in normal flow since payer splits are always auto-settled. It acts as a defensive check against data inconsistency — if a payer's split is accidentally marked unsettled, the balance calculation still produces correct results.

---

## Known Limitations

- **Equal splits only** — custom unequal splits not yet supported
- **No partial payments** — splits are either fully settled or not
- **Single currency** — no multi-currency support
- **No database migrations** — schema changes require manual DB deletion
- **Usernames are unique** — two users cannot share the same name
- **Server restarts per command** — no persistent daemon process

---

## Future Enhancements

- Custom expense splits (unequal amounts)
- Partial payment tracking
- Edit expense description and amount
- Database migrations using `golang-migrate`
- Idle timeout — auto shutdown server after inactivity
- Multi-currency support
- SQLite error code checking instead of string matching
