package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shrikar007/rest-api-go/dberror"
	"github.com/shrikar007/rest-api-go/requests"
	"github.com/shrikar007/rest-api-go/response"
	"github.com/shrikar007/rest-api-go/structs"
	"log"
	"net/http"
	"time"
)

var expenses structs.Expenses
var expense structs.Expense

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
			r.Use(ExpenseCtx)
			r.Get("/", ListOneExpense)
			r.Put("/", UpdateExpense)
			r.Delete("/", DeleteExpense)
		})
	})

	log.Fatal(http.ListenAndServe(":8083", r))
}
func ExpenseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		expenseID := chi.URLParam(r, "id")
		Db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")
		if err != nil {
			err=errors.New("unable to open database")
			render.Render(w, r, dberror.ErrRender(err))

		}
		db:=Db.Table("expenses").Where("id = ?",expenseID ).Find(&expense)
		ctx := context.WithValue(r.Context(), "expense", db )
		next.ServeHTTP(w, r.WithContext(ctx))


		if db.RowsAffected==0{
			err=errors.New("ID not Found")
			render.Render(w, r, dberror.ErrRender(err))
		}
	})
}


func CreateExpense(writer http.ResponseWriter,request *http.Request){

	var req requests.CreateExpenseRequest
	err:=render.Bind(request,&req)
	if err != nil {
		log.Println(err)
		return
	}
	expense=*req.Expense
	expense.CreatedOn=time.Now()
	expense.UpdatedOn=time.Now()

	Db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")
	if err != nil {
		err=errors.New("unable to open database")
		render.Render(writer, request, dberror.ErrRender(err))

	}
	db1:=Db.Create(&expense)
	if db1.RowsAffected!=0{
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		return

	}else{
		err:=errors.New("Unable to update")
		render.Render(writer,request,dberror.ErrRender(err))
		return
	}
	Db.Close()

}

func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
	db:=request.Context().Value("expense").(*gorm.DB)

	Db:=db.Find(&expense)
	if db.RowsAffected!=0{
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		render.Render(writer, request, response.Getoneresponse(&expense))

	}else{
		err:=errors.New("Unable fetch data")
		render.Render(writer,request,dberror.ErrRender(err))
		return
	}
	Db.Close()
}

func ListAllExpense(writer http.ResponseWriter, request *http.Request) {

	Db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println("invalid")
	}
	db1:=Db.Find(&expenses)

	if db1.RowsAffected!=0{
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		render.Render(writer,request,response.Getallresponse(&expenses))
		return

	}else{
		err:=errors.New("Unable to update")
		render.Render(writer,request,dberror.ErrRender(err))
		return
	}

	Db.Close()


}


func UpdateExpense(writer http.ResponseWriter, request *http.Request) {
	db:=request.Context().Value("expense").(*gorm.DB)

	var req requests.UpdateExpenseRequest
	err:=render.Bind(request,&req)
	if err != nil {
		log.Println(err)
		return
	}
	temp:=*req.Expense
	temp.UpdatedOn=time.Now()

	if err != nil {
		log.Println(err)
		return
	}
	db1:=db.Update(&temp)
	if db1.RowsAffected!=0{
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		return

	}else{
		err:=errors.New("Unable to update")
		render.Render(writer,request,dberror.ErrRender(err))
		return
	}

   db.Close()
}

func DeleteExpense(writer http.ResponseWriter, request *http.Request) {
	db:=request.Context().Value("expense").(*gorm.DB)

	Db:=db.Delete(&expense)
	if Db.RowsAffected!=0{
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		return

	}else{
		err:=errors.New("Unable to delete")
		render.Render(writer,request,dberror.ErrRender(err))
		return

	}

}
