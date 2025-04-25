package service

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	realdomain "github.com/kal997/banking/domain"
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

func Test_NewAccount_should_return_a_validation_error_response_when_the_request_is_not_validated(t *testing.T) {
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

func Test_NewAccount__should_return_an_error_from_the_server_side_if_the_new_account_cannot_be_created(t *testing.T) {
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

	account := realdomain.NewAccount(req.CustomerId, req.AccountType, req.Amount)

	accMockRepo.EXPECT().Save(account).Return(nil, errs.NewUnexpectedError("some db error"))

	//Act
	_, appErr := service.NewAccount(req)

	// Assert: assert on the retuern
	if appErr == nil {
		t.Error("Test failed while validating new account")
	}

}

func Test_NewAccount__should_return_new_account_response_when_a_new_account_is_saved_successfully(t *testing.T) {
	//AAA

	//Arrange
	teardown := setup(t)
	defer teardown()

	req := dto.NewAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      6000,
	}

	account := realdomain.NewAccount(req.CustomerId, req.AccountType, req.Amount)

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

func Test_MakeTransaction_should_return_a_validation_error_response_when_the_request_is_not_validated(t *testing.T) {
	//AAA
	//Arrange
	teardown := setup(t)
	defer teardown()

	req := dto.TransactionRequest{
		AccountId:       "100",
		Amount:          -5,
		TransactionType: "deposit",
		TransactionDate: time.Now().Format(dbTSLayout),
		CustomerId:      "2",
	}
	// this hasn't to be called (expect), because the call to SaveTransaction will not happen because the
	// MakeTrabsaction will return after Validate call (because of the validation error)

	//accMockRepo.EXPECT().SaveTransaction(req).Return(nil, nil)
	service.MakeTransaction(req)

	//Act
	err := req.Validate()

	//Assert
	if err == nil {
		t.Error("Test failed while validating a new transaction request")
	}

}
func Test_MakeTransaction_should_return_a_validation_error_response_when_the_request_balance_is_not_sufficient(t *testing.T) {

	//AAA
	//Arrage
	teardown := setup(t)
	defer teardown()

	req := dto.TransactionRequest{
		AccountId:       "100",
		Amount:          500,
		TransactionType: "withdrawal",
		TransactionDate: time.Now().Format(dbTSLayout),
		CustomerId:      "2",
	}

	foundAcc := realdomain.Account{
		AccountId:   req.AccountId,
		CustomerId:  req.CustomerId,
		OpeningDate: req.TransactionDate,
		AccountType: "saving",
		Amount:      200,
		Status:      "1",
	}

	accMockRepo.EXPECT().FindBy(req.AccountId).Return(&foundAcc, nil)
	//Act
	_, appErr := service.MakeTransaction(req)

	//Assert
	if appErr == nil {
		t.Error("Test failed while creating a new transaction request")
	}
}
func Test_MakeTransaction_should_return_a_new_transaction_response_when_the_request_is_saved_successfully(t *testing.T) {
	//AAA
	//Arrage
	teardown := setup(t)
	defer teardown()

	req := dto.TransactionRequest{
		AccountId:       "100",
		Amount:          500,
		TransactionType: "withdrawal",
		TransactionDate: time.Now().Format(dbTSLayout),
		CustomerId:      "2",
	}

	foundAcc := realdomain.Account{
		AccountId:   req.AccountId,
		CustomerId:  req.CustomerId,
		OpeningDate: req.TransactionDate,
		AccountType: "saving",
		Amount:      600,
		Status:      "1",
	}

	transactionRequest := realdomain.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: req.TransactionDate,
	}

	transactionResponse := realdomain.Transaction{
		TransactionId:   " ",
		AccountId:       req.AccountId,
		Amount:          foundAcc.Amount - req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: req.TransactionDate,
	}

	accMockRepo.EXPECT().FindBy(req.AccountId).Return(&foundAcc, nil)
	accMockRepo.EXPECT().SaveTransaction(transactionRequest).Return(&transactionResponse, nil)
	//Act
	expectedTransactionResponse, appErr := service.MakeTransaction(req)

	//Assert
	if appErr != nil {
		t.Error("Test failed while creating a new transaction request")
	}

	if expectedTransactionResponse.Amount != foundAcc.Amount-req.Amount {
		t.Error("Test failed while creating a new transaction request")
	}
}
func Test_MakeTransaction_should_return_a_repo_error_when_theres_a_internal_db_error(t *testing.T) {
	//AAA
	//Arrage
	teardown := setup(t)
	defer teardown()

	req := dto.TransactionRequest{
		AccountId:       "100",
		Amount:          500,
		TransactionType: "withdrawal",
		TransactionDate: time.Now().Format(dbTSLayout),
		CustomerId:      "2",
	}

	foundAcc := realdomain.Account{
		AccountId:   req.AccountId,
		CustomerId:  req.CustomerId,
		OpeningDate: req.TransactionDate,
		AccountType: "saving",
		Amount:      600,
		Status:      "1",
	}

	transactionRequest := realdomain.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: req.TransactionDate,
	}

	accMockRepo.EXPECT().FindBy(req.AccountId).Return(&foundAcc, nil)
	accMockRepo.EXPECT().SaveTransaction(transactionRequest).Return(nil, errs.NewUnexpectedError("Unexpected database error"))

	//Act
	_, appErr := service.MakeTransaction(req)

	if appErr.Code != http.StatusInternalServerError {
		t.Error("Test failed while creating a new transaction request")
	}

}

func Test_MakeTransaction_should_return_a_account_not_found_when_the_request_account_is_not_found_in_db(t *testing.T) {
	//AAA
	//Arrage
	teardown := setup(t)
	defer teardown()

	req := dto.TransactionRequest{
		AccountId:       "100",
		Amount:          500,
		TransactionType: "withdrawal",
		TransactionDate: time.Now().Format(dbTSLayout),
		CustomerId:      "2",
	}

	accMockRepo.EXPECT().FindBy(req.AccountId).Return(nil, errs.NewNotFoundError("account not found"))
	//Act
	_, appErr := service.MakeTransaction(req)

	//Assert
	if appErr.Code != http.StatusNotFound {
		t.Error("Test failed while creating a new transaction request")
	}

}
