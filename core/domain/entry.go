package domain

type Type string

const (
	Income  Type = "INCOME"
	Expense Type = "EXPENSE"
)

type Entry struct {
	ID          string
	Description string
	Amount      float64
	Tags        []Tag
	Type        Type
}
