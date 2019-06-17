package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/shrikar007/rest-api-go/requests"
	"github.com/shrikar007/rest-api-go/response"
	"github.com/shrikar007/rest-api-go/structs"
	"log"
	"net/http"
	"strconv"
)

var expenses structs.Expenses

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", CreateExpense)
		r.Get("/", ListAllExpense)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ListOneExpense)
			r.Put("/", UpdateExpense)
			r.Delete("/", DeleteExpense)
		})
	})

	log.Fatal(http.ListenAndServe(":8083", r))
}

func CreateExpense(writer http.ResponseWriter,request *http.Request){

	var req requests.CreateExpenseRequest
	err:=render.Bind(request,&req)
	if err != nil {
		log.Println(err)
		return
	}
	expenses =append(expenses,*req.Expense)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	_, _ = fmt.Fprintln(writer, `{"success": true}`)
	render.Render(writer, request, response.Getoneresponse(req.Expense))

}

func ListOneExpense(writer http.ResponseWriter, request *http.Request) {

	for _,expense:=range expenses{
		if strconv.Itoa(expense.Id)==chi.URLParam(request, "id"){
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			render.Render(writer, request, response.Getoneresponse(&expense))
			return
		}
	}
}

func ListAllExpense(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	render.Render(writer,request,response.Getallresponse(&expenses))

}

func UpdateExpense(writer http.ResponseWriter, request *http.Request) {

	var req requests.UpdateExpenseRequest
	err:=render.Bind(request,&req)
	if err != nil {
		log.Println(err)
		return
	}
	for updateindex,expense:=range expenses{

		if  strconv.Itoa(expense.Id)==chi.URLParam(request, "id") {

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			expenses[updateindex] = *req.Expense
			_, _ = fmt.Fprintln(writer, `{"success": true}`)
			render.Render(writer, request, response.Getoneresponse(req.Expense))

		}
	}
}

func DeleteExpense(writer http.ResponseWriter, request *http.Request) {

	for deleteindex,expense:=range expenses{

		if strconv.Itoa(expense.Id)==chi.URLParam(request, "id"){
			expenses=append(expenses[:deleteindex],expenses[deleteindex+1:]...)

			_, _ = fmt.Fprintln(writer, `{"success": true}`)
			return
		}
	}


}
