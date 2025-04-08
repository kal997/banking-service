package app

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	
)

//same data transfer objects and it can handle different encodings/decodings

type Customer struct {
	Name    string `json:"name" xml:"name"`
	City    string `json:"city" xml:"city"`
	ZipCode string `json:"zipcode" xml:"zipcode"`
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers := []Customer{
		{Name: "khaled", City: "egypt", ZipCode: "1123"},
		{Name: "ahmed", City: "italy", ZipCode: "7899"},
	}

	if "application/xml" == r.Header.Get("Content-Type") {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(customers)

	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)

	}

}
