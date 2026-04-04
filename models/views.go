package models

type ExpenseView struct {
	ExpenseID    int64  `json:"expense_id"`
	Amount       string `json:"amount"`
	Description  string `json:"description"`
	PaidByUserID int    `json:"paid_by_user_id"`
	GroupID      int    `json:"group_id"`
	CreatedAt    string `json:"created_at"`
}

type SplitView struct {
	ExpenseID int64  `json:"expense_id"`
	UserID    int    `json:"user_id"`
	Amount    string `json:"amount"`
	Settled   bool   `json:"settled"`
}

type ExpenseDetail struct {
	Description string        `json:"description"`
	Amount      int64         `json:"amount"`
	PaidBy      string        `json:"paid_by"`
	GroupName   string        `json:"group_name"`
	CreatedAt   string        `json:"created_at"`
	Splits      []SplitDetail `json:"splits"`
}

type SplitDetail struct {
	UserName string `json:"username"`
	Amount   int64  `json:"amount"`
	Settled  bool   `json:"settled"`
}

type UserBalance struct {
	User       User  `json:"user"`
	NetBalance int64 `json:"net_balance"`
}

type Transaction struct {
	From   User  `json:"from"`
	To     User  `json:"to"`
	Amount int64 `json:"amount"`
}

type SplitWithExpense struct {
	PaidByUserID int   `json:"paid_by_user_id"`
	UserID       int   `json:"user_id"`
	Amount       int64 `json:"amount"`
	Settled      bool  `json:"settled"`
}

type ExpenseDetailView struct {
	Description string            `json:"description"`
	Amount      string            `json:"amount"`
	PaidBy      string            `json:"paid_by"`
	GroupName   string            `json:"group_name"`
	CreatedAt   string            `json:"created_at"`
	Splits      []SplitDetailView `json:"splits"`
}

type SplitDetailView struct {
	UserName string `json:"username"`
	Amount   string `json:"amount"`
	Settled  bool   `json:"settled"`
}

type UserBalanceView struct {
	User       User   `json:"user"`
	NetBalance string `json:"net_balance"`
}

type TransactionView struct {
	From   User   `json:"from"`
	To     User   `json:"to"`
	Amount string `json:"amount"`
}
