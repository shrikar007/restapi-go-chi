package response

import(
	"github.com/shrikar007/rest-api-go/structs"
	"net/http"
)

type Getallstruct struct {
	 *structs.Expenses
}

func (Getallstruct) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Getallresponse(expenses *structs.Expenses) *Getallstruct{
	return &Getallstruct{Expenses: expenses}
}

