package models

type User struct {
	UserID    int    `json:"user_id"`
	UserName  string `json:"username"`
	DeletedAt string `json:"deleted_at,omitempty"`
}

type Group struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
}

type GroupMember struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}

type Expense struct {
	ExpenseID    int64  `json:"expense_id"`
	Amount       int64  `json:"amount"`
	Description  string `json:"description"`
	PaidByUserID int    `json:"paid_by_user_id"`
	GroupID      int    `json:"group_id"`
	CreatedAt    string `json:"created_at"`
}

type Split struct {
	ExpenseID int64 `json:"expense_id"`
	UserID    int   `json:"user_id"`
	Amount    int64 `json:"amount"`
	Settled   bool  `json:"settled"`
}
