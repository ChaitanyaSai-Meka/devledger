package models

type ExpenseDetail struct {
	Description string
	Amount int64
	PaidBy string
	GroupName string
	CreatedAt string
	Splits []SplitDetail
}

type SplitDetail struct {
	UserName string
	Amount int64
	Settled bool
}