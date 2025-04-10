package domain

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kal997/banking/errs"
	"github.com/kal997/banking/logger"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func (s CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	findOne := "SELECT customer_id, name, city, zipcode, date_of_birth, status from customers WHERE customer_id = ?"
	var c Customer
	err := s.client.Get(&c, findOne, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer Not Found")
		}
		logger.Error("error while scan customer, " + err.Error())
		return nil, errs.NewNotExpectedError("unexpected database error")
	}
	return &c, nil
}

func (s CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {

	var err error
	customers := make([]Customer, 0)

	if status == "" {
		findAllSql := "SELECT customer_id, name, city, zipcode, date_of_birth, status from customers"
		err = s.client.Select(&customers, findAllSql)
	} else {
		findAllSql := "SELECT customer_id, name, city, zipcode, date_of_birth, status from customers WHERE status = ?"
		err = s.client.Select(&customers, findAllSql, status)
	}

	if err != nil {
		logger.Error("error while query customer table, " + err.Error())
		return nil, errs.NewNotExpectedError("database Connection Error")
	}

	return customers, nil
}

func NewCustomerRepositoryDb() CustomerRepositoryDb {

	client, err := sqlx.Open("mysql", "root:0000@/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return CustomerRepositoryDb{client}
}
