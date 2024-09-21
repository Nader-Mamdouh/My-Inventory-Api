package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)
var tst APP
func TestMain(m *testing.M){
	err:=tst.Initialise(DBUser, DBPassword, "test")
	if err !=nil{
		log.Fatal("error while creating database")
	}
	createtable()
	m.Run()
}
func createtable(){
	createtablequery:=`CREATE TABLE IF NOT EXISTS products(
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	quantity INT,
	price float(10,7),
	PRIMARY KEY (id));`

	_,err:=tst.DB.Exec(createtablequery)
	if err !=nil{
		log.Fatal(err)
	}
}

func ClearTable(){
	tst.DB.Exec("DELETE from products;")
	tst.DB.Exec("ALTER TABLE products AUTO_INCREMENT=1")
}
func AddProduct(name string,quantity int,price float64){
	query:= fmt.Sprintf("insert into products(name,quantity,price) values('%v',%v,%v)",name,quantity,price)
	tst.DB.Exec(query)
}
func TestGetPoduct(t *testing.T){
	ClearTable()
	AddProduct("Keyboeard",100,400.00)
	request,_:=http.NewRequest("GET","/products/1",nil)
	response:=sendRequest(request)
	checkstatuscode(t,http.StatusOK,response.Code)
}
func TestPostPoduct(t *testing.T){
	ClearTable()
	var product =[]byte(`{"name":"chair","quantity":100,"price":200.00}`)
	req,_:=http.NewRequest("POST","/products/",bytes.NewBuffer(product))
	req.Header.Set("Content-Type","application/json")
	response:=sendRequest(req)
	checkstatuscode(t,http.StatusCreated,response.Code)
}

func TestDeleteProduct(t *testing.T){
	ClearTable()
	AddProduct("mouse",100,400.00)
	request,_:=http.NewRequest("GET","/products/1",nil)
	response:=sendRequest(request)
	checkstatuscode(t,http.StatusOK,response.Code)

	request,_=http.NewRequest("DELETE","/products/1",nil)
	response=sendRequest(request)
	checkstatuscode(t,http.StatusOK,response.Code)

	request,_=http.NewRequest("DELETE","/products/1",nil)
	response=sendRequest(request)
	checkstatuscode(t,http.StatusNotFound,response.Code)
}

func TestUpdateProduct(t *testing.T){
	ClearTable()
	AddProduct("mouse",100,400.00)
	request,_:=http.NewRequest("GET","/products/1",nil)
	response:=sendRequest(request)
	checkstatuscode(t,http.StatusOK,response.Code)

	var oldvalue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(),&oldvalue)

	var product =[]byte(`{"name":"chair","quantity":100,"price":400.00}`)
	request,_=http.NewRequest("PUT","/products/1",bytes.NewBuffer(product))
	request.Header.Set("Content-Type","application/json")
	response=sendRequest(request)

	var newvalue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(),&newvalue)
	if oldvalue["id"]!=newvalue["id"]{
		t.Errorf("ecpected %v found %v",oldvalue["id"],newvalue["id"])
	}
	if oldvalue["name"]==newvalue["name"]{
		t.Errorf("ecpected %v found %v",oldvalue["name"],newvalue["name"])
	}
	if oldvalue["quantity"]!=newvalue["quantity"]{//we actually wanted to change thus value
		t.Errorf("ecpected %v found %v",oldvalue["quantity"],newvalue["quantity"])
	}	
	if oldvalue["price"]!=newvalue["price"]{
		t.Errorf("ecpected %v found %v",oldvalue["price"],newvalue["price"])
	}
}
func checkstatuscode(t *testing.T,expectedcode int,actualcode int){
	if expectedcode!=actualcode{
		t.Errorf("expected code : %v but actual code : %v",expectedcode,actualcode)
	}
}
func sendRequest(request *http.Request) *httptest.ResponseRecorder{
	recorder:=httptest.NewRecorder()
	tst.Router.ServeHTTP(recorder,request)
	return recorder
}