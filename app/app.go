package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kal997/banking-lib/logger"
	"github.com/kal997/banking/domain"
	"github.com/kal997/banking/service"
)

func Start() {

	sanityCheck()

	router := mux.NewRouter()

	dbClient := getDbClient()

	customerRepoDb := domain.NewCustomerRepositoryDb(dbClient)
	customerService := service.NewCustomerService(customerRepoDb)

	accountRepoDb := domain.NewAccountRepositoryDb(dbClient)
	accountService := service.NewAccountService(accountRepoDb)

	authRepoDb := domain.NewRemoteAuthRepository(dbClient)
	authService := service.NewAuthService(authRepoDb, domain.GetRolePermissions())

	ch := CustomerHandler{service: customerService}
	ah := AccountHandler{service: accountService}
	authh := AuthMiddleware{service: authService}

	router.
		HandleFunc("/customers", ch.GetAllCustomers).
		Methods(http.MethodGet).
		Name("GetAllCustomers")
	router.
		HandleFunc("/customers/{customer_id:[0-9]+}", ch.GetCustomer).
		Methods(http.MethodGet).
		Name("GetCustomer")

	router.
		HandleFunc("/customers/{customer_id:[0-9]+}/account", ah.NewAccount).
		Methods(http.MethodPost).
		Name("NewAccount")
	router.
		HandleFunc("/customers/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", ah.MakeTransaction).
		Methods(http.MethodPost).
		Name("NewTransaction")

	router.
		HandleFunc("/auth/login", authh.Login).
		Methods(http.MethodPost).
		Name("Login")

	router.
		HandleFunc("/auth/verify", authh.Verify).
		Methods(http.MethodGet).
		Name("Verify")

	// adds authorization middleware to recieve the requests, dispatch the user request + verification info
	// to Verify endpoint, if the user is authorized, will pass the control to the user requested endpoint, else
	// return unauthorized
	router.Use(authh.authorizationHandler())

	// starting server
	logger.Info("starting server ..")
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	logger.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router).Error())

}

func sanityCheck() {
	if os.Getenv("SERVER_ADDRESS") == "" {
		log.Fatal("SERVER_ADDRESS is missing..")

	}

	if os.Getenv("SERVER_PORT") == "" {
		log.Fatal("SERVER_PORT is missing..")
	}

	if os.Getenv("DB_USER") == "" {
		log.Fatal("DB_USER is missing..")
	}

	if os.Getenv("DB_PASSWD") == "" {
		log.Fatal("DB_PASSWD is missing..")
	}
	if os.Getenv("DB_ADDR") == "" {
		log.Fatal("DB_ADDR is missing..")
	}
	if os.Getenv("DB_PORT") == "" {
		log.Fatal("DB_PORT is missing..")
	}
	if os.Getenv("DB_NAME") == "" {
		log.Fatal("DB_NAME is missing..")
	}

}

func getDbClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasWD := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasWD, dbAddr, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client

}
