package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kal997/banking/domain"
	"github.com/kal997/banking/service"
)

func Start() {

	// wiring
	// we choose :
		// stub repo (instead of DB)
		// Default Service (the main business logic)

	ch := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryStub())} 

	router := mux.NewRouter()
	router.HandleFunc("/customers", ch.GetAllCustomers).Methods(http.MethodGet) //  method matcher
	
	
	fmt.Println("starting server ..")
	log.Fatal(http.ListenAndServe(":8001", router))

}
