package structs

import "google.golang.org/genproto/googleapis/type/date"

type Expense struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	CreatedOn   date.Date `json:"created_on" `
	UpdatedOn   date.Date `json:"updated_on"`
}

type Expenses []Expense
