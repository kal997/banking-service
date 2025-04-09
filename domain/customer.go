package domain // bussiness side/domain -> repo interface

import "github.com/kal997/banking/errs"

type Customer struct {
	Id          string
	Name        string
	City        string
	Zipcode     string
	DataofBirth string
	Status      string
}

// the interface == port

type CustomerRepository interface {

	// the method that has to be implemented in order to satisfy the interface
	FindAll(string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}
