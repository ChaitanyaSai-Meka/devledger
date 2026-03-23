package models

type User struct {
	UserID   int
	UserName string
}

type Group struct {
	GroupID   int
	GroupName string
}

type GroupMember struct {
	GroupID int
	UserID  int
}

type Expense struct {
	ExpenseID    int64
	Amount       int64
	Description  string
	PaidByUserID int
	GroupID      int
	CreatedAt    string
}

type Split struct {
	ExpenseID int64
	UserID    int
	Amount    int64
	Settled   bool
}
