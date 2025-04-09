package domain

import (
	"database/sql"
	"log"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type CustomerRepositoryDb struct {
	client    *sql.DB
}


func (s CustomerRepositoryDb)ById(id string) (*Customer, error){
	findOne := "SELECT customer_id, name, city, zipcode, date_of_birth, status from customers WHERE customer_id = ?"
	row := s.client.QueryRow(findOne, id)
	var c Customer
	err := row.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.DataofBirth, &c.Status)
	if err != nil{
		log.Println("error while scan customer, ", err.Error())
		return nil, err
	}
	return &c, nil
}


func (s CustomerRepositoryDb) FindAll() ([]Customer, error) {

	findAllSql := "SELECT customer_id, name, city, zipcode, date_of_birth, status from customers"
	rows, err := s.client.Query(findAllSql)
	if err != nil {
		log.Println("error while query customer table, ", err.Error())
		return nil, err
	}

	customers := make([]Customer, 0)
	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.DataofBirth, &c.Status)
		if err != nil {
			log.Println("error while query scanning customers, ", err.Error())
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}


func NewCustomerRepositoryDb() CustomerRepositoryDb {

	client, err := sql.Open("mysql", "root:0000@/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return CustomerRepositoryDb{client}
}
