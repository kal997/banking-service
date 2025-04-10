package service // user -> the domain interface

import (
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
)

// the interface that has to be implemented by the core logic to able user to interact with
// the interface that the user will expect from the business logic
type CustomerService interface {

	// one of the behaviours
	GetAllCustomers(string) ([]dto.CustomerResponse, *errs.AppError)

	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
}
