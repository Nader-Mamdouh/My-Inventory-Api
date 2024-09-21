package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type APP struct {
	Router *mux.Router
	DB *sql.DB
}

func (app *APP) Initialise( DBUser string, DBPassword string, DBName string)error{
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DBUser, DBPassword, DBName)
    var err error
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
        fmt.Println("Error connecting to database:", err)
        return err
    }
	app.Router =mux.NewRouter().StrictSlash(true)
	app.handleRequests()
	return nil

}
func (app *APP)Run(address string){
	log.Fatal(http.ListenAndServe(address,app.Router))
}
func (app *APP) sendresponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func (app *APP) senderror(w http.ResponseWriter, statusCode int, err string) {
	errorMessage := map[string]string{"error": err}
	app.sendresponse(w, statusCode, errorMessage)  
}

func (app *APP) getproducts(w http.ResponseWriter, r *http.Request) {
	products, err := getproducts(app.DB)
	if err != nil {
		app.senderror(w, http.StatusInternalServerError, err.Error())  
		return 
	}
	app.sendresponse(w, http.StatusOK, products)  
}
func (app *APP) getproduct(w http.ResponseWriter, r *http.Request) {
	vars:=mux.Vars(r)
	key,err :=strconv.Atoi(vars["id"]) 
	if err !=nil{
		app.senderror(w,http.StatusBadRequest,"invalid product id")
		return
	}
	p:=product{ID:key}
	err =p.getproduct(app.DB)
	if err !=nil{
		switch err{
		case sql.ErrNoRows:
			app.senderror(w,http.StatusNotFound,"product not found")
		default:
			app.senderror(w,http.StatusInternalServerError,err.Error())	
		}
		return
	}
	app.sendresponse(w,http.StatusOK,p)
}

func (app *APP) createproduct(w http.ResponseWriter, r *http.Request){
	var p product
	err:=json.NewDecoder(r.Body).Decode(&p)
	if err !=nil{
		app.senderror(w,http.StatusBadRequest,"Invalid request payload")
		return
	}
	err=p.createproduct(app.DB)
	if err != nil {
		app.senderror(w, http.StatusInternalServerError, err.Error())  
		return 
	}
	app.sendresponse(w,http.StatusCreated,p)
}
func (app *APP) updateproduct(w http.ResponseWriter, r *http.Request){
	vars:=mux.Vars(r)
	key,err :=strconv.Atoi(vars["id"]) 
	if err !=nil{
		app.senderror(w,http.StatusBadRequest,"invalid product id")
		return
	}
	var p product
	err=json.NewDecoder(r.Body).Decode(&p)
	if err !=nil{
		app.senderror(w,http.StatusBadRequest,"Invalid request payload")
		return
	}
	p.ID=key
	err=p.updateproduct(app.DB)
	if err != nil {
		app.senderror(w, http.StatusInternalServerError, err.Error())  
		return 
	}
	app.sendresponse(w,http.StatusOK,p)
}

func (app *APP) deleteproduct(w http.ResponseWriter, r *http.Request){
	vars:=mux.Vars(r)
	key,err :=strconv.Atoi(vars["id"]) 
	if err !=nil{
		app.senderror(w,http.StatusBadRequest,"invalid product id")
		return
	}
	p:= product{ID:key}
	err=p.deleteproduct(app.DB)
	if err != nil {
        if err.Error() == "No Such Rows Affected" {
            app.senderror(w, http.StatusNotFound, "Product not found")
        } else {
            app.senderror(w, http.StatusInternalServerError, "Error deleting product")
        }
        return
    }
	app.sendresponse(w,http.StatusOK,map[string]string{"result":"product deleted succefully"})
}
func (app *APP)handleRequests(){
	app.Router.HandleFunc("/products",app.getproducts).Methods("GET")
	app.Router.HandleFunc("/products/{id}",app.getproduct).Methods("GET")
	app.Router.HandleFunc("/products/",app.createproduct).Methods("POST")
	app.Router.HandleFunc("/products/{id}",app.updateproduct).Methods("PUT")
	app.Router.HandleFunc("/products/{id}",app.deleteproduct).Methods("DELETE")
}