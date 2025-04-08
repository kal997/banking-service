package domain // bussiness side/domain -> repo interface

type Customer struct {
	Id          string
	Name        string
	City        string
	Zipcode     string
	DataofBirth string
	Status      string
}


// the interface == port

type CustomerRepository interface{

	// the method that has to be implemented in order to satisfy the interface
	FindAll() ([]Customer, error)
}



