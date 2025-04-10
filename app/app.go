package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kal997/banking/domain"
	"github.com/kal997/banking/logger"
	"github.com/kal997/banking/service"
)

func Start() {

	// wiring
	// we choose :
	// stub repo (instead of DB)
	// Default Service (the main business logic)

	//ch := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryStub())}
	ch := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryDb())}

	router := mux.NewRouter()
	router.HandleFunc("/customers", ch.GetAllCustomers).Methods(http.MethodGet)                  //  method matcher
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.GetCustomer).Methods(http.MethodGet) //  method matcher

	logger.Info("starting server ..")
	logger.Fatal(http.ListenAndServe(":8001", router).Error())

}
