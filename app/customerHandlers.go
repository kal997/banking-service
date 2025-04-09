package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kal997/banking/service"
)

// a handler must have dependecy on the service (the interface)
// we will create a concrete impementation
type CustomerHandler struct {
	service service.CustomerService // depends on service port
}

// we pass CustomerHandler as a reciver to the GetAllCustomers http handler
// and inside it we use the service interface to get the customers
func (ch *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {

	status := r.URL.Query().Get("status")
	customers, appError := ch.service.GetAllCustomers(status)
	if appError != nil {
		writeResponse(w, appError.Code, appError.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, customers)
	}
}

func (ch *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["customer_id"]
	customer, appErr := ch.service.GetCustomer(id)
	if appErr != nil {
		writeResponse(w, appErr.Code, appErr.AsMessage())
	} else {

		writeResponse(w, http.StatusOK, customer)
	}

}


func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Encode ", data, "failed with err ", err)
		panic(err)
	}
}
