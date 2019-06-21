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
	"github.com/shrikar007/rest-api-go/crud_interface"
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

type Mysql struct {

	Db *gorm.DB
}



func main() {

	db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")
	if err != nil {

	}
	sel:=&Mysql{db}

	Init(sel)
}
func Init(d crud_interface.Database){
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", d.CreateExpense)
		r.Get("/", d.GetAll)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(d.ExpenseCtx)
			r.Get("/", d.GetId)
			r.Put("/", d.UpdateExpense)
			r.Delete("/", d.DeleteExpense)
		})
	})

	log.Fatal(http.ListenAndServe(":8083", r))
}


func (db *Mysql)ExpenseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		expenseID := chi.URLParam(r, "id")
		//Db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")

		db:=db.Db.Table("expenses").Where("id = ?",expenseID ).Find(&expense)
		ctx := context.WithValue(r.Context(), "expense", db )
		next.ServeHTTP(w, r.WithContext(ctx))


		if db.RowsAffected==0{
			err=errors.New("ID not Found")
			render.Render(w, r, dberror.ErrRender(err))
		}
	})
}


func  (db *Mysql)CreateExpense(writer http.ResponseWriter,request *http.Request){

	var req requests.CreateExpenseRequest
	err:=render.Bind(request,&req)
	if err != nil {
		log.Println(err)
		return
	}
	expense=*req.Expense
	expense.CreatedOn=time.Now()
	expense.UpdatedOn=time.Now()

//	Db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")

	db1:=db.Db.Create(&expense)
	if db1.RowsAffected!=0{
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		return

	}else{
		err:=errors.New("Unable to update")
		render.Render(writer,request,dberror.ErrRender(err))
		return
	}
	//db.Db.Close()

}

func  (db *Mysql)GetId(writer http.ResponseWriter, request *http.Request) {
	db.Db=request.Context().Value("expense").(*gorm.DB)

	Db:=db.Db.Find(&expense)
	if Db.RowsAffected!=0{
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		render.Render(writer, request, response.Getoneresponse(&expense))

	}else{
		err:=errors.New("Unable fetch data")
		render.Render(writer,request,dberror.ErrRender(err))
		return
	}
	//Db.Close()
}

func (db *Mysql)GetAll(writer http.ResponseWriter, request *http.Request) {

	//Db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")

	db1:=db.Db.Find(&expenses)

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

	//db.Db.Close()


}


func (db *Mysql)UpdateExpense(writer http.ResponseWriter, request *http.Request) {
	db.Db=request.Context().Value("expense").(*gorm.DB)

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
	db2:=db.Db.Update(&temp)
	if db2.RowsAffected!=0{
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		return

	}else{
		err:=errors.New("Unable to update")
		render.Render(writer,request,dberror.ErrRender(err))
		return
	}

   //db.Close()
}

func (db *Mysql)DeleteExpense(writer http.ResponseWriter, request *http.Request) {
	db.Db=request.Context().Value("expense").(*gorm.DB)

	Db:=db.Db.Delete(&expense)
	if Db.RowsAffected!=0{
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		return

	}else{
		err:=errors.New("Unable to delete")
		render.Render(writer,request,dberror.ErrRender(err))
		return

	}

}
