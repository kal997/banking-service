package app


import (
	
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)


func Start(){
	
	router := mux.NewRouter()
	router.HandleFunc("/customers", getAllCustomers)

	fmt.Println("starting server ..")
	log.Fatal(http.ListenAndServe(":8001", router))

}