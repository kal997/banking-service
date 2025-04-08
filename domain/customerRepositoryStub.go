package domain


// adapter == a type that satisfy the nterface
// stub adapter is simpler taht mock adapter
// the normal case will be database adapter

type CustomerRepositoryStub struct{
	customers []Customer
}



func (cs CustomerRepositoryStub)FindAll() ([]Customer, error){
	return cs.customers, nil
}



// helper function to get a usable instance of the concrete type
func NewCustomerRepositoryStub() CustomerRepositoryStub{
	customers := []Customer{
		{"10001", "khaled", "egypt", "100011", "04-06-1997", "1"},
		{"10021", "ahmed", "egypt", "100011", "01-01-2000", "1"},
	}
	return CustomerRepositoryStub{customers}
}



