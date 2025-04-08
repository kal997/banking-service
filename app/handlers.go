package app

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

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
