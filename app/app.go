package app


import (
	
	"fmt"
	"log"
	"net/http"
)


func Start(){
	
	mux := http.NewServeMux()
	mux.HandleFunc("/customers", getAllCustomers)

	fmt.Println("starting server ..")
	log.Fatal(http.ListenAndServe(":8001", mux))

}