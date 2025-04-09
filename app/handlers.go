package app

import (
	"encoding/json"
	"encoding/xml"
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

	customers, _ := ch.service.GetAllCustomers()

	if r.Header.Get("Content-Type") == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(customers)

	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)

	}

}

func (ch *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["customer_id"]
	customer, appErr := ch.service.GetCustomer(id)
	if appErr != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(appErr.Code)
		json.NewEncoder(w).Encode(appErr.AsMessage())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(customer)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
