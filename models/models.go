package models

type User struct {
	UserID string
	UserName string
}

type Group struct {
	GroupID string
	GroupName string
}

type GroupMember  struct {
	GroupID string
	UserID string
}

type Expense struct {
	ExpenseID string
	Amount float32
	Description string
	PaidByUserID string
	GroupID string
}

type Split struct {
	ExpenseID string
	UserID string
	Amount float32
	Settled bool
}



