package models

type ExpenseDetail struct {
	Description string
	Amount      int64
	PaidBy      string
	GroupName   string
	CreatedAt   string
	Splits      []SplitDetail
}

type SplitDetail struct {
	UserName string
	Amount   int64
	Settled  bool
}

type UserBalance struct {
	User       User
	NetBalance int64
}

type Transaction struct {
	From   User
	To     User
	Amount int64
}
