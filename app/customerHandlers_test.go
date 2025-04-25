package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/kal997/banking-lib/errs"
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/mocks/service"
)

// testing getallcustomers as a unit, we will inject the mock (the only external dependency)
// and we expect a slice of customers if the request is valid, and an error one if not

var mockService *service.MockCustomerService
var ch CustomerHandler
var router *mux.Router

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockCustomerService(ctrl)

	ch = CustomerHandler{mockService}
	router = mux.NewRouter()
	return func() {
		router = nil
		defer ctrl.Finish()
	}

}

func Test_should_return_customer_with_status_code_200(t *testing.T) {
	//AAA

	//Arange (mock setup)
	teardown := setup(t)
	defer teardown()

	// what is needed, is to implement the behaviour we want for positive scenario test, i.e. returning a slice of customers from the mock
	dummyCustomers := []dto.CustomerResponse{
		{Id: "10001", Name: "khaled", City: "egypt", Zipcode: "100011", DataofBirth: "04-06-1997", Status: "1"},
		{Id: "10021", Name: "ahmed", City: "egypt", Zipcode: "100011", DataofBirth: "01-01-2000", Status: "1"},
	}
	mockService.EXPECT().GetAllCustomers("").Return(dummyCustomers, nil)
	router.HandleFunc("/customers", ch.GetAllCustomers)

	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act (create a http request)
	// since we need a response writer for the test, not a normal one
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request) // test recorder will bring the response from the handler

	//Assert
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code positive scenario")
	}

}

func Test_should_return_status_code_500_with_error_message(t *testing.T) {
	//AAA

	//Arange (mock setup)
	teardown := setup(t)
	defer teardown()
	mockService.EXPECT().GetAllCustomers("").Return(nil, errs.NewUnexpectedError("some db errors"))
	router.HandleFunc("/customers", ch.GetAllCustomers)

	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act (create a http request)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request) // test recorder will bring the response from the handler

	//Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code negative scenario")
	}

}

func Test_should_return_one_customer_with_status_code_200(t *testing.T) {
	//AAA

	//Arange (mock setup)
	teardown := setup(t)
	defer teardown()

	dummyCustomer := dto.CustomerResponse{
		Id:          "10001",
		Name:        "khaled",
		City:        "egypt",
		Zipcode:     "100011",
		DataofBirth: "04-06-1997",
		Status:      "1",
	}

	mockService.EXPECT().GetCustomer(dummyCustomer.Id).Return(&dummyCustomer, nil)
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.GetCustomer)

	request, _ := http.NewRequest(http.MethodGet, "/customers/10001", nil)

	// Act (create a http request)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request) // test recorder will bring the response from the handler

	//Assert
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code positive scenario")
	}

}

func Test_should_return_error_customer_not_found(t *testing.T) {
	//AAA

	//Arange (mock setup)
	teardown := setup(t)
	defer teardown()

	mockService.EXPECT().GetCustomer("500").Return(nil, errs.NewNotFoundError("customer not found"))
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.GetCustomer)

	request, _ := http.NewRequest(http.MethodGet, "/customers/500", nil)

	// Act (create a http request)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request) // test recorder will bring the response from the handler

	//Assert
	if recorder.Code == http.StatusOK {
		t.Error("Failed while testing the status code negative scenario")
	}

}
