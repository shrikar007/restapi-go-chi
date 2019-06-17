package requests

import (
	"errors"
	"net/http"
)

type UpdateExpenseRequest struct {
	*CreateExpenseRequest
}

func (u *UpdateExpenseRequest) Bind(r *http.Request) error {
	if u.Id == 0 {
		return errors.New("id is empty or invalid")
	}

	return u.CreateExpenseRequest.Bind(r)
}
