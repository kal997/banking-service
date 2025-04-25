package service // user -> the domain interface

import (
	"github.com/kal997/banking-lib/errs"
	"github.com/kal997/banking/dto"
)

// the interface that has to be implemented by the core logic to able user to interact with
// the interface that the user will expect from the business logic
// go gen tool tag, to invoke go gen tool
//
//go:generate mockgen -destination=../mocks/service/mockCustomerService.go -package=service github.com/kal997/banking/service CustomerService
type CustomerService interface {

	// one of the behaviours
	GetAllCustomers(string) ([]dto.CustomerResponse, *errs.AppError)

	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
}
