package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	QUANTITY int     `json:"quantity"`
	PRICE    float64 `json:"price"`
}

func getproducts(db *sql.DB) ([]product, error) {
	query := "SELECT * FROM products;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()  

	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.QUANTITY, &p.PRICE) 
		if err != nil {
			return nil, err  
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err  
	}

	return products, nil  
}

func (p *product)getproduct(db *sql.DB)error{
	query:= fmt.Sprintf("select name , quantity , price from products where id=%v",p.ID)
	row:=db.QueryRow(query)
	err := row.Scan(&p.Name, &p.QUANTITY, &p.PRICE)
	if err != nil {
		return err  
	}
	return nil
}

func (p *product)createproduct(db *sql.DB)error{
	query:= fmt.Sprintf("insert into products(name,quantity,price) values('%v',%v,%v)",p.Name,p.QUANTITY,p.PRICE)
	result,err:=db.Exec(query)
	if err != nil {
		return err  
	}
	id,err:=result.LastInsertId()
	p.ID=int(id)
	return nil
}

func (p *product)updateproduct(db *sql.DB)error{
	query:= fmt.Sprintf("update products set name='%v' , quantity=%v , price=%v where id=%v",p.Name,p.QUANTITY,p.PRICE,p.ID)
	result,err:=db.Exec(query)
	rows_affected,err:=result.RowsAffected()
	if rows_affected==0{
		return errors.New("No Such Rows Affected")
	}
	return err
	
}
func (p *product)deleteproduct(db *sql.DB)error{
	query:= fmt.Sprintf("delete from products where id = %v",p.ID)
	result,err:=db.Exec(query)
	rows_affected,err:=result.RowsAffected()
	if rows_affected==0{
		return errors.New("No Such Rows Affected")
	}
	
	return err
	
}