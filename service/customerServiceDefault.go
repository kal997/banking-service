package service

import (
	"github.com/kal997/banking/domain"
	"github.com/kal997/banking/errs"
)

// the implementaion of the CustomerService interface
type DefaultCustomerService struct {

	// it also depends on the domain port, not the concrete implementation, mock or db
	repo domain.CustomerRepository
}

func (dcs DefaultCustomerService) GetAllCustomers()([]domain.Customer, error) {
	return dcs.repo.FindAll()
}

// again, helper function to get us a instance of ready to use Service, that is already has it's dependency satisfied
// kunna momkin n3ml instance we i-inject by accessing the memebr , bas kda a7sn
func NewCustomerService(repo domain.CustomerRepository) DefaultCustomerService{
	return DefaultCustomerService{repo}
}




func (dcs DefaultCustomerService) GetCustomer(id string) (*domain.Customer, *errs.AppError){
	return dcs.repo.ById(id)
}
