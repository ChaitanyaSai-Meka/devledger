package models

type User struct {
	UserID string
	UserName string
	Balance int
}

type Group struct {
	GroupID string
	GroupName string
	MembersIDs []string
}

type Expense struct {
	ExpenseID string
	Amount int
	Service string
	PaidBy string
	GroupID string
}

type Split struct {
	UserID string
	Amount int
}



