package crud_interface

import (
	"net/http"
)

type Database interface {
	CreateExpense(w http.ResponseWriter,r *http.Request)
	UpdateExpense(w http.ResponseWriter,r *http.Request)
	DeleteExpense(w http.ResponseWriter,r *http.Request)
	GetId(w http.ResponseWriter,r *http.Request)
	GetAll(w http.ResponseWriter,r *http.Request)
	ExpenseCtx(next http.Handler) http.Handler
}
