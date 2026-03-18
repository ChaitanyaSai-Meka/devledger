package models

type User struct {
	UserID int
	UserName string
}

type Group struct {
	GroupID int
	GroupName string
}

type GroupMember  struct {
	GroupID int
	UserID int
}

type Expense struct {
	ExpenseID int
	Amount float32
	Description string
	PaidByUserID int
	GroupID int
}

type Split struct {
	ExpenseID int
	UserID int
	Amount float32
	Settled bool
}



