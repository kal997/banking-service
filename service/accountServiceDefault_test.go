package service

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	realdomian "github.com/kal997/banking/domain"
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
	"github.com/kal997/banking/mocks/domain"
)

var accMockRepo *domain.MockAccountRepository
var service DefaultAccountService

func setup(t *testing.T) func() {

	ctrl := gomock.NewController(t)

	// create a new mock account service
	accMockRepo = domain.NewMockAccountRepository(ctrl)
	service = NewAccountService(accMockRepo)

	return func() {

		defer ctrl.Finish()

	}
}

func Test_should_return_a_validation_error_response_when_the_request_is_not_validated(t *testing.T) {
	// we want to fail validate, and see how NewAccount will behave
	//AAA
	//Arrange
	req := dto.NewAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      0,
	}
	service := NewAccountService(nil) // this test doesn't need to mock repo as it dependes on Validate failure

	// Act
	_, appErr := service.NewAccount(req)

	//Assert
	if appErr.Code != http.StatusUnprocessableEntity {
		t.Error("Failed while testing a new account validation")
	}
}

func Test_should_return_an_error_from_the_server_side_if_the_new_account_cannot_be_created(t *testing.T) {
	// should return error from server side if validate passed, but the db failed, so we need to mock
	// the server side (repo)

	//AAA

	//Arrange
	teardown := setup(t)
	defer teardown()

	// create a valid dto object
	req := dto.NewAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      6000,
	}

	account := realdomian.NewAccount(req.CustomerId, req.AccountType, req.Amount)

	accMockRepo.EXPECT().Save(account).Return(nil, errs.NewUnexpectedError("some db error"))

	//Act
	_, appErr := service.NewAccount(req)

	// Assert: assert on the retuern
	if appErr == nil {
		t.Error("Test failed while validating new account")
	}

}

func Test_should_return_new_account_response_when_a_new_account_is_saved_successfully(t *testing.T) {
	//AAA

	//Arrange
	teardown := setup(t)
	defer teardown()

	req := dto.NewAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      6000,
	}

	account := realdomian.NewAccount(req.CustomerId, req.AccountType, req.Amount)

	createdAccount := account
	createdAccount.AccountId = "500"
	accMockRepo.EXPECT().Save(account).Return(&createdAccount, nil)

	//Act
	response, appErr := service.NewAccount(req)

	if appErr != nil {
		t.Error("Test failed while creating new account, wrong app err recieved")
	}

	if response.AccountId != createdAccount.AccountId {
		t.Error("Test failed while creating new account, wrong account id")
	}

}
