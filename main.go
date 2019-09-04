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

type Mysql struct {
	Db *gorm.DB
}

func main() {

	dba, err := gorm.Open("mysql", "root:root@tcp(sqldb:3306)/")
	dba.Exec("CREATE DATABASE IF NOT EXISTS"+" crudexpenses")
	dba.Close()

	db, err := gorm.Open("mysql", "root:root@tcp(sqldb:3306)/crudexpenses?charset=utf8&parseTime=True")

	if err != nil {
		fmt.Println(err)
	}
	if (!db.HasTable(&structs.Expense{})) {
		db.AutoMigrate(&structs.Expense{})
	}
	set := &Mysql{db}
	Init(set)
}
func Init(d crud_interface.Database) {
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
	log.Fatal(http.ListenAndServe(":8086", r))
}

func (db *Mysql) ExpenseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var temp structs.Expense
		expenseID := chi.URLParam(r, "id")
		DB:= db.Db.Table("expenses").Where("id = ?", expenseID).Find(&temp)
		//fmt.Println(temp)
		if DB.RowsAffected == 0{
			err:=errors.New("ID not Found")
			render.Render(w, r, dberror.ErrRender(err))
			return
		} else{
			ctx := context.WithValue(r.Context(), "expense", temp)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (db *Mysql) CreateExpense(writer http.ResponseWriter, request *http.Request) {
	//var expense structs.Expense

	var req requests.CreateExpenseRequest
	err := render.Bind(request, &req)
	if err != nil {
		log.Println(err)
		return
	}
	expense := *req.Expense
	expense.CreatedOn = time.Now()
	expense.UpdatedOn = time.Now()

	//	Db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")

	db1 := db.Db.Create(&expense)
	if db1.RowsAffected != 0 {
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		return

	} else {
		err := errors.New("Unable to update")
		render.Render(writer, request, dberror.ErrRender(err))
		return
	}
	//db.Db.Close()

}

func (db *Mysql) GetId(writer http.ResponseWriter, request *http.Request) {
//	var expense structs.Expense
	expen:= request.Context().Value("expense").(structs.Expense)
	fmt.Println(db.Db.Value)

	//Db1 := db.Db.Find(&expense)

		render.Render(writer, request, response.Getoneresponse(&expen))


	//Db.Close()
}

func (db *Mysql) GetAll(writer http.ResponseWriter, request *http.Request) {
	var expenses structs.Expenses

	//Db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")

	db1 := db.Db.Find(&expenses)

	if db1.RowsAffected != 0 {
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		render.Render(writer, request, response.Getallresponse(&expenses))
		return

	} else {
		err := errors.New("Unable to update")
		render.Render(writer, request, dberror.ErrRender(err))
		return
	}

	//db.Db.Close()

}

func (db *Mysql) UpdateExpense(writer http.ResponseWriter, request *http.Request) {
	expe := request.Context().Value("expense").(structs.Expense)

	var req requests.UpdateExpenseRequest
	err := render.Bind(request, &req)
	if err != nil {
		log.Println(err)
		return
	}

	temp:= *req.Expense
	temp.CreatedOn=expe.CreatedOn
	temp.UpdatedOn = time.Now()

	if err != nil {
		log.Println(err)
		return
	}
	db2 :=db.Db.Model(&expe).Update(&temp)

	if db2.RowsAffected != 0 {
		_, _ = fmt.Fprintln(writer, `{"success": true}`,expe)
		return

	} else {
		err := errors.New("Unable to update")
		render.Render(writer, request, dberror.ErrRender(err))
		return
	}


}

func (db *Mysql) DeleteExpense(writer http.ResponseWriter, request *http.Request) {
	//var expense structs.Expense
	exp:= request.Context().Value("expense").(structs.Expense)

	Db := db.Db.Delete(&exp)
	if Db.RowsAffected != 0 {
		_, _ = fmt.Fprintln(writer, `{"success": true}`)
		return
	} else {
		err := errors.New("Unable to delete")
		render.Render(writer, request, dberror.ErrRender(err))
		return
	}
}
