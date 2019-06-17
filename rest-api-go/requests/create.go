package requests

import (
	"errors"
	"github.com/shrikar007/rest-api-go/structs"
	"net/http"
)

type CreateExpenseRequest struct {
	*structs.Expense
}
func (c *CreateExpenseRequest) Bind(r *http.Request) error {
	if c.Description == "" {
		return errors.New("description is either empty or invalid")
	}
	if c.Amount == 0 {
		return errors.New("amount is either empty or invalid")
	}

	if c.Type == "" {
		return errors.New("type is either empty or invalid")
	}

	return nil
}
