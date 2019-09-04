package response

import (
	"github.com/shrikar007/rest-api-go/structs"
	"net/http"
)




type GetOneStruct struct {
	*structs.Expense
}

func (GetOneStruct) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}



func Getoneresponse(expense *structs.Expense) *GetOneStruct {
	return &GetOneStruct{Expense: expense}

}
