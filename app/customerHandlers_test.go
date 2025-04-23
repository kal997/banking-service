package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
	"github.com/kal997/banking/mocks/service"
)

// testing getallcustomers as a unit, we will inject the mock (the only external dependency)
// and we expect a slice of customers if the request is valid, and an error one if not

func Test_should_return_customer_with_status_code_200(t *testing.T) {
	//AAA

	//Arange (mock setup)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockCustomerService(ctrl)

	ch := CustomerHandler{mockService}

	router := mux.NewRouter()
	router.HandleFunc("/customers", ch.GetAllCustomers)

	// untill now, what we have done is creating a new router with one route, and the handler of this route will use the
	// mock service instead of the default one

	// what is needed, is to implement the behaviour we want for positive scenario test, i.e. returning a slice of customers from the mock
	dummyCustomers := []dto.CustomerResponse{
		{Id: "10001", Name: "khaled", City: "egypt", Zipcode: "100011", DataofBirth: "04-06-1997", Status: "1"},
		{Id: "10021", Name: "ahmed", City: "egypt", Zipcode: "100011", DataofBirth: "01-01-2000", Status: "1"},
	}
	mockService.EXPECT().GetAllCustomers("").Return(dummyCustomers, nil)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockCustomerService(ctrl)

	ch := CustomerHandler{mockService}

	router := mux.NewRouter()
	router.HandleFunc("/customers", ch.GetAllCustomers)

	// untill now, what we have done is creating a new router with one route, and the handler of this route will use the
	// mock service instead of the default one

	// what is needed, is to implement the behaviour we want for positive scenario test, i.e. returning a slice of customers from the mock

	mockService.EXPECT().GetAllCustomers("").Return(nil, errs.NewUnexpectedError("some db errors"))
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act (create a http request)
	// since we need a response writer for the test, not a normal one
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request) // test recorder will bring the response from the handler

	//Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code negative scenario")
	}

}
