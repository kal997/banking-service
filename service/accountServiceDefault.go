package service

import (
	"time"

	"github.com/kal997/banking-lib/errs"
	"github.com/kal997/banking/domain"
	"github.com/kal997/banking/dto"
)

const dbTSLayout = "2006-01-02 15:04:05"

type DefaultAccountService struct {
	repo domain.AccountRepository
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {

	// validate the request
	if appErr := req.Validate(); appErr != nil {
		return nil, appErr
	}

	// in request: transform from dto to domain
	a := domain.Account{
		AccountId:   "",
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format("2006-01-02 15:04:05"),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	if newAccount, err := s.repo.Save(a); err != nil {
		return nil, err
	} else {

		// in response; transform from domain to dto
		return newAccount.ToNewAccountResponseDto(), nil

	}
}

func (s DefaultAccountService) MakeTransaction(dtoRequest dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError) {

	// Incoming request validation
	err := dtoRequest.Validate()
	if err != nil {
		return nil, err
	}

	// server side validation, if withdrawal, if the account has a sufficient balance
	if dtoRequest.IsTransactionTypeWithdrawal() {
		acc, err := s.repo.FindBy(dtoRequest.AccountId)
		if err != nil {
			return nil, err
		}

		// we need if the balance is sufficient
		if !acc.CanWithdraw(dtoRequest.Amount) {
			return nil, errs.NewValidationError("Insufficient balance in the account")
		}
	}

	// the transaction is valid by server side, if deposit ok, if withdrawal, he has a sufficient balance
	// transform from dto to domain
	domainTransactionRequest := domain.Transaction{
		AccountId:       dtoRequest.AccountId,
		Amount:          dtoRequest.Amount,
		TransactionType: dtoRequest.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	// DO BUSINESS LOGIC
	// request from domain
	domainTransactionResponse, appErr := s.repo.SaveTransaction(domainTransactionRequest)
	if appErr != nil {
		return nil, appErr
	}

	// transform response from domain to dto
	dtoResponse := domainTransactionResponse.ToDto()

	// END BUSINESS LOGIC

	// return dto
	return &dtoResponse, nil
}

func NewAccountService(repo domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{repo}
}
