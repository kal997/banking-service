package domain // bussiness side/domain -> repo interface

import (
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
)

type Customer struct {
	Id          string `db:"customer_id"`
	Name        string
	City        string
	Zipcode     string
	DataofBirth string `db:"date_of_birth"`
	Status      string
}

func (c Customer) statusAsText() string{
	statusAsText := "Active"
	if c.Status == "0"{
		statusAsText = "Inactive"
	} 
	return statusAsText
}

func (c Customer) ToDto() dto.CustomerResponse{

	
	return dto.CustomerResponse{
		Id: c.Id,
		Name: c.Name,
		City: c.City,
		Zipcode: c.Zipcode,
		DataofBirth: c.DataofBirth,
		Status: c.statusAsText(),
	}
}

// the interface == port

type CustomerRepository interface {

	// the method that has to be implemented in order to satisfy the interface
	FindAll(string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}
