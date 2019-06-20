package structs

import (
	"time"
)

type Expense struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Typeofaccount        string    `json:"typeofaccount"`
	Amount      float64   `json:"amount"`
	CreatedOn   time.Time `json:"created_on" `
	UpdatedOn   time.Time `json:"updated_on"`
}

type Expenses []Expense
